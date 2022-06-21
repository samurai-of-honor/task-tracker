package main

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"strconv"
	"strings"
	tm "task-manager"
	"task-manager/tg/localization"
	"time"
)

func NotificationOn(ok bool) {
	if ok == true {
		tick := time.NewTicker(time.Hour * 12)
		defer tick.Stop()

		tasks := tm.Create()

		for {
			select {
			case <-tick.C:
				filesFromDir, err := ioutil.ReadDir("./tg/taskBases/")
				if err != nil {
					panic(err)
				}

				for _, file := range filesFromDir {
					tasks.Load("./tg/taskBases/" + file.Name())
					var msgStr string

					// Read all tasks
					for _, v := range *tasks {
						if v.Complete != true {
							dLine, _ := tm.TimeParser(v.Deadline)
							// Calculate the number of hours until the deadline
							dur := time.Until(dLine).Hours()
							if dur > 0 && dur <= 24 || dur > 48 && dur <= 72 {
								msgStr += tm.Show(v)
							}
						}
					}

					// If completed tasks are found
					if msgStr != "" {
						userID, err := strconv.Atoi(strings.TrimSuffix(file.Name(), ".json"))
						if err != nil {
							continue
						}

						dontForgetNotifMsg1 := tg.NewMessage(int64(userID), localization.DontForgetMsg)
						Send(dontForgetNotifMsg1)
						dontForgetNotifMsg2 := tg.NewMessage(int64(userID), msgStr)
						Send(dontForgetNotifMsg2)
					}
				}
			}
		}
	}
}
