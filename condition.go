package main

import (
	"fmt"
	"strings"
)

type Condition struct {
	expression      ConditionExp
	operand         int      // goes with some conditions
	cellIdentifiers []string // cell identifiers
}

type ConditionExp int

const (
	ConditionExpSumEquals ConditionExp = iota
	ConditionExpSumLessThan
	ConditionExpSumGreaterThan
	ConditionExpEquivalent
	ConditionExpDistinct
)

func (c Condition) String() string {
	s := "Cell"
	switch len(c.cellIdentifiers) {
	case 0:
		panic("no cells for condition")
	case 1:
		s += " " + c.cellIdentifiers[0]
	case 2:
		s += fmt.Sprintf("s %s and %s", c.cellIdentifiers[0], c.cellIdentifiers[1])
	default: // > 2 cells
		s += fmt.Sprintf("s %s", strings.Join(c.cellIdentifiers, ", "))
		if idx := strings.LastIndex(s, ","); idx != -1 {
			sBeforeAnd := s[0 : idx+1]
			sAfterAnd := s[idx+1:]
			s = sBeforeAnd + " and" + sAfterAnd
		}
	}
	s += " must "
	switch c.expression {
	case ConditionExpSumEquals:
		s += fmt.Sprintf("add up to %d", c.operand)
	case ConditionExpSumLessThan:
		s += fmt.Sprintf("add up to less than %d", c.operand)
	case ConditionExpSumGreaterThan:
		s += fmt.Sprintf("add up to greater than %d", c.operand)
	case ConditionExpEquivalent:
		s += "all be the same"
	case ConditionExpDistinct:
		s += "all be different"
	default:
		panic("unhandled expression type")
	}
	s += "\n"
	return s
}
