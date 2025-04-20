#### Overview ####
Code untuk melakukan pengecekan bidang persegi panjang dari suatu gambar

#### How to Run ####
1. Clone this repo
2. Open the project directory
3. Run `go mod tidy`
4. Run `make run`

#### Pseudo Code Image Processing (Cropped Rectangle) ####

1. Buka file gambar input (image.png)
2. Decode gambar PNG dan ambil ukurannya
3. Siapin variabel buat nyimpan koordinat border (minX, minY, maxX, maxY)
4. Bikin file log buat nyatet proses debug (debug.log)

5. Bikin fungsi bantu isBlackPixel(r, g, b):
   - Balikin true kalau piksel dianggap hitam (RGB < 30)

6. Langkah 1: Scan horizontal (per baris) buat cari border atas & bawah:
   - Hitung jumlah piksel hitam per baris
   - Kalau banyak piksel hitam, kemungkinan itu border → update minY / maxY
   - Juga update minX / maxX kalau nemu piksel hitam

7. Langkah 2: Scan vertikal (per kolom) buat cari border kiri & kanan:
   - Hitung jumlah piksel hitam per kolom
   - Log kolom yang kelihatannya jadi border vertikal

8. Langkah 3: Deteksi deretan piksel hitam yang panjang (streak):
   - Kalau nemu deretan piksel hitam yang panjangnya lebih dari 1/3 lebar gambar,
     kemungkinan itu garis border → update minY / maxY

9. Langkah 4: Validasi hasil deteksi
   - Kalau tinggi hasil deteksi >75% tinggi gambar, mungkin border bawah salah
   - Coba scan ulang dari bawah buat cari border bawah yang bener

10. Langkah 5: Crop gambar berdasarkan koordinat border yang ditemukan:
    - Bikin gambar baru dengan ukuran sesuai area crop
    - Salin piksel dari gambar asli ke gambar baru

11. Simpan hasil crop ke file output.png

12. Selesai!
