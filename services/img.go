package services

import (
	"color-pallete/cmd"
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"os"
	"strings"
)

type PaintFunc func(image.Image, []Tile, *image.RGBA, []Tile) image.Image

type GPResult struct {
	path string
	mode cmd.Mode
	err error
}

func ProcessFiles(config cmd.Config) []error {
	errs := make([]error, 0)
	filesCount := len(config.Modes) * len(config.InputFiles)
	imageProcessingCh := make(chan GPResult, filesCount)

	for _, path := range config.InputFiles {
		for _, m := range config.Modes {
			go ProcessFileAsync(path, config, cmd.Mode(m), imageProcessingCh)
		}
	}

	for i := 0; i < filesCount; i++ {
		res := <-imageProcessingCh
		if res.err != nil {
			errs = append(errs, res.err)
			continue
		}
		fmt.Printf("done: [ %s ] for %s (%d / %d)\n", string(res.mode), res.path, i + 1, filesCount)
	}
	
	return errs
}

func ProcessFileAsync(path string, config cmd.Config, mode cmd.Mode, ch chan GPResult) {
	img, _, err := ReadImage(path)
	if err != nil {
		ch <- GPResult{ path, mode, err}
	}
	colors := GetColors(img)
	inTiles := MakeTiles(len(colors[0]), len(colors), config.GridRows, config.GridCols)

	dstBounds := img.Bounds()
	outTiles := inTiles
	shouldUseConfigBounds := (mode == cmd.PALLETE) && (config.OutputHeight > 0 && config.OutputWidth > 0)
	if shouldUseConfigBounds {
		dstBounds = image.Rectangle{
			image.Point{ 0, 0 },
			image.Point{ config.OutputWidth, config.OutputHeight },
		}
		outTiles = MakeTiles(config.OutputWidth, config.OutputHeight, config.GridRows, config.GridCols)
	}
	dst := image.NewRGBA(dstBounds)
	
	var Paint PaintFunc
	switch mode {
	case cmd.GRID:
		Paint = DrawGrid
	case cmd.PALLETE:
		Paint = DrawPallete
	default:
		ch <- GPResult{ path, mode, errors.New("invalid paint mode, expected [GRID | PALLETE], got " + string(mode)) }
	}
	output := Paint(img, inTiles, dst, outTiles)

	suffix := strings.ToLower(string(mode))
	err = SaveImage(output, makePath(path, suffix))
	if err != nil {
		ch <- GPResult{ path, mode, err }
	}

	ch <- GPResult{ path, mode, nil}
}

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

func DrawPallete(src image.Image, inTiles []Tile, dst *image.RGBA, outTiles []Tile) image.Image {
	for i := range inTiles {
		DrawTile(src, inTiles[i], dst, outTiles[i])
	}

	return dst
}

func DrawTile(src image.Image, inTile Tile, dst *image.RGBA, outTile Tile) {
	var rTotal, gTotal, bTotal uint64
	totalPixels := uint64((inTile.XEnd - inTile.XStart) * (inTile.YEnd - inTile.YStart))
	for y := inTile.YStart; y < inTile.YEnd; y++ {
		for x := inTile.XStart; x < inTile.XEnd; x++ {
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

	for y := outTile.YStart; y < outTile.YEnd; y++ {
		for x := outTile.XStart; x < outTile.XEnd; x++ {
			dst.Set(x, y, avgColor)
		}
	}
}

func DrawGrid(src image.Image, tiles []Tile, dst *image.RGBA, _ []Tile) image.Image {
	lineColor := color.RGBA { 0, 0, 0, 255 }
	for _, tile := range tiles {
		// copy non-grid colors
		for y := tile.YStart; y < tile.YEnd; y++ {
			for x := tile.XStart; x < tile.XEnd; x++ {
				dst.Set(x, y, src.At(x, y))
			}
		}

		// paint vertical lines
		if tile.XEnd != src.Bounds().Max.X {
			for y := tile.YStart; y < tile.YEnd; y++ {
				dst.Set(tile.XEnd - 1, y, lineColor)
			}
		}

		// paint horizontal lines
		if tile.YEnd != src.Bounds().Max.Y {
			for x := tile.XStart; x < tile.XEnd; x++ {
				dst.Set(x, tile.YEnd - 1, lineColor)
			}
		}
	}

	return dst
}

func SaveImage(img image.Image, path string) error {
	outFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if err = png.Encode(outFile, img); err != nil {
		return err
	}
	return nil
}

func makePath(original, suffix string) string {
	parts := strings.Split(original, ".")
	name := join(parts[:len(parts) - 1], ".")
	name += "-" + suffix
	return name + "." + parts[len(parts) - 1]
}

func join(parts []string, ch string) string {
	var sb strings.Builder
	for _, part := range parts[:len(parts) - 1] {
		sb.WriteString(part)
		sb.WriteString(ch)
	}
	sb.WriteString(parts[len(parts) - 1])
	return sb.String()
}