# LocalNet

- `anvil --block-time 0.05`
- From inside `scripts`: `forge script DeployHyperlane.s.sol --rpc-url http://localhost:8545 --private-key 0x2a871d0798f97d79848a013d4936a73bf4cc922c825d33c1cf7073dff6d409c6 --broadcast`
- From inside `indexer`: `pnpm dev`
- From inside `scripts`: `go run get_metadata.go`
