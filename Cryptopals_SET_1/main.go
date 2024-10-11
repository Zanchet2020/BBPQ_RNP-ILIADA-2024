package main

import (
	"fmt"
	"encoding/hex"
	"encoding/base64"
	"errors"
	"strings"
	"os"
	"bufio"
	"sync"
)

var word_weight_dictionary=map[string]uint{
        "of":142,
        " a ":131,
        "the":142,
        "oo":131,
        "is":131,
        "th":131,
        "ing":200,
	"and":200,
	"in":160,
	"tion":150,
	"to":170,
}
var char_weight_dictionary=map[string]uint{
        "e":116,
        "t":85,
        "a":82,
        "o":71,
        "i":70,
        "n":68,
        "s":64,
        "r":60,
        "h":47,
        "l":40,
        "d":37,
        "c":30,
        "u":25,
        "m":24,
        "f":20,
        "g":20,
        "p":19,
        "w":17,
        "y":16,
        "b":15,
        "v":10,
        "k":7,
        "j":2,
        "x":1,
        "z":1,
        "q":1,
}

func hex_to_base64(input string) string{
	// Decodificando hex
	bytes, err := hex.DecodeString(input)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	// Codificando para base64
	var output string = base64.StdEncoding.EncodeToString(bytes)
	return output
}


// Realiza XOR de duas slices de bytes
func xor_byte_slice(a []byte, b []byte) []byte{
	var output []byte = make([]byte, len(a))
	for i := range output{
		output[i] = a[i] ^ b[i]
	}
	return output
}

// Realiza o XOR de duas strings hex
func fixed_xor(input1 string, input2 string) (string, error) {
	bytes1, err1 := hex.DecodeString(input1)
	if err1 != nil {
		return "", err1
	}

	bytes2, err2 := hex.DecodeString(input2)
	if err2 != nil {
		return "", err2
	}

	if len(bytes1) != len(bytes2){
		return "", errors.New("Strings must have the same lenght")
	}

	var output []byte = xor_byte_slice(bytes1, bytes2)
	return hex.EncodeToString(output), nil
}

// Realiza XOR de todos os bytes de uma string com um byte único
func char_xor(input string, char int) string{
	bytes, err := hex.DecodeString(input)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	var output []byte = make([]byte, len(bytes))
	for i := 0; i < len(bytes); i++{
		output[i] = bytes[i] ^ byte(char)
	}
	return hex.EncodeToString(output)
}

type return_val struct{
	message string
	score uint
	key byte
}

func break_single_byte_XOR_cypher_routine(input string, key byte, return_list *[]return_val, wg *sync.WaitGroup, m *sync.Mutex){
	defer wg.Done()
	var xored_string string = char_xor(input, int(key))
	var bytes, _ = hex.DecodeString(xored_string)
	var str = string(bytes)
	var score uint
	// Dois loops para somar a "pontuação" da string
	// Somando a pontuação em relação às letras
	for index := range str{
		score = 0
		score += char_weight_dictionary[strings.ToLower(string(str[index]))]
	}
	// Somando a pontuação em relação aos fonemas
	for word, key := range word_weight_dictionary{
		if strings.Contains(str, word){
			score += key
		}
	}

	new := return_val{str, score, key}
	m.Lock()
	*return_list = append(*return_list, new)
	m.Unlock()
}

// Quebra uma codificação XOR de uma string através do teste da frequência de letras e fonemas do inglês
func break_single_byte_XOR_cypher(input string) (string, byte, uint){
	const byte_size int = 256
	var max uint = 0
	var message string
	var cypher byte
	// Loop por todos os bytes possíveis

	var m sync.Mutex
	var wg sync.WaitGroup
	wg.Add(byte_size)
	var return_list = make([]return_val, byte_size)
	
	for i := 0; i < byte_size; i++{
		go break_single_byte_XOR_cypher_routine(input, byte(i), &return_list, &wg, &m)
	}
	wg.Wait()
	for _, x := range return_list{
		if x.score > max{
			max = x.score
			message = x.message
			cypher = x.key
		}
	}
	
	return message, byte(cypher), max
}


func repeating_key_XOR_cypher(input string, key string) string {
	key_bytes := []byte(key)
	var output = make([]byte, len(input))
	for i, c := range input{
		output[i] = byte(c) ^ key_bytes[i % len(key)]
	}
	return hex.EncodeToString(output)
}


func main(){
	
	// Convert hex to base 64 ====================================================
	{
		fmt.Println("===========================")
		fmt.Println(">>>> Convert hex to base 64:")
		const input = "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
		fmt.Println(hex_to_base64(input), "\n===========================\n\n")
	}


	
	// Fixed XOR ================================================================
	{
		fmt.Println("===========================")
		fmt.Println(">>>> Fixed XOR")
		const input1 = "1c0111001f010100061a024b53535009181c"
		const input2 = "686974207468652062756c6c277320657965"
		var output, err = fixed_xor(input1, input2)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(output, "\n===========================\n\n")
	}

	
	// Single Byte XOR cypher =====================================================
	{
		fmt.Println("===========================")
		fmt.Println(">>>> Single Byte XOR Cypher")
		const hex_encoded_string = "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"
		message, cypher_key, _ := break_single_byte_XOR_cypher(hex_encoded_string)	
		fmt.Println("Message:", message, "\nCypher:", string(cypher_key))
		fmt.Println("===========================\n\n")
	}


	// Detect single character XOR
	{
		fmt.Println("===========================")
		fmt.Println(">>>> Detect single byte XOR Cypher")
		file, err := os.Open("4.txt")
		if err!=nil{
			fmt.Println(err)
		}

		scanner := bufio.NewScanner(file)
		var max_score uint = 0
		var message string
		var key byte
		for scanner.Scan(){
			msg, cypher_key, score := break_single_byte_XOR_cypher(scanner.Text())
			if score > max_score {
				message = msg
				key = cypher_key
				max_score = score
			}
		}

		file.Close()
		fmt.Println("Message:", strings.Trim(message, "\n "), "\nCypher:", string(key))
		fmt.Println("===========================\n\n")
	}


	// Implement repeating-key XOR
	{
		fmt.Println("===========================")
		fmt.Println(">>>> Implement repeating-key XOR")
		var output string = repeating_key_XOR_cypher(`Burning 'em, if you ain't quick and nimble
			I go crazy when I hear a cymbal`, "ICE")
		fmt.Println(output)
		fmt.Println("===========================\n\n")
	}
}
