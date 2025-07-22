package main

import (
	"fmt"
	"slices"
	"strings"
)

func (verbs Verbs) GetVerb(verb int) (Verb, error) {
	idx := slices.IndexFunc(verbs, func(v Verb) bool {
		return v.Verb == verb
	})
	if idx == -1 {
		return Verb{}, ErrInvalidVerb
	}
	return verbs[idx], nil

}

func (p *Parser) Print() {
	fmt.Printf("Verb: %d\n", p.Verb.Verb)
	fmt.Printf("Args:")
	for _, s := range p.Args {
		fmt.Printf(" %d", s)
	}
	fmt.Printf("\nKwArgs:\n")
	for key, value := range p.Kwargs {
		fmt.Printf("%d: %s\n", key, value)
	}
}

func (p *Parser) ToString() string {
	str := fmt.Sprintf("%d -", p.Verb.Verb)
	for _, s := range p.Args {
		str = fmt.Sprintf("%s %d", str, s)
	}
	for key, value := range p.Kwargs {
		str = fmt.Sprintf("%s %d:%v", str, key, value)
	}
	return str
}

func NewParser(defaultVerb int, defaultArgs []int, defaultKwargs map[int]interface{}) (Parser, error) {
	var err error
	var p = Parser{
		Verb:   Verb{},
		Args:   make([]int, 0, 2),
		Kwargs: map[int]interface{}{},
	}
	p.Verb, err = verbs.GetVerb(defaultVerb)
	if err != nil {
		return p, err
	}

	for _, arg := range defaultArgs {
		if !slices.Contains(p.Verb.ValidArgs, arg) {
			return p, ErrInvalidDefaultArg
		}
	}
	for key := range defaultKwargs {
		if !slices.Contains(p.Verb.ValidArgs, key) {
			return p, ErrInvalidDefaultKwarg
		}
	}
	p.Args = defaultArgs
	p.Kwargs = defaultKwargs
	return p, nil
}

func (p *Parser) Parse(args []string) error {
	if len(args) == 0 {
		if p.Verb.Verb == X_NIL {
			return ErrNoArguments
		}
		return nil
	}
	p.Args = make([]int, 0, 2)
	p.Kwargs = map[int]interface{}{}

	var err error
	p.Verb, err = verbs.GetVerb(verbMap[args[0]])
	if err != nil {
		return fmt.Errorf("%w \"%s\"", ErrInvalidVerb, args[0])
	}

	argsStart := 1
	if p.Verb.RequiredValue != 0 {
		if len(args) < 2 {
			return fmt.Errorf("%w: \"%s\"", ErrVerbRequiresValue, args[0])
		}
		verbVal, err := validatorMap[p.Verb.RequiredValue](args[1])
		if err != nil {
			return err
		}
		p.Kwargs[p.Verb.RequiredValue] = verbVal
		argsStart = 2
	}

	for _, arg := range args[argsStart:] {
		kwarg := strings.SplitN(arg, "=", 2)
		if len(kwarg) == 1 {
			sw := argMap[kwarg[0]]
			if sw == X_NIL {
				return fmt.Errorf("%w: %s", ErrInvalidArgument, arg)
			}
			if !slices.Contains(p.Verb.ValidArgs, sw) {
				return fmt.Errorf("%w: %s", ErrInvalidVerbArgument, arg)
			}
			p.Args = append(p.Args, sw)
			if len(p.Args) > p.Verb.MaxArgs {
				return ErrTooManyArguments
			}
		} else {
			key := kwargMap[kwarg[0]]
			if key == X_NIL {
				err := fmt.Errorf("%w: %s", ErrInvalidArgument, arg)
				return err
			}
			value := kwarg[1]
			if !slices.Contains(p.Verb.ValidKwargs, key) {
				return fmt.Errorf("%w: %s", ErrInvalidVerbArgument, arg)
			}
			val, err := validatorMap[key](value)
			if err != nil {
				return err
			}
			p.Kwargs[key] = val
		}
	}
	return nil
}

func (p *Parser) ArgIsPresent(arg int) bool {
	if slices.Contains(p.Args, arg) {
		return true
	}
	return false
}

func (p *Parser) GetArg(index int) int {
	if len(p.Args) == 0 {
		return X_NIL
	}
	if index < 0 || index >= len(p.Args) {
		return X_NIL
	}
	return p.Args[index]
}
