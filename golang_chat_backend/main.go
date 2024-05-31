package main

import (
	"log"
	"websocket/network"
)

func init() {
	log.Print("init first")
}

func main() {

	//defer func() {
	//	r := recover()
	//	fmt.Println("recover", r)
	//}()
	n := network.NewServer()
	n.StartServer()
	//err := n.StartServer()
	//if err != nil {
	//	log.Panic("err", err)
	//}

}
