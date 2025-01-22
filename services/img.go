package services

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

func ReadImage(path string) (image.Image, string, error) {
	srcFile, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	defer srcFile.Close()

	// Decode the source image
	img, format, err := image.Decode(srcFile)
	if err != nil {
		return nil, "", err
	}

	return img, format, nil
}

func GetColors(img image.Image) [][]color.Color {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	result := make([][]color.Color, height)
	for y := 0; y < height; y++ {
		row := make([]color.Color, width)
		result[y] = row
		for x := 0; x < width; x++ {
			result[y][x] = img.At(x, y)
		}
	}

	return result
}

type Tile struct {
	YStart, XStart, YEnd, XEnd int
}

func MakeTiles(width, height int, rows, cols int) []Tile {
	tileWidth, extraWidth := width / cols, width % cols
	tileHeight, extraHeight := height / rows, height % rows

	tiles := make([]Tile, 0)
	extraWidthCopy := extraWidth
	for h := 0; h < height; {
		// calc row height
		rowEnd := h + tileHeight
		if (extraHeight > 0) {
			rowEnd++
			extraHeight--
		}

		for w := 0; w < width; {
			tile := Tile { YStart: h, XStart: w }
			// add horizontal pixels
			w += tileWidth
			if (extraWidth > 0) {
				w++
				extraWidth--
			}
			tile.YEnd = rowEnd
			tile.XEnd = w
			tiles = append(tiles, tile)
		}
		
		// update row & reset extra width pixels
		extraWidth = extraWidthCopy
		h = rowEnd
	}

	return tiles
}

func DrawPallete(src image.Image, tiles []Tile) image.Image {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)

	for _, tile := range tiles {
		DrawTile(src, dst, tile)
	}

	return dst
}

func DrawTile(src image.Image, dst *image.RGBA, tile Tile) {
	var rTotal, gTotal, bTotal uint64
	totalPixels := uint64((tile.XEnd - tile.XStart) * (tile.YEnd - tile.YStart))
	for y := tile.YStart; y < tile.YEnd; y++ {
		for x := tile.XStart; x < tile.XEnd; x++ {
			r, g, b, _ := src.At(x, y).RGBA()
			rTotal += uint64(r >> 8)
			gTotal += uint64(g >> 8)
			bTotal += uint64(b >> 8)
		}
	}
	
	avgColor := color.RGBA {
		R: uint8(rTotal / totalPixels),
		G: uint8(gTotal / totalPixels),
		B: uint8(bTotal / totalPixels),
		A: 255,
	}

	for y := tile.YStart; y < tile.YEnd; y++ {
		for x := tile.XStart; x < tile.XEnd; x++ {
			dst.Set(x, y, avgColor)
		}
	}
}

func SaveImage(img image.Image, path string) {
	outFile, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	if err := png.Encode(outFile, img); err != nil {
		log.Fatal(err)
	}
}

