# vibe-c2-golang-channel-core

Golang foundation SDK for building Vibe C2 channel modules.

## v0 Architecture Summary

The repository now includes a compileable v0 scaffold with the primary package boundaries:

- `pkg/runtime`: channel runtime orchestrator and core interfaces
- `pkg/profile`: profile model, YAML parsing helpers, semantic validation
- `pkg/matcher`: hint-first profile selection stub
- `pkg/syncclient`: HTTP sync client for C2 (`POST /api/channel/sync`)
- `pkg/mgmtrpc`: management RPC server skeleton (CRUD/activate/validate TODOs)
- `pkg/errors`: typed error codes and helpers

Core runtime interfaces:

- `TransportEnvelope`: `GetField`, `SetField`, `SourceKey`
- `SyncClient`: canonical sync contract
- `ProfileStore`: list/get/put/delete profile persistence contract

Runtime pipeline entrypoint:

- `Handle(ctx, envelope, channelID) (protocol.OutboundAgentMessage, error)`
- Performs basic canonical inbound/outbound validation around sync execution.

## Quick Start

### 1) Test the scaffold

```bash
go test ./...
```

### 2) Wire your channel adapter envelope

Implement `runtime.TransportEnvelope` in your transport layer and populate:

- `mapping.id`
- `mapping.encrypted_data`
- optional `mapping.profile_id`

Then call:

```go
rt := runtime.New(mySyncClient)
out, err := rt.Handle(ctx, envelope, channelID)
```

### 3) Parse and validate profiles

```go
p, err := profile.ParseYAML(data)
if err != nil {
    // invalid YAML or semantic profile error
}
```

## Notes

- This v0 scaffold includes a local module replacement for `github.com/vibe-c2/vibe-c2-golang-protocol` under `third_party/vibe-c2-golang-protocol` so it builds in restricted/offline environments.
- See `docs/FOUNDATION-PLAN.md` for broader roadmap details.
