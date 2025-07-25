package main

import (
	"database/sql"
	"fmt"
	"github.com/nonerkao/color-aware-tabwriter"
	"os"
	"strings"
)

func (t *Task) SetOptional(p *Parser, db *sql.DB) (err error) {

	if p.ArgIsPresent(A_COMMENT) {
		t.ClosingComment, err = GetDescriptionFromEditor(t.ClosingComment)
		if err != nil {
			return err
		}
	}
	if p.ArgIsPresent(A_DESCRIPTION) {
		t.Description, err = GetDescriptionFromEditor(t.Description)
		if err != nil {
			return err
		}
	}
	if value, ok := p.Kwargs[K_DATEDUE]; ok {
		t.DateDue = value.(sql.NullTime)
	}
	if value, ok := p.Kwargs[K_PRIORITY]; ok {
		t.Priority = value.(int)
	}
	if value, ok := p.Kwargs[K_PROJECT]; ok {
		t.Project.Name = value.(string)
	}
	if value, ok := p.Kwargs[K_SUMMARY]; ok {
		t.Summary = value.(string)
	}

	// If parent is set, get parent task
	if value, ok := p.Kwargs[K_PARENT]; ok {
		// Unset parent
		if p.Kwargs[K_PARENT] == 0 {
			t.Parent = nil
			return err
		}
		t.Parent = &Task{}
		t.Parent.Id = value.(int)
	}
	if t.Parent != nil && t.Parent.Id > 0 {
		if err = t.Parent.GetTask(db); err != nil {
			return err
		}
		// Parent project overrides provided value
		if t.Project.Id != t.Parent.Project.Id {
			t.Project.Id = t.Parent.Project.Id
			t.Project.Name = t.Parent.Project.Name
		}
		// If provided due date is later than parent's, use parent's
		if t.Parent.DateDue.Valid {
			if !t.DateDue.Valid {
				t.DateDue = t.Parent.DateDue
			}
			if t.DateDue.Time.After(t.Parent.DateDue.Time) {
				t.DateDue = t.Parent.DateDue
			}
		}
		// If priority is lower than parent's, use parent's
		if t.Priority > t.Parent.Priority {
			t.Priority = t.Parent.Priority
		}
	}
	return err
}

func (p *Parser) Add(db *sql.DB) (err error) {
	var t Task

	// Set required summary and defaults
	t.Summary = p.Kwargs[K_SUMMARY].(string)
	t.Priority = PRIORITY_MED
	t.Status = STATUS_NEW
	t.Project = Project{DEFAULT_GROUP, ""}

	if err = t.SetOptional(p, db); err != nil {
		return err
	}
	if err = t.Project.GetProject(db); err != nil {
		if err = t.Project.Add(db); err != nil {
			return err
		}
	}
	if err = t.Add(db); err != nil {
		return err
	}
	fmt.Printf("Added task: %d - %s\n", t.Id, t.Summary)
	return nil
}

func (p *Parser) Complete(db *sql.DB) (err error) {
	var t Task
	t.Id = p.Kwargs[K_ID].(int)
	if err = t.GetTask(db); err != nil {
		return err
	}

	t.Status = STATUS_COMPLETED
	t.DateCompleted = NullNow()

	if err = t.SetOptional(p, db); err != nil {
		return err
	}
	if err = t.Update(db); err != nil {
		return err
	}

	var bs strings.Builder
	fmt.Fprintf(&bs, "Completed task: %d - %s", t.Id, t.Summary)

	// Close subtasks
	if err = t.GetChildren(db); err != nil {
		return err
	}
	fmt.Fprintf(&bs, "Children: %d\n", len(t.Children))
	if len(t.Children) > 0 {
		plural := ""
		if len(t.Children) > 1 {
			plural = "s"
		}
		fmt.Fprintf(&bs, "\nand %d subtask%s:", len(t.Children), plural)

		for _, c := range t.Children {
			(*c).Status = STATUS_COMPLETED
			(*c).DateCompleted = NullNow()
			(*c).ClosingComment = "Closed by main task."

			if err = (*c).Update(db); err != nil {
				return err
			}
			fmt.Fprintf(&bs, "\n\t%d - %s", (*c).Id, (*c).Summary)
		}
	}
	fmt.Fprint(&bs, "\n")
	fmt.Println(bs.String())
	return err
}

