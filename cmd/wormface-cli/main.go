package main

import (
	"fmt"
	"log"
	"os"

	"github.com/OrkWard/wormface/internal/cmd"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
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
