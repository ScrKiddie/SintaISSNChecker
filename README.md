# SintaISSNChecker

<p align="center">
<img src="https://github.com/user-attachments/assets/43d20aad-79bd-48c4-a92d-8b53a7347809" alt="SintaISSNChecker" width="400">
</p>

**SintaISSNChecker** adalah sebuah aplikasi berbasis Go yang digunakan untuk memeriksa status akreditasi jurnal berdasarkan ISSN di situs SINTA (Science and Technology Index) dan mengelola file PDF terkait. Aplikasi ini secara otomatis mengganti nama file PDF jika ISSN yang terdeteksi memiliki status akreditasi di SINTA.

## Cara Kerja

- Ekstraksi nomor ISSN dari file PDF.
- Cek status akreditasi ISSN di situs SINTA.
- Penamaan ulang file PDF berdasarkan status akreditasi di SINTA.

## Cara Pakai

1. **Compile Program**

2. **Masukkan Key UniPDF**
   <br> Edit file `.env` di root folder proyek dan masukkan **UNIPDF_KEY** Anda yang didapatkan dari [UniPDF Metered License](https://www.unidoc.io/unipdf/):
      ```
      UNIPDF_KEY=YOUR_UNIPDF_KEY
      ```

3. **Isi Folder PDF**
   <br>Siapkan file-file PDF jurnal yang ingin diproses dan masukkan ke dalam folder `pdf` di root proyek.

4. **Jalankan Program**
   <br>Program akan memproses semua file PDF dalam folder `pdf`, mengekstrak ISSN, mengecek status akreditasi di SINTA, dan mengganti nama file jika ISSN terdaftar di SINTA.