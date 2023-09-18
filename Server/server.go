package Server

import (
	"github.com/tarantool/go-tarantool"
	"log"
)

func Server() *tarantool.Connection {
	conn, err := tarantool.Connect("127.0.0.1:3301", tarantool.Opts{
		User: "ex",
		Pass: "secret",
	})
	if err != nil {
		log.Fatalf("Connection refused")
	}

	return conn
}
