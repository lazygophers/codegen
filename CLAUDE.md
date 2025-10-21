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

## Development Guidelines

### Core Development Principles

- **Memory Management**: When users say "remember xxx" or repeatedly emphasize certain rules, behaviors, etc., automatically update relevant memory banks and documentation
- **Code Analysis**: Before coding, analyze at least 3 existing implementations or patterns to identify reusable interfaces and constraints
- **Dependency Mapping**: Draw dependencies and integration points, confirming input/output protocols, configuration, and environment requirements
- **Code Consistency**: Understand existing testing frameworks, naming conventions, and formatting rules to ensure output aligns with the codebase
- **Documentation Priority**: Use context7 to query programming library documentation first, avoid over-reliance on web searches or assumptions
- **Best Practice Learning**: Use github.search_code to search open-source implementation examples and learn best practices
- **Complete Implementation**: Absolutely prohibit MVP, minimal implementations, or placeholders; full functionality and data paths must be completed before submission
- **Code Cleanup**: Must improve all MVP, minimal implementations, and placeholders to complete, concrete code implementations
- **Active Cleanup**: Must proactively delete outdated, duplicate, or escape-style code to maintain clean implementations
- **Style Consistency**: Must always adhere to programming language standard code styles and project-specific style specifications
- **Breaking Changes**: Do not maintain backward compatibility for breaking changes, while providing migration steps or rollback plans
- **Disruptive Strategy**: Must always adopt disruptive breaking change strategy, absolutely no backward compatibility
- **Quality Assurance**: Must follow best practices to ensure code quality and maintainability
- **Language-Specific Patterns**: Absolutely prohibit using Java-style coding standards to write code or documentation in other languages
- **Performance Evaluation**: Must evaluate time complexity, memory usage, and I/O impact during design to avoid unnecessary consumption
- **Optimization Awareness**: After identifying potential bottlenecks, provide monitoring or optimization suggestions to ensure sustainable iteration
- **Dependency Assessment**: Prohibit introducing unevaluated expensive dependencies or blocking operations

## 🔧 Project Integration Rules

### Learning the Codebase
- Must find at least 3 similar features or components to understand their design and reuse patterns
- Must identify common patterns and conventions in the project and apply them in new implementations
- Must prioritize using existing libraries, tools, or helper functions
- Must follow existing testing arrangements and沿用 assertion and fixture structures

### Tools
- Must use the project's existing build system, no private addition of scripts
- Must use the project's established testing framework and execution methods
- Must use the project's formatting/static checking settings
- If there is indeed a need for new tools, must provide sufficient justification and obtain documented approval

## ⚠️ Important Reminders

**Absolutely Prohibited:**
- Making assumptions without evidence; all conclusions must cite existing code or documentation

**Must Be Done:**
- Complete detailed planning and documentation before implementing complex tasks
- Generate task decomposition for cross-module or more than 5 subtasks work
- Maintain TODO lists for complex tasks and update progress in a timely manner
- Verify planning documents are confirmed before starting development
- Maintain small-step delivery, ensuring each commit is in a usable state
- Synchronously update planning documents and progress records during execution
- Proactively learn the advantages and disadvantages of existing implementations and reuse or improve them
- Must pause operations after three consecutive failures and reassess strategy

## 🎯 Content Uniqueness Rules

- Each level must be self-consistent in mastering its own abstract scope, prohibiting cross-level content mixing
- Must reference materials from other levels rather than copy-pasting, maintaining unique information sources
- Each level must describe the system from the corresponding perspective, avoiding overreaching details
- Prohibit stacking implementation details in high-level documents, ensuring clear architecture and implementation boundaries