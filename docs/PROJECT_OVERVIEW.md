# Intent CLI Project Overview

## Objective
Ship a reliable CLI for the core IntentRegistry lifecycle:
`create -> run -> package -> verify -> publish -> search -> install`.

The CLI is "done" when this lifecycle is stable, documented, and repeatable across local dev, CI, and release channels (GitHub Releases + Homebrew).

## EPICS / Phases

| Phase | Epic | Goal | Status |
|---|---|---|---|
| P1 | Core CLI Lifecycle | Implement the full command surface for the package lifecycle | DONE |
| P2 | Package Integrity | Deterministic `.itpkg`, signatures, and verification command | DONE |
| P3 | Developer Experience | Init scaffolding, completion, docs, mock API, examples | DONE |
| P4 | CI/CD & Distribution | Green CI, automated release, Homebrew publishing, smoke verification | DONE |
| P5 | Product Hardening | Keep behavior/docs/tests aligned, remove drift, simplify operations | IN PROGRESS |
| P6 | Scale Enablement | Freeze scope, define maintenance mode, free bandwidth for other repos | NEXT |

## Current Focus (P5)
1. Keep docs/plan/tasks/status aligned in every release.
2. Prioritize reliability fixes over new features unless explicitly approved.
3. Track only high-impact work that protects the core lifecycle.

## Out of Scope (until P6 is complete)
- New major command families unrelated to the core lifecycle.
- Large refactors without direct reliability or release impact.
- Nice-to-have UX polish that does not reduce operational risk.
