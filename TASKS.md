# Intent CLI Tasks

## Active Tasks (Now)
- [x] Publish and verify `v0.3.18` (release + Homebrew + smoke checks)
- [x] Align `README.md`, `STATUS.md`, plan, and task tracking
- [x] Add concise EPICS/PHASES high-level overview
- [x] Define and enforce "maintenance mode" scope for this repo
- [x] Create a short handoff checklist for cross-repo focus

## Maintenance Mode Scope (Proposed)
Only accept:
- Security fixes
- Reliability regressions in core lifecycle
- CI/release breakages
- Critical docs drift affecting usage or release operations

## Deferred Until Explicitly Approved
- New non-critical feature expansions
- Broad refactors without operational impact
- Additional command families outside lifecycle scope

## Scope Gate Implementation
- Pull request template: `/Users/nectios/nekotori/dev/intent-ecosystem/intent-cli/.github/PULL_REQUEST_TEMPLATE.md`
- Issue templates:
  - `/Users/nectios/nekotori/dev/intent-ecosystem/intent-cli/.github/ISSUE_TEMPLATE/bug-report.yml`
  - `/Users/nectios/nekotori/dev/intent-ecosystem/intent-cli/.github/ISSUE_TEMPLATE/change-request.yml`
  - `/Users/nectios/nekotori/dev/intent-ecosystem/intent-cli/.github/ISSUE_TEMPLATE/config.yml`

## Handoff
- Checklist: `/Users/nectios/nekotori/dev/intent-ecosystem/intent-cli/docs/HANDOFF_CHECKLIST.md`
