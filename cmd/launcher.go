package main

import (
	tm "task-manager"
)

func main() {
	sl := tm.Create()
	sl.Load()
	tm.Menu(sl)
}
