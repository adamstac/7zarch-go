# 7EP-0018: Static Blog Generator

**Status:** ✅ COMPLETED  
**Author(s):** CC (Claude Code)  
**Assignment:** Completed  
**Difficulty:** 2 (simple - focused scope, clear implementation)  
**Created:** 2025-08-14  
**Updated:** 2025-08-15  
**Completed:** 2025-08-15 (PR29 + PR30 security fixes)  

## Executive Summary

Build a dead-simple static blog generator in Go that converts our markdown posts to a beautiful, fast website that looks like a respectable engineering blog. Clear, bold design that works perfectly on both mobile and desktop.

## Problem Statement

We're writing excellent blog content about AI-augmented development but it's trapped in GitHub's code view. We need:
- Professional presentation layer for our content
- Fast, accessible reading experience  
- Code syntax highlighting that makes examples shine
- Zero JavaScript bloat
- Self-contained solution (Go project → Go generator)

## Proposed Solution

### Core Architecture
```
blog/
├── posts/           # Markdown source files
├── static/          # CSS, fonts, assets
├── templates/       # Go HTML templates
├── generator.go     # ~200 line static site generator
└── public/          # Generated HTML output
```

### Technical Implementation

**Core Generator Flow**:
1. Read all `.md` files from `posts/` directory
2. Parse YAML frontmatter (title, date, author, summary)
3. Convert markdown content to HTML with syntax highlighting
4. Apply HTML templates for consistent layout
5. Generate individual post pages and index
6. Create RSS feed for subscribers
7. Copy static assets (CSS, fonts, images)

**Implementation Details**: See supporting files in `docs/7eps/7ep-0018-supporting/`:
- `generator.go` - Complete Go implementation (~200 lines)
- `post.html` - HTML template for individual posts
- `style.css` - Complete CSS with typography and color system
- `github-actions.yml` - Automated deployment configuration

**Key Technical Decisions**:
- **Goldmark** for markdown processing (GitHub-compatible)
- **Chroma** for syntax highlighting (GitHub theme)
- **Go templates** for HTML generation (no JavaScript needed)
- **Frontmatter** in YAML for post metadata

### Design System

**Typography Philosophy**:
- Massive, confident headers (3.5rem) with tight letter spacing
- Readable body text (1.125rem) with generous line height (1.75)
- Monospace fonts for code: Berkeley Mono → JetBrains Mono → SF Mono fallback

