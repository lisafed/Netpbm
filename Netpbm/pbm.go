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
	var PBMPfour PBM // New variable for P4 format

	// Iterate through the lines of the file one by one
	for scanner.Scan() {
		line = scanner.Text()

		// Skip comments
		if strings.HasPrefix(line, "#") {
			continue
		}

		// The magic number and image dimensions
		if pbm.magicNumber == "" {
			if strings.HasPrefix(line, "P") && (line[1] == '1' || line[1] == '4') {
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

			// If it's P4 format, initialize the second PBM struct
			if pbm.magicNumber == "P4" {
				PBMPfour = PBM{
					data:        make([][]bool, pbm.height),
					width:       pbm.width,
					height:      pbm.height,
					magicNumber: "P4",
				}
				for i := range PBMPfour.data {
					PBMPfour.data[i] = make([]bool, PBMPfour.width)
				}
			}

			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error reading file: %v", err)
	}

	// Read the pixel and update a data matrix
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
		} else if pbm.magicNumber == "P4" {
			if pbm.magicNumber == "P4" {
				var char int32
				var rowdata []bool
				for j := 0; j < pbm.width*2; j++ {
					if len(rowdata) == pbm.width {
						pbm.data[i] = rowdata
						i++
						for l := 0; l < len(rowdata); l++ {
							rowdata = []bool{}
						}
					}
					fmt.Printf("j %d\n", j)
					char = int32(line[j])
					fmt.Printf("char %d\n", char)
					binary := fmt.Sprintf("%08b", char)
					fmt.Println(boolget(binary))
					boleanlist := boolget(binary)
					for k := 0; k < len(binary); k++ {
						rowdata = append(rowdata, boleanlist[k])
						if len(rowdata) == pbm.width {
							k = len(binary)
						}
					}
					fmt.Println(rowdata)

				}
				if len(rowdata) == pbm.width {
					pbm.data[i] = rowdata
					i++
					for l := 0; l < len(rowdata); l++ {
						rowdata = []bool{}
					}
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
func boolget(binaryString string) []bool {
	var boolData []bool
	for _, char := range binaryString {
		if char == '1' {
			boolData = append(boolData, true)
		} else if char == '0' {
			boolData = append(boolData, false)
		}
	}
	//fmt.Println(boolData)
	return boolData
}

// Set sets the value of the pixel at (x, y).
func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[y][x] = value //true=1,false=0
}

// Save saves the PBM image to a file and returns an error if there was a problem.
func (pbm *PBM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	newFile := bufio.NewWriter(file)

	// Write the magic number and the dimensions to the file
	fmt.Fprintf(newFile, "%s\n", pbm.magicNumber)
	fmt.Fprintf(newFile, "%d %d\n", pbm.width, pbm.height)

	// Write the pixel data to the file
	for _, row := range pbm.data {
		for _, pixel := range row {
			if pixel {
				fmt.Fprintf(newFile, "1 ")
			} else {
				fmt.Fprintf(newFile, "0 ")
			}
		}
		fmt.Fprintln(newFile) // Move to the next line after each row
	}
	newFile.Flush()
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
