# 7EP-0007: Enhanced MAS Operations

**Status:** 🎯 Ready for Implementation  
**Author(s):** Claude Code (CC), Augment Code (AC)  
**Assignment:** AC Lead (Query System, CLI Integration) + CC Support (Search Engine, Batch Core)  
**Difficulty:** 3 (moderate - builds on 7EP-0004 foundation)  
**Created:** 2025-08-12  
**Updated:** 2025-08-13 (Amp architectural review + 7EP-0014 foundation integration)  
**Foundation Status:** ✅ 7EP-0014 Complete - All dependencies satisfied  

## Executive Summary

Extend the MAS (Managed Archive Storage) foundation with advanced operations including batch processing, full-text search, saved queries, and enhanced workflow commands to provide a complete archive management experience.

## 🏛️ Amp (Sourcegraph) Architectural Review

**Overall Assessment**: ✅ **EXCELLENT DESIGN** - Well-architected, ready for implementation with 7EP-0014 foundation complete.

**Foundation Status**: 🎯 **PERFECT TIMING** - 7EP-0014 delivered all required foundation components:
- ✅ Database migration system enables safe schema evolution
- ✅ Machine-readable output (JSON/CSV) enables batch processing integration  
- ✅ Shell completion foundation provides baseline UX
- ✅ Complete trash lifecycle integrates with batch operations

**Strategic Impact**: This 7EP transforms 7zarch-go from **basic archive manager** → **power user command center** with enterprise-grade capabilities.

### 🎯 Implementation Priority Recommendations

**AC Implementation Focus** (User-Facing Power Features):
1. **Query System** - Highest ROI, enables all other features
2. **Search Integration** - Massive discoverability improvement
3. **Batch CLI Integration** - Workflow transformation
4. **Advanced List Enhancements** - Query save/load workflow

**CC Support Areas** (Infrastructure):
1. **Search Engine Performance** - Indexing optimization critical for UX
2. **Batch Processing Core** - Progress tracking and error handling
3. **Performance Validation** - Ensure benchmarks met under load

### 🚀 Key Architectural Strengths

**1. Excellent Layered Design**
- Query → Search → Batch → Integration layers well-separated
- Clean interfaces enable independent development and testing
- Builds perfectly on proven 7EP-0004 foundation patterns

**2. Performance-First Approach**  
- Realistic benchmarks (<100ms queries, <500ms search, <50ms completion)
- In-memory + persistent cache strategy for search index
- Progress tracking for batch operations ensures responsive UX

**3. User Experience Excellence**
- Saved queries address real power-user pain points
- Batch operations with confirmation/rollback provide safety
- Search across metadata fields enables content discovery

### ⚡ Implementation Acceleration Opportunities

**Leverage 7EP-0014 Foundation:**
- **Query Storage**: Use new migration system for schema changes
- **Machine Output**: Integrate with existing JSON/CSV for batch stdin
- **Error Patterns**: Follow established error handling from foundation
- **Testing**: Build on proven test patterns from migration system

## Evidence & Reasoning

**User feedback/pain points:**
- Need to perform operations on multiple archives at once (bulk move, delete, upload)
- Want to save complex filter combinations for repeated use
- Difficulty discovering archives by content or metadata beyond name/ID
- Shell completion missing for ULID prefixes and archive names
- Workflow gaps between basic operations (create → organize → upload cycles)

**Current limitations:**
- All operations work on single archives only
- No search beyond exact name/ID/prefix matching
- Filter combinations must be re-typed each session
- No shell integration for auto-completion
- Missing workflow automation for common patterns

**Why now:**
- 7EP-0004 MAS foundation provides stable base for advanced features
- User feedback indicates single-archive operations are insufficient
- Competition with other tools requires advanced discoverability features
- Shell completion expected by CLI power users

## Use Cases

### Primary Use Case: Batch Operations
```bash
# Select archives by filter and apply operations
7zarch-go list --profile=media --larger-than=100MB --save-query media-large
7zarch-go batch --query=media-large move --to=/backup/media/
7zarch-go batch --query=media-large upload --provider=s3

# Quick batch operations without saved queries
7zarch-go list --older-than=1y | 7zarch-go batch delete --confirm
```

