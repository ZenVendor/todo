package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/nonerkao/color-aware-tabwriter"
)

func main() {
	var conf Config
	var db *sql.DB
	args := os.Args[1:]

	parser, err := NewParser(V_LIST, []int{A_OPEN}, map[int]interface{}{})
	if err != nil {
		log.Fatal(err)
	}
	if err = parser.Parse(args); err != nil {
		log.Fatal(err)
	}

	db, err = conf.Prepare()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	msg, err := parser.Verb.Call(&parser, db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("TODO:", msg)

	// ===== Left for reference

	//	if parser.Verb.Verb == V_LIST {
	//		tl, err := ListTasks(db)
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//
	//		w := tabwriter.NewWriter(os.Stdout, 4, 0, 2, ' ', 0)
	//		fmt.Fprintf(w, "%s\tID\tGroup\tStatus\tDue\tDescription%s\n", C_BOLD, C_RESET)
	//
	//		for _, t := range tl {
	//			lColor := C_WHITE
	//			tStatus := "Open"
	//
	//			tDue := HumanDue(t.DateDue.Time, conf.DateFormat)
	//			if !t.DateDue.Valid {
	//				tDue = ""
	//			}
	//
	//			fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%s\t%s%s\n", lColor, t.Id, t.Group.Name, tStatus, tDue, t.Description, C_RESET)
	//		}
	//		w.Flush()
	//	}
}