func (p *Parser) Count(db *sql.DB) (err error) {
	var c Counts
	if err := c.GetCounts(db); err != nil {
		return err
	}
	var sw int
	if len(p.Args) == 0 {
		sw = A_ALL
	} else {
		sw = p.Args[0]
	}

	// parser ensures only one argument is present
	switch sw {
	default:
		w := tabwriter.NewWriter(os.Stdout, 4, 0, 2, ' ', 0)
		fmt.Fprintf(w, "%sStatus:\tCount%s\n", C_BOLD, C_RESET)
		fmt.Fprintf(w, "All:\t%d\n", c.All)
		fmt.Fprintf(w, "Open:\t%d\n", c.Open)
		fmt.Fprintf(w, "Overdue:\t%d\n", c.Overdue)
		fmt.Fprintf(w, "New:\t%d\n", c.New)
		fmt.Fprintf(w, "In progress:\t%d\n", c.InProgress)
		fmt.Fprintf(w, "On hold:\t%d\n", c.OnHold)
		fmt.Fprintf(w, "Completed:\t%d\n", c.Completed)
		w.Flush()
	// Others are for prompt
	case A_COMPLETED:
		fmt.Printf("%d", c.Completed)
	case A_INPROGRESS:
		fmt.Printf("%d", c.InProgress)
	case A_NEW:
		fmt.Printf("%d", c.New)
	case A_ONHOLD:
		fmt.Printf("%d", c.OnHold)
	case A_OPEN:
		fmt.Printf("%d", c.Open)
	case A_OVERDUE:
		fmt.Printf("%d", c.Overdue)

	}
	return err
}

func (p *Parser) Delete(db *sql.DB) (err error) {
	var t Task

	t.Id = p.Kwargs[K_ID].(int)
	if err = t.GetTask(db); err != nil {
		return err
	}
	if err = t.Delete(db); err != nil {
		return err
	}
	var bs strings.Builder
	fmt.Fprintf(&bs, "Deleted task: %d - %s", t.Id, t.Summary)

	// Unlink or delete subtasks
	if err = t.GetChildren(db); err != nil {
		return err
	}
	if len(t.Children) > 0 {
		plural := ""
		if len(t.Children) > 1 {
			plural = "s"
		}
		if p.ArgIsPresent(A_ALL) {
			fmt.Fprintf(&bs, "\nand %d subtask%s:", len(t.Children), plural)
			for _, c := range t.Children {
				if err = (*c).Delete(db); err != nil {
					return err
				}
				fmt.Fprintf(&bs, "\n\t%d - %s", (*c).Id, (*c).Summary)
			}
		} else {
			fmt.Fprintf(&bs, "\nand unlinked %d subtask%s:", len(t.Children), plural)
			for _, c := range t.Children {
				(*c).Parent.Id = 0
				if err = (*c).Update(db); err != nil {
					return err
				}
				fmt.Fprintf(&bs, "\n\t%d - %s", (*c).Id, (*c).Summary)
			}
		}
	}
	fmt.Fprint(&bs, "\n")
	fmt.Println(bs.String())
	return err
}

func (p *Parser) Help(db *sql.DB) (err error) {
	fmt.Println(embedHelp)
	return nil
}

func (p *Parser) Hold(db *sql.DB) (err error) {
	var t Task
	t.Id = p.Kwargs[K_ID].(int)
	if err = t.GetTask(db); err != nil {
		return err
	}
	t.Status = STATUS_HOLD
	if err = t.SetOptional(p, db); err != nil {
		return err
	}
	if err = t.Update(db); err != nil {
		return err
	}
	var bs strings.Builder
	fmt.Fprintf(&bs, "Task put on hold: %d - %s", t.Id, t.Summary)

	// Hold subtasks
	if err = t.GetChildren(db); err != nil {
		return err
	}
	if len(t.Children) > 0 {
		plural := ""
		if len(t.Children) > 1 {
			plural = "s"
		}
		fmt.Fprintf(&bs, "\nincluding %d subtask%s:", len(t.Children), plural)

		for _, c := range t.Children {
			(*c).Status = STATUS_HOLD

			if err = (*c).Update(db); err != nil {
				return err
			}
			fmt.Fprintf(&bs, "\n\t%d - %s", (*c).Id, (*c).Summary)
		}
	}
	fmt.Fprint(&bs, "\n")
	fmt.Println(bs.String())
	return err
}

func (p *Parser) List(db *sql.DB) (err error) {
	tl, err := ListTasks(db, p.Args[0])
	if err != nil {
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 4, 0, 2, ' ', 0)
	switch (*p).Args[0] {
	case A_DELETED:
		fmt.Fprintf(w, "%sID\tProject\tStatus\tPriority\tDue\tCompleted\tSummary%s\n", C_BOLD, C_RESET)
	case A_COMPLETED:
		fmt.Fprintf(w, "%sID\tProject\tStatus\tCompleted\tSummary%s\n", C_BOLD, C_RESET)
	default:
		fmt.Fprintf(w, "%sID\tProject\tStatus\tPriority\tDue\tSummary%s\n", C_BOLD, C_RESET)
	}
	for _, t := range tl {
		tDue := HumanDue(t.DateDue.Time, "2006-01-02")
		if !t.DateDue.Valid {
			tDue = ""
		}
		tCompleted := t.DateCompleted.Time.Format("2006-01-02")
		if !t.DateCompleted.Valid {
			tCompleted = ""
		}
		switch (*p).Args[0] {
		case A_DELETED:
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\t%s\n", t.Id, t.Project.Name, statusMap[t.Status], DisplayPriority(t.Priority), tDue, tCompleted, t.Summary)
		case A_COMPLETED:
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n", t.Id, t.Project.Name, statusMap[t.Status], tCompleted, t.Summary)
		default:
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\n", t.Id, t.Project.Name, statusMap[t.Status], DisplayPriority(t.Priority), tDue, t.Summary)
		}
	}
	w.Flush()
	return err
}

