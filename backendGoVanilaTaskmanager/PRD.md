# Product Requirements Document (PRD)
# Vanilla Go Task Manager API

---

## Document Metadata
- **Product Name:** Go Vanilla Task Manager API
- **Version:** 1.0.0 (Current As-Is & To-Be Roadmap)
- **Module Name:** `taskmanager`
- **Language / Stack:** Go (Standard Library `net/http`), PostgreSQL (`github.com/lib/pq`)
- **Status:** Draft / Active Development
- **Target Frontend:** Svelte SPA (Planned)

---

## 1. Executive Summary & Vision

### 1.1 Ringkasan Produk
**Go Vanilla Task Manager** adalah RESTful API manajemen tugas (*task management backend*) berkinerja tinggi yang dibangun menggunakan pendekatan **Vanilla Go (Murni / Standard Library)** tanpa framework web pihak ketiga (seperti Gin, Fiber, atau Echo). 

Tujuan utama produk ini adalah:
1. Menyediakan backend REST API yang ringan, cepat, dan terstruktur rapi dengan prinsip *clean / layered architecture*.
2. Memfasilitasi operasi manajemen tugas lengkap mulai dari CRUD dasar, filtering, searching, pagination, hingga soft-delete dan lifecycle status tugas.
3. Menjadi fondasi backend yang solid dan siap dihubungkan (*plug-and-play*) dengan aplikasi frontend modern (khususnya Svelte).

---

## 2. Analisis Fitur Eksisting (As-Is Architecture)

Berdasarkan audit menyeluruh terhadap kode sumber di repositori `backendGoVanilaTaskmanager`:

### 2.1 Arsitektur Kode Saat Ini
Sistem telah mengimplementasikan **Layered Architecture** dengan pemisahan tanggung jawab (*separation of concerns*):
- **`main.go`**: Entrypoint aplikasi, inisialisasi koneksi DB, dependency injection (Model -> Service -> Controller), registrasi route, dan server HTTP listener (`:8080`).
- **`config/database.go`**: Konfigurasi koneksi PostgreSQL via `database/sql` dan driver `github.com/lib/pq`.
- **`models/task.go`**: Definisi entitas data `Task`, query SQL mentah (*raw SQL*), dan interaksi database.
- **`services/task_service.go`**: Lapisan logika bisnis (*business logic*), validasi input (validasi tanggal *due date*, validasi *title*), dan kalkulasi metadata paginasi.
- **`controllers/task_controller.go`**: Handler HTTP, ekstraksi parameter URL/query/body JSON, dan pemetaan response.
- **`routes/routes.go`**: Definisi routing menggunakan standard `http.ServeMux`.
- **`utils/response.go`**: Standarisasi format response JSON (`APIResponse`).
- **`migrations/`**: Script migrasi database manual SQL.

```
📁 backendGoVanilaTaskmanager/
├── 📁 config/        # database.go (PostgreSQL connection)
├── 📁 controllers/   # task_controller.go (HTTP Request/Response Handler)
├── 📁 env/           # config.yaml, config.example.yaml (Belum terhubung ke main)
├── 📁 middlewares/   # .gitkeep (Masih kosong)
├── 📁 migrations/    # 001.init.up.sql, 002.add_deleted_at.up.sql
├── 📁 models/        # task.go (Task struct & DB operations)
├── 📁 routes/        # routes.go (Multiplexer & HTTP method dispatching)
├── 📁 services/      # task_service.go (Business validation & pagination logic)
├── 📁 tests/         # .gitkeep (Masih kosong)
├── 📁 utils/         # response.go (Standard Success/Error JSON builder)
├── 📄 .air.toml      # Konfigurasi Live-Reload development
├── 📄 go.mod / sum   # Module manifest (Go 1.25, lib/pq)
└── 📄 main.go        # Server entrypoint
```

---

### 2.2 Inventaris Fitur yang Sudah Dibuat

