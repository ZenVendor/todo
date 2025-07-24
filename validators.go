package main

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (p *Parser) validateInt(value string) (interface{}, error) {
	if value == "" {
		return 0, nil
	}
	return strconv.Atoi(value)
}

func (p *Parser) validateProject(value string) (interface{}, error) {
	proj := strings.TrimSpace(value)
	nameLen := len([]rune(proj))
	if nameLen > p.Conf.ProjectNameLength {
		return "", fmt.Errorf("%w: Project: %d/%d", ErrStringLength, nameLen, p.Conf.ProjectNameLength)
	}
	return proj, nil
}
func (p *Parser) validatePriority(value string) (interface{}, error) {
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

func (p *Parser) validateSummary(value string) (interface{}, error) {
	summary := strings.TrimSpace(value)
	sumLen := len([]rune(summary))
	if sumLen > p.Conf.SummaryLength {
		return "", fmt.Errorf("%w: Summary: %d/%d", ErrStringLength, sumLen, p.Conf.SummaryLength)
	}
	return summary, nil
}

func (p *Parser) validateDate(value string) (interface{}, error) {
	if value == "" {
		return sql.NullTime{Time: time.Now(), Valid: false}, nil
	}
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
