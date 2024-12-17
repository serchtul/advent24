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

const MinSafe = 1
const MaxSafe = 3

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	reports := make([][]int, 0, 1000) // Capacity _could_ be ommited
	for scanner.Scan() {
		report := strings.Fields(scanner.Text())

		levels := make([]int, len(report))
		for i, v := range report {
			// Assume numbers can always be parsed
			l, _ := strconv.Atoi(v)
			levels[i] = l
		}
		reports = append(reports, levels)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Safe Reports:", safeReports(reports))
	fmt.Println("Safe Dampened Reports:", safeDampenedReports(reports))
}

func safeReports(reports [][]int) int {
	safeReports := 0
	for _, levels := range reports {
		if safeDampenedReport(levels, true) {
			fmt.Println(levels, "is safe")
			safeReports++
		}
	}

	return safeReports
}

func safeDampenedReports(reports [][]int) int {
	safeReports := 0
	for _, levels := range reports {
		if safeDampenedReport(levels, false) {
			fmt.Println(levels, "is dampened-safe")
			safeReports++
		}
	}

	return safeReports
}

func safeDampenedReport(levels []int, dampened bool) bool {
	safeReport := true

	increasing := levels[1]-levels[0] > 0
	for i := 0; i < len(levels)-1 && safeReport; {
		diff := levels[i+1] - levels[i]
		distance := abs(diff)

		if diff > 0 != increasing || distance < MinSafe || distance > MaxSafe {
			if !dampened {
				fmt.Println(levels, "dampened", "i", i, "diff", diff, "increasing", increasing)

				// There's probably a better way to solve this part
				canRecover := false
				for j := max(i-1, 0); j < min(len(levels), i+2) && !canRecover; j++ {
					canRecover = safeDampenedReport(slices.Delete(slices.Clone(levels), j, j+1), true)
				}
				return canRecover
			}
			safeReport = false
			fmt.Println(levels, "is unsafe.", "i", i, "diff", diff, "increasing", increasing)
		}
		i++
	}

	return safeReport
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
