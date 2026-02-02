# GitHub Push Guide

Karena GitHub tidak support password authentication lagi, kamu perlu Personal Access Token.

## Generate GitHub Token (Sekali Aja)

1. **Login ke GitHub** → https://github.com
2. **Klik profile (kanan atas)** → **Settings**
3. **Scroll ke bawah** → **Developer settings**
4. **Personal access tokens** → **Tokens (classic)**
5. **Generate new token (classic)**
6. **Isi form**:
   - Note: `Railway Deployment`
   - Expiration: `90 days` atau `No expiration`
   - Scopes: **Centang ✅ repo** (full control of private repositories)
7. **Generate token**
8. **COPY TOKEN** (contoh: `ghp_abc123xyz...`) - **SIMPAN** karena hanya muncul sekali!

## Push ke GitHub

### Cara 1: Pakai Git Credential Manager (Recommended)

```bash
cd c:\Whaleee\Intern\TA

# Push - akan muncul prompt Windows untuk login
git push -u origin main

# Windows akan tampilkan:
# 1. Browser popup untuk login GitHub
# 2. Atau prompt username/password
#    - Username: AntonyPutra
#    - Password: PASTE_TOKEN_DISINI (token yang kamu copy)
```

### Cara 2: Embed Token di URL (Quick but less secure)

```bash
cd c:\Whaleee\Intern\TA

# Ganti remote URL dengan token
git remote set-url origin https://ghp_YOUR_TOKEN_HERE@github.com/AntonyPutra/ProjectTAIntern.git

# Push
git push -u origin main
```

Ganti `ghp_YOUR_TOKEN_HERE` dengan token kamu!

## Verify Push Berhasil

Setelah push berhasil, cek di browser:
```
https://github.com/AntonyPutra/ProjectTAIntern
```

Kamu harus lihat semua file:
- mini-oms-backend/
- mini-oms-frontend/
- DEPLOYMENT_GUIDE.md
- dll

## Next: Deploy ke Railway

Setelah push berhasil:

1. **Buka https://railway.app**
2. **Sign Up** dengan GitHub
3. **New Project** → **Deploy from GitHub repo**
4. **Select**: ProjectTAIntern
5. **Add** PostgreSQL database
6. **Configure** environment variables
7. **Deploy!**

Railway akan auto-detect Golang dan React, deploy keduanya.
