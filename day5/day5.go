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
		pages := strings.Split(scanner.Text(), "|")

		// The dependency list ends whenever we encounter an empty line
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
	fmt.Println("Sequencing requirements", sequencing)

	firstTryResults := 0
	reorderedResults := 0
	for scanner.Scan() {
		pageUpdates := strings.Split(scanner.Text(), ",")

		firstTry := true // Whether the sequence was valid on the first check
		var validOrder bool
		for {
			var idx int
			var page string
			var conflictingIdx int

			validOrder = true
			for idx, page = range pageUpdates {
				// Find whether there's a conflicting page (one that should be after the current one, instead of before)
				conflictingIdx = slices.IndexFunc(pageUpdates[:idx], func(p string) bool {
					_, conflictExists := sequencing[page][p]
					return conflictExists
				})

				if conflictingIdx != -1 {
					fmt.Println("failed at page", page, "| seen pages", pageUpdates[:idx])

					validOrder = false
					firstTry = false
					break
				}
			}

			if validOrder {
				fmt.Println("Valid sequence order", pageUpdates)
				break
			}

			fmt.Println("Invalid sequence order", pageUpdates)
			fmt.Println(page, "must be before", pageUpdates[conflictingIdx], "on index", conflictingIdx)

			// Re-arrange current item to be before its first dependency
			pageUpdates = slices.Delete(pageUpdates, idx, idx+1)
			pageUpdates = slices.Insert(pageUpdates, conflictingIdx, page)
			fmt.Println("rearranged conflicting update", pageUpdates)
		}

		// Assumes the number in the file will correctly be parsed.
		// This also assumes the number of updates is always odd
		middle, _ := strconv.Atoi(pageUpdates[(len(pageUpdates)-1)/2])
		if firstTry {
			firstTryResults += middle
		} else {
			reorderedResults += middle
		}
	}

	fmt.Println("First-try sequences sum", firstTryResults)
	fmt.Println("Reordered sequences sum", reorderedResults)
}
