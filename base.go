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

func Show(i int, val Task) string {
	var mark string
	if val.Complete == true {
		mark = "X"
	} else {
		mark = " "
	}
	str := fmt.Sprintf("%d. %s  %s  [%s] %s\n%s\n\n", i+1, val.Deadline, val.Title,
		mark, val.CompleteDate, val.Description)
	return str
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

func (sl *SlTasks) ShowAll() string {
	var str string
	for i, val := range *sl {
		str += Show(i, val)
	}
	return str
}

func (sl *SlTasks) ShowUncompleted() string {
	st := *sl
	now := time.Now()
	for _, val := range st {
		dLine, _ := timeParser(val.Deadline)

		if val.Complete == true || dLine.Before(now) {
			st.Delete(val.Title)
		}
	}
	if len(st) == 0 {
		return "There's nothing"
	}

	for i := 0; i < len(st)-1; i++ {
		for j := 0; j < len(st)-i-1; j++ {
			curDLine, _ := timeParser(st[j].Deadline)
			nextDLine, _ := timeParser(st[j+1].Deadline)

			if curDLine.After(nextDLine) {
				st[j], st[j+1] = st[j+1], st[j]
			}
		}
	}

	var str string
	for i, val := range st {
		str += Show(i, val)
	}
	return str
}

func (sl *SlTasks) ShowOverdue() string {
	now := time.Now()
	var str string
	for i, val := range *sl {
		dLine, _ := timeParser(val.Deadline)
		if val.Complete == false && dLine.Before(now) {
			str += Show(i, val)
		}
	}
	if str == "" {
		return "There's nothing"
	}
	return str
}

//------------ DB CONTROL -----------------

func (sl *SlTasks) Find(title string) Task {
	for _, val := range *sl {
		if val.Title == title {
			return val
		}
	}
	return Task{Title: "undefined"}
}

func (sl *SlTasks) Load(db string) {
	file, openErr := os.Open(db)
	if openErr != nil {
		sl.Save(db)
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

func (sl *SlTasks) Save(db string) {
	data, convErr := json.MarshalIndent(sl, "", "  ")
	if convErr != nil {
		panic(convErr)
	}

	file, openErr := os.OpenFile(db, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if openErr != nil {
		panic(openErr)
	}

	w := io.Writer(file)
	if _, writeErr := w.Write(data); writeErr != nil {
		panic(writeErr)
	}
	fmt.Println("Saved!")
}
