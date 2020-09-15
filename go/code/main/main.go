package main

import "fmt"

func main() {
	//Map := make(map[int]int)
	Map := map[string]int{
		"n":64,
	}

	for i := 0; i < 100000; i++ {
		//go writeMap(Map, i, i)
		go readMap(Map)
	}

}

func readMap(Map map[string]int) {
	fmt.Println(Map["n"])
	//return Map[key]
}


func writeMap(Map map[int]int, key int, value int) {
	Map[key] = value
}
