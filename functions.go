package main

import (
	"database/sql"
	"fmt"
	"github.com/nonerkao/color-aware-tabwriter"
	"os"
	"strings"
)

func (p *Parser) Add(db *sql.DB, conf *Config) (msg string, err error) {
	// Required: K_SUMMARY
	// Optional: K_DESCRIPTION, K_DATEDUE, K_GROUP, K_PARENT, K_PRIORITY

	var t Task

	// Set required summary and defaults
	t.Summary = p.Kwargs[K_SUMMARY].(string)
	t.Priority = PRIORITY_MED
	t.Status = STATUS_NEW
	t.Group = Group{DEFAULT_GROUP, ""}

	// Optional values
	if value, ok := p.Kwargs[K_DESCRIPTION]; ok {
		t.Description = value.(string)
	}
	if value, ok := p.Kwargs[K_DATEDUE]; ok {
		t.DateDue = value.(sql.NullTime)
	}
	if value, ok := p.Kwargs[K_PROJECT]; ok {
		t.Group.Name = value.(string)
	}
	if value, ok := p.Kwargs[K_PRIORITY]; ok {
		t.Priority = value.(int)
	}

	// If parent is set, get parent task
	if value, ok := p.Kwargs[K_PARENT]; ok {
		t.Parent.Id = value.(int)
		if err = t.Parent.GetTask(db); err != nil {
			return "", err
		}
		// If provided due date is later than parent's, use parent's
		if t.Parent.DateDue.Valid && (t.DateDue.Time.After(t.Parent.DateDue.Time) ||
			!t.DateDue.Valid) {
			t.DateDue = t.Parent.DateDue
		}

		// Parent group overrides provided value
		t.Group.Id = t.Parent.Group.Id
		t.Group.Name = t.Parent.Group.Name

		// If priority is lower than parent's, use parent's
		if t.Priority > t.Parent.Priority {
			t.Priority = t.Parent.Priority
		}
	}
	if err = t.Group.GetGroup(db); err != nil {
		if err = t.Group.Add(db); err != nil {
			return "", err
		}
	}
	if err = t.Add(db); err != nil {
		return "", err
	}

	msg = fmt.Sprintf("\t%sAdded task:% d - %s%s\n", C_GREEN, t.Id, t.Summary, C_RESET)
	return msg, nil
}

func (p *Parser) Complete(db *sql.DB, conf *Config) (msg string, err error) {
	// Required: K_ID
	// Optional: K_COMMENT
	var t Task
	t.Id = p.Kwargs[K_ID].(int)
	if err = t.GetTask(db); err != nil {
		return msg, err
	}

	t.Status = STATUS_COMPLETED
	t.DateCompleted = NullNow()

	if value, ok := p.Kwargs[K_COMMENT]; ok {
		t.ClosingComment = value.(string)
	}
	if err = t.Update(db); err != nil {
		return msg, err
	}

	var bs strings.Builder
	bs.WriteString(fmt.Sprintf("Completed task: %d - %s", t.Id, t.Summary))

	// Close subtasks
	if err = t.GetChildren(db); err != nil {
		return msg, err
	}
	fmt.Printf("Children: %d\n", len(t.Children))
	if len(t.Children) > 0 {
		plural := ""
		if len(t.Children) > 1 {
			plural = "s"
		}
		bs.WriteString(fmt.Sprintf("\nand %d subtask%s:", len(t.Children), plural))

		for _, c := range t.Children {
			(*c).Status = STATUS_COMPLETED
			(*c).DateCompleted = NullNow()
			(*c).ClosingComment = "Closed by main task."

			if err = (*c).Update(db); err != nil {
				return msg, err
			}
			bs.WriteString(fmt.Sprintf("\n\t%d - %s", (*c).Id, (*c).Summary))
		}
	}
	bs.WriteString("\n")

	return bs.String(), err
}

