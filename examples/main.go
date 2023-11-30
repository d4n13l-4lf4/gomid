package main

import (
	"log"
	"reflect"
)

func main() {
	a := 1
	b := func(a int) {
		log.Println(a)
	}

	aVal := reflect.ValueOf(&a)
	bVal := reflect.ValueOf(b)

	log.Println(aVal.Type().Kind())
	log.Println(bVal.Type().Kind())
}
