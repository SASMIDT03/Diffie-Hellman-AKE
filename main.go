package main

import (
	"diffie-hellman-ake/peer"
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello World")

	alice := peer.NewPeer("Alice", "50001")
	bob := peer.NewPeer("Bob", "50002")

	listenerAddrAlice := alice.GetListenerAddr()
	fmt.Println("listener addr of Alice: ", listenerAddrAlice)
	fmt.Println("listener addr of Bob: ", bob.GetListenerAddr())

	bob.Connect(listenerAddrAlice)

	time.Sleep(3 * time.Second)

	bob.SendMsg("Hello Alice")

	// cursed, I know
	for {
	}
}
