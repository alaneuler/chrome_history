package main

import (
	"os"

	"me.alaneuler/chrome_history/history"
)

func main() {
	if os.Getenv("alfred_workflow_bundleid") != "" {
		run()
	} else {
		history.Query("")
	}
}
