package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const (
	domain = "https://ciphersprint.pulley.com"
	email  = "ojaswadhwani098@gmail.com"
)

type Challenge struct {
	Challenger       string `json:"challenger"`
	EncryptedPath    string `json:"encrypted_path"`
	EncryptionMethod string `json:"encryption_method"`
	ExpiresIn        string `json:"expires_in"`
	Hint             string `json:"hint"`
	Instructions     string `json:"instructions"`
	Level            int    `json:"level"`
}

func MakeGetRequest(url string) (Challenge, error) {
	// Make the GET request
	response, err := http.Get(url)
	if err != nil {
		return Challenge{}, err
	}
	defer response.Body.Close()

	// Read the response body
	var body Challenge
	err = json.NewDecoder(response.Body).Decode(&body)
	if err != nil {
		return Challenge{}, err
	}

	fmt.Println("result", body)

	return body, nil
}

// Challenge: converted to a JSON array of ASCII
func convertJSONASCIIArraytoString(encrypted_path string) (string, error) {
	// encrypted_path := "task_[55,56,97,55,101,100,51,54,51,54,52,102,98,100,51,57,52,99,102,97,48,99,56,101,99,50,100,55,57,56,48,51]"
	arrayString := strings.TrimPrefix(encrypted_path, "task_")

	numbersStr := strings.Trim(arrayString, "[]")
	numbers := strings.Split(numbersStr, ",")

	var letters []string
	for _, numStr := range numbers {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			fmt.Println("Error converting number:", err)
			return "", err
		}
		letter := string(num)
		letters = append(letters, letter)
	}

	result := strings.Join(letters, "")
	// fmt.Println("result", fmt.Sprintf("task_%s", result))
	return fmt.Sprintf("task_%s", result), nil
}

// Challenge: inserted some non-hex characters: task_4ea91aj110l447h6ba7k439i5gb9e6e8cb015e
func removeNonHex(input string) string {
	nonHexString := strings.TrimPrefix(input, "task_")

	// Define a regular expression to match hexadecimal characters
	re := regexp.MustCompile("[^0-9a-fA-F]+")

	// Remove non-hex characters from the input string
	return fmt.Sprintf("task_%s", re.ReplaceAllString(nonHexString, ""))
}

// Challenge: added x to ASCII value of each character: task_)*.],]-0,-1()(+0,,]^+\]-\*[Y)Y]]
func addXToASCII(input string, buffer int) string {
	trimmedString := strings.TrimPrefix(input, "task_")

	result := ""
	for _, char := range trimmedString {
		newChar := char - rune(buffer)
		result += string(newChar)
	}

	return fmt.Sprintf("task_%s", result)
}

func extractBufferFromMethod(method string) int {
	words := strings.Fields(method)

	number, _ := strconv.Atoi(words[1])

	return number
}

// hex decoded, encrypted with XOR, hex encoded again: task_2d7ea6d5368f39e840bc6c8c41fe3c25
// key seems to be secret
func xorDecrypt(data []byte, key []byte) []byte {
	decrypted := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		decrypted[i] = data[i] ^ key[i%len(key)]
	}
	return decrypted
}

func decodeDecryptEncode(input string) (string, error) {
	trimmedString := strings.TrimPrefix(input, "task_")

	// Decode with hex
	decodedBytes, err := hex.DecodeString(trimmedString)
	if err != nil {
		return "", err
	}

	fmt.Println("decodedBytes", decodedBytes)

	keyStr := "secret"

	// Decrypt with XOR
	decryptedData := xorDecrypt(decodedBytes, []byte(keyStr))

	fmt.Println("decryptedData", decryptedData)

	// Encode with hex
	originalString := hex.EncodeToString(decryptedData)

	fmt.Println("originalString", originalString)
	return fmt.Sprintf("task_%s", originalString), nil
}

// task_1834beed3b41e5e333a6d4d8512f742b scrambled! original positions as base64 encoded messagepack: 3AAgAxcMCgEWBhoICQACHxwUHRELDRkPBRgEHhITBxsQFQ4=

func main() {
	fmt.Println("Hey There, Pulley!")
	//https://ciphersprint.pulley.com/task_7f3671e3cf343511fe14a4f81f8dd50

	// fmt.Println("request url", fmt.Sprintf("%s%s", domain, email))
	firstChallenge, err := MakeGetRequest(fmt.Sprintf("%s/%s", domain, email))
	if err != nil {
		panic(err)
	}

	encrypted_path := firstChallenge.EncryptedPath
	path := encrypted_path

	secondChallenge, err := MakeGetRequest(fmt.Sprintf("%s/%s", domain, path))
	if err != nil {
		panic(err)
	}

	encrypted_path = secondChallenge.EncryptedPath
	path, err = convertJSONASCIIArraytoString(encrypted_path)
	if err != nil {
		panic(err)
	}

	thirdChallenge, err := MakeGetRequest(fmt.Sprintf("%s/%s", domain, path))
	if err != nil {
		panic(err)
	}

	path = removeNonHex(thirdChallenge.EncryptedPath)
	fourthChallenge, err := MakeGetRequest(fmt.Sprintf("%s/%s", domain, path))
	if err != nil {
		panic(err)
	}

	buffer := extractBufferFromMethod(fourthChallenge.EncryptionMethod)
	path = addXToASCII(fourthChallenge.EncryptedPath, buffer)

	// fmt.Printf("buffer and path: %d %s", buffer, path)
	fifthChallenge, err := MakeGetRequest(fmt.Sprintf("%s/%s", domain, path))
	if err != nil {
		panic(err)
	}

	path, err = decodeDecryptEncode(fifthChallenge.EncryptedPath)
	sixthChallenge, err := MakeGetRequest(fmt.Sprintf("%s/%s", domain, path))
	if err != nil {
		panic(err)
	}

	fmt.Println("sixthChallenge", sixthChallenge)
}
