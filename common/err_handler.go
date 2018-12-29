package common

import "log"

// CheckErr checks if an error is nil and logs and fatals it
func CheckErr(err error, message string) {
	if err != nil {
		log.Println(message)
		log.Fatal("ERROR:", err)
	}
}
