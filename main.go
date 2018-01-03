package main

import (
	"bufio"
	"fmt"
	"gopkg.in/klaidliadon/next.v1"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

var n int

type Vertex struct {
	lat  float64
	long float64
}

func (v Vertex) String() string {
	return fmt.Sprintf("\nlat:\t%f\nlong:\t%f\n\n", v.lat, v.long)
}

var BitMap = make(map[int]uint64)
var SetMap = make(map[uint64][]int)
var VertexMap = make(map[int]*Vertex)
var SolutionsMap = make(map[uint64][]float64)

func main() {

	readFile(os.Args[1])

	// Quick look ups for binary representation of number
	for i := 1; i <= n; i++ {
		BitMap[i] = uint64(math.Pow(2, float64(i)))
	}

	// Identify which sets are of which size
	// which controls TSP subproblems
	setsBySize := make([][]uint64, n)

	for i := range setsBySize {
		setsBySize[i] = make([]uint64, 0)
	}

	// Used to create combinations from 2 to n
	combination := []interface{}{}

	for i := 2; i <= n; i++ {
		combination = append(combination, i)
	}

	var SetMapKeyList []uint64

	for i := 1; i < n; i++ {

		for v := range next.Combination(combination, i, false) {

			b := make([]int, len(v))

			// SetMapKey uses binary numbers to represent elements in set
			// this offers efficient lookups of S - {j}
			var SetMapKey uint64
			for i := range v {
				SetMapKey += BitMap[v[i].(int)]
				b[i] = v[i].(int)
			}

			w, ok := SetMap[SetMapKey]

			if ok {
				log.Fatalf("YOU ARE BAD AT MATH %v", w)
			} else {
				SetMap[SetMapKey] = b
				SolutionsMap[SetMapKey] = make([]float64, n+1)
				setsBySize[i] = append(setsBySize[i], SetMapKey)
				SetMapKeyList = append(SetMapKeyList, SetMapKey)
			}
		}
	}

	// Quick lookups for distances between vertices
	distanceMatrix := make([][]float64, n+1)

	for i := range distanceMatrix {
		distanceMatrix[i] = make([]float64, n+1)
	}

	// Fill in distanceMatrix with all pairs of vertices' distance
	for i := 2; i <= n; i++ {
		thisDistance := distance(VertexMap[1], VertexMap[i])
		distanceMatrix[1][i] = thisDistance
		distanceMatrix[i][1] = thisDistance
	}

	for _, x := range setsBySize[2] {
		thisDistance := distance(VertexMap[SetMap[x][0]], VertexMap[SetMap[x][1]])
		distanceMatrix[SetMap[x][0]][SetMap[x][1]] = thisDistance
		distanceMatrix[SetMap[x][1]][SetMap[x][0]] = thisDistance
	}

	// Fill in SolutionsMap for 0 and 1 to prevent nonsense minTour
	for i := range SolutionsMap {
		SolutionsMap[i][0] = math.Inf(1)
		SolutionsMap[i][1] = math.Inf(1)
	}

	// fmt.Println(solutionsMatrix)
	// Fill in solution for 1 to j
	for _, x := range setsBySize[1] {
		SolutionsMap[x][SetMap[x][0]] = distanceMatrix[1][SetMap[x][0]]
	}

	numOfSets := len(SetMap)
	//TSP
	for m := 2; m < n; m++ {

		// x is the SetMapKey
		for _, x := range setsBySize[m] {

			// j = the VertexMap ID
			for _, j := range SetMap[x] {

				// y is the SetMapKey of S - {j}
				y := x - BitMap[j]

				// Find minCost ranging over sets[y] to j
				minCost := math.Inf(1)
				thisCost := math.Inf(1)
				for _, k := range SetMap[y] {
					// k is the VertexMap ID
					thisCost = SolutionsMap[y][k] + distanceMatrix[k][j]
					minCost = math.Min(minCost, thisCost)
				}
				// set solutionsMatrix[x][j] = minCost
				SolutionsMap[x][j] = minCost
			}
		}
	}

	minTour := math.Inf(1)

	for i, x := range SolutionsMap[SetMapKeyList[len(SetMapKeyList)-1]] {
		minTour = math.Min(minTour, distanceMatrix[1][i]+x)
	}

	fmt.Println(int(minTour))
}

func distance(a, b *Vertex) float64 {
	return math.Sqrt(math.Pow(a.lat-b.lat, 2) + math.Pow(a.long-b.long, 2))
}

func readFile(filename string) {

	i := 1

	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Scan first line
	if scanner.Scan() {
		n, err = strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatalf("couldn't convert number: %v\n", err)
		}
	}

	for scanner.Scan() {

		thisLine := strings.Fields(scanner.Text())

		thisLat, err := strconv.ParseFloat(thisLine[0], 64)
		thisLong, err := strconv.ParseFloat(thisLine[1], 64)

		if err != nil {
			log.Fatal(err)
		}

		VertexMap[i] = &Vertex{thisLat, thisLong}
		i++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
