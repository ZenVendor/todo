package main

import (
	"errors"
	"strconv"
	"testing"
)

var testCases = []struct {
	inputArgs []string
	result    Parser
	err       error
}{
	{[]string{}, Parser{Verb{}, []int{}, map[int]interface{}{}}, ErrNoArguments},
	{[]string{"open"}, Parser{Verb{}, []int{}, map[int]interface{}{}}, ErrInvalidVerb},
	{[]string{"version"}, Parser{verbs[10], []int{}, map[int]interface{}{}}, nil},
	{[]string{"v"}, Parser{verbs[10], []int{}, map[int]interface{}{}}, nil},
	{[]string{"version", "-o"}, Parser{verbs[10], []int{}, map[int]interface{}{}}, ErrInvalidArgument},
	{[]string{"help"}, Parser{verbs[5], []int{}, map[int]interface{}{}}, nil},
	{[]string{"h"}, Parser{verbs[5], []int{}, map[int]interface{}{}}, nil},
	{[]string{"help", "--reset"}, Parser{verbs[5], []int{}, map[int]interface{}{}}, ErrInvalidArgument},
	{[]string{"show"}, Parser{verbs[8], []int{}, map[int]interface{}{}}, nil},
	{[]string{"show", "--id=12"}, Parser{verbs[8], []int{}, map[int]interface{}{}}, ErrVerbRequiresValue},
	{[]string{"s", "12"}, Parser{verbs[8], []int{}, map[int]interface{}{K_ID: 12}}, nil},
	{[]string{"configure"}, Parser{verbs[2], []int{}, map[int]interface{}{}}, nil},
	{[]string{"configure", "--due"}, Parser{verbs[2], []int{}, map[int]interface{}{}}, ErrInvalidArgument},
	{[]string{"reopen"}, Parser{verbs[7], []int{}, map[int]interface{}{}}, ErrVerbRequiresValue},
	{[]string{"r", "12"}, Parser{verbs[7], []int{}, map[int]interface{}{K_ID: 12}}, nil},
	{[]string{"reopen", "--due"}, Parser{verbs[7], []int{}, map[int]interface{}{}}, ErrInvalidArgument},
	{[]string{"list"}, Parser{verbs[6], []int{}, map[int]interface{}{}}, ErrVerbRequiresArgument},
	{[]string{"l", "--due"}, Parser{verbs[6], []int{A_DUE}, map[int]interface{}{}}, nil},
	{[]string{"l", "--reset"}, Parser{verbs[6], []int{}, map[int]interface{}{}}, ErrInvalidArgument},
	{[]string{"list", "-d", "-o"}, Parser{verbs[6], []int{A_DUE}, map[int]interface{}{}}, ErrTooManyArguments},
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
