#!/bin/bash

# Check if a workflow file name is provided
if [ $# -eq 0 ]; then
    echo "Please provide a workflow file name as an argument."
    echo "Usage: $0 <workflow-file-name.yml>"
    exit 1
fi

WORKFLOW_FILE="$1"
BATCH_SIZE=100
TOTAL_DELETED=0

while true; do
    # List workflow runs and extract run IDs
    RUN_IDS=$(gh run list --workflow "$WORKFLOW_FILE" --limit $BATCH_SIZE --json databaseId --jq '.[].databaseId')

    # Check if there are any runs to delete
    if [ -z "$RUN_IDS" ]; then
        break
    fi

    # Delete each run
    for ID in $RUN_IDS; do
        if gh run delete "$ID"; then
            TOTAL_DELETED=$((TOTAL_DELETED + 1))
            echo "Deleted run $ID"
        else
            echo "Failed to delete run $ID"
        fi
    done
done

echo "Total workflow runs deleted: $TOTAL_DELETED"
