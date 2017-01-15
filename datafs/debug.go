package datafs

import "log"

func debug(args ...interface{}) {
	log.Println(args...)
}

func debugf(s string, args ...interface{}) {
	log.Printf(s, args...)
}
