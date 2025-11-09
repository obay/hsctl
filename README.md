# HSCTL

<div align="center">

![Build Status](https://img.shields.io/github/actions/workflow/status/obay/hsctl/test.yml?branch=main&label=build&logo=github&logoColor=white)
![Release](https://img.shields.io/github/v/release/obay/hsctl?label=release&logo=github&logoColor=white&sort=semver)
![Downloads](https://img.shields.io/github/downloads/obay/hsctl/total?label=downloads&logo=github&logoColor=white)
![Go Version](https://img.shields.io/github/go-mod/go-version/obay/hsctl?logo=go&logoColor=white)
![License](https://img.shields.io/github/license/obay/hsctl?label=license&logo=github&logoColor=white)
![Platform](https://img.shields.io/badge/platform-linux%20%7C%20macos%20%7C%20windows-lightgrey?logo=github&logoColor=white)
![Go Report Card](https://goreportcard.com/badge/github.com/obay/hsctl)
![Code Size](https://img.shields.io/github/languages/code-size/obay/hsctl?logo=github&logoColor=white)
![Contributors](https://img.shields.io/github/contributors/obay/hsctl?logo=github&logoColor=white)
![Last Commit](https://img.shields.io/github/last-commit/obay/hsctl/main?logo=github&logoColor=white)
![GitHub stars](https://img.shields.io/github/stars/obay/hsctl?style=social&logo=github)
![GitHub forks](https://img.shields.io/github/forks/obay/hsctl?style=social&logo=github)

[![Made with Go](https://img.shields.io/badge/Made%20with-Go-00ADD8?style=flat&logo=go&logoColor=white)](https://golang.org/)
[![Powered by HubSpot](https://img.shields.io/badge/Powered%20by-HubSpot-FF7A59?style=flat&logo=hubspot&logoColor=white)](https://www.hubspot.com/)
[![CLI Tool](https://img.shields.io/badge/CLI-Tool-4EC820?style=flat&logo=gnu-bash&logoColor=white)](https://github.com/obay/hsctl)

</div>

**HSCTL** is a powerful command-line interface (CLI) tool written in Go for managing HubSpot contacts. It provides a simple, efficient way to perform CRUD operations on your HubSpot contacts directly from your terminal, eliminating the need to navigate through the HubSpot web interface for routine contact management tasks.

[Features](#benefits) • [Installation](#installation) • [Usage](#usage) • [Documentation](#command-reference) • [Contributing](#contributing)

## Why HSCTL Exists

Managing contacts in HubSpot through the web interface can be time-consuming, especially when you need to:
- Bulk update contact properties
- Quickly search and filter contacts
- Automate contact management workflows
- Integrate HubSpot operations into scripts and automation pipelines

HSCTL bridges this gap by providing a fast, scriptable interface to HubSpot's contact management capabilities, making it perfect for developers, sales operations teams, and anyone who prefers working from the command line.

## Benefits

- ⚡ **Fast**: Direct API access means instant results without browser overhead
- 🔧 **Scriptable**: Integrate HubSpot operations into your automation workflows
- 📊 **Flexible Output**: Choose between human-readable tables or JSON for programmatic use
- 🎯 **Focused**: Streamlined interface for contact management without UI distractions
- 🔄 **Extensible**: Built with a modular architecture ready for future HubSpot object support
- 🛡️ **Secure**: Supports multiple authentication methods including environment variables

## Installation

### macOS

#### Using Homebrew (Recommended)

```bash
brew install obay/hsctl/hsctl
```

#### Manual Installation

1. Download the latest release for macOS from the [Releases page](https://github.com/obay/hsctl/releases)
2. Extract the archive:
   ```bash
   tar -xzf hsctl_darwin_amd64_v1.0.0.tar.gz
   ```
3. Move the binary to your PATH:
   ```bash
   sudo mv hsctl /usr/local/bin/
   ```
4. Verify installation:
   ```bash
   hsctl --version
   ```

### Windows

#### Using Scoop (Recommended)

```powershell
scoop bucket add hsctl https://github.com/obay/hsctl-scoop
scoop install hsctl
```

#### Manual Installation

1. Download the latest release for Windows from the [Releases page](https://github.com/obay/hsctl/releases)
2. Extract the ZIP file
3. Add the extracted directory to your system PATH
4. Verify installation:
   ```powershell
   hsctl --version
   ```

### Linux

#### Using Package Manager (Debian/Ubuntu)

```bash
# Download and install .deb package
wget https://github.com/obay/hsctl/releases/latest/download/hsctl_linux_amd64.deb
sudo dpkg -i hsctl_linux_amd64.deb
```

#### Manual Installation

1. Download the latest release for Linux from the [Releases page](https://github.com/obay/hsctl/releases)
2. Extract the archive:
   ```bash
   tar -xzf hsctl_linux_amd64_v1.0.0.tar.gz
   ```
3. Move the binary to your PATH:
   ```bash
   sudo mv hsctl /usr/local/bin/
   ```
4. Verify installation:
   ```bash
   hsctl --version
   ```

### Build from Source

If you prefer to build from source:

```bash
git clone https://github.com/obay/hsctl.git
cd hsctl
go build -o hsctl .
sudo mv hsctl /usr/local/bin/
```

## Configuration

HSCTL requires a HubSpot API key to authenticate. You can provide it in three ways:

1. **Command-line flag** (recommended for testing):
   ```bash
   hsctl contacts list --api-key YOUR_API_KEY
   ```

2. **Environment variable** (recommended for production):
   ```bash
   export HUBSPOT_API_KEY=YOUR_API_KEY
   hsctl contacts list
   ```

3. **Config file** (optional):
   Create `~/.hsctl.yaml`:
   ```yaml
   api-key: YOUR_API_KEY
   ```

### Getting Your HubSpot API Key

1. Log in to your HubSpot account
2. Navigate to **Settings** → **Integrations** → **Private Apps**
3. Create a new private app or use an existing one
4. Ensure the app has the following scopes:
   - `crm.objects.contacts.read`
   - `crm.objects.contacts.write`
5. Copy the API key (starts with `pat-`)

## Usage

### List Contacts

List all contacts with their properties:

```bash
# List first 100 contacts
hsctl contacts list --api-key YOUR_API_KEY

# List with custom limit
hsctl contacts list --limit 50 --api-key YOUR_API_KEY

# List all contacts (paginated)
hsctl contacts list --all --api-key YOUR_API_KEY

# Output as JSON
hsctl contacts list --format json --api-key YOUR_API_KEY
```

### List Properties

View all available contact properties:

```bash
hsctl contacts properties --api-key YOUR_API_KEY

# Output as JSON
hsctl contacts properties --format json --api-key YOUR_API_KEY
```

### Create a Contact

Create a new contact:

```bash
# Basic contact
hsctl contacts create \
  --email "john.doe@example.com" \
  --firstname "John" \
  --lastname "Doe" \
  --api-key YOUR_API_KEY

# With lifecycle stage
hsctl contacts create \
  --email "jane@example.com" \
  --firstname "Jane" \
  --lastname "Smith" \
  --lifecycle-stage "customer" \
  --api-key YOUR_API_KEY

# With custom properties
hsctl contacts create \
  --email "bob@example.com" \
  --firstname "Bob" \
  --properties "company=Acme Inc,phone=555-1234" \
  --api-key YOUR_API_KEY
```

### Update a Contact

Update an existing contact:

```bash
# Update lifecycle stage
hsctl contacts update CONTACT_ID \
  --lifecycle-stage "customer" \
  --api-key YOUR_API_KEY

# Update multiple properties
hsctl contacts update CONTACT_ID \
  --firstname "John" \
  --lastname "Updated" \
  --properties "company=New Company,phone=555-9999" \
  --api-key YOUR_API_KEY
```

### Search/Query Contacts

Search for contacts:

```bash
# Search by email
hsctl contacts query "email=john@example.com" --api-key YOUR_API_KEY

# Search by property
hsctl contacts query "lifecyclestage=customer" --api-key YOUR_API_KEY

# Limit results
hsctl contacts query "email=example" --limit 10 --api-key YOUR_API_KEY
```

### Delete a Contact

Delete a contact (with confirmation):

```bash
hsctl contacts delete CONTACT_ID --api-key YOUR_API_KEY

# Skip confirmation prompt
hsctl contacts delete CONTACT_ID --force --api-key YOUR_API_KEY
```

## Examples

### Bulk Update Lifecycle Stage

```bash
# List all leads
hsctl contacts query "lifecyclestage=lead" --format json --api-key YOUR_API_KEY | \
  jq -r '.[] | .id' | \
  while read id; do
    hsctl contacts update "$id" --lifecycle-stage "customer" --api-key YOUR_API_KEY
  done
```

### Export Contacts to CSV

```bash
hsctl contacts list --format json --api-key YOUR_API_KEY | \
  jq -r '.[] | [.id, .properties.email, .properties.firstname, .properties.lastname] | @csv' > contacts.csv
```

### Find Contacts by Domain

```bash
hsctl contacts list --format json --api-key YOUR_API_KEY | \
  jq '.[] | select(.properties.email | contains("@example.com"))'
```

## Command Reference

### Global Flags

- `--api-key string`: HubSpot API key (or set HUBSPOT_API_KEY env var)
- `--config string`: Config file path (default: `$HOME/.hsctl.yaml`)
- `-h, --help`: Show help information

### Contacts Commands

#### `hsctl contacts list`
List all contacts.

**Flags:**
- `-l, --limit int`: Maximum number of contacts to retrieve (default: 100)
- `-a, --all`: Retrieve all contacts (paginate through all pages)
- `-f, --format string`: Output format - `table` or `json` (default: `table`)

#### `hsctl contacts properties`
List all available contact properties.

**Flags:**
- `-f, --format string`: Output format - `table` or `json` (default: `table`)

#### `hsctl contacts create`
Create a new contact.

**Flags:**
- `-e, --email string`: Email address
- `-f, --firstname string`: First name
- `-l, --lastname string`: Last name
- `--lifecycle-stage string`: Lifecycle stage (e.g., `lead`, `customer`)
- `-p, --properties string`: Additional properties (format: `key1=value1,key2=value2`)

#### `hsctl contacts update [contact-id]`
Update an existing contact.

**Flags:**
- `-e, --email string`: Email address
- `-f, --firstname string`: First name
- `-l, --lastname string`: Last name
- `--lifecycle-stage string`: Lifecycle stage
- `-p, --properties string`: Additional properties (format: `key1=value1,key2=value2`)

#### `hsctl contacts query [search-query]`
Search for contacts.

**Flags:**
- `-l, --limit int`: Maximum number of results (default: 100)
- `-f, --format string`: Output format - `table` or `json` (default: `table`)

**Query Format:**
- Property-based: `property=value` (e.g., `email=john@example.com`)
- Text search: `text` (searches in email field)

#### `hsctl contacts delete [contact-id]`
Delete a contact.

**Flags:**
- `--force`: Skip confirmation prompt

## Troubleshooting

### Authentication Errors

If you see authentication errors:
1. Verify your API key is correct
2. Ensure your private app has the required scopes
3. Check that the API key hasn't expired

### Rate Limiting

HubSpot has rate limits. If you encounter rate limit errors:
- Reduce the frequency of requests
- Use pagination (`--limit`) instead of `--all` for large datasets
- Implement delays in scripts

### Property Not Found

If a property update fails:
1. Use `hsctl contacts properties` to list available properties
2. Ensure property names match exactly (case-sensitive)
3. Check that the property type matches the value you're setting

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For issues, questions, or contributions, please open an issue on [GitHub](https://github.com/obay/hsctl/issues).

## Roadmap

- [ ] Support for other HubSpot objects (Deals, Companies, etc.)
- [ ] Batch operations for bulk updates
- [ ] Advanced filtering and querying options
- [ ] Import/export functionality
- [ ] Interactive mode for easier exploration

---

Made with ❤️ for the HubSpot community

