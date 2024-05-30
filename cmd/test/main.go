package main

import (
	"fmt"
	"github.com/nats-io/stan.go"
	"os"
)

func main() {
	sc, err := stan.Connect("test-cluster", "client5")
	if err != nil {
		fmt.Println("Ошибка при подключении к серверу nats streaming: ", err)
		return
	}
	bytes, _ := os.ReadFile(os.Args[1])
	fmt.Println("Data: ", string(bytes))
	err = sc.Publish("foo", bytes)
	if err != nil {
		fmt.Println("Err: ", err)
		return
	}
	fmt.Println("Shipped")

	sc.Close()
}
