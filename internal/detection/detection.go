package detection

import (
	"GoDupeDetector/internal/parsing"
	"errors"
	"math"
	"sync"
)

type CloneSet struct {
	Clones []Clone
}

type Clone struct {
	FunctionId1 uint
	FunctionId2 uint
}

func DetectClones(pset *parsing.ParseSet, threshold float64) (*CloneSet, error) {
	if threshold < 0.0 || threshold > 1.0 {
		return nil, errors.New("invalid threshold value")
	}

	var wg sync.WaitGroup
	wg.Add(1) // Hold waitgroup to not end immediately

	cclones := make(chan Clone, 100)

	go func() {
		wg.Wait()
		close(cclones)
	}()

	functions := make([]parsing.Function, 0, len(pset.Functions))
	for _, function := range pset.Functions {
		functions = append(functions, function)
	}

	for i := 0; i < len(pset.Functions); i++ {
		f1 := functions[i]
		wg.Add(1)
		go func(f1 parsing.Function, i int, functions []parsing.Function, thresold float64) {
			defer wg.Done()

			for j := i + 1; j < len(functions); j++ {
				f2 := functions[j]
				if isClone(f1, f2, threshold) {
					cclones <- Clone{
						FunctionId1: f1.Id,
						FunctionId2: f2.Id,
					}
				}
			}
		}(f1, i, functions, threshold)
	}
	wg.Done() // Release extra hold

	clones := make([]Clone, 0, 100)
	for clone := range cclones {
		clones = append(clones, clone)
	}

	return &CloneSet{Clones: clones}, nil
}

func isClone(f1 parsing.Function, f2 parsing.Function, threshold float64) bool {
	len1 := len(f1.PrettyPrintBody)
	len2 := len(f2.PrettyPrintBody)
	reqshared := int(math.Ceil(float64(max(len1, len2)) * threshold))

	if len1 < reqshared || len2 < reqshared {
		return false
	}

	if lcslength(f1.PrettyPrintBody, f2.PrettyPrintBody) >= reqshared {
		return true
	} else {
		return false
	}
}

func lcslength(x []string, y []string) int {
	lenx := len(x)
	leny := len(y)
	c := matrix(lenx, leny)

	for i := 0; i < lenx; i++ {
		c[i][0] = 0
	}
	for j := 0; j < leny; j++ {
		c[0][j] = 0
	}

	for i := 1; i < lenx; i++ {
		for j := 1; j < leny; j++ {
			if x[i] == y[j] {
				c[i][j] = c[i-1][j-1] + 1
			} else {
				c[i][j] = max(c[i][j-1], c[i-1][j])
			}
		}
	}

	return c[lenx-1][leny-1]
}

func max(x int, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
}

func matrix(a int, b int) [][]int {
	matrix := make([][]int, a)
	for i := 0; i < a; i++ {
		matrix[i] = make([]int, b)
	}
	return matrix
}
