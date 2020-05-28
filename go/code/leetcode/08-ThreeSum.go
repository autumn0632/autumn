package src

import (
	"sort"
)

/*
给你一个包含 n 个整数的数组 nums，判断 nums 中是否存在三个元素 a，b，c ，使得 a + b + c = 0 ？
请你找出所有满足条件且不重复的三元组。

*/
func ThreeSum(nums []int) [][]int {
	result := [][]int{}
	length := len(nums)
	if len(nums) < 3 {
		return result
	}
	sort.Ints(nums) // 先排序
	temp := 0
	for start, value := range nums { // 第一个数的移动逻辑
		if value > 0 {  // 对于排完序的队列，第一个数大于0
			return result
		}
		if start > 0 && value == nums[start-1] {
			continue
		}
		left, right := start+1, length-1 // 确定第一个第二个数
		for left < right { // 第二个数和第三个数的移动逻辑
			temp = value + nums[left] + nums[right]
			if temp == 0 {
				result = append(result, append([]int{}, start, nums[left], nums[right]))
				for left < right && nums[left] == nums[left+1] {
					left += 1
				}
				for left < right && nums[right] == nums[right-1] {
					right -= 1
				}
				left += 1
				right -= 1
			} else if temp > 0 {
				right -= 1
			} else if temp < 0 {
				left += 1
			}
		}
	}
	return result
}
