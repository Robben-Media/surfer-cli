# surfer-cli

A CLI tool for [Surfer SEO](https://surferseo.com/) built with Go. Create content editors, run SEO audits, and check content scores from the command line.

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap builtbyrobben/tap
brew install surfer-cli
```

### Download Binary

Download the latest release from [GitHub Releases](https://github.com/builtbyrobben/surfer-cli/releases).

### Build from Source

```bash
git clone https://github.com/builtbyrobben/surfer-cli.git
cd surfer-cli
make build
```

## Authentication

### Set API Key

```bash
# Interactive (secure, recommended)
surfer-cli auth set-key --stdin

# From environment variable
echo $SURFER_API_KEY | surfer-cli auth set-key --stdin

# From argument (discouraged - exposes in shell history)
surfer-cli auth set-key YOUR_API_KEY
```

### Check Status

```bash
surfer-cli auth status
```

### Remove Credentials

```bash
surfer-cli auth remove
```

### Environment Variables

- `SURFER_API_KEY` - Override stored credentials (also used by Surfer's official API)
- `SURFER_CLI_KEYRING_BACKEND` - Force keyring backend (auto/keychain/file)
- `SURFER_CLI_KEYRING_PASS` - Password for file backend (headless systems)

## Usage

### Content Editors

```bash
# List all content editors
surfer-cli editors list

# Create a new content editor
surfer-cli editors create --keywords "surfer seo guide,content optimization"

# Create with options
surfer-cli editors create --keywords "keyword1,keyword2" --language en --location "United States"

# Get editor details
surfer-cli editors get <editor-id>

# Get content score
surfer-cli editors score <editor-id>
```

### Audits

```bash
# List all audits
surfer-cli audits list

# Create a new audit
surfer-cli audits create --url https://example.com/page
```

### Output Formats

```bash
# JSON output (for scripting)
surfer-cli editors list --json

# Plain text output (TSV)
surfer-cli editors list --plain
```

## API Reference

This CLI wraps the [Surfer SEO API](https://surferseo.com/api/).

- **Base URL:** `https://app.surferseo.com/api/v1`
- **Auth:** `API-KEY` header

## Development

### Prerequisites

- Go 1.22+
- Make

### Commands

```bash
make build        # Build binary
make test         # Run tests
make lint         # Run linter
make ci           # Run full CI suite
make tools        # Install dev tools
```

## License

MIT

## Contributing

Contributions are welcome! Please read our contributing guidelines before submitting PRs.
