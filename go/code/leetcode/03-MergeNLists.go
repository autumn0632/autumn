package src

/*
合并 k 个排序链表，返回合并后的排序链表。
请分析和描述算法的复杂度。

example:
输入:
[
  1->4->5,
  1->3->4,
  2->6
]
输出: 1->1->2->3->4->4->5->6

思路：
1. 递归算法，循环将前两条链表合并后的结果，与第三条合并
2. 新建一个空链表进行比较，方便循环
*/

func MergeNLists(lists []*ListNode) *ListNode {
	var temp *ListNode
	length := len(lists)
	for i := 0; i < length; i ++ {
		temp = mergeTwoLists(temp, lists[i])
	}
	return temp
}

func mergeTwoLists(l1 *ListNode, l2 *ListNode) *ListNode {
	l := new(ListNode)
	temp := l

	for l1 != nil && l2 != nil {
		if l1.Val < l2.Val {
			temp.Next = l1
			l1 = l1.Next
		} else {
			temp.Next = l2
			l2 = l2.Next
		}
		temp = temp.Next
	}
	if l1 != nil {
		temp.Next = l1
	}
	if l2 != nil {
		temp.Next = l2
	}
	return l.Next
}