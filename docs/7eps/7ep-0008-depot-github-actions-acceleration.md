# 7EP-0008: Depot GitHub Actions Acceleration

**Status**: Draft  
**Author**: Claude Code (CC)  
**Created**: 2025-08-12  
**Updated**: 2025-08-12  

## Summary

Implement Depot's managed GitHub Actions runners to accelerate CI/CD pipeline performance by 2-6x while reducing costs and improving developer productivity. This enhancement addresses current CI bottlenecks and provides foundation for faster feedback loops as the project scales.

## Background

### Current CI Performance Analysis

The 7zarch-go project utilizes an intensive CI pipeline with multiple workflows:

- **Multi-matrix testing**: 3 OS platforms × 3 Go versions = 9 parallel jobs
- **Cross-platform builds**: Linux, macOS, Windows with artifact generation
- **Comprehensive testing suite**: Unit, integration, MAS operations, edge cases, race detection
- **Multiple workflows**: CI, Test, Build, Quality, Release

### Performance Pain Points

1. **Slow job pickup**: GitHub-hosted runners typically take 10-30 seconds to start
2. **Limited cache performance**: 10GB cache limit with slow upload/download speeds
3. **Resource constraints**: Standard runners may be underpowered for compilation-heavy workloads
4. **Per-minute billing rounding**: 15-second job billed as full minute

## Goals

### Primary Objectives
- **Reduce CI pipeline duration by 50-70%**
- **Improve developer feedback loop speed**
- **Maintain or reduce CI costs despite performance gains**
- **Establish foundation for future scaling**

### Success Metrics
- CI job completion time (target: <3 minutes for standard workflows)
- Cache upload/download speeds (target: >100 MiB/s)
- Job pickup time (target: <5 seconds)
- Cost per CI minute (target: maintain current or lower)

## Design

### Implementation Strategy

#### Phase 1: Validation (Week 1)
- Start with most resource-intensive workflow (`test.yml`)
- Use 7-day free trial to validate performance gains
- Monitor and measure improvements

#### Phase 2: Full Migration (Week 2)
- Migrate remaining workflows if validation successful
- Configure organization-level settings
- Implement monitoring and alerting

### Technical Approach

#### Runner Selection
- **Primary**: `depot-ubuntu-22.04` for Linux workloads
- **macOS**: `depot-macos-13` for cross-platform compatibility
- **Windows**: `depot-windows-2022` for Windows builds
- **Fallback**: Keep existing GitHub runners as backup

#### Configuration Changes

```yaml
# Before (current)
runs-on: ${{ matrix.os }}
strategy:
  matrix:
    os: [ubuntu-latest, macos-latest, windows-latest]

# After (with Depot)
runs-on: ${{ matrix.os }}
strategy:
  matrix:
    os: [depot-ubuntu-22.04, depot-macos-13, depot-windows-2022]
```

#### Workflow-Specific Optimizations

**CI Workflow (`ci.yml`)**
```yaml
jobs:
  ci:
    runs-on: depot-ubuntu-22.04  # Single runner change
```

**Test Workflow (`test.yml`)**
```yaml
strategy:
  matrix:
    os: [depot-ubuntu-22.04, depot-macos-13]  # Depot runners
    go-version: ['1.21', '1.22', '1.23']
```

**Build Workflow (`build.yml`)**
```yaml
strategy:
  matrix:
    os: [depot-ubuntu-22.04, depot-macos-13, depot-windows-2022]
```

### Enhanced Features

