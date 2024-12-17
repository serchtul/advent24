package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	list1 := make([]int, 0, 1000) // Capacity _could_ be ommited
	list2 := make([]int, 0, 1000)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		// Assume numbers can always be parsed
		n1, _ := strconv.Atoi(line[0])
		n2, _ := strconv.Atoi(line[1])

		pos1, _ := slices.BinarySearch(list1, n1)
		list1 = slices.Insert(list1, pos1, n1)

		pos2, _ := slices.BinarySearch(list2, n2)
		list2 = slices.Insert(list2, pos2, n2)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Total distance:", distance(list1, list2))
	fmt.Println("Total similarity:", similarity(list1, list2))
}

// Lists must be the same size
func distance(list1 []int, list2 []int) int {
	var distance int = 0
	for i, v1 := range list1 {
		v2 := list2[i]

		d := abs(v1 - v2)

		distance += d
		// fmt.Println("line", i, "items:", v1, v2, "distance:", d, "total:", distance)
	}
	return distance
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// Lists must be the same size
func similarity(list1 []int, list2 []int) int {
	var similarity int = 0
	for _, val := range list1 {
		startIdx, existsInBoth := slices.BinarySearch(list2, val)

		if existsInBoth {
			endIdx := startIdx
			for list2[endIdx] == val {
				endIdx++
			}
			repetitions := int(endIdx - startIdx)
			similarity += val * repetitions
			// fmt.Println(val, "repetitions:", repetitions, "similarity:", similarity)
		}
	}
	return similarity
}
