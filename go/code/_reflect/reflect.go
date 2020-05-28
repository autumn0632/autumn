package _reflect

import (
	"fmt"
	"reflect"
)

type Score struct {
	match, music, biology float64
}

type Info struct {
	Name     string
	Age      int
	Favorite []string
	Score
}

func (c *Info) GetName() string {
	return c.Name
}
func (c *Info) GetAge() int {
	return c.Age
}

func PrintInfo() {
	_info := Info{
		Name:     "克林顿",
		Age:      18,
		Favorite: []string{"骑马", "画画", "喝酒"},
		Score: Score{
			10, 11, 12,
		},
	}
	PrintStruct(_info)
}

func PrintStruct(data interface{}) {
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	if t.Kind() == reflect.Struct {
		numField := t.NumField()
		for i := 0; i < numField; i++ {
			f := t.Field(i)
			if f.Type.Kind() == reflect.Struct {
				secondT := f.Type
				secondV := reflect.ValueOf(f)
				fmt.Printf("name: %v - type: %v - field(s):%v \n", secondT.Name(), secondT.Kind(), secondT.NumField())
				for i := 0; i < secondT.NumField(); i++ {
					t := secondT.Field(i)
					val := secondV.Field(i).Interface()
					fmt.Printf("name: %v - type: %v - value:%v \n", t.Name, t.Type, val)
				}
				break
			}
			val := v.Field(i).Interface()
			fmt.Printf("name: %v - type: %v - value: %v\n", f.Name, f.Type, val)
		}
	}
}
