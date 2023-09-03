package inserter

import (
	"bufio"
	"fmt"
	util "hashlookup/util"
	"hashlookup/util/optional"
	"log"
	"math/big"
	"os"
	"strings"
	"time"
)

var chainLength = 1000
var alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"

func Reduce(hash string, chainIndex int) string {
	hashInt := new(big.Int)
	hashInt.SetString(util.SHA256Hash(hash), 16)

	maxLen := 20
	minLen := 6

	keyLength := maxLen - chainIndex
	if keyLength < minLen {
		keyLength = minLen
	} else if keyLength > maxLen {
		keyLength = maxLen
	}

	s := ""
	for i, b := range hashInt.Bytes() {
		if i >= keyLength {
			break
		}
		s = s + string(alphabet[int(b)%len(alphabet)])
	}

	return s
}

func createRainbowFromFileContent(filePath string, algorithmName string, hashFn util.HashFunc) {
	file, err := os.Open(filePath)
	util.CheckErrorPanic(err)

	algorithmId, err := util.GetAlgorithmIdFromName(algorithmName)
	util.CheckErrorPanic(err)

	log.Printf("Inserting for algorithm %s id %d\n", algorithmName, algorithmId)

	var keyList []string

	lineCount := 0
	scanner := bufio.NewScanner(file)
	log.Println("Reading file")
	for scanner.Scan() {
		lineCount += 1
		line := scanner.Text()
		line = strings.TrimSpace(line)
		keyList = append(keyList, line)
	}

	log.Printf("Creating rainbow table with chain length %d\n", chainLength)
	hashesToInsert := []util.HashesToInsertType{}

	reduceTimer := time.Now()
	for keyIndex, key := range keyList {
		// fmt.Printf("Processing key %s\n", key)
		reducedHash := key
		for i := 0; i < chainLength; i++ {
			hashedKey := hashFn(reducedHash)
			reducedHash = Reduce(hashedKey, i)
		}

		hashesToInsert = append(hashesToInsert, util.HashesToInsertType{
			Id:     algorithmId,
			Input:  key,
			Output: reducedHash,
		})

		if len(hashesToInsert) >= 10000 {
			util.InsertMultipleRainbowEntries(hashesToInsert)
			hashesToInsert = []util.HashesToInsertType{}
		}

		if (keyIndex+1)%100 == 0 {
			log.Printf("Processed %d keys\n", keyIndex)
		}
	}

	if len(hashesToInsert) > 0 {
		util.InsertMultipleRainbowEntries(hashesToInsert)
	}

	fmt.Println(time.Since(reduceTimer))
}

func FindHashInChain(hashFn util.HashFunc, chainInput string, findHash string) *optional.Optional[string] {
	reducedHash := chainInput
	for i := 0; i < chainLength; i++ {
		hashedKey := hashFn(reducedHash)
		if hashedKey == findHash {
			return optional.Create[string](reducedHash)
		}
		reducedHash = Reduce(hashedKey, i)
	}

	return optional.Empty[string]()
}

func CreateRainbow(filePath string, algorithmName string) {
	hashFn, err := util.GetHashFuncFromName(algorithmName)
	if err != nil {
		fmt.Println(err)
		return
	}
	startTime := time.Now()
	createRainbowFromFileContent(filePath, algorithmName, hashFn)
	log.Printf("Finished. Creating rainbow table took %s.\n", time.Since(startTime))
}

func CheckIfHashInRainbow(algorithmName string, hash string) (*optional.Optional[string], error) {
	hashFn, err := util.GetHashFuncFromName(algorithmName)

	if err != nil {
		return optional.Empty[string](), fmt.Errorf("invalid algorithm name %s", algorithmName)
	}

	var possibleKeys []string

	for offset := 0; offset < chainLength-1; offset++ {
		hashedKey := hash
		reducedHash := ""
		for i := offset; i < chainLength; i++ {
			reducedHash = Reduce(hashedKey, i)
			hashedKey = hashFn(reducedHash)
		}
		possibleKeys = append(possibleKeys, reducedHash)
	}

	inputs, err := util.GetRainbowInputsFromOutputs(algorithmName, possibleKeys)

	if err != nil {
		return optional.Empty[string](), err
	}

	for _, input := range *inputs {
		opt := FindHashInChain(hashFn, input, hash)
		if opt.IsValueSet {
			return opt, nil
		}
	}

	return optional.Empty[string](), nil
}

func CheckIfMessageInRainbow(algorithmName string, message string) (*optional.Optional[string], error) {
	hashFn, err := util.GetHashFuncFromName(algorithmName)

	if err != nil {
		return optional.Empty[string](), fmt.Errorf("invalid algorithm name %s", algorithmName)
	}

	return CheckIfHashInRainbow(algorithmName, hashFn(message))
}
