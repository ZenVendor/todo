package main

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (c *Config) validateInt(value string) (interface{}, error) {
	if value == "" {
		return 0, nil
	}
	return strconv.Atoi(value)
}

func (c *Config) validateProject(value string) (interface{}, error) {
	proj := strings.TrimSpace(value)
	nameLen := len([]rune(proj))
	if nameLen > c.ProjectNameLength {
		return "", fmt.Errorf("%w: Project: %d/%d", ErrStringLength, nameLen, c.ProjectNameLength)
	}
	return proj, nil
}
func (c *Config) validatePriority(value string) (interface{}, error) {
	if value == "" {
		return PRIORITY_NONE, nil
	}
	if intPty, err := strconv.Atoi(value); err == nil {
		return intPty, nil
	}
	if intPty, ok := priorityMap[value]; ok {
		return intPty, nil
	}
	return 0, fmt.Errorf("%w: %s", ErrInvalidPriority, value)
}

func (c *Config) validateSummary(value string) (interface{}, error) {
	summary := strings.TrimSpace(value)
	sumLen := len([]rune(summary))
	if sumLen > c.SummaryLength {
		return "", fmt.Errorf("%w: Summary: %d/%d", ErrStringLength, sumLen, c.SummaryLength)
	}
	return summary, nil
}

func (c *Config) validateDate(value string) (interface{}, error) {
	if value == "" {
		return value, nil
	}
	if _, err := parseDate(value); err == nil {
		return value, nil
	}
	if _, err := parseDuration(sql.NullTime{Time: time.Now(), Valid: true}, value); err == nil {
		return value, nil
	}
	return "", ErrInvalidDate
}

func (t *Task) setDueDate(value string) {
	if value == "" {
		t.DateDue = sql.NullTime{Time: time.Now(), Valid: false}
		return
	}
	if date, err := parseDate(value); err == nil {
		t.DateDue = date
		return
	}
	if date, err := parseDuration(t.DateDue, value); err == nil {
		t.DateDue = date
		return
	}
}

func parseDate(value string) (sql.NullTime, error) {
	re := regexp.MustCompile(`(\d{4})-{0,1}(\d{2})-{0,1}(\d{2})`)
	result := re.FindAllStringSubmatch(value, -1)
	if len(result) == 0 {
		err := fmt.Errorf("%w: %s", ErrInvalidDate, value)
		return sql.NullTime{Time: time.Now(), Valid: false}, err
	}
	if len(result[0]) != 4 {
		err := fmt.Errorf("%w: %s", ErrInvalidDate, value)
		return sql.NullTime{Time: time.Now(), Valid: false}, err
	}

	dateStr := fmt.Sprintf("%s-%s-%s", result[0][1], result[0][2], result[0][3])
	date, err := time.ParseInLocation(time.DateOnly, dateStr, time.Local)
	return sql.NullTime{Time: date, Valid: true}, err
}

func parseDuration(oldDate sql.NullTime, dur string) (newDate sql.NullTime, err error) {
	re := regexp.MustCompile(`([+-]{0,1})(\d{0,3})([dwmy]{0,1})`)
	result := re.FindAllStringSubmatch(dur, -1)
	if len(result) == 0 {
		err := fmt.Errorf("%w: %s", ErrInvalidDuration, dur)
		return sql.NullTime{Time: time.Now(), Valid: false}, err
	}
	if len(result) > 1 {
		err := fmt.Errorf("%w: %s", ErrInvalidDuration, dur)
		return sql.NullTime{Time: time.Now(), Valid: false}, err
	}
	vals := result[0][1:]
	plus := true
	if vals[0] == "-" {
		plus = false
	}
	num := 1
	if vals[1] != "" {
		num, _ = strconv.Atoi(vals[1])
	}
	if !plus {
		num = -num
	}
	duration := "d"
	if vals[2] != "" {
		duration = vals[2]
	}
	var years, months, days int
	switch duration {
	case "d":
		days = num
	case "w":
		days = 7 * num
	case "m":
		months = num
	case "y":
		years = num
	}
	date := oldDate.Time
	if !oldDate.Valid {
		date = time.Now()
	}
	newDate = sql.NullTime{
		Time:  date.AddDate(years, months, days),
		Valid: true,
	}
	return newDate, nil
}
