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

func (sl *SlTasks) Find(str string) Task {
	for _, val := range *sl {
		title0 := strings.TrimPrefix(str, "🔹 ")
		title := strings.TrimPrefix(title0, "🔸 ")
		if val.Title == title {
			return val
		}
	}
	return Task{Title: ""}
}

var timeLayout = "02.01.2006"

func TimeParser(strDate string) (time.Time, error) {
	date, err := time.Parse(timeLayout, strDate)
	if err != nil {
		fmt.Println("Time parse error!")
		return date, err
	}
	return date, nil
}

//------------ CHANGES -----------------

func (sl *SlTasks) Add(title, desc, dLine string) string {
	if _, err := TimeParser(dLine); err != nil {
		return "deadline error"
	}
	for _, val := range *sl {
		if val.Title == title {
			return "title error"
		}
	}
	t := Task{title, desc, dLine, false, strings.Repeat("_", 11)} // strings.Repeat("_", 16
	*sl = append(*sl, t)
	return ""
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
				deadline, err := TimeParser(dLine)
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

func Show(val Task) string {
	var mark, titleSuffix string
	if val.Complete == true {
		mark = "\U0001F7E2"
	} else {
		mark = "🔴"

		dLine, _ := TimeParser(val.Deadline)
		// Calculate the number of hours until the deadline
		dur := time.Until(dLine).Hours()
		if dur > 0 && dur <= 24 {
			titleSuffix = "🔥"
		} else if dur <= 0 {
			titleSuffix = "🚫"
		}
	}

	str := fmt.Sprintf("🏷 %s %s\n📝 %s\n⏰ %s\n%s %s\n\n", val.Title, titleSuffix, val.Description, val.Deadline,
		mark, val.CompleteDate)
	return str
}

func (sl *SlTasks) ShowAll() string {
	var str string
	for _, val := range *sl {
		str += Show(val)
	}
	return str
}

func (sl *SlTasks) ShowUncompleted() string {
	st := *sl
	now := time.Now()
	for _, val := range st {
		dLine, _ := TimeParser(val.Deadline)

		if val.Complete == true || dLine.Before(now) {
			st.Delete(val.Title)
		}
	}
	if len(st) == 0 {
		return ""
	}

	for i := 0; i < len(st)-1; i++ {
		for j := 0; j < len(st)-i-1; j++ {
			curDLine, _ := TimeParser(st[j].Deadline)
			nextDLine, _ := TimeParser(st[j+1].Deadline)

			if curDLine.After(nextDLine) {
				st[j], st[j+1] = st[j+1], st[j]
			}
		}
	}

	var str string
	for _, val := range st {
		str += Show(val)
	}
	return str
}

func (sl *SlTasks) ShowOverdue() string {
	now := time.Now()
	var str string
	for _, val := range *sl {
		dLine, _ := TimeParser(val.Deadline)
		if val.Complete == false && dLine.Before(now) {
			str += Show(val)
		}
	}
	return str
}

//------------ DB CONTROL -----------------

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
