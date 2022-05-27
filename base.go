package control

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Tasks struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Deadline     string `json:"deadline"`
	Complete     bool   `json:"done"`
	CompleteDate string `json:"completeDate"`
}

func Create() *[]Tasks {
	sl := make([]Tasks, 0)
	return &sl
}

func separators() {
	fmt.Println(strings.Repeat("-", 50))
}

func timeParser(strDate string) time.Time {
	date, err := time.Parse("02-01-2006 15:04", strDate)
	if err != nil {
		panic(err)
	}
	return date
}

//------------ CHANGES -----------------

func Add(sl *[]Tasks, title, desc, dLine string) {
	timeParser(dLine)
	t := Tasks{title, desc, dLine, false, strings.Repeat("_", 16)}
	*sl = append(*sl, t)
}

func Delete(sl *[]Tasks, title string) {
	st := *sl
	for i, val := range st {
		if val.Title == title {
			st[0], st[i] = st[i], st[0]
			st = st[1:]
		}
	}
	*sl = st
}

func Mark(sl *[]Tasks, title string) {
	st := *sl
	now := time.Now().Format("02-01-2006 15:04")

	for i, val := range st {
		if val.Title == title {
			st[i].Complete = true
			st[i].CompleteDate = now
		}
	}
	st = *sl
}

func Change(sl *[]Tasks, title, newTitle, desc, dLine string) {
	st := *sl
	for i, val := range st {
		if val.Title == title {
			if newTitle != "" {
				st[i].Title = newTitle
			}
			if desc != "" {
				st[i].Description = desc
			}
			if dLine != "" {
				deadline := timeParser(dLine).Format("02-01-2006 15:04")
				st[i].Deadline = deadline
			}
		}
	}
	*sl = st
}

//------------ SHOWS -----------------

func ShowAll(sl *[]Tasks) {
	var mark string
	for i, val := range *sl {
		if val.Complete == true {
			mark = "X"
		} else {
			mark = " "
		}
		fmt.Printf("%d. %s  %s  [%s] %s\n%s\n", i+1, val.Deadline, val.Title, mark, val.CompleteDate, val.Description)
		separators()
	}
}

func ShowUncompleted(sl *[]Tasks) {
	st := *sl
	now := time.Now()
	for _, val := range st {
		dLine := timeParser(val.Deadline)
		if val.Complete == true || dLine.Before(now) {
			Delete(&st, val.Title)
		}
	}
	if len(st) == 0 {
		return
	}

	for i := 0; i < len(st)-1; i++ {
		for j := 0; j < len(st)-i-1; j++ {
			curDLine := timeParser(st[j].Deadline)
			nextDLine := timeParser(st[j+1].Deadline)
			if curDLine.After(nextDLine) {
				st[j], st[j+1] = st[j+1], st[j]
			}
		}
	}

	for _, val := range st {
		fmt.Printf(" %s   %s\n%s\n", val.Deadline, val.Title, val.Description)
		separators()
	}
}

func ShowOverdue(sl *[]Tasks) {
	now := time.Now()
	for _, val := range *sl {
		dLine := timeParser(val.Deadline)
		if val.Complete == false && dLine.Before(now) {
			fmt.Printf(" %s    %s\n%s\n", val.Title, val.Deadline, val.Description)
			separators()
		}
	}
}

//------------ DB CONTROL -----------------

func Load(sl *[]Tasks) {
	file, openErr := os.Open("db.json")
	if openErr != nil {
		Save(sl)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if len(data) == 0 {
		return
	}
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, sl); err != nil {
		panic(err)
	}

}

func Save(sl *[]Tasks) {
	data, convErr := json.MarshalIndent(sl, "", "  ")
	if convErr != nil {
		panic(convErr)
	}

	file, openErr := os.OpenFile("db.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if openErr != nil {
		panic(openErr)
	}

	w := io.Writer(file)
	if _, writeErr := w.Write(data); writeErr != nil {
		panic(writeErr)
	}
	fmt.Println("Saved!")
}
