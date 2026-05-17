# Contribution Ledger

Track your private work contributions and generate a public portfolio.

**Problem:** Many developers do significant work in private repositories or on closed-source projects. Your GitHub profile doesn't reflect this work, making your contributions invisible to potential employers and collaborators.

**Solution:** Contribution Ledger scans your private repositories, code reviews, and work, then generates an anonymized public portfolio that showcases your skills and productivity.

## Features

- **Local Git Scanning** — Analyze commits across all your private repos
- **ADO Integration** — Fetch pull requests and code review metrics from Azure DevOps
- **Code Review Analytics** — Track reviews authored, feedback given, and impact
- **AI Summarization** — Claude generates high-level project summaries (sanitized)
- **Anonymous Portfolio** — Projects named "Alpha", "Beta", etc. with no company/file names
- **Contribution Graph** — Automatic commits to populate your GitHub contribution graph
- **Terminal Distribution** — Single Go binary, easy to share and run

## Installation

### From Source

```bash
git clone https://github.com/nathan815/contribution-ledger
cd contribution-ledger
go build -o contribution-ledger .
sudo mv contribution-ledger /usr/local/bin/
```

### Prebuilt Binary

(Coming soon)

## Quick Start

### 1. Initialize

```bash
contribution-ledger init
```

Prompts for:
- Azure DevOps organization
- GitHub username
- Directories to scan
- Company name (for anonymization)

Saves to `~/.contribution-ledger/config.json`

### 2. Scan Repositories

```bash
contribution-ledger scan
```

Scans all git repos in configured directories:
- Counts commits by language, date, time-of-day
- Extracts file changes, additions, deletions
- Outputs to temporary JSON file

### 3. Fetch ADO Data

```bash
contribution-ledger ado-pull-requests
```

Authenticates via `az` CLI (OAuth):
- Queries all PRs authored by you
- Fetches all PRs reviewed by you
- Extracts code review metrics
- Merges with scan data

### 4. Summarize Contributions

```bash
contribution-ledger summarize
```

Uses Claude API to generate sanitized summaries:
- "Built authentication service in Go, 347 commits"
- "Led code reviews on payment system, 38 PRs reviewed"

### 5. Push to Portfolio

```bash
contribution-ledger push
```

Creates commits in `contribution-ledger-data` repo:
- Timestamps match real activity
- Fills your GitHub contribution graph
- Updates portfolio website

## Configuration

`~/.contribution-ledger/config.json`:

```json
{
  "ado": {
    "org": "dev.azure.com/your-org",
    "projects": ["Project A", "Project B"]
  },
  "portfolioRepo": {
    "owner": "nathan815",
    "name": "contribution-ledger-data",
    "path": "/path/to/contribution-ledger-data"
  },
  "scanPaths": ["~/dev", "~/work"],
  "excludeRepos": ["node_modules", ".git", "vendor"],
  "anonymization": {
    "projectNames": {
      "RealProjectName": "Project Alpha"
    },
    "companyName": "Work",
    "stripFileNames": true
  }
}
```

## Data & Privacy

### What Gets Shared

- ✓ Programming languages (Go, TypeScript, etc.)
- ✓ Frameworks and platforms (Kubernetes, AWS, etc.)
- ✓ High-level impact ("Reduced latency by 40%")
- ✓ Commit frequency and timeline
- ✓ Code review metrics

### What's Anonymized

- ✗ Project names (→ "Project Alpha", "Project Beta")
- ✗ Company names (→ configured name)
- ✗ File paths and specific file names
- ✗ Exact commit messages
- ✗ Internal tools and APIs

## Roadmap

- [x] Phase 1: Local git scanning + CLI scaffolding
- [ ] Phase 2: ADO integration + code review analytics
- [ ] Phase 3: Claude AI summarization
- [ ] Phase 4: Portfolio repo + commit generation
- [ ] Phase 5: Website + visualization
- [ ] Phase 6: GitHub Actions auto-refresh workflow

## Development

### Build

```bash
go build -o contribution-ledger .
```

### Run

```bash
./contribution-ledger help
./contribution-ledger init
./contribution-ledger scan
```

### Dependencies

- `github.com/go-git/go-git/v5` — Git operations
- `github.com/spf13/cobra` — CLI framework
- `github.com/Azure/azure-sdk-for-go` — ADO authentication
- `github.com/anthropic-ai/sdk-go` (coming Phase 3) — Claude API

## Related

- [repos-portfolio](https://github.com/nathan815/repos-portfolio) — Showcase public repos
- [copilot-proxy](https://github.com/nathan815/copilot-proxy) — GitHub Copilot API proxy

## License

MIT

## Author

Nathan

---

**For details, see** [PLAN.md](./PLAN.md)
