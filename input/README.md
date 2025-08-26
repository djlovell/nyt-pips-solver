# Game Solver Input Specification (JSON Format)

This document describes how to create JSON input files for the game solver.  

> [!WARNING]
> This README *was* AI assisted and generated based on a sample input file and a previous hard coded Go input specification format. I'm sorry.

### Coordinate System

- The game board is described as a **uniform grid** of cells.
- Coordinates are expressed as **(x, y)**:
  - **x** = column index (0-based, increasing to the right)  
  - **y** = row index (0-based, increasing downward)

---

# JSON Structure

An input file must define the following top-level attributes:

- **`cells`** - A 2D grid describing playable and blocked cells.  
- **`conditions`** - Rules that apply to groups of cells.  
- **`dominoes`** - A list of available domino pieces.

Example:
```json
{
  "cells": [...],
  "conditions": [...],
  "dominoes": [...]
}
```

## Cells (`cells`)
Represents the game board layout.
It must be a rectangular 2D array where each element is one of:
- "O" - playable cell (domino can be placed here)
- "X" - blocked cell (not part of the game)

Example:
```json
{
  "cells": [
    ["X", "O", "O"],
    ["X", "O", "O"],
    ["O", "O", "X"]
  ]
}
```

## Conditions (`conditions`)
Conditions are rules that specific groups of cells must meet after dominoes are placed.
Each condition object has the following fields:

`expression` (string) - The type of rule.

- "N" - sum of domino values (pips) covering the cells must equal a specified value (operand)
- "=" - all domino faces must be equal
- "!=" - all domino faces must be distinct
- ">N" - sum must be greater than operand
- "<N" - sum must be less than operand

`operand` (optional, integer) - Used only for expressions that compare cell values to a scalar value (i.e. expressions that contain "N").

`cells` - Array of `{ "x": int, "y": int }` objects listing the positions of cells whose values must satisfy the condition.

Example:

```json
{
  "conditions": [
    {
      "expression": "N",
      "operand": 0,
      "cells": [{ "x": 0, "y": 2 }]
    },
    {
      "expression": "=",
      "cells": [
        { "x": 1, "y": 1 },
        { "x": 1, "y": 2 }
      ]
    },
    {
      "expression": "N",
      "operand": 10,
      "cells": [
        { "x": 2, "y": 0 },
        { "x": 2, "y": 1 }
      ]
    }
  ]
}
```

### Notes:
- A condition may cover any number of cells.
- Not every cell needs a condition.
- Conditions may overlap.

## Dominoes (`dominoes`)
Represents the set of dominoes available to solve the puzzle.
Each domino is an object with:
- `val1` - integer from 0-6
- `val2` - integer from 0-6

*Note: Polarity does not matter for input (e.g. 5,3 is the same as 3,5)*

### Example:

```json
{
  "dominoes": [
    { "val1": 5, "val2": 5 },
    { "val1": 0, "val2": 2 },
    { "val1": 2, "val2": 3 }
  ]
}
```

## Full Example
```json
{
  "cells": [
    ["X", "O", "O"],
    ["X", "O", "O"],
    ["O", "O", "X"]
  ],
  "conditions": [
    {
      "expression": "N",
      "operand": 0,
      "cells": [{ "x": 0, "y": 2 }]
    },
    {
      "expression": "=",
      "cells": [
        { "x": 1, "y": 1 },
        { "x": 1, "y": 2 }
      ]
    },
    {
      "expression": "N",
      "operand": 10,
      "cells": [
        { "x": 2, "y": 0 },
        { "x": 2, "y": 1 }
      ]
    }
  ],
  "dominoes": [
    { "val1": 5, "val2": 5 },
    { "val1": 0, "val2": 2 },
    { "val1": 2, "val2": 3 }
  ]
}
```
