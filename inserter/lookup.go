package inserter

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	util "hashlookup/util"

	_ "github.com/lib/pq"
)

func createHashesFromFileContent(filePath string, algorithmName string, hashFn util.HashFunc) {
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
		hash := hashFn(line)

		hashesToInsert = append(hashesToInsert, util.HashesToInsertType{
			Id:     algorithmId,
			Input:  line,
			Output: hash,
		})

		if len(hashesToInsert) >= 10000 {
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

func CreateLookup(filePath string, algorithmName string) {
	hashFn, err := util.GetHashFuncFromName(algorithmName)
	if err != nil {
		fmt.Println(err)
		return
	}
	startTime := time.Now()
	createHashesFromFileContent(filePath, algorithmName, hashFn)
	log.Printf("Finished. Creating lookup table took %s.\n", time.Since(startTime))
}
