package main

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"slices"

	"gonum.org/v1/gonum/stat/combin"
)

type vect struct {
	x int
	y int
}

func (a vect) Add(b vect) vect {
	return vect{a.x + b.x, a.y + b.y}
}

func (a vect) Diff(b vect) vect {
	return vect{a.x - b.x, a.y - b.y}
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	grid := buildGrid(bufio.NewScanner(file))

	freqMap := classifyFrequencies(grid)
	fmt.Println("Coordinates grouped by frequency", freqMap)

	gridSize := vect{len(grid[0]), len(grid)}

	// Part 1
	fmt.Println("==== Checking for antinodes", gridSize)
	antinodes := getAntinodes(freqMap, gridSize, false)

	// Part 2
	fmt.Println("==== Checking for antinodes, accounting for resonant harmonics")
	resonantAntinodes := getAntinodes(freqMap, gridSize, true)

	fmt.Println("Total antinodes", len(antinodes), antinodes)
	fmt.Println("Total antinodes accounting for resonant harmonics", len(resonantAntinodes), resonantAntinodes)
}

func buildGrid(scanner *bufio.Scanner) [][]rune {
	grid := make([][]rune, 0)
	// Input file is a rectangular grid
	for scanner.Scan() {
		line := scanner.Text()
		grid = append(grid, []rune(line))
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return grid
}

func classifyFrequencies(grid [][]rune) map[rune][]vect {
	result := make(map[rune][]vect)

	for y, row := range grid {
		for x, r := range row {
			if r != '.' {
				result[r] = append(result[r], vect{x, y})
			}
		}
	}

	return result
}

func getAntinodes(freqMap map[rune][]vect, gridSize vect, expand bool) []vect {
	antinodes := make(map[vect]struct{})

	for _, coords := range freqMap {
		for _, indexes := range combin.Combinations(len(coords), 2) {
			p1 := coords[indexes[0]]
			p2 := coords[indexes[1]]

			slope := p2.Diff(p1) // Assuming a positive run & raise to calculate the antinodes, but these equations hold without loss of generality
			fmt.Println("checking location pair", []vect{p1, p2}, "slope", slope)

			// Expand in the "upward" direction
			for antinode := p1.Diff(slope); !isOutOfBounds(antinode, gridSize); antinode = antinode.Diff(slope) {
				fmt.Println("Antinode at", antinode)

				antinodes[antinode] = struct{}{}
				if !expand {
					break
				}
			}

			// And then in the "downward" one
			for antinode := p2.Add(slope); !isOutOfBounds(antinode, gridSize); antinode = antinode.Add(slope) {
				fmt.Println("Antinode at", antinode)

				antinodes[antinode] = struct{}{}
				if !expand {
					break
				}
			}

			// When we expand for harmonics, the points are also antinodes
			if expand {
				antinodes[p1] = struct{}{}
				antinodes[p2] = struct{}{}
			}
		}
	}

	return slices.Collect(maps.Keys(antinodes))
}

func isOutOfBounds(p vect, gridSize vect) bool {
	// fmt.Println("Found tentative antinode", p)
	return p.x < 0 || p.y < 0 || p.x >= gridSize.x || p.y >= gridSize.y
}
