package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type Trie struct {
	isLeaf   bool
	children map[int]Trie
	next     []int
}

func main() {
	// input := "77 515 6779622 6 91370 959685 0 9861"
	input := "125 17"
	pebblesStr := strings.Split(input, " ")
	pebbles := make([]int, len(pebblesStr))
	for i, s := range pebblesStr {
		pebbles[i], _ = strconv.Atoi(s)
	}
	blinks := 15
	rootNode := Trie{children: map[int]Trie{}}
	rootNode.children[0] = Trie{isLeaf: true, next: []int{1}}

	for b := range blinks {
		for i := 0; i < len(pebbles); i++ {
			if pebbles[i] == 0 {
				pebbles[i] = 1
				pebblesStr[i] = "1"
			} else if len(pebblesStr[i])%2 == 0 {
				pStr := pebblesStr[i]
				newPStr := pStr[len(pStr)/2:]
				oldPStr := pStr[:len(pStr)/2]

				oldP, _ := strconv.Atoi(oldPStr)
				pebblesStr[i] = strconv.Itoa(oldP)
				pebbles[i] = oldP

				newP, _ := strconv.Atoi(newPStr)
				pebbles = slices.Insert(pebbles, i+1, newP)
				pebblesStr = slices.Insert(pebblesStr, i+1, strconv.Itoa(newP))
				i++
			} else {
				pebbles[i] *= 2024
				pebblesStr[i] = strconv.Itoa(pebbles[i])
			}
		}
		fmt.Println("After", b+1, "blink(s), stones", len(pebbles), pebbles)
	}
}
