package main

import (
	"fmt"

	"me.alaneuler/chrome_history/history"
)

func main() {
	entries := history.Query()
	fmt.Println(entries)
}
