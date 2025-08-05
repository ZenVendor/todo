package main

var priorityMap = map[string]int{
	"none":     PRIORITY_NONE,
	"low":      PRIORITY_LOW,
	"l":        PRIORITY_LOW,
	"medium":   PRIORITY_MED,
	"mid":      PRIORITY_MED,
	"m":        PRIORITY_MED,
	"high":     PRIORITY_HIGH,
	"hi":       PRIORITY_HIGH,
	"h":        PRIORITY_HIGH,
	"critical": PRIORITY_CRIT,
	"crit":     PRIORITY_CRIT,
	"c":        PRIORITY_CRIT,
}

var statusMap = map[int]string{
	STATUS_NEW:       "New",
	STATUS_INPROG:    "In progress",
	STATUS_HOLD:      "On hold",
	STATUS_COMPLETED: "Completed",
}

var verbMap = map[string]int{
	"add":      V_ADD,
	"complete": V_COMPLETE,
	"count":    V_COUNT,
	"delete":   V_DELETE,
	"help":     V_HELP,
	"hold":     V_HOLD,
	"list":     V_LIST,
	"reopen":   V_REOPEN,
	"show":     V_SHOW,
	"start":    V_START,
	"update":   V_UPDATE,
	"undelete": V_UNDELETE,
	"version":  V_VERSION,
}
var argMap = map[string]int{
	"--all":         A_ALL,
	"--comment":     A_COMMENT,
	"--completed":   A_COMPLETED,
	"--deleted":     A_DELETED,
	"--description": A_DESCRIPTION,
	"--desc":        A_DESCRIPTION,
	"--inprog":      A_INPROGRESS,
	"--new":         A_NEW,
	"--onhold":      A_ONHOLD,
	"--open":        A_OPEN,
	"--overdue":     A_OVERDUE,
}
var kwargMap = map[string]int{
	"--due":      K_DATEDUE,
	"--parent":   K_PARENT,
	"--pid":      K_PARENT,
	"--priority": K_PRIORITY,
	"--pty":      K_PRIORITY,
	"--project":  K_PROJECT,
	"--proj":     K_PROJECT,
	"--summary":  K_SUMMARY,
	"--sum":      K_SUMMARY,
}
var validatorMap = map[int]func(*Config, string) (interface{}, error){
	K_DATEDUE:  (*Config).validateDate,
	K_ID:       (*Config).validateInt,
	K_PARENT:   (*Config).validateInt,
	K_PRIORITY: (*Config).validatePriority,
	K_PROJECT:  (*Config).validateProject,
	K_SUMMARY:  (*Config).validateSummary,
}
