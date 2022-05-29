package control

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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

func readStr(r *bufio.Reader) string {
	str, _, err := r.ReadLine()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(str)
}

func readArgs(r *bufio.Reader) []string {
	str := readStr(r)
	strSl := strings.Split(str, "/")
	if len(strSl) != 3 {
		fmt.Println("Input error!")
		return nil
	}
	return strSl
}

func Menu(sl *[]Tasks) {
	printMenu()
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter options number: ")
		n := readStr(r)
		separators()

		switch n {
		case "1":
			ShowUncompleted(sl)
		case "2":
			ShowAll(sl)
		case "3":
			fmt.Print("Enter task title for mark: ")
			Mark(sl, readStr(r))
		case "4":
			fmt.Print("Enter task info in format:\ntitle/description/01-01-2022 13:00\n")
			strSl := readArgs(r)
			if strSl == nil {
				break
			}
			Add(sl, strSl[0], strSl[1], strSl[2])
		case "5":
			fmt.Print("Enter task title for change: ")
			title := readStr(r)
			fmt.Print("Enter changed info and skip unchanged:\nExample for title: new title//\n")
			strSl := readArgs(r)
			if strSl == nil {
				break
			}
			Change(sl, title, strSl[0], strSl[1], strSl[2])
		case "6":
			ShowOverdue(sl)
		case "7":
			fmt.Print("Enter task title for delete: ")
			Delete(sl, readStr(r))
		case "8":
			Save(sl)
		case "9":
			printMenu()
		case "0":
			os.Exit(0)
		}
	}
}
