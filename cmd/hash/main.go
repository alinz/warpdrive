package main

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

var (
	flags = flag.NewFlagSet("passhash", flag.ExitOnError)
	pass  = flags.String("pass", "", "password")
	salt  = flags.String("salt", "", "salt value")
)

func main() {
	flags.Parse(os.Args[1:])

	fmt.Printf("Password: %s\n", *pass)
	fmt.Printf("salt: %s\n", *salt)

	hashpass, err := bcrypt.GenerateFromPassword([]byte(*pass+*salt), 0)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(hashpass))
}
