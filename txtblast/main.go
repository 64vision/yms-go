package main

import (
	"fmt"

	"gollux/sms"
)

func main() {
	fmt.Println("Test")
	sms.Send("***", "Test message from . Please ignore.")
}

//09317143427
