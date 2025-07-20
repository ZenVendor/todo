package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"strings"
)

//go:embed help.txt
var helpString string

func (p *Parser) Add(db *sql.DB) (msg string, err error) {
	// Required: K_SUMMARY
	// Optional: K_DESCRIPTION, K_DUEDATE, K_GROUP, K_PARENT, K_PRIORITY
	var t Task
	t.Summary = p.Kwargs[K_SUMMARY].(string)

	// Optional values
	if value, ok := p.Kwargs[K_DUEDATE]; ok {
		t.DateDue = value.(sql.NullTime)
	}
	if value, ok := p.Kwargs[K_DESCRIPTION]; ok {
		t.Description = value.(string)
	}
	if value, ok := p.Kwargs[K_PRIORITY]; ok {
		t.Priority = value.(int)
	} else {
		t.Priority = PRIORITY_MED
	}

	// If parent is set, get parent task
	// Group is set from the parent, the argument or default
	if value, ok := p.Kwargs[K_PARENT]; ok {
		t.Parent.Id = value.(int)
		if err = t.Parent.GetById(db); err != nil {
			return "", err
		}
		t.Group.Id = t.Parent.Group.Id
		t.Group.Name = t.Parent.Group.Name
	} else {
		if value, ok := p.Kwargs[K_GROUP]; ok {
			t.Group.Name = value.(string)
		} else {
			t.Group.Id = DEFAULT_GROUP
		}
		if err = t.Group.GetByName(db); err != nil {
			if err = t.Group.Add(db); err != nil {
				return "", err
			}
		}
	}

	if err = t.Add(db); err != nil {
		return "", err
	}

	msg = fmt.Sprintf("Added task:\n\tGroup: %s\n\t%d - %s", t.Group.Name, t.Id, t.Summary)
	return msg, nil
}

func (p *Parser) Complete(db *sql.DB) (msg string, err error) {
	// Required: K_ID
	// Optional: K_COMMENT
	var t Task
	t.Id = p.Kwargs[K_ID].(int)
	if err = t.GetById(db); err != nil {
		return "", err
	}
	if value, ok := p.Kwargs[K_COMMENT]; ok {
		t.ClosingComment = value.(string)
	}
	t.DateCompleted = NullNow()
	if err = t.Update(db); err != nil {
		return "", err
	}
	msg = fmt.Sprintf("Completed task: %d - %s", t.Id, t.Summary)
	return msg, nil
}

