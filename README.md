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
- **Common Name (CN)**: The fully qualified domain name (FQDN) for the certificate.
- **Organization (O)**: The name of your organization.
- **Organizational Unit (OU)**: The division within your organization.
- **Locality (L)**: Your city.
- **State (ST)**: Your state or province.
- **Country (C)**: The two-letter ISO country code.
- **Key File Path**: The file path where the private key will be saved.
- **CSR File Path**: The file path where the CSR will be saved.

## Examples

```bash
$ make-csr
? Enter Common Name (CN): example.com
? Enter Organization (O): My Company
? Enter Organizational Unit (OU): IT Department
? Enter Locality (L): San Francisco
? Enter State (ST): California
? Enter Country (C): US
? Enter key file path: ./example-key.pem
? Enter CSR file path: ./example-csr.pem
```

### Options
- `--cn`: Common Name (CN)
- `--o`: Organization (O)
- `--ou`: Organizational Unit (OU)
- `--l`: Locality (L)
- `--st`: State (ST)
- `--c`: Country (C)
- `--key`: Path to save the private key
- `--csr`: Path to save the CSR

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
