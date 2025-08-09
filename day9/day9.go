package main

import (
	"fmt"
	"os"
	"slices"
)

func main() {
	fileContents, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	digits := []rune(string(fileContents))
	diskMap := make([]int, 0)
	fileSizes := make([]int, 0)
	spaceSizes := make([]int, 0)
	currId := 0
	for i, digit := range digits[:len(digits)-1] {
		r := -1 // This represents an empty space throughout this solution
		repetitions := int(digit - 0x30)
		if i%2 == 0 {
			r = currId
			fileSizes = append(fileSizes, r)
			currId++
		} else {
			spaceSizes = append(spaceSizes, repetitions)
		}

		diskMap = append(diskMap, slices.Repeat([]int{r}, repetitions)...)
	}

	fmt.Println("diskMap", diskMap)

	// Part 1
	fmt.Println("==== Compacting disk map by blocks")
	compacted1 := compactByBlock(slices.Clone(diskMap))
	fmt.Println("Checksum is", computeChecksum(compacted1))

	// Part 2
	fmt.Println("==== Compacting disk map by file")
	compacted2 := compactByFile(slices.Clone(diskMap))
	fmt.Println("Checksum is", computeChecksum(compacted2))
}

func compactByBlock(diskMap []int) []int {
	leftPtr := 0
	rightPtr := len(diskMap) - 1
	for {
		for ; diskMap[leftPtr] != -1; leftPtr++ {
		}
		for ; diskMap[rightPtr] == -1; rightPtr-- {
		}
		if leftPtr >= rightPtr {
			break
		}

		diskMap[leftPtr], diskMap[rightPtr] = diskMap[rightPtr], diskMap[leftPtr]
		// fmt.Println("diskMap", diskMap)
	}
	return diskMap
}

func compactByFile(diskMap []int) []int {
	startFilePtr := len(diskMap) - 1
	endFilePtr := startFilePtr

	var startSpacePtr, endSpacePtr int
	for endFilePtr > 0 {
		for ; diskMap[startFilePtr] == -1; startFilePtr-- {
		}
		for endFilePtr = startFilePtr - 1; endFilePtr > 0 && diskMap[endFilePtr] == diskMap[startFilePtr]; endFilePtr-- {
		}
		fileSize := startFilePtr - endFilePtr

		endSpacePtr = 0
		startSpacePtr = 0
		for endSpacePtr-startSpacePtr < fileSize && endSpacePtr < startFilePtr {
			if startSpacePtr != 0 {
				startSpacePtr++
			}
			for ; diskMap[startSpacePtr] != -1; startSpacePtr++ {
			}
			for endSpacePtr = startSpacePtr + 1; endSpacePtr < len(diskMap) && diskMap[endSpacePtr] == -1; endSpacePtr++ {
			}
		}
		// Ensure no end index is out of the expected bounds
		if endSpacePtr >= len(diskMap) || endSpacePtr >= startFilePtr || endFilePtr < 0 {
			// fmt.Println("file", diskMap[startFilePtr], "cannot be compacted")
			startFilePtr = endFilePtr // Skip file, start checking for the next one
			continue
		}

		fmt.Println("Compacting file", diskMap[startFilePtr], "to index", startSpacePtr)
		copy(diskMap[startSpacePtr:endSpacePtr], diskMap[endFilePtr+1:startFilePtr+1])
		copy(diskMap[endFilePtr+1:startFilePtr+1], slices.Repeat([]int{-1}, fileSize))
		// fmt.Println("diskMap", diskMap)
	}
	return diskMap
}

func computeChecksum(diskMap []int) int {
	checksum := 0
	for i, fileId := range diskMap {
		if fileId == -1 {
			continue
		}
		checksum += i * fileId
	}

	return checksum
}
