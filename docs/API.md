# Todo API Reference

Kelima server (`chi-todo`, `echo-todo`, `fiber-todo`, `gin-todo`, dan `native-todo`) menyediakan kontrak REST API yang sama dan memakai tabel `todos.todos`.

## Base URL

| Implementasi | Base URL |
| --- | --- |
| Chi | `http://localhost:8081` |
| Echo | `http://localhost:8082` |
| Fiber | `http://localhost:8083` |
| Gin | `http://localhost:8084` |
| Go native | `http://localhost:8085` |

## Format Todo

```json
{
  "id": 1,
  "user_id": 1,
  "title": "Bandingkan framework",
  "description": "Implementasi CRUD Todo",
  "completed": false,
  "due_date": "2026-10-01",
  "created_at": "2026-09-13T10:00:00Z",
  "updated_at": "2026-09-13T10:00:00Z"
}
```

`user_id` dan `title` wajib diisi saat membuat atau memperbarui Todo. `description` dan `due_date` boleh bernilai `null`. Format `due_date` adalah `YYYY-MM-DD`.

## Health check

`GET /health`

Memeriksa koneksi aplikasi ke MariaDB.

```sh
curl http://localhost:8081/health
```

Respons `200`:

```json
{"status":"ok"}
```

## Daftar Todo

`GET /todos`

Parameter query opsional:

| Parameter | Keterangan |
| --- | --- |
| `user_id` | Menampilkan Todo untuk pengguna tertentu. |

```sh
curl 'http://localhost:8081/todos?user_id=1'
```

Respons `200` adalah array Todo. Jika tidak ada data, responsnya `[]`.

## Detail Todo

`GET /todos/{id}`

```sh
curl http://localhost:8081/todos/1
```

Respons `200` berisi satu objek Todo. Jika ID tidak ada, responsnya:

```json
{"error":"todo not found"}
```

Status: `404`.

## Membuat Todo

`POST /todos`

```sh
curl -X POST http://localhost:8081/todos \
  -H 'Content-Type: application/json' \
  -d '{
    "user_id": 1,
    "title": "Bandingkan framework",
    "description": "Implementasi CRUD Todo",
    "completed": false,
    "due_date": "2026-10-01"
  }'
```

Respons `201` berisi Todo yang baru dibuat.

## Memperbarui Todo

`PUT /todos/{id}`

Endpoint ini melakukan penggantian penuh. Kirim semua field yang dapat diubah, termasuk `user_id` dan `title`.

```sh
curl -X PUT http://localhost:8081/todos/1 \
  -H 'Content-Type: application/json' \
  -d '{
    "user_id": 1,
    "title": "Bandingkan framework - selesai",
    "description": "CRUD sudah diuji",
    "completed": true,
    "due_date": null
  }'
```

Respons `200` berisi Todo terbaru. Jika ID tidak ada, responsnya `404`.

## Menghapus Todo

`DELETE /todos/{id}`

```sh
curl -i -X DELETE http://localhost:8081/todos/1
```

Respons berhasil adalah `204 No Content`. Jika ID tidak ada, responsnya `404`.

## Error umum

| Status | Arti |
| --- | --- |
| `400` | JSON tidak valid, atau `user_id`/`title` kosong. |
| `404` | Todo tidak ditemukan. |
| `500` | Terjadi kesalahan saat menjalankan operasi database. |
| `503` | Database tidak dapat dihubungi oleh endpoint `/health`. |
