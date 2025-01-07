package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	allMulOps := readAllMul(file)
	fmt.Println("Result of all mul operations is", multiplyNumbers(allMulOps))

	file.Seek(0, 0)

	doOnlyMulOps := readDoOnlyMul(file)
	fmt.Println("Result of do-only mul operations is", multiplyNumbers(doOnlyMulOps))
}

// Part 1
func readAllMul(file *os.File) []string {
	scanner := bufio.NewScanner(file)

	mulOp := regexp.MustCompile("mul\\((\\d{1,3},\\d{1,3})\\)")
	scanner.Split(func(data []byte, atEOF bool) (int, []byte, error) {
		found := mulOp.FindSubmatchIndex(data)
		if found == nil {
			return 0, nil, nil
		}
		return found[1], data[found[2]:found[3]], nil
	})

	result := make([]string, 0, 1000)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return result
}

// Part 2
func readDoOnlyMul(file *os.File) []string {
	scanner := bufio.NewScanner(file)

	validTokens := regexp.MustCompile("do\\(\\)|don't\\(\\)|mul\\(\\d{1,3},\\d{1,3}\\)")
	scanner.Split(func(data []byte, atEOF bool) (int, []byte, error) {
		found := validTokens.FindIndex(data)
		if found == nil {
			return 0, nil, nil
		}
		return found[1], data[found[0]:found[1]], nil
	})

	result := make([]string, 0, 1000)
	includeOp := true
	for scanner.Scan() {
		op := scanner.Text()
		if op == "do()" {
			includeOp = true
			continue
		} else if op == "don't()" {
			includeOp = false
		}

		if !includeOp {
			continue
		}

		numsStr, _ := strings.CutPrefix(op, "mul(")
		numsStr, _ = strings.CutSuffix(numsStr, ")")

		result = append(result, numsStr)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return result
}

// Expects an array of strings, each containing a pair of comma-separated numbers
func multiplyNumbers(operations []string) int {
	result := 0
	for _, numsStr := range operations {
		numbers := strings.Split(numsStr, ",")

		n1, _ := strconv.Atoi(numbers[0])
		n2, _ := strconv.Atoi(numbers[1])

		result += n1 * n2
		fmt.Println("included mul", numbers)
	}
	return result
}
