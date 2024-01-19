package Netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type PBM struct {
	data          [][]bool //Matrix of pixels
	width, height int
	magicNumber   string //P1 or P4 (BPM formats)
}

// ReadPBM reads a file in PBM format and returns a struct representing the image.
func ReadPBM(filename string) (*PBM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var pbm PBM
	var line string
	//iterate the line of the file one by one
	for scanner.Scan() {
		line = scanner.Text()
		//Skip comments
		if string(line[0]) == "#" {
			continue
		}
		//the magic number and image dimensions
		if pbm.magicNumber == "" {
			if string(line[0]) == "P" && (line[1] == '1' || line[1] == '4') {
				pbm.magicNumber = line
			} else {
				return nil, fmt.Errorf("Invalid magic number: %s", line)
			}
		} else if pbm.width == 0 || pbm.height == 0 {
			var width, height int
			// Use "fmt.Sscanf" to read the dimensions.
			if _, err := fmt.Sscanf(line, "%d %d", &width, &height); err != nil {
				return nil, fmt.Errorf("Error converting dimensions: %v", err)
			}
			pbm.width, pbm.height = width, height
			// Initialize the data matrix.
			pbm.data = make([][]bool, pbm.height)
			for i := range pbm.data {
				pbm.data[i] = make([]bool, pbm.width)
			}
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error reading file: %v", err)
	}
	//Read the pixel and update a data matrix
	for i := 0; i < pbm.height; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("Unexpected end of file")
		}
		line := scanner.Text()
		if pbm.magicNumber == "P1" {
			fieldLine := strings.Fields(line)
			if len(fieldLine) != pbm.width {
				return nil, fmt.Errorf("Invalid number of pixels in the line")
			}
			for j := 0; j < pbm.width; j++ {
				if fieldLine[j] == "1" {
					pbm.data[i][j] = true
				} else if fieldLine[j] == "0" {
					pbm.data[i][j] = false
				}
			}
		}
	}

	return &pbm, nil
}

// Size returns the width and height of the PBM image.
func (pbm *PBM) Size() (int, int) {
	return pbm.width, pbm.height
}

// At returns the value of the pixel at (x, y).
func (pbm *PBM) At(x, y int) bool {
	if y < 0 || y >= len(pbm.data) || x < 0 || x >= len(pbm.data[y]) {
		return false
	}
	return pbm.data[y][x]
}

// Set sets the value of the pixel at (x, y).
func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[y][x] = value //true=1,false=0
}

// Save saves the PBM image to a file and returns an error if there was a problem.
func (pbm *PBM) Save(filename string) error {
	file, err := os.Create(filename) //create a new file
	if err != nil {
		return err // return a error if we found one when we create  the file
	}
	//Write a magic number and the dimensions in the file
	_, err = file.WriteString(pbm.magicNumber + "\n") //Write tha magic number end return to the next line
	_, err = fmt.Fprintf(file, "%d %d\n", pbm.width, pbm.height)
	if err != nil {
		return err
	}
	//P1 format
	if pbm.magicNumber == "P1" {
		for _, list := range pbm.data {
			for _, pixel := range list {
				if pixel == true {
					_, err = file.WriteString("1")
				} else {
					_, err = file.WriteString("0")
				}
				if err != nil {
					return err // return an error if we encounter one when writing to the file
				}
			}
		}
	} else if pbm.magicNumber == "P4" {
	} else {
		return fmt.Errorf("Unsupported magic number: %s", pbm.magicNumber)
	}
	return nil
}

// Invert inverts the colors of the PBM image.
func (pbm *PBM) Invert() {
	newData := make([][]bool, len(pbm.data)) //Create a new matrix  to hold the inverted
	for i := range pbm.data {
		newData[i] = make([]bool, len(pbm.data[i])) //create the same Width
		for j := range pbm.data[i] {
			newData[i][j] = !pbm.data[i][j] //to invert the pixel
		}
	}
	pbm.data = newData //update the originam data
}

// Flip inverts the PBM image horizontally
func (pbm *PBM) Flip() {
	newData := make([][]bool, len(pbm.data))
	for i := range newData {
		newData[i] = make([]bool, len(pbm.data[i]))
		for j := 0; j < len(pbm.data[i]); j++ {
			newData[i][j] = pbm.data[i][len(pbm.data[i])-1-j]
		}
	}
	pbm.data = newData
}

// /////////////////////////////////////////////////////

// Flop inverts the PBM image vertically.
func (pbm *PBM) Flop() {
	newData := make([][]bool, len(pbm.data))
	for i := range pbm.data {
		newData[len(pbm.data)-i-1] = make([]bool, len(pbm.data[i]))
		copy(newData[len(pbm.data)-i-1], pbm.data[i])
	}
	pbm.data = newData
}

// /////////////////////////////////////////////////////////////////////////

// SetMagicNumber sets the magic number of the PBM image.
func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
	return

}

/*func main() {
	// Test de la fonction ReadPBM
	testPBM, readErr := ReadPBM("example.pbm")
	if readErr != nil {
		fmt.Printf("ReadPBM Test: Error reading file: %v\n", readErr)
		return
	}

	// Test de la fonction Size
	width, height := testPBM.Size()
	fmt.Printf("Size Test: Width=%d, Height=%d\n", width, height)

	// Test de la fonction At
	value := testPBM.At(1, 1)
	fmt.Printf("At Test: Value at (1, 1) = %v\n", value)

	// Test de la fonction Set
	testPBM.Set(0, 0, false)
	fmt.Println("Set Test: Updated data matrix:")
	for i := 0; i < testPBM.height; i++ {
		fmt.Println(testPBM.data[i])
	}

	// Test de la fonction Save
	saveErr := testPBM.Save("test_output.pbm")
	if saveErr != nil {
		fmt.Printf("Save Test: Error saving file: %v\n", saveErr)
	} else {
		fmt.Println("Save Test: File saved successfully.")
	}

	// Test de la fonction Invert
	testPBM.Invert()
	fmt.Println("Invert Test: Inverted data matrix:")
	for i := 0; i < testPBM.height; i++ {
		fmt.Println(testPBM.data[i])
	}

	// Test de la fonction Flip
	testPBM.Flip()
	fmt.Println("Flip Test: Flipped data matrix:")
	for i := 0; i < testPBM.height; i++ {
		fmt.Println(testPBM.data[i])
	}

	// Test de la fonction Flop
	testPBM.Flop()
	fmt.Println("Flop Test: Flopped data matrix:")
	for i := 0; i < testPBM.height; i++ {
		fmt.Println(testPBM.data[i])
	}

	// Test de la fonction SetMagicNumber
	testPBM.SetMagicNumber("P4")
	fmt.Printf("SetMagicNumber Test: Magic Number=%s\n", testPBM.magicNumber)

}*/