### Secondary Use Cases

#### Saved Query Management
```bash
# Save complex filter combinations
7zarch-go query save "my-docs" --profile=documents --pattern="*-2024*" --managed
7zarch-go query list
7zarch-go query run my-docs

# Query composition
7zarch-go query save "large-old" --larger-than=1GB --older-than=6m
7zarch-go list --query=large-old --status=missing
```

#### Full-Text Search
```bash
# Search across all metadata fields
7zarch-go search "project backup 2024"
7zarch-go search --field=path "/Users/*/Documents"
7zarch-go search --field=name --regex ".*\.sql$"

# Combined with filters
7zarch-go search "important" --profile=documents --managed
```

#### Shell Completion & Workflow
```bash
# Tab completion for IDs and names
7zarch-go show <TAB>  # Shows ULID prefixes and archive names
7zarch-go move 01K2E3<TAB>  # Completes to full ULID

# Workflow automation
7zarch-go workflow create "daily-backup" \
  --steps="create,upload,cleanup" \
  --pattern="/home/user/projects/*" \
  --schedule="@daily"
```

## Technical Design

### Overview
Build advanced operations on top of the existing MAS foundation (7EP-0004) using a layered approach:
1. **Query Layer**: Saved queries and full-text search
2. **Batch Layer**: Multi-archive operations 
3. **Shell Layer**: Completion and workflow automation

### Component Architecture

#### 1. Query System (`internal/query/`)
```go
// Query management and execution
type Query struct {
    Name    string            `json:"name"`
    Filters map[string]string `json:"filters"`
    Created time.Time         `json:"created"`
    Used    int              `json:"used"`
}

type QueryManager struct {
    storage QueryStorage
    resolver *storage.Resolver
}

func (qm *QueryManager) Save(name string, filters ListFilters) error
func (qm *QueryManager) Run(name string) ([]*storage.Archive, error)
func (qm *QueryManager) List() ([]*Query, error)
```

#### 2. Search Engine (`internal/search/`)
```go
// Full-text search across archive metadata
type SearchIndex struct {
    registry *storage.Registry
    index    map[string][]string // term -> archive UIDs
}

func (si *SearchIndex) Search(query string) ([]*storage.Archive, error)
func (si *SearchIndex) SearchField(field, query string) ([]*storage.Archive, error)
func (si *SearchIndex) Reindex() error
```

#### 3. Batch Operations (`internal/batch/`)
```go
// Multi-archive operations with progress tracking
type BatchProcessor struct {
    manager *storage.Manager
    progress ProgressTracker
}

func (bp *BatchProcessor) Move(archives []*Archive, dest string) error
func (bp *BatchProcessor) Delete(archives []*Archive, confirm bool) error
func (bp *BatchProcessor) Upload(archives []*Archive, provider string) error
```

#### 4. Advanced Query Integration (`internal/query/`)
```go
// Advanced query operations building on 7EP-0014 foundation
type AdvancedQueryManager struct {
    basic    QueryManager      // From foundation
    search   *SearchEngine     // For query + search combinations
    batch    *BatchProcessor   // For query + batch operations
}

func (aqm *AdvancedQueryManager) SaveWithSearch(name string, searchTerms []string, filters ListFilters) error
func (aqm *AdvancedQueryManager) RunBatch(queryName string, operation BatchOperation) error
```

**Note**: Basic shell completion now provided by 7EP-0014 foundation. This 7EP adds advanced completion for saved queries.

### API Changes

#### New Commands
```bash
# Query management
7zarch-go query save <name> [filters...]
7zarch-go query list
7zarch-go query run <name>
7zarch-go query delete <name>

# Search operations  
7zarch-go search <query> [filters...]
7zarch-go reindex  # Rebuild search index

# Batch operations
7zarch-go batch <operation> [options...]
7zarch-go batch --query=<name> <operation>
7zarch-go batch --stdin <operation>  # Read UIDs from stdin

# Shell completion setup
7zarch-go completion bash|zsh|fish
```

