# make-csr-cli

`make-csr-cli` is a command-line tool for generating Certificate Signing Requests (CSRs) easily and interactively. This tool simplifies the process of creating CSRs by providing a user-friendly prompt for entering required details.

## Features
- Interactive command-line prompts for input.
- Generates a CSR and private key in one step.
- Saves the CSR and key to specified files.
- Works cross-platform: macOS, Linux, and Windows.

## Installation

Install `make-csr-cli` globally using npm:

```bash
npm install -g make-csr-cli
```

## Usage

Run the command in your terminal:

```bash
make-csr
```

### Interactive Prompts
When you run `make-csr`, you will be prompted to enter the following details:
- **Country (C)**: The two-letter ISO country code.
- **State (ST)**: Your state or province.
- **Locality (L)**: Your locality name.
- **Organization (O)**: The name of your organization.
- **Organizational Unit (OU)**: The division within your organization.

- **Common Name (CN)**: The fully qualified domain name (FQDN) for the certificate.
## Examples

```bash
$ make-csr
Configuration file not found. Please provide the following details:
? Enter Country Code (C) [e.g., JP]:
? Enter State or Province Name (ST) [e.g., Tokyo]:
? Enter Locality Name (L) [e.g., CHUOU-KU]:
? Enter Organization Name (O) [e.g., BRIDGE CO.,LTD.]:
? Enter Organizational Unit Name (OU) [e.g., BRIDGE CO.,LTD.]:
Common name not found. Do you want to provide content? [Y/n]:
```

### Options
- `--cn`: Common Name (CN)
- `--o`: Organization (O)
- `--ou`: Organizational Unit (OU)
- `--l`: Locality (L)
- `--st`: State (ST)
- `--c`: Country (C)

## Compatibility

`make-csr-cli` works on the following platforms:
- macOS
- Linux
- Windows

## Common Issues

### Problem: "Command not found"
**Solution:** Ensure `make-csr-cli` is installed globally. If the issue persists, check that the npm global bin directory is in your `PATH`.

```bash
echo $PATH
```

### Problem: "Permission denied" when saving files
**Solution:** Verify that you have write permissions for the specified file paths.

## License

`make-csr-cli` is open-source software licensed under the MIT License. Contributions are welcome!
