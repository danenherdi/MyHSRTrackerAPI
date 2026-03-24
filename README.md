# MyHSRTrackerAPI - In Progress (Not Tested Yet)

This is your local backend for Honkai Star Rail pull tracking. It connects to the HoYoverse API using an `AuthURL` to privately download your gacha history into a local SQLite database.

## How to Run the Server

1. Open your terminal in this directory.
2. Run the command:
   ```bash
   go run main.go
   ```
3. The server will start locally on `http://localhost:8080`.

## How to Get your AuthURL

1. Open Honkai Star Rail on PC.
2. Open the **Warp** page -> Click **View Details** -> Click **Records**.
3. Open Windows PowerShell and run this private, local script (this privately grabs the hidden URL from your local game logs without downloading anything):
   ```powershell
   To be added...
   ```

---

## API Endpoints 

### 1. Import Warp History
**POST** `/api/warp/import`

**Request Body:**
```json
{
  "auth_url": "...<YOUR_AUTH_KEY_HERE>"
}
```

**What it does:**
It will iterate through all 4 banner types, fetching up to your oldest pull and inserting it securely into the local database. Rate limits are automatically managed.

### 2. View Warp History 
**GET** `/api/warp/list`

**Query Parameters (Optional):**
- `page` (default `1`)
- `size` (default `20`)
- `gacha_type` (e.g. `11`, `1`, `2`, `12`)
- `uid` (filter by specific account id)

**Example:**
`http://localhost:8080/api/warp/list?gacha_type=12&page=1`

### 3. View Pity Stats
**GET** `/api/warp/stats?uid={YOUR_UID}`

**What it does:**
Calculates the exact total pulls and current pity (pulls since your last 5-Star) entirely locally without relying on any third-party websites or external trackers!
