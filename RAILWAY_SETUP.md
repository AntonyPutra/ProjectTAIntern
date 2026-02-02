# Railway Deployment - Step by Step Fix

## Masalah

Railway gagal detect karena di root ada 2 project:
- `mini-oms-backend/` (Golang)
- `mini-oms-frontend/` (React)

Railway tidak tahu mau deploy yang mana.

## Solusi: Deploy Backend dan Frontend Terpisah

### Step 1: Delete Service Yang Gagal

1. Di Railway dashboard, klik service **ProjectTAIntern** yang error
2. Klik **Settings** (tab paling kanan)
3. Scroll ke bawah → **Danger** section
4. Klik **Delete Service**
5. Confirm

### Step 2: Deploy Backend (Golang)

1. **Kembali ke Project dashboard**
2. Klik **+ New** → **GitHub Repo**
3. **Select**: ProjectTAIntern
4. **Configure Service**:
   - Name: `mini-oms-backend`
   - **PENTING**: Klik **Settings** → **Service**
   - **Root Directory**: `mini-oms-backend`
   - **Build Command**: (kosongkan, Railway auto-detect Go)
   - **Start Command**: (kosongkan, Railway pakai `./main`)

5. **Add Environment Variables** (klik **Variables** tab):
   ```
   PORT=8080
   ENV=production
   JWT_SECRET=railway-production-secret-key-12345
   JWT_EXPIRY_HOURS=24
   ```
   
   Database variables akan diisi setelah create PostgreSQL.

6. **Deploy** → Railway mulai build

### Step 3: Create PostgreSQL Database

1. Klik **+ New** → **Database** → **Add PostgreSQL**
2. Name: `mini-oms-db`
3. Region: **us-west1** (sama dengan backend)
4. **Create**

Railway auto-generate credentials dan expose sebagai environment variables:
- `PGHOST`
- `PGPORT`
- `PGUSER`
- `PGPASSWORD`
- `PGDATABASE`

### Step 4: Connect Backend ke Database

1. Klik service **mini-oms-backend**
2. Klik **Variables** tab
3. **Add References** (klik + Reference):
   - `DB_HOST` = `${{Postgres.PGHOST}}`
   - `DB_PORT` = `${{Postgres.PGPORT}}`
   - `DB_USER` = `${{Postgres.PGUSER}}`
   - `DB_PASS` = `${{Postgres.PGPASSWORD}}`
   - `DB_NAME` = `${{Postgres.PGDATABASE}}`

4. Railway auto-redeploy backend

### Step 5: Get Backend URL

1. Klik service **mini-oms-backend**
2. Klik **Settings** tab
3. Scroll ke **Networking** section
4. Klik **Generate Domain**
5. Dapat URL: `mini-oms-backend.up.railway.app`
6. **COPY URL INI** untuk frontend

### Step 6: Deploy Frontend (React + Vite)

1. **Kembali ke Project dashboard**
2. Klik **+ New** → **GitHub Repo**
3. **Select**: ProjectTAIntern (lagi, tapi beda root directory)
4. **Configure Service**:
   - Name: `mini-oms-frontend`
   - **Settings** → **Service**
   - **Root Directory**: `mini-oms-frontend`
   - **Build Command**: `npm install && npm run build`
   - **Start Command**: `npm run preview -- --host 0.0.0.0 --port $PORT`

5. **Add Environment Variables**:
   ```
   VITE_API_BASE_URL=https://mini-oms-backend.up.railway.app
   ```
   
   Ganti dengan URL backend dari Step 5!

6. **Deploy** → Railway build frontend

### Step 7: Get Frontend URL

1. Klik service **mini-oms-frontend**
2. **Settings** → **Networking**
3. **Generate Domain**
4. Dapat URL: `mini-oms-frontend.up.railway.app`

### Step 8: Update CORS di Backend (Optional)

Jika CORS masih bermasalah, update backend:

1. Di local, edit `mini-oms-backend/cmd/api/main.go`
2. Update CORS AllowOrigins:
   ```go
   AllowOrigins: []string{
       "https://mini-oms-frontend.up.railway.app",
       "*", // atau specific domain saja
   },
   ```
3. Commit & push:
   ```bash
   git add .
   git commit -m "fix: update CORS for Railway domain"
   git push origin main
   ```
4. Railway auto-redeploy

## Verify Deployment

### Test Backend

```bash
curl https://mini-oms-backend.up.railway.app/api/products
```

Expected:
```json
{
  "success": true,
  "message": "Products retrieved successfully",
  "data": []
}
```

### Test Frontend

Buka browser:
```
https://mini-oms-frontend.up.railway.app
```

Seharusnya:
- Halaman login/register muncul
- Bisa register user baru
- Bisa login
- Bisa browse products
- Bisa create order & payment

## Troubleshooting

### Backend Error: "Failed to connect to database"

**Fix:**
1. Pastikan PostgreSQL service running
2. Cek environment variables `DB_HOST`, `DB_PORT`, dll sudah di-set
3. Cek di **Logs** ada error message apa

### Frontend Error: "Network Error" saat API call

**Fix:**
1. Cek `VITE_API_BASE_URL` sudah benar (https, bukan http)
2. Cek backend URL bisa diakses
3. Cek CORS di backend

### Build Failed: "go: cannot find main module"

**Fix:**
1. Pastikan **Root Directory** = `mini-oms-backend`
2. Pastikan ada file `go.mod` di `mini-oms-backend/`

### Frontend Build Failed: "vite: not found"

**Fix:**
1. Pastikan **Root Directory** = `mini-oms-frontend`
2. Build command: `npm install && npm run build`
3. Start command: `npm run preview -- --host 0.0.0.0 --port $PORT`

## Alternative: Nixpacks Config (Advanced)

Jika mau deploy dari root directory (tidak recommended), buat file:

**nixpacks.toml** di root:
```toml
[phases.build]
nixPkgs = ["go_1_21", "nodejs_20"]

[phases.setup]
cmds = ["cd mini-oms-backend && go mod download"]

[phases.build]
cmds = ["cd mini-oms-backend && go build -o main cmd/api/main.go"]

[start]
cmd = "cd mini-oms-backend && ./main"
```

Tapi lebih mudah pakai approach terpisah (Step 1-7).

## Summary

Railway config:
- **Service 1**: Backend (root: `mini-oms-backend/`)
- **Service 2**: Frontend (root: `mini-oms-frontend/`)
- **Database**: PostgreSQL (linked ke backend)

Total usage free tier:
- Backend: ~$2/month
- Frontend: ~$1/month
- Database: $1/month
- **Total**: ~$4/month (masih di bawah $5 limit)

Deployment berhasil!