func (p *Parser) Count(db *sql.DB, conf *Config) (msg string, err error) {
	// Optional: A_ALL, A_COMPLETED, A_DUE, A_INPROGRESS, A_ONHOLD, A_OPEN, A_OVERDUE
	return msg, err
}

func (p *Parser) Delete(db *sql.DB, conf *Config) (msg string, err error) {
	// Required: K_ID
	// Optional: A_ALL
	var t Task
	t.Id = p.Kwargs[K_ID].(int)
	if err = t.GetTask(db); err != nil {
		return msg, err
	}
	if err = t.Delete(db); err != nil {
		return msg, err
	}
	var bs strings.Builder
	bs.WriteString(fmt.Sprintf("Deleted task: %d - %s", t.Id, t.Summary))

	// Unlink or delete subtasks
	if err = t.GetChildren(db); err != nil {
		return msg, err
	}
	if len(t.Children) > 0 {
		plural := ""
		if len(t.Children) > 1 {
			plural = "s"
		}
		if p.ArgIsPresent(A_ALL) {
			bs.WriteString(fmt.Sprintf("\nand %d subtask%s:", len(t.Children), plural))
			for _, c := range t.Children {
				if err = (*c).Delete(db); err != nil {
					return msg, err
				}
				bs.WriteString(fmt.Sprintf("\n\t%d - %s", (*c).Id, (*c).Summary))
			}
		} else {
			bs.WriteString(fmt.Sprintf("\nand unlinked %d subtask%s:", len(t.Children), plural))
			for _, c := range t.Children {
				(*c).Parent.Id = 0
				if err = (*c).Update(db); err != nil {
					return msg, err
				}
				bs.WriteString(fmt.Sprintf("\n\t%d - %s", (*c).Id, (*c).Summary))
			}
		}
	}
	bs.WriteString("\n")
	return bs.String(), err
}
func (p *Parser) Help(db *sql.DB, conf *Config) (msg string, err error) {
	return embedHelp, nil
}
func (p *Parser) Hold(db *sql.DB, conf *Config) (msg string, err error) {
	// Required: K_ID
	var t Task
	t.Id = p.Kwargs[K_ID].(int)
	if err = t.GetTask(db); err != nil {
		return "", err
	}
	t.Status = STATUS_HOLD
	if err = t.Update(db); err != nil {
		return "", err
	}

	var bs strings.Builder
	bs.WriteString(fmt.Sprintf("Task put on hold: %d - %s", t.Id, t.Summary))

	// Hold subtasks
	if err = t.GetChildren(db); err != nil {
		return msg, err
	}
	if len(t.Children) > 0 {
		plural := ""
		if len(t.Children) > 1 {
			plural = "s"
		}
		bs.WriteString(fmt.Sprintf("\nincluding %d subtask%s:", len(t.Children), plural))

		for _, c := range t.Children {
			(*c).Status = STATUS_HOLD

			if err = (*c).Update(db); err != nil {
				return msg, err
			}
			bs.WriteString(fmt.Sprintf("\n\t%d - %s", (*c).Id, (*c).Summary))
		}
	}
	bs.WriteString("\n")

	return bs.String(), err
}
func (p *Parser) List(db *sql.DB, conf *Config) (msg string, err error) {
	//	Optional: A_ALL, A_COMPLETED, A_DELETED, A_DUE, A_INPROGRESS, A_NEW, A_ONHOLD, A_OPEN, A_OVERDUE
	var tl TaskList

	switch (*p).Args[0] {
	case A_ALL:
		tl, err = ListTasksAll(db)
	case A_COMPLETED:
		tl, err = ListTasksCompleted(db)
	case A_DELETED:
		tl, err = ListTasksDeleted(db)
	case A_INPROGRESS:
		tl, err = ListTasksInProgress(db)
	case A_NEW:
		tl, err = ListTasksNew(db)
	case A_ONHOLD:
		tl, err = ListTasksOnHold(db)
	case A_OPEN, A_DUE:
		tl, err = ListTasksOpen(db)
	case A_OVERDUE:
		tl, err = ListTasksOverdue(db)
	}
	if err != nil {
		return msg, err
	}
	w := tabwriter.NewWriter(os.Stdout, 4, 0, 2, ' ', 0)
	switch (*p).Args[0] {
	case A_ALL, A_DELETED:
		fmt.Fprintf(w, "%s\tID\tGroup\tStatus\tDate Due\tDate Completed\tSummary%s\n", C_BOLD, C_RESET)
	case A_COMPLETED:
		fmt.Fprintf(w, "%s\tID\tGroup\tStatus\tDate Completed\tSummary%s\n", C_BOLD, C_RESET)
	default:
		fmt.Fprintf(w, "%s\tID\tGroup\tStatus\tDate Due\tSummary%s\n", C_BOLD, C_RESET)
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
		case A_ALL, A_DELETED:
			fmt.Fprintf(w, "\t%d\t%s\t%s\t%s\t%s\t%s\n", t.Id, t.Group.Name, statusMap[t.Status], tDue, tCompleted, t.Summary)
		case A_COMPLETED:
			fmt.Fprintf(w, "\t%d\t%s\t%s\t%s\t%s\n", t.Id, t.Group.Name, statusMap[t.Status], tCompleted, t.Summary)
		default:
			fmt.Fprintf(w, "\t%d\t%s\t%s\t%s\t%s\n", t.Id, t.Group.Name, statusMap[t.Status], tDue, t.Summary)
		}
	}
	w.Flush()
	return msg, err
}
func (p *Parser) Reopen(db *sql.DB, conf *Config) (msg string, err error) {
	// Required: K_ID,
	var t Task
	t.Id = p.Kwargs[K_ID].(int)
	if err = t.GetTask(db); err != nil {
		return "", err
	}
	t.Status = STATUS_INPROG
	t.DateCompleted.Valid = false

	if err = t.Update(db); err != nil {
		return "", err
	}

	var bs strings.Builder
	bs.WriteString(fmt.Sprintf("Task resumed: %d - %s", t.Id, t.Summary))

	// Hold subtasks
	if err = t.GetChildren(db); err != nil {
		return msg, err
	}
	if len(t.Children) > 0 {
		plural := ""
		if len(t.Children) > 1 {
			plural = "s"
		}
		bs.WriteString(fmt.Sprintf("\nincluding %d subtask%s:", len(t.Children), plural))

		for _, c := range t.Children {
			(*c).Status = STATUS_INPROG
			(*c).DateCompleted.Valid = false

			if err = (*c).Update(db); err != nil {
				return msg, err
			}
			bs.WriteString(fmt.Sprintf("\n\t%d - %s", (*c).Id, (*c).Summary))
		}
	}
	bs.WriteString("\n")

	return bs.String(), err
}
func (p *Parser) Show(db *sql.DB, conf *Config) (msg string, err error) {
	// Required: K_ID,
	var t Task
	t.Id = p.Kwargs[K_ID].(int)
	if err = t.GetTask(db); err != nil {
		return "", err
	}
	if err = t.GetChildren(db); err != nil {
		return "", err
	}
	var bs strings.Builder
	bs.WriteString(fmt.Sprintf("%sTASK %d%s\n", C_BOLD, t.Id, C_RESET))
	if t.SysStatus == SYS_DELETED {
		bs.WriteString(fmt.Sprintf("\t%s%sDELETED%s\n", C_RED, C_BOLD, C_RESET))
	}
	bs.WriteString(fmt.Sprintf("\tSummary:\t%s\n", t.Summary))
	bs.WriteString(fmt.Sprintf("\tPriority:\t%d\n", t.Priority))
	bs.WriteString(fmt.Sprintf("\tProject\t\t%s\n", t.Group.Name))
	bs.WriteString(fmt.Sprintf("\tStatus:\t\t%s\n", statusMap[t.Status]))
	if t.DateDue.Valid {
		bs.WriteString(fmt.Sprintf("\tDue Date:\t%s\n", t.DateDue.Time.Format("2006-01-02")))
	}
	if t.Description != "" {
		bs.WriteString(fmt.Sprintf("\tDescription:\t%s\n", t.Description))
	}
	if t.Status == STATUS_COMPLETED {
		bs.WriteString(fmt.Sprintf("\tCompleted:\t %s\n", t.DateCompleted.Time.Format("2006-01-02")))
		bs.WriteString(fmt.Sprintf("\tComment:\t%s\n", t.ClosingComment))
	}
	if t.Parent.Id != 0 {
		bs.WriteString(fmt.Sprintf("\tParent: %d - %s\n", t.Parent.Id, t.Parent.Summary))
	}
	if len(t.Children) > 0 {
		bs.WriteString("Linked tasks:\n")
		for _, c := range t.Children {
			bs.WriteString(fmt.Sprintf(" - %d: %s\n", (*c).Id, (*c).Summary))
		}
	}
	bs.WriteString("\n")
	bs.WriteString(fmt.Sprintf("\tCreated:\t %s\n", t.DateCreated.Time.Format("2006-01-02 15:04:05")))
	bs.WriteString(fmt.Sprintf("\tUpdated:\t %s\n", t.DateUpdated.Time.Format("2006-01-02 15:04:05")))
	return bs.String(), err
}