| Modul | Endpoint | HTTP Method | Deskripsi & Fungsionalitas | Status |
| :--- | :--- | :--- | :--- | :--- |
| **System** | `/health` | `GET` | Health check endpoint untuk memverifikasi API aktif. | ✅ Berjalan |
| **Tasks** | `/tasks` | `POST` | Membuat tugas baru. Validasi: `title` tidak boleh kosong, `due_date` tidak boleh di masa lalu. Default `completed = false`. | ✅ Berjalan |
| **Tasks** | `/tasks` | `GET` | Mengambil daftar tugas (aktif / belum dihapus) dengan pagination (`page`, `limit`), filter status (`completed=true/false`), dan pencarian (`search`). | ✅ Berjalan |
| **Tasks** | `/tasks/:id` | `GET` | Mengambil detail 1 tugas berdasarkan ID (hanya yang `deleted_at IS NULL`). | ✅ Berjalan |
| **Tasks** | `/tasks/:id` | `PUT` | Mengupdate `title` dan status `completed` tugas berdasarkan ID. | ✅ Berjalan |
| **Tasks** | `/tasks/:id` | `DELETE` | Menghapus tugas dengan metode **Soft Delete** (`deleted_at = NOW()`). | ✅ Berjalan |
| **Tasks** | `/tasks/:id/complete` | `PATCH` | Menandai tugas sebagai selesai (`completed = true`). | ✅ Berjalan |
| **Tasks** | `/tasks/:id/uncomplete` | `PATCH` | Menandai tugas sebagai belum selesai (`completed = false`). | ✅ Berjalan |
| **Tasks** | `/tasks/:id/restore` | `PATCH` | Memulihkan tugas yang telah di-soft-delete (`deleted_at = NULL`). | ✅ Berjalan |

---

### 2.3 Standar Format Response Saat Ini
Struktur JSON standar yang diterapkan di `utils/response.go`:

#### Success Response:
```json
{
  "status": "success",
  "message": "Task created successfully",
  "data": {
    "id": 1,
    "title": "Belajar Golang",
    "completed": false,
    "created_at": "2026-08-28T14:00:00Z"
  }
}
```

#### Paginated List Response:
```json
{
  "status": "success",
  "message": "success",
  "data": {
    "data": [
      {
        "id": 1,
        "title": "Belajar Golang",
        "completed": false,
        "created_at": "2026-08-28T14:00:00Z",
        "deleted_at": null
      }
    ],
    "meta": {
      "page": 1,
      "limit": 5,
      "total_data": 1,
      "total_page": 1,
      "has_next": false,
      "has_prev": false
    }
  }
}
```

#### Error Response:
```json
{
  "status": "error",
  "message": "Title tidak boleh kosong",
  "data": null
}
```

---

## 3. Gap Analysis & Temuan Masalah Teknis (Technical Debt)

Dari hasil investigasi mendalam, ditemukan beberapa ketidaksesuaian (*mismatch*) dan celah teknis yang harus diperbaiki:

1. **Sinkronisasi Skema Database vs Struct Model (`Task`):**
   - Di `models/task.go`, struct `Task` memiliki field: `SubTitle`, `Description`, `DueDate`, `UpdatedAt`.
   - Di file migrasi (`migrations/001.init.up.sql` dan `002.add_deleted_at.up.sql`), tabel `tasks` hanya memiliki kolom `id`, `title`, `completed`, `created_at`, `deleted_at`. Kolom `sub_title`, `description`, `due_date`, dan `updated_at` belum ada di tabel database!
   - Di query `GetAll`: `WHERE (title LIKE $1 OR description LIKE $1)` akan menimbulkan error SQL di PostgreSQL jika kolom `description` tidak ada di database fisik.
   - Query `Create` dan `Update` saat ini belum menyimpan / mengupdate `sub_title`, `description`, `due_date`, atau `updated_at`.

2. **Konfigurasi Database Hardcoded di `main.go`:**
   - File `env/config.yaml` sudah dibuat, namun `main.go` masih meng-hardcode:
     `config.NewPostgresDB("localhost", 5432, "postgres", "berjuang02", "taskmanager", "disable")`.

3. **Routing Manual String Slicing:**
   - Ekstraksi ID dilakukan manual: `r.URL.Path[len("/tasks/"):]` dan `strings.Contains(r.URL.Path, "/complete")`.
   - Karena Go versi 1.22+ telah mendukung routing pattern modern pada `http.ServeMux` (misal: `GET /tasks/{id}`, `PATCH /tasks/{id}/complete`), routing saat ini rentan bug jika ada trailing slash atau nested path.

4. **Middleware Belum Terimplementasi (`middlewares/` kosong):**
   - Belum ada **CORS Middleware** (krusial untuk frontend Svelte yang berjalan di beda port).
   - Belum ada **Logger Middleware** dan **Panic Recovery Middleware**.

