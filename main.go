package main

import (
	"fmt"
	"os"

	"github.com/Can-U-Join-Us/CUJU-Backend/modules/server" // Host serve when loaded
)

var mode = 0

func main() {
	fmt.Println("Args : ", os.Args)
	if len(os.Args) > 0 && os.Args[1] == `dev` {
		mode = 1
	}
	server.Serve(mode)
}
