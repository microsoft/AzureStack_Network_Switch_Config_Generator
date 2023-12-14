package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
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

func subnetToIPMask(subnet string) (string, int) {
	ip, ipnet, err := net.ParseCIDR(subnet)
	if err != nil {
		log.Fatalln(err)
	}
	mask, _ := ipnet.Mask.Size()
	return ip.String(), mask
}

func subnetToIPRange(subnet string) (string, string, string, string, int) {
	ip, ipnet, err := net.ParseCIDR(subnet)
	if err != nil {
		log.Fatalln(err)
	}
	mask, _ := ipnet.Mask.Size()
	firstIP := ip.Mask(ipnet.Mask)
	lastIP := net.IP(make([]byte, 4))
	for i := 0; i < 4; i++ {
		lastIP[i] = firstIP[i] | ^ipnet.Mask[i]
	}
	ones, bits := ipnet.Mask.Size()
	availableIPs := (1 << uint(bits-ones)) - 2
	firstHost := net.IP(make([]byte, 4))
	copy(firstHost, firstIP)
	lastHost := net.IP(make([]byte, 4))
	copy(lastHost, lastIP)
	if availableIPs > 0 {
		firstHost[3]++
		lastHost[3]--
	}
	return firstIP.String(), lastIP.String(), firstHost.String(), lastHost.String(), mask
}

// optimizeArray takes an array of integers and returns a string in the format "1-4,6,8-10,100"
// // Example usage
// arr := []int{1, 2, 3, 4, 6, 8, 9, 10, 100}
// result := optimizeArray(arr)
// fmt.Println(result) // Output: "1-4,6,8-10,100"

func optimizeArray(arr []int) string {
	var result strings.Builder
	start := arr[0]
	end := arr[0]

	// Iterate over the array and keep track of the start and end of each range of consecutive numbers
	for i := 1; i < len(arr); i++ {
		if arr[i] == end+1 {
			end = arr[i]
		} else {
			// If a gap is encountered, add the current range to the result string and start a new range
			if start == end {
				result.WriteString(strconv.Itoa(start) + ",")
			} else {
				result.WriteString(strconv.Itoa(start) + "-" + strconv.Itoa(end) + ",")
			}
			start = arr[i]
			end = arr[i]
		}
	}

	// Add the last range to the result string
	if start == end {
		result.WriteString(strconv.Itoa(start))
	} else {
		result.WriteString(strconv.Itoa(start) + "-" + strconv.Itoa(end))
	}

	return result.String()
}
