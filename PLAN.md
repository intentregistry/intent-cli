# Intent CLI Execution Plan

## Planning Rule
Only plan work that moves one of these outcomes:
1. Core lifecycle reliability
2. Release/distribution reliability
3. Team throughput across repositories

## Phase Plan

| Phase | Outcome | Exit Criteria | Status |
|---|---|---|---|
| P1 Core Lifecycle | Commands for create/run/package/publish/install/search/test/login | Commands implemented and tested | DONE |
| P2 Integrity & Trust | Verifiable package integrity/signatures | `intent verify` + signing flows validated | DONE |
| P3 DevEx & Local Ops | Fast local iteration and support docs | Mock API + guides + examples + completion | DONE |
| P4 CI/CD & Release | Predictable releases across OS + Homebrew | CI green, release workflow green, brew smoke green | DONE |
| P5 Hardening & Drift Control | Docs/tests/code remain synchronized | No drift after releases; concise operational docs | IN PROGRESS |
| P6 Maintenance Mode | Minimal active scope for this repo | Backlog trimmed to critical-only items | NEXT |

## Current Execution Order
1. Complete P5 drift control and simplify docs.
2. Enter P6 maintenance mode with strict scope gate.
3. Redirect engineering capacity to other repos.
