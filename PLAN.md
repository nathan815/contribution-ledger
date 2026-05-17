# Contribution Ledger — Project Plan

**Goal:** CLI tool to track private work contributions (commits + code reviews) and generate a shareable public portfolio with anonymized summaries and contribution graph.

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│  CLI Tool (contribution-ledger-cli)                              │
│  - Scan local git repos                                          │
│  - Query ADO (via `az`) for PRs, code reviews                   │
│  - Extract commit stats, review metrics                          │
│  - AI summarization (Copilot SDK / Claude)                      │
│  - Anonymize sensitive data                                      │
│  - Generate JSON output                                          │
└────────────────────┬────────────────────────────────────────────┘
                     │
                     ↓ (push)
┌─────────────────────────────────────────────────────────────────┐
│  Public Data Repo (contribution-ledger-data)                    │
│  - Store anonymized summaries                                    │
│  - Commit history mirrors private activity                      │
│  - JSON + index.html for visualization                          │
└────────────────────┬────────────────────────────────────────────┘
                     │
                     ↓ (render)
┌─────────────────────────────────────────────────────────────────┐
│  Website (GitHub Pages)                                         │
│  - Timeline of work                                              │
│  - Contribution graph                                            │
│  - Stats dashboard                                               │
│  - Code review highlights                                        │
└─────────────────────────────────────────────────────────────────┘
```

---

## CLI Commands

### 1. Initialize
```bash
contribution-ledger init
# Prompts for:
# - Azure DevOps organization URL
# - (Optional) GitHub token for private repos
# Saves to ~/.contribution-ledger/config.json
# Verifies auth with `az login`
```

### 2. Scan Local Repos
```bash
contribution-ledger scan [--since 2025-01-01] [--dir ~/dev]
# Scans all git repos in directory
# Extracts:
#  - Commit count by language, date, time-of-day
#  - Files touched, lines changed
#  - Commit frequency patterns
# Output: temp JSON
```

### 3. Fetch ADO Data
```bash
contribution-ledger ado-pull-requests [--project "Project Alpha"]
# Uses `az` CLI to get token
# Queries ADO REST API for:
#  - PRs authored (count, avg review time)
#  - PRs reviewed (count, comments, feedback type)
#  - Code review velocity (comments per PR)
#  - Time spent reviewing
# Anonymizes project/repo names
# Output: merged with scan data
```

### 4. AI Summarize
```bash
contribution-ledger summarize
# Passes commit/PR data to Copilot SDK (Claude)
# Generates high-level summaries like:
#  - "Built authentication service (Go, TypeScript), 145 commits over 2 months"
#  - "Led code reviews on performance improvements, 38 PRs reviewed"
# Strips technical debt mentions, company names, etc.
# Output: summary.json
```

### 5. Push to Portfolio Repo
```bash
contribution-ledger push --repo contribution-ledger-data
# Creates commits in the data repo
# Commit timestamps match real private activity
# Fills contribution graph
# Updates JSON summaries
```

---

## Data Structure

### Local Config (`~/.contribution-ledger/config.json`)
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
  "scanPaths": [
    "~/dev",
    "~/work"
  ],
  "excludeRepos": [
    "node_modules",
    ".git",
    "terraform"
  ],
  "anonymization": {
    "projectNames": {
      "RealProjectName": "Project Alpha"
    },
    "companyName": "Work",
    "stripFileNames": true
  },
  "ai": {
    "provider": "copilot",
    "apiKey": ""
  }
}
```

### Output Schema (`repos.json` in portfolio repo)

```json
{
  "metadata": {
    "generatedAt": "2026-05-17T05:42:00Z",
    "period": "2025-01-01 to 2026-05-17",
    "totalCommits": 1247,
    "totalCodeReviewsLeft": 156
  },
  "projects": [
    {
      "id": "project-alpha-auth",
      "name": "Project Alpha",
      "company": "Work",
      "period": "2025-11 to 2026-04",
      "type": "backend",
      "stats": {
        "commits": 347,
        "filesChanged": 89,
        "linesAdded": 12450,
        "linesDeleted": 3200,
        "frequency": "5-10 commits/week",
        "languages": ["Go", "TypeScript", "Python"],
        "primaryLanguage": "Go"
      },
      "aiSummary": "Built distributed authentication service with OAuth2 support. Implemented caching layer (Redis), optimized database queries, reduced auth latency by 40%.",
      "technologies": ["Go", "PostgreSQL", "Redis", "Kubernetes"],
      "impact": "Improved login reliability from 98% to 99.9% uptime"
    }
  ],
  "codeReviews": {
    "totalReviewed": 156,
    "totalComments": 892,
    "averageCommentsPerPR": 5.7,
    "highlights": [
      {
        "date": "2026-04-15",
        "project": "Project Alpha",
        "summary": "Reviewed critical database optimization PR. Suggested indexes, reducing query time by 60%."
      }
    ]
  },
  "timeline": {
    "2026-05": { "commits": 87, "prsReviewed": 12 },
    "2026-04": { "commits": 156, "prsReviewed": 18 }
  }
}
```

---

