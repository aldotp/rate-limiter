package main

import "github.com/aldotp/rate-limiter/cmd"

func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
