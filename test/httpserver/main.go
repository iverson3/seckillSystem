package main

import (
	bytes2 "bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"math"
	"math/big"
	"net/http"
	"strconv"
	"unicode"
	"unicode/utf8"
)


type Person struct {
	Name string `json:"name";db:"username"`
	Age int8 `json:"age";db:"age"`
	Score float32 `json:"score";db:"score"`
}

func (p Person) work() {
	fmt.Printf("%s is working...\n", p.Name)
}

type Skyer interface {
	Fly() bool
	Stop(string) int
}

func (p Person) Fly() bool {
	return true
}
func (p Person) Stop(str string) int {
	a := string(p.Age) + str
	i, err := strconv.Atoi(a)
	//i, err := strconv.ParseInt(a, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func doSome(s Skyer) {
	fmt.Printf("fly: %v\n", s.Fly())
}

func main() {
	//http.HandleFunc("/", handler)
	//err := http.ListenAndServe("localhost:8888", nil)
	//if err != nil {
	//	panic(err)
	//}
	//
	//logs.Info("httpserver start to listen...")


	var sky1 Skyer
	sky1 = Person{Name: "xxx", Age: 10, Score: 88}
	fmt.Println(sky1)
	doSome(sky1)

	p2 := Person{}
	doSome(p2)



	// 超大的数，int64 uint64都无法表示的数

	//var bigN int = 3e28
	//fmt.Println(bigN)

	//bigNum := big.NewInt(3e28)
	//fmt.Println(bigNum)

	bigInt := new(big.Int)
	bigInt.SetString("30000000000000000000000000000", 10)
	fmt.Println(bigInt)

	num1 := big.NewInt(3000000)

	seconds := new(big.Int)
	seconds.Div(bigInt, num1)
	fmt.Println(seconds)

	marshal, err := json.Marshal(111)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v \n", marshal)
	res := string(marshal)
	fmt.Printf("result: %s\n", res)

	p1 := Person{}
	p1.Name = "stefan"
	p1.Age  = 27
	p1.Score= 98
	p1.work()


	// golang中string底层是通过byte数组实现的。中文字符在unicode下占2个字节，在utf-8编码下占3个字节，而golang默认编码正好是utf-8。
	// byte 等同于uint8，常用来处理ascii字符
	// rune 等同于int32，常用来处理unicode或utf-8字符

	// Unicode是ASCII的超集，它们前128个code points是一样的；(包括数字 英文字母 常用标点符号)
	a := '第'
	fmt.Println(a)
	fmt.Printf("%s \n", string(a))

	b := `第x繓`
	fmt.Println(len(b))
	fmt.Println(utf8.RuneCountInString(b))

	runes := []rune(b)
	fmt.Println(len(runes))
	fmt.Println(runes[0])
	fmt.Println(string(runes))
	runes[0] = 'w'
	fmt.Println(string(runes))

	bytes := []byte(b)
	fmt.Println(len(bytes))
	fmt.Println(bytes[0])
	fmt.Println(string(bytes))
	bytes[0] = 'w'
	fmt.Println(string(bytes))


	var f1 float64 = 32768.0
	// 类型转换前要确保转换后的数值大小在目标类型的数值范围内，否则就会得到错误的转换结果
	if f1 <= math.MaxInt16 && f1 >= math.MinInt16 {
		i2 := int16(f1)
		fmt.Println(i2)
	} else {
		fmt.Println("f1 can not convert to int16")
	}

	inString, size := utf8.DecodeLastRuneInString("撒X伺")
	fmt.Printf("%s, %d \n", string(inString), size)

	//unicode.ASCII_Hex_Digit

	fmt.Printf("%d \n", unicode.MaxASCII)
	fmt.Printf("%d \n", math.MaxInt8)

	newStr := strconv.AppendFloat([]byte("float64: "), 3.14159, 'E', 3, 64)
	fmt.Println(string(newStr))

	res2 := bytes2.Contains([]byte("xxx"), []byte("x"))
	fmt.Println(res2)
	//res3 := bytes2.SplitAfter([]byte("a,b,c"), []byte(","))
	res3 := bytes2.SplitN([]byte("a,b,c"), []byte(","), -1)
	for _, v := range res3 {
		fmt.Println(string(v))
	}
	for i := 0; i < len(res3); i++ {
		fmt.Println(string(res3[i]))
	}

}

func handler(w http.ResponseWriter, req *http.Request)  {
	defer req.Body.Close()

	var buf []byte
	_, err := req.Body.Read(buf)
	if err != nil {
		logs.Warn("read from client failed! error: ", err)
		return
	}

	logs.Info("got from client: ", string(buf))

	_, err = w.Write([]byte("response from httpserver"))
	if err != nil {
		logs.Warn("write to client failed! error: ", err)
	}
	logs.Info("write to client successfully!")
}