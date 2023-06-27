package main

import "log"

func ConfigureLogger() {
	log.SetFlags(log.Lmicroseconds | log.Llongfile)
}
