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
│  - AI summarization (Claude)                                     │
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
# Passes commit/PR data to Claude API
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
    },
    {
      "id": "project-beta-migrations",
      "name": "Project Beta",
      "company": "Work",
      "period": "2025-08 to 2025-10",
      "type": "infrastructure",
      "stats": {
        "commits": 156,
        "filesChanged": 45,
        "frequency": "8-12 commits/week",
        "languages": ["Go", "Terraform"],
        "primaryLanguage": "Terraform"
      },
      "aiSummary": "Led infrastructure migration from legacy datacenter to Kubernetes. Designed multi-region failover strategy.",
      "technologies": ["Terraform", "Kubernetes", "AWS"]
    }
  ],
  "codeReviews": {
    "totalReviewed": 156,
    "totalComments": 892,
    "averageCommentsPerPR": 5.7,
    "feedbackTypes": {
      "approved": 78,
      "requestedChanges": 42,
      "commented": 36
    },
    "highlights": [
      {
        "date": "2026-04-15",
        "project": "Project Alpha",
        "type": "Performance review",
        "summary": "Reviewed critical database optimization PR. Suggested indexes on hot tables, reducing query time by 60%."
      },
      {
        "date": "2026-03-22",
        "project": "Project Beta",
        "type": "Architecture review",
        "summary": "Reviewed service mesh migration proposal. Provided feedback on observability patterns."
      }
    ]
  },
  "timeline": {
    "2026-05": { "commits": 87, "prsReviewed": 12, "focus": ["Project Alpha"] },
    "2026-04": { "commits": 156, "prsReviewed": 18, "focus": ["Project Alpha", "Project Beta"] },
    "2026-03": { "commits": 203, "prsReviewed": 24, "focus": ["Project Alpha"] }
  },
  "skills": {
    "languages": ["Go", "TypeScript", "Python", "Terraform"],
    "platforms": ["Kubernetes", "AWS", "PostgreSQL", "Redis"],
    "practices": ["Code Review", "Architecture Design", "Performance Optimization"]
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

### GitHub (Token in config, or prompt)
```bash
# Optional: Store GitHub token in config for private repos
GITHUB_TOKEN=$(cat ~/.contribution-ledger/config.json | jq -r .github.token)
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
   → Calls Claude: "Summarize this work in one sentence: 347 commits in Go, focused on auth..."
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

**CLI:**
- **Language:** Node.js (TypeScript) or Go
  - Node: Better integration with Claude API, faster to prototype
  - Go: Faster, smaller binary, easier distribution
  - Recommendation: **Node/TypeScript** for v1, can be rewritten in Go later
- **Dependencies:**
  - `@azure/identity` or `az CLI` for ADO auth
  - `@anthropic-ai/sdk` for Claude summarization
  - `simple-git` for git operations
  - `commander` for CLI
  - `fs-extra` for file ops

**Data Repo:**
- Plain Git repo with JSON + HTML
- GitHub Pages for website
- GitHub Actions for optional auto-refresh

**Website:**
- Vanilla JS (load repos.json, render)
- Or React for more interactivity
- Terminal theme (consistent with repos-portfolio)

---

## Implementation Phases

### Phase 1: MVP (Local Scanning + Anonymization)
- [ ] Project setup (package.json, tsconfig, etc.)
- [ ] CLI scaffolding (commander)
- [ ] Config file structure
- [ ] Git scanning logic
- [ ] Anonymization module
- [ ] Save to JSON

### Phase 2: ADO Integration
- [ ] `az` CLI auth
- [ ] ADO REST API wrapper
- [ ] PR/review querying
- [ ] Code review stats extraction
- [ ] Merge with scan data

### Phase 3: AI Summarization
- [ ] Claude API integration
- [ ] Prompt engineering for sanitized summaries
- [ ] Batch summarization
- [ ] Error handling/retries

### Phase 4: Portfolio Push
- [ ] Create commits with realistic timestamps
- [ ] Push to data repo
- [ ] GitHub Pages setup

### Phase 5: Website
- [ ] Build HTML + CSS
- [ ] Render repos.json dynamically
- [ ] Timeline visualization
- [ ] Contribution graph (ASCII art or Chart.js)
- [ ] Stats dashboard

### Phase 6: Polish + Open Source
- [ ] Documentation
- [ ] Tests
- [ ] Example config
- [ ] Publish to npm/GitHub

---

## Next Steps

1. Confirm tech stack (Node.js TypeScript?)
2. Confirm data repo name (`contribution-ledger-data`?)
3. Start Phase 1: scaffolding + config + git scanning

---

**Created:** May 17, 2026
**Status:** Planning
