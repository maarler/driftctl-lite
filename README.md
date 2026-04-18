# driftctl-lite

> Lightweight CLI to detect config drift between live infrastructure and declared state

---

## Installation

```bash
go install github.com/your-org/driftctl-lite@latest
```

Or build from source:

```bash
git clone https://github.com/your-org/driftctl-lite.git
cd driftctl-lite
go build -o driftctl-lite .
```

---

## Usage

Run a drift check by pointing the CLI at your state file and a provider:

```bash
driftctl-lite scan --provider aws --state ./terraform.tfstate
```

**Example output:**

```
[✓] aws_s3_bucket.my-bucket       — in sync
[✗] aws_security_group.web        — DRIFT DETECTED: ingress rules changed
[✗] aws_iam_role.lambda-exec      — DRIFT DETECTED: policy attachment missing

2 resources drifted out of 3 scanned.
```

### Flags

| Flag | Description |
|------|-------------|
| `--provider` | Cloud provider (`aws`, `gcp`, `azure`) |
| `--state` | Path to your state file |
| `--output` | Output format: `text` (default), `json` |
| `--quiet` | Only report drifted resources |

---

## Requirements

- Go 1.21+
- Valid cloud credentials configured in your environment

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any major changes.

---

## License

[MIT](LICENSE)