#### Enhanced Existing Commands
```bash
# List command extensions (builds on 7EP-0014 machine output)
7zarch-go list --save-query=<name>    # Save current filters
7zarch-go list --query=<name>         # Use saved query
# Note: --output=json|csv already provided by 7EP-0014

# Show command extensions  
7zarch-go show --related              # Show similar archives
7zarch-go show --usage                # Show access history
# Note: --output=json already provided by 7EP-0014
```

### Data Model Changes

#### Query Storage (SQLite)
```sql
CREATE TABLE queries (
    name TEXT PRIMARY KEY,
    filters TEXT NOT NULL,  -- JSON-encoded filters
    created INTEGER NOT NULL,
    last_used INTEGER,
    use_count INTEGER DEFAULT 0
);
```

#### Search Index (In-Memory + Persistent Cache)
```sql
CREATE TABLE search_index (
    term TEXT,
    archive_uid TEXT,
    field TEXT,  -- name, path, etc.
    PRIMARY KEY (term, archive_uid, field)
);
```

## Implementation Plan

### Phase 1: Query Foundation (AC Lead)
- [ ] **Query Storage System** (AC)
  - [ ] SQLite schema for saved queries
  - [ ] Query CRUD operations
  - [ ] Query execution against existing filters
  - [ ] Query management commands (`query save/list/run/delete`)

- [ ] **List Command Integration** (AC)  
  - [ ] `--save-query` flag for saving current filters
  - [ ] `--query` flag for using saved queries
  - [ ] JSON/CSV output formats for scripting

### Phase 2: Search & Discovery (CC Lead)
- [ ] **Full-Text Search Engine** (CC)
  - [ ] Search index building and maintenance
  - [ ] Search query parsing and execution
  - [ ] Field-specific search capabilities
  - [ ] Search command implementation

- [ ] **Enhanced Show Command** (CC)
  - [ ] Related archive discovery
  - [ ] Usage history tracking
  - [ ] Metadata expansion

### Phase 3: Batch Operations (Shared)
- [ ] **Batch Processing Core** (CC)
  - [ ] Multi-archive operation framework
  - [ ] Progress tracking and reporting
  - [ ] Error handling and rollback

- [ ] **Batch Command Integration** (AC)
  - [ ] Batch command with query integration
  - [ ] Stdin processing for piped operations
  - [ ] Confirmation prompts for destructive operations

### Phase 4: Advanced Integration (AC Lead)
- [ ] **Query + Search Integration** (AC)
  - [ ] Combine saved queries with search terms
  - [ ] Search result saving as queries
  - [ ] Complex query composition workflows
  - [ ] Advanced completion for saved query names (builds on 7EP-0014 foundation)

### Dependencies
- **7EP-0004**: MAS Foundation ✅ (completed) - provides resolver, registry, and basic commands  
- **7EP-0001**: Trash Management ✅ (completed via 7EP-0014) - batch delete integrates with trash
- **7EP-0014**: Critical Foundation Gaps ✅ (completed) - provides machine output, database migrations, shell completion baseline

## Testing Strategy

### Acceptance Criteria
- [ ] Can save and reuse complex filter combinations
- [ ] Full-text search finds archives by any metadata field
- [ ] Batch operations work on 100+ archive sets with progress tracking
- [ ] Shell completion works for ULID prefixes and archive names
- [ ] All operations complete in <5s for typical registries (<10K archives)
- [ ] Query system handles 100+ saved queries efficiently

### Test Scenarios

#### Query System Testing
- Query saving with various filter combinations
- Query execution with registry changes
- Query management (list, delete, rename)
- Query performance with large registries

#### Search Testing
- Full-text search across all metadata fields
- Field-specific searches with regex support
- Search performance with 10K+ archives
- Index rebuilding and consistency

#### Batch Operations Testing  
- Large batch operations (1000+ archives)
- Mixed operation types with error handling
- Progress tracking accuracy
- Rollback on partial failures

#### Shell Integration Testing
- Completion accuracy across different shells
- Performance of completion queries
- Integration with existing shell environments

