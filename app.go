package main

import (
	"bufio"
	"fmt"
	"os"
)

func Menu(sl *[]Tasks) {
	var n int
	fmt.Print(`Options:
  1 Show uncompleted tasks
  2 Show all tasks
  3 Mark the task completed
  4 Add task
  5 Change task
  6 Show overdue tasks
  7 Delete task
  8 Save changes
  0 Exit
`)
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter options number: ")
		fmt.Scanln(&n)
		switch n {
		case 1:
			ShowUncompleted(sl)
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
			fmt.Print("Enter task title for delete: ")
			title, _, err := r.ReadLine()
			if err != nil {
				panic(err)
			}
			Delete(sl, string(title))
		case 8:
			Save(sl)
		case 0:
			os.Exit(0)
		}
	}
}

func main() {
	sl := New()
	Load(sl)
	// Add(sl, "1", "test1", "02-01-2006 15:04", false, "-")
	Menu(sl)
}
