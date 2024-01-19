package Netpbm

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           uint8
}

// ReadPGM reads a file in PGM format and returns a struct representing the image
func ReadPGM(filename string) (*PGM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var width, height, max int
	var data [][]uint8

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	magicNumber := scanner.Text()
	if magicNumber != "P2" && magicNumber != "P5" {
		return nil, errors.New("type de fichier non pris en charge")
	}
	//Read the dimensiosn of the image
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") { // skips comments
			_, err := fmt.Sscanf(line, "%d %d", &width, &height)
			if err == nil {
				break
			} else {
				fmt.Println("Largeur ou hauteur invalide :", err)
			}
		}
	}
	//read a maximum pixel value
	scanner.Scan()
	max, err = strconv.Atoi(scanner.Text())
	if err != nil {
		return nil, errors.New("valeur maximale de pixel invalide")
	}
	//Read a pixel value and make a data matrix
	for scanner.Scan() {
		line := scanner.Text()
		if magicNumber == "P2" {
			row := make([]uint8, 0)
			for _, char := range strings.Fields(line) {
				pixel, err := strconv.Atoi(char)
				if err != nil {
					fmt.Println("Erreur de conversion en entier :", err)
				}
				if pixel >= 0 && pixel <= max {
					row = append(row, uint8(pixel))
				} else {
					fmt.Println("Valeur de pixel invalide :", pixel)
				}
			}
			data = append(data, row)
		}
	}
	return &PGM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
		max:         uint8(max),
	}, nil
}

// Size returns the width and height of the image
func (pgm *PGM) Size() (int, int) {
	return pgm.width, pgm.height
}

// At returns the value of the pixel at (x, y)
func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[y][x]
}

// Set sets the value of the pixel at (x, y).
func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[y][x] = value
}

// Save saves the PGM image to a file and returns an error if there was a problem.
// Save saves the PGM image to a file and returns an error if there was a problem.
func (pgm *PGM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%s\n", pgm.magicNumber)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(file, "%d %d\n", pgm.width, pgm.height)
	if err != nil {
		return err
	}
	// Write pixel values to the file for P2 format
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			_, err = fmt.Fprintf(file, "%d ", pgm.data[i][j])
			if err != nil {
				return err
			}
		}
		_, err = fmt.Fprintln(file) // Use Fprintln to write a newline
		if err != nil {
			return err
		}
	}

	return nil
}

// Invert inverts the colors of the PBM image.
func (pgm *PGM) Invert() {
	newData := make([][]uint8, len(pgm.data)) //Create a new matrix  to hold the inverted
	for i := range pgm.data {
		newData[i] = make([]uint8, len(pgm.data[i])) //create the same largeur
		for j := range pgm.data[i] {
			newData[i][j] = uint8(pgm.max) - pgm.data[i][j] //retirer le max to invert the pixel
		}
	}
	pgm.data = newData //update the originam data
}

// Flip flips the PGM image horizontally.
func (pgm *PGM) Flip() {
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width/2; j++ {
			pgm.data[i][j], pgm.data[i][pgm.width-1-j] = pgm.data[i][pgm.width-1-j], pgm.data[i][j]
		}
	}
}

// Flop inverts the PGM image vertically by creating a new matrix.
func (pgm *PGM) Flop() {
	newData := make([][]uint8, len(pgm.data))
	for i := range pgm.data {
		newData[len(pgm.data)-i-1] = make([]uint8, len(pgm.data[i]))
		copy(newData[len(pgm.data)-i-1], pgm.data[i])
	}
	pgm.data = newData
}

// SetMagicNumber sets the magic number of the PGM image.
func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

// SetMaxValue sets the max value of the PPM image.
func (pgm *PGM) SetMaxValue(maxValue uint8) {
	pgm.max = uint8(maxValue)
}

// ////////////////
// Rotate90CW rotates the PGM image 90° clockwise.
func (pgm *PGM) Rotate90CW() {
	NumRows := pgm.width
	NumColums := pgm.height
	for i := 0; i < len(pgm.data); i++ {
		var temp uint8
		for j := i + 1; j < len(pgm.data[0]); j++ {
			temp = pgm.data[i][j]
			pgm.data[i][j] = pgm.data[j][i]
			pgm.data[j][i] = temp
		}
	}
	for i := 0; i < NumColums; i++ {
		for j := 0; j < NumRows/2; j++ {
			temp := pgm.data[i][j]
			pgm.data[i][j] = pgm.data[i][NumRows-j-1]
			pgm.data[i][NumRows-j-1] = temp
		}
	}
}

// /////////
// ToPBM converts the PGM image to PBM.
func (pgm *PGM) ToPBM() *PBM {
	pbm := &PBM{
		magicNumber: "P1",
		width:       pgm.width,
		height:      pgm.height,
		data:        make([][]bool, pgm.height),
	}

	for i := range pbm.data {
		pbm.data[i] = make([]bool, pgm.width)
		for j := 0; j < pgm.width; j++ {
			// Convert to pixel value
			pbm.data[i][j] = pgm.data[i][j] > pgm.max/2
		}
	}

	return pbm
}
