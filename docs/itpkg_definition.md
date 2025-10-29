What is .itpkg?

A signed, versioned Intent package: a single distributable archive that bundles one ITML project (its .itml intents plus metadata, policies, tests, and optional assets) so any AI Runtime or the intent tool can install, verify, and execute it deterministically. Think “the npm tarball / Docker image of ITML,” but purpose-first and policy-aware.
	•	Why do we need it?
ITML apps aren’t just code—they’re executable intentions (purpose + contracts + policies). We need a package that preserves those semantics, capabilities, tests, and signatures end-to-end across registries and runtimes. Existing formats don’t model capability gates, privacy/energy policies, and declarative workflows as first-class, portable artifacts.

⸻

What's inside a .itpkg

A .itpkg is a **single tar.gz archive** with a flat structure (no nested archives). The root contains:

**Required files:**
/itpkg.json                # Package metadata (name, version, entry, policies)
/MANIFEST.sha256           # File list with SHA256 checksums
/SIGNATURE                 # ed25519 signature over MANIFEST.sha256 (or "UNSIGNED")

**Project structure:**
/project.app.itml          # Entrypoint (required for app packages)
/intents/**/*.itml         # Intents, layouts, views, components (required directory)
/policies/**/*.itml        # Security/privacy/energy policies (required directory)
/schemas/**/*.itml         # Data contracts (recommended)
/tests/**/*.itml           # Declarative test cases (recommended)
/.ci/*.itml                # Lint/format tasks (optional)
/assets/**                  # Optional static assets

**Archive Layout:**
```
.itpkg (tar.gz)
├── itpkg.json
├── MANIFEST.sha256
├── SIGNATURE
├── project.app.itml
├── intents/
│   └── ...
├── policies/
│   └── ...
└── [other project files...]
```

**Key points:**
	•	**Flat structure**: All files at archive root (no nested payload.tar.gz)
	•	**project.app.itml**: Entrypoint the runtime reads to resolve imports, routes, and policies
	•	**Policies & capabilities**: Encoded declaratively in itpkg.json (deny-by-default, allowlists, PII/energy modes)
	•	**Tests**: Portable, runtime-executable tests in /tests/**/*.itml
	•	**itpkg.json**: Required authoritative manifest with validation rules
	•	**MANIFEST.sha256**: Contains checksums for all files (except itself)
	•	**SIGNATURE**: ed25519 signature over MANIFEST.sha256 content

How it's produced & used (MVP flow):
	1.	intent package [path] → builds the .itpkg, signs it with ed25519, validates structure and policies.
	2.	intent publish @scope/pkg@version → uploads signed artifact to registry (validates itpkg.json, signature, policies).
	3.	intent install @scope/name → resolves, verifies ed25519 signature, validates checksums, extracts into project.
	4.	intent verify ./pkg.itpkg → verifies signature and integrity (supports --legacy-hmac for old artifacts).

**CLI Commands:**
```bash
# Package with scaffold (generates itpkg.json if missing)
intent package . --scaffold --unsigned

# Package with ed25519 signing
intent package . --sign-key ~/.ssh/intent_sign_key

# Package using environment variable
INTENT_SIGN_KEY=/path/to/key intent package .

# Verify package integrity
intent verify package.itpkg
intent verify package.itpkg --legacy-hmac  # For old HMAC-signed packages
```

⸻

Similarities & differences vs other package types

Format	What it bundles	Manifest	Integrity	Runtime model	Policies/Capabilities	Tests inside	Primary install UX
.itpkg	ITML intents + schemas + policies + tests	itpkg.json + project.app.itml	ed25519 signature + MANIFEST.sha256	AI Runtime reads ITML to execute workflows	First-class (security/privacy/energy, deny-by-default)	Yes, declarative tests:	intent install / viewer
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

itpkg.json Specification (v0.1)

**Required Fields:**
- `name`: Package identifier (e.g., `"@scope/pkg"`)
- `version`: Semantic version (e.g., `"0.3.1"`)
- `description`: Human-readable description
- `itmlVersion`: ITML format version (e.g., `"0.1"`)
- `policies`: Policy object (see below)

**App Packages (type: "app" or default):**
- `entry`: Required, path to entrypoint file (typically `"project.app.itml"`)
- `policies.security.network`: Required, must define outbound network policy

**Library Packages (type: "lib"):**
- `entry`: Optional, omitted for library packages
- `policies.security.network`: Still required, typically `{ "outbound": { "deny": ["*"] } }`

**Optional Fields:**
- `type`: `"app"` (default) or `"lib"`
- `capabilities`: Array of capability strings (e.g., `["ui.render", "http.outbound"]`)
- `meta.signature`: Signature metadata (algorithm, keyId)

