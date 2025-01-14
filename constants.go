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
	V_LIST
	V_REOPEN
	V_SHOW
	V_UPDATE
	V_VERSION
	K_COMMENT
	K_DUEDATE
	K_GROUP
	K_ID
	K_LONG
	K_PRIORITY
	K_SHORT
	K_TASKID
	A_ALL
	A_CLOSED
	A_DUE
	A_LOCAL
	A_OPEN
	A_OVERDUE
	A_RESET
)

func verbMap() map[string]int {
	return map[string]int{
		"add":       V_ADD,
		"a":         V_ADD,
		"complete":  V_COMPLETE,
		"c":         V_COMPLETE,
		"configure": V_CONFIGURE,
		"count":     V_COUNT,
		"delete":    V_DELETE,
		"help":      V_HELP,
		"h":         V_HELP,
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
}
func argMap() map[string]int {
	return map[string]int{
		"--all":     A_ALL,
		"-a":        A_ALL,
		"--closed":  A_CLOSED,
		"-c":        A_CLOSED,
		"--due":     A_DUE,
		"-d":        A_DUE,
		"--local":   A_LOCAL,
		"--open":    A_OPEN,
		"-o":        A_OPEN,
		"--overdue": A_OVERDUE,
		"-v":        A_OVERDUE,
		"--reset":   A_RESET,
	}
}
func kwargMap() map[string]int {
	return map[string]int{
		"--comment":  K_COMMENT,
		"-c":         K_COMMENT,
		"--due":      K_DUEDATE,
		"-d":         K_DUEDATE,
		"--group":    K_GROUP,
		"-g":         K_GROUP,
		"--id":       K_ID,
		"--long":     K_LONG,
		"-l":         K_LONG,
		"--priority": K_PRIORITY,
		"-p":         K_PRIORITY,
		"--short":    K_SHORT,
		"-s":         K_SHORT,
		"--taskid":   K_TASKID,
		"-t":         K_TASKID,
	}
}

func kwargValidatorMap() map[int]func(string) (interface{}, error) {
	return map[int]func(string) (interface{}, error){
		K_COMMENT:  validateString,
		K_DUEDATE:  validateDate,
		K_GROUP:    validateGroup,
		K_ID:       validateInt,
		K_LONG:     validateString,
		K_PRIORITY: validatePriority,
		K_SHORT:    validateShort,
		K_TASKID:   validateInt,
	}
}

var verbs = Verbs{
	Verb{V_ADD, K_SHORT, []int{K_DUEDATE, K_GROUP, K_LONG, K_PRIORITY, K_ID, K_TASKID}, 1},
	Verb{V_COMPLETE, K_ID, []int{K_COMMENT}, 0},
	Verb{V_CONFIGURE, 0, []int{A_LOCAL, A_RESET}, 2},
	Verb{V_COUNT, 0, []int{A_ALL, A_CLOSED, A_DUE, A_OPEN, A_OVERDUE}, 1},
	Verb{V_DELETE, K_ID, []int{}, 0},
	Verb{V_HELP, 0, []int{}, 0},
	Verb{V_LIST, 0, []int{A_ALL, A_CLOSED, A_DUE, A_OPEN, A_OVERDUE}, 1},
	Verb{V_REOPEN, K_ID, []int{}, 0},
	Verb{V_SHOW, K_ID, []int{}, 0},
	Verb{V_UPDATE, K_ID, []int{K_DUEDATE, K_GROUP, K_LONG, K_PRIORITY, K_SHORT, K_ID, K_TASKID}, 1},
	Verb{V_VERSION, 0, []int{}, 0},
}
