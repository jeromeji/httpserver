package main

import "fmt"

func syc(arry *[]int) {
	for i := 0; i < len(*arry)-1; i++ {
		for j := 0; j < len(*arry)-i-1; j++ {
			if (*arry)[j] > (*arry)[j+1] {
				tmp := (*arry)[j]
				(*arry)[j] = (*arry)[j+1]
				(*arry)[j+1] = tmp
			}
		}
	}

	fmt.Println(*arry)
}

func main() {
	arry := []int{1, 3, 2, 100, 60}
	syc(&arry)
}
