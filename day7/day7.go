package main

import (
	"bufio"
	"fmt"
	"iter"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Part 1
	fmt.Println("---- Testing * and + operators")
	result1 := getTotalCalibrationResult(scanner, []rune{'*', '+'})

	file.Seek(0, 0) // Re-read the file

	// Part 2
	fmt.Println("---- Testing |, *, and + operators")
	result2 := getTotalCalibrationResult(scanner, []rune{'|', '*', '+'})

	fmt.Println("The calibration result with * and + is", result1)
	fmt.Println("The calibration result with |, *, and + is", result2)
}

func getTotalCalibrationResult(scanner *bufio.Scanner, operators []rune) int {
	calibrationResult := 0

	for scanner.Scan() {
		targetNumber, numbers := parseLine(scanner.Text())

		fmt.Println("Testing target", targetNumber, "with numbers", numbers)
		for operations := range SequenceOperations(len(numbers)-1, operators) {
			value := numbers[0]
			for i, op := range operations {
				switch op {
				case '|':
					var err error
					// I was pleasantly surprised to (empirically) find out that this is performant enough
					value, err = strconv.Atoi(strconv.Itoa(value) + strconv.Itoa(numbers[i+1]))
					if err != nil { // We shouldn't have any errors, but you can never be too sure
						panic(err)
					}
				case '*':
					value *= numbers[i+1]
				case '+':
					value += numbers[i+1]
				}

				if value > targetNumber {
					break
				}
			}

			if targetNumber == value {
				calibrationResult += targetNumber
				fmt.Println("Test is valid with the following operations:", strings.Split(string(operations), ""))
				break
			}
		}
	}

	return calibrationResult
}

func parseLine(line string) (int, []int) {
	values := strings.SplitN(line, ": ", 2)
	targetNumber, _ := strconv.Atoi(values[0])
	numbersStr := strings.Split(values[1], " ")

	numbers := make([]int, 0, len(numbersStr))
	for _, str := range numbersStr {
		number, _ := strconv.Atoi(str)
		numbers = append(numbers, number)
	}

	return targetNumber, numbers
}

/*
Returns an iterator that yields an array with `nOps` operations (each one represented by a rune).
The iterator loops through every valid `nOps`-tuple of operators, in a predictable order.
*/
func SequenceOperations(nOps int, operators []rune) iter.Seq[[]rune] {
	mapToOperators := func(indexes []int) []rune {
		result := make([]rune, 0, len(indexes))
		for _, n := range indexes {
			result = append(result, operators[n])
		}
		return result
	}

	return func(yield func([]rune) bool) {
		if nOps <= 0 {
			return
		}

		indexes := make([]int, nOps)
		for yield(mapToOperators(indexes)) && incrementArray(&indexes, len(operators), nOps-1) {
		}
	}
}

func incrementArray(arr *[]int, end int, idx int) bool {
	if idx < 0 {
		return false
	}

	if (*arr)[idx] < end-1 {
		(*arr)[idx]++
		return true
	}

	for i := idx; i < len(*arr); i++ {
		(*arr)[i] = 0
	}
	return incrementArray(arr, end, idx-1)
}