#### Automatic Caching Acceleration
- **10x faster cache performance** with no configuration changes
- **Unlimited cache storage** (vs GitHub's 10GB limit)
- **Automatic cache orchestration** between jobs

#### Infrastructure Benefits
- **Single-tenant EC2 instances** for security and performance
- **3x faster disk access** via RAM disk acceleration
- **12.5 Gbps network throughput** for faster dependency downloads
- **Per-second billing** eliminates minute-rounding waste

## Implementation Plan

### Setup Requirements

#### 1. Account Configuration
```bash
# Organization setup at depot.dev
# - Create Depot organization account
# - Connect GitHub organization
# - Install Depot GitHub App
# - Configure runner permissions for public repos
```

#### 2. Workflow Updates
```yaml
# File: .github/workflows/test.yml (Phase 1 validation)
jobs:
  test:
    name: Test (${{ matrix.os }}, Go ${{ matrix.go-version }})
-   runs-on: ${{ matrix.os }}
+   runs-on: ${{ matrix.os }}
    strategy:
      matrix:
-       os: [ubuntu-latest, macos-latest]
+       os: [depot-ubuntu-22.04, depot-macos-13]
        go-version: ['1.21', '1.22', '1.23']
```

#### 3. Gradual Migration Path

**Week 1: Validation**
- Update `test.yml` only
- Monitor performance metrics
- Compare costs and timing
- Document results

**Week 2: Full Deployment**
- Update remaining workflows (`ci.yml`, `build.yml`, `quality.yml`)
- Configure retention and egress policies
- Set up monitoring dashboard
- Train team on new capabilities

### Configuration Management

#### Depot Organization Settings
```yaml
# Cache retention policy
cache_retention_days: 14  # Default, can extend to 30

# Egress filtering (if needed for security)
egress_rules:
  default: allow  # Start permissive, tighten as needed
  
# Runner settings
auto_auth_registry: true  # For Docker operations
containerd_layer_store: true  # Performance optimization
```

#### Rollback Strategy
```yaml
# Keep GitHub runners available for emergencies
# Can switch back by reverting runs-on values
runs-on: ${{ matrix.os }}
strategy:
  matrix:
    # Emergency fallback to GitHub runners
    os: [ubuntu-latest, macos-latest, windows-latest]
```

## Expected Impact

### Performance Improvements
- **Job startup time**: 10-30s → <5s (5-6x faster)
- **Cache operations**: ~10 MiB/s → >100 MiB/s (10x faster)
- **Overall pipeline**: 8-15 minutes → 3-5 minutes (2-3x faster)
- **Dependency downloads**: Faster via 12.5 Gbps network

### Cost Analysis
- **Per-second billing**: Eliminates waste from minute rounding
- **Higher performance**: Faster jobs = lower total runtime costs
- **Reduced waiting**: Faster feedback = higher developer productivity
- **Expected cost**: Neutral to 20% reduction despite premium infrastructure

### Developer Experience
- **Faster feedback loops**: Quicker PR validation
- **Reduced context switching**: Less waiting for CI results
- **Better reliability**: Enterprise-grade infrastructure
- **Scalability foundation**: Ready for team growth

## Risks and Mitigation

### Technical Risks
| Risk | Impact | Mitigation |
|------|--------|------------|
| Service outage | CI blocked | Keep GitHub runners as fallback |
| Performance regression | Slower than expected | 7-day trial allows validation |
| Configuration errors | Failed deployments | Gradual rollout with monitoring |

### Business Risks
| Risk | Impact | Mitigation |
|------|--------|------------|
| Higher costs | Budget overrun | Per-second billing often reduces costs |
| Vendor lock-in | Dependency concern | Standard GitHub Actions API, easy migration |
| Compliance issues | Security audit | Single-tenant infrastructure, egress controls |

### Operational Considerations
- **Team training**: New dashboard and monitoring tools
- **Support**: Depot provides dedicated support channels
- **Documentation**: Update CI/CD documentation for new capabilities

## Success Criteria

### Quantitative Metrics
- [ ] CI pipeline duration reduced by 50%+ (from ~10min to <5min)
- [ ] Job startup time under 5 seconds consistently
- [ ] Cache operations >100 MiB/s upload/download speeds
- [ ] Cost per CI run maintained or reduced

### Qualitative Goals
- [ ] Improved developer satisfaction with CI speed
- [ ] Reduced time to merge for PRs
- [ ] Foundation established for future scaling
- [ ] Team comfortable with new tooling

## Future Enhancements

### Phase 2 Possibilities
- **Custom runner configurations**: Tailor CPU/memory for specific workloads
- **Advanced caching strategies**: Cross-repository cache sharing
- **Container acceleration**: If Docker builds are added
- **Analytics dashboard**: Detailed performance monitoring

### Integration Opportunities
- **7EP-0002 CI/CD Enhancement**: Build on existing pipeline improvements
- **7EP-0005 Test Dataset System**: Faster test execution for large datasets
- **Future performance testing**: Rapid feedback for optimization work

## Conclusion

Implementing Depot GitHub Actions runners represents a strategic investment in developer productivity and CI infrastructure. With minimal configuration changes, the project can achieve significant performance improvements while establishing a foundation for future scaling.

The phased rollout approach minimizes risk while maximizing learning opportunities. The 7-day free trial provides a perfect validation period to confirm expected benefits before committing to the solution.

This enhancement aligns with 7zarch-go's focus on performance and developer experience, providing immediate benefits while positioning the project for continued growth and efficiency gains.

---

## Appendix A: Workflow File Changes

### test.yml Changes
```diff
 jobs:
   test:
     name: Test (${{ matrix.os }}, Go ${{ matrix.go-version }})
-    runs-on: ${{ matrix.os }}
+    runs-on: ${{ matrix.os }}
     strategy:
       matrix:
-        os: [ubuntu-latest, macos-latest]
+        os: [depot-ubuntu-22.04, depot-macos-13]
         go-version: ['1.21', '1.22', '1.23']
```

### ci.yml Changes
```diff
 jobs:
   ci:
-    runs-on: ubuntu-latest
+    runs-on: depot-ubuntu-22.04
```

### build.yml Changes
```diff
 strategy:
   matrix:
-    os: [ubuntu-latest, macos-latest, windows-latest]
+    os: [depot-ubuntu-22.04, depot-macos-13, depot-windows-2022]
```

## Appendix B: Performance Benchmarks

### Expected Improvements
| Metric | Current | With Depot | Improvement |
|--------|---------|------------|-------------|
| Job startup | 15-30s | <5s | 3-6x faster |
| Cache download | 5-10 MiB/s | 100+ MiB/s | 10-20x faster |
| Total pipeline | 8-15 min | 3-5 min | 2-3x faster |
| Cost per minute | GitHub pricing | 50% of GitHub | ~50% reduction |

### Real-World Examples
- **grpc/grpc**: 24m 18s → 3m 32s (6.9x improvement)
- **apache/kafka**: 12m 16s → 4m 42s (2.6x improvement)
- **zed-industries/zed**: 33m 28s → 24m 17s (1.4x improvement)