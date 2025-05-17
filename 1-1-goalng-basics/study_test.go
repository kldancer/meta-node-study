package __1_goalng_basics

import (
	"fmt"
	"testing"
)

func TestIsPalindrome(t *testing.T) {
	a := 1221
	b := 1234
	fmt.Println(isPalindrome(a))
	fmt.Println(isPalindrome(b))
}

func TestIsValid(t *testing.T) {
	a := "()(){[XY0"
	b := "()(){[]}"
	fmt.Println(isValid(a))
	fmt.Println(isValid(b))
}

func TestLongestCommonPrefix(t *testing.T) {
	a := []string{"flower", "flow", "flight"}
	fmt.Println(longestCommonPrefix(a))
}

func TestRemoveDuplicates(t *testing.T) {
	nums := []int{1, 1, 2, 2, 3}
	length := removeDuplicates(nums)
	fmt.Println(nums[:length]) // 输出: [1 2 3]
}

func TestPlusOne(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{9, 9, 9}
	fmt.Println(plusOne(a))
	fmt.Println(plusOne(b))
}

func TestTwoSum(t *testing.T) {
	num := []int{2, 7, 11, 15}
	target := 13
	fmt.Println(twoSum(num, target))
}
