# Deployment Guide - Render.com

## Masalah Payment di Render & Solusinya

### Masalah Umum:
1. **CORS Error** - Frontend tidak bisa hit backend API
2. **401 Unauthorized** - JWT token tidak valid
3. **500 Internal Server Error** - Database connection atau environment variables salah

---

## Step-by-Step Deploy ke Render.com

### 1. Deploy Backend (Golang + PostgreSQL)

#### A. Create PostgreSQL Database
1. Login ke **render.com**
2. Klik **New +** → **PostgreSQL**
3. Isi:
   - Name: `mini-oms-db`
   - Region: **Singapore** (paling dekat)
   - Instance Type: **Free**
4. Klik **Create Database**
5. **PENTING**: Copy semua kredensial:
   - Internal Database URL
   - PSQL Command
   - Database: `mini_oms_db`
   - User: `mini_oms_db_user`
   - Password: `xxxxx`
   - Host: `xxx.postgres.render.com`
   - Port: `5432`

#### B. Deploy Backend Service
1. Klik **New +** → **Web Service**
2. Connect ke GitLab:
   - Authorize Render to access GitLab
   - Pilih repository: **projecttaintern**
3. Konfigurasi:
   - **Name**: `mini-oms-backend`
   - **Region**: Singapore
   - **Branch**: main
   - **Root Directory**: `mini-oms-backend`
   - **Runtime**: Go
   - **Build Command**: 
     ```bash
     go build -o main cmd/api/main.go
     ```
   - **Start Command**:
     ```bash
     ./main
     ```
   - **Instance Type**: Free

4. **Environment Variables** (PENTING!):
   Klik **Advanced** → **Add Environment Variable**
   
   ```
   PORT=8080
   ENV=production
   
   # Database (dari PostgreSQL yang dibuat tadi)
   DB_HOST=<host-dari-render-postgres>.postgres.render.com
   DB_PORT=5432
   DB_USER=mini_oms_db_user
   DB_NAME=mini_oms_db
   DB_PASS=<password-dari-render-postgres>
   
   # JWT Secret (bikin random string panjang)
   JWT_SECRET=production-super-secret-key-random-12345-abcdef
   JWT_EXPIRY_HOURS=24
   ```
   
5. Klik **Create Web Service**
6. Tunggu deploy (5-10 menit)
7. Setelah selesai, dapat URL: `https://mini-oms-backend.onrender.com`

#### C. Test Backend
Buka browser atau Postman:
```
https://mini-oms-backend.onrender.com/api/products
```

Jika return `{"success":true, ...}` berarti backend sudah jalan!

---

### 2. Deploy Frontend (React + Vite)

#### A. Create Static Site
1. Klik **New +** → **Static Site**
2. Connect GitLab → **projecttaintern**
3. Konfigurasi:
   - **Name**: `mini-oms-frontend`
   - **Branch**: main
   - **Root Directory**: `mini-oms-frontend`
   - **Build Command**:
     ```bash
     npm install && npm run build
     ```
   - **Publish Directory**: `dist`

4. **Environment Variables**:
   ```
   VITE_API_BASE_URL=https://mini-oms-backend.onrender.com
   ```
   
   **PENTING**: Ganti dengan URL backend kamu yang sebenarnya!

5. Klik **Create Static Site**
6. Tunggu build (5 menit)
7. Dapat URL: `https://mini-oms-frontend.onrender.com`

---

### 3. Update CORS di Backend (SUDAH DILAKUKAN)

File `cmd/api/main.go` sudah diupdate dengan:
```go
e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins: []string{"*"},
    AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS},
    AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
}))
```

Setelah ini, **push ke Git** lagi:
```bash
cd c:\Whaleee\Intern\TA
git add .
git commit -m "fix: update CORS config for production deployment"
git push
```

Render akan auto-redeploy backend dalam 2-3 menit.

---

## Troubleshooting Payment Issues

### Problem 1: Payment endpoint return 401 Unauthorized

**Penyebab:**
- JWT token tidak dikirim
- JWT Secret berbeda antara local dan production

**Solusi:**
1. Cek JWT_SECRET di environment variables Render sama dengan yang kamu pakai local
2. Cek token tersimpan di localStorage:
   ```javascript
   // Di browser console
   console.log(localStorage.getItem('access_token'));
   ```
3. Pastikan login ulang di production untuk dapat token baru

### Problem 2: CORS Error - "Access-Control-Allow-Origin"

**Penyebab:**
- Backend tidak allow request dari frontend domain

**Solusi:**
SUDAH FIXED dengan update CORS config di main.go. Pastikan sudah push ke Git dan redeploy.

