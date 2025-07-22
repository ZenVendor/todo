package main

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func NullNow() sql.NullTime {
	return sql.NullTime{Time: time.Now(), Valid: true}
}

func CheckDB(db *sql.DB) (err error) {
	query := "select cs_db, cs_trigger, cs_view from SysVersion;"

	var csDB, csTrig, csView string
	if err := db.QueryRow(query).Scan(&csDB, &csTrig, &csView); err != nil {
		return ErrDBVersion
	}
	if csDB != DB_VERSION || csTrig != TRIG_VERSION || csView != VIEW_VERSION {
		return ErrDBVersion
	}
	return err
}

func (t *Task) Add(db *sql.DB) (err error) {
	result, err := db.Exec(`
        INSERT INTO Task (
            summary
            , priority
            , date_due
            , description
            , group_id
            , parent_id
        ) VALUES (?, ?, ?, ?, ?, ?);
        `,
		t.Summary,
		t.Priority,
		t.DateDue,
		t.Description,
		t.Group.Id,
		t.Parent.Id,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	t.Id = int(id)
	return err
}

func (g *Group) Add(db *sql.DB) (err error) {
	result, err := db.Exec(
		"INSERT INTO TaskGroup (group_name) VALUES (?);",
		g.Name,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	g.Id = int(id)
	return err
}

func (t *Task) GetById(db *sql.DB) (err error) {
	if t.Group == nil {
		t.Group = &Group{}
	}
	if t.Status == nil {
		t.Status = &Status{}
	}
	if t.Parent == nil {
		t.Parent = &Task{}
	}
	query := `
        SELECT 
            id 
            , summary
            , priority
            , date_due
            , date_completed
            , description
            , closing_comment
            , status_id
            , status_name
            , group_id
            , group_name
            , parent_id
            , sys_created
            , sys_updated
            , sys_status
        FROM task_list_all
        WHERE id = ?;
    `
	if err = db.QueryRow(query, t.Id).Scan(
		&t.Id,
		&t.Summary,
		&t.Priority,
		&t.DateDue,
		&t.DateCompleted,
		&t.Description,
		&t.ClosingComment,
		&t.Status.Id,
		&t.Status.Name,
		&t.Group.Id,
		&t.Group.Name,
		&t.Parent.Id,
		&t.DateCreated,
		&t.DateUpdated,
		&t.SysStatus,
	); err != nil {
		return err
	}
	if t.Parent.Id != 0 {
		if err = t.Parent.GetById(db); err != nil {
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
            , status_id
            , status_name
            , group_id
            , group_name
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
		ct.Status = &Status{}
		ct.Group = &Group{}
		ct.Parent = &Task{}

		if err = rows.Scan(
			&ct.Id,
			&ct.Summary,
			&ct.Priority,
			&ct.DateDue,
			&ct.DateCompleted,
			&ct.Description,
			&ct.ClosingComment,
			&ct.Status.Id,
			&ct.Status.Name,
			&ct.Group.Id,
			&ct.Group.Name,
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

func (g *Group) GetById(db *sql.DB) (err error) {
	err = db.QueryRow(
		"SELECT id, group_name FROM TaskGroup WHERE id = ?;",
		g.Id,
	).Scan(&g.Id, &g.Name)
	return err
}

func (g *Group) GetByName(db *sql.DB) (err error) {
	err = db.QueryRow(
		"SELECT id, group_name FROM TaskGroup WHERE group_name = ?;",
		g.Name,
	).Scan(&g.Id, &g.Name)
	return err
}

func (g *Group) GetGroup(db *sql.DB) (err error) {
	if g.Name != "" {
		return g.GetByName(db)
	}
	return g.GetById(db)
}

func (t Task) Update(db *sql.DB) (err error) {
	query := `
        UPDATE Task SET 
            summary = ?
            , priority = ?
            , date_due = ?
            , date_completed = ?
            , description = ?
            , closing_comment = ?
            , status_id = ?
            , group_id = ?
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
		t.Status.Id,
		t.Group.Id,
		t.Parent.Id,
		t.Id,
	)
	return err
}

func (g Group) Update(db *sql.DB) (err error) {
	_, err = db.Exec(
		"UPDATE TaskGroup SET group_name = ?, sys_updated = current_timestamp WHERE id = ?;",
		g.Name,
		g.Id,
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

func (g Group) Delete(db *sql.DB) (err error) {
	_, err = db.Exec(
		"UPDATE TaskGroup SET sys_status = 0, sys_updated = current_timestamp WHERE id = ?;",
		g.Id,
	)
	return err
}

func ListTasks(db *sql.DB) (tl TaskList, err error) {
	query := `
        SELECT 
            id 
            , summary
            , priority
            , date_due
            , date_completed
            , description
            , closing_comment
            , status_id
            , status_name
            , group_id
            , group_name
            , parent_id
            , sys_created
            , sys_updated
            , sys_status
        FROM task_list_all;
    `
	rows, err := db.Query(query, NullNow())
	if err != nil {
		return tl, err
	}
	defer rows.Close()

	for rows.Next() {
		var t Task
		if err = rows.Scan(
			&t.Id,
			&t.Summary,
			&t.Priority,
			&t.DateDue,
			&t.DateCompleted,
			&t.Description,
			&t.ClosingComment,
			&t.Status.Id,
			&t.Status.Name,
			&t.Group.Id,
			&t.Group.Name,
			&t.Parent.Id,
			&t.DateCreated,
			&t.DateUpdated,
			&t.SysStatus,
		); err != nil {
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
        FROM task_counts
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

func (g *Group) GetCounts(db *sql.DB) (err error) {
	query := `
        SELECT 
            count_all
            , count_new
            , count_in_progress
            , count_on_hold
            , count_completed
            , count_open
            , count_overdue
        FROM group_counts
        WHERE id = ?;
    `
	err = db.QueryRow(query, g.Id).Scan(
		&g.Counts.All,
		&g.Counts.New,
		&g.Counts.InProgress,
		&g.Counts.OnHold,
		&g.Counts.Completed,
		&g.Counts.Open,
		&g.Counts.Overdue,
	)
	return err
}

func (s *Status) GetCounts(db *sql.DB) (err error) {
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
        WHERE id = ?;
    `
	err = db.QueryRow(query, s.Id).Scan(
		&s.Counts.All,
		&s.Counts.New,
		&s.Counts.InProgress,
		&s.Counts.OnHold,
		&s.Counts.Completed,
		&s.Counts.Open,
		&s.Counts.Overdue,
	)
	return err
}