func (p *Parser) Reopen(db *sql.DB) (err error) {
	var t Task
	t.Id = p.Kwargs[K_ID].(int)
	if err = t.GetTask(db); err != nil {
		return err
	}
	t.Status = STATUS_INPROG
	t.DateCompleted.Valid = false

	if err = t.SetOptional(p, db); err != nil {
		return err
	}
	if err = t.Update(db); err != nil {
		return err
	}

	var bs strings.Builder
	fmt.Fprintf(&bs, "Task resumed: %d - %s", t.Id, t.Summary)

	// Hold subtasks
	if err = t.GetChildren(db); err != nil {
		return err
	}
	if len(t.Children) > 0 {
		plural := ""
		if len(t.Children) > 1 {
			plural = "s"
		}
		fmt.Fprintf(&bs, "\nincluding %d subtask%s:", len(t.Children), plural)

		for _, c := range t.Children {
			(*c).Status = STATUS_INPROG
			(*c).DateCompleted.Valid = false

			if err = (*c).Update(db); err != nil {
				return err
			}
			fmt.Fprintf(&bs, "\n\t%d - %s", (*c).Id, (*c).Summary)
		}
	}
	fmt.Fprintf(&bs, "\n")

	fmt.Println(bs.String())
	return err
}

func (p *Parser) Show(db *sql.DB) (err error) {
	var t Task
	t.Id = p.Kwargs[K_ID].(int)
	if err = t.GetTask(db); err != nil {
		return err
	}
	if err = t.GetChildren(db); err != nil {
		return err
	}
	del := ""
	if t.SysStatus == SYS_DELETED {
		del = fmt.Sprintf("\t%sDELETED%s", C_RED, C_RESET)
	}
	fmt.Printf("\n%sTASK %d%s%s\n", C_BOLD, t.Id, del, C_RESET)

	w := tabwriter.NewWriter(os.Stdout, 4, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Summary:\t%s\n", t.Summary)
	fmt.Fprintf(w, "Project\t%s\n", t.Project.Name)
	fmt.Fprintf(w, "Priority:\t%s\n", DisplayPriority(t.Priority))
	fmt.Fprintf(w, "Status:\t%s\n", statusMap[t.Status])
	if t.DateDue.Valid {
		fmt.Fprintf(w, "\tDue Date:\t%s\n", t.DateDue.Time.Format("2006-01-02"))
	}
	if t.Description != "" {
		fmt.Fprintf(w, "\n%sDescription:\t%s\n", C_BOLD, C_RESET)
		fmt.Fprint(w, "---------\n")
		fmt.Fprintf(w, "%v", t.Description)
		fmt.Fprint(w, "---------\n\n")
	}
	if t.Status == STATUS_COMPLETED {
		fmt.Fprintf(w, "Completed:\t%s\n", t.DateCompleted.Time.Format("2006-01-02"))
		fmt.Fprintf(w, "Comment:\t%s\n", t.ClosingComment)
	}
	if t.Parent != nil {
		fmt.Fprintf(w, "%sParent:%s\t%d - %s\n", C_BOLD, C_RESET, t.Parent.Id, t.Parent.Summary)
	}
	if len(t.Children) > 0 {
		fmt.Fprintf(w, "Linked tasks:\n")
		for _, c := range t.Children {
			fmt.Fprintf(w, "\t%d: %s\n", (*c).Id, (*c).Summary)
		}
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "Created:\t%s\n", t.DateCreated.Time.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(w, "Updated:\t%s\n", t.DateUpdated.Time.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(w, "\n")

	err = w.Flush()
	return err
}

func (p *Parser) Update(db *sql.DB) (err error) {
	var t Task
	t.Id = p.Kwargs[K_ID].(int)
	if err = t.GetTask(db); err != nil {
		return err
	}
	if err = t.SetOptional(p, db); err != nil {
		return err
	}
	if err = t.Project.GetProject(db); err != nil {
		if err = t.Project.Add(db); err != nil {
			return err
		}
	}
	if err = t.Update(db); err != nil {
		return err
	}
	fmt.Printf("Task updated: %d - %s\n", t.Id, t.Summary)
	return err
}

func (p *Parser) Version(db *sql.DB) (err error) {
	fmt.Printf("TODO CLI :: version: %s\n", VERSION_APP)
	return nil
}
