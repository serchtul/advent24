package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type coord struct {
	x int
	y int
}

func (a coord) Add(b coord) coord {
	return coord{a.x + b.x, a.y + b.y}
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	topography := [][]int{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		numberStrings := strings.Split(line, "")
		numbers := make([]int, len(line))

		for i, num := range numberStrings {
			numbers[i], _ = strconv.Atoi(num) // Assume input file's format will always be correct, hence no errors
		}
		topography = append(topography, numbers)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	scoresSum := 0
	rankingsSum := 0
	for y, row := range topography {
		for x, height := range row {
			if height == 0 {
				trailhead := coord{x, y}
				paths := getPathsToPeak(trailhead, trailhead, topography, nil)

				score := getTrailheadScore(paths)
				ranking := len(paths)

				scoresSum += score
				rankingsSum += ranking
				fmt.Println("Found trailhead at", trailhead, "with score", score, "and ranking", ranking)
			}
		}
	}

	fmt.Println("Map's sum of scores is", scoresSum)
	fmt.Println("Map's sum of rankings is", rankingsSum)
}

func getTrailheadScore(paths [][]coord) int {
	uniquePeaks := map[coord]struct{}{}

	for _, path := range paths {
		uniquePeaks[path[len(path)-1]] = struct{}{}
	}

	return len(uniquePeaks)
}

func getPathsToPeak(here coord, prev coord, topography [][]int, path []coord) [][]coord {
	if isOutOfBounds(here, topography) {
		return nil // We're out of the grid before reaching the top
	}
	if slices.Contains(path, here) {
		return nil // We've been here before
	}

	currValue := topography[here.y][here.x]
	// fmt.Println("curr", currValue, "prev", topography[prev.y][prev.x], "here", here)
	if currValue != topography[prev.y][prev.x]+1 && here != prev {
		return nil // We didn't go up exactly one
	}

	path = append(path, here)
	if currValue == 9 {
		fmt.Println("Found path to peak", here, path)
		return [][]coord{path}
	}

	res := getPathsToPeak(here.Add(coord{-1, 0}), here, topography, slices.Clone(path))
	res = append(res, getPathsToPeak(here.Add(coord{1, 0}), here, topography, slices.Clone(path))...)
	res = append(res, getPathsToPeak(here.Add(coord{0, 1}), here, topography, slices.Clone(path))...)
	res = append(res, getPathsToPeak(here.Add(coord{0, -1}), here, topography, slices.Clone(path))...)

	return res
}

func isOutOfBounds(location coord, topography [][]int) bool {
	return location.x < 0 || location.y < 0 || location.x >= len(topography[0]) || location.y >= len(topography)
}
