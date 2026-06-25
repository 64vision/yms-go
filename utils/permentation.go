package utils

import (
	"fmt"
	"strings"
)

func GenerateCombinations(nums []int) [][]int {
	var tombok [][]int
	var sahud [][]int
	var combinations [][]int
	// Loop through each pair of numbers in the slice
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			tombok = append(tombok, []int{nums[i], nums[j]})
		}
	}
	for _, nums := range tombok {
		sahud = append(sahud, []int{nums[1], nums[0]})
	}
	combinations = append(tombok, sahud...)
	return combinations
}

func Permutations(arr []int) [][]int {
	return filterComninations(DoPermutations(arr))
}
func filterComninations(arr [][]int) [][]int {
	res := [][]int{}
	isExist := false
	for _, combi := range arr {
		isExist = false
		for _, checkcombi := range res {
			stringcheckcombi := strings.Trim(strings.Replace(fmt.Sprint(checkcombi), " ", " ", -1), "[]")
			stringcombi := strings.Trim(strings.Replace(fmt.Sprint(combi), " ", " ", -1), "[]")
			if stringcheckcombi == stringcombi {
				isExist = true
			}
		}
		if !isExist {
			res = append(res, combi)
		}

	}
	return res
}

func DoPermutations(arr []int) [][]int {
	var helper func([]int, int)
	res := [][]int{}
	helper = func(arr []int, n int) {
		if n == 1 {
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)

		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}

func HasDuplicates(arr []string) bool {
	seen := make(map[string]bool)

	for _, num := range arr {
		if seen[num] {
			return true // Duplicate found
		}
		seen[num] = true
	}
	return false // No duplicates
}