**Example App Package:**
```json
{
  "name": "@shop/checkout",
  "version": "0.3.1",
  "description": "Checkout flow (intents + views + policies)",
  "entry": "project.app.itml",
  "itmlVersion": "0.1",
  "capabilities": ["ui.render", "http.outbound"],
  "policies": {
    "security": {
      "network": {
        "outbound": {
          "allow": ["api.payments.com:443"],
          "deny": ["*"]
        }
      }
    },
    "privacy": {
      "pii": {
        "export": "deny"
      }
    },
    "energy": {
      "mode": "optimize-render"
    }
  }
}
```

**Example Library Package:**
```json
{
  "name": "@acme/ui-kit",
  "version": "0.4.0",
  "description": "Common intents & views",
  "type": "lib",
  "itmlVersion": "0.1",
  "capabilities": [],
  "policies": {
    "security": {
      "network": {
        "outbound": {
          "deny": ["*"]
        }
      }
    }
  }
}
```

**Validation Rules:**
- `name`, `version`, `itmlVersion` are required
- App packages must have `entry` field
- App packages must have `policies.security.network` defined
- All packages must have `policies` object
- Entry file must exist if specified

⸻

MANIFEST.sha256 Format

Contains SHA256 checksums for all files in the package (except MANIFEST.sha256 itself):

```
<sha256-hex>  <relative/path>
<sha256-hex>  <relative/path>
...
```

Example:
```
f2c7a8b9c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8  itpkg.json
a91e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9  project.app.itml
0d77e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6  intents/home.view.itml
...
```

Entries are sorted by path for deterministic output. The signature covers this entire manifest content.

⸻

SIGNATURE Format

**Signed packages:**
- Contains ed25519 signature bytes (64 bytes) over MANIFEST.sha256 content
- Signature is hex-encoded when stored as file
- Public key can be stored in `itpkg.json.meta.signature.keyId` for discovery

**Unsigned packages:**
- Contains literal string `"UNSIGNED"` (allowed only with `--unsigned` flag)
- Not recommended for production use

⸻

Directory Structure Requirements

**Required (ERROR if missing):**
- `/intents/` - Must exist as directory (may be empty in edge cases)
- `/policies/` - Must exist as directory (can contain minimal stubs)
- `/project.app.itml` - Required for app packages (must match `itpkg.json.entry`)

**Recommended (WARNING if missing for app packages):**
- `/schemas/` - Data contract definitions
- `/tests/` - Test cases (warn if missing for app packages)
- `/.ci/` - CI/lint tasks
- `/assets/` - Static assets

**Validation Output:**
- ERROR: Missing required directories/files → package fails
- WARN: Missing recommended items → package succeeds with warning
- INFO: Extra files outside canonical tree → logged but allowed

⸻

CLI Behaviors & Error Codes

**intent package [path]**
- Fails if `itpkg.json` missing (unless `--scaffold`)
- Fails if entry missing for type !== "lib"
- Fails if required directories (`intents/`, `policies/`) missing
- Validates policies.security.network for app packages
- Outputs: `{name}-{version}.itpkg`

**intent publish @scope/pkg@1.2.3**
- Requires valid `itpkg.json`
- Validates signature (if registry requires signed packages)
- Validates policy structure
- Checks compatibility (itmlVersion)

**intent verify ./pkg.itpkg**
- Verifies ed25519 signature over MANIFEST.sha256
- Validates all file checksums
- Supports `--legacy-hmac` for old HMAC-signed artifacts (with deprecation warning)

**Exit Codes:**
- `0` - Success
- `1` - Validation error (manifest/structure)
- `2` - Signature/Integrity error
- `3` - Policy violation (deny-by-default failures)
- `4` - Compatibility error (itmlVersion mismatch)

⸻

Migration Notes (v0.1)

**Changes from initial implementation:**
- ❌ Removed: Nested `payload.tar.gz` structure
- ❌ Removed: `sha256.txt` and `signature.txt` files
- ❌ Removed: HMAC-SHA256 signing (replaced with ed25519)
- ❌ Removed: tar.gz format option (always .itpkg now)
- ✅ Added: Flat tar.gz structure with files at root
- ✅ Added: `MANIFEST.sha256` with sorted file list
- ✅ Added: `SIGNATURE` file with ed25519 signature
- ✅ Added: Required `itpkg.json` manifest
- ✅ Added: Directory structure validation
- ✅ Added: `--scaffold` flag for manifest generation

**Backward Compatibility:**
- Old HMAC-signed packages can be verified with `--legacy-hmac` flag
- Support for legacy format will be maintained for 2 minor versions (e.g., 0.1 → 0.3)
- New packages always use ed25519 signing

⸻

TL;DR
	•	.itpkg is the portable, signed container for ITML projects—intents + policies + tests—so they can be installed, verified, and executed by AI runtimes with predictable behavior.
	•	It’s similar to npm/Docker/Helm in distribution & versioning, but different in that it preserves purpose, contracts, and capability policies as first-class runtime inputs, not just docs or build metadata.
	•	We need it to make Software by Intentions practical: searchable, composable, secure, and AI-native across the entire ecosystem (CLI, API, Exec, Web, IDE, Agents).