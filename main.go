package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"github.com/nonerkao/color-aware-tabwriter"
	"log"
	"os"
	"slices"
	"time"
)

const VERSION = "1.0.0"

const C_WHITE = "\033[37m"
const C_GREY = "\033[38;5;7m"
const C_RED = "\033[38;5;9m"
const C_ORANGE = "\033[38;5;214m"
const C_YELLOW = "\033[38;5;226m"
const C_GREEN = "\033[38;5;40m"
const C_BOLD = "\033[1m"
const C_RESET = "\033[0m"

//go:embed help.txt
var helpString string

func PrintVersion() {
	fmt.Printf("TODO CLI\tversion: %s\n", VERSION)
}
func PrintHelp() {
	fmt.Println(helpString)
}

func Color(text, color string, bold bool) string {
	if bold {
		color = fmt.Sprintf("%s%s", color, C_BOLD)
	}
	return fmt.Sprintf("%s%s%s", color, text, C_RESET)
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
		fmt.Printf("%s%s tasks: %d%s\n", C_BOLD, MapArgumentDescription()[sw], count, C_RESET)

		w := tabwriter.NewWriter(os.Stdout, 4, 0, 2, ' ', 0)
		fmt.Fprintf(w, "%s\tID\tGroup\tStatus\tDue\tDescription%s\n", C_BOLD, C_RESET)

		for _, t := range tl {
			lColor := C_WHITE
			tStatus := "Open"

			if t.Done == 1 {
				tStatus = "Closed"
				lColor = C_GREEN
			}
			if t.Done == 0 && t.Due.Valid {
				if t.Due.Time.Sub(time.Now()).Hours() > 120 && t.Due.Time.Sub(time.Now()).Hours() <= 240 {
					lColor = C_YELLOW
				}
				if t.Due.Time.Sub(time.Now()).Hours() <= 120 {
					lColor = C_ORANGE
				}
				if t.Due.Time.Before(time.Now()) {
					tStatus = "Overdue"
					lColor = C_RED
				}
			}
			tDue := HumanDue(t.Due.Time, conf.DateFormat)
			if !t.Due.Valid {
				tDue = ""
			}

			fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%s\t%s%s\n", lColor, t.Id, t.Group.Name, tStatus, tDue, t.Description, C_RESET)
		}
		w.Flush()
	}
	if cmd == CMD_SHOW {
		var t Task

		t.Id = vals.ReadValue(cmd, SW_DEFAULT, conf.DateFormat).(int)
		if err = t.Select(db); err != nil {
			log.Fatal(err)
		}

		tStatus := "Open"
		if t.Done == 0 && t.Due.Valid && t.Due.Time.Before(time.Now()) {
			tStatus = Color("Overdue", C_RED, true)
		}
		if t.Done == 1 {
			tStatus = Color("Closed", C_GREEN, false)
		}

		fmt.Printf("Task ID: %d\n", t.Id)
		fmt.Printf("Status: %s\n", tStatus)
		fmt.Printf("Group: %s\n", t.Group.Name)

		fmt.Printf("Priority: %d\n", t.Priority)
		if t.Due.Valid {
			fmt.Printf("Due: %s\n", t.Due.Time.Format(conf.DateFormat))
		}
		fmt.Printf("\n%s\n\n", Color(t.Description, C_WHITE, true))
		meta := fmt.Sprintf("Created: %s\nUpdated: %s\n", t.Created.Time.Format(conf.DateFormat), t.Updated.Time.Format(conf.DateFormat))
		fmt.Print(Color(meta, C_GREY, false))
	}
}
