# Svelte Task Manager Frontend

Frontend Svelte SPA untuk Go Vanilla Task Manager API.

## Struktur Folder

```
svelte-taskmanager-frontend/
├── src/
│   ├── lib/
│   │   ├── components/
│   │   │   └── ui/
│   │   │       ├── auth/          # Komponen autentikasi (Login, Register)
│   │   │       ├── tasks/         # Komponen manajemen tugas
│   │   │       ├── categories/    # Komponen manajemen kategori
│   │   │       ├── dashboard/     # Komponen dashboard utama
│   │   │       └── shared/        # Komponen reusable umum
│   │   ├── stores/                # Svelte stores untuk state management
│   │   ├── services/              # API client services
│   │   ├── utils/                 # Helper functions dan utilities
│   │   ├── types/                 # TypeScript type definitions
│   │   └── hooks/                 # Custom Svelte hooks
│   └── routes/
│       ├── auth/                  # Halaman autentikasi
│       ├── tasks/                 # Halaman manajemen tugas
│       ├── categories/            # Halaman manajemen kategori
│       ├── dashboard/             # Halaman dashboard
│       └── settings/              # Halaman pengaturan
├── static/
│   └── assets/
│       ├── images/                # Gambar dan aset visual
│       └── icons/                 # Icon library
├── tests/                         # Unit dan integration tests
└── (Config files: package.json, vite.config.js, etc.)
```

## Fitur Utama

- **Autentikasi User**: Login, Register, Logout dengan JWT
- **Manajemen Tugas**: CRUD, pagination, filter, search, soft-delete
- **Kategori & Tagging**: Organisasi tugas dengan kategori dan tags
- **Dashboard**: Statistik dan overview tugas
- **Subtasks**: Checklist items dalam tugas
- **Activity Log**: Riwayat perubahan tugas
- **Responsive Design**: Mendukung berbagai ukuran layar

## Tech Stack

- **Framework**: SvelteKit / Svelte
- **Build Tool**: Vite
- **Language**: TypeScript
- **State Management**: Svelte Stores
- **Styling**: CSS / Tailwind CSS (opsional)
- **API Integration**: Fetch API / Axios

## API Backend

Backend API tersedia di: `http://localhost:8080`

Dokumentasi API tersedia di: `http://localhost:8080/swagger/index.html`

## Setup Project

1. Install dependencies:
```bash
npm install
```

2. Start development server:
```bash
npm run dev
```

3. Build untuk production:
```bash
npm run build
```

## Environment Variables

Buat file `.env` dengan konfigurasi:
```
VITE_API_BASE_URL=http://localhost:8080
VITE_APP_TITLE=Task Manager
```
