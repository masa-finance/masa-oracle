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

echo "Searching for workflows with name: '$WORKFLOW_NAME'"
WORKFLOW_IDS=$(gh api "/repos/$REPO/actions/workflows" | jq -r ".workflows[] | select(.name == \"$WORKFLOW_NAME\") | .id")

if [ -z "$WORKFLOW_IDS" ]; then
    echo "No workflows found with name '$WORKFLOW_NAME'"
    exit 1
fi

echo "Found workflow(s) with ID(s): $WORKFLOW_IDS"

delete_runs() {
    local workflow_id=$1
    local page=$2
    echo "Fetching runs for workflow $workflow_id (page $page)..."
    RUNS=$(gh api "/repos/$REPO/actions/workflows/$workflow_id/runs?per_page=100&page=$page")
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

for WORKFLOW_ID in $WORKFLOW_IDS; do
    echo "Processing workflow ID: $WORKFLOW_ID"
    page=1
    while delete_runs $WORKFLOW_ID $page; do
        ((page++))
    done

    echo "All runs deleted for workflow ID $WORKFLOW_ID"

    echo "Deleting the workflow itself..."
    gh api -X DELETE "/repos/$REPO/actions/workflows/$WORKFLOW_ID"

    echo "Workflow with ID $WORKFLOW_ID has been deleted."
done

echo "All workflows named '$WORKFLOW_NAME' and their runs have been deleted."
