# presto_assignment

This service manages Time-of-Use (TOU) pricing schedules for individual EV chargers.
It supports schedule creation, retrieval, timezone-aware evaluation, and atomic bulk updates.

📌 Features
1) Charger-specific TOU pricing
2) Versioned pricing schedules (effective_from, effective_to)
3) Daily schedule retrieval
4) Timezone-aware datetime-based pricing lookup
5) Atomic bulk schedule updates
6) Transaction-safe operations
7) Interval semantics using [start_time, end_time)

🏗 Architecture Overview
1) The service follows a layered architecture:
 Handler → Service → Repository → Database
2) Handler Layer: HTTP request/response handling
3) Service Layer: Business logic, validation, normalization
4) Repository Layer: Database interactions
5) Database: PostgreSQL with normalized relational schema
6) All schedule modifications are executed inside database transactions to ensure consistency.

🛠 Running the Service
1️⃣ Set Environment Variable
Create a .env file:
paste in your configured DATABASE_URL
e.g. DATABASE_URL=postgres://user-name:password@localhost:5432/db_name?sslmode=disable

2️⃣ Run Migrations
migrate -path migrations -database "postgres://user-name:password@localhost:5432/db_name?sslmode=disable" up

3️⃣ Start Server
go run main.go

Sample-CURLS for testing
1) CREATE ENDPOINT: (MAKE sure a charger with a charging_station exists in database)
curl --location 'http://localhost:8080/chargers/22222222-2222-2222-2222-222222222222/pricing-schedules' \
--header 'Content-Type: application/json' \
--data '{
    "effective_from": "2026-03-10",
    "periods": [
        {
            "start_time": "00:00",
            "end_time": "12:00",
            "price_per_kwh": 0.15
        },
        {
            "start_time": "12:00",
            "end_time": "23:59",
            "price_per_kwh": 0.25
        }
    ]
}'

2) GET ENDPOINT: 
curl --location 'http://localhost:8080/chargers/22222222-2222-2222-2222-222222222222/pricing?date=2026-03-05'

3) Bulk Update ENDPOINT:
curl --location 'http://localhost:8080/pricing-schedules/bulk' \
--header 'Content-Type: application/json' \
--data '{
  "charger_ids": [
    "68aa0000-abc6-48d7-8c76-ab6d2c3d2901",
    "a2755236-7994-4df3-ab33-3e2ea4e60352",
    "99999999-9999-9999-9999-999999999999"
  ],
  "effective_from": "2026-03-12",
  "periods": [
    {
      "start_time": "00:00",
      "end_time": "06:00",
      "price_per_kwh": 0.15
    },
    {
      "start_time": "06:00",
      "end_time": "23:59",
      "price_per_kwh": 0.25
    }
  ]
}'