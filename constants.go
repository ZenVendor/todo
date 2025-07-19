package main

const VERSION = "1.0.0"

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
	INVALID_ARG = iota
	V_ADD
	V_COMPLETE
	V_CONFIGURE
	V_COUNT
	V_DELETE
	V_HELP
    V_HOLD
	V_LIST
	V_REOPEN
	V_SHOW
	V_UPDATE
	V_VERSION
	K_COMMENT
	K_DUEDATE
	K_GROUP
	K_ID
	K_DESCRIPTION
	K_PRIORITY
	K_SUMMARY
	K_PARENT
	A_ALL
    A_COMPLETED
    A_DUE
    A_GROUPS
    A_INPROGRESS
	A_LOCAL
    A_NEW
	A_ONGOING
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
		"help":      V_HELP,
		"h":         V_HELP,
        "hold":      V_HOLD,
        "ho":       V_HOLD,
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
		"--all":     A_ALL,
		"-a":        A_ALL,
		"--closed":  A_COMPLETED,
		"-c":        A_COMPLETED,
		"--due":     A_DUE,
		"-d":        A_DUE,
        "--groups":  A_GROUPS,
        "-g":        A_GROUPS,
        "--inprogress": A_INPROGRESS,
        "-ip":      A_INPROGRESS,
		"--local":   A_LOCAL,
        "--new":    A_NEW,
        "-n":    A_NEW,
		"--ongoing":    A_ONGOING,
		"-o":        A_ONGOING,
        "--onhold": A_ONHOLD,
        "-h": A_ONHOLD,
		"--overdue": A_OVERDUE,
		"-v":        A_OVERDUE,
		"--reset":   A_RESET,
	}
var kwargMap = map[string]int{
		"--comment":  K_COMMENT,
		"-c":         K_COMMENT,
		"--due":      K_DUEDATE,
		"-d":         K_DUEDATE,
		"--group":    K_GROUP,
		"-g":         K_GROUP,
		"--id":       K_ID,
		"--description":     K_DESCRIPTION,
        "--parent":   K_PARENT,
		"--priority": K_PRIORITY,
		"-p":         K_PRIORITY,
		"--summary":    K_SUMMARY,
		"-s":         K_SUMMARY,
	}

func kwargValidatorMap() map[int]func(string) (interface{}, error) {
	return map[int]func(string) (interface{}, error){
		K_COMMENT:  validateString,
		K_DUEDATE:  validateDate,
		K_GROUP:    validateGroup,
		K_ID:       validateInt,
		K_DESCRIPTION:     validateString,
		K_PRIORITY: validatePriority,
		K_SUMMARY:    validateShort,
		K_PARENT:   validateInt,
	}
}

var verbs = Verbs{
	Verb{V_ADD, K_SUMMARY, []int{K_DUEDATE, K_GROUP, K_DESCRIPTION, K_PRIORITY, K_ID, K_PARENT}, 1},
	Verb{V_COMPLETE, K_ID, []int{K_COMMENT}, 0},
	Verb{V_CONFIGURE, 0, []int{A_LOCAL, A_RESET}, 2},
	Verb{V_COUNT, 0, []int{A_ALL, A_COMPLETED, A_DUE, A_ONGOING, A_OVERDUE}, 1},
	Verb{V_DELETE, K_ID, []int{}, 0},
	Verb{V_HELP, 0, []int{}, 0},
	Verb{V_LIST, 0, []int{A_ALL, A_COMPLETED, A_DUE, A_ONGOING, A_OVERDUE}, 1},
	Verb{V_REOPEN, K_ID, []int{}, 0},
	Verb{V_SHOW, K_ID, []int{}, 0},
	Verb{V_UPDATE, K_ID, []int{K_DUEDATE, K_GROUP, K_DESCRIPTION, K_PRIORITY, K_SUMMARY, K_ID, K_PARENT}, 1},
	Verb{V_VERSION, 0, []int{}, 0},
}
