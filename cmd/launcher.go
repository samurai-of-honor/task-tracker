package main

import (
	"bufio"
	"flag"
	"os"
	tm "task-manager"
)

var dbFile = flag.String("f", "", "File with database or desired file name")

func main() {
	flag.Parse()
	var db string
	if *dbFile != "" {
		db = *dbFile
	} else {
		w := bufio.NewWriter(os.Stdout)
		_, _ = w.WriteString("Enter database file or skip for default: ")
		_ = w.Flush()
		db = tm.ReadStr(bufio.NewReader(os.Stdin))
		if db == "" {
			db = "db.json"
		}
	}
	sl := tm.Create()
	sl.Load(db)
	tm.Menu(sl, db)
}
