package src

/*
题目：
给定一个整数数组 nums 和一个目标值 target，
请你在该数组中找出和为目标值的那 两个 整数，并返回他们的数组下标。
你可以假设每种输入只会对应一个答案。但是，数组中同一个元素不能使用两遍。

示例：
nums = [2, 7, 11, 15], target = 9
返回[0,1]

*/

// 方法一：暴力破解，时间复杂度o(n^2)，空间复杂度o(1)
func TwoSum_1(nums []int, target int) []int {
	for i := 0; i < len(nums)-1; i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				return []int{i, j}
			}
		}
	}
	return []int{}
}

// 方法二：map法
// 为了对时间复杂度进行优化，需要更有效的方法来检查数组中是否存在目标元素。
// 通过用hash表的方式，用空间换取时间
func TwoSum_2(nums []int, target int) []int {
	m := make(map[int]int)
	for i, k := range nums {
		if value, exist := m[target-k]; exist {
			return []int{value, i}
		}
		m[k] = i // 将数据保存到map中
	}
	return nil
}
