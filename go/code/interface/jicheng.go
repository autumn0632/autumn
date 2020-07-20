package inter

import (
	"fmt"
	"math"
)

type ShapeInterface interface {
	Area() float64
	GetName() string
	PrintArea()
}

// 标准形状，面积为0.0
type Shape struct {
	name string
}

func (s *Shape) Area() float64{
	return 0.0
}

func (s *Shape) GetName() string {
	return s.name
}

func (s *Shape) PrintArea() {
	fmt.Printf("%s: Area  %v\n", s.name, s.Area())
}

// 矩形
type Rectangle struct {
	Shape
	w, h float64
}

func (r *Rectangle) Area() float64 {
	return r.w * r.h
}

// 圆形  : 重新定义 Area 和PrintArea 方法
type Circle struct {
	Shape
	r float64
}

func (c *Circle) Area() float64 {
	return c.r * c.r * math.Pi
}

func (c *Circle) PrintArea() {
	fmt.Printf("%s : Area %v\n", c.GetName(), c.Area())
}

func Run() {

	//c和r实例没有重写name方法，c实例没有重写PrintArea()方法。没有重写的方法均调用的s的方法

	s := Shape{name:"Shape"}
	c := Circle{Shape: Shape{name:"Circle"}, r:10}
	r := Rectangle{Shape:Shape{name:"Rectangle"}, w:5, h:4}



	listShape := []ShapeInterface{&s, &c, &r}

	for _, si := range listShape {
		si.PrintArea()
	}



}