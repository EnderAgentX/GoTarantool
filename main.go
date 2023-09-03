package main

import (
	"fmt"
	"github.com/tarantool/go-tarantool"
	"log"
)

func main() {
	conn, err := tarantool.Connect("127.0.0.1:3301", tarantool.Opts{
		User: "example",
		Pass: "secret",
	})
	if err != nil {
		log.Fatalf("Connection refused")
	}
	defer conn.Close()

	resp, err := conn.Select("example", "primary", 0, 1, tarantool.IterEq, []interface{}{2})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)

	spaceName := "tester"
	indexName := "scanner"
	idFn := conn.Schema.Spaces[spaceName].Fields["id"].Id
	bandNameFn := conn.Schema.Spaces[spaceName].Fields["band_name"].Id

	var tuplesPerRequest uint32 = 2
	cursor := []interface{}{}

	for {
		resp, err := conn.Select(spaceName, indexName, 0, tuplesPerRequest, tarantool.IterGt, cursor)
		if err != nil {
			log.Fatalf("Failed to select: %s", err)
		}

		if resp.Code != tarantool.OkCode {
			log.Fatalf("Select failed: %s", resp.Error)
		}

		if len(resp.Data) == 0 {
			break
		}

		fmt.Println("Iteration")

		tuples := resp.Tuples()
		for _, tuple := range tuples {
			fmt.Printf("\t%v\n", tuple)
		}

		lastTuple := tuples[len(tuples)-1]
		cursor = []interface{}{lastTuple[idFn], lastTuple[bandNameFn]}
	}

}