func (p *Parser) Update(db *sql.DB, conf *Config) (msg string, err error) {
	// Required: K_ID,
	//  Optional: K_DATEDUE, K_GROUP, K_DESCRIPTION, K_PRIORITY, K_SUMMARY, K_PARENT
	var t Task
	t.Id = p.Kwargs[K_ID].(int)
	if err = t.GetTask(db); err != nil {
		return "", err
	}
	// Optional values
	if value, ok := p.Kwargs[K_DATEDUE]; ok {
		t.DateDue = value.(sql.NullTime)
	}
	if value, ok := p.Kwargs[K_DESCRIPTION]; ok {
		t.Description = value.(string)
	}
	if value, ok := p.Kwargs[K_PROJECT]; ok {
		t.Group.Name = value.(string)
	}
	if value, ok := p.Kwargs[K_PRIORITY]; ok {
		t.Priority = value.(int)
	}
	if value, ok := p.Kwargs[K_SUMMARY]; ok {
		t.Summary = value.(string)
	}
	// If parent is set, get parent task

	if value, ok := p.Kwargs[K_PARENT]; ok {
		t.Parent.Id = value.(int)
		if value.(int) != 0 {
			if err = t.Parent.GetTask(db); err != nil {
				return "", err
			}
			// If provided due date is later than parent's, use parent's
			if t.Parent.DateDue.Valid && (t.DateDue.Time.After(t.Parent.DateDue.Time) ||
				!t.DateDue.Valid) {
				t.DateDue = t.Parent.DateDue
			}
			// Parent group overrides provided value
			t.Group.Id = t.Parent.Group.Id
			t.Group.Name = t.Parent.Group.Name

			// If priority is lower than parent's, use parent's
			if t.Priority > t.Parent.Priority {
				t.Priority = t.Parent.Priority
			}
		}
	}
	if err = t.Group.GetGroup(db); err != nil {
		if err = t.Group.Add(db); err != nil {
			return "", err
		}
	}

	if err = t.Update(db); err != nil {
		return "", err
	}

	msg = fmt.Sprintf("Task updadted: %d - %s", t.Id, t.Summary)
	return msg, err
}
func (p *Parser) Version(db *sql.DB, conf *Config) (msg string, err error) {
	msg = fmt.Sprintf("TODO CLI\t::\tversion: %s\n", VERSION_APP)
	return msg, nil
}
