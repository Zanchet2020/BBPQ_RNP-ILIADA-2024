package main

import (
	"fmt"
	"encoding/hex"
	"encoding/base64"
	"errors"
	"strings"
	"os"
	"bufio"
	//	"sort"
)

// var word_weight_dictionary = map[string]uint{
// 	" of ": 650000000,
// 	" a ": 600000000,
// 	" the ": 650000000,
// 	"oo": 600000000,
// 	" is ":600000000,
// 	"th":600000000,
// 	"ing ":700000000,
// }
// var char_weight_dictionary = map[string]uint{
// 	"e": 529117365,
// 	"t": 390965105,
// 	"a": 374061888,
// 	"o": 326627740,
// 	"i": 320410057,
// 	"n": 313720540,
// 	"s": 294300210,
// 	"r": 277000841,
// 	"h": 216768975,
// 	"l": 183996130,
// 	"d": 169330528,
// 	"c": 138416451,
// 	"u": 117295780,
// 	"m": 110504544,
// 	"f": 95422055,
// 	"g": 91258980,
// 	"p": 90376747,
// 	"w": 79843664,
// 	"y": 75294515,
// 	"b": 70195826,
// 	"v": 46337161,
// 	"k": 35373464,
// 	"j": 9613410,
// 	"x": 8369915,
// 	"z": 4975847,
// 	"q": 4550166,
// }


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
	bytes, err := hex.DecodeString(input)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	var output string = base64.StdEncoding.EncodeToString(bytes)
	return output
}

func xor_byte_slice(a []byte, b []byte) []byte{
	var output []byte = make([]byte, len(a))
	for i := range output{
		output[i] = a[i] ^ b[i]
	}
	return output
}

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


func break_single_byte_XOR_cypher(input string) (string, byte, uint){
	var scores = make([]uint, 256)
	var max uint = 0
	var message string
	var cypher byte
	for i := 0; i < 256; i++{
		var bytes, _ = hex.DecodeString(char_xor(input, i))
		var str = string(bytes)
		for index := range str{
			scores[i] = 0
			scores[i] += char_weight_dictionary[strings.ToLower(string(str[index]))]
		}
		for word, key := range word_weight_dictionary{
			if strings.Contains(str, word){
				scores[i] += key
			}
		}
		if scores[i] > max {
			max = scores[i]
			message = str
			cypher = byte(i)
		}
	}
	return message, byte(cypher), max
}



// Cool litte algorithm to count the number of bits in a uint32 in constant time found in this website:
// https://web.archive.org/web/20151229003112/http://blogs.msdn.com/b/jeuge/archive/2005/06/08/hakmem-bit-count.aspx
 func BitCount(n uint32) int{
	var count uint32;
    
	count = n - ((n >> 1) & 033333333333) - ((n >> 2) & 011111111111)
	return int(((count + (count >> 3)) & 030707070707) % 63)
 }

func hamming_distance(a string, b string) int {
	var count int
	
	var xored []byte = xor_byte_slice([]byte(a), []byte(b))
	
	for _, byte := range xored{
		count += BitCount(uint32(byte))
	} 
	return count
}

func edit_distance(a string, b string) int {	
	return hamming_distance(a, b)
}

func repeating_key_XOR_cypher(input string, key string) string {
	key_bytes := []byte(key)
	var output = make([]byte, len(input))
	for i, c := range input{
		output[i] = byte(c) ^ key_bytes[i % len(key)]
	}
	return hex.EncodeToString(output)
}



// func break_repeating_XOR(input string) (string, string){
// 	input_bytes := []byte(input)

// 	type distance_keysize_pair struct{
// 		distance float32
// 		keysize int
// 	}

// 	const max_keysize = 40
// 	const min_keysize = 2
// 	const keysize_guesses int = 3
// 	var norm_distances = make([]distance_keysize_pair, max_keysize - min_keysize)
	
