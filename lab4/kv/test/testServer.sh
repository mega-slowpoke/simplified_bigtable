#!/bin/bash

iteration=1

# Infinite loop to keep running the test until a FAIL is detected
while true; do
    echo "Iteration #$iteration" >> output.txt
    echo "Running iteration #$iteration"
    
    # Run the go test and append the output to output.txt
    go test -run=TestServer -v >> output.txt 2>&1

    # Check if the output contains the word "FAIL"
    if grep -q "FAIL" output.txt; then
        echo "Test failed on iteration #$iteration. Stopping..."
        break
    fi

    # Increment the iteration counter
    ((iteration++))
done
