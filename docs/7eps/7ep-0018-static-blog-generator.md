# 7EP-0018: Static Blog Generator

**Status:** Draft  
**Author(s):** CC (Claude Code)  
**Assignment:** Unassigned  
**Difficulty:** 2 (simple - focused scope, clear implementation)  
**Created:** 2025-08-14  
**Updated:** 2025-08-14  

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

### Phase 1: MVP Generator
- [ ] Generator reads markdown files from `blog/posts/`
- [ ] Converts markdown to HTML with proper formatting
- [ ] Applies syntax highlighting to code blocks
- [ ] Generates individual post pages
- [ ] Creates index page listing all posts
- [ ] Outputs static files to `blog/public/`

### Phase 2: Design Excellence  
- [ ] CSS achieves "classy hacker" aesthetic
- [ ] Typography optimized for technical content
- [ ] Code blocks are beautiful and readable
- [ ] Mobile responsive without complexity
- [ ] Page loads in <100ms on slow connection

### Phase 3: Production Features
- [ ] RSS feed generation
- [ ] GitHub Pages deployment workflow
- [ ] Optional: ASCII art header support
- [ ] Optional: Dark mode with CSS only
- [ ] Optional: Reading time estimates

## Use Cases

### Primary: Publishing Blog Posts
```bash
# Write post
vim blog/posts/003-shipping-with-ai.md

# Generate site
cd blog && go run generator.go

# Preview locally
python -m http.server 8080 --directory public

# Deploy (automatic via GitHub Actions)
git push origin main
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

**Mitigation**:
- Start with working prototype in 1 day
- Iterate on design after core works
- Keep generator simple enough to rewrite

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

Study these sites for their taste level and design decisions, but create our own unique aesthetic:
- https://coder.com/blog - Clean technical blog with excellent typography
- https://brandur.org - Minimal, content-focused design
- https://danluu.com - Speed-obsessed, zero-bloat approach

**Goal**: Build a respectable engineering blog that engineers would be proud to read and share. Learn from what makes these sites readable, professional, and fast, then build something uniquely ours. We're not copying - we're understanding why certain design decisions work for technical content and applying those principles in our own way.

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

_[Pending approval]_

Once approved, this becomes our blog platform, showcasing both our content and our philosophy: build simple tools that do one thing exceptionally well.

---

*7-Zip with a brain deserves a blog with soul.*