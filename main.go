package main

import (
	"bufio"
	"fmt"
	"github.com/tarantool/go-tarantool"
	"log"
	"os"
)

func main() {
	conn, err := tarantool.Connect("127.0.0.1:3301", tarantool.Opts{
		User: "ex",
		Pass: "secret",
	})
	if err != nil {
		log.Fatalf("Connection refused")
	}
	defer conn.Close()

	resp, err := conn.Select("user", "primary", 0, 1, tarantool.IterEq, []interface{}{2})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("пользователь", resp)

	//playerName := ""
	//fmt.Scanln(&playerName)
	//funcres, err := conn.Call("mm.user_guild", []interface{}{playerName})
	//fmt.Println(funcres)

	guild := 1
	myscanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Введите сообщение: ")
		myscanner.Scan()
		msg := myscanner.Text()
		fmt.Println(msg)
		funcres, _ := conn.Call("mm.new_msg", []interface{}{msg, guild})
		fmt.Println(funcres)
	}

	//spaceName := "user"
	//indexName := "primary"
	//idFn := conn.Schema.Spaces[spaceName].Fields["user_id"].Id
	//bandNameFn := conn.Schema.Spaces[spaceName].Fields["user_name"].Id
	//
	//var tuplesPerRequest uint32 = 2
	//cursor := []interface{}{}
	//
	//for {
	//	resp, err := conn.Select(spaceName, indexName, 0, tuplesPerRequest, tarantool.IterGt, cursor)
	//	if err != nil {
	//		log.Fatalf("Failed to select: %s", err)
	//	}
	//
	//	if resp.Code != tarantool.OkCode {
	//		log.Fatalf("Select failed: %s", resp.Error)
	//	}
	//
	//	if len(resp.Data) == 0 {
	//		break
	//	}
	//
	//	fmt.Println("Iteration")
	//
	//	tuples := resp.Tuples()
	//	for _, tuple := range tuples {
	//		fmt.Printf("\t%v\n", tuple)
	//	}
	//
	//	lastTuple := tuples[len(tuples)-1]
	//	cursor = []interface{}{lastTuple[idFn], lastTuple[bandNameFn]}
	//}

}