### Performance Benchmarks
- **Query execution**: <100ms for complex queries on 10K archives
- **Search operations**: <500ms for full-text search on 10K archives
- **Batch operations**: Progress updates every 1-2 seconds for large sets
- **Completion queries**: <50ms for ULID prefix completion

## Migration/Compatibility

### Breaking Changes
None - all new functionality building on existing commands.

### Upgrade Path
- Existing commands continue working unchanged
- New features opt-in through new flags and commands
- Query system starts empty, users build saved queries over time

### Backward Compatibility
All existing 7EP-0004 functionality preserved exactly.

## Alternatives Considered

**External search tools**: Considered integrating with `fzf` or `ripgrep` but decided native search provides better integration and doesn't require external dependencies.

**Query language**: Evaluated SQL-like syntax but decided flag-based queries are more CLI-native and easier to save/compose.

**Batch processing via shell pipes**: Considered Unix-style piping (`7zarch-go list | xargs 7zarch-go delete`) but decided native batch operations provide better error handling and progress tracking.

## AC/CC Implementation Split

### AC (Augment Code) Responsibilities - User-Facing Features
- **Query Management System**: Storage, CRUD operations, query execution
- **List Command Enhancement**: Save/load queries, output formats  
- **Batch Command Integration**: Query integration, confirmation flows
- **User Experience**: Command interfaces, help text, error messages
- **CLI Integration**: Flag design, command composition, workflow patterns

### CC (Claude Code) Responsibilities - Infrastructure & Performance
- **Search Engine**: Full-text indexing, search execution, reindexing
- **Batch Processing Core**: Multi-archive operations, progress tracking
- **Shell Completion**: Completion provider, shell script generation
- **Performance Optimization**: Search indexing, batch operation efficiency
- **Testing Infrastructure**: Benchmarks, performance tests, edge case coverage

### Shared Responsibilities
- **API Design**: Command interfaces and flag naming (AC leads, CC reviews)
- **Integration Testing**: Cross-component workflow validation
- **Error Handling**: Consistent error patterns across components
- **Documentation**: User guides (AC), technical architecture (CC)

### Coordination Points
1. **Query Filter Integration**: How saved queries map to existing ListFilters (AC designs, CC implements backend)
2. **Batch Operation Interface**: How batch processor integrates with existing commands (AC designs CLI, CC implements engine)
3. **Search Index Schema**: What metadata fields to index and how (CC designs, AC provides user requirements)
4. **Completion Data Source**: How completion provider accesses registry efficiently (CC implements, AC defines user experience)

### Communication Protocol
- **Weekly Planning**: AC and CC coordinate feature priorities and dependencies
- **PR Cross-Review**: AC reviews CC infrastructure PRs, CC reviews AC user experience PRs
- **Integration Points**: Dedicated coordination sessions when components need to integrate
- **User Feedback Loop**: AC gathers and prioritizes user needs, CC ensures technical feasibility

## Future Considerations

- **Query Sharing**: Export/import saved queries between users
- **Advanced Search**: Fuzzy matching, similarity scoring, machine learning
- **Workflow Automation**: Scheduled operations, trigger-based actions
- **Web Interface**: Browser-based archive management dashboard
- **API Server**: REST API for external integrations

## 🎯 AC Implementation Guidance (Amp Strategic Direction)

### **Phase 1: Query Foundation - START HERE** (Estimated: 3-4 days)

**HIGHEST ROI** - Enables all subsequent features. Focus on user-facing query workflows first.

#### Quick Wins (Day 1-2):
```bash
# Core query commands to implement
7zarch-go query save "my-docs" --profile=documents --managed
7zarch-go query list
7zarch-go query run my-docs
7zarch-go query delete my-docs

# List integration 
7zarch-go list --save-query=my-docs    # Save current filters as query
7zarch-go list --query=my-docs         # Use saved query
```

**Implementation Strategy:**
1. **Use 7EP-0014 migration system** for query table schema
2. **Leverage existing ListFilters** - serialize to JSON for storage
3. **Build on machine output** - queries can export JSON for automation
4. **Follow resolver patterns** - query names resolve like archive IDs

