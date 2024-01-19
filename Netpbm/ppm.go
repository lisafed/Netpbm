package Netpbm

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PPM struct {
	data          [][]Pixel
	width, height int
	magicNumber   string
	max           uint8
}
type Pixel struct {
	R, G, B uint8
}

// ReadPPM reads a PPM image from a file and returns a struct that represents the image.
func ReadPPM(filename string) (*PPM, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var ppm PPM
	var headerSize int
	var counter int

	splitFile := strings.SplitN(string(file), "\n", -1)

	for _, line := range splitFile {
		if strings.HasPrefix(line, "#") && ppm.max != 0 {
			headerSize = counter
		}
		fields := strings.Fields(line)
		switch {
		case ppm.magicNumber == "" && (fields[0] == "P3" || fields[0] == "P6"):
			ppm.magicNumber = fields[0]
		case ppm.width == 0 && ppm.height == 0 && len(fields) >= 2:
			ppm.width, _ = strconv.Atoi(fields[0])
			ppm.height, _ = strconv.Atoi(fields[1])
			headerSize = counter
		case ppm.max == 0 && ppm.width != 0:
			maxval, err := strconv.Atoi(fields[0])
			if err != nil {
				return nil, err // gestion de l'erreur
			}

			if maxval < 0 || maxval > 255 {
				return nil, errors.New("Valeur max invalide pour uint8")
			}

			ppm.max = uint8(maxval)
			headerSize = counter
		}
		counter++
	}

	data := make([][]Pixel, ppm.height)

	for i := 0; i < ppm.height; i++ {
		data[i] = make([]Pixel, ppm.width)
	}

	var splitData []string

	if counter > headerSize {
		for i := 0; i < ppm.height; i++ {
			splitData = strings.Fields(splitFile[headerSize+1+i])
			for j := 0; j < ppm.width*3; j += 3 {
				r, _ := strconv.Atoi(splitData[j])
				g, _ := strconv.Atoi(splitData[j+1])
				b, _ := strconv.Atoi(splitData[j+2])
				data[i][j/3] = Pixel{R: uint8(r), G: uint8(g), B: uint8(b)}
			}
		}
	}
	return &PPM{data: data, width: ppm.width, height: ppm.height, magicNumber: ppm.magicNumber, max: uint8(ppm.max)}, nil
}

// Size returns the width and height of the image.
func (ppm *PPM) Size() (int, int) {
	return ppm.width, ppm.height
}

// At returns the value of the pixel at (x, y).
func (ppm *PPM) At(x, y int) Pixel {
	return ppm.data[y][x]
}

// Set sets the value of the pixel at (x, y).
func (ppm *PPM) Set(x, y int, value Pixel) {
	ppm.data[y][x] = value
}

// Invert inverts the colors of the PPM image.
func (ppm *PPM) Invert() {
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			ppm.data[y][x].R = ppm.max - ppm.data[y][x].R
			ppm.data[y][x].G = ppm.max - ppm.data[y][x].G
			ppm.data[y][x].B = ppm.max - ppm.data[y][x].B
		}
	}
}

// Flip flips the PPM image horizontally.
func (ppm *PPM) Flip() {
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width/2; x++ {
			ppm.data[y][x], ppm.data[y][ppm.width-x-1] = ppm.data[y][ppm.width-x-1], ppm.data[y][x]
		}
	}
}

// Flop flops the PPM image vertically.
func (ppm *PPM) Flop() {
	for y := 0; y < ppm.height/2; y++ {
		ppm.data[y], ppm.data[ppm.height-y-1] = ppm.data[ppm.height-y-1], ppm.data[y]
	}
}

// SetMagicNumber sets the magic number of the PPM image.
func (ppm *PPM) SetMagicNumber(magicNumber string) {
	ppm.magicNumber = magicNumber
}

// SetMaxValue sets the max value of the PPM image.
func (ppm *PPM) SetMaxValue(maxValue uint8) {
	ppm.max = uint8(maxValue)
}

// Rotate90CW rotates the PPM image 90° clockwise.
func (ppm *PPM) Rotate90CW() {
	rotated := make([][]Pixel, ppm.width)
	for x := 0; x < ppm.width; x++ {
		rotated[x] = make([]Pixel, ppm.height)
		for y := 0; y < ppm.height; y++ {
			rotated[x][y] = ppm.data[ppm.height-y-1][x]
		}
	}
	ppm.width, ppm.height = ppm.height, ppm.width
	ppm.data = rotated
}

// ToPBM converts the PPM image to PBM.
func (ppm *PPM) ToPBM() *PBM {
	pbm := &PBM{
		width:       ppm.width,
		height:      ppm.height,
		magicNumber: "P1",
		data:        make([][]bool, ppm.height),
	}

	// Convert color pixels to binary pixels
	for y := 0; y < ppm.height; y++ {
		pbm.data[y] = make([]bool, ppm.width)
		for x := 0; x < ppm.width; x++ {
			// Convert color pixel to binary pixel (true if color is not black)
			pbm.data[y][x] = ppm.data[y][x].R != 0 || ppm.data[y][x].G != 0 || ppm.data[y][x].B != 0
		}
	}

	return pbm
}

// ToPGM converts the PPM image to PGM.
func (ppm *PPM) ToPGM() *PGM {
	pgm := &PGM{
		width:       ppm.width,
		height:      ppm.height,
		magicNumber: "P2",
		max:         ppm.max,
		data:        make([][]uint8, ppm.height),
	}

	// Convert color pixels to grayscale pixels
	for y := 0; y < ppm.height; y++ {
		pgm.data[y] = make([]uint8, ppm.width)
		for x := 0; x < ppm.width; x++ {
			// Convert color pixel to grayscale pixel
			grayscale := (uint32(ppm.data[y][x].R) + uint32(ppm.data[y][x].G) + uint32(ppm.data[y][x].B)) / 3
			pgm.data[y][x] = uint8(grayscale)
		}
	}

	return pgm
}
func display(data [][]Pixel) {
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[0]); j++ {
			fmt.Print(data[i][j], " ")
		}
		fmt.Println()
	}
}