5. **Koneksi Database & Graceful Shutdown:**
   - Belum ada konfigurasi *Connection Pool* (`SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime`).
   - Server belum menangani *graceful shutdown* (`os.Interrupt`, `syscall.SIGTERM`).

6. **Pengujian (Tests):**
   - Folder `tests/` masih kosong (`.gitkeep`).

---

## 4. Rekomendasi Fitur & Roadmap Pengembangan (To-Be)

Roadmap dibagi menjadi 4 Fase strategis:

```
┌──────────────────────────────────────────────────────────────────────────┐
│                             ROADMAP PRODUK                               │
├─────────────────┬─────────────────┬───────────────────┬──────────────────┤
│     FASE 1      │     FASE 2      │      FASE 3       │      FASE 4      │
│ Pondasi & Stabil│  Auth & Keamanan│ Fitur Lanjutan    │ Integrasi & Scale│
├─────────────────┼─────────────────┼───────────────────┼──────────────────┤
│• Schema Sync    │• User Auth (JWT)│• Kategori & Tag   │• Svelte Frontend │
│• YAML Config    │• Multi-Tenancy  │• Prioritas Tugas  │• OpenAPI/Swagger │
│• Go 1.22+ Mux   │• User Roles     │• Subtasks/Check   │• Stats/Dashboard │
│• Middlewares    │• Password Hash  │• Due Date Reminder│• WebSocket Notif │
│• Unit Testing   │• Profil User    │• Activity Audit   │• Export/Import   │
└─────────────────┴─────────────────┴───────────────────┴──────────────────┘
```

---

### 4.1 FASE 1: Pondasi, Stabilitas & Pembersihan Teknis (High Priority - Immediate)

#### 1.1 Sinkronisasi Skema Database & Model
- **Deskripsi:** Tambahkan migrasi `003.update_task_columns.up.sql` untuk menambahkan kolom `sub_title VARCHAR(255)`, `description TEXT`, `due_date TIMESTAMP NULL`, `updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP`.
- **Update Query:** Perbarui query SQL `Create`, `Update`, `GetByID`, `GetAll` agar konsisten membaca dan menulis semua field tersebut.

#### 1.2 Dynamic Configuration Loader
- **Deskripsi:** Buat package `config/config.go` untuk membaca `env/config.yaml` dengan fallback ke Environment Variables (`os.Getenv`), sehingga kredensial aman dan dapat dikonfigurasi per environment (Dev, Staging, Prod).

#### 1.3 Modernisasi Routing (Go 1.22+ Standard Mux)
- **Deskripsi:** Manfaatkan native path value `r.PathValue("id")` dan method routing bawaan Go 1.22+:
  - `GET /tasks`
  - `POST /tasks`
  - `GET /tasks/{id}`
  - `PUT /tasks/{id}`
  - `DELETE /tasks/{id}`
  - `PATCH /tasks/{id}/complete`
  - `PATCH /tasks/{id}/uncomplete`
  - `PATCH /tasks/{id}/restore`

#### 1.4 Core Middlewares
- **CORS Middleware:** Mengizinkan request dari frontend (Origin: `http://localhost:5173`, `http://localhost:3000`, methods: GET, POST, PUT, DELETE, PATCH, OPTIONS).
- **Logging Middleware:** Mencatat HTTP Method, Path, Remote IP, Status Code, dan Durasi Eksekusi.
- **Recovery Middleware:** Menangkap panic runtime dan mengembalikan HTTP 500 JSON tanpa membuat server crash.

#### 1.5 Database Connection Pool & Graceful Shutdown
- Konfigurasi `SetMaxOpenConns(25)`, `SetMaxIdleConns(5)`, `SetConnMaxLifetime(5 * time.Minute)`.
- Tangkap sinyal OS (`SIGINT`, `SIGTERM`) dengan `server.Shutdown(ctx)`.

#### 1.6 Unit & Integration Tests
- Unit testing pada `services/task_service_test.go` dan `controllers/task_controller_test.go` menggunakan `net/http/httptest` dan sql mock / test DB.

---

### 4.2 FASE 2: Autentikasi & Multi-Tenancy (Security & User Isolation)

