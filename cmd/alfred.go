package main

import (
	"os"
	"strconv"
	"strings"

	aw "github.com/deanishe/awgo"
	"me.alaneuler/chrome_history/history"
)

func runWithAlfred(wf *aw.Workflow) {
	limit := 100
	limitConfig := os.Getenv("LIMIT")
	if limitConfig != "" {
		limit, _ = strconv.Atoi(limitConfig)
	}

	entries := history.Query(strings.Join(wf.Args(), " "), limit)
	if len(entries) == 0 {
		wf.NewItem("No history entries found")
	} else {
		for _, entry := range entries {
			item := wf.NewItem(entry.Title)
			item.UID(strconv.FormatInt(entry.ID, 10))
			item.Subtitle(entry.URL)
			item.Arg(entry.URL)
			item.Valid(true)
			item.Icon(entry.Icon)
		}
	}
	wf.SendFeedback()
}

func run() {
	wf := aw.New()
	wf.Run(func() {
		runWithAlfred(wf)
	})
}
