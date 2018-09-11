package main

import (
	"fmt"

	"github.com/hypebeast/go-osc/osc"
)

func main() {
	client := osc.NewClient("localhost", 4559)

	msg := osc.NewMessage("/trigger/prophet")
	msg.Append(int32(70))
	msg.Append(int32(100))
	msg.Append(int32(8))
	err := client.Send(msg)
	fmt.Println(err)
}
