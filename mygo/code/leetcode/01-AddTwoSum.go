package src

/*
题目：
给出两个 非空 的链表用来表示两个非负的整数。
它们各自的位数是按照 逆序 的方式存储的，并且它们的每个节点只能存储 一位 数字。
如果，我们将这两个数相加起来，则会返回一个新的链表来表示它们的和。

example：
输入：(2 -> 4 -> 3) + (5 -> 6 -> 4)
输出：7 -> 0 -> 8
原因：342 + 465 = 807

思路：
1. 新建哨兵节点，返回值返回哨兵节点
2. 新建变量，初始值等于哨兵节点，后续根据逻辑移动
3. 进位的处理
*/
func AddTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	l := new(ListNode)
	if l1 == nil && l2 == nil {
		return l
	}
	newNode := l
	carry := 0
	for l1 != nil || l2 != nil || carry > 0 {
		a, b := 0, 0
		newNode.Next = new(ListNode) // 新建节点,连接到哨兵节点后面
		newNode = newNode.Next       // newNode 移动
		if l1 != nil {
			a = l1.Val
		}
		if l2 != nil {
			b = l2.Val
		}
		newNode.Val = (a + b + carry) % 10
		carry = (a + b + carry) / 10

		if l1 != nil {
			l1 = l1.Next
		}
		if l2 != nil {
			l2 = l2.Next
		}
	}
	return l.Next
}