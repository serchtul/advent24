package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type direction struct {
	x int
	y int
}

func main() {
	file, err := os.Open("input.txt")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	matrix := buildMatrix(file)

	xmasSequence := []rune{'X', 'M', 'A', 'S'}
	xmasDirections := [][]direction{
		{{0, 0}, {1, 0}, {1, 0}, {1, 0}},
		{{0, 0}, {0, 1}, {0, 1}, {0, 1}},
		{{0, 0}, {-1, 0}, {-1, 0}, {-1, 0}},
		{{0, 0}, {0, -1}, {0, -1}, {0, -1}},
		{{0, 0}, {1, 1}, {1, 1}, {1, 1}},
		{{0, 0}, {1, -1}, {1, -1}, {1, -1}},
		{{0, 0}, {-1, -1}, {-1, -1}, {-1, -1}},
		{{0, 0}, {-1, 1}, {-1, 1}, {-1, 1}},
	}
	fmt.Println("Total # XMAS:", countOccurrences(matrix, xmasSequence, xmasDirections))

	xMasSequence := []rune{'M', 'A', 'S', 'M', 'S'}
	xMasDirections := [][]direction{
		{{0, 0}, {1, 1}, {1, 1}, {-2, 0}, {2, -2}},
		{{0, 0}, {1, 1}, {1, 1}, {0, -2}, {-2, 2}},
		{{0, 0}, {-1, -1}, {-1, -1}, {2, 0}, {-2, 2}},
		{{0, 0}, {-1, -1}, {-1, -1}, {0, 2}, {2, -2}},
	}
	fmt.Println("Total # X-MAS:", countOccurrences(matrix, xMasSequence, xMasDirections))
}

func buildMatrix(file *os.File) [][]rune {
	matrix := make([][]rune, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		row := make([]rune, 0)
		for _, r := range line {
			row = append(row, r)
		}

		matrix = append(matrix, row)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return matrix
}

// Assumes matrix is a rectangular two-dimensional array
func countOccurrences(matrix [][]rune, targetSequence []rune, directionsSet [][]direction) int {
	result := 0
	numCols := len(matrix[0])
	numRows := len(matrix)
	for i := 0; i < numRows; i++ {
		for j := 0; j < numCols; j++ {
			// Search for the word, starting from a given cell, trying all possible directions
			for _, directions := range directionsSet {
				iScan := i
				jScan := j

				var k int
				for k = 0; k < len(directions); k++ {
					iScan += directions[k].y
					jScan += directions[k].x
					// Stop if the next direction is out of bounds
					if iScan < 0 || iScan >= numRows || jScan < 0 || jScan >= numCols {
						break
					}
					if matrix[iScan][jScan] != targetSequence[k] {
						break
					} else {
						// For debugging
						// fmt.Println("match!", iScan, jScan, string(matrix[iScan][jScan]), string(targetSequence[k]), k)
					}
				}
				// If we successfully scanned the whole target sequence, then we have a match
				if k == len(directions) {
					result++
				}
			}

		}
	}

	return result
}
