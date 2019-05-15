package main

import (
	"github.com/affiliator/mgw/cmd"
	"github.com/affiliator/mgw/storage"
)

func main() {
	cmd.Execute()

	defer storage.Close()
}
