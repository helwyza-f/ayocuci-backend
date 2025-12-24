# Ayo Cuci Backend (Clean Architecture)

Backend untuk **Ayo Cuci – Laundry Management System**, dibangun menggunakan **Golang** dan **Gin Framework**.  
Backend ini dirancang modular, scalable, dan siap production untuk mendukung operasional laundry **multi-outlet** dengan sistem autentikasi yang aman.

---

## 🛠️ Tech Stack

- **Bahasa**: Go (1.22+)
- **Framework**: Gin Gonic
- **ORM**: GORM
- **Database**:
  - SQLite (Development)
  - PostgreSQL (Production)
- **Autentikasi**: JWT (JSON Web Token)
- **Konfigurasi**: Environment Variable (`.env`)

---

## 🧱 Arsitektur

Project ini menggunakan pendekatan **Clean Architecture** agar kode mudah dirawat dan dikembangkan.

### Struktur Folder

```
.
├── cmd/
│   └── server/
│       └── main.go        # Entry point aplikasi
├── internal/
│   └── module/
│       ├── auth/          # Modul autentikasi
│       ├── outlet/        # Manajemen outlet
│       ├── employee/      # Manajemen karyawan & role
│       ├── service/       # Layanan laundry & harga
│       └── order/         # Workflow order laundry
├── middleware/
│   ├── auth.go            # Middleware JWT
│   └── role.go            # Middleware role-based access
├── config/
│   └── database.go        # Inisialisasi database
├── .env.example
├── go.mod
└── README.md
```

### Pembagian Layer

- **Handler / Controller**  
  Menangani request dan response HTTP
- **Service / Usecase**  
  Berisi logika bisnis
- **Repository**  
  Akses dan manipulasi database
- **Entity / Model**  
  Representasi data dan domain

---

## 🔐 Autentikasi & Otorisasi

- Menggunakan **JWT** untuk autentikasi
- Mendukung **Role-Based Access Control**:
  - Owner
  - Employee
- Validasi token dilakukan melalui middleware Gin

---

## 🌐 Environment Variable

Buat file `.env` di root project:

```
APP_ENV=development
APP_PORT=8080
JWT_SECRET=secret_jwt_kamu

# Database PostgreSQL
DB_DRIVER=postgres
DB_URL=postgresql://username:password@host:port/nama_database
```

Untuk development, aplikasi juga dapat dijalankan menggunakan SQLite.

---

## ▶️ Cara Menjalankan (Lokal)

### 1. Clone Repository

```bash
git clone https://github.com/helwyza-f/ayocuci-backend.git
cd ayocuci-backend
```

### 2. Install Dependency

```bash
go mod download
```

### 3. Setup Environment

```bash
cp .env.example .env
```

Lalu sesuaikan isi file `.env`.

### 4. Jalankan Server

```bash
go run ./cmd/server/main.go
```

Server akan berjalan di:

```
http://localhost:8080
```

---

## 📡 Fitur API (MVP)

- Login & autentikasi user
- Manajemen outlet
- Manajemen karyawan
- Katalog layanan laundry
- Workflow order laundry
- Rekap keuangan dasar

---

## 🚀 Siap Production

- Mendukung PostgreSQL
- Struktur modular (Clean Architecture)
- JWT Authentication
- Mudah dikembangkan untuk fitur lanjutan

---

## 👨‍💻 Developer

**Helwiza Fahry**

> _Membangun backend yang scalable dengan clean architecture._

---

## 📌 Catatan

- Pastikan database sudah berjalan sebelum server dijalankan
- Gunakan Postman / Insomnia untuk testing endpoint API

---

✨ **Backend siap diintegrasikan dengan aplikasi Flutter Ayo Cuci.**
