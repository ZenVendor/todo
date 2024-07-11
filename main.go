package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"slices"
	"time"
)

const VERSION = "0.9.1"

//go:embed help.txt
var helpString string

func PrintVersion() {
	fmt.Printf("TODO CLI\tversion: %s\n", VERSION)
}
func PrintHelp() {
	fmt.Println(helpString)
}

func main() {
	var conf Config

	cmd, sw, vals, err := ParseArgs(CMD_LIST, SW_OPEN)

	if err != nil || cmd == CMD_HELP {
		PrintVersion()
		PrintHelp()
		log.Fatal(err)
	}
	if cmd == CMD_VERSION {
		PrintVersion()
		return
	}
	if cmd == CMD_PREPARE || cmd == CMD_RESET {
		local := false
		if sw == SW_LOCAL {
			local = true
		}
		reset := false
		if cmd == CMD_RESET {
			reset = true
		}
		conf.Prepare(local, reset)
		return
	}

	conf.ReadConfig()
	db, err := conf.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if cmd == CMD_COUNT {
		count, err := CountTask(db, sw)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%d", count)
		return
	}

	if slices.Contains([]int{CMD_COMPLETE, CMD_DELETE, CMD_REOPEN, CMD_UPDATE}, cmd) {
		var t Task
		addGroup := false

		t.Id = vals.ReadValue(cmd, SW_DEFAULT, conf.DateFormat).(int)
		if err = t.Select(db); err != nil {
			log.Fatal(err)
		}
		switch cmd {
		case CMD_COMPLETE:
			t.Done = 1
			t.Completed = NullNow()
		case CMD_REOPEN:
			t.Done = 0
			t.Completed.Valid = false
		case CMD_UPDATE:
			if vals.ValueIsSet(SW_DESCRIPTION) {
				t.Description = vals.ReadValue(cmd, SW_DESCRIPTION, conf.DateFormat).(string)
			}
			if vals.ValueIsSet(SW_PRIORITY) {
				t.Priority = vals.ReadValue(cmd, SW_PRIORITY, conf.DateFormat).(int)
			}
			if vals.ValueIsSet(SW_DUE) {
				t.Due = vals.ReadValue(cmd, SW_DUE, conf.DateFormat).(sql.NullTime)
			}
			if vals.ValueIsSet(SW_GROUP) {
				t.Group.Name = vals.ReadValue(cmd, SW_GROUP, conf.DateFormat).(string)
				if err = t.Group.Select(db, true); err == sql.ErrNoRows {
					addGroup = true
				}
			}
		}
		t.Updated = NullNow()

		if cmd == CMD_DELETE {
			err = t.Delete(db)
		} else {
			if addGroup {
				t.Group.Created = NullNow()
				t.Group.Updated = NullNow()
				if err = t.Group.Insert(db); err != nil {
					log.Fatal(err)
				}
			}
			err = t.Update(db)
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Task %d has been %sd.\n", t.Id, MapCommandDescription()[cmd])
	}

	if cmd == CMD_ADD {
		var t Task
		//t.Id
		t.Description = vals.ReadValue(cmd, SW_DEFAULT, conf.DateFormat).(string)
		t.Priority = 1
		if vals.ValueIsSet(SW_PRIORITY) {
			t.Priority = vals.ReadValue(cmd, SW_PRIORITY, conf.DateFormat).(int)
		}
		t.Done = 0
		//t.Due
		if vals.ValueIsSet(SW_DUE) {
			t.Due = vals.ReadValue(cmd, SW_DUE, conf.DateFormat).(sql.NullTime)
		}
		//t.Completed
		t.Created = NullNow()
		t.Updated = NullNow()
		t.Group.Name = "Default"
		if vals.ValueIsSet(SW_GROUP) {
			t.Group.Name = vals.ReadValue(cmd, SW_GROUP, conf.DateFormat).(string)
		}
		if err = t.Group.Select(db, true); err == sql.ErrNoRows {
			t.Group.Created = NullNow()
			t.Group.Updated = NullNow()
			if err = t.Group.Insert(db); err != nil {
				log.Fatal(err)
			}
		}
		if err = t.Insert(db); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Added task %d: %s\n", t.Id, t.Description)
	}

	if cmd == CMD_LIST {
		count, err := CountTask(db, sw)
		if err != nil {
			log.Fatal(err)
		}
		tl, err := List(db, sw)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s tasks: %d\n", MapArgumentDescription()[sw], count)

		for _, t := range tl {
			tStatus := "Open"
			if t.Done == 0 && t.Due.Valid && t.Due.Time.Before(time.Now()) {
				tStatus = "Overdue"
			}
			if t.Done == 1 {
				tStatus = "Closed"
			}
			if t.Done == 0 {
				if !t.Due.Valid {
					fmt.Printf("\t%s :: %s %d: %s\n", t.Group.Name, tStatus, t.Id, t.Description)
				} else {
					fmt.Printf("\t%s :: %s %d: %s, due date: %s\n", t.Group.Name, tStatus, t.Id, t.Description, t.Due.Time.Format(conf.DateFormat))
				}
			} else {
				fmt.Printf("\t%s :: %s %d: %s, completed: %s\n", t.Group.Name, tStatus, t.Id, t.Description, t.Completed.Time.Format(conf.DateFormat))
			}
		}
	}
}
