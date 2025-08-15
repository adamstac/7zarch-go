# 7zarch Blog Management Guide

**Purpose**: Dead-simple blog workflow for Adam to manage content creation and publication  
**Framework**: Draft-first workflow with safe publication controls  
**No Accidents**: Impossible to accidentally publish drafts to production

## 🎯 Quick Commands (Copy-Paste Ready)

### Draft Content
```
Draft a blog post about [topic]
```
**Result**: Creates `blog/drafts/[slug].md` - safe to iterate and commit

### Review Content
```
Show me the draft about [topic]
Preview the blog locally
What blog drafts do we have?
```
**Result**: Review content before publication decision

### Publish Content  
```
Publish the [topic] post
Publish this as post #5
```
**Result**: Moves to `blog/published/` and triggers automatic deployment

### Status & Management
```
What's published on the blog?
What blog drafts do we have?
Deploy the blog manually
```

## 📁 Directory Structure

```
blog/
├── README.md           # This file - your blog management guide
├── drafts/            # SAFE WORKSPACE - never published
│   ├── exploring-ai.md   # Work in progress
│   ├── team-ideas.md     # Ideas and notes  
│   └── half-finished.md  # Safe to abandon
├── published/         # AUTO-DEPLOYS to production
│   ├── 001-ddd-framework.md    # Live on blog
│   └── 002-ai-agents.md        # Live on blog
├── templates/         # HTML templates (when implemented)
├── static/           # CSS, fonts, assets (when implemented)  
└── generator.go      # Static site generator (when implemented)
```

## 🚦 Publishing Workflow

### Phase 1: Safe Drafting
```bash
# AI Agent creates draft
vim blog/drafts/my-topic.md

# Safe to commit and push
git add blog/drafts/my-topic.md
git commit -m "draft: exploring my topic"  
git push  # ← NO DEPLOYMENT - completely safe
```

**Key Point**: Everything in `drafts/` is safe. Commit, push, iterate without any risk of publication.

### Phase 2: Review & Approval
```bash
# Preview locally (includes drafts)
cd blog && go run generator.go --source=posts --output=dev-public
python -m http.server 8080 --directory dev-public
# Visit: http://localhost:8080
```

**Review Process**: Read the draft, make edits, ensure it's ready for public consumption.

### Phase 3: Publication Decision
```bash
# Move to published (triggers deployment)
git mv blog/drafts/my-topic.md blog/published/003-my-topic.md
git commit -m "publish: my topic exploration"
git push  # ← DEPLOYMENT TRIGGERED - goes live automatically
```

**Key Point**: Only moving files to `published/` triggers deployment. This is intentional and explicit.

## 🛡️ Safety Features

### No Accidental Publishing
- **Drafts directory**: Never deployed, always safe
- **Published directory**: Only place that triggers deployment  
- **Explicit file move**: Publishing requires intentional `git mv` command
- **Git history**: Clear record of publication decisions

### Manual Override Available
```
Deploy the blog manually
```
**When to use**: Emergency deployments, template changes, or when automation fails

**How it works**: Uses GitHub Actions manual trigger with deployment reason tracking

### Emergency Unpublish
```bash
# Move back to drafts (rarely needed)
git mv blog/published/003-bad-post.md blog/drafts/003-bad-post.md
git commit -m "unpublish: need more work"
git push  # Removes from live blog
```

## 📝 Content Creation Patterns

### Blog Post Frontmatter
Every post needs YAML frontmatter:
```yaml
---
title: "Document-Driven Development: How We Ship 4,300 Lines in 3 Days"
date: 2025-08-15
author: Adam Stacoviak
slug: document-driven-development
summary: The future of software development isn't AI writing all the code. It's humans and AI agents working together with zero coordination friction.
---
```

### Content Ideas
- **AI-Augmented Development**: Our experience building with AI agents
- **Document-Driven Development**: The 7EP framework and coordination patterns
- **Technical Deep Dives**: Architecture decisions, performance improvements
- **Team Coordination**: Multi-agent development workflows
- **Project Evolution**: From basic tool to enterprise platform

