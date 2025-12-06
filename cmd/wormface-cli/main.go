package main

import (
	"fmt"
	"os"

	"github.com/OrkWard/wormface/internal/cmd"
	"github.com/OrkWard/wormface/internal/utils"
)

func main() {
	if err := utils.InitConfig(); err != nil {
		panic(err)
	}

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Usage: wormface-cli [twitter | pixiv] ...")
		os.Exit(1)
	}

	switch args[1] {
	case "twitter":
		cmd.CmdTwitter(args[2:])
	case "pixiv":
		cmd.CmdPixiv(args[2:])
	default:
		fmt.Println("Usage: wormface-cli [twitter | pixiv] ...")
		os.Exit(1)
	}
}
