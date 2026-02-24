# Intent CLI Handoff Checklist

Use this checklist when switching focus from `intent-cli` to other repositories.

## Release and CI Health
- [ ] Latest `main` commit has green `CI` and `dev-build`
- [ ] If a release was cut, release workflow is green
- [ ] Homebrew formula is updated for the latest released version

## Scope and Planning Hygiene
- [ ] `docs/PROJECT_OVERVIEW.md` still reflects current phase
- [ ] `PLAN.md` reflects current execution order
- [ ] `TASKS.md` contains only active work (no stale items)
- [ ] `STATUS.md` snapshot is current

## Drift Control
- [ ] `README.md` command surface matches implemented CLI behavior
- [ ] `docs/USER_GUIDE.md` and `docs/LOCAL_DEVELOPMENT.md` are consistent with workflows
- [ ] Any behavior change includes tests and doc updates in same PR

## Handoff Note (copy/paste)
"Intent CLI is in maintenance mode scope (security/reliability/CI-release/docs-critical only). Current phase and active tasks are up to date in `docs/PROJECT_OVERVIEW.md`, `PLAN.md`, and `TASKS.md`."
