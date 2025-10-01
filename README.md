# SintaISSNChecker

<p align="center">
<img src="https://github.com/user-attachments/assets/c2df2bfb-5032-469c-89f6-2821deb58a2b" alt="SintaISSNChecker" width="400">
</p>

**SintaISSNChecker** adalah aplikasi berbasis Go untuk memeriksa status akreditasi jurnal secara massal melalui situs SINTA (Science and Technology Index).

## Cara Kerja

Aplikasi ini bekerja dalam beberapa langkah sederhana:

1.  **Pilih Folder**: Pilih file PDF atau folder yang berisi kumpulan file PDF jurnal.
2.  **Ekstraksi ISSN**: Aplikasi akan memindai setiap file PDF untuk mengekstrak nomor ISSN.
3.  **Pengecekan SINTA**: Setiap ISSN yang ditemukan akan diperiksa status akreditasinya di situs web SINTA.
4.  **Ubah Nama File**: File PDF dari jurnal yang terakreditasi SINTA akan diganti namanya secara otomatis dengan format `[Peringkat SINTA] - [Nama File Asli].pdf`. Contoh: `SINTA 2 Accredited - Nama Jurnal Asli.pdf`.

## Cara Penggunaan

1. **Build Aplikasi**
   <br>Buka terminal di direktori utama proyek dan jalankan perintah berikut:
   ```bash
   go run ./cmd/app/main.go
   ```

2. **Pilih Folder & Mulai Proses**
   - Setelah aplikasi terbuka, klik tombol **"Pilih Folder"** untuk memilih direktori yang berisi file-file PDF Anda.
   - Klik tombol **"Mulai Proses"** untuk memulai pemindaian dan penggantian nama file.
   - Anda dapat memantau progres dan melihat log detail proses langsung di jendela aplikasi.
