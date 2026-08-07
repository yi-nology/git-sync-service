# git-platform-sdk Upgrade Summary

## Overview

Successfully upgraded `git-platform-sdk` from v0.34.0 to v0.35.0 with all required breaking changes addressed.

## Changes Made

### 1. SDK Version Upgrade (Task 1)
- Updated `go.mod` to use `github.com/yi-nology/git-platform-sdk v0.35.0`
- Removed local replace directive

### 2. CryptoManager Migration (Task 2)
- Replaced all `credential.EncryptGCM()` / `credential.DecryptGCM()` calls with `CryptoManager`
- `CryptoManager` is initialized via `credential.NewCryptoManager()` (reads `ENCRYPTION_KEY` env var)
- Affected files:
  - `internal/dao/repo_dao.go` - DAO layer encryption/decryption

### 3. Platform Detection Error Handling (Task 3)
- Updated error handling to use `provider.ErrPlatformNotSupported` for unsupported platforms
- Affected files:
  - `internal/service/repo.go` - Service layer platform detection

### 4. Environment Variable Configuration (Task 4)
- Added documentation for required `ENCRYPTION_KEY` environment variable
- `ENCRYPTION_KEY` must be 32 bytes (64 hex characters) for AES-256 encryption
- Created `docs/ENVIRONMENT_VARIABLES.md`

### 5. Test Updates (Task 5-6)
- Added comprehensive tests for `CryptoManager` (encrypt/decrypt roundtrip, empty string, missing key)
- Added integration tests for end-to-end encryption through DAO layer
- Added tests for `ErrPlatformNotSupported` handling
- Fixed lint issues: replaced `os.Setenv`/`os.Unsetenv` with `t.Setenv` in tests

### 6. Final Verification (Task 7)
- Ran `go vet ./...` - passed
- Ran `go mod tidy` - passed
- Ran `golangci-lint run` - 0 issues
- Ran `go test ./...` - all tests pass

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `ENCRYPTION_KEY` | Yes | 32-byte hex key for AES-256 encryption of access tokens |

## Breaking Changes Addressed

1. **CryptoManager API**: `credential.EncryptGCM`/`DecryptGCM` replaced with `CryptoManager.Encrypt`/`Decrypt`
2. **Platform Detection**: `provider.DetectPlatform` now returns `ErrPlatformNotSupported` for unsupported hosts

## Verification

All quality checks pass:
- `go vet ./...` - clean
- `golangci-lint run` - 0 issues
- `go test ./...` - all tests pass
