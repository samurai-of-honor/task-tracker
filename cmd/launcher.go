package main

import (
	tm "task-manager"
)

func main() {
	sl := tm.Create()
	tm.Load(sl)
	tm.Menu(sl)
}
