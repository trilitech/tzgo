# Feature Specification: tz5 Address Support (PKH only)

**Feature Branch**: `001-tz5-pkh-support`

**Created**: 2026-06-10

**Status**: Draft

**Input**: User description: "Address item 2 from BIN-878 (v025/Ushuaia cryptography). For now only do the fast win: tz5 PKH support only. Majority of applications will be able to get away with tz5 as PKH only — tx parsing, tx recipients etc."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Parse tz5 addresses from chain data (Priority: P1)

An application reading Tezos blocks, operations, or contract storage encounters a
tz5 address (a post-quantum ML-DSA-44 account introduced in protocol v025) as a
sender, recipient, delegate, or stored value. The SDK recognises the address in
both its text form (`tz5…`) and its on-chain binary form, and surfaces it as a
normal typed address — exactly as it does today for tz1/tz2/tz3/tz4.

**Why this priority**: Without this, reading any block or operation that touches a
tz5 account fails or silently mis-decodes. This is the breaking exposure named in
BIN-878 and blocks indexers and wallets the moment the feature flag activates on
any network.

**Independent Test**: Feed the SDK a tz5 address string and the equivalent binary
encodings; verify both decode to the same typed address and re-encode to the exact
original bytes/string.

**Acceptance Scenarios**:

1. **Given** a valid `tz5…` address string, **When** the application parses it,
   **Then** it gets a valid address whose type identifies the new account kind and
   whose text form round-trips to the identical string.
2. **Given** a 21-byte binary implicit address with the new on-chain tag (4),
   **When** the application decodes it, **Then** it gets the same address as the
   text form, and re-encoding produces the identical 21 bytes.
3. **Given** a 22-byte padded binary address (as found in contract storage),
   **When** decoded, **Then** the same address results.
4. **Given** an invalid tz5 string (wrong checksum or length), **When** parsed,
   **Then** the application receives an explicit error, never a wrong address.

---

### User Story 2 - Use tz5 addresses as transaction recipients (Priority: P2)

A wallet or service constructs a transaction whose destination is a tz5 address.
The SDK encodes the destination correctly so a Tezos node accepts the operation.

**Why this priority**: Sending funds *to* tz5 accounts is the most common
interaction ordinary applications will have with post-quantum accounts; it does
not require any new cryptography on the sender side.

**Independent Test**: Build a transaction with a tz5 destination and verify the
binary encoding matches the protocol's expected destination bytes.

**Acceptance Scenarios**:

1. **Given** a transaction whose destination is a tz5 address, **When** the
   operation is serialized, **Then** the destination field contains the correct
   tagged binary form and the operation round-trips decode→encode unchanged.
2. **Given** a tz5 address used where an entity classification is needed,
   **When** the application asks "is this an implicit (externally-owned)
   account?", **Then** the answer is yes.

---

### Edge Cases

- A blinded (btz1) address previously shared the internal tag slot now assigned
  on-chain to tz5. Binary encoding of blinded addresses must not be confusable
  with tz5; blinded addresses never legitimately appear in tagged binary form, so
  they must now fail loudly rather than alias to tz5.
- ML-DSA public keys, secret keys, and signatures are **out of scope**; if chain
  data contains an ML-DSA public key or signature (e.g. a tz5 reveal), the SDK
  must return an explicit "unsupported" error rather than mis-parse — matching
  current behaviour for unknown key tags.
- Addresses with entrypoint suffixes (`tz5…%ep`) must parse like other types.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The SDK MUST parse base58 `tz5…` strings (prefix bytes 6,161,169;
  20-byte hash; 36-char encoding) into a typed address, and format that address
  back to the identical string.
- **FR-002**: The SDK MUST decode 21-byte tagged binary implicit addresses with
  tag 4 as tz5 addresses, and encode tz5 addresses to that exact binary form;
  the 22-byte zero-padded variant MUST behave equivalently.
- **FR-003**: tz5 addresses MUST be classified as implicit (externally-owned)
  accounts, comparable and usable as map keys like all other address types.
- **FR-004**: Binary encoding of blinded (btz1) addresses MUST NOT produce the
  tz5 tag; the previous internal aliasing of tag 4 to blinded MUST be removed.
- **FR-005**: Operations (e.g. transactions) that reference tz5 addresses as
  source/destination MUST round-trip through binary encode/decode without loss.
- **FR-006**: ML-DSA keys and signatures remain unsupported in this slice; any
  attempt to decode them MUST fail with an explicit error (no silent mis-parse).
- **FR-007**: Existing behaviour for tz1/tz2/tz3/tz4/KT1/txr1/sr1/btz1 text
  parsing MUST be unchanged (backward compatibility).

### Assumptions

- The tz5 base58 prefix bytes and the binary tag value are taken from the
  protocol implementation source (octez `base58.ml`: `"\006\161\169"`;
  `signature_v3.ml`: PKH union Tag 4) — both verified, not guessed.
- Full ML-DSA-44 support (keys, signatures, reveal proofs, signing) is a
  follow-up slice of BIN-878 item 2 and intentionally excluded here.
- tz5 accounts are feature-flagged off on mainnet at v025 activation, so this
  slice carries no immediate mainnet risk while unblocking testnet usage.

## Success Criteria *(mandatory)*

- **SC-001**: 100% of valid tz5 address strings round-trip parse→format
  byte-identically, and invalid ones produce explicit errors.
- **SC-002**: 100% of tagged binary tz5 encodings (21- and 22-byte forms)
  round-trip decode→encode byte-identically.
- **SC-003**: Applications can classify tz5 addresses as implicit accounts and
  use them as transaction recipients with no application-side special-casing.
- **SC-004**: The full existing test suite continues to pass (no regression for
  other address types).

## Key Entities

- **tz5 address**: a 20-byte public key hash of an ML-DSA-44 (post-quantum)
  public key, text-encoded with prefix `tz5`, binary-encoded as implicit-address
  tag 4. Introduced by protocol v025 (Ushuaia) behind the `tz5_account_enable`
  feature flag.
- **Blinded address (btz1)**: commitment-style address that previously occupied
  the SDK-internal tag slot now used on-chain by tz5; never legitimately appears
  in tagged binary form.
