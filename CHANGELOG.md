# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Phase 8.2: First integration test with CSV reader bento example
- Integration test framework with helper utilities
- 100% CLI command test coverage (15 tests)
- Comprehensive DEVELOPMENT.md guide for Go development
- Version command shorthand (`bento v`)
- OS-specific installation instructions for glow
- Test fixtures and mock servers (Figma API, Blender)
- Node.js to Go translation guide
- Dependency management philosophy documentation

### Changed
- Improved error messages for missing glow installation
- Renamed `helpers.go` to `test_utilities.go` for Bento Box compliance
- Enhanced package documentation for integration tests

### Fixed
- Linting errors in progress tests
- Error handling in test utilities

## [0.1.0] - 2025-10-19

### Added

#### Phase 7: CLI Implementation
- Complete bento CLI with playful sushi-themed commands
- `savor` command - Execute bento workflows
- `sample` command - Validate bentos without execution
- `menu` command - List available bentos
- `box` command - Create new bento templates
- `recipe` command - View documentation with glow
- `version` command - Display version information
- Charm CLI integration with beautiful output
- Progress display with Bubbletea TUI

#### Phase 7: Visual Feedback (Miso Package)
- Bubbletea-based progress display
- Daemon-combo pattern for running Bubbletea in background
- Multi-step progress tracking
- Spinner and progress bar components
- Status update system

#### Phase 6: Orchestration Engine (Itamae Package)
- Bento orchestration engine
- Group execution with dependency resolution
- Context passing between neta
- Edge-based execution flow
- Error handling and validation
- Template variable substitution

#### Phase 5: Neta Registry (Pantry Package)
- Neta type registry system
- Factory pattern for neta instantiation
- Built-in neta types registration
- Extensible architecture for custom neta

#### Phases 2-4: Core Infrastructure
- **Shoyu** (pkg/shoyu): Logging with charm/log integration
- **Omakase** (pkg/omakase): Configuration management
- **Hangiri** (pkg/hangiri): Core models and definitions

#### Neta Implementations
- **spreadsheet**: Read/write CSV and Excel files
- **http-request**: HTTP client with full request support
- **transform**: Data transformation with templates
- **edit-fields**: Set/modify field values
- **file-system**: File and directory operations
- **shell-command**: Execute shell commands
- **image**: Image processing and optimization
- **loop**: Iterate over collections
- **parallel**: Concurrent execution
- **group**: Orchestrate multiple neta

### Documentation
- Bento Box Principle (coding philosophy)
- Package naming conventions
- Go standards review
- Complete node inventory
- Status word guidelines
- Emoji usage standards
- Charm stack integration guide
- Phase strategy documents (Phases 1-8)

### Infrastructure
- Go module setup
- golangci-lint configuration
- Git workflow and hooks
- Project structure following Go best practices

## Version Numbering

Bento follows [Semantic Versioning](https://semver.org/):

- **MAJOR** version (1.0.0): Incompatible API changes
- **MINOR** version (0.1.0): New functionality in a backward compatible manner
- **PATCH** version (0.0.1): Backward compatible bug fixes

### Pre-1.0.0 Versioning

While in 0.x.x versions (pre-1.0):
- **MINOR** version increments may include breaking changes
- **PATCH** version increments are for backward compatible changes
- Version 1.0.0 will be released when the API is stable

### Current Status

**Current Version:** 0.1.0 (Development)

We're currently in **Phase 8** of development, building integration tests and validating the complete system with real-world workflows.

### Upcoming Versions

- **0.2.0**: Phase 8 completion (real-world integration tests)
- **0.3.0**: Additional neta types and features
- **0.4.0**: Performance optimizations
- **0.5.0**: Enhanced error handling and validation
- **1.0.0**: Stable API, production-ready release

## How to Update This Changelog

When completing a phase or adding features:

1. **Add to [Unreleased] section** during development
2. **Create a new version section** when releasing
3. **Use these categories:**
   - `Added` - New features
   - `Changed` - Changes to existing functionality
   - `Deprecated` - Soon-to-be removed features
   - `Removed` - Removed features
   - `Fixed` - Bug fixes
   - `Security` - Security fixes

4. **Follow this format:**
   ```markdown
   ### Added
   - Feature description with context
   - Another feature with details
   ```

## Git Tags and Releases

To create a new release:

```bash
# Update CHANGELOG.md with version and date
# Commit the changelog
git add CHANGELOG.md
git commit -m "chore: Release v0.2.0"

# Create and push tag
git tag -a v0.2.0 -m "Release v0.2.0"
git push origin v0.2.0
git push origin main

# GitHub will automatically create a release from the tag
```

## Links

- [Keep a Changelog](https://keepachangelog.com/)
- [Semantic Versioning](https://semver.org/)
- [Go Modules](https://go.dev/blog/using-go-modules)

---

**Note:** This changelog started with version 0.1.0. Earlier development history is preserved in git commits but not detailed here for brevity.
