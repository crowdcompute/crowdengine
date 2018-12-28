package common

import "log"

func CheckErr(err error, message string) {
	if err != nil {
		log.Println(message)
		log.Fatal("ERROR:", err)
	}
}
