package main

import (
	"fmt"
	"os"

	"me.alaneuler/chrome_history/history"
)

func main() {
	if os.Getenv("alfred_workflow_bundleid") != "" {
		run()
	} else {
		entries := history.Query("", 5, false)
		for _, entry := range entries {
			fmt.Println(entry)
		}
	}
}
