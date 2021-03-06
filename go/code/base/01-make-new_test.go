package base

import (
	"fmt"
	"log"
	"testing"
)

func Test(t *testing.T) {

	// 声明引用类型、符合类型等，只是声明，并没有分配内存
	var a map[string]string
	var b map[string]interface{}

	//a["a"] = "a"  // panic, 分配给了空map
	fmt.Printf("1: %#v, %#v\n", a, a["a"])
	fmt.Printf("2: %#v, %#v\n", b, b["b"])

	test := make(map[string]string)
	//test["a"] = "a"
	fmt.Printf("3: %p, %#v, %#v\n", test, test, test["aaa"]) // test["aaa"] 默认值为""

	test1 := make(map[string]interface{})
	fmt.Printf("4: %p, %#v, %#v\n", test1, test1, test1["aaa"]) // test1["aaa"] 默认值为nil

	test2 := make(map[string]map[string]interface{})
	fmt.Printf("5: %p, %#v, %#v, %#v\n", test2, test2, test2["aaa"], test2["aaa"]["aaa"]) // test1["aaa"] 默认值为nil
	//test2["aaa"]["aaa"] = "test"    // 无法给nil 赋值
	test2["aaa"] = make(map[string]interface{})
	test2["aaa"]["aaa"] = "test"

	c := new(map[string]string)
	fmt.Printf("new map:%p %#v, %#v\n", c, c, *c)

	d := new(string) // 值类型已经默认分配好内存
	fmt.Printf("new string: %#v\n", d)

}

func TestForRange(t *testing.T) {

	slice := []int{0, 1, 2, 3}
	myMap := make(map[int]*int)

	for index, value := range slice {
		log.Printf("%d, %p", value, &value)
		myMap[index] = &value
	}
	fmt.Println("=====new map=====")
	for key, value := range myMap {
		fmt.Printf("map[%v]=%v\n", key, *value)
	}

}

func TestForRange2(t *testing.T) {

	slice := []int{0, 1, 2, 3}
	myMap := make(map[int]*int)

	for index, value := range slice {
		name := value
		log.Printf("%d, %p", name, &name)
		myMap[index] = &name
	}
	fmt.Println("=====new map=====")
	for key, value := range myMap {
		fmt.Printf("map[%v]=%v\n", key, *value)
	}

}

func TestA(t *testing.T) {
	a := func(base int) (func(int) int, func(int) int) {
		add := func(i int) int {
			base += i
			return base
		}
		sub := func(i int) int {
			base -= i
			return base

		}

		return add, sub
	}

	add, sub := a(100)

	fmt.Println(sub(200), add(200))

	fmt.Println(add(200), sub(200))

}

func TestMap(t *testing.T) {
	var te map[string]map[string]float64

	fmt.Printf("%v\n", te)
	fmt.Printf("%p\n", te)
	if te == nil {
		fmt.Println("te is nil")
	}
	te = map[string]map[string]float64{

	}
	fmt.Printf("%v\n", te)
	fmt.Printf("%p\n", te)
	if te == nil {
		fmt.Println("te is nil")
	}

}

func TestString(t *testing.T) {
	var data string

	fmt.Println(len(data))

}

func TestSlice(t *testing.T) {
	data := make([]int, 1)

	data = append(data, 1,2,3,4,5)

	fmt.Println(data)

}