#### 2.1 User Management & Authentication
- **User Registration & Login:** `POST /auth/register`, `POST /auth/login`, `POST /auth/logout`.
- **Password Hashing:** Menggunakan `golang.org/x/crypto/bcrypt`.
- **JWT (JSON Web Token) Implementation:** Stateless authentication via header `Authorization: Bearer <token>`.
- **Auth Middleware:** Verifikasi JWT dan injeksi `user_id` ke dalam `r.Context()`.

#### 2.2 Multi-Tenancy (Data Isolation)
- Tambahkan kolom `user_id INT REFERENCES users(id)` pada tabel `tasks`.
- Setiap operasi CRUD hanya dapat mengakses tugas milik user yang sedang terautentikasi (mencegah IDOR vulnerability).

---

### 4.3 FASE 3: Fitur Manajemen Tugas Tingkat Lanjut (Advanced Task Features)

#### 3.1 Prioritas Tugas (*Task Priority*)
- Enum/Tipe: `LOW`, `MEDIUM`, `HIGH`, `URGENT`.
- Filter tugas berdasarkan prioritas: `GET /tasks?priority=HIGH`.
- Sorting tugas berdasarkan prioritas atau due date: `GET /tasks?sort_by=due_date&order=asc`.

#### 3.2 Kategori & Label/Tagging
- Modul Kategori: `POST /categories`, `GET /categories`, `DELETE /categories/{id}`.
- Relasi tugas ke kategori atau banyak tag (*many-to-many*).

#### 3.3 Sub-tasks / Checklist Items
- Setiap tugas utama dapat memiliki checklist sub-tugas.
- Endpoint: `POST /tasks/{id}/subtasks`, `PATCH /subtasks/{subtaskId}/toggle`.
- Auto-calculate progress persentase tugas (misal: 3 dari 4 selesai = 75%).

#### 3.4 Activity Log / Riwayat Perubahan
- Audit log riwayat status (misal: "User A mengubah status menjadi Completed pada 14:00").

---

### 4.4 FASE 4: Dokumentasi API, Analitik, & Kesiapan Frontend Svelte

#### 4.1 Statistik & Dashboard Overview Endpoint
- `GET /tasks/analytics/summary`:
  - Total tugas aktif
  - Total tugas selesai
  - Total tugas overdue (lewat tenggat waktu)
  - Persentase penyelesaian
  - Distribusi per prioritas

#### 4.2 OpenAPI / Swagger Documentation
- Menyediakan endpoint `/swagger/index.html` atau spec `openapi.yaml` agar tim frontend (Svelte) dapat mengintegrasikan API dengan mudah.

#### 4.3 Bulk Operations & Export/Import
- `POST /tasks/bulk-delete` & `POST /tasks/bulk-complete`.
- `GET /tasks/export/csv` untuk download laporan tugas.

---

## 5. Spesifikasi Teknis & API Contract Detail (To-Be)

### 5.1 Skema Database Baru (Relational Schema)

```sql
-- Users Table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Categories Table
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    color_hex VARCHAR(7) DEFAULT '#3B82F6',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tasks Table (Enhanced)
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id INT NULL REFERENCES categories(id) ON DELETE SET NULL,
    title VARCHAR(255) NOT NULL,
    sub_title VARCHAR(255) NULL,
    description TEXT NULL,
    priority VARCHAR(20) DEFAULT 'MEDIUM' CHECK (priority IN ('LOW', 'MEDIUM', 'HIGH', 'URGENT')),
    completed BOOLEAN DEFAULT false,
    due_date TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Subtasks Table
CREATE TABLE subtasks (
    id SERIAL PRIMARY KEY,
    task_id INT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    completed BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

### 5.2 Daftar Rencana Endpoint Lengkap (Target API Spec)

#### 🔐 Auth Endpoints
| Method | Path | Deskripsi |
| :--- | :--- | :--- |
| `POST` | `/api/v1/auth/register` | Mendaftarkan akun baru |
| `POST` | `/api/v1/auth/login` | Login & mendapatkan JWT token |
| `GET` | `/api/v1/auth/me` | Mendapatkan profil user yang sedang login |

#### 📋 Tasks Endpoints
| Method | Path | Query Params | Deskripsi |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/v1/tasks` | `page, limit, completed, search, priority, category_id, sort_by, order` | List tasks berpaginasi |
| `POST` | `/api/v1/tasks` | - | Membuat tugas baru |
| `GET` | `/api/v1/tasks/{id}` | - | Detail tugas |
| `PUT` | `/api/v1/tasks/{id}` | - | Update seluruh data tugas |
| `DELETE` | `/api/v1/tasks/{id}` | `force=true/false` | Soft delete atau permanent delete |
| `PATCH` | `/api/v1/tasks/{id}/complete` | - | Tandai selesai |
| `PATCH` | `/api/v1/tasks/{id}/uncomplete` | - | Tandai belum selesai |
| `PATCH` | `/api/v1/tasks/{id}/restore` | - | Restore dari soft-deleted |
| `GET` | `/api/v1/tasks/trash` | `page, limit` | Melihat daftar tugas yang di-soft-delete |
| `GET` | `/api/v1/tasks/summary` | - | Statistik ringkasan tugas |

