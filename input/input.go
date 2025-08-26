package input

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Game struct {
	Cells      *[][]string  `json:"cells"`
	Conditions *[]Condition `json:"conditions"`
	Dominoes   *[]Domino    `json:"dominoes"`
}

type Condition struct {
	Expression *string `json:"expression"`
	Operand    *int    `json:"operand"`
	Cells      []conditionCell
}

type conditionCell struct {
	X *int `json:"x"`
	Y *int `json:"y"`
}

type Domino struct {
	Val1 *int `json:"val1"`
	Val2 *int `json:"val2"`
}

func ReadFile(filename string) (*Game, error) {
	// load the game board from input
	inputJSON, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.New("Error: fail to read JSON file")
	}

	g := new(Game)
	if err := json.Unmarshal(inputJSON, g); err != nil {
		return nil, fmt.Errorf("Error: JSON parse failed with following error - %w", err)
	}

	return g, nil
}
