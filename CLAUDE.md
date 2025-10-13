# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`codegen` is a Go-based code generation tool that generates code from Protocol Buffer (protobuf) definitions. It parses `.proto` files and generates Go code with enhanced struct tags (gorm, validate, role), implementation scaffolding, database tables, cache layers, i18n support, and more.

**Documentation**: https://lazygophers.pages.dev/
**Template Repository**: https://github.com/lazygophers/template

## Architecture

### Core Components

1. **codegen/** - Core code generation logic
   - `pbparse.go` - Protobuf file parsing using github.com/emicklei/proto
   - `generate_pb_file.go` - Generates .pb.go files via protoc
   - `generate_struct_tag.go` - Adds struct tags (gorm, validate, role) to generated structs
   - `generate_impl.go` - Generates service implementation boilerplate
   - `generate_table.go`, `generate_cache.go`, `generate_conf.go` - State layer generators
   - `generate_add_rpc.go` - Batch RPC endpoint generation

2. **cli/** - CLI interface using cobra
   - `root.go` - Root command initialization
   - `gen/gen.go` - Main generation command group with subcommands (pb, impl, table, cache, etc.)
   - `i18n/` - i18n translation commands

3. **i18n/** - Internationalization framework
   - `translater.go` - Translation providers (Google Free, DeepL Free)
   - `tran.go` - Translation orchestration with caching
   - `storage.go` - BoltDB-based translation cache
   - `localize.go` - Multi-format serialization (JSON, YAML, TOML)

4. **state/** - Configuration and state management
   - `config.go` - Main configuration structure
   - `load.go` - Configuration loading from multiple sources
   - `lazy_config.go` - Project-specific .lazygophers config

### Data Flow

1. Parse `.proto` file → `PbPackage` structure
2. Generate `.pb.go` using protoc
3. Parse generated `.pb.go` to add custom tags
4. Generate implementation files (impl, routes, clients)
5. Generate state files (tables, cache, config) if enabled

## Common Commands

### Build
```bash
# Build using goreleaser
make build

# Output: ./dist/cli_<platform>_<arch>/codegen
```

### Run Tests
```bash
go test ./...
```

### Linting
```bash
make lint
# Requires: golangci-lint
```

### Generate i18n Translations
```bash
# Generate translations from source file
make tran

# Or directly:
codegen i18n tran --generate-const=true -s ./state/localize/zh.yaml
```

### Generate Code from Proto
```bash
# Generate .pb.go file
codegen gen pb -i path/to/file.proto

# Generate all (pb + impl + state)
codegen gen all -i path/to/file.proto

# Generate specific components
codegen gen impl -i path/to/file.proto
codegen gen table -i path/to/file.proto
codegen gen cache -i path/to/file.proto
```

### Configuration

Configuration is loaded from (in order):
1. Environment variable `LAZYGO_CODEGEN_CONFIG_FILE`
2. System config directory: `~/Library/Application Support/lazygophers/codegen.cfg.yaml` (macOS)
3. Project-level `.lazygophers` file
4. Command-line flags

**Key configuration fields**:
- `protoc-path` - Path to protoc binary
- `protogen-go-path` - Path to protoc-gen-go plugin
- `output-path` - Base directory for generated code
- `proto-files` - Additional proto import paths
- `default-tag` - Default struct tags (gorm, validate)
- `state.table/cache/i18n` - Enable/disable state generation

See `example.codegen.cfg.yaml` for full configuration options.

## Development Patterns

### Proto Comment Tags

The tool recognizes special comment tags in .proto files:

```protobuf
message User {
  // @gorm: primary_key
  // @role: admin:wr;user:r
  // @validate: required
  string id = 1;
}
```

These are parsed and converted to Go struct tags in the generated code.

### Code Generation Flow

When adding new generators:
1. Create `generate_<name>.go` in `codegen/`
2. Add CLI command in `cli/gen/gen_<name>.go`
3. Register command in `cli/gen/gen.go:initGen()`
4. Add localized strings in `state/localize/`

### Testing Generation

```bash
# Build the tool
make build

# Test on example proto
./dist/cli_darwin_arm64/codegen gen pb -i ../example/example.proto
```

## Windows Compatibility

The codebase has Windows path handling fixes (see commit 329f560, 451b9c0). When working with file paths:
- Always use `filepath.ToSlash()` for cross-platform paths
- Use `filepath.Join()` instead of manual path construction
- Test on Windows when modifying path handling

## i18n System

The i18n translation system supports:
- Multiple translation providers (Google Free, DeepL Free)
- Translation caching with BoltDB to avoid re-translating unchanged text
- Automatic translation of nested YAML/JSON/TOML structures
- Batch translation with configurable overwrite policies

When translating:
- Source file: typically `state/localize/zh.yaml`
- Cache file: `.i18n_tran_cache.db` (configurable)
- Supports 100+ languages

## Common Pitfalls

1. **Proto file location**: The tool auto-detects proto files by walking up the directory tree. Use `-i` flag if auto-detection fails.

2. **Module prefix**: If `go_package` in .proto doesn't match your actual module path, use `--go-module-prefix` flag or configure `go-module-prefix` in config.

3. **Overwrite protection**: By default, impl files won't be overwritten. Use `--overwrite` flag with caution.

4. **Lazy config**: The tool looks for `.lazygophers` file in the project root (determined by go_package). This configures project-specific generation options.