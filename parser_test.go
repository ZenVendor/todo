package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

var conf = Config{
	DBLocation:        ".",
	DBName:            "todo.db",
	DateFormat:        "2006-01-02",
	DefaultProject:    "General",
	ProjectNameLength: 24,
	SummaryLength:     255,
}

var testCases = []struct {
	inputArgs []string
	result    Parser
	err       error
}{
	{[]string{}, Parser{Verb{}, []int{}, map[int]interface{}{}, &conf}, ErrNoArguments},
	{[]string{"open"}, Parser{Verb{}, []int{}, map[int]interface{}{}, &conf}, ErrInvalidVerb},
	{[]string{"version"}, Parser{verbs[10], []int{}, map[int]interface{}{}, &conf}, nil},
	{[]string{"v"}, Parser{verbs[10], []int{}, map[int]interface{}{}, &conf}, nil},
	{[]string{"version", "-o"}, Parser{verbs[10], []int{}, map[int]interface{}{}, &conf}, ErrInvalidArgument},
	{[]string{"help"}, Parser{verbs[5], []int{}, map[int]interface{}{}, &conf}, nil},
	{[]string{"h"}, Parser{verbs[5], []int{}, map[int]interface{}{}, &conf}, nil},
	{[]string{"help", "--reset"}, Parser{verbs[5], []int{}, map[int]interface{}{}, &conf}, ErrInvalidArgument},
	{[]string{"show"}, Parser{verbs[8], []int{}, map[int]interface{}{}, &conf}, nil},
	{[]string{"show", "--id=12"}, Parser{verbs[8], []int{}, map[int]interface{}{}, &conf}, ErrVerbRequiresValue},
	{[]string{"s", "12"}, Parser{verbs[8], []int{}, map[int]interface{}{K_ID: 12}, &conf}, nil},
	{[]string{"configure"}, Parser{verbs[2], []int{}, map[int]interface{}{}, &conf}, nil},
	{[]string{"configure", "--due"}, Parser{verbs[2], []int{}, map[int]interface{}{}, &conf}, ErrInvalidArgument},
	{[]string{"reopen"}, Parser{verbs[7], []int{}, map[int]interface{}{}, &conf}, ErrVerbRequiresValue},
	{[]string{"r", "12"}, Parser{verbs[7], []int{}, map[int]interface{}{K_ID: 12}, &conf}, nil},
	{[]string{"reopen", "--due"}, Parser{verbs[7], []int{}, map[int]interface{}{}, &conf}, ErrInvalidArgument},
	{[]string{"list"}, Parser{verbs[6], []int{}, map[int]interface{}{}, &conf}, ErrVerbRequiresArgument},
	{[]string{"l", "--due"}, Parser{verbs[6], []int{A_DUE}, map[int]interface{}{}, &conf}, nil},
	{[]string{"l", "--reset"}, Parser{verbs[6], []int{}, map[int]interface{}{}, &conf}, ErrInvalidArgument},
	{[]string{"list", "-d", "-o"}, Parser{verbs[6], []int{A_DUE}, map[int]interface{}{}, &conf}, ErrTooManyArguments},
}

func TestParse(t *testing.T) {
	for i, c := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var prs = Parser{}
			err := prs.Parse(c.inputArgs)
			if prs.ToString() != c.result.ToString() && !errors.Is(err, c.err) {
				t.Errorf("Expected: %s, %s || Got: %s, %s", c.result.ToString(), c.err, prs.ToString(), err)
			}
		})
	}

}

// Helper functions
func (p *Parser) ToString() string {
	var bs strings.Builder
	fmt.Fprintf(&bs, "%d -", p.Verb.Verb)
	for _, s := range p.Args {
		fmt.Fprintf(&bs, " %d", s)
	}
	for key, value := range p.Kwargs {
		fmt.Fprintf(&bs, " %d:%v", key, value)
	}
	return bs.String()
}
