import os
import sys
import httpx
import zipfile
import io

# --- Configuration ---
GITHUB_TOKEN = os.environ["GITHUB_TOKEN"]
REPO_OWNER = os.environ["REPO_OWNER"]
REPO_NAME = os.environ["REPO_NAME"]
WORKFLOW_RUN_ID = os.environ["WORKFLOW_RUN_ID"]
HEAD_SHA = os.environ["HEAD_SHA"]
API_URL = "https://api.github.com"

# --- Main Logic ---

def get_pr_for_workflow_run(client: httpx.Client, run_id: str, head_sha: str) -> dict | None:
    """Fetches the pull request associated with a workflow run."""
    print(f"Fetching workflow run {run_id}...")
    response = client.get(
        f"/repos/{REPO_OWNER}/{REPO_NAME}/actions/runs/{run_id}"
    )
    response.raise_for_status()
    run_data = response.json()

    # Method 1: Check the pull_requests array in the workflow_run event
    if run_data.get("pull_requests"):
        print("Found PR via workflow_run.pull_requests array.")
        return run_data["pull_requests"][0]
    
    # Method 2: If the above is empty, find PR associated with the head SHA
    print(f"Could not find PR in workflow_run event. Trying to find PR for SHA: {head_sha}")
    # This API is in preview, so we need to provide a custom media type.
    headers = {"Accept": "application/vnd.github.groot-preview+json"}
    response = client.get(
        f"/repos/{REPO_OWNER}/{REPO_NAME}/commits/{head_sha}/pulls",
        headers=headers
    )
    response.raise_for_status()
    prs_data = response.json()

    if not prs_data:
        print(f"No pull requests found for SHA {head_sha}.")
        return None

    print(f"Found {len(prs_data)} PR(s) for SHA {head_sha}. Using the first one.")
    return prs_data[0]


def get_failed_jobs(client: httpx.Client, run_id: str) -> list[dict]:
    """Gets all failed jobs for a given workflow run."""
    print(f"Fetching jobs for workflow run {run_id}...")
    response = client.get(
        f"/repos/{REPO_OWNER}/{REPO_NAME}/actions/runs/{run_id}/jobs"
    )
    response.raise_for_status()
    jobs = response.json().get("jobs", [])
    
    failed_jobs = [
        job for job in jobs if job.get("conclusion") == "failure"
    ]
    print(f"Found {len(failed_jobs)} failed jobs.")
    return failed_jobs

def get_job_logs(client: httpx.Client, job_id: int) -> str:
    """Downloads and extracts logs for a specific job."""
    print(f"Fetching logs for job {job_id}...")
    log_zip_response = client.get(
        f"/repos/{REPO_OWNER}/{REPO_NAME}/actions/jobs/{job_id}/logs"
    )
    log_zip_response.raise_for_status()

    log_content = ""
    with zipfile.ZipFile(io.BytesIO(log_zip_response.content)) as zf:
        for filename in zf.namelist():
            # We are interested in the log file itself, not sub-steps
            if "/" not in filename: 
                with zf.open(filename) as f:
                    try:
                        log_content += f.read().decode("utf-8", errors="ignore")
                        log_content += "\n" # Add a newline for readability between files if multiple
                    except Exception as e:
                        log_content += f"Could not read log file {filename}: {e}\n"
    
    return log_content.strip()

def post_comment(client: httpx.Client, pr_number: int, comment: str):
    """Posts a comment to a pull request."""
    print(f"Posting comment to PR #{pr_number}...")
    response = client.post(
        f"/repos/{REPO_OWNER}/{REPO_NAME}/issues/{pr_number}/comments",
        json={"body": comment},
    )
    try:
        response.raise_for_status()
        print("Successfully posted comment.")
    except httpx.HTTPStatusError as e:
        print(f"Error posting comment: {e.response.status_code}")
        print(f"Response: {e.response.text}")
        sys.exit(1)

def get_log_for_failed_step(full_log: str, step_name: str) -> str:
    """Parses a full log to extract the output of a specific step."""
    log_lines = full_log.splitlines()
    start_marker = f"##[group]Run {step_name}"
    # Sometimes the step name in the log is just the name, not "Run {name}"
    # e.g. for `actions/checkout` it's just "Checkout repository"
    alt_start_marker = f"##[group]{step_name}"
    end_marker = "##[endgroup]"
    
    step_log_lines = []
    in_step_log = False

    for line in log_lines:
        if start_marker in line or alt_start_marker in line:
            in_step_log = True
            continue
        if end_marker in line and in_step_log:
            break
        if in_step_log:
            step_log_lines.append(line)

    if not step_log_lines:
        # Fallback if group markers aren't found, return last 50 lines
        return "\n".join(log_lines[-50:])
        
    return "\n".join(step_log_lines)


def main():
    """Main function to orchestrate the process."""
    headers = {
        "Authorization": f"Bearer {GITHUB_TOKEN}",
        "Accept": "application/vnd.github.v3+json",
    }
    with httpx.Client(base_url=API_URL, headers=headers, timeout=30.0) as client:
        # 1. Find the PR
        pr = get_pr_for_workflow_run(client, WORKFLOW_RUN_ID, HEAD_SHA)
        if not pr:
            sys.exit(0)
        pr_number = pr["number"]

        # 2. Get failed jobs
        failed_jobs = get_failed_jobs(client, WORKFLOW_RUN_ID)
        if not failed_jobs:
            print("Workflow run failed, but no specific jobs were marked as failed. Exiting.")
            sys.exit(0)

        # 3. Build the comment
        comment_body = f"## 🤖 GitHub Actions Failure Analysis for run <{failed_jobs[0]['run_url']}|#{failed_jobs[0]['run_attempt']}>\n\n"
        comment_body += "A CI job failed. Here is a summary of the failing step:\n\n"

        for job in failed_jobs:
            job_id = job["id"]
            job_name = job["name"]
            
            # Find the failing step
            failing_step = None
            for step in job.get("steps", []):
                if step.get("conclusion") == "failure":
                    failing_step = step
                    break
            
            if not failing_step:
                comment_body += f"### ⚠️ Job: `{job_name}`\n"
                comment_body += f"Could not determine the exact failing step. <{job['html_url']}|View Full Log>\n\n"
                continue

            full_logs = get_job_logs(client, job_id)
            step_logs = get_log_for_failed_step(full_logs, failing_step['name'])

            comment_body += f"### ❌ Job: `{job_name}` / Step: `{failing_step['name']}`\n\n"
            comment_body += "```\n"
            comment_body += step_logs.strip()
            comment_body += "\n```\n"
            comment_body += f"<{job['html_url']}|View Full Log>\n\n"


        # 4. Post the comment
        post_comment(client, pr_number, comment_body)


if __name__ == "__main__":
    main()
