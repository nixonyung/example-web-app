package servicebase

import "log"

func HandleErr(err error) {
	if err != nil {
		log.Fatalf("HandleErr: %s", err)
	}
}
