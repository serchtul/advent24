package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Children map[int]*Node

type Node struct {
	children  Children
	value     int
	nextValue []int
}

type PebbleCounter map[int]int

func main() {
	input := "77 515 6779622 6 91370 959685 0 9861"
	pebblesStr := strings.Split(input, " ")
	pebbles := make([]int, len(pebblesStr))
	for i, s := range pebblesStr {
		pebbles[i], _ = strconv.Atoi(s) // This assumes the input has the expected format
	}

	root := Node{children: Children{}}
	for i := range 10 {
		root.children[i] = &Node{value: i, children: Children{}}
	}
	for i := range 10 {
		buildBlinksTree(root.children[i], &root)
	}

	pebbleCount := PebbleCounter{}
	var newPebbleCount PebbleCounter

	// Use a pebble map instead of an array, to avoid calculating the next state for the same pebble multiple times
	for _, p := range pebbles {
		if _, exists := root.children[p]; !exists {
			root.children[p] = &Node{value: p, children: Children{}}
			// Pre-compute all possible next blinks for the given input
			buildBlinksTree(root.children[p], &root)
		}
		pebbleCount[p]++
	}

	blinks := 75 // This works for both Part 1 & Part 2
	for b := range blinks {
		newPebbleCount = PebbleCounter{}
		for p, count := range pebbleCount {
			for _, nextP := range root.children[p].nextValue {
				newPebbleCount[nextP] += count
			}
		}
		pebbleCount = newPebbleCount

		pebbles := 0
		for _, count := range pebbleCount {
			pebbles += count
		}

		fmt.Println("After", b+1, "blink(s), there are", pebbles, "pebbles")
	}
}

func buildBlinksTree(node *Node, root *Node) {
	next := blink(node.value)
	node.nextValue = next
	// fmt.Println("Current", node.value, "Next", node.nextValue)

	for _, p := range next {
		if _, nodeExists := root.children[p]; nodeExists {
			node.children[p] = root.children[p]
			// fmt.Println("End of path", n)
		} else {
			node.children[p] = &Node{value: p, children: Children{}}
			root.children[p] = node.children[p]
			buildBlinksTree(node.children[p], root)
		}
	}
}

func blink(n int) []int {
	if n == 0 {
		return []int{1}
	} else if nStr := strconv.Itoa(n); len(nStr)%2 == 0 {
		left := nStr[:len(nStr)/2]
		right := nStr[len(nStr)/2:]

		nLeft, _ := strconv.Atoi(left)
		nRight, _ := strconv.Atoi(right)

		return []int{nLeft, nRight}
	}

	return []int{n * 2024}
}
