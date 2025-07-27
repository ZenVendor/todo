package main

import "database/sql"

type Config struct {
	ProjectNameLength int    `yaml:"projectNameLength"`
	SummaryLength     int    `yaml:"summaryLength"`
	DBLocation        string `yaml:"dbLocation"`
	DBName            string `yaml:"dbName"`
	DateFormat        string `yaml:"dateFormat"`
	DefaultProject    string `yaml:"defaultProject"`
	Editor            string `yaml:"editor"`
}

type Task struct {
	Id             int
	Summary        string
	Priority       int
	DateDue        sql.NullTime
	DateCompleted  sql.NullTime
	Description    string
	ClosingComment string
	Status         int
	Project        Project
	Parent         *Task
	Children       TaskList
	DateCreated    sql.NullTime
	DateUpdated    sql.NullTime
	SysStatus      int
}
type TaskList []*Task

type Project struct {
	Id   int
	Name string
}

type Counts struct {
	All        int
	New        int
	InProgress int
	OnHold     int
	Completed  int
	Open       int
	Overdue    int
}

type Verb struct {
	Verb          int
	RequiredValue int
	ValidArgs     []int
	ValidKwargs   []int
	MaxArgs       int
	Call          func(*Parser, *sql.DB) error
}
type Verbs []Verb

type Parser struct {
	Verb   Verb
	Args   []int
	Kwargs map[int]interface{}
	Conf   *Config
}

var verbs = Verbs{
	Verb{
		V_ADD,
		K_SUMMARY,
		[]int{A_DESCRIPTION},
		[]int{K_DATEDUE, K_PARENT, K_PRIORITY, K_PROJECT},
		1,
		(*Parser).Add,
	},
	Verb{
		V_COMPLETE,
		K_ID,
		[]int{A_COMMENT},
		[]int{},
		1,
		(*Parser).Complete,
	},
	Verb{
		V_COUNT,
		X_NIL,
		[]int{A_ALL, A_COMPLETED, A_INPROGRESS, A_NEW, A_ONHOLD, A_OPEN, A_OVERDUE},
		[]int{},
		1,
		(*Parser).Count,
	},
	Verb{
		V_DELETE,
		K_ID,
		[]int{A_ALL},
		[]int{},
		1,
		(*Parser).Delete,
	},
	Verb{
		V_HELP,
		X_NIL,
		[]int{},
		[]int{},
		0,
		(*Parser).Help,
	},
	Verb{
		V_HOLD,
		K_ID,
		[]int{A_DESCRIPTION},
		[]int{},
		1,
		(*Parser).Hold,
	},
	Verb{
		V_LIST,
		X_NIL,
		[]int{A_COMPLETED, A_DELETED, A_INPROGRESS, A_NEW, A_ONHOLD, A_OPEN, A_OVERDUE},
		[]int{},
		1,
		(*Parser).List,
	},
	Verb{
		V_REOPEN,
		K_ID,
		[]int{A_DESCRIPTION},
		[]int{},
		1,
		(*Parser).Reopen,
	},
	Verb{
		V_SHOW,
		K_ID,
		[]int{},
		[]int{},
		0,
		(*Parser).Show,
	},
	Verb{
		V_START,
		K_ID,
		[]int{A_DESCRIPTION},
		[]int{},
		1,
		(*Parser).Start,
	},
	Verb{
		V_UPDATE,
		K_ID,
		[]int{A_COMMENT, A_DESCRIPTION},
		[]int{K_DATEDUE, K_PARENT, K_PRIORITY, K_SUMMARY},
		1,
		(*Parser).Update,
	},
	Verb{
		V_UNDELETE,
		K_ID,
		[]int{},
		[]int{},
		0,
		(*Parser).Undelete,
	},
	Verb{
		V_VERSION,
		X_NIL,
		[]int{},
		[]int{},
		0,
		(*Parser).Version,
	},
}