#### Technical Foundation:
```go
// Query storage building on 7EP-0014 migration foundation
type QueryStorage struct {
    db *sql.DB  // Use existing registry database
}

// Migration file: 0004_query_system.sql
CREATE TABLE queries (
    name TEXT PRIMARY KEY,
    filters TEXT NOT NULL,  -- JSON-encoded ListFilters
    created INTEGER NOT NULL,
    last_used INTEGER,
    use_count INTEGER DEFAULT 0
);
```

### **Phase 2: Search Engine - HIGH IMPACT** (Estimated: 2-3 days)

**MASSIVE DISCOVERABILITY** - Users can find archives by any metadata field.

```bash
# Search implementation priority
7zarch-go search "project backup 2024"    # Cross-field search
7zarch-go search --name "backup.*2024"     # Field-specific with regex
7zarch-go reindex                          # Rebuild search index
```

**Performance Critical:**
- **<500ms search target** for 10K archives requires optimized indexing
- **In-memory index** with persistent cache for performance
- **Incremental updates** when archives added/modified

### **Phase 3: Batch Operations - WORKFLOW TRANSFORMATION** (Estimated: 3-4 days)

**POWER USER ENABLEMENT** - Multi-archive operations with safety and progress.

```bash
# Batch implementation priority  
7zarch-go batch --query=media-large move --to=/backup/media/
7zarch-go batch --stdin delete --confirm   # Pipe from list output
7zarch-go list --older-than=1y --output=json | jq -r '.[].uid' | 7zarch-go batch delete --confirm
```

**Safety Requirements:**
- **Confirmation prompts** for destructive operations
- **Progress tracking** every 1-2 seconds for large sets
- **Rollback capability** on partial failures
- **Integration with trash system** from 7EP-0014

### **Phase 4: Advanced Integration - POLISH** (Estimated: 1-2 days)

**PROFESSIONAL FINISH** - Query + search combinations, advanced workflows.

```bash
# Advanced combinations
7zarch-go search "important" --save-query=important-files
7zarch-go query run important-files --output=json | jq -r '.[].uid' | 7zarch-go batch upload --provider=s3
```

### 🛠️ AC Technical Implementation Notes

**Database Integration:**
- **Use 7EP-0014 migration system** - create migration files, don't modify schema directly
- **Follow established patterns** - error handling, ULID resolution, registry operations
- **Performance validation** - test with realistic datasets (1000+ archives)

**CLI Design:**
- **Consistent flag naming** - follow patterns from existing commands
- **Machine output integration** - ensure all new commands support `--output json`
- **Help text clarity** - comprehensive examples and error guidance

**Testing Strategy:**
- **Unit tests** for all query, search, and batch operations
- **Integration tests** combining query + search + batch workflows
- **Performance tests** validating benchmark targets
- **CLI workflow tests** ensuring smooth user experience

### 🚨 Critical Implementation Decisions for AC

**1. Query Storage Schema** - Use JSON encoding of ListFilters for maximum compatibility
**2. Search Index Strategy** - Start with in-memory, add persistence for performance
**3. Batch Operation Safety** - Always require confirmation for destructive ops
**4. CLI Verb Structure** - Keep `query`, `search`, `batch` as top-level commands

### 📊 Success Metrics

**User Experience:**
- Query save/load workflow reduces repeated typing by 80%
- Search enables finding archives without knowing exact names
- Batch operations handle 100+ archives efficiently with progress feedback

**Performance:**
- Query execution <100ms for complex queries on 10K archives  
- Search operations <500ms for full-text search on 10K archives
- Batch operations with progress updates every 1-2 seconds

**Adoption:**
- Power users can accomplish multi-step workflows in single commands
- External tools can integrate via machine-readable output
- CLI feels professional and efficient compared to basic file managers

## References

- **Builds on**: 7EP-0004 MAS Foundation Implementation, 7EP-0014 Critical Foundation Gaps ✅
- **Integrates with**: 7EP-0001 Trash Management System ✅ 
- **Enables**: Advanced archive discovery and bulk management workflows
- **Related**: CLI completion patterns from tools like `kubectl`, `docker`, `git`