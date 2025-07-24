package main

import (
	"fmt"
	"time"
)

func Color(text, color string, bold bool) string {
	if bold {
		color = fmt.Sprintf("%s%s", color, C_BOLD)
	}
	return fmt.Sprintf("%s%s%s", color, text, C_RESET)
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

func DisplayPriority(priority int) string {
	if priority >= 0 && priority < 10 {
		return fmt.Sprintf("Critical [%d]", priority)
	}
	if priority >= 10 && priority < 100 {
		return fmt.Sprintf("High [%d]", priority)
	}
	if priority >= 100 && priority < 1000 {
		return fmt.Sprintf("Medium [%d]", priority)
	}
	if priority >= 1000 && priority < 10000 {
		return fmt.Sprintf("Low [%d]", priority)
	}
	return fmt.Sprintf("Reminder", priority)
}
