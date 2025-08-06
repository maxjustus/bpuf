# Communication and Development Guidelines

## Communication Guidelines

When working on projects, use measured, specific language:

**Avoid overly confident terms**:
- Don't use: "comprehensive", "production ready", "robust", "enterprise-grade", "bulletproof", "seamless", "cutting-edge"
- Instead use: specific descriptions of what the code does

**Use factual, measured language**:
- "implements X pattern" rather than "elegantly implements"
- "handles Y scenario" rather than "comprehensively handles" 
- "supports Z formats" rather than "full support for Z formats"
- "processes N objects" rather than "efficiently processes large volumes"

**Focus on specifics**:
- Include actual numbers, concrete capabilities, and factual descriptions
- Describe what exists rather than aspirational qualities
- Use precise technical terms rather than marketing language

**Documentation style**:
- No emojis in README files or documentation
- Use plain text checkmarks and formatting instead of emoji indicators
- Remember: "I'm not trying to prove to the world that I'm competent"
- Avoid sections that show off (architecture diagrams, test coverage lists, etc.)
- Focus on what users need, not demonstrating technical prowess

## Git Commit Guidelines

**Commit Message Format**: Use standard, concise commit messages without any tool attribution:

```bash
git commit -m "feat: add pattern matching for selective deletion"
git commit -m "fix: resolve compatibility issue"
git commit -m "test: add unit tests"
```

**Important**: Do NOT include Claude Code attribution, co-author tags, or any generated-by notices in commit messages. Use clean, professional commit messages that focus on the actual changes made.

## Development Process

### Incremental Development Rules

1. **Small Changes Only**: Each change should be under 50 lines. If bigger, break it down.
2. **Test Every Change**: Run the full validation cycle after every code modification.
3. **No Exceptions**: Every warning must be fixed. No "I'll fix it later."

### Required Validation Cycle

Run this after every code change, in order:

```bash
go fmt ./...       # Format code
golangci-lint run  # Lint and catch issues
go test ./...      # Run all tests
```

All must pass with zero warnings/errors before proceeding.

### Git Workflow

```bash
# After validation cycle passes:
git add .
git commit -m "short: what changed"
```

Commit criteria:
- **Logical unit of work complete** (function, struct, test, etc.)
- **All validation passes**
- **Code actually works**

### Code Style Preferences

- **Simple over clever**: Readable code beats impressive code
- **Explicit over implicit**: Clear intent over brevity
- **Boring is good**: Standard patterns, no custom solutions unless necessary
- **Flat structure**: Avoid deep nesting, prefer early returns
- **Small functions**: 20-30 lines max, single responsibility

### Example Development Flow

```bash
# 1. Add basic struct
# Edit: Add Config struct with 3 fields
[validation cycle]
git add . && git commit -m "config: add basic Config struct"

# 2. Add validation  
# Edit: Add Config.Validate() method
[validation cycle]
git add . && git commit -m "config: add validation method"

# 3. Add tests
# Edit: Add 5 unit tests for Config
[validation cycle]
git add . && git commit -m "config: add unit tests"
```

## README Guidelines

Write focused, user-oriented documentation:

**Include**:
- Brief description of what the tool does
- Installation instructions
- Configuration examples
- Usage examples
- Development setup (for contributors)
- License information

**Avoid**:
- Architecture diagrams
- Test coverage statistics
- Performance benchmarks
- Technical implementation details
- Sections that demonstrate competence rather than provide utility

**Approach**: Show, don't tell. Each feature includes example commands. Designed for users who want to understand quickly.

## Testing Strategy

- **Unit Tests**: Core logic and individual components
- **Integration Tests**: End-to-end operations with real services
- **Keep tests focused**: Test behavior, not implementation details
- **Avoid test coverage boasting**: Tests exist for quality, not metrics