package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var existsStruct struct{}

func main() {
	file, err := os.Open("input.txt")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Adjancency list of page dependencies
	sequencing := make(map[string]map[string]struct{}) // Using strings to avoid parsing every number

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		pages := strings.Split(line, "|")

		// We stop whenever we encounter an empty line
		if len(pages) == 1 {
			break
		}

		dependents, exists := sequencing[pages[0]]
		if !exists {
			dependents = make(map[string]struct{})
			sequencing[pages[0]] = dependents
		}
		dependents[pages[1]] = existsStruct
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Requirements", sequencing)

	result := 0
	for scanner.Scan() {
		line := scanner.Text()
		updates := strings.Split(line, ",")
		seenPages := make(map[string]struct{})

		validOrder := true
		for _, pageNo := range updates {
			fmt.Println("seen", seenPages, "current page", pageNo)

			for seen := range seenPages {
				// The sequence is not in the correct order, because we saw a page we should have seen *after* this one
				if _, exists := sequencing[pageNo][seen]; exists {
					fmt.Println("Invalid sequence order", line)
					validOrder = false
					break
				}
			}
			seenPages[pageNo] = existsStruct
			if !validOrder {
				break
			}
		}

		if validOrder {
			fmt.Println("Valid sequence order", line)
			// Assumes the number will correctly be parsed.
			// This also assumes the number of updates is always odd
			middle, _ := strconv.Atoi(updates[(len(updates)-1)/2])
			result += middle
		}
	}

	fmt.Println("Result is", result)
}
