package main

import "github.com/Ablebil/eco-sample/internal/bootstrap"

func main() {
	if err := bootstrap.Start(); err != nil {
		panic(err)
	}
}
