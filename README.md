# nyt-pips-solver

[![CI](https://github.com/djlovell/nyt-pips-solver/actions/workflows/ci.yml/badge.svg?branch=main&event=push)](https://github.com/djlovell/nyt-pips-solver/actions/workflows/ci.yml)

## Overview
This is my attempt at writing a Go based solver for the NYT Pips game released on 2025-08-18, and is strictly for fun! Accordingly, comments, unit tests, design plans, etc., if supplied, may not live up to professional standards.

I take no credit for the Pips game. Maybe for the first solver on GitHub (none appeared on Google when I started the project).

Update: As of 8/26/2025, I have seen a few nice ones online that are much more usable!

[NYT Game Link Here](https://www.nytimes.com/games/pips)

## On AI Usage
With the exception of some documentation & GH Actions setup, this was **NOT** an exercise in "Vibe Coding", "Prompt Engineering", or speed running development with enhanced auto-complete (i.e. Copilot). That tooling is exciting and useful, but I wanted to approach this like a classic open-docs coding challenge.

Perhaps down the road I may leverage more resources to explore optimizations, create a proper GUI for generating input files, or crafting fancy terminal visualizations using print statements - tasks I have far less interest in - but for now I am shooting for an unaided MVP that can reliably solve these puzzles.

## Requirements
- a working Go environment (>=1.25.x) - [Download Here](https://go.dev/doc/install)

## Quick Start
1. clone the repo
2. copy a known working example JSON file from `/test_files` or create one using the README in `/input`
3. in a terminal, run `go run . -f {{your file name}}.json`
4. test one of the found solutions in your NYT Games App or on the website

## Feedback
I would love to hear your feedback on my solution, optimization ideas, potential bugs, and the like. Contact info should be on my profile!

## Usage/Contribution
Feel free to use, share, copy (credit appreciated), or whatever else. 

If you want to contribute tests, an input file creation UI, or something else - email me or create a PR.
