package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	util "hashlookup/util"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

func createHashesFromFileContent(filePath string, algorithmName string) {
	hashesToInsert := []util.HashesToInsertType{}

	file, err := os.Open(filePath)
	util.CheckErrorPanic(err)

	algorithmId, err := util.GetAlgorithmIdFromName(algorithmName)
	util.CheckErrorPanic(err)

	log.Printf("Inserting for algorithm %s id %d\n", algorithmName, algorithmId)

	lineCount := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineCount += 1
		line := scanner.Text()
		line = strings.TrimSpace(line)
		hash := util.MD5Hash(line)

		hashesToInsert = append(hashesToInsert, util.HashesToInsertType{
			Id:     algorithmId,
			Input:  line,
			Output: hash,
		})

		if len(hashesToInsert) >= 1000 {
			util.InsertMultipleAlgorithmHashes(hashesToInsert)
			hashesToInsert = []util.HashesToInsertType{}
		}

		if lineCount%1000 == 1 {
			log.Printf("Processed %d lines\n", lineCount)
		}
	}

	if len(hashesToInsert) > 0 {
		util.InsertMultipleAlgorithmHashes(hashesToInsert)
	}

}

func main() {
	log.Println("Running...")
	godotenv.Load()

	//filePath := "./data/passwords/10-million-password-list-top-1000000.txt"
	filePath := "./data/passwords/xato-net-10-million-passwords.txt"
	algorithmName := "MD-5"

	createHashesFromFileContent(filePath, algorithmName)
	log.Println("Finished.")
}
