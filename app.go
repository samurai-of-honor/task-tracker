package main

import (
	"fmt"
)

func Menu(sl *[]Tasks) {
	var n int
	for {
		fmt.Print(`Enter options number:
  1 Show uncompleted tasks
  2 Show all tasks
  3 Mark the task completed
  4 Add task
  5 Change task
  6 Show overdue tasks
  7 Delete task
  8 Save changes
`)

		fmt.Scanln(&n)
		switch n {
		case 1:
			// ShowUncompleted()
		case 2:
			ShowAll(sl)
		case 3:
			// Mark()
		case 4:
			// Add()
		case 5:
			// Change()
		case 6:
			// ShowOverdue()
		case 7:
			// Delete()
		case 8:
			Save(sl)
		}
	}
}

func main() {
	sl := New()
	Add(sl, "1", "test1", "26.05.2022 21:00", false, "-")
	Add(sl, "2", "test2", "26.05.2022 22:00", true, "-")
	Menu(sl)

	/*
		dLine, err := time.Parse("02-01-2006 15:04", inDeadline)
		if err != nil {
			panic(err)
		}
		deadline := dLine.Format("02-01-2006 15:04")
	*/
}
