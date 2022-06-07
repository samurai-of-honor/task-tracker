package main

import (
	"flag"
	tm "task-manager"
)

var (
	help   = flag.Bool("h", false, "Show help")
	dbFile = flag.String("f", "db.json", "File with database or desired file name")

	show            = flag.NewFlagSet("show", flag.PanicOnError)
	showAll         = show.Bool("all", false, "Show all tasks")
	showUncompleted = show.Bool("uncompleted", false, "Show uncompleted tasks")
	showOverdue     = show.Bool("overdue", false, "Show overdue tasks")
)

func main() {
	flag.Parse()
	if *help {
		flag.Usage()
		show.Usage()
		return
	}
	sl := tm.Create()
	sl.Load(*dbFile)
	switch {
	case *showAll == true:
		sl.ShowAll()
	case *showUncompleted == true:
		sl.ShowUncompleted()
	case *showOverdue == true:
		sl.ShowOverdue()
	}
}
