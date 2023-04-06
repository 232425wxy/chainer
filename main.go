package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	random := rand.New(rand.NewSource(time.Now().Unix()))
	max := []float64{547.0, 642.0, 756.0, 903.0, 1095.0, 1237.0}
	min := []float64{453.0, 538.0, 632.0, 758.0, 925.0, 1111.0}
	fmt.Print("[")
	for k := 0; k < 100; k++ {
		for i := 0; i < 6; i++ {
			diff := max[i] - min[i]
			val := random.Float64()*diff + min[i]
			if i != 5 {
				fmt.Printf("%.0f ", val)
			} else {
				fmt.Printf("%.0f", val)
			}
		}
		if k != 99 {
			fmt.Println()
		}
	}
	fmt.Print("]\n")
}
