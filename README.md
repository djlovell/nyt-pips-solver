# nyt-pips-solver

## Overview

This is my attempt at writing a Go based solver for the NYT Pips game released on 2025-08-18, and is strictly for fun! Comments, unit tests, design plans, etc., if supplied, are not going to adhere to professional standards.

Game Link: https://www.nytimes.com/games/pips

I take no credit for the Pips game.

## On AI Usage
Besides some documentation work, this was **NOT** an exercise in "Vibe Coding", "Promp Engineering", or even speeding up development time with AI auto-complete (i.e. Copilot). This tooling is exciting and useful, but I wanted to approach this more like a classic, semi closed web coding challenge. Perhaps down the road I might use my resources to explore potential optimizations, create a GUI to input/create game boards, or create fancy terminal visualizations with print statements - tasks I have far less interest in - but for now I am shooting for an unaided MVP that can reliably solve these puzzles. I saw none on Github at the time I started this project.

## Requirements
- a working Go environment (>=1.24.6)
  - download link: https://go.dev/doc/install

## How to use
1. clone the repo
2. copy a known working example JSON file from `/test_files` or create one using the README in `/input`
3. in a terminal, run `go run . -f {{your file name}}.json`
4. test one of the found solutions in your NYT Games App or on the website
