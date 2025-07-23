package main

const (
	VERSION      = "2.0.0"
	DB_VERSION   = "fc4247a1b0cbfbfcff856355ba6ca58e"
	TRIG_VERSION = "84bd34b5905b45c7f0bde28a28c64143"
	VIEW_VERSION = "6dfef78f0d49b3f5843f655ecdff3965"
)

// Config defaults
const (
	CONFIG_FILE        = "todo_config.yml"
	CONFIG_DBNAME      = "todo.db"
	CONFIG_DATE_FORMAT = "2006-01-02"
	CONFIG_GROUP_NAME  = "General"
)

// Default values
const DEFAULT_GROUP = 1

const (
	STATUS_NEW       = 1
	STATUS_INPROG    = 2
	STATUS_HOLD      = 3
	STATUS_COMPLETED = 4
)
const (
	SYS_DELETED = 0
	SYS_ACTIVE  = 1
)
const (
	PRIORITY_CRIT = 5
	PRIORITY_HIGH = 50
	PRIORITY_MED  = 500
	PRIORITY_LOW  = 5000
)

// CLI arguments
const (
	X_NIL = iota
	// Verbs
	V_ADD
	V_COMPLETE
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
	A_NEW
	A_OPEN
	A_ONHOLD
	A_OVERDUE
)

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
