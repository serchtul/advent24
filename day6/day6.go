package main

import (
	"bufio"
	"fmt"
	"log"
	"maps"
	"os"
	"slices"
)

type coord struct {
	x int
	y int
}

type UnitVector struct {
	location  coord
	direction coord
}

var existsStruct struct{}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var initialLocation coord
	obstacles := make([][]bool, 0)

	scanner := bufio.NewScanner(file)
	y := 0
	for scanner.Scan() {
		line := scanner.Text()

		obstacles = append(obstacles, make([]bool, len(line)))

		for x, c := range line {
			switch c {
			case '#':
				obstacles[y][x] = true
			case '^':
				initialLocation = coord{x, y}
			}
		}
		y++
	}

	path := make([]UnitVector, 0)

	// Calculate the guard's original path
	_, visitedMap := calculatePath(initialLocation, coord{0, -1}, obstacles, &path)

	// Calculate all alternate paths (I couldn't think of anything more clever than to brute force all obstacles that the guard could have found along the original path)
	cycles := 0
	tested := make(map[coord]struct{})
	for i := range path {
		obstacles := cloneObstacles(obstacles)

		newObstacle := getObstacleLocation(path[i].location, path[i].direction, obstacles)
		if _, alreadyTested := tested[newObstacle]; alreadyTested || newObstacle.x == -1 && newObstacle.y == -1 {
			fmt.Println("Skiping test for obstacle, location", path[i].location, "direction", path[i].direction)
			continue
		}
		obstacles[newObstacle.y][newObstacle.x] = true

		location := path[0].location
		direction := path[0].direction
		path := make([]UnitVector, 0) // We need to re-calculate the path, to ensure the new obstacle doesn't modify previous parts of it

		fmt.Println("Testing alternate path for obstacle at", newObstacle)
		hasCycle, _ := calculatePath(location, direction, obstacles, &path)
		if hasCycle {
			fmt.Println("!!! Found cycle")
			cycles++
		}

		tested[newObstacle] = existsStruct
		// writeToFile("out/day6_"+strconv.Itoa(i)+"_"+strconv.FormatBool(hasCycle)+".out", serializeMap(initialLocation, visitedMap, path, obstacles, newObstacle))
	}

	fmt.Println("Visited", len(visitedMap), "different positions, found", cycles, "cycles")
	// writeToFile("out/day6.out", serializeMap(initialLocation, visitedMap, path, obstacles, coord{-1, -1}))
}

// Returns a boolean specifying whether the path contains a cycle and a map of all the visited locations.
func calculatePath(initialLocation coord, initialDirection coord, obstacles [][]bool, path *[]UnitVector) (hasCycle bool, visitedMap map[coord]map[coord]struct{}) {
	location := initialLocation
	direction := initialDirection

	visitedMap = make(map[coord]map[coord]struct{})
	for !hasCycle && !isOutOfBounds(location, len(obstacles), len(obstacles[0])) {
		location, direction, hasCycle = visitLocation(location, direction, obstacles, path, visitedMap)
	}

	return
}

func visitLocation(location coord, direction coord, obstacles [][]bool, path *[]UnitVector, visited map[coord]map[coord]struct{}) (nextLocation coord, nextDirection coord, hasCycle bool) {
	visitedDirs, exists := visited[location]
	if !exists {
		visitedDirs = make(map[coord]struct{})
		visited[location] = visitedDirs
	}
	_, hasCycle = visitedDirs[direction]

	visitedDirs[direction] = existsStruct
	*path = append(*path, UnitVector{location, direction})

	nextLocation = coord{location.x + direction.x, location.y + direction.y}
	// Make sure we don't collide with obstacles
	if !isOutOfBounds(nextLocation, len(obstacles), len(obstacles[0])) && obstacles[nextLocation.y][nextLocation.x] {
		// fmt.Println("Obstacle at location", location, "& direction", direction)
		nextDirection = turnRight90Deg(direction)
		nextLocation = location // Don't move, just rotate
	} else {
		nextDirection = direction
	}

	// fmt.Println("Next location", nextLocation, " direction", nextDirection, "cycle", hasCycle)
	return
}

func isOutOfBounds(location coord, sizeX int, sizeY int) bool {
	return location.x < 0 || location.x >= sizeX || location.y < 0 || location.y >= sizeY
}

// This is the result of multiplying the coordinates vector with the corresponding rotation matrix
func turnRight90Deg(current coord) coord {
	return coord{-current.y, current.x}
}

/*
Given a location and direction, returns where the obstacle would be placed for the guard to turn right.
Returns coord{-1,-1} if the location can't exist (out of bounds or an obstacle already exists there)
*/
func getObstacleLocation(location coord, direction coord, obstacles [][]bool) coord {
	obstacleLocation := coord{location.x + direction.x, location.y + direction.y}

	if isOutOfBounds(obstacleLocation, len(obstacles), len(obstacles[0])) || obstacles[obstacleLocation.y][obstacleLocation.x] {
		return coord{-1, -1}
	}
	return obstacleLocation
}

func cloneObstacles(obstacles [][]bool) [][]bool {
	newObstacles := make([][]bool, len(obstacles))
	for i := range obstacles {
		newObstacles[i] = slices.Clone(obstacles[i])
	}
	return newObstacles
}

func serializeMap(initialLocation coord, visitedMap map[coord]map[coord]struct{}, path []UnitVector, obstacles [][]bool, obstacleLoc coord) []string {
	visited := make([][]rune, len(obstacles))
	for y := range obstacles {
		visited[y] = slices.Repeat([]rune{'.'}, len(obstacles[0]))
		for x, obstacle := range obstacles[y] {
			if obstacle {
				visited[y][x] = '#'
			}
		}
	}
	for _, vector := range path {
		loc := vector.location
		visited[loc.y][loc.x] = getDirectionRune(visitedMap[loc])
	}

	visited[initialLocation.y][initialLocation.x] = '^'
	if !isOutOfBounds(obstacleLoc, len(obstacles), len(obstacles[0])) {
		visited[obstacleLoc.y][obstacleLoc.x] = 'O'
	}

	res := make([]string, 0, len(obstacles))
	for _, row := range visited {
		res = append(res, string(row))
	}
	return res
}

func getDirectionRune(directionMap map[coord]struct{}) rune {
	visitedDirs := slices.Collect(maps.Keys(directionMap))

	result := '.'
	switch len(visitedDirs) {
	case 4, 3:
		result = '+'
	case 1:
		if visitedDirs[0].x != 0 {
			result = '-'
		} else {
			result = '|'
		}
	case 2:
		switch {
		case slices.Contains(visitedDirs, coord{1, 0}) && slices.Contains(visitedDirs, coord{-1, 0}):
			result = '-'
		case slices.Contains(visitedDirs, coord{0, 1}) && slices.Contains(visitedDirs, coord{0, -1}):
			result = '|'
		default:
			result = '+'
		}
	}
	return result
}

func writeToFile(fileName string, rows []string) {
	outF, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer outF.Close()

	for _, row := range rows {
		if _, err := outF.WriteString(string(row) + "\n"); err != nil {
			log.Fatal(err)
		}
	}
	if err := outF.Sync(); err != nil {
		log.Fatal(err)
	}
}
