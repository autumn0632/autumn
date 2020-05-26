package src

import (
	"fmt"
	"testing"
)

var (
	l4 = &ListNode{9, nil}
 	l3 = &ListNode{5, l4}
 	l2 = &ListNode{3, l3}
 	l = &ListNode{1, l2}

	m4 = &ListNode{8, nil}
	m3 = &ListNode{7, m4}
	m2 = &ListNode{4, m3}
	m = &ListNode{0, m2}

	n3 = &ListNode{4, nil}
	n2 = &ListNode{3, n3}
	n = &ListNode{1, n2}
 	)

func PrintList(l *ListNode) {
	for l != nil {
		fmt.Printf("%d \t", l.Val)
		l = l.Next
	}
	fmt.Println()
}

func TestTwoSum_1(t *testing.T) {
	nums := []int{1, 2, 6, 9, 19}
	target := 28

	res := TwoSum_1(nums, target)

	fmt.Println(nums, target, res)
}

func TestTwoSum_2(t *testing.T) {
	nums := []int{1, 2, 6, 9, 19}
	target := 28

	res := TwoSum_2(nums, target)

	fmt.Println(nums, target, res)
}

func TestThreeSum(t *testing.T) {
	nums := []int{-1, 0, 1, 2, -1, -4}

	fmt.Println(ThreeSum(nums))
}

func TestAddTwoNumbers(t *testing.T) {
	PrintList(l)
	PrintList(m)
	PrintList(AddTwoNumbers(l, m))

}

func TestMergeTwoLists(t *testing.T) {

	PrintList(l)
	PrintList(m)
	PrintList(MergeTwoLists(l, m))
}
func TestRemoveNthFromEnd(t *testing.T) {
	PrintList(l)
	PrintList(RemoveNthFromEnd(l, 4))
}

func TestMergeNLists(t *testing.T) {
	var newL []*ListNode
	newL = append(newL, l)
	newL = append(newL, m)
	newL = append(newL, n)

	PrintList(l)
	PrintList(m)
	PrintList(n)
	PrintList(MergeNLists(newL))
}

func TestSwapPairs(t *testing.T) {
	l4 := &ListNode{4, nil}
	l3 := &ListNode{3, l4}
	l2 := &ListNode{2, l3}
	l := &ListNode{1, l2}
	SwapPairs(l)
}

func TestIsPalindrome(t *testing.T) {
	i := 1221
	fmt.Printf("%d is 回文数?  %v\n", i, IsPalindromeInt(i))

	str := "abcncbam"
	fmt.Printf("%s is 回文字符串?  %v\n", str, IsPalindromeStr(str))

}



















