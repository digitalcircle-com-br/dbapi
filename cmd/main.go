package main

import "github.com/digitalcircle-com-br/dbapi/lib"

func main() {
	err := lib.Run()
	if err != nil {
		panic(err)
	}
}
