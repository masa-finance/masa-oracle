#!/bin/bash

set -e

if [ $# -eq 0 ]; then
    echo "Error: No workflow name provided."
    echo "Usage: $0 <workflow-name>"
    exit 1
fi

WORKFLOW_NAME="$1"

echo "Fetching repository information..."
REPO_INFO=$(gh repo view --json nameWithOwner)
REPO=$(echo $REPO_INFO | jq -r .nameWithOwner)
echo "Repository: $REPO"

echo "Searching for workflow: '$WORKFLOW_NAME'"
WORKFLOW_ID=$(gh api "/repos/$REPO/actions/workflows" | jq -r ".workflows[] | select(.name == \"$WORKFLOW_NAME\") | .id")

if [ -z "$WORKFLOW_ID" ]; then
    echo "No workflow found with name '$WORKFLOW_NAME'"
    exit 1
fi

echo "Found workflow with ID: $WORKFLOW_ID"

delete_runs() {
    local page=$1
    echo "Fetching runs (page $page)..."
    RUNS=$(gh api "/repos/$REPO/actions/workflows/$WORKFLOW_ID/runs?per_page=100&page=$page")
    RUN_COUNT=$(echo $RUNS | jq '.workflow_runs | length')
    
    if [ "$RUN_COUNT" -eq 0 ]; then
        return 1
    fi

    echo "Deleting $RUN_COUNT runs..."
    echo $RUNS | jq -r '.workflow_runs[].id' | while read -r run_id; do
        echo "Deleting run $run_id"
        gh api -X DELETE "/repos/$REPO/actions/runs/$run_id"
    done

    return 0
}

page=1
while delete_runs $page; do
    ((page++))
done

echo "All runs deleted for workflow '$WORKFLOW_NAME'"

echo "Deleting the workflow itself..."
gh api -X DELETE "/repos/$REPO/actions/workflows/$WORKFLOW_ID"

echo "Workflow '$WORKFLOW_NAME' and all its runs have been deleted."
