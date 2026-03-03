EV Charger Time-of-Use (TOU) Pricing API

This document defines the REST API contract for managing and retrieving Time-of-Use pricing schedules for EV chargers.

Base URL
http://localhost:8080
=========================================================
1️⃣ Create Pricing Schedule
Endpoint:
POST /chargers/{chargerID}/pricing-schedules

Description:
1) Creates a new TOU pricing schedule for a specific charger.
2) If a previously active schedule exists, its effective_to date will be updated to the day before the new schedule’s effective_from, ensuring no overlap and maintaining temporal versioning.
3) Schedules are versioned using effective_from and effective_to.

Path Parameters:
Parameter	Type	Required    Description
chargerID	UUID	Yes	        Unique identifier of the charger

Request Body:
{
  "effective_from": "2026-03-01",
  "periods": [
    {
      "start_time": "00:00", //Accepts time in formats such as: 
      "end_time": "06:00",
      "price_per_kwh": 0.15
    },
    {
      "start_time": "06:00",
      "end_time": "18:00",
      "price_per_kwh": 0.25
    },
    {
      "start_time": "18:00",
      "end_time": "23:59",
      "price_per_kwh": 0.20
    }
  ]
}

Field Definitions:
Field	        Type	Format	    Description
effective_from	String	YYYY-MM-DD	Date from which schedule becomes active
periods	        Array	-	        List of pricing periods
start_time	    String	HH:MM	    Inclusive start time
end_time	    String	HH:MM	    Exclusive end time
price_per_kwh	Float	≥ 0	        Cost per kWh

start_time and end_time must represent a valid time of day.
The API accepts common time formats (e.g., HH:MM, HH:MM:SS, 3 PM, 3:04 PM).
All times are normalized internally to HH:MM (24-hour format).

Business Rules:
1) Periods must not overlap.
2) start_time must be less than end_time.
3) Time intervals follow [start_time, end_time) rule.
4) Only one active schedule per charger at a time.
5) If a schedule already exists for the same effective_from, it will be rejected (or replaced depending on logic) //CAN BE IMPLEMENTED BASED ON THE REQUIREMENT.

Success Response:
{
  "status": "success",
  "message": "Pricing schedule created successfully"
}
HTTP Status: 201 Created

Error Response:
{
  "status": "error",
  "message": "charger not found"
}
HTTP Status: 400 Bad Request or 404 Not Found

=========================================================
=========================================================
2️⃣ Get Daily Pricing Schedule
Endpoint:
GET /chargers/{chargerID}/pricing?date=YYYY-MM-DD

Description:
1) Retrieves the applicable pricing schedule for a specific charger and date.
2) Returns all time periods valid for that day.

Path Parameters:
Parameter	Type	Required
chargerID	UUID	Yes


Query Parameters:
Parameter	Type	Required	Format
date	    String	Yes	        YYYY-MM-DD

Success Response:
{
  "status": "success",
  "message": "pricing retrieved successfully",
  "data": {
    "charger_id": "22222222-2222-2222-2222-222222222222",
    "effective_from": "2026-03-01",
    "periods": [
      {
        "start_time": "00:00",
        "end_time": "06:00",
        "price_per_kwh": 0.15
      },
      {
        "start_time": "06:00",
        "end_time": "18:00",
        "price_per_kwh": 0.25
      }
    ]
  }
}
HTTP Status: 200 OK

Error Response:
{
  "status": "error",
  "message": "date must be in YYYY-MM-DD format"
}

Interval Semantics: 
All time intervals follow: [start_time, end_time)
Meaning:
1) Start time is inclusive
2) End time is exclusive

Example:
00:00 - 18:00
18:00 - 23:59
A transaction at 18:00 belongs to the second interval.

Schedule Versioning Logic:
Schedules are selected using:
1) effective_from <= date
2) AND (effective_to IS NULL OR effective_to >= date)

Only one active schedule per charger is allowed at a time.
Assumptions:
1) No weekday/weekend variation.
2) No timezone conversion (local station time assumed).
3) Time precision is minute-level.
4) Cross-midnight intervals are currently not supported.

Time Zone Interpretation:
1) The date parameter is interpreted in the charger’s local timezone.
2) The system retrieves pricing based on the charger’s local calendar date.
3) Clients must provide the calendar date corresponding to the charger’s physical location.

=========================================================
=========================================================
📄 Bulk Pricing Schedule API Contract
Bulk Create Pricing Schedule
Endpoint:
POST /pricing-schedules/bulk

Description:
1) Applies the same TOU pricing schedule to multiple chargers in a single request.
2) The operation is atomic i.e. If any charger update fails, the entire transaction is rolled back.
3) The response will include the charger_id where the first failure occurred.

Request Body:
{
  "charger_ids": [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222"
  ],
  "effective_from": "2026-03-10",
  "periods": [
    {
      "start_time": "00:00",
      "end_time": "06:00",
      "price_per_kwh": "0.15"
    },
    {
      "start_time": "06:00",
      "end_time": "18:00",
      "price_per_kwh": "0.25"
    }
  ]
}
Field Definitions
Field	        Type	    Required	Description
charger_ids	    Array<UUID>	Yes	        List of chargers to apply schedule to
effective_from	String	    Yes	        Date from which schedule becomes active (YYYY-MM-DD)
periods	        Array	    Yes	        List of pricing periods
start_time	    String	    Yes	        Start time (inclusive)
end_time	    String	    Yes	        End time (exclusive)
price_per_kwh	Float	    Yes	        Cost per kWh

Business Rules/Assumptions:
1) All charger_ids must exist.
2) Periods must not overlap.
3) Time intervals follow [start_time, end_time) semantics.
4) Only one active schedule per charger at a time.
5) Entire request is processed within a single transaction.

Success Response:
{
  "status": "success",
  "message": "Bulk pricing schedule created successfully",
}
HTTP Status: 201 Created

Failure Response (Atomic Rollback):
If any charger fails:
{
  "status": "error",
  "message": "Bulk update failed",
  "data": {
    "failed_charger_id": "22222222-2222-2222-2222-222222222222"
  }
}
HTTP Status: 400 Bad Request or 500 Internal Server Error
All changes are rolled back.

Time Zone Handling:
1) effective_from is interpreted in the charger’s local timezone.
2) Schedule evaluation follows charger-local calendar semantics.

=========================================================
Future Improvements (Optional):
1) Instead of using SQL TIME, store start_time and end_time as integer minute offsets from midnight (e.g., 00:00 → 0, 24:00 → 1440).
This simplifies comparisons, avoids edge cases at 23:59, and provides cleaner interval evaluation using [start, end) semantics.
2) We have used float64 for price_per_kwh but using float64 may cause rounding errors in financial calculations.
In production, pricing should use fixed-point precision (e.g., NUMERIC in DB and decimal handling in Go or string representation in API) to ensure accurate billing.
3) Expose a new API to support exact timestamp-based pricing lookup (e.g., ?datetime=RFC3339).
This enables precise billing and audit capabilities for charging sessions.

