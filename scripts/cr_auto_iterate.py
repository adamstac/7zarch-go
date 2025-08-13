#!/usr/bin/env python3
import json
import os
import re
import subprocess
from pathlib import Path

try:
    import yaml  # type: ignore
except Exception as e:
    print(f"::error::pyyaml not available: {e}")
    exit(1)

GITHUB_TOKEN = os.environ.get("GITHUB_TOKEN") or os.environ.get("GH_TOKEN")
PR_NUMBER = os.environ.get("PR_NUMBER")
REPO = os.environ.get("REPO")

if not (GITHUB_TOKEN and PR_NUMBER and REPO):
    print("::error::Missing env GITHUB_TOKEN/PR_NUMBER/REPO")
    exit(1)


def gh_json(args):
    cmd = ["gh", "api"] + args
    env = os.environ.copy()
    env["GH_TOKEN"] = GITHUB_TOKEN
    res = subprocess.run(cmd, capture_output=True, text=True, env=env)
    if res.returncode != 0:
        print(res.stdout)
        print(res.stderr)
        raise RuntimeError(f"gh api failed: {' '.join(args)}")
    return json.loads(res.stdout)


def get_pr():
    return gh_json([
        f"repos/{REPO}/pulls/{PR_NUMBER}",
        "-H", "Accept: application/vnd.github+json",
    ])


def list_comments():
    comments = gh_json([
        f"repos/{REPO}/issues/{PR_NUMBER}/comments",
        "-H", "Accept: application/vnd.github+json",
    ])
    # Also review comments if needed later
    return comments


def load_yaml(path: Path):
    if not path.exists():
        return None, False
    with path.open("r", encoding="utf-8") as f:
        data = yaml.safe_load(f)
    return data, True


def dump_yaml(path: Path, data):
    with path.open("w", encoding="utf-8") as f:
        yaml.safe_dump(data, f, sort_keys=False)


# Whitelisted fixes for .coderabbit.yaml

def fix_coderabbit_yaml(path: Path, comments: list) -> bool:
    data, exists = load_yaml(path)
    if not exists or data is None:
        return False

    changed = False

    # Ensure version: 2
    if data.get("version") != 2:
        data["version"] = 2
        changed = True

    reviews = data.get("reviews") or {}

    # Align keys: profile, auto_review.enabled, auto_review.drafts
    prof = reviews.get("profile")
    if prof not in ("chill", "assertive"):
        reviews["profile"] = "chill"
        changed = True

    auto = reviews.get("auto_review") or {}
    if auto.get("enabled") is not True:
        auto["enabled"] = True
        changed = True
    if auto.get("drafts") not in (True, False):
        auto["drafts"] = False
        changed = True
    reviews["auto_review"] = auto

    # path_instructions presence (non-blocking if missing)
    if "path_instructions" not in reviews:
        reviews["path_instructions"] = [
            {"path": "internal/storage/**", "instructions": "Use assertive scrutiny for storage internals (correctness, concurrency, durability)."},
            {"path": "cmd/**", "instructions": "Normal scrutiny; ensure CLI UX and error handling."},
            {"path": "**/*_test.go", "instructions": "Chill profile acceptable for tests/tooling."},
        ]
        changed = True

    data["reviews"] = reviews

    # Remove unsupported top-level keys if present
    for k in list(data.keys()):
        if k in ("fail_conditions", "labels", "comments", "ui"):
            del data[k]
            changed = True

    # summaries placement
    summaries = data.get("summaries") or {"enabled": True, "placement": "PR_BODY"}
    if summaries.get("enabled") is not True:
        summaries["enabled"] = True
        changed = True
    if summaries.get("placement") not in ("PR_BODY", "PR_COMMENT"):
        summaries["placement"] = "PR_BODY"
        changed = True
    data["summaries"] = summaries

    # ignore.files defaults
    ignore = data.get("ignore") or {}
    files = set(ignore.get("files") or [])
    defaults = {"dist/**", "vendor/**", "**/*.sum", "**/*.min.*", ".claude/**", ".coderabbit.yaml"}
    if not defaults.issubset(files):
        files |= defaults
        ignore["files"] = sorted(list(files))
        data["ignore"] = ignore
        changed = True

    if changed:
        dump_yaml(path, data)
    return changed


def main():
    pr = get_pr()
    labels = [l["name"] for l in pr.get("labels", [])]
    if "cr:auto-iterate" not in labels:
        print("No cr:auto-iterate label; exiting.")
        return

    comments = list_comments()

    # Only act on same-repo branches for safety
    head = pr.get("head", {})
    if head.get("repo", {}).get("full_name") != REPO:
        print("PR is from a fork; skipping.")
        return

    root = Path('.')
    any_changes = False

    # Apply .coderabbit.yaml fixes if needed
    any_changes |= fix_coderabbit_yaml(root / '.coderabbit.yaml', comments)

    if any_changes:
        os.environ["CR_CHANGES"] = "true"
    else:
        print("No changes applied.")

if __name__ == "__main__":
    main()

