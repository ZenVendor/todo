package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"time"
)

const INVALID = 0

const (
	CMD_LIST = iota + 1
	CMD_COUNT
	CMD_ADD
	CMD_UPDATE
	CMD_DELETE
	CMD_COMPLETE
	CMD_REOPEN
	CMD_PROMPT
	CMD_HELP
	CMD_VERSION
	CMD_PREPARE
	CMD_RESET
	CMD_SHOW
)

func MapCommand() map[string]int {
	return map[string]int{
		"list":      CMD_LIST,
		"l":         CMD_LIST,
		"count":     CMD_COUNT,
		"add":       CMD_ADD,
		"a":         CMD_ADD,
		"update":    CMD_UPDATE,
		"u":         CMD_UPDATE,
		"delete":    CMD_DELETE,
		"del":       CMD_DELETE,
		"complete":  CMD_COMPLETE,
		"c":         CMD_COMPLETE,
		"reopen":    CMD_REOPEN,
		"open":      CMD_REOPEN,
		"prompt":    CMD_PROMPT,
		"help":      CMD_HELP,
		"h":         CMD_HELP,
		"--help":    CMD_HELP,
		"-h":        CMD_HELP,
		"version":   CMD_VERSION,
		"v":         CMD_VERSION,
		"--version": CMD_VERSION,
		"-v":        CMD_VERSION,
		"prepare":   CMD_PREPARE,
		"prep":      CMD_PREPARE,
		"reset":     CMD_RESET,
		"show":      CMD_SHOW,
		"s":         CMD_SHOW,
	}
}

func MapCommandDescription() map[int]string {
	return map[int]string{
		CMD_LIST:     "list",
		CMD_COUNT:    "count",
		CMD_ADD:      "add",
		CMD_UPDATE:   "update",
		CMD_DELETE:   "delete",
		CMD_COMPLETE: "complete",
		CMD_REOPEN:   "reopen",
		CMD_PREPARE:  "prepare",
		CMD_RESET:    "reset",
	}
}

const (
	SW_OPEN = iota + 1
	SW_CLOSED
	SW_ALL
	SW_OVERDUE
	SW_DESCRIPTION
	SW_DUE
	SW_PRIORITY
	SW_GROUP
	SW_DEFAULT
	SW_LOCAL
)

func MapArgument() map[string]int {
	return map[string]int{
		"--open":        SW_OPEN,
		"--completed":   SW_CLOSED,
		"--closed":      SW_CLOSED,
		"-c":            SW_CLOSED,
		"--all":         SW_ALL,
		"-a":            SW_ALL,
		"--overdue":     SW_OVERDUE,
		"--od":          SW_OVERDUE,
		"-o":            SW_OVERDUE,
		"--description": SW_DESCRIPTION,
		"--desc":        SW_DESCRIPTION,
		"--due":         SW_DUE,
		"--priority":    SW_PRIORITY,
		"--pty":         SW_PRIORITY,
		"--group":       SW_GROUP,
		"--local":       SW_LOCAL,
	}
}
func MapArgumentDescription() map[int]string {
	return map[int]string{
		SW_OPEN:        "open",
		SW_CLOSED:      "closed",
		SW_ALL:         "all",
		SW_OVERDUE:     "overdue",
		SW_DESCRIPTION: "description",
		SW_DUE:         "due date",
		SW_PRIORITY:    "priority",
		SW_GROUP:       "group",
		SW_DEFAULT:     "default",
		SW_LOCAL:       "local",
	}
}

func MapPriorityValues() map[string]int {
	return map[string]int{
		"urgent":    0,
		"important": 25,
		"standard":  50,
		"low":       75,
	}
}
func MapPriorityStrings() map[int]string {
	return map[int]string{
		0:  "urgent",
		25: "important",
		50: "standard",
		75: "low",
	}
}

type ArgVal struct {
	sw    int
	value string
}
type ArgVals []ArgVal

