package main

import (
	"flag"
	tm "task-manager"
)

var dbFile = flag.String("f", "", "File with database or desired file name")

func main() {
	flag.Parse()
	var db string
	if *dbFile != "" {
		db = *dbFile
	} else {
		db = "db.json"
	}
	sl := tm.Create()
	sl.Load(db)
	tm.Menu(sl, db)
}
