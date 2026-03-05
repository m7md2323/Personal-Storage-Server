package auth

import (
	"crypto/rand"
	"fmt"
	"math/big"
)
var CurrentAccessCode string
//This function will print a 10 digits code on the laptop server terminal and the user should have access to the laptop,
// this will make only trusted people access the website.
func GenerateStorageCode() string {
	var code string
	for i := 0; i < 3; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(1000))
		code += fmt.Sprintf("%03d", num)
		if i < 2 {
			code += "-"
		}
	}
	
	lastDigit, _ := rand.Int(rand.Reader, big.NewInt(10))
	code += fmt.Sprintf("%d", lastDigit)

	fmt.Println("*********************************")
	fmt.Printf(" YOUR ACCESS CODE IS: %s \n", code)
	fmt.Println("*********************************")
	
	CurrentAccessCode = code
	return code
}