### Writing Voice
- **Technical but accessible**: Engineers are the audience
- **Story-driven**: What we built, why, and what we learned  
- **Honest about challenges**: Real problems and real solutions
- **Future-focused**: Where AI-augmented development is heading

## 🚀 Deployment Details

### GitHub Pages Integration
- **Automatic**: Published content deploys in ~2-3 minutes
- **Free hosting**: GitHub Pages handles CDN, HTTPS, custom domains
- **Fast loading**: Static HTML, minimal CSS, zero JavaScript bloat
- **Professional URLs**: `https://[username].github.io/7zarch-go/` or custom domain

### Performance Targets
- **Page loads**: <100ms on slow connections
- **Build time**: All posts in <1 second  
- **CSS size**: <10KB total
- **Mobile-first**: Works beautifully on all devices

### Monitoring & Analytics
- **None**: No tracking, no analytics, respect for readers
- **GitHub insights**: Basic traffic data available in repo settings
- **Focus on content**: Let great writing speak for itself

## 🛠️ Technical Implementation Status

### Current State (Pre-Implementation)
- [ ] Blog directory structure
- [ ] Static site generator (`generator.go`)
- [ ] HTML templates
- [ ] CSS design system  
- [ ] GitHub Actions deployment
- [ ] Local development workflow

### Implementation Plan (7EP-0018)
1. **Week 1**: Core generator (markdown → HTML pipeline)
2. **Week 2**: Design implementation (typography, syntax highlighting)  
3. **Week 3**: Production polish (deployment, RSS, performance)

**Estimated Timeline**: 2-3 weeks for complete blog system

## 💡 Content Strategy

### Launch Content Ideas
1. **"Building 7zarch with AI Agents"** - Project origin story
2. **"Document-Driven Development Framework"** - Our coordination breakthrough
3. **"4,300 Lines in 3 Days"** - How AI agents and humans collaborate
4. **"The Future of Archive Management"** - Where we're heading

### Ongoing Content Themes
- **Technical Deep Dives**: Architecture, performance, design decisions
- **AI Development Insights**: What works, what doesn't, lessons learned  
- **Open Source Journey**: Building in public, community feedback
- **Industry Commentary**: Archive management, developer tools, AI trends

## 🎯 Success Metrics

### Content Quality
- **Readable by engineers**: Technical depth without complexity bloat
- **Shareable**: Content engineers want to share with their teams
- **Actionable**: Readers can apply insights to their own projects
- **Authentic**: Honest about challenges and failures, not just successes

### Technical Excellence  
- **Fast**: <100ms load times globally
- **Beautiful**: Code examples more readable than GitHub
- **Simple**: Single Go file generates entire blog
- **Maintainable**: Easy to modify, extend, and debug

### Strategic Impact
- **Project visibility**: Showcase 7zarch capabilities and philosophy
- **Community building**: Attract contributors and users
- **Industry influence**: Shape conversations about AI-augmented development
- **Team documentation**: Preserve and share our development insights

---

## 🚨 Emergency Procedures

### Blog Not Deploying
1. Check GitHub Actions: **Actions** tab → **Deploy Blog** workflow
2. Manual trigger: **Actions** → **Deploy Blog** → **Run workflow**
3. Check file paths: Ensure published files moved correctly
4. Contact CC: Technical deployment issues

### Wrong Content Published  
1. **Immediate**: Move back to drafts
2. **Quick fix**: Edit published file directly for typos
3. **Full revert**: Use git revert if major issues

### Generator Issues
1. **Local testing**: `cd blog && go run generator.go --help`
2. **Dependencies**: `go mod tidy` in blog directory
3. **Fallback**: Manual HTML generation if needed

---

**Remember**: The blog system is designed for safety first, speed second. You cannot accidentally publish anything. When in doubt, everything goes to `drafts/` first.

*Building a blog with soul, one intentional post at a time.*