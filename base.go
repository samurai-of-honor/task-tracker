package control

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Task struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Deadline     string `json:"deadline"`
	Complete     bool   `json:"done"`
	CompleteDate string `json:"completeDate"`
}

type SlTasks []Task

func Create() *SlTasks {
	sl := new(SlTasks)
	return sl
}

func separators() {
	fmt.Println(strings.Repeat("-", 50))
}

func show(i int, val Task) {
	var mark string
	if val.Complete == true {
		mark = "X"
	} else {
		mark = " "
	}
	fmt.Printf("%d. %s  %s  [%s] %s\n%s\n", i+1, val.Deadline, val.Title, mark, val.CompleteDate, val.Description)
	separators()
}

var timeLayout = "02-01-2006 15:04"

func timeParser(strDate string) (time.Time, error) {
	date, err := time.Parse(timeLayout, strDate)
	if err != nil {
		fmt.Println("Time parse error!")
		return date, err
	}
	return date, nil
}

//------------ CHANGES -----------------

func (sl *SlTasks) Add(title, desc, dLine string) {
	if _, err := timeParser(dLine); err != nil {
		return
	}
	t := Task{title, desc, dLine, false, strings.Repeat("_", 16)}
	*sl = append(*sl, t)
}

func (sl *SlTasks) Delete(title string) {
	st := *sl
	for i, val := range st {
		if val.Title == title {
			st[0], st[i] = st[i], st[0]
			st = st[1:]
		}
	}
	*sl = st
}

func (sl *SlTasks) Mark(title string) {
	st := *sl
	now := time.Now().Format(timeLayout)

	for i, val := range st {
		if val.Title == title {
			st[i].Complete = true
			st[i].CompleteDate = now
		}
	}
	*sl = st
}

func (sl *SlTasks) Change(title, newTitle, desc, dLine string) {
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
				deadline, err := timeParser(dLine)
				if err != nil {
					return
				}
				st[i].Deadline = deadline.Format(timeLayout)
			}
		}
	}
	*sl = st
}

//------------ SHOWS -----------------

func (sl *SlTasks) ShowAll() {
	for i, val := range *sl {
		show(i, val)
	}
}

func (sl *SlTasks) ShowUncompleted() {
	st := *sl
	now := time.Now()
	for _, val := range st {
		dLine, err := timeParser(val.Deadline)
		if err != nil {
			return
		}
		if val.Complete == true || dLine.Before(now) {
			st.Delete(val.Title)
		}
	}
	if len(st) == 0 {
		return
	}

	for i := 0; i < len(st)-1; i++ {
		for j := 0; j < len(st)-i-1; j++ {
			curDLine, err1 := timeParser(st[j].Deadline)
			if err1 != nil {
				return
			}
			nextDLine, err2 := timeParser(st[j+1].Deadline)
			if err2 != nil {
				return
			}
			if curDLine.After(nextDLine) {
				st[j], st[j+1] = st[j+1], st[j]
			}
		}
	}

	for i, val := range st {
		show(i, val)
	}
}

func (sl *SlTasks) ShowOverdue() {
	now := time.Now()
	for i, val := range *sl {
		dLine, err := timeParser(val.Deadline)
		if err != nil {
			return
		}
		if val.Complete == false && dLine.Before(now) {
			show(i, val)
		}
	}
}

//------------ DB CONTROL -----------------

func (sl *SlTasks) Load() {
	file, openErr := os.Open("db.json")
	if openErr != nil {
		sl.Save()
		return
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

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

func (sl *SlTasks) Save() {
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
