# SintaISSNChecker

<p align="center">
<img src="https://github.com/user-attachments/assets/e2f15cb6-a52f-4bab-95ae-7ea13d1a3514" alt="SintaISSNChecker" width="400">
</p>

**SintaISSNChecker** adalah sebuah aplikasi berbasis Go yang digunakan untuk memeriksa status akreditasi jurnal berdasarkan ISSN di situs SINTA (Science and Technology Index) dan mengelola file PDF terkait. Aplikasi ini secara otomatis mengganti nama file PDF jika ISSN yang terdeteksi memiliki status akreditasi di SINTA.

## Cara Kerja

- Ekstraksi nomor ISSN dari file PDF.
- Cek status akreditasi ISSN di situs SINTA.
- Penamaan ulang file PDF berdasarkan status akreditasi di SINTA.

## Cara Penggunaan

1. **Siapkan File PDF**
   <br>Masukkan semua file PDF jurnal yang ingin Anda proses ke dalam folder `pdf`. Jika folder tersebut belum ada, silakan buat terlebih dahulu di direktori utama proyek.

2. **Jalankan Program**
   <br>Buka terminal di direktori utama proyek dan jalankan perintah berikut:
   ```bash
   go run ./cmd/app/main.go
   ```
   Program akan secara otomatis memproses semua file di dalam folder `pdf`. File yang jurnalnya terakreditasi SINTA akan diganti namanya, contohnya: `SINTA 2 Accredited - Nama Jurnal Asli.pdf`.

