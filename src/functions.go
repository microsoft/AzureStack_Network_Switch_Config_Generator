package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
)

func parseInputJson(inputJsonFile string) InputData {
	inputObj := &InputType{}
	bytes, err := ioutil.ReadFile(inputJsonFile)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(bytes, inputObj)
	if err != nil {
		log.Fatalln(err)
	}
	return inputObj.InputData
}

func generateRandomString(stringLength, minUpperChar, minNum, minSpecialChar int) string {

	specialCharSet := "!@#$%&*"
	lowerCharSet := listCharSet('a', 'z')
	upperCharSet := listCharSet('A', 'Z')
	numberSet := listCharSet('0', '9')

	remainLength := stringLength - minUpperChar - minNum - minSpecialChar
	if remainLength < 0 {
		log.Fatalln("[Error] Please input valid length for string to be generated.")
	}

	// Build the string
	var randString strings.Builder
	randCharStrings(&randString, upperCharSet, minUpperChar)
	randCharStrings(&randString, numberSet, minNum)
	randCharStrings(&randString, specialCharSet, minSpecialChar)
	randCharStrings(&randString, lowerCharSet, remainLength)

	// Shuffle the string
	randRune := []rune(randString.String())
	rand.Shuffle(len(randRune), func(i, j int) {
		randRune[i], randRune[j] = randRune[j], randRune[i]
	})

	return string(randRune)
}

func listCharSet(start, end rune) string {
	var charSet strings.Builder
	for i := start; i <= end; i++ {
		fmt.Fprint(&charSet, string(i))
	}
	return charSet.String()
}

func randCharStrings(randString *strings.Builder, charset string, minCharNum int) {
	if minCharNum == 0 || len(charset) == 0 {
		return
	}
	for i := 0; i < minCharNum; i++ {
		randNum := rand.Intn(len(charset))
		randString.WriteString(string(charset[randNum]))
	}
}

func createFolder(folderPath string) {
	_, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(folderPath, 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	}
}