**Color Philosophy**:
- Terminal-inspired monochrome base (#0a0a0a to #fafafa)
- Classic terminal green accent (#10b981)
- Syntax highlighting optimized for readability

**Layout Principles**:
- Single column design, 680px content width
- Code blocks extend wider (920px) for better readability
- Massive whitespace, no sidebar distractions
- Mobile-first responsive design that works beautifully on all devices
- Bold, clear visual hierarchy that feels professional

**Complete Design System**: See `docs/7eps/7ep-0018-supporting/style.css` for full CSS implementation with all typography, colors, spacing, and responsive rules.

## Acceptance Criteria

### Phase 1: MVP Generator ✅ COMPLETED
- [x] Generator reads markdown files from `blog/posts/`
- [x] Converts markdown to HTML with proper formatting
- [x] Applies syntax highlighting to code blocks
- [x] Generates individual post pages
- [x] Creates index page listing all posts
- [x] Outputs static files to `blog/public/`

### Phase 2: Design Excellence ✅ COMPLETED
- [x] CSS achieves "classy hacker" aesthetic
- [x] Typography optimized for technical content
- [x] Code blocks are beautiful and readable
- [x] Mobile responsive without complexity
- [x] Page loads in <100ms on slow connection

### Phase 3: Production Features ✅ COMPLETED
- [x] RSS feed generation with Dublin Core compliance
- [x] Enhanced GitHub Pages deployment workflow with draft/published separation
- [x] Command line flags for flexible source/output directories
- [x] Safe deployment strategy preventing accidental publication
- [x] Security fixes: proper file permissions (0750), XSS protection
- [x] Reading time estimates
- [ ] Optional: ASCII art header support
- [ ] Optional: Dark mode with CSS only

## Use Cases

### Primary: Publishing Blog Posts
```bash
# Write draft (safe, won't deploy)
vim blog/drafts/shipping-with-ai.md
git add blog/drafts/shipping-with-ai.md
git commit -m "draft: new post about AI development"
git push  # ← No deployment

# Preview locally with drafts
cd blog && go run generator.go --source=posts
python -m http.server 8080 --directory public

# Publish when ready (triggers deployment)
git mv blog/drafts/shipping-with-ai.md blog/published/003-shipping-with-ai.md
git commit -m "publish: shipping with AI post"
git push  # ← Automatic deployment to production
```

### Secondary: Documentation Site
Could extend to generate docs/ as beautiful HTML, replacing need for external documentation tools.

## Non-Goals

- **No JavaScript frameworks** - Static HTML only
- **No complex build pipeline** - Single Go file
- **No external services** - Self-contained
- **No analytics/tracking** - Respect readers
- **No comments system** - Link to GitHub discussions

## Implementation Plan

### Week 1: Core Generator
1. Implement basic markdown → HTML pipeline
2. Add frontmatter parsing
3. Create minimal template system
4. Test with existing blog posts

### Week 2: Design Implementation
1. Implement clean, minimal design aesthetic
2. Implement syntax highlighting
3. Optimize typography and spacing
4. Mobile responsiveness

### Week 3: Production Polish
1. GitHub Pages deployment
2. RSS feed generation
3. Performance optimization
4. Documentation

## Technical Decisions

### Why Go Instead of Hugo/Jekyll?
- **Self-contained**: Our Go project generates its own blog
- **No dependencies**: Don't need Ruby, Node, or Python
- **Learning**: Building > configuring
- **Control**: Every line of code is ours
- **Simple**: 200 lines of Go vs 1000s of config

### Why Static Instead of Server?
- **Speed**: Nothing beats static files
- **Hosting**: GitHub Pages is free and fast
- **Security**: No server, no vulnerabilities
- **Focus**: Write content, not maintain servers

## Success Metrics

- Blog loads in <100ms globally
- Code examples are more readable than GitHub
- Generator builds all posts in <1 second
- Total implementation <500 lines of Go
- CSS smaller than 10KB

## Risk Assessment

**Low Risk**:
- Scope is intentionally minimal
- Technology (Go + HTML + CSS) is proven
- Can fallback to existing tools if needed

**Deployment Risk - IDENTIFIED**:
- Current design auto-publishes any `blog/` changes to production
- No staging environment or content review process
- Risk of accidental publication of drafts or work-in-progress

**Mitigation**:
- Start with working prototype in 1 day
- Iterate on design after core works
- Keep generator simple enough to rewrite
- **Enhanced deployment strategy** (see Deployment Strategy section below)

## Alternative Approaches Considered

1. **Hugo** - Too much configuration, not our code
2. **Next.js** - JavaScript complexity we don't need  
3. **GitHub Pages with Jekyll** - Ruby dependency
4. **README rendering** - Not professional enough

## References

- Markdown parser: https://github.com/yuin/goldmark
- Syntax highlighter: https://github.com/alecthomas/chroma
- Similar approach: https://github.com/brandur/sorg

## Design Inspiration

Reference sites for design principles:
- https://coder.com/blog - Clean technical blog with excellent typography
- https://brandur.org - Minimal, content-focused design
- https://danluu.com - Speed-obsessed, zero-bloat approach

**Goal**: Build a respectable engineering blog that engineers would be proud to read and share.

## Deployment Strategy

### Current Design Problem
The initial GitHub Actions workflow (`docs/7eps/7ep-0018-supporting/github-actions.yml`) triggers on **any push to main that touches `blog/`**, meaning:
- ❌ Draft posts publish immediately 
- ❌ No content review process
- ❌ Work-in-progress goes public
- ❌ No staging environment

### Enhanced Deployment Options

#### Option 1: Draft/Published Directory Structure ⭐ **RECOMMENDED**
```
blog/
├── drafts/          # Safe to edit, never published
│   ├── wip-post.md
│   └── ideas.md
├── published/       # Only these deploy to production
│   ├── 001-ddd.md
│   └── 002-agents.md
└── templates/
```

**Deployment Trigger**:
```yaml
on:
  push:
    branches: [main]
    paths: ['blog/published/**']  # Only published directory
```

**Publishing Workflow**:
```bash
# Write draft
vim blog/drafts/my-post.md

# Review and approve
git add blog/drafts/my-post.md
git commit -m "draft: new post about X"

# Publish when ready
git mv blog/drafts/my-post.md blog/published/003-my-post.md
git commit -m "publish: my post about X"
git push  # Triggers deployment
```

**Benefits**:
- ✅ Clear draft vs published separation
- ✅ Safe to iterate on drafts in main
- ✅ Intentional publishing via file move
- ✅ Git history shows publish decisions

#### Option 2: Manual Deployment Trigger
```yaml
on:
  workflow_dispatch:  # Manual trigger only
    inputs:
      reason:
        description: 'Reason for deployment'
        required: true
```

**Publishing Workflow**:
1. Push blog changes to main
2. Go to Actions tab in GitHub
3. Click "Deploy Blog" → "Run workflow"
4. Enter deployment reason → "Run workflow"

**Benefits**:
- ✅ Complete control over when content goes live
- ✅ Deployment reason for audit trail
- ✅ No accidental publications
- ❌ Extra step required for each publish

#### Option 3: Release-Based Publishing
```yaml
on:
  release:
    types: [published]
```

**Publishing Workflow**:
```bash
# Tag blog content release
git tag blog-v1.2.0
git push origin blog-v1.2.0

# Create release in GitHub UI
# → Triggers deployment
```

**Benefits**:
- ✅ Versioned blog content
- ✅ Release notes for content updates
- ✅ Clear publication history
- ❌ Heavy process for individual posts

#### Option 4: Hybrid Approach
```yaml
on:
  workflow_dispatch:      # Manual trigger always available
  push:
    branches: [main]
    paths: ['blog/published/**']  # Auto-deploy published content
```

**Benefits**:
- ✅ Automatic for published content
- ✅ Manual override available
- ✅ Best of both approaches

### Implementation Recommendation

**Use Option 1 (Draft/Published Directories)** because:
- Maintains development velocity (drafts in main)
- Clear publication intent (file move)
- Prevents accidental publication
- Simple mental model
- Preserves git-based workflow

### Updated GitHub Actions Workflow

Replace existing `docs/7eps/7ep-0018-supporting/github-actions.yml` with:

```yaml
name: Deploy Blog

on:
  workflow_dispatch:
    inputs:
      reason:
        description: 'Deployment reason'
        required: false
  push:
    branches: [main]
    paths: ['blog/published/**', 'blog/templates/**', 'blog/static/**']

jobs:
  build-and-deploy:
    runs-on: ubuntu-latest
    
    permissions:
      contents: read
      pages: write
      id-token: write
    
    steps:
      - name: Checkout
        uses: actions/checkout@v4
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Build blog
        run: |
          cd blog
          mkdir -p public/static
          # Only process published posts
          go run generator.go --source=published --output=public
      
      - name: Setup Pages
        uses: actions/configure-pages@v3
      
      - name: Upload artifact
        uses: actions/upload-pages-artifact@v2
        with:
          path: ./blog/public
      
      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v2
        
      - name: Log deployment
        run: |
          echo "::notice::Blog deployed successfully"
          echo "::notice::Trigger: ${{ github.event_name }}"
          if [ "${{ github.event_name }}" = "workflow_dispatch" ]; then
            echo "::notice::Reason: ${{ github.event.inputs.reason }}"
          fi
```

### Content Organization Benefits

**For Writers**:
```bash
# Safe experimentation in drafts
blog/drafts/exploring-idea.md      # Won't go live
blog/drafts/half-finished.md       # Safe to commit WIP

# Intentional publishing
git mv blog/drafts/ready.md blog/published/004-ready.md
```

**For Reviewers**:
- All drafts visible in main branch
- Publishing requires explicit file move
- Clear intent in git history
- No surprise publications

**For Automation**:
- Templates/static changes deploy immediately (safe)
- Only published content triggers production build
- Manual override always available
- Deployment logging for audit trail

This strategy balances development flexibility with publication control, ensuring no accidental content releases while maintaining a smooth writing workflow.

### Enhanced Generator Requirements

To support the draft/published directory structure, the `generator.go` needs these enhancements:

**Command Line Flags**:
```go
var (
    sourceDir = flag.String("source", "posts", "Source directory for markdown files") 
    outputDir = flag.String("output", "public", "Output directory for generated files")
    watch     = flag.Bool("watch", false, "Watch for file changes and rebuild")
)
```

**Updated Main Function**:
```go
func main() {
    flag.Parse()
    fmt.Printf("Building blog from %s/ to %s/\n", *sourceDir, *outputDir)
    
    // Load posts from configurable source directory
    posts := loadPosts(*sourceDir + "/")
    
    // Generate to configurable output directory  
    generateSite(posts, *outputDir)
}
```

**Usage Examples**:
```bash
# Development: process all posts (drafts + published)
go run generator.go --source=posts --output=dev-public

# Production: only published posts  
go run generator.go --source=published --output=public

# Local development with drafts
go run generator.go --source=posts --watch
```

This enhancement maintains backward compatibility (`posts/` default) while enabling the safe deployment strategy.

## Appendix A: Example Frontmatter

```yaml
---
title: "Document-Driven Development: How We Ship 4,300 Lines in 3 Days"
date: 2025-08-14
author: Claude Code (CC)
slug: document-driven-development
summary: The future of software development isn't AI writing all the code. It's humans and AI agents working together with zero coordination friction.
---
```

## Appendix B: Supporting Files Overview

All detailed implementation files are provided in `docs/7eps/7ep-0018-supporting/`:

1. **`generator.go`** - Complete static site generator implementation
2. **`post.html`** - HTML template for individual blog posts  
3. **`style.css`** - Full CSS design system with typography and colors
4. **`github-actions.yml`** - Automated deployment to GitHub Pages

These files contain the complete, production-ready implementation of the blog system described in this 7EP.

## Appendix D: Development Workflow

```bash
# Local development with file watching
cd blog
go run generator.go --watch

# Serves on localhost:8080 with live reload
# Changes to posts/*.md trigger rebuild
# Changes to templates/*.html trigger rebuild
# Changes to static/*.css are instant
```

## Appendix E: AI Image Generation Options

### Option 1: Prompt-Based Placeholders
```markdown
![AI: A terminal window floating in cyberspace with green text cascading like rain, cyberpunk aesthetic, dark and moody]
```

Generator detects `![AI: ...]` patterns and:
1. Generates subtle gradient placeholder
2. Adds prompt as alt text
3. Future: Could hook to DALL-E/Midjourney API

### Option 2: ASCII Art Headers
```markdown
---
ascii-header: |
  ╔════════════════════════╗
  ║  DOCUMENT >> DRIVEN    ║
  ║      DEVELOPMENT       ║
  ╚════════════════════════╝
---
```

More authentic to hacker aesthetic, zero dependencies.

## Appendix F: Progressive Enhancement Roadmap

**Phase 1 (MVP)**: Static HTML/CSS only
- No JavaScript at all
- Fast everywhere

**Phase 2 (Enhanced)**: Minimal JS for QOL
- Copy code button (10 lines of JS)
- Keyboard navigation (j/k for next/prev)

**Phase 3 (Optional)**: Dark mode
- CSS-only using prefers-color-scheme
- Optional toggle with localStorage

**Never**: 
- React, Vue, or any framework
- Analytics or tracking
- External dependencies
- Build complexity

## Decision

✅ **APPROVED AND IMPLEMENTED**

**Implementation Summary:**
- **PR29:** Core implementation merged (blog generator, design system, GitHub Actions)
- **PR30:** Security fixes merged (gosec compliance, XSS protection, RSS improvements)
- **Production Status:** Live at https://7zarch.com with secure, compliant blog generator

This is now our production blog platform, showcasing both our content and our philosophy: build simple tools that do one thing exceptionally well.

---

*7-Zip with a brain deserves a blog with soul.*