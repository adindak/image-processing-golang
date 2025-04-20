package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
)

func main() {
    // Buka file gambar input
    inputFile, err := os.Open("image.png")
    if err != nil {
        log.Fatalf("Failed to open input image: %v", err)
    }
    defer inputFile.Close()

    // Decode file PNG menjadi gambar
    img, err := png.Decode(inputFile)
    if err != nil {
        log.Fatalf("Failed to decode image: %v", err)
    }

    // Ambil ukuran gambar
    bounds := img.Bounds()
    width, height := bounds.Max.X, bounds.Max.Y

    // Siapkan variabel untuk menyimpan koordinat border
    minX, minY := width, height  // Mulai dari nilai maksimum
    maxX, maxY := 0, 0           // Mulai dari nilai minimum

    // Buat file log untuk mencatat proses deteksi
    logFile, err := os.Create("debug.log")
    if err != nil {
        log.Fatalf("Failed to create log file: %v", err)
    }
    defer logFile.Close()

    // Fungsi untuk ngecek apakah suatu piksel itu hitam atau nggak
    isBlackPixel := func(r, g, b uint32) bool {
        // Ubah dari format 16-bit ke 8-bit biar lebih mudah
        r8 := r >> 8
        g8 := g >> 8
        b8 := b >> 8
        
        // Piksel dianggap hitam kalau nilai RGB-nya semua di bawah 30
        return r8 < 30 && g8 < 30 && b8 < 30
    }

    // Tulis header di file log
    fmt.Fprintln(logFile, "Starting border detection...")

    // Langkah 1: Scan secara horizontal untuk nemuin border atas dan bawah
    for y := 0; y < height; y++ {
        blackPixelCount := 0
        for x := 0; x < width; x++ {
            r, g, b, _ := img.At(x, y).RGBA()
            
            if isBlackPixel(r, g, b) {
                blackPixelCount++
                
                // Update koordinat min/max
                if x < minX {
                    minX = x
                }
                if x > maxX {
                    maxX = x
                }
                if y < minY {
                    minY = y
                }
                if y > maxY {
                    maxY = y
                }
                
                // Catat di log
                fmt.Fprintf(logFile, "Border pixel at (%d, %d) with RGB: (%d, %d, %d)\n", 
                    x, y, r>>8, g>>8, b>>8)
            }
        }
        
        // Kalau ada banyak piksel hitam di baris ini, mungkin ini border horizontal
        if blackPixelCount > width/10 {  // Threshold 10% dari lebar
            fmt.Fprintf(logFile, "Potential horizontal border at y=%d (black pixel count: %d)\n", 
                y, blackPixelCount)
        }
    }
    
    // Langkah 2: Scan secara vertikal untuk nemuin border kiri dan kanan
    for x := 0; x < width; x++ {
        blackPixelCount := 0
        for y := 0; y < height; y++ {
            r, g, b, _ := img.At(x, y).RGBA()
            
            if isBlackPixel(r, g, b) {
                blackPixelCount++
                // Koordinat min/max udah diupdate di scan horizontal
            }
        }
        
        // Kalau ada banyak piksel hitam di kolom ini, mungkin ini border vertikal
        if blackPixelCount > height/10 {  // Threshold 10% dari tinggi
            fmt.Fprintf(logFile, "Potential vertical border at x=%d (black pixel count: %d)\n", 
                x, blackPixelCount)
        }
    }

	// Langkah 3: Cari deretan piksel hitam horizontal untuk deteksi border
	for y := 0; y < height; y++ {
		streak := 0
		maxRun := 0
		
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			
			if r == 0 && g == 0 && b == 0 {
				streak++
			} else {
				maxRun = max(maxRun, streak)
				streak = 0
			}
		}
		
		// Cek piksel hitam di akhir baris
		maxRun = max(maxRun, streak)
		
		// Jika ditemukan garis hitam panjang
		if maxRun > width/3 {
			if y < minY {
				minY = y  // Border atas
				fmt.Fprintf(logFile, "Found top border at y=%d\n", y)
			} else if y > maxY {
				maxY = y  // Border bawah
				fmt.Fprintf(logFile, "Found bottom border at y=%d\n", y)
			}
		}
	}

    // Cek hasil deteksi koordinat
    fmt.Fprintf(logFile, "Detected bounds: (%d, %d) to (%d, %d)\n", minX, minY, maxX, maxY)
    
    // Langkah 4: Double-check kalo deteksi border sepertinya salah
    // Kalau tinggi area yang terdeteksi > 75% gambar, kemungkinan border bawah salah
    if maxY - minY > height * 3/4 {
        fmt.Fprintln(logFile, "Warning: Detected height seems too large, looking for bottom border...")
        
        // Cari border bawah dari bawah ke atas
        bottomBorderFound := false
        for y := height - 1; y > height/2 && y > minY; y-- {
            blackPixelCount := 0
            for x := minX; x <= maxX; x++ {
                r, g, b, _ := img.At(x, y).RGBA()
                if isBlackPixel(r, g, b) {
                    blackPixelCount++
                }
            }
            
            // Kalau baris ini punya cukup banyak piksel hitam, ini mungkin border bawah
            if blackPixelCount > (maxX - minX) / 3 {
                maxY = y
                bottomBorderFound = true
                fmt.Fprintf(logFile, "Found bottom border at y=%d (black pixels: %d)\n", 
                    y, blackPixelCount) 
                break
            }
        }
        
        if !bottomBorderFound {
            fmt.Fprintln(logFile, "Warning: Could not find a clear bottom border")
        }
    }
    
    // Langkah 5: Buat gambar baru dengan ukuran hasil cropping
    croppedWidth := maxX - minX + 1
    croppedHeight := maxY - minY + 1
    croppedImage := image.NewRGBA(image.Rect(0, 0, croppedWidth, croppedHeight))

    // Salin piksel dari gambar asli ke gambar hasil crop
    for y := minY; y <= maxY; y++ {
        for x := minX; x <= maxX; x++ {
            croppedImage.Set(x-minX, y-minY, img.At(x, y))
        }
    }

    // Simpan hasilnya ke file output.png
    outputFile, err := os.Create("output.png")
    if err != nil {
        log.Fatalf("Failed to create output file: %v", err)
    }
    defer outputFile.Close()

    // Encode gambar hasil crop ke format PNG
    if err := png.Encode(outputFile, croppedImage); err != nil {
        log.Fatalf("Failed to encode output image: %v", err)
    }

    // Tampilkan pesan sukses
    fmt.Println("Successfully cropped image to border bounds and saved as output.png")
    fmt.Println("Border detection details saved to border_detection.log")
}