### Problem 3: 500 Internal Server Error saat Create Payment

**Penyebab:**
- Database connection error
- Order tidak ditemukan
- Validasi error

**Solusi:**
1. Cek logs di Render:
   - Buka service **mini-oms-backend**
   - Klik tab **Logs**
   - Lihat error message
   
2. Common errors:
   - `order not found` → Order ID salah atau tidak ada
   - `duplicate payment` → Payment sudah ada untuk order tersebut
   - `connection refused` → Database credentials salah

3. Fix database connection:
   - Pastikan DB_HOST, DB_USER, DB_PASS, DB_NAME benar
   - Test connection manual:
     ```bash
     # Di Render Shell (buka dari Dashboard)
     psql -h <db-host> -U <db-user> -d <db-name>
     ```

### Problem 4: Payment Created tapi Tidak Muncul di Order Detail

**Penyebab:**
- Frontend tidak refresh data
- Payment relation tidak di-preload

**Solusi:**
Pastikan order handler preload payment:
```go
// Seharusnya sudah ada di code
db.Preload("Payment").First(&order, "id = ?", id)
```

### Problem 5: Request Timeout / Slow

**Penyebab:**
- Free tier Render sleep setelah 15 menit tidak dipakai
- Database query lambat

**Solusi:**
1. **Auto-wake service**:
   - Pakai cron job gratis (UptimeRobot, Cron-job.org)
   - Ping backend setiap 10 menit
   
2. **Optimize query**:
   - Pastikan index sudah benar di database
   - Gunakan Preload untuk relasi

---

## Cek Status Deployment

### Backend Health Check
```bash
curl https://mini-oms-backend.onrender.com/api/products
```

Expected:
```json
{
  "success": true,
  "message": "Products retrieved successfully",
  "data": []
}
```

### Frontend Access
Buka browser:
```
https://mini-oms-frontend.onrender.com
```

Seharusnya bisa register, login, lihat products, create order, dan **create payment**.

---

## Testing Payment Flow di Production

1. **Register/Login** di frontend production
2. **Browse products** → Add to cart
3. **Create order**
4. **Navigate to order detail**
5. **Click "Create Payment"**
6. **Fill form** dengan:
   - Payment method: bank_transfer
   - Payment proof URL: https://i.imgur.com/example.jpg
   - Notes: "Test payment"
7. **Submit**

Jika berhasil:
- Redirect ke order detail
- Payment info muncul
- Status: pending

Jika error:
- Buka browser console (F12)
- Lihat error di tab Network
- Screenshot dan kasih tahu error message

---

## Update deployment (setelah fix code)

```bash
cd c:\Whaleee\Intern\TA

# Commit changes
git add .
git commit -m "fix: payment issue on production"
git push

# Render auto-detect push dan redeploy dalam 2-3 menit
```

---

## Environment Variables Template

### Backend (.env for local)
```env
PORT=8080
ENV=development
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=postgres
DB_NAME=mini_oms
JWT_SECRET=local-development-secret-key
JWT_EXPIRY_HOURS=24
```

### Backend (Render Environment Variables)
```
PORT=8080
ENV=production
DB_HOST=<dari-render-postgres>.postgres.render.com
DB_PORT=5432
DB_USER=<dari-render-postgres>
DB_PASS=<dari-render-postgres>
DB_NAME=<dari-render-postgres>
JWT_SECRET=production-random-very-long-secret-key-123456789
JWT_EXPIRY_HOURS=24
```

### Frontend (.env for local)
```env
VITE_API_BASE_URL=http://localhost:8080
```

### Frontend (Render Environment Variables)
```
VITE_API_BASE_URL=https://mini-oms-backend.onrender.com
```

---

## Monitoring & Debugging

### Render Dashboard
- **Logs**: Real-time logs dari backend/frontend
- **Metrics**: CPU, Memory, Request count
- **Events**: Deploy history
- **Shell**: Terminal akses ke container

### Browser DevTools (F12)
- **Console**: Error messages JavaScript
- **Network**: API requests & responses
- **Application**: localStorage untuk debug JWT token

---

## Next Steps After Successful Deploy

1. **Custom Domain** (opsional):
   - Beli domain (Namecheap, Cloudflare)
   - Setup di Render Settings → Custom Domains
   
2. **HTTPS** (otomatis di Render):
   - Render provide SSL certificate gratis
   
3. **Database Backup**:
   - Render auto-backup database daily
   - Manual backup: Download dari Render dashboard
   
4. **Monitoring**:
   - Setup UptimeRobot untuk monitor uptime
   - Email alert jika service down

---

Jika masih ada error payment di Render, screenshot error messagenya dan kasih tahu!
