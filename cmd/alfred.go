package main

import (
	"strconv"
	"strings"

	aw "github.com/deanishe/awgo"
	"me.alaneuler/chrome_history/history"
)

func runWithAlfred(wf *aw.Workflow) {
	entries := history.Query(strings.Join(wf.Args(), " "), 0)
	if len(entries) == 0 {
		wf.NewItem("No history entries found")
	} else {
		for _, entry := range entries {
			item := wf.NewItem(entry.Title)
			item.UID(strconv.FormatInt(entry.ID, 10))
			item.Subtitle(entry.URL)
			item.Arg(entry.URL)
			item.Valid(true)
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
