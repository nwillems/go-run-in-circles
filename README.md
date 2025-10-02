# go-run-in-circles
A simple system to keep track of runners, doing laps

## Lap Tracking System: Complete Architecture
### 1. System Components
Frontend (Basic HTML/JS)
Single-Page Application:
Runner Registration:
Officials input bib number, name, and upload a photo (embedded in DB).


Lap Recording & Dashboard (Combined):
Officials record laps via bib number input (manual or barcode scanner).
Table displays all runners, lap counts, and metrics.
Auto-refreshes every 5 seconds.


Display View:
Full-screen, read-only table of all runners and their progress for the extra screen.

Backend (Golang)
API Endpoints:
* POST /runners: Register a runner (bib number, name, photo as binary data).
* POST /laps: Record a lap (bib number, timestamp).
* GET /dashboard: Return all runners, laps, and calculated metrics.


*Database Selection:* SQLite database file selected at startup (no reset functionality).

Database (SQLite)
Schema:
runners: bib_number (PK, TEXT), name (TEXT), photo (BLOB).
laps: id (INTEGER PK), bib_number (TEXT, FK), lap_number (INTEGER), timestamp (DATETIME).
events: id (INTEGER PK), name (TEXT), date (DATE), lap_distance (INTEGER).


### 2. Data Flow
Startup:
Official selects SQLite database file.
System loads event and runner data.


Registration:
Official adds runners; photos are stored as BLOBs in the database.

Lap Recording:
Official enters/scans bib number → system records lap with timestamp.


Dashboard:

/dashboard endpoint returns JSON:
```
[
  {
    "bib_number": "A001",
    "name": "Nicolai Willems",
    "total_laps": 5,
    "last_lap_time": "01:23",
    "average_lap_time": "01:25",
    "fastest_lap": "01:20",
    "slowest_lap": "01:30",
    "avg_pace": "05:45/km",
    "photo": "base64encoded..."
  },
  ...
]
```


Display:
Extra screen shows all runners, sorted by bib number, with metrics updated on refresh.

### 3. Example Workflow
Official Actions:
Register runners, record laps, monitor progress.

Public Display:
Always shows all runners, even those with zero laps.




### 4. Technical Implementation Notes

Photo Handling:
Frontend encodes uploaded photos as base64, sends to backend, stored as BLOB in SQLite.


Auto-Refresh:
Frontend uses setInterval to poll /dashboard every 5 seconds.


Portability:
Single Go binary + SQLite file; no external dependencies.




Ready to build!
Would you like a basic code skeleton for the Go backend or HTML frontend to get started?
