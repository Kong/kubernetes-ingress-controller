#!/bin/bash

# For the last (NPAGES * PERPAGE) workflow runs in REPO, obtain logs for all jobs matching JOBFILTER.
# Requires `hub` and `jq`.

set -e

REPO=Kong/kubernetes-ingress-controller
NPAGES=5
PERPAGE=100
WORKFLOWNAME="Integration Tests"
JOBFILTER="select(.name == \"integration-test-postgres\" or .name == \"integration-test-dbless\")"

for page in $(seq 1 $NPAGES); do
	echo "Getting runs (page $page out of $NPAGES)..."
	hub api "/repos/$REPO/actions/runs?per_page=$PERPAGE&page=$page" | jq ".workflow_runs" > "01_repo_runs_$page.json"
done

echo "Concatenating the list of workflow runs..."
jq -s 'add' $(for page in $(seq 1 $NPAGES); do echo "01_repo_runs_$page.json"; done) > "02_repo_runs.json"

echo "Filtering for workflow name $WORKFLOWNAME..."
jq ".[] | select(.name == \"$WORKFLOWNAME\")" "02_repo_runs.json" > "03_repo_runs_for_workflow.json"

echo "Extracting workflow run IDs..."
jq -r ".id" "03_repo_runs_for_workflow.json" > "04_workflow_run_ids.txt"

nruns="$(wc -l 04_workflow_run_ids.txt | cut -d' ' -f1)"
echo "Retrieving workflow run metadata ($nruns runs...)"
for id in $(cat 04_workflow_run_ids.txt); do
	echo "Retrieving metadata for workflow run $id"
	hub api "https://api.github.com/repos/Kong/kubernetes-ingress-controller/actions/runs/$id/jobs" > "05_workflow_run_$id.json"
done

echo "Joining jobs into one big file..."
jq -s 'map(.jobs[])' 05_workflow_run_*.json > "06_all_jobs.json"

echo "Filtering jobs..."
jq ".[] | $JOBFILTER" "06_all_jobs.json" > "07_jobs_filtered.json"

echo "Getting logs for each job..."
for jobid in $(jq -r '.id' "07_jobs_filtered.json"); do
	echo "Getting logs for job $jobid..."
	hub api "https://api.github.com/repos/$REPO/actions/jobs/$jobid/logs" > "08_job_$jobid.log"
done
