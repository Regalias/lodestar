// Package utils provides internal utility functions for the lodestar library.
package utils

// ReverseSliceInPlace reverses a slice of any type T in place with two pointers
func ReverseSliceInPlace[T any](list []T) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}
