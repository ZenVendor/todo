package main

import (
	"database/sql"
	"fmt"
	"github.com/nonerkao/color-aware-tabwriter"
	"log"
	"os"
	"slices"
	"time"
)

func main() {
	var conf Config
	args := os.Args[1:]
	parser, err := NewParser(V_LIST, []int{A_ONGOING}, map[int]interface{}{})
	if err != nil {
		log.Fatal(err)
	}
	if err = parser.Parse(args); err != nil {
		log.Fatal(err)
	}

	if parser.Verb.Verb == V_HELP {
		PrintVersion()
		PrintHelp()
		return
	}
	if parser.Verb.Verb == V_VERSION {
		PrintVersion()
		return
	}
	if parser.Verb.Verb == V_CONFIGURE {
		conf.Prepare(
			parser.ArgIsPresent(A_LOCAL),
			parser.ArgIsPresent(A_RESET),
		)
		return
	}

	conf.ReadConfig()
	db, err := conf.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if parser.Verb.Verb == V_COUNT {
		count, err := CountTask(db, parser.GetArg(0))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%d", count)
		return
	}

	if slices.Contains([]int{V_COMPLETE, V_DELETE, V_REOPEN, V_UPDATE}, parser.Verb.Verb) {
		var t Task
		addGroup := false

		t.Id = parser.Kwargs[K_ID].(int)
		if err = t.Select(db); err != nil {
			log.Fatal(err)
		}
		switch parser.Verb.Verb {
		case V_COMPLETE:
			t.Done = 1
			t.DateCompleted = NullNow()
		case V_REOPEN:
			t.Done = 0
			t.DateCompleted.Valid = false
		case V_UPDATE:
			if short, ok := parser.Kwargs[K_SHORT]; ok {
				t.Description = short.(string)
			}
			if priority, ok := parser.Kwargs[K_PRIORITY]; ok {
				t.Priority = priority.(int)
			}
			if duedate, ok := parser.Kwargs[K_DUEDATE]; ok {
				t.DateDue = duedate.(sql.NullTime)
			}
			if group, ok := parser.Kwargs[K_GROUP]; ok {
				t.Group.Name = group.(string)
				if err = t.Group.Select(db, true); err == sql.ErrNoRows {
					addGroup = true
				}
			}
		}
		t.SysDateUpdated = NullNow()

		if parser.Verb.Verb == V_DELETE {
			err = t.Delete(db)
		} else {
			if addGroup {
				t.Group.SysDateCreated = NullNow()
				t.Group.SysDateUpdated = NullNow()
				if err = t.Group.Insert(db); err != nil {
					log.Fatal(err)
				}
			}
			err = t.Update(db)
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Task %d has been added.\n", t.Id)
	}

	if parser.Verb.Verb == V_ADD {
		var t Task
		//t.Id
		t.Description = parser.Kwargs[K_SHORT].(string)
		t.Priority = 1
		if priority, ok := parser.Kwargs[K_PRIORITY]; ok {
			t.Priority = priority.(int)
		}
		t.Done = 0
		//t.Due
		if duedate, ok := parser.Kwargs[K_DUEDATE]; ok {
			t.DateDue = duedate.(sql.NullTime)
		}
		//t.Completed
		t.SysDateCreated = NullNow()
		t.SysDateUpdated = NullNow()
		t.Group.Name = "Default"
		if group, ok := parser.Kwargs[K_GROUP]; ok {
			t.Group.Name = group.(string)
		}
		if err = t.Group.Select(db, true); err == sql.ErrNoRows {
			t.Group.SysDateCreated = NullNow()
			t.Group.SysDateUpdated = NullNow()
			if err = t.Group.Insert(db); err != nil {
				log.Fatal(err)
			}
		}
		if err = t.Insert(db); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Added task %d: %s\n", t.Id, t.Description)
	}

	if parser.Verb.Verb == V_LIST {
		count, err := CountTask(db, parser.GetArg(0))
		if err != nil {
			log.Fatal(err)
		}
		tl, err := List(db, parser.GetArg(0))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%stasks: %d%s\n", C_BOLD, count, C_RESET)

		w := tabwriter.NewWriter(os.Stdout, 4, 0, 2, ' ', 0)
		fmt.Fprintf(w, "%s\tID\tGroup\tStatus\tDue\tDescription%s\n", C_BOLD, C_RESET)

		for _, t := range tl {
			lColor := C_WHITE
			tStatus := "Open"

			if t.Done == 1 {
				tStatus = "Closed"
				lColor = C_GREEN
			}
			if t.Done == 0 && t.DateDue.Valid {
				if t.DateDue.Time.Sub(time.Now()).Hours() > 120 && t.DateDue.Time.Sub(time.Now()).Hours() <= 240 {
					lColor = C_YELLOW
				}
				if t.DateDue.Time.Sub(time.Now()).Hours() <= 120 {
					lColor = C_ORANGE
				}
				if t.DateDue.Time.Before(time.Now()) {
					tStatus = "Overdue"
					lColor = C_RED
				}
			}
			tDue := HumanDue(t.DateDue.Time, conf.DateFormat)
			if !t.DateDue.Valid {
				tDue = ""
			}

			fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%s\t%s%s\n", lColor, t.Id, t.Group.Name, tStatus, tDue, t.Description, C_RESET)
		}
		w.Flush()
	}
	if parser.Verb.Verb == V_SHOW {
		var t Task

		t.Id = parser.Kwargs[K_ID].(int)
		if err = t.Select(db); err != nil {
			log.Fatal(err)
		}

		tStatus := "Open"
		if t.Done == 0 && t.DateDue.Valid && t.DateDue.Time.Before(time.Now()) {
			tStatus = Color("Overdue", C_RED, true)
		}
		if t.Done == 1 {
			tStatus = Color("Closed", C_GREEN, false)
		}

		fmt.Printf("Task ID: %d\n", t.Id)
		fmt.Printf("Status: %s\n", tStatus)
		fmt.Printf("Group: %s\n", t.Group.Name)

		fmt.Printf("Priority: %d\n", t.Priority)
		if t.DateDue.Valid {
			fmt.Printf("Due: %s\n", t.DateDue.Time.Format(conf.DateFormat))
		}
		fmt.Printf("\n%s\n\n", Color(t.Description, C_WHITE, true))
		meta := fmt.Sprintf("Created: %s\nUpdated: %s\n", t.SysDateCreated.Time.Format(conf.DateFormat), t.SysDateUpdated.Time.Format(conf.DateFormat))
		fmt.Print(Color(meta, C_GREY, false))
	}
}
