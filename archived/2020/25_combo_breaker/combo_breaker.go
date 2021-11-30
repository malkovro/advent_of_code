package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func Encrypt(value int, subject int, loop int) (encryptedKey int) {
	encryptedKey = value
	for i := 0; i < loop; i++ {
		encryptedKey = encryptedKey * subject % 20201227
	}
	return
}

var debug = false

func main() {
	defer timeTrack(time.Now(), "Solving Time")

	doorPublicKey := flag.Int("door", -1, "the door public key")
	cardPublicKey := flag.Int("card", -1, "the card public key")
	flag.Parse()

	if *doorPublicKey == -1 || *cardPublicKey == -1 {
		panic("Please provide a door and card public key")
	}

	subject := 7

	doorLoop := -1
	encrypted := 1

	loop := 0
	for doorLoop == -1 {
		loop += 1
		encrypted = Encrypt(encrypted, subject, 1)

		if encrypted == *doorPublicKey {
			doorLoop = loop
		}
	}

	fmt.Println("Encryption key:", Encrypt(1, *cardPublicKey, doorLoop))
}
