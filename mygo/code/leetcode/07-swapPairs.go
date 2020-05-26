package src

/*
给定一个链表，两两交换其中相邻的节点，并返回交换后的链表。
不能只是单纯的改变节点内部的值，而是需要实际的进行节点交换。

example:
给定 1->2->3->4, 返回 2->1->4->3.

题解：
1. 使用递归，当只有两个节点时开始处理
2. 处理时返回第二个节点（第二个节点指向下一个交互的第一个节点）

*/
func SwapPairs(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}

	firstNode := head
	secondNode := head.Next

	// 交换
	firstNode.Next = SwapPairs(secondNode.Next)
	secondNode.Next = firstNode

	return secondNode
}