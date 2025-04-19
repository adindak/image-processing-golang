package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

func isBlack(c color.Color) bool {
	r, g, b, _ := c.RGBA()
	return r == 0 && g == 0 && b == 0
}

func main() {
	// Buka file image
	// Letakkan sejajar dengan main.go
	file, err := os.Open("image.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Decode gambar
	// Olah supaya bisa jadi objek image.Image
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	// Ambil batas minimum dan maksimum dari koordinat border yg berwarna hitam
	bounds := img.Bounds()
	minX, minY := bounds.Max.X, bounds.Max.Y
	maxX, maxY := bounds.Min.X, bounds.Min.Y

	// Simpan dalam bentuk log
	logFile, _ := os.Create("zerolog.log")
	defer logFile.Close()

	// Loopong pixel satu persatu
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// pengecekan apakah pixelnya berwarna Hitam
			if isBlack(img.At(x, y)) {
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
				logFile.WriteString((fmt.Sprintf("Black Pixel: (%d, %d)\n", x, y)))
			}
		}
	}

	// Membuat objek jadi image.Rectangle
	croppedRectangle := image.Rect(minX, minY, maxX+1, maxY+1)
	croppedImg := image.NewRGBA(croppedRectangle)

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			croppedImg.Set(x, y, img.At(x, y))
		}
	}

	// Simpan
	outFile, err := os.Create("output.png")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	err = png.Encode(outFile, croppedImg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("✅ Cropped image saved to output.png")
}
