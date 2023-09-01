package main

import (
	"flag"
	"fmt"
	"hashlookup/api"
	"hashlookup/inserter"
	"hashlookup/util"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Invalid command. Missing command.")
		return
	}

	command := args[0]

	switch command {
	case "create-tables":
		util.CreateTables()
	case "delete-tables":
		util.DeleteTables()
	case "create-lookup":
		fileName := flag.String("f", "", "Input file name")
		algorithmName := flag.String("h", "", "Hashing algorithm")

		flag.CommandLine.Parse(args[1:])

		inserter.CreateLookup(*fileName, *algorithmName)
	case "create-rainbow":

	case "start-api":
		port := flag.Int("p", 3000, "Port to listen on")
		flag.CommandLine.Parse(args[1:])

		api.StartAPI(*port)
	default:
		fmt.Printf("Invalid command. %s is not a valid command.\n", command)
	}
}
