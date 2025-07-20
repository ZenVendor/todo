package main

const VERSION = "2.0.0"

// Colour ANSI codes
const (
	C_WHITE  = "\033[37m"
	C_GREY   = "\033[38;5;7m"
	C_RED    = "\033[38;5;9m"
	C_ORANGE = "\033[38;5;214m"
	C_YELLOW = "\033[38;5;226m"
	C_GREEN  = "\033[38;5;40m"
	C_BOLD   = "\033[1m"
	C_RESET  = "\033[0m"
)

// CLI arguments
const (
	X_NIL = iota
	// Verbs
	V_ADD
	V_COMPLETE
	V_CONFIGURE
	V_COUNT
	V_DELETE
	V_GROUP
	V_HELP
	V_HOLD
	V_LIST
	V_REOPEN
	V_SHOW
	V_UPDATE
	V_VERSION
	// Key-value args
	K_COMMENT
	K_DUEDATE
	K_GROUP
	K_ID
	K_DESCRIPTION
	K_PRIORITY
	K_SUMMARY
	K_PARENT
	// Switches
	A_ALL
	A_COMPLETED
	A_DELETED
	A_DUE
	A_GROUPS
	A_INPROGRESS
	A_LOCAL
	A_NEW
	A_OPEN
	A_ONHOLD
	A_OVERDUE
	A_RESET
)

var verbMap = map[string]int{
	"add":       V_ADD,
	"a":         V_ADD,
	"complete":  V_COMPLETE,
	"c":         V_COMPLETE,
	"configure": V_CONFIGURE,
	"count":     V_COUNT,
	"delete":    V_DELETE,
	"group":     V_GROUP,
	"g":         V_GROUP,
	"help":      V_HELP,
	"h":         V_HELP,
	"hold":      V_HOLD,
	"ho":        V_HOLD,
	"list":      V_LIST,
	"l":         V_LIST,
	"reopen":    V_REOPEN,
	"r":         V_REOPEN,
	"show":      V_SHOW,
	"s":         V_SHOW,
	"update":    V_UPDATE,
	"u":         V_UPDATE,
	"version":   V_VERSION,
	"v":         V_VERSION,
}
var argMap = map[string]int{
	"--all":        A_ALL,
	"-a":           A_ALL,
	"--completed":  A_COMPLETED,
	"-c":           A_COMPLETED,
	"-deleted":     A_DELETED,
	"--due":        A_DUE,
	"-d":           A_DUE,
	"--inprogress": A_INPROGRESS,
	"-ip":          A_INPROGRESS,
	"--local":      A_LOCAL,
	"--new":        A_NEW,
	"-n":           A_NEW,
	"--open":       A_OPEN,
	"-o":           A_OPEN,
	"--onhold":     A_ONHOLD,
	"-h":           A_ONHOLD,
	"--overdue":    A_OVERDUE,
	"-od":          A_OVERDUE,
	"--reset":      A_RESET,
}
var kwargMap = map[string]int{
	"--comment":     K_COMMENT,
	"-c":            K_COMMENT,
	"--due":         K_DUEDATE,
	"-d":            K_DUEDATE,
	"--group":       K_GROUP,
	"-g":            K_GROUP,
	"--id":          K_ID,
	"--description": K_DESCRIPTION,
	"--parent":      K_PARENT,
	"--priority":    K_PRIORITY,
	"-p":            K_PRIORITY,
	"--summary":     K_SUMMARY,
	"-s":            K_SUMMARY,
}

const (
	PRIORITY_CRIT = 5
	PRIORITY_HIGH = 50
	PRIORITY_MED  = 500
	PRIORITY_LOW  = 5000
)

var priorityMap = map[string]int{
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

const DEFAULT_GROUP = 1
const (
	STATUS_NEW       = 1
	STATUS_INPROG    = 2
	STATUS_HOLD      = 3
	STATUS_COMPLETED = 4
)

var validatorMap = map[int]func(string) (interface{}, error){
	K_COMMENT:     validateString,
	K_DUEDATE:     validateDate,
	K_GROUP:       validateGroup,
	K_ID:          validateInt,
	K_DESCRIPTION: validateString,
	K_PRIORITY:    validatePriority,
	K_SUMMARY:     validateSummary,
	K_PARENT:      validateInt,
}

var verbs = Verbs{
	Verb{
		V_ADD,
		K_SUMMARY,
		[]int{},
		[]int{K_DUEDATE, K_GROUP, K_DESCRIPTION, K_PRIORITY, K_PARENT},
		1,
	},
	Verb{
		V_COMPLETE,
		K_ID,
		[]int{},
		[]int{K_COMMENT},
		0,
	},
	Verb{
		V_CONFIGURE,
		X_NIL,
		[]int{A_LOCAL, A_RESET},
		[]int{},
		2,
	},
	Verb{
		V_COUNT,
		X_NIL,
		[]int{A_ALL, A_COMPLETED, A_DUE, A_INPROGRESS, A_ONHOLD, A_OPEN, A_OVERDUE},
		[]int{},
		1,
	},
	Verb{
		V_DELETE,
		K_ID,
		[]int{A_ALL},
		[]int{},
		0,
	},
	Verb{
		V_GROUP,
		K_ID,
		[]int{},
		[]int{},
		1,
	},
	Verb{
		V_HELP,
		X_NIL,
		[]int{},
		[]int{},
		0,
	},
	Verb{
		V_HOLD,
		K_ID,
		[]int{},
		[]int{},
		0,
	},
	Verb{
		V_LIST,
		X_NIL,
		[]int{A_ALL, A_COMPLETED, A_DELETED, A_DUE, A_GROUPS, A_INPROGRESS, A_ONHOLD, A_OPEN, A_OVERDUE},
		[]int{},
		1,
	},
	Verb{
		V_REOPEN,
		K_ID,
		[]int{},
		[]int{},
		0,
	},
	Verb{
		V_SHOW,
		K_ID,
		[]int{},
		[]int{},
		0,
	},
	Verb{
		V_UPDATE,
		K_ID,
		[]int{},
		[]int{K_DUEDATE, K_GROUP, K_DESCRIPTION, K_PRIORITY, K_SUMMARY, K_PARENT},
		1,
	},
	Verb{
		V_VERSION,
		X_NIL,
		[]int{},
		[]int{},
		0,
	},
}
