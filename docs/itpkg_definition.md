What is .itpkg?

A signed, versioned Intent package: a single distributable archive that bundles one ITML project (its .itml intents plus metadata, policies, tests, and optional assets) so any AI Runtime or the intent tool can install, verify, and execute it deterministically. Think “the npm tarball / Docker image of ITML,” but purpose-first and policy-aware.
	•	Why do we need it?
ITML apps aren’t just code—they’re executable intentions (purpose + contracts + policies). We need a package that preserves those semantics, capabilities, tests, and signatures end-to-end across registries and runtimes. Existing formats don’t model capability gates, privacy/energy policies, and declarative workflows as first-class, portable artifacts.

⸻

What’s inside a .itpkg

A typical .itpkg (conceptually a compressed archive) contains an ITML project:

/project.app.itml          # manifest & routes
/intents/**/*.itml         # intents, layouts, views, components
/schemas/**/*.itml         # data contracts
/policies/**/*.itml        # security/privacy/energy policies
/tests/**/*.itml           # declarative test cases
/.ci/*.itml                # lint/format tasks
/assets/**                  # optional static assets
/itpkg.json                # package metadata (name, version, compatibility)
/SIGNATURE                 # detached or embedded (ed25519)

	•	project.app.itml: entrypoint the runtime reads to resolve imports, routes, and policies.  ￼
	•	Policies & capabilities: encoded declaratively (deny-by-default, allowlists, PII/energy modes).  ￼
	•	Tests: portable, runtime-executable (tests:) for conformance & CI.  ￼
	•	itpkg.json: metadata the toolchain uses when you run intent package / intent publish.

How it’s produced & used (MVP flow):
	1.	intent package → builds the .itpkg, signs it, emits trace info.
	2.	intent publish → uploads artifact to the registry (S3-compatible backend).
	3.	intent install @scope/name → resolves, verifies signature, extracts into project.  ￼

⸻

Similarities & differences vs other package types

Format	What it bundles	Manifest	Integrity	Runtime model	Policies/Capabilities	Tests inside	Primary install UX
.itpkg	ITML intents + schemas + policies + tests	itpkg.json + project.app.itml	Signature + checksums	AI Runtime reads ITML to execute workflows	First-class (security/privacy/energy, deny-by-default)	Yes, declarative tests:	intent install / viewer
npm tarball (.tgz)	JS/TS source/dist	package.json	shasum in lockfiles	Node/JS module system	No standard capability model	Optional scripts/tests, not portable	npm i
Python wheel (.whl)	Python dist	METADATA (PEP 621)	wheel RECORD hash	Python interpreter	No	Tests separate	pip install
Go module (.zip in cache)	Go source	go.mod	sumdb checks	Go toolchain	No	Separate	go get
Docker/OCI image	Filesystem layers + entrypoint	OCI manifest	content digests, sigstore (opt)	Container runtime	Linux caps/seccomp (host-specific)	Not portable across runtimes	docker pull/run
Helm chart (.tgz)	K8s manifests	Chart.yaml	digest/sig (opt)	K8s API server	Limited via K8s policy	Tests as hooks (cluster-coupled)	helm install
JAR (.jar)	JVM bytecode/resources	MANIFEST.MF	signing (opt)	JVM	No	Separate	mvn/gradle
WASM module (.wasm)	Bytecode	n/a	digest/sig (opt)	WASM runtime + host	Host caps policy (external)	Separate	wasmtime/embed

Key deltas for .itpkg:
	•	Purpose-first: ships declarative intentions (what to do), not just code to be called.  ￼
	•	Policy-aware by design: embeds capabilities, network allow/deny, PII export, energy modes with deny-by-default expectations the runtime enforces.  ￼
	•	Portable tests: tests: blocks travel with the package and run the same in CLI, exec service, IDE viewer, or agents.  ￼
	•	AI-readable semantics: the content is ITML, which is both human- and LLM-parsable, enabling composition and reasoning beyond language/framework boundaries.

⸻

Why a new package at all?
	1.	Preserve semantics end-to-end
ITML encodes intent, inputs/outputs, workflow, rules, tests, and policies. Conventional package formats treat these as disparate docs or build-time concerns; .itpkg keeps them first-class and executable at install/run time.  ￼
	2.	Deterministic, policy-enforced execution
The runtime is deny-by-default and must know capabilities and network policy up front. .itpkg bundles those decisions with the artifact so resolution ≠ surprise side effects.  ￼
	3.	Cross-runtime portability
An ITML project should run in the headless executor, the viewer, a browser extension, or through MCP agents—without rewriting to a framework-specific app. A purpose-centric package lets the ecosystem (CLI/API/exec/web/IDE) agree on one contract.  ￼
	4.	Verifiable provenance
Signed artifacts + registry metadata (downloads/installs metrics, audit logs) make intent packages traceable and trustworthy—similar to modern supply-chain practices, but tailored to intent artifacts.  ￼
	5.	AI-native composition
Agents (via MCP/SDKs) can search, install, run, and compose .itpkg units as building blocks, because what’s inside is structured purpose, not opaque binaries.  ￼

⸻

Minimal itpkg.json (illustrative)

{
  "name": "@shop/checkout",
  "version": "0.3.1",
  "description": "Checkout flow (intents + views + policies)",
  "entry": "project.app.itml",
  "itmlVersion": "0.1",
  "capabilities": ["ui.render", "http.outbound"],
  "policies": {
    "security": { "network": { "outbound": { "allow": ["api.payments.com:443"], "deny": ["*"] } } },
    "privacy": { "pii": { "export": "deny" } },
    "energy": { "mode": "optimize-render" }
  }
}

This mirrors what the runtime needs before execution and what the registry shows on the package page (inputs/outputs, routes, policies, install snippet).

⸻

TL;DR
	•	.itpkg is the portable, signed container for ITML projects—intents + policies + tests—so they can be installed, verified, and executed by AI runtimes with predictable behavior.
	•	It’s similar to npm/Docker/Helm in distribution & versioning, but different in that it preserves purpose, contracts, and capability policies as first-class runtime inputs, not just docs or build metadata.
	•	We need it to make Software by Intentions practical: searchable, composable, secure, and AI-native across the entire ecosystem (CLI, API, Exec, Web, IDE, Agents).