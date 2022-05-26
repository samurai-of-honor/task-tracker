package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type Tasks struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Deadline     string `json:"deadline"`
	Complete     bool   `json:"done"`
	CompleteDate string `json:"completeDate"`
}

func New() *[]Tasks {
	sl := make([]Tasks, 0)
	return &sl
}

func Add(sl *[]Tasks, title, desc, dLine string, done bool, compDate string) {
	t := Tasks{title, desc, dLine, done, compDate}
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

func ShowAll(sl *[]Tasks) {
	st := *sl
	var mark string
	for i, val := range st {
		if val.Complete == true {
			mark = "X"
		} else {
			mark = " "
		}
		fmt.Printf("%d. [%s]%s %s    %s\n%s\n", i+1, mark, val.CompleteDate, val.Title, val.Deadline, val.Description)
	}
}

func ShowUncompleted(sl *[]Tasks) {
	st := *sl
	now := time.Now()
	for _, val := range st {
		if val.Complete == true {
			Delete(&st, val.Title)
		} else {
			dLine, err := time.Parse("02-01-2006 15:04", val.Deadline)
			if err != nil {
				panic(err)
			}
			if dLine.Before(now) {
				Delete(&st, val.Title)
			}
		}
	}

	current, err := time.Parse("02-01-2006 15:04", st[0].Deadline)
	if err != nil {
		panic(err)
	}
	for i, val := range st {
		dLine, err1 := time.Parse("02-01-2006 15:04", val.Deadline)
		if err1 != nil {
			panic(err1)
		}

		if dLine.Before(current) {
			st[i], st[0] = st[0], st[i]
			current = dLine
		}
	}

	for _, val := range st {
		fmt.Printf(" %s    %s\n%s\n", val.Title, val.Deadline, val.Description)
	}
}

func Load(sl *[]Tasks) {
	file, openErr := os.Open("db.json")
	if openErr != nil {
		Save(sl)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
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