func (p *Parser) Configure() (err error) {
	// Optional: A_LOCAL, A_RESET
	return err
}
func (p *Parser) Count() (err error) {
	// Optional: A_ALL, A_COMPLETED, A_DUE, A_INPROGRESS, A_ONHOLD, A_OPEN, A_OVERDUE
	return err
}
func (p *Parser) Delete(db *sql.DB) (msg string, err error) {
	// Required: K_ID
	// Optional: A_ALL
	var t Task
	t.Id = p.Kwargs[K_ID].(int)
	if err = t.GetById(db); err != nil {
		return "", err
	}
	if err = t.Delete(db); err != nil {
		return "", err
	}
	msg = fmt.Sprintf("Deleted task: %d - %s", t.Id, t.Summary)

	if p.ArgIsPresent(A_ALL) {
		rows, err := t.DeleteChildren(db)
		if err != nil {
			return "", err
		}
		msg = fmt.Sprintf("%s\n\tand %d subtasks", msg, rows)
	}
	return msg, err
}
func (p *Parser) Group() (err error) {
	// RequiredD: K_ID
	return err
}
func (p *Parser) Help() (msg string, err error) {
	return helpString, nil
}
func (p *Parser) Hold(db *sql.DB) (msg string, err error) {
	// Required: K_ID
	var t Task
	t.Id = p.Kwargs[K_ID].(int)
	if err = t.GetById(db); err != nil {
		return "", err
	}
	t.Status.Id = STATUS_HOLD
	if err = t.Update(db); err != nil {
		return "", err
	}

	msg = fmt.Sprintf("Task put on hold: %d - %s", t.Id, t.Summary)
	return msg, err
}
func (p *Parser) List() (err error) {
	//	Optional: A_ALL, A_COMPLETED, A_DELETED, A_DUE, A_GROUPS, A_INPROGRESS, A_ONHOLD, A_OPEN, A_OVERDUE
	return err
}
func (p *Parser) Reopen(db *sql.DB) (msg string, err error) {
	// Required: K_ID,
	var t Task
	t.Id = p.Kwargs[K_ID].(int)
	if err = t.GetById(db); err != nil {
		return "", err
	}
	t.Status.Id = STATUS_INPROG
	t.DateCompleted.Valid = false
	if err = t.Update(db); err != nil {
		return "", err
	}

	msg = fmt.Sprintf("Task resumed: %d - %s", t.Id, t.Summary)
	return msg, err
}
func (p *Parser) Show(db *sql.DB) (msg string, err error) {
	// Required: K_ID,
	var t Task
	t.Id = p.Kwargs[K_ID].(int)
	if err = t.GetById(db); err != nil {
		return "", err
	}
	var bs strings.Builder
	bs.WriteString(fmt.Sprintf("%sTASK %d%s\n", C_BOLD, t.Id, C_RESET))
	if t.SysStatus == SYS_DELETED {
		bs.WriteString(fmt.Sprintf("\t%s%sDELETED%s\n", C_RED, C_BOLD, t.Id, C_RESET))
	}
	bs.WriteString(fmt.Sprintf("\tGroup:\t%s\n", t.Group.Name))
	bs.WriteString(fmt.Sprintf("\tSummary:\t%s\n", t.Summary))
	bs.WriteString(fmt.Sprintf("\tStatus:\t%s\n", t.Status.Name))
	bs.WriteString(fmt.Sprintf("\tDue Date:\t%s\n", t.DateDue.Time))
	bs.WriteString(fmt.Sprintf("\tPriority:\t%s\n", t.Priority))
	bs.WriteString(fmt.Sprintf("\tDescription:\t%s\n", t.Description))
	if t.Parent != nil {
		bs.WriteString(fmt.Sprintf("\tParent: %d - %s\n", t.Parent.Id, t.Parent.Summary))
	}
	if t.Status.Id == STATUS_COMPLETED {
		bs.WriteString(fmt.Sprintf("\tCompleted: %s\n", t.DateCompleted.Time))
		bs.WriteString(fmt.Sprintf("\tComment: %s\n", t.ClosingComment))
	}
	msg = bs.String()
	return msg, err
}
func (p *Parser) Update(db *sql.DB) (msg string, err error) {
	// Required: K_ID,
	//  Optional: K_DUEDATE, K_GROUP, K_DESCRIPTION, K_PRIORITY, K_SUMMARY, K_PARENT
	var t Task
	t.Id = p.Kwargs[K_ID].(int)
	if err = t.GetById(db); err != nil {
		return "", err
	}
	// Optional values
	if value, ok := p.Kwargs[K_DUEDATE]; ok {
		t.DateDue = value.(sql.NullTime)
	}
	if value, ok := p.Kwargs[K_DESCRIPTION]; ok {
		t.Description = value.(string)
	}
	if value, ok := p.Kwargs[K_PRIORITY]; ok {
		t.Priority = value.(int)
	}
	if value, ok := p.Kwargs[K_SUMMARY]; ok {
		t.Summary = value.(string)
	}
	// If parent is set, get parent task
	// Group is set from the parent or the argument
	if value, ok := p.Kwargs[K_PARENT]; ok {
		t.Parent.Id = value.(int)
		if err = t.Parent.GetById(db); err != nil {
			return "", err
		}
		t.Group.Id = t.Parent.Group.Id
		t.Group.Name = t.Parent.Group.Name
	} else {
		if value, ok := p.Kwargs[K_GROUP]; ok {
			t.Group.Name = value.(string)
			if err = t.Group.GetByName(db); err != nil {
				if err = t.Group.Add(db); err != nil {
					return "", err
				}
			}
		}
	}
	if err = t.Update(db); err != nil {
		return "", err
	}

	msg = fmt.Sprintf("Task updadted: %d - %s", t.Id, t.Summary)
	return msg, err
}
func (p *Parser) Version() (msg string, err error) {
	msg = fmt.Sprintf("TODO CLI\t::\tversion: %s\n", VERSION)
	return msg, nil
}
