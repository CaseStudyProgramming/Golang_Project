backend Golang Vanila atau Murni Taskmanager

layered architecture
my-backend-app/
├── config/             # Konfigurasi database, env, dan library pihak ketiga
├── controllers/        # Menangani HTTP request & response (Validator input)
├── models/             # Skema database / ORM (Object-Relational Mapping)
├── routes/             # Definisi URL / Endpoint API
├── services/           # Logika bisnis utama (Business Logic)
├── middlewares/        # Fungsi perantara (Autentikasi, Logging, Error Handling)
├── utils/              # Fungsi pembantu / Helper (Format tanggal, enkripsi)
├── tests/              # File pengujian (Unit testing, Integration testing)
├── .env                # Variabel lingkungan rahasia (DB password, API key)
├── server.js (atau app.js) # Entry point utama aplikasi
└── package.json        # Manifest proyek dan daftar dependensi
