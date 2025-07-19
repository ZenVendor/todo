package main

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func validateInt(value string) (interface{}, error) {
	return strconv.Atoi(value)
}

func validateGroup(value string) (interface{}, error) {
	group := strings.TrimSpace(value)
	if len([]rune(group)) > 24 {
		return "", fmt.Errorf("%w: %d/24", ErrStringLength, len([]rune(group)))
	}
	return group, nil
}
func validatePriority(value string) (interface{}, error) {
	return strconv.Atoi(value)
}

func validateSummary(value string) (interface{}, error) {
	short := strings.TrimSpace(value)
	if len([]byte(short)) > 255 {
		return "", fmt.Errorf("%w: %d/255", ErrStringLength, len([]byte(short)))
	}
	return short, nil
}

func validateString(value string) (interface{}, error) {
	return strings.TrimSpace(value), nil
}

func validateDate(value string) (interface{}, error) {
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
