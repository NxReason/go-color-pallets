package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeTilesCount(t *testing.T) {
	width, height, rows, cols := 100, 150, 4, 5
	want := rows * cols

	tiles := MakeTiles(width, height, rows, cols)
	if len(tiles) != want {
		t.Fatalf(`MakeTiles(_, _, 4, 5) should return %d tiles, got %d`, want, len(tiles))
	}
}

func TestMakeTiles(t *testing.T) {
	width, height, rows, cols := 11, 10, 3, 3
	want := []Tile {
		// row 1
		{ 0, 0, 4, 4 },
		{ 0, 4, 4, 8 },
		{ 0, 8, 4, 11 },
		// row 2
		{ 4, 0, 7, 4 },
		{ 4, 4, 7, 8 },
		{ 4, 8, 7, 11 },
		// row 3
		{ 7, 0, 10, 4 },
		{ 7, 4, 10, 8 },
		{ 7, 8, 10, 11 },
	}
	tiles := MakeTiles(width, height, rows, cols)

	for i, tile := range tiles {
		if tile != want[i] {
			t.Fatalf("MakeTiles(11, 10, 3, 3) want Tile: %v at index: %d, got: %v", want[i], i, tile)
		}
	}
}

func TestMakePath_Filename(t *testing.T) {
	want := "filename-pallete.jpg"
	got := makePath("filename.jpg", "pallete")

	assert.Equal(t, want, got)

	// dots in filename
	want = "file.name.test-grid.jpg"
	got = makePath("file.name.test.jpg", "grid")

	assert.Equal(t, want, got)
}

func TestMakePath_FileInsideFolder(t *testing.T) {
	want := "fold/filename-grid.jpg"
	got := makePath("fold/filename.jpg", "grid")

	assert.Equal(t, want, got)
}