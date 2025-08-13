# 7EP-0011: Re-tighten golangci-lint rules in stages

Status: 🟡 Draft
Assignees: CC
Difficulty: 2
Created: 2025-08-13

## Summary
Gradually re-tighten golangci-lint rules to improve code quality without blocking delivery. Start from a minimal passing baseline (post-7EP-0002) and re-enable stricter checks in small, measurable increments with fast feedback.

## Motivation
- CI stability: PR #11 established Quality checks but initial strictness caused slowdowns. We need a balanced path.
- Developer velocity: Avoid long feedback loops by staging rule enforcement.
- Quality goals: Catch real issues (err handling, inefficiency, dead code) while minimizing noise and false positives.

## Goals
- Keep CI green at all times while incrementally increasing lint coverage
- Document rationale and scope for each rule introduced
- Ensure local reproduction (make lint) matches CI behavior
- Minimize noisy/false-positive rules or confine them via targeted excludes

## Non-Goals
- Enforce security-irrelevant stylistic nits that slow velocity with little ROI
- Introduce rules that don’t match our project patterns without prior discussion

## Current Baseline
- golangci-lint action pinned: v1.60.3
- Passing subset enabled; typecheck coordination issues previously observed with yaml/progressbar
- gosec integrated and green

## Staged Plan

Stage 0 (baseline) — already in main
- Keep CI and gosec passing. Ensure golangci-lint runs with `--modules-download-mode=mod` in CI.

Stage 1 — Core safety and correctness (low noise)
- Enable: govet, staticcheck, ineffassign, errcheck, gosimple, unused
- Excludes:
  - Line-level `// nolint:errcheck` or `_ =` for intentional best-effort calls (progress bars, cleanup)
  - Targeted excludes for generated or test-only code that cannot conform (document why)
- Acceptance: CI passes in all packages; 0 new warnings introduced

Stage 2 — API hygiene and maintainability
- Enable: revive (selected rules), misspell, unparam (review noise), prealloc (informational)
- Configure revive with a light rule set: exported, unused-params, error-naming, receiver-naming
- Acceptance: No more than 5 actionable findings; fix or exclude with rationale

Stage 3 — Complexity and style (opt-in)
- Evaluate: cyclop, gocognit, gocyclo (thresholds set generously, e.g., 20–25)
- Evaluate: dupl for clones in tests only (informational)
- Acceptance: Findings do not fail CI initially; report-only for 1–2 weeks, then gate if actionable fixes land

Stage 4 — Project-specific rules
- Add custom exclude-rules for known patterns (managed paths, deterministic test RNG) to keep gosec/golangci-lint quiet where appropriate
- Document patterns in docs/development/linting.md

## Implementation Details

- Config file: `.golangci.yml` maintained in repo
- CI: `golangci-lint-action` with `args: --timeout=5m --modules-download-mode=mod`
- Local: `make deps && make lint` should mirror CI behavior
- Exclusions: Prefer line-level `// nolint:<linter>` with short justification; avoid broad file-level disables

### Initial .golangci.yml (Stage 1)
```yaml
run:
  timeout: 5m
  issues-exit-code: 1
  tests: true
  modules-download-mode: mod

linters:
  enable:
    - govet
    - staticcheck
    - ineffassign
    - errcheck
    - gosimple
    - unused

issues:
  exclude-rules:
    - path: cmd/(create|test)\.go
      linters: [errcheck]
      text: 'bar\.(Add|Finish)\('
```

## Risks and Mitigations
- False positives (typecheck in mixed module/test contexts): Use modules-download-mode and targeted excludes; prefer fixing imports over excluding
- Developer friction: Stage increases; communicate changes in PR; provide autofix guidance

## Success Metrics
- CI duration remains < 5 minutes for Quality job
- Reduction in linter findings over time; track per-stage delta
- No rollback of stages for 4 consecutive weeks

## Rollout
- Create PR with Stage 1 config and fixes
- After 1–2 weeks of stable green builds, enable Stage 2 and address findings
- Iterate until Stage 3/4 are at desired levels

## Alternatives
- Make all rules required at once (causes delays and broad exclusions)
- Disable linting entirely (quality regression risk)

## References
- 7EP-0002: CI Integration & Automation
- golangci-lint docs: https://golangci-lint.run
- gosec: https://github.com/securego/gosec

