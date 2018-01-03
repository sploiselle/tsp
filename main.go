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

var VertexMap = make(map[int]*Vertex)

func main() {

	readFile(os.Args[1])

	// fmt.Println(VertexMap)

	sets := make([][]int, 0)

	for i := range sets {
		sets[i] = make([]int, 1)
	}

	setsBySize := make([][]int, n)
	setsCount := 0

	for i := range setsBySize {
		setsBySize[i] = make([]int, 0)
	}

	combination := []interface{}{}

	for i := 2; i <= n; i++ {
		combination = append(combination, i)
	}

	for i := 1; i < n; i++ {

		for v := range next.Combination(combination, i, false) {

			b := make([]int, len(v))

			for i := range v {
				b[i] = v[i].(int)
			}

			sets = append(sets, b)
			setsBySize[i] = append(setsBySize[i], setsCount)
			setsCount++
		}
	}

	//	fmt.Println(sets)
	//	fmt.Println(setsBySize)

	distanceMatrix := make([][]float64, n+1)

	// Quick lookups for distances between vertices
	for i := range distanceMatrix {
		distanceMatrix[i] = make([]float64, n+1)
	}

	// Fill in distanceMatrix with all pairs of vertices' distance
	for _, x := range setsBySize[2] {

		// fmt.Printf("\nsets[%d]:\t%v", x, sets[x])

		// fmt.Printf("\na:%v\t\nb:%v\t", sets[x][0], sets[x][1])

		thisDistance := distance(VertexMap[sets[x][0]], VertexMap[sets[x][1]])
		distanceMatrix[sets[x][0]][sets[x][1]] = thisDistance
		distanceMatrix[sets[x][1]][sets[x][0]] = thisDistance

		// for _, y := range sets[x] {

		// 	fmt.Println("y", y)
		// 	// thisDistance := distance(VertexMap[y[0]], VertexMap[y[1]])
		// 	// distanceMatrix[y[0]][y[1]] = thisDistance
		// 	// distanceMatrix[y[1]][y[0]] = thisDistance
		// }
		// fmt.Println(x)

	}

	for i := 2; i <= n; i++ {
		thisDistance := distance(VertexMap[1], VertexMap[i])
		distanceMatrix[1][i] = thisDistance
		distanceMatrix[i][1] = thisDistance
	}

	// fmt.Println(distanceMatrix)

	// [number of sets][solution for adding j to S-{j}]
	solutionsMatrix := make([][]float64, len(sets))

	// Quick lookups for distances between vertices
	for i := range solutionsMatrix {
		solutionsMatrix[i] = make([]float64, n+1)
		solutionsMatrix[i][0] = math.Inf(1)
		solutionsMatrix[i][1] = math.Inf(1)
	}

	// fmt.Println(solutionsMatrix)
	// Fill in solution for 1 to j
	for _, x := range setsBySize[1] {
		solutionsMatrix[x][sets[x][0]] = distanceMatrix[1][sets[x][0]]
	}

	numOfSets := len(sets)
	//TSP
	for m := 2; m < n; m++ {
		for _, x := range setsBySize[m] {
			// x is the setID
			fmt.Println(numOfSets - x)
			for _, j := range sets[x] {
				// j = the VertexMap ID
				//				fmt.Println("sets[x]", sets[x])
				//				fmt.Println("j", j)

				// Remove j from sets[x]

				setMinusJ := removeElementFromSlice(sets[x], j)
				//				fmt.Println("setMinusJ", setMinusJ)

				// Search setsBySize[m-1] for the sets ID that matches sets[x]-j;
				//		call it sets[y]

				for _, y := range setsBySize[m-1] {
					// y is the setID
					if testEq(setMinusJ, sets[y]) {

						// Find minCost ranging over sets[y] to j
						minCost := math.Inf(1)
						thisCost := math.Inf(1)

						for _, k := range sets[y] {
							// k is the VertexMap ID
							thisCost = solutionsMatrix[y][k] + distanceMatrix[k][j]

							minCost = math.Min(minCost, thisCost)
						}

						// set solutionsMatrix[x][j] = minCost
						solutionsMatrix[x][j] = minCost
					}
				}
			}
		}
	}

	// // Print solutionsMatrix
	// for i := range solutionsMatrix {
	// 	fmt.Println(solutionsMatrix[i])
	// }

	minTour := math.Inf(1)

	for i, x := range solutionsMatrix[len(sets)-1] {
		minTour = math.Min(minTour, distanceMatrix[1][i]+x)
	}

	fmt.Println(int(minTour))
}

func testEq(a, b []int) bool {
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}
func removeElementFromSlice(s []int, v int) []int {

	for p, x := range s {
		if x == v {
			r := append([]int{}, s[:p]...)
			r = append(r, s[p+1:]...)
			return r
		}
	}
	log.Fatalf("Cannot remove %d from %v", v, s)

	return s
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

// func Factorial(n int) (result int) {
// 	if n > 0 {
// 		result = n * Factorial(n-1)
// 		return result
// 	}
// 	return 1
// }

// useThis := &c
// fmt.Println("Trying to add", useThis)
// fmt.Println("currently", i+1, setsBySize[i+1])
// setsBySize[i+1] = append(setsBySize[i+1], *useThis)
// fmt.Println("after assignment", i+1, setsBySize[i+1])
