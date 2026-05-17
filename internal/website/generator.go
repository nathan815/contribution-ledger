package website

import (
	"fmt"
	"os"
	"path/filepath"
)

// Generator creates a website from portfolio data
type Generator struct {
	outputPath string
}

// NewGenerator creates a new website generator
func NewGenerator(outputPath string) *Generator {
	return &Generator{outputPath: outputPath}
}

// Generate creates the portfolio website
func (g *Generator) Generate(portfolioJSONPath string) error {
	fmt.Println("🌐 Generating portfolio website...")

	// Read portfolio JSON
	jsonData, err := os.ReadFile(portfolioJSONPath)
	if err != nil {
		return fmt.Errorf("failed to read portfolio data: %w", err)
	}

	// Create index.html
	html := g.buildHTML(string(jsonData))

	indexPath := filepath.Join(g.outputPath, "index.html")
	if err := os.WriteFile(indexPath, []byte(html), 0644); err != nil {
		return fmt.Errorf("failed to write index.html: %w", err)
	}

	fmt.Printf("✅ Website generated at: %s\n", indexPath)
	return nil
}

// buildHTML constructs the HTML page
func (g *Generator) buildHTML(portfolioJSON string) string {
	htmlTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Work Portfolio</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 40px 20px;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
        }
        
        header {
            text-align: center;
            color: white;
            margin-bottom: 60px;
        }
        
        h1 {
            font-size: 3em;
            margin-bottom: 10px;
            font-weight: 700;
        }
        
        .subtitle {
            font-size: 1.2em;
            opacity: 0.9;
        }
        
        .projects-grid {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
            gap: 20px;
            margin-bottom: 40px;
        }
        
        .project-card {
            background: white;
            border-radius: 8px;
            padding: 24px;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
            transition: transform 0.3s, box-shadow 0.3s;
        }
        
        .project-card:hover {
            transform: translateY(-4px);
            box-shadow: 0 12px 20px rgba(0,0,0,0.15);
        }
        
        .project-title {
            font-size: 1.3em;
            font-weight: 600;
            margin-bottom: 8px;
            color: #333;
        }
        
        .project-summary {
            color: #666;
            font-size: 0.95em;
            line-height: 1.5;
            margin-bottom: 12px;
        }
        
        .stats {
            display: flex;
            justify-content: space-between;
            margin: 12px 0;
            padding: 12px 0;
            border-top: 1px solid #eee;
            border-bottom: 1px solid #eee;
        }
        
        .stat {
            text-align: center;
        }
        
        .stat-number {
            font-weight: 700;
            color: #667eea;
            font-size: 1.5em;
        }
        
        .stat-label {
            color: #999;
            font-size: 0.8em;
        }
        
        .languages {
            display: flex;
            gap: 8px;
            margin-top: 12px;
            flex-wrap: wrap;
        }
        
        .language-badge {
            background: #f0f0f0;
            color: #333;
            padding: 4px 12px;
            border-radius: 20px;
            font-size: 0.85em;
        }
        
        .stats-section {
            background: white;
            border-radius: 8px;
            padding: 30px;
            text-align: center;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        }
        
        .stats-number {
            font-size: 2.5em;
            font-weight: 700;
            color: #667eea;
            margin-bottom: 10px;
        }
        
        @media (max-width: 768px) {
            h1 { font-size: 2em; }
            .projects-grid { grid-template-columns: 1fr; }
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>Work Portfolio</h1>
            <p class="subtitle">Professional contributions & impact</p>
        </header>
        
        <div class="projects-grid" id="projects"></div>
        
        <div class="stats-section">
            <div class="stats-number" id="total-commits">-</div>
            <p style="font-size: 1.1em; color: #666;">Total Commits</p>
        </div>
    </div>
    
    <script>
        const portfolioData = ` + portfolioJSON + `;
        
        document.addEventListener('DOMContentLoaded', function() {
            // Render projects
            const projectsContainer = document.getElementById('projects');
            portfolioData.projects.forEach(project => {
                const card = document.createElement('div');
                card.className = 'project-card';
                
                const languagesHTML = project.technologies
                    .map(lang => '<span class="language-badge">' + lang + '</span>')
                    .join('');
                
                card.innerHTML = '<div class="project-title">' + project.name + '</div>' +
                    '<div class="project-summary">' + project.aiSummary + '</div>' +
                    '<div class="stats">' +
                        '<div class="stat">' +
                            '<div class="stat-number">' + project.stats.commits + '</div>' +
                            '<div class="stat-label">Commits</div>' +
                        '</div>' +
                        '<div class="stat">' +
                            '<div class="stat-number">' + project.stats.filesChanged + '</div>' +
                            '<div class="stat-label">Files</div>' +
                        '</div>' +
                    '</div>' +
                    '<div class="languages">' + languagesHTML + '</div>';
                
                projectsContainer.appendChild(card);
            });
            
            // Update stats
            document.getElementById('total-commits').textContent = portfolioData.metadata.totalCommits;
        });
    </script>
</body>
</html>`

	return htmlTemplate
}
