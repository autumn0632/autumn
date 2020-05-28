package src

/*
回文数：判断一个整数是否是回文数。回文数是指正序（从左向右）和倒序（从右向左）读都是一样的整数。


*/

func IsPalindromeInt(x int) bool {
	if x < 0 || (x % 10 == 0 && x != 0) {
		return false
	}
	revertNum := 0
	for x > revertNum {
		revertNum = revertNum * 10 + x % 10  // 取余
		x /=  10 // 整除
	}
	return x == revertNum || x == revertNum / 10 // 去掉奇数位
}
func IsPalindromeStr(str string) bool {
	length := len(str)

	for i:= 0; i < length / 2; i ++ {
		if str[i] != str[length-i-1] {
			return false
		}
	}
	return true
}
