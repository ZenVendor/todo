package main

import (
	"database/sql"
	"fmt"
)

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
	// Required K_ID
	// Optional A_ALL
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
func (p *Parser) Help() (err error) {
	return err
}
func (p *Parser) Hold() (err error) {
	// Required: K_ID
	return err
}
func (p *Parser) List() (err error) {
	//	Optional: A_ALL, A_COMPLETED, A_DELETED, A_DUE, A_GROUPS, A_INPROGRESS, A_ONHOLD, A_OPEN, A_OVERDUE
	return err
}
func (p *Parser) Reopen() (err error) {
	// Required: K_ID,
	return err
}
func (p *Parser) Show() (err error) {
	// Required: K_ID,
	return err
}
func (p *Parser) Update() (err error) {
	// Required: K_ID,
	//  Optional: K_DUEDATE, K_GROUP, K_DESCRIPTION, K_PRIORITY, K_SUMMARY, K_PARENT
	return err
}
func (p *Parser) Version() (err error) {
	return err
}
