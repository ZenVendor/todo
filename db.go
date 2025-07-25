package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func NullNow() sql.NullTime {
	return sql.NullTime{Time: time.Now(), Valid: true}
}

func CheckDB(db *sql.DB) (err error) {
	var csDB string
	query := "select cs_db from SysVersion where id = 0;"
	if err := db.QueryRow(query).Scan(&csDB); err != nil {
		return ErrDBVersion
	}
	if csDB != VERSION_DB {
		return ErrDBVersion
	}
	return err
}

func (t *Task) Add(db *sql.DB) (err error) {
	parentId := 0
	if t.Parent != nil {
		parentId = t.Parent.Id
	}

	result, err := db.Exec(`
        INSERT INTO Task (
            summary
            , priority
            , date_due
            , description
            , project_id
            , parent_id
        ) VALUES (?, ?, ?, ?, ?, ?);
        `,
		t.Summary,
		t.Priority,
		t.DateDue,
		t.Description,
		t.Project.Id,
		parentId,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	t.Id = int(id)
	return err
}

func (p *Project) Add(db *sql.DB) (err error) {
	result, err := db.Exec(
		"INSERT INTO Project (project_name) VALUES (?);",
		p.Name,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	p.Id = int(id)
	return err
}

func (t *Task) GetTask(db *sql.DB) (err error) {
	query := `
        SELECT 
            id 
            , summary
            , priority
            , date_due
            , date_completed
            , description
            , closing_comment
            , status
            , project_id
            , project_name
            , parent_id
            , sys_created
            , sys_updated
            , sys_status
        FROM task_list_all
        WHERE id = ?;
    `
	var parentId int
	if err = db.QueryRow(query, t.Id).Scan(
		&t.Id,
		&t.Summary,
		&t.Priority,
		&t.DateDue,
		&t.DateCompleted,
		&t.Description,
		&t.ClosingComment,
		&t.Status,
		&t.Project.Id,
		&t.Project.Name,
		&parentId,
		&t.DateCreated,
		&t.DateUpdated,
		&t.SysStatus,
	); err != nil {
		return err
	}
	if parentId != 0 {
		t.Parent = &Task{Id: parentId}
		if err = t.Parent.GetTask(db); err != nil {
			return err
		}
	}
	return err
}
func (t *Task) GetChildren(db *sql.DB) (err error) {
	query := `
        SELECT 
            id 
            , summary
            , priority
            , date_due
            , date_completed
            , description
            , closing_comment
            , status
            , project_id
            , project_name
            , parent_id
            , sys_created
            , sys_updated
            , sys_status
        FROM task_list_all
        WHERE parent_id = ?;
    `
	rows, err := db.Query(query, t.Id)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var ct Task
		ct.Parent = &Task{}

		if err = rows.Scan(
			&ct.Id,
			&ct.Summary,
			&ct.Priority,
			&ct.DateDue,
			&ct.DateCompleted,
			&ct.Description,
			&ct.ClosingComment,
			&ct.Status,
			&ct.Project.Id,
			&ct.Project.Name,
			&ct.Parent.Id,
			&ct.DateCreated,
			&ct.DateUpdated,
			&ct.SysStatus,
		); err != nil {
			return err
		}
		t.Children = append(t.Children, &ct)
	}
	err = rows.Err()
	return err
}

func (p *Project) GetById(db *sql.DB) (err error) {
	err = db.QueryRow(
		"SELECT id, project_name FROM Project WHERE id = ?;",
		p.Id,
	).Scan(&p.Id, &p.Name)
	return err
}

func (p *Project) GetByName(db *sql.DB) (err error) {
	err = db.QueryRow(
		"SELECT id, project_name FROM Project WHERE project_name = ?;",
		p.Name,
	).Scan(&p.Id, &p.Name)
	return err
}

func (p *Project) GetProject(db *sql.DB) (err error) {
	if p.Name != "" {
		return p.GetByName(db)
	}
	return p.GetById(db)
}

func (t Task) Update(db *sql.DB) (err error) {
	parentId := 0
	if t.Parent != nil {
		parentId = t.Parent.Id
	}
	query := `
        UPDATE Task SET 
            summary = ?
            , priority = ?
            , date_due = ?
            , date_completed = ?
            , description = ?
            , closing_comment = ?
            , status = ?
            , project_id = ?
            , parent_id = ?
            , sys_updated = current_timestamp
        WHERE id = ?;
    `
	_, err = db.Exec(
		query,
		t.Summary,
		t.Priority,
		t.DateDue,
		t.DateCompleted,
		t.Description,
		t.ClosingComment,
		t.Status,
		t.Project.Id,
		parentId,
		t.Id,
	)
	return err
}

func (p Project) Update(db *sql.DB) (err error) {
	_, err = db.Exec(
		"UPDATE Project SET project_name = ?, sys_updated = current_timestamp WHERE id = ?;",
		p.Name,
		p.Id,
	)
	return err
}

func (t Task) Delete(db *sql.DB) (err error) {
	_, err = db.Exec(
		"UPDATE Task SET sys_status = 0, sys_updated = current_timestamp WHERE id = ?;",
		t.Id,
	)
	return err
}

func (t Task) Undelete(db *sql.DB) (err error) {
	_, err = db.Exec(
		"UPDATE Task SET sys_status = 1, sys_updated = current_timestamp WHERE id = ?;",
		t.Id,
	)
	return err
}

func (t Task) DeleteChildren(db *sql.DB) (rows int, err error) {
	result, err := db.Exec(`
		UPDATE Task 
        SET sys_status = 0, sys_updated = current_timestamp 
        WHERE parent_id = ?
        `,
		t.Id,
	)
	if err != nil {
		return 0, err
	}
	r, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(r), err
}

func (p Project) Delete(db *sql.DB) (err error) {
	_, err = db.Exec(
		"UPDATE Project SET sys_status = 0, sys_updated = current_timestamp WHERE id = ?;",
		p.Id,
	)
	return err
}

func ListTasks(db *sql.DB, view int) (tl TaskList, err error) {
	var vName string
	switch view {
	default:
		vName = "task_list_open"
	case A_ALL:
		vName = "task_list_all"
	case A_COMPLETED:
		vName = "task_list_completed"
	case A_DELETED:
		vName = "task_list_deleted"
	case A_INPROGRESS:
		vName = "task_list_in_progress"
	case A_NEW:
		vName = "task_list_new"
	case A_ONHOLD:
		vName = "task_list_on_hold"
	case A_OVERDUE:
		vName = "task_list_overdue"
	}
	query := fmt.Sprintf(`
        SELECT id 
            , summary
            , priority
            , date_due
            , date_completed
            , status
            , project_id
            , project_name 
        FROM %s;
        `, vName)
	rows, err := db.Query(query)
	if err != nil {
		return tl, err
	}
	defer rows.Close()

	for rows.Next() {
		var t Task

		err = rows.Scan(
			&t.Id,
			&t.Summary,
			&t.Priority,
			&t.DateDue,
			&t.DateCompleted,
			&t.Status,
			&t.Project.Id,
			&t.Project.Name,
		)
		if err != nil {
			return tl, err
		}
		tl = append(tl, &t)
	}
	err = rows.Err()
	return tl, err
}

func (c *Counts) GetCounts(db *sql.DB) (err error) {
	query := `
        SELECT 
            count_all
            , count_new
            , count_in_progress
            , count_on_hold
            , count_completed
            , count_open
            , count_overdue
        FROM status_counts
    `
	err = db.QueryRow(query).Scan(
		&c.All,
		&c.New,
		&c.InProgress,
		&c.OnHold,
		&c.Completed,
		&c.Open,
		&c.Overdue,
	)
	return err
}

func (c *Counts) ProjectCounts(projectId int, db *sql.DB) (err error) {
	query := `
        SELECT 
            count_all
            , count_new
            , count_in_progress
            , count_on_hold
            , count_completed
            , count_open
            , count_overdue
        FROM project_counts
        WHERE id = ?;
    `
	err = db.QueryRow(query, projectId).Scan(
		c.All,
		c.New,
		c.InProgress,
		c.OnHold,
		c.Completed,
		c.Open,
		c.Overdue,
	)
	return err
}
