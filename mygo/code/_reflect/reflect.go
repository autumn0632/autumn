package _reflect

import (
	"fmt"
	"reflect"
)

type Score struct {
	match, music, biology float64
}

type Info struct {
	Name string
	Age int
	Favorite []string
	Score
}


func PrintInfo() {
	_info := Info{
		Name:     "克林顿",
		Age:      18,
		Favorite: []string{"骑马", "画画", "喝酒"},
		Score:    Score{
			10,11,12,
		},
	}
	PrintStruct(_info)
}

func PrintStruct(data interface{})  {
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	numField := t.NumField()
	for i:= 0; i < numField; i ++ {
		f := t.Field(i)
		val := v.Field(i).Interface()
		fmt.Printf("name: %v - type: %v - value: %v\n", f.Name, f.Type, val)
	}
}