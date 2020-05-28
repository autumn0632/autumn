package json

import "testing"

func TestStructToJson(t *testing.T) {
	//StructToJson()
	StrToJson()
	t.Log("ok")
}


func BenchmarkStrToJson(b *testing.B) {
	for n := 0; n < b.N; n++ {
		StrToJson()
	}
}