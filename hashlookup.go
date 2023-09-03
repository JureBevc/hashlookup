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
	case "reset-tables":
		util.DeleteTables()
		util.CreateTables()
	case "create-lookup":
		fileName := flag.String("file", "", "Input file name")
		algorithmName := flag.String("alg", "", "Hashing algorithm")
		flag.CommandLine.Parse(args[1:])

		*algorithmName = util.FormatAlgorithmName(*algorithmName)

		inserter.CreateLookup(*fileName, *algorithmName)
	case "check-lookup":
		hash := flag.String("hash", "", "Input file name")
		algorithmName := flag.String("alg", "", "Hashing algorithm")
		flag.CommandLine.Parse(args[1:])

		*algorithmName = util.FormatAlgorithmName(*algorithmName)

		message, err := util.GetMessageByHash(*algorithmName, *hash)

		if err == nil && message != "" {
			fmt.Println("Hash found! Message:")
			fmt.Println(message)
		} else {
			fmt.Println("Hash not found.")
		}

	case "create-rainbow":
		fileName := flag.String("file", "", "Input file name")
		algorithmName := flag.String("alg", "", "Hashing algorithm")
		flag.CommandLine.Parse(args[1:])

		inserter.CreateRainbow(*fileName, *algorithmName)
	case "check-rainbow":
		hash := flag.String("hash", "", "Input file name")
		algorithmName := flag.String("alg", "", "Hashing algorithm")
		flag.CommandLine.Parse(args[1:])

		*algorithmName = util.FormatAlgorithmName(*algorithmName)

		opt, err := inserter.CheckIfHashInRainbow(*algorithmName, *hash)
		util.CheckErrorPanic(err)

		opt.IfSet(func() {
			fmt.Println("Hash found! Message:")
			fmt.Println(opt.Value)
		}).IfNotSet(func() {
			fmt.Println("Hash not found.")
		})
	case "start-api":
		port := flag.Int("p", 3000, "Port to listen on")
		flag.CommandLine.Parse(args[1:])

		api.StartAPI(*port)
	default:
		fmt.Printf("Invalid command. %s is not a valid command.\n", command)
	}
}
