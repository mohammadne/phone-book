package utils

import "golang.org/x/exp/constraints"

func Sort[T constraints.Ordered](array []T) []T {
	for pointer := 0; pointer < len(array)-1; pointer++ {
		for index := pointer + 1; index > 0; index-- {
			if array[index-1] <= array[index] {
				break
			}
			array[index-1], array[index] = array[index], array[index-1]
		}
	}

	return array
}