// 	// Getting lowest distance
// 	for KEYSIZE := min_keysize; KEYSIZE < max_keysize; KEYSIZE++{
// 		var d1 int = edit_distance(string(input_bytes[:KEYSIZE]), string(input_bytes[KEYSIZE:2*KEYSIZE]))
// 		var d2 int = edit_distance(string(input_bytes[KEYSIZE:2*KEYSIZE]), string(input_bytes[2*KEYSIZE:3*KEYSIZE]))
// 		var d3 int = edit_distance(string(input_bytes[2*KEYSIZE:3*KEYSIZE]), string(input_bytes[3*KEYSIZE:4*KEYSIZE]))
// 		norm_distances[KEYSIZE - min_keysize].distance = float32(d1 + d2 + d3) / float32(KEYSIZE * 3)
// 		norm_distances[KEYSIZE - min_keysize].keysize = KEYSIZE + min_keysize 
// 	}

// 	var probable_keysizes = make([]distance_keysize_pair, keysize_guesses)
	
// 	for i, d := range distances{
// 		if probable_keysizes[0].distance == 0{
// 			probable_keysizes[0].distance = d
// 			probable_keysizes[0].keysize = i + min_keysize
// 		} else if d < probable_keysizes[keysize_guesses - 1].distance {
// 			probable_keysizes[keysize_guesses - 1].distance = d
// 			probable_keysizes[keysize_guesses - 1].keysize = i + min_keysize
// 		}

// 		sort.Slice(probable_keysizes, func (p, q int) bool{
// 			return probable_keysizes[p].distance < probable_keysizes[q].distance
// 		})
// 	}

// 	// dkp = distance_keysize_pair
// 	for _, dkp := range probable_keysizes{
// 		var key []byte = make([]byte, dkp.keysize)
// 		for i := 0; i < dkp.keysize; i++{
// 			var block []byte = make([]byte, len(input_bytes) / dkp.keysize)
// 			for j := i; j < len(input_bytes); j+=dkp.keysize{
// 				block = append(block, input_bytes[j])
// 			}
// 			_, cypher, _ := break_single_byte_XOR_cypher(hex.EncodeToString(block))
// 			//fmt.Println("Message:", message, "   Cypher:", cypher, "   Max:", max)
// 			key = append(key, cypher)
// 		}
// 		fmt.Println(string(key))
// 		res, _ := hex.DecodeString(repeating_key_XOR_cypher(input, string(key)))
// 		fmt.Println(string(res))
// 	}
	
// 	// fmt.Println(distances)

// 	// for _, d := range probable_keysizes{
// 	// 	fmt.Println(d.distance, d.keysize)
// 	// }
	
// 	// for i := range probable_keysizes{
// 	// 	probable_
// 	// }
	
// 	return "", ""
// }


func main(){
	{
		// Convert hex to base 64 ====================================================
		fmt.Println("===========================")
		fmt.Println(">>>> Convert hex to base 64:")
		var input = "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
		fmt.Println(hex_to_base64(input), "\n===========================\n\n")
	}


	
	// Fixed XOR ================================================================
	{
		fmt.Println("===========================")
		fmt.Println(">>>> Fixed XOR")
		var input1 = "1c0111001f010100061a024b53535009181c"
		var input2 = "686974207468652062756c6c277320657965"
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
	

	// // Break Repeating Key XOR
	// {
	// 	fmt.Println("===========================")
	// 	fmt.Println(">>>> Break Repeating Key XOR")
		
	// 	//fmt.Println("Lev Distance:", edit_distance("this is a test", "wokka wokka!!!"))

	// 	file_content, err := os.ReadFile("6.txt")
	// 	if err != nil{
	// 		fmt.Println(err)
	// 	}

	// 	b64_decoded, err := base64.StdEncoding.DecodeString(string(file_content))
	// 	if err != nil {
	// 		fmt.Println("Error decoding Base64:", err)
	// 		return
	// 	}
		
	// 	var decrypted, key string = break_repeating_XOR(string(b64_decoded))
	// 	fmt.Println("Decrypted message:\n", decrypted, "\nKey:", key)

	// 	fmt.Println("===========================\n\n")
	// }
}
