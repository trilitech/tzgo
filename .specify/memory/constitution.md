<!--
SYNC IMPACT REPORT
==================
Version change: (template, unversioned) → 1.0.0
Bump rationale: Initial ratification — placeholder template replaced with concrete principles.

Modified principles: (none — first concrete definition)
Added sections:
  - Core Principles I–V (Correctness & Protocol Compliance; Backward Compatibility &
    Stable Interfaces; Test-Verified Behavior; Performance & Efficiency; Protocol
    Coverage & Isolation)
  - Library & Dependency Constraints
  - Development Workflow & Quality Gates
  - Governance
Removed sections: (none)

Templates requiring updates:
  ✅ .specify/templates/plan-template.md — generic "Constitution Check" gate; no hardcoded
     principles, references this file dynamically. No edit required.
  ✅ .specify/templates/spec-template.md — no constitution-coupled sections. No edit required.
  ✅ .specify/templates/tasks-template.md — tests are opt-in per template; consistent with
     Principle III (test discipline applies to public-API/encoding work). No edit required.
  ⚠ .specify/templates/commands/ — directory not present; nothing to propagate.

Follow-up TODOs: (none)
-->

# TzGo Constitution

## Core Principles

### I. Correctness & Protocol Compliance (NON-NEGOTIABLE)

TzGo MUST behave identically to Tezos mainnet for every supported protocol. Binary and
JSON encodings MUST round-trip without loss: any value decoded from the chain re-encodes
to the exact same bytes, and any value encoded by TzGo is accepted by a Tezos node.
Hashes, addresses, keys, signatures, Micheline data, and operation formats MUST match the
protocol specification and real on-chain data. When the SDK cannot represent something
correctly, it MUST return an explicit error rather than produce an approximate or silently
wrong result.

Rationale: TzGo is a low-level SDK other applications sign and broadcast transactions with;
a single encoding defect can lose funds or corrupt chain reads. Correctness is the product.

### II. Backward Compatibility & Stable Interfaces

The public Go API and on-the-wire encodings are a contract. Breaking changes to exported
identifiers or to serialization behavior MUST trigger a MAJOR version bump and be recorded
in the changelog. Within a major version, existing call sites MUST keep compiling and
behaving as before. Historic block data MUST remain readable across all supported
protocols even as newer protocols are added.

Rationale: The README promises stable interfaces and long-term support; downstream commercial
and non-commercial users depend on upgrades being non-disruptive.

### III. Test-Verified Behavior (NON-NEGOTIABLE)

Every change to encoding/decoding, type parsing, operation construction, or signing MUST be
covered by tests before it is considered complete. Tests MUST include round-trip cases and,
where available, real on-chain test vectors. Bug fixes MUST add a regression test that fails
before the fix and passes after. `go test ./...` MUST pass on every change.

Rationale: Correctness (Principle I) is only credible when it is continuously and
automatically verified against concrete vectors, not asserted by inspection.

### IV. Performance & Efficiency

TzGo targets high-performance applications. Hot paths (base58, hashing, Micheline encode/
decode, hash maps) MUST avoid unnecessary allocations and copies. Performance-sensitive code
MUST be supported by benchmarks, and changes that regress a benchmark MUST be justified by a
corresponding correctness or compatibility requirement. Optimizations MUST NOT compromise
Principle I.

Rationale: The SDK is positioned for reliable, high-throughput chain access; efficiency is a
stated design goal, but never at the expense of correctness.

### V. Protocol Coverage & Isolation

New Tezos protocols MUST be added as soon as practically feasible. Protocol-specific
behavior (constants, opcodes, encoding quirks) MUST be isolated and version-gated so that
adding a protocol does not alter behavior for existing ones. Reading support SHOULD span all
historic protocols; binary encoding and signing MAY be limited to recent protocols, and any
such limitation MUST be documented.

Rationale: Compliance is a moving target across protocol upgrades; isolation keeps each
protocol correct and makes the multi-protocol support contract (Principle II) maintainable.

## Library & Dependency Constraints

TzGo is a library-first SDK organized into focused, independently usable packages (`tezos`,
`micheline`, `rpc`, `codec`, `contract`, `signer`, `base58`, `hash`). Each package MUST have
a clear, single purpose and MUST be importable without pulling in unrelated concerns.
Third-party dependencies MUST be kept minimal and justified; prefer the standard library and
existing in-tree helpers over adding a dependency. The module path and Go version baseline
declared in `go.mod` are part of the compatibility contract and change only with explicit
review. CLI tools (`tzgen`, `tzcompose`) build on the libraries and MUST NOT cause the core
libraries to depend on CLI-only code.

## Development Workflow & Quality Gates

All changes land via pull request. Before merge, a change MUST: build cleanly, pass
`go vet ./...` and `go test ./...`, and include tests for new or changed behavior per
Principle III. Public-API or encoding changes MUST state their version impact (PATCH / MINOR
/ MAJOR) and update the changelog. Reviewers MUST verify compliance with the Core Principles,
especially round-trip correctness and backward compatibility. Complexity that violates a
principle MUST be justified in the PR description or rejected.

## Governance

This constitution supersedes other development practices for TzGo. Amendments MUST be made
via pull request that updates this document, states the rationale, and bumps the constitution
version per the policy below; dependent templates and guidance docs MUST be re-synced in the
same change.

Versioning policy (semantic):

- MAJOR: Backward-incompatible removal or redefinition of a principle or governance rule.
- MINOR: A new principle or section is added, or guidance is materially expanded.
- PATCH: Clarifications, wording, or non-semantic refinements.

Compliance review: every PR review MUST confirm the change upholds the Core Principles.
Runtime and agent guidance lives in `CLAUDE.md` and the active feature plan; those MUST stay
consistent with this constitution.

**Version**: 1.0.0 | **Ratified**: 2026-06-09 | **Last Amended**: 2026-06-09
