package main

import (
	"fmt"
	"reflect"
)

type Demo struct {
	Name string
}

func main() {
	d := &Demo{}

	fmt.Println(reflect.TypeOf(d).Kind())
	fmt.Println(reflect.ValueOf(d).Kind())
	fmt.Println(reflect.ValueOf(d).Type())
}
