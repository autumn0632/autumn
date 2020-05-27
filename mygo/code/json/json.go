package json

import (
	"encoding/json"
	"fmt"
	"github.com/json-iterator/go"

)

type Info struct {
	Name     string `json:"学生姓名"`
	Age      int	`json:"学生年龄"`
	Favorite []string	`json:"兴趣爱好"`
	Score 			`json:"学生成绩"`
}

type Score struct {
	Match float64 `json:"数学"`
	Music float64  `json:"音乐"`
	Biology float64 `json:"生物"`
}

var _info = Info{
	Name:     "克林顿",
	Age:      18,
	Favorite: []string{"骑马", "画画", "喝酒"},
	Score: Score{
		10, 11, 12,
	},
}

func StructToJson() string {
	b, err := json.Marshal(_info)
	if err != nil {
		fmt.Printf("%v\n", err)
	}else {
		fmt.Printf("%v\n", string(b))
	}
	return string(b)
}


func StrToJson(){

	//var d interface{}
	var info Info
	//data := "{'学生姓名':'克林顿', '学生年龄':18, '兴趣爱好':['骑马', '画画', '喝酒'], '学生成绩':{'数学':10, '音乐':11, '生物':12}}"
	//data := "{\"Name\":\"克林顿\", \"Age\":18, \"Favorite\":[\"骑马\", \"画画\", \"喝酒\"], \"Score\":{\"Match\":10, \"Music\":11, \"Biology\":12}}"
	data := "{\"学生姓名\":\"克林顿\", \"学生年龄\":18, \"兴趣爱好\":[\"骑马\", \"画画\", \"喝酒\"], \"学生成绩\":{\"数学\":10, \"音乐\":11, \"生物\":12}}"

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal([]byte(data), &info)
	//err := json.Unmarshal([]byte(data), &d)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(info)

}