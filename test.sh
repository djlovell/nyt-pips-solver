#!/bin/bash

SCRIPT_DIR="$(dirname "$(realpath "$0")")"
TEST_FILE_DIR="test_files"

# determines from solve output if a puzzle was successfully solved
SUCCESS_GREP="grep \"NYT Pips Solver Completed\" | grep -q \"Found\""

test_passed="true"

# test each file
echo -e "Running known working test files...\n"
for file in "$SCRIPT_DIR/$TEST_FILE_DIR/"*.json; do
    echo -e "Checking "$file"...\n"

    run_output=$(go run . --f "$file")
    echo -e "${run_output}\n"

    if echo "$run_output" | (eval $SUCCESS_GREP); then
        echo -e "File success...\n"
    else
        echo -e "File failure...\n"
        test_passed="false"
    fi
done

# return the overall success/failure status
if [[ "$test_passed" == "false" ]]; then
    echo "Result: FAIL"
    exit 1 # failure
fi
echo "Result: PASS"
exit 0 # success
