package __1_goalng_basics

import "fmt"

/*
回文数,

考察：数字操作、条件判断
题目：判断一个整数是否是回文数
*/
func isPalindrome(x int) bool {
	s := fmt.Sprintf("%d", x)
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		if s[i] != s[j] {
			return false
		}
	}
	return true
}

/*
有效的括号 ,

考察：字符串处理、栈的使用

题目：给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串，判断字符串是否有效

链接：https://leetcode-cn.com/problems/valid-parentheses/
*/
func isValid(s string) bool {
	base := map[string]struct{}{
		"(": {},
		")": {},
		"{": {},
		"}": {},
		"[": {},
		"]": {},
	}
	for _, str := range s {
		if _, ok := base[string(str)]; !ok {
			return false
		}
	}
	return true
}

/*
最长公共前缀,

考察：字符串处理、循环嵌套

题目：查找字符串数组中的最长公共前缀

链接：https://leetcode-cn.com/problems/longest-common-prefix/
*/
func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}

	prefix := strs[0]

	for i := 1; i < len(strs); i++ {
		j := 0
		for j < len(prefix) && j < len(strs[i]) && prefix[j] == strs[i][j] {
			j++
		}
		prefix = prefix[:j]
		if prefix == "" {
			return ""
		}
	}

	return prefix
}

/*
删除排序数组中的重复项 ,

难度：简单

考察：数组/切片操作

题目：给定一个排序数组，你需要在原地删除重复出现的元素

链接：https://leetcode-cn.com/problems/remove-duplicates-from-sorted-array/
*/
func removeDuplicates(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	i := 0 // 指向不重复部分的最后一个位置
	for j := 1; j < len(nums); j++ {
		if nums[j] != nums[i] {
			i++
			nums[i] = nums[j]
		}
	}

	return i + 1 // 返回去重后数组的长度
}

/*
加一 ,

难度：简单

考察：数组操作、进位处理

题目：给定一个由整数组成的非空数组所表示的非负整数，在该数的基础上加一

链接：https://leetcode-cn.com/problems/plus-one/
*/
func plusOne(digits []int) []int {
	n := len(digits)
	for i := n - 1; i >= 0; i-- {
		digits[i]++
		if digits[i] < 10 {
			return digits
		}
		digits[i] = 0
	}
	res := make([]int, n+1)
	res[0] = 1
	return res
}

/*
基础,
两数之和 ,

考察：数组遍历、map使用

题目：给定一个整数数组 nums 和一个目标值 target，请你在该数组中找出和为目标值的那两个整数

链接：https://leetcode-cn.com/problems/two-sum/
*/
func twoSum(nums []int, target int) []int {
	hashMap := make(map[int]int)

	for i, num := range nums {
		complement := target - num
		if j, ok := hashMap[complement]; ok {
			return []int{j, i}
		}
		hashMap[num] = i
	}

	return nil // 没有找到符合条件的两个数
}