func ParseArgs(defaultCmd, defaultSw int) (cmd int, sw int, argvals ArgVals, err error) {
	args := os.Args[1:]

	if len(args) == 0 {
		cmd = defaultCmd
		sw = defaultSw
		return
	}

	if len(args) > 0 {
		cmd = MapCommand()[args[0]]
		if cmd == INVALID {
			err = fmt.Errorf("Invalid command: %s\n", args[0])
			return
		}
	}

	if len(args) > 1 {
		swSet := false
		for i := 1; i < len(args); i++ {
			arg := MapArgument()[args[i]]
			if arg == INVALID {
				if i > 1 {
					err = fmt.Errorf("Invalid argument: %s\n", args[i])
					return
				} else {
					argvals = append(argvals, ArgVal{SW_DEFAULT, args[i]})
				}
				continue
			}
			if slices.Contains([]int{SW_OPEN, SW_CLOSED, SW_ALL, SW_OVERDUE, SW_LOCAL}, arg) {
				if swSet {
					err = fmt.Errorf("Multiple switches: %s, %s\n", MapArgumentDescription()[sw], MapArgumentDescription()[arg])
					return
				}
				sw = arg
				swSet = true
				continue
			}
			if slices.Contains([]int{SW_DESCRIPTION, SW_DUE, SW_PRIORITY, SW_GROUP}, arg) && len(args) > i+1 {
				argvals = append(argvals, ArgVal{arg, args[i+1]})
				i++
			}
		}
	}
	return
}

func (vals ArgVals) ValueIsSet(sw int) bool {
	idx := slices.IndexFunc(vals, func(v ArgVal) bool {
		return v.sw == sw
	})
	if idx == -1 {
		return false
	}
	return true
}

func (vals ArgVals) ReadValue(cmd, sw int, dateFormat string) interface{} {
	idx := slices.IndexFunc(vals, func(v ArgVal) bool {
		return v.sw == sw
	})
	if idx == -1 {
		log.Printf("Value not found")
		return nil
	}
	switch sw {
	case SW_DESCRIPTION, SW_GROUP:
		return vals[idx].value
	case SW_PRIORITY:
		n, err := strconv.Atoi(vals[idx].value)
		if err != nil {
			log.Printf("Invalid PRIORITY value: %s", vals[idx].value)
			return nil
		}
		return n
	case SW_DUE:
		d, err := time.Parse(dateFormat, vals[idx].value)
		if err != nil {
			log.Printf("Invalid DUE value: %s", vals[idx].value)
			return nil
		}
		return sql.NullTime{Time: d, Valid: true}
	case SW_DEFAULT:
		switch cmd {
		case CMD_ADD:
			return vals[idx].value
		case CMD_UPDATE, CMD_DELETE, CMD_REOPEN, CMD_COMPLETE, CMD_SHOW:
			n, err := strconv.Atoi(vals[idx].value)
			if err != nil {
				log.Printf("Invalid ID value: %s", vals[idx].value)
				return nil
			}
			return n
		}
	}
	return nil
}

func HumanDue(due time.Time, dateFormat string) string {

	today, _ := time.Parse(dateFormat, time.Now().Format(dateFormat))
	_, tWeek := today.ISOWeek()
	_, dueWeek := due.ISOWeek()

	days := int(due.Sub(today).Hours() / 24)
	weeks := dueWeek - tWeek
	if weeks < 0 {
		weeks = 52 + weeks
	}
	months := due.Month() - today.Month()
	if months < 0 {
		months = 12 + months
	}

	result := due.Format(dateFormat)

	switch {
	case days < 0:
		break
	case days == 0:
		result = "today"
	case days == 1:
		result = "tomorrow"
	case days == 2:
		result = "in two days"
	case weeks == 0:
		result = "this week"
	case weeks == 1:
		result = "next week"
	case weeks == 2:
		result = "in two weeks"
	case months == 0:
		result = "this month"
	case months == 1:
		result = "next month"
	}
	return result
}
