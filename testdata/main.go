package main

import (
	"fmt"

	"github.com/lack-io/cli"
)

func main() {
	ips := make([]string, 0)

	cli.CommandLine.StringSliceVarP(&ips, "ips", "", ips, "", "")
	cli.CommandLine.Run([]string{"", "--ips=111,333"})

	fmt.Println(ips)
}