#### 🏷️ Categories Endpoints
| Method | Path | Deskripsi |
| :--- | :--- | :--- |
| `GET` | `/api/v1/categories` | List semua kategori user |
| `POST` | `/api/v1/categories` | Buat kategori baru |
| `PUT` | `/api/v1/categories/{id}` | Ubah nama / warna kategori |
| `DELETE` | `/api/v1/categories/{id}` | Hapus kategori |

---

## 6. Non-Functional Requirements (NFR)

1. **Performance & Concurrency:**
   - Waktu respon API p95 < 50ms untuk operasi CRUD standar.
   - Pemanfaatan *goroutines* dan *connection pooling* yang aman dari *race condition*.
2. **Security:**
   - Perlindungan password via bcrypt (cost factor >= 10).
   - SQL Injection Prevention (Semua query wajib menggunakan parameterized query `$1, $2`).
   - CORS policy terkonfigurasi ketat untuk lingkungan produksi.
3. **Reliability & Observability:**
   - Graceful shutdown dengan timeout 10 detik untuk memastikan request yang sedang berlangsung selesai.
   - Logging standar dengan timestamp, status code, dan request path.
4. **Code Quality:**
   - Zero external framework (menjaga dependensi seminimal mungkin: hanya driver database `lib/pq` atau alternatif modern `pgx`).
   - Mengikuti kaidah idiomatic Go (`gofmt`, *table-driven tests*, error wrapping).

---

## 7. Action Plan & Checklist Pengembangan (TODO List)

Berikut adalah checklist langkah eksekusi yang direkomendasikan:

### Task Checklist
- [ ] **[Fase 1.1] Perbaiki Database Migration:**
  - Buat migration SQL baru untuk menambahkan kolom `sub_title`, `description`, `due_date`, `updated_at`.
  - Update method `Create`, `Update`, `GetByID`, `GetAll` di `models/task.go`.
- [ ] **[Fase 1.2] Muat Konfigurasi dari `config.yaml` / `.env`:**
  - Buat parser YAML atau env loader di `config/`.
  - Hapus hardcoded DB string di `main.go`.
- [ ] **[Fase 1.3] Refactor Routing ke Go 1.22+ Standard Mux:**
  - Ganti manual URL string splitting dengan `r.PathValue("id")`.
- [ ] **[Fase 1.4] Implementasi Middleware:**
  - Buat CORS middleware di `middlewares/cors.go`.
  - Buat Logging middleware di `middlewares/logger.go`.
  - Buat Recovery middleware di `middlewares/recovery.go`.
- [ ] **[Fase 1.5] Connection Pooling & Graceful Shutdown:**
  - Tambahkan pool settings di `config/database.go`.
  - Tambahkan `signal.Notify` dan `server.Shutdown` di `main.go`.
- [ ] **[Fase 1.6] Tulis Unit Testing:**
  - Tulis test suite di `tests/` atau unit test per layer.
- [ ] **[Fase 2.1] Modul User & JWT Authentication:**
  - Tambahkan tabel `users`, implementasi register & login.
  - Tambahkan JWT Auth middleware.
  - Relasikan `tasks` dengan `user_id`.
- [ ] **[Fase 3.1] Fitur Kategori & Prioritas:**
  - Tambahkan kolom prioritas dan relasi kategori.
- [ ] **[Fase 4.1] Kesiapan Integrasi Frontend Svelte:**
  - Endpoint summary/analytics dan dokumentasi OpenAPI.

---
*Dokumen ini dirancang sebagai panduan tunggal (Single Source of Truth) untuk pengembangan fitur dan arsitektur Go Vanilla Task Manager.*
