# Blog Tone & Style Guide

*How we talk about our work on 7zarch-go*

## Core Voice Principles

### 1. Technical Confidence Without Arrogance
- **Do**: "We shipped 4,338 lines in 72 hours. Here's how."
- **Don't**: "We're 10x developers changing the world."
- **Why**: Let the work speak for itself. Facts over hype.

### 2. Accessible to Everyone
- **Do**: "7-Zip with a brain - smart compression that just works"
- **Don't**: "Enterprise-grade solutions for digital transformation initiatives"
- **Why**: We're building for everyone who cares about their data, not just corporations.

### 3. Show the Process, Not Just Results
- **Do**: "We tried X, it failed. Then we discovered Y."
- **Don't**: "Everything worked perfectly the first time."
- **Why**: Real development is messy. Authenticity builds trust.

## Language Guidelines

### Words We Use
- **Smart** - Our tool makes intelligent decisions
- **Advanced** - Professional features accessible to all
- **Everyone/Anyone** - Inclusive, not elite
- **Ship/Shipped** - We deliver working code
- **Production-ready** - It works in the real world

### Words We Avoid
- **Enterprise** - We're not corporate software
- **Revolutionary/Game-changing** - Let others say this
- **Leverage/Synergy** - No business jargon
- **Best/Perfect** - Be specific about what's good
- **Obviously/Simply** - Nothing is obvious to everyone

### Personality Touches
OK to use sparingly for emphasis:
- "Pretty damn powerful"
- "Productive as hell"
- "That's the entire point"

Avoid:
- Excessive profanity
- Memes or dated references
- Forced humor

## Technical Writing

### Code Examples
```go
// Show real code from the project
// Include comments that explain why, not what
// Keep examples focused and runnable
```

### Metrics & Claims
- **Always specific**: "60-100μs search time" not "blazing fast"
- **Always verifiable**: Link to PRs, commits, benchmarks
- **Always honest**: "Usually 60μs, occasionally 100μs" 

### Structure Patterns

**Problem → Solution → Impact**
1. Here's what sucked
2. Here's what we built
3. Here's why it matters

**Technical Deep Dives**
1. Set context quickly
2. Jump to the technical meat
3. Show real implementation
4. Extract the lesson

## Philosophical Stance

### On AI Development
- AI agents are tools, not replacements
- Human-AI collaboration as peers
- Document-driven, not meeting-driven
- Ship code, not promises

### On Open Source
- Show everything - successes and failures
- Explain decisions, not just outcomes
- Credit influences and inspirations
- Help others implement our patterns

## Blog Post Checklist

Before publishing:
- [ ] Would a developer find this useful?
- [ ] Would a non-developer understand the impact?
- [ ] Are all claims backed by evidence?
- [ ] Does it sound like a human wrote it?
- [ ] Is it something we'd want to read?

## Example Opening Lines

**Good**:
- "You're using Claude, Cursor, or GitHub Copilot. Maybe all three. You're productive as hell for about 2 hours."
- "Something interesting happened today. Adam asked me to compare our DDD framework against Claude's defaults."
- "We solved this. Not with fancy memory systems or vector databases. With documents."

**Avoid**:
- "In today's fast-paced development environment..."
- "As we all know, AI is transforming software..."
- "Let me explain why our solution is revolutionary..."

## The One Rule

**Write like you're explaining to a smart friend over coffee, not presenting to a boardroom.**

---

*This guide itself follows these principles. Notice it's direct, specific, and includes real examples from our [actual blog posts](posts/).*