package src

/*
给定一个链表，删除链表的倒数第 n 个节点，并且返回链表的头结点。

example:
input: 1->2->3->4->5, 2
output: 1->2->3->5

题解：
1. 使用快慢指针,当快指针走完n的节点时，慢指针开始走。当快指针走到尾结点时，慢指针指向倒数第n个
*/

func RemoveNthFromEnd(head *ListNode, n int) *ListNode {
	l := new(ListNode)
	l.Next = head // 返回l.Next
	cur := l
	var pre *ListNode

	i := 1
	for head != nil {
		// 慢指针开始移动，等head移动到最后一个，pre指向倒数第N个节点
		if i >= n {
			pre = cur
			cur = cur.Next
		}
		head = head.Next
		i ++
	}
	// 删除动作：快指针一直移动，当快指针head移动到最后一个，慢指针pre指向倒数第N个
	pre.Next = pre.Next.Next
	return l.Next
}
