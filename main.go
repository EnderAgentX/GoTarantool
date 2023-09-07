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

	//resp, err := conn.Select("user", "primary", 0, 1, tarantool.IterEq, []interface{}{2})
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println("пользователь", resp)

	//playerName := ""
	//fmt.Scanln(&playerName)
	//funcres, err := conn.Call("mm.user_guild", []interface{}{playerName})
	//fmt.Println(funcres)

	myscanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Введите имя пользователя: ")
	myscanner.Scan()
	myUser := myscanner.Text()
	fmt.Println(myUser)

	info, _ := conn.Call("mm.login", []interface{}{myUser})
	tuples := info.Tuples()
	userId := tuples[0][0]
	guildId := tuples[1][0]
	fmt.Println(userId)
	guildName, _ := conn.Call("mm.user_guild", []interface{}{myUser})

	for {
		fmt.Print("Введите сообщение: ")
		myscanner.Scan()
		msg := myscanner.Text()
		fmt.Printf("%s(%s): %s", myUser, guildName.Tuples()[0][0], msg)
		fmt.Println("")
		_, _ = conn.Call("mm.new_msg", []interface{}{msg, guildId, userId})
		//fmt.Println(funcres)
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
