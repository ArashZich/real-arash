package main

import (
	"gitag.ir/armogroup/armo/services/reality/cmd/chef/cmd"
	"gitag.ir/armogroup/armo/services/reality/config"
)

func main() {
	config.Load()

	cmd.Execute()
}
