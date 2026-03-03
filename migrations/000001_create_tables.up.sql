-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create ENUM type
DO $$ BEGIN
    CREATE TYPE status_enum AS ENUM ('ACTIVE', 'INACTIVE');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Charging Stations
CREATE TABLE charging_stations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    timezone VARCHAR(64) NOT NULL,
    status status_enum DEFAULT 'ACTIVE',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Chargers
CREATE TABLE chargers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    station_id UUID NOT NULL REFERENCES charging_stations(id) ON DELETE CASCADE,
    status status_enum DEFAULT 'ACTIVE',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- TOU Pricing Schedules
CREATE TABLE tou_pricing_schedules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    charger_id UUID NOT NULL REFERENCES chargers(id) ON DELETE CASCADE,
    effective_from DATE NOT NULL,
    effective_to DATE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- TOU Pricing Periods
CREATE TABLE tou_pricing_periods (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    schedule_id UUID NOT NULL REFERENCES tou_pricing_schedules(id) ON DELETE CASCADE,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    price_per_kwh NUMERIC(10,4) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
-- CREATE INDEX idx_schedule_lookup
-- ON tou_pricing_schedules (charger_id, effective_from, effective_to);

-- CREATE INDEX idx_period_lookup
-- ON tou_pricing_periods (schedule_id, start_time);