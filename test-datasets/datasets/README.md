# Generated Datasets Directory

This directory is used for caching generated test datasets when needed. All contents except this README and .gitignore are excluded from version control.

## Purpose

While the test dataset system primarily works with metadata-only archives (no actual files), this directory can be used for:

- Caching generated test metadata for large scenarios
- Storing benchmark results for performance tracking
- Temporary test output during development

## Usage

The directory is automatically managed by the test system. No manual intervention required.

## Cleanup

Test data in this directory is ephemeral and can be safely deleted at any time:

```bash
# Clean all generated test data
rm -rf test-datasets/datasets/*
```

The test system will regenerate any needed data on next run.