## Anonymization Strategy

### Data Stripped/Anonymized
- ✗ Actual project names → "Project Alpha", "Project Beta", etc.
- ✗ Company names → "Work" (configurable)
- ✗ File paths and specific file names
- ✗ Exact commit messages (summarized instead)
- ✗ Private APIs, internal tools

### Data Kept (Public-safe)
- ✓ Languages (Go, TypeScript, etc.)
- ✓ Frameworks/platforms (Kubernetes, AWS, etc.)
- ✓ High-level impact ("Reduced latency by 40%")
- ✓ Commit frequency
- ✓ Review metrics
- ✓ Timeline of work

---

## Auth Flow

### ADO (OAuth via `az` CLI)

```bash
# 1. CLI checks if already logged in
az account show 2>/dev/null

# 2. If not, prompt user
az login

# 3. Get ADO token
TOKEN=$(az account get-access-token --resource 499b84ac-1321-427f-aa17-267ca6975798 --query accessToken -o tsv)

# 4. Use token for ADO REST API queries
curl -H "Authorization: Bearer $TOKEN" \
  https://dev.azure.com/{org}/_apis/git/repositories
```

### Copilot SDK (GitHub authentication)

```bash
# Use GitHub token from config or environment
# SDK handles authentication automatically
# Falls back to Claude SDK if Copilot not available
```

---

## Data Flow Example

```
1. User runs: contribution-ledger scan --since 2025-01-01
   → Scans ~/dev for all git repos
   → Finds: prscope, pr-reviewer, dotfiles, etc.
   → Extracts commits by author, language, date

2. User runs: contribution-ledger ado-pull-requests
   → Uses `az` to auth
   → Queries ADO org for all PRs authored/reviewed by user
   → Merges with commit data

3. User runs: contribution-ledger summarize
   → Calls Copilot SDK: "Summarize this work in one sentence: 347 commits in Go..."
   → Returns: "Built authentication service with OAuth2 support..."
   → Generates: projects[].aiSummary

4. User runs: contribution-ledger push --repo ~/code/contribution-ledger-data
   → For each project, creates fake commit(s) with timestamps matching real activity
   → Dates spread across 2025-11 to 2026-04
   → Commit messages: "Work: Project Alpha (347 commits, Go/TypeScript)"
   → Updates repos.json
   → Pushes to GitHub → contribution graph populates

5. Website renders ~/code/contribution-ledger-data/index.html
   → Loads repos.json
   → Shows timeline, languages, code review stats
   → Contribution graph visible in GitHub profile
```

---

## Tech Stack

**Language:** Go 1.21+
- Faster execution, smaller binary
- Easy distribution (single executable)
- Better for CLI tools

**Dependencies:**
- `github.com/go-git/go-git/v5` — Git operations
- `github.com/spf13/cobra` — CLI framework
- `github.com/mitchellh/go-homedir` — Cross-platform paths
- `github.com/Azure/azure-sdk-for-go/sdk/identity` — ADO OAuth
- **`github.com/copilot-sdk/go-copilot`** — Claude API via GitHub Copilot (default)
- `github.com/anthropic-ai/sdk-go` — Claude API (optional fallback)

**Website:**
- Vanilla JS (load repos.json, render)
- Terminal theme (consistent with repos-portfolio)
- GitHub Pages deployment

---

## Implementation Phases

### Phase 1: MVP (Local Scanning + Anonymization) ✅
- [x] Project setup (Go module, structure)
- [x] CLI scaffolding (Cobra)
- [x] Config file system
- [x] Git scanning logic
- [x] Anonymization framework
- [x] Save to JSON
- Tests: 9 passing

### Phase 2: ADO Integration ✅
- [x] `az` CLI auth
- [x] ADO REST API client
- [x] PR/review querying
- [x] Code review stats extraction
- [x] Concurrent API calls
- [x] Merge with scan data
- Tests: 11 passing, 9% coverage

### Phase 3: AI Summarization (IN PROGRESS)
- [ ] Copilot SDK integration
- [ ] Prompt engineering for sanitized summaries
- [ ] Batch summarization
- [ ] Error handling/retries
- [ ] Claude SDK as fallback option

### Phase 4: Portfolio Push
- [ ] Git commit generation with real timestamps
- [ ] Data repo management
- [ ] Contribution graph population
- [ ] GitHub Pages setup

### Phase 5: Website
- [ ] Build HTML + CSS
- [ ] Render repos.json dynamically
- [ ] Timeline visualization
- [ ] Contribution graph (ASCII art or Chart.js)
- [ ] Stats dashboard

### Phase 6: Polish + Open Source
- [ ] Documentation
- [ ] Full test suite
- [ ] Example config
- [ ] Publish releases

---

## Next Steps

1. ✅ Phase 1 complete
2. ✅ Phase 2 complete
3. → Phase 3: Copilot SDK summarization
4. → Phase 4: Portfolio push
5. → Phase 5: Website
6. → Phase 6: Polish

---

**Created:** May 17, 2026
**Updated:** May 17, 2026 (Phases 1-2 complete, Copilot SDK default)
**Status:** Building Phase 3
