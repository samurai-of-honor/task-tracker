package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
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
		fmt.Printf("%d. [%s] %s    %s\n%s\n", i+1, mark, val.Title, val.Deadline, val.Description)
	}
}

func Save(sl *[]Tasks) {
	data, convErr := json.MarshalIndent(sl, "", "  ")
	if convErr != nil {
		panic(convErr)
	}

	file, openErr := os.Create("db.json")
	if openErr != nil {
		panic(openErr)
	}
	defer file.Close()

	w := io.Writer(file)
	if _, writeErr := w.Write(data); writeErr != nil {
		panic(writeErr)
	}
}
