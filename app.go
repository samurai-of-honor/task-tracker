package main

import (
	"bufio"
	"fmt"
	"os"
)

func printMenu() {
	fmt.Print(`Options:
  1 Show uncompleted tasks
  2 Show all tasks
  3 Mark the task completed
  4 Add task
  5 Change task
  6 Show overdue tasks
  7 Delete task
  8 Save changes
  9 Menu
  0 Exit
`)
}

func readArgs(r *bufio.Reader) string {
	str, _, err := r.ReadLine()
	if err != nil {
		panic(err)
	}
	return string(str)
}

func Menu(sl *[]Tasks) {
	var n int
	printMenu()
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
			fmt.Print("Enter task title for mark: ")
			// Mark(sl, readArgs(r))
		case 4:
			// Add()
		case 5:
			// Change()
		case 6:
			ShowOverdue(sl)
		case 7:
			fmt.Print("Enter task title for delete: ")
			Delete(sl, readArgs(r))
		case 8:
			Save(sl)
		case 9:
			printMenu()
		case 0:
			os.Exit(0)
		}
	}
}

func main() {
	sl := Create()
	Load(sl)
	// Add(sl, "1", "test1", "02-01-2006 15:04", false, "-")
	Menu(sl)
}
