# Setup PostgreSQL di Podman

## Langkah-langkah Migrasi dari Native PostgreSQL ke Podman

### 1. Stop PostgreSQL Native (Wajib dilakukan oleh user)

```bash
sudo systemctl stop postgresql@15
sudo systemctl stop postgresql@18
sudo systemctl disable postgresql@15
sudo systemctl disable postgresql@18
```

Jika systemctl tidak berhasil, gunakan perintah force:
```bash
sudo pkill -9 postgres
```

### 2. Jalankan PostgreSQL di Podman

```bash
# Untuk pengguna Podman
podman-compose up -d

# Untuk pengguna Docker
docker-compose up -d
```

### 3. Verifikasi PostgreSQL berjalan di Podman

```bash
# Cek container
podman ps

# Cek log
podman logs taskmanager-postgres

# Test koneksi dari host
PGPASSWORD=postgres psql -h localhost -p 5432 -U postgres -c "\l"

# Test koneksi dari dalam container
podman exec -it taskmanager-postgres psql -U postgres -c "\l"
```

### 4. Setup Environment Variables untuk Backend

Anda bisa menggunakan environment variables atau file config.yaml. Environment variables akan meng-override nilai di config file.

**Option 1: Environment Variables (Recommended)**
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=taskmanager
export DB_SSLMODE=disable
```

**Option 2: Update config.yaml**
Update file `backendGoVanilaTaskmanager/env/config.yaml`:
```yaml
database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "postgres"
  dbname: "taskmanager"
  sslmode: "disable"
```

### 5. Test Backend Connection

```bash
cd backendGoVanilaTaskmanager
go run main.go
```

Jika berhasil, akan muncul:
```
✅ Connected to PostgreSQL
Server running at :8080
```

### 6. Restore database (jika ada backup)

Jika Anda memiliki backup database, restore dengan:

```bash
# Copy backup ke container
podman cp backup.sql taskmanager-postgres:/tmp/backup.sql

# Restore
podman exec -it taskmanager-postgres psql -U postgres -d taskmanager -f /tmp/backup.sql
```

## Perintah Podman Umum

```bash
# Start containers
podman-compose up -d

# Stop containers
podman-compose down

# View logs
podman-compose logs -f postgres

# Masuk ke container
podman exec -it taskmanager-postgres bash

# Connect ke PostgreSQL
podman exec -it taskmanager-postgres psql -U postgres

# List databases
podman exec -it taskmanager-postgres psql -U postgres -c "\l"

# Connect ke database spesifik
podman exec -it taskmanager-postgres psql -U postgres -d taskmanager
```

## Perintah Docker (untuk tim yang menggunakan Docker)

Semua perintah Podman di atas juga bekerja dengan Docker:

```bash
docker-compose up -d
docker-compose down
docker ps
docker logs taskmanager-postgres
docker exec -it taskmanager-postgres psql -U postgres
```

## Troubleshooting

### Port 5432 sudah digunakan
```bash
# Cek proses yang menggunakan port 5432
sudo netstat -tlnp | grep 5432

# Stop native PostgreSQL
sudo systemctl stop postgresql@15
sudo systemctl stop postgresql@18

# Atau kill proses
sudo pkill -9 postgres
```

### Container tidak bisa start
```bash
# Cek logs untuk error
podman logs taskmanager-postgres

# Remove container dan volume (WARNING: akan menghapus data)
podman-compose down -v
podman-compose up -d
```

### Backend tidak bisa connect
```bash
# Test koneksi dari host
PGPASSWORD=postgres psql -h localhost -p 5432 -U postgres -d taskmanager -c "\l"

# Cek apakah container berjalan
podman ps

# Cek config.yaml atau environment variables
```

## Kompatibilitas

Konfigurasi ini kompatibel dengan:
- Podman (penggunaan utama)
- Docker (untuk tim yang menggunakan Docker)

File `docker-compose.yml` dapat digunakan oleh kedua platform.

## Keuntungan Setup Ini

1. **Cross-platform compatibility** - Tim di Mac, Windows, Linux bisa gunakan config yang sama
2. **Version consistency** - Semua tim menggunakan PostgreSQL 15 yang sama
3. **Isolated environment** - Tidak konflik dengan native PostgreSQL
4. **Easy backup** - Volume Podman mudah di-backup dan restore
5. **Team collaboration** - Environment yang konsisten untuk semua developer
