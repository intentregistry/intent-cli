# Intent CLI Status

Last updated: 2026-02-24

## Project Objective
Deliver a stable CLI for the IntentRegistry lifecycle:
`create -> run -> package -> verify -> publish -> search -> install`.

## Snapshot
- Release: `v0.3.18` published
- CI: green (`CI`, `dev-build`, `release`)
- Distribution: GitHub Releases + Homebrew validated
- Scope gate: PR and issue templates enforce maintenance-mode classification
- Current phase: **P5 - Hardening & Drift Control**

## Phase Progress
| Phase | Status |
|---|---|
| P1 Core Lifecycle | DONE |
| P2 Integrity & Trust | DONE |
| P3 DevEx & Local Ops | DONE |
| P4 CI/CD & Release | DONE |
| P5 Hardening & Drift Control | IN PROGRESS |
| P6 Maintenance Mode | NEXT |

## Current Priorities
1. Keep code/tests/docs/release behavior aligned with zero drift.
2. Restrict changes to high-impact reliability work.
3. Transition this repo to maintenance mode to free capacity for other repos.

## Source of Truth Docs
- Overview: `/Users/nectios/nekotori/dev/intent-ecosystem/intent-cli/docs/PROJECT_OVERVIEW.md`
- Plan: `/Users/nectios/nekotori/dev/intent-ecosystem/intent-cli/PLAN.md`
- Tasks: `/Users/nectios/nekotori/dev/intent-ecosystem/intent-cli/TASKS.md`
