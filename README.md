# Go Todo API Framework Comparison

Satu REST API Todo yang diimplementasikan dengan lima pilihan HTTP framework Go. Semua server memakai kontrak endpoint, validasi, query, dan tabel MariaDB yang sama, sehingga fokus perbandingan ada pada cara masing-masing framework menangani HTTP.

## Implementasi

| Implementasi | Folder | Port bawaan |
| --- | --- | --- |
| [Chi](https://github.com/go-chi/chi) | [`chi-todo`](./chi-todo) | `8081` |
| [Echo](https://github.com/labstack/echo) | [`echo-todo`](./echo-todo) | `8082` |
| [Fiber](https://github.com/gofiber/fiber) | [`fiber-todo`](./fiber-todo) | `8083` |
| [Gin](https://github.com/gin-gonic/gin) | [`gin-todo`](./gin-todo) | `8084` |
| Go native (`net/http`) | [`native-todo`](./native-todo) | `8085` |

## Apa yang dibandingkan?

Kode yang sama untuk semua server dipusatkan di [`common/todo.go`](./common/todo.go): konfigurasi MariaDB, model, validasi input, dan operasi SQL CRUD. Folder implementasi hanya berisi adapter HTTP masing-masing framework.

Setiap `main.go` memiliki komentar `*-specific` pada bagian yang perlu diperhatikan saat membandingkan:

- cara mendeklarasikan route;
- cara mengambil parameter `id`;
- cara membaca/binding JSON request;
- tipe context dan pembuatan respons.

Struktur proyek:

```text
.
├── common/          # kode domain dan database yang dibagi bersama
├── chi-todo/        # adapter HTTP Chi
├── echo-todo/       # adapter HTTP Echo
├── fiber-todo/      # adapter HTTP Fiber
├── gin-todo/        # adapter HTTP Gin
├── native-todo/     # adapter HTTP net/http
├── sql/             # inisialisasi MariaDB
└── docs/            # referensi API
```

## Prasyarat

- Go `1.26` atau lebih baru
- MariaDB/MySQL yang dapat diakses dari host
- Database `todos` dan tabel `todos` (dapat dibuat menggunakan skrip di bawah)

## Menyiapkan database

Skrip idempoten [sql/001_todos.sql](./sql/001_todos.sql) membuat database `todos` dan tabel `todos` tanpa menghapus data yang telah ada.

Untuk MariaDB dalam container Docker bernama `mariadb`:

```sh
docker exec -i mariadb mariadb -h 127.0.0.1 -uadmin -p'jamurkembang' < sql/001_todos.sql
```

Untuk MariaDB yang terpasang langsung di host:

```sh
mariadb -h 127.0.0.1 -P 3306 -uadmin -p < sql/001_todos.sql
```

## Konfigurasi

Seluruh implementasi membaca variabel environment berikut. Nilai bawaan cocok dengan setup lokal yang digunakan proyek ini, tetapi sebaiknya password diatur sebagai environment variable saat dipakai di environment lain.

| Variabel | Nilai bawaan |
| --- | --- |
| `DB_HOST` | `127.0.0.1` |
| `DB_PORT` | `3306` |
| `DB_USER` | `admin` |
| `DB_PASSWORD` | `jamurkembang` |
| `DB_NAME` | `todos` |
| `PORT` | port sesuai implementasi |

Contoh override:

```sh
DB_PASSWORD='password-lokal-anda' PORT=9000 go run ./chi-todo
```

## Menjalankan server

Jalankan dari root repository:

```sh
go run ./chi-todo
go run ./echo-todo
go run ./fiber-todo
go run ./gin-todo
go run ./native-todo
```

Server dapat dijalankan bersamaan karena port bawaannya berbeda. Verifikasi salah satunya:

```sh
curl http://localhost:8081/health
```

## Endpoint ringkas

| Method | Path | Keterangan |
| --- | --- | --- |
| `GET` | `/health` | Memeriksa koneksi MariaDB. |
| `GET` | `/todos` | Mengambil Todo; mendukung `?user_id=1`. |
| `GET` | `/todos/{id}` | Mengambil satu Todo. |
| `POST` | `/todos` | Membuat Todo. |
| `PUT` | `/todos/{id}` | Mengganti Todo secara penuh. |
| `DELETE` | `/todos/{id}` | Menghapus Todo. |

Lihat [docs/API.md](./docs/API.md) untuk referensi lengkap: payload, respons, status error, dan contoh `curl`.

## Verifikasi build

```sh
go test ./common ./chi-todo ./echo-todo ./fiber-todo ./gin-todo ./native-todo
```
