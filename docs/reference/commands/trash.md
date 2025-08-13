# Trash Management (Coming Soon)

Note: This feature is being implemented in PR #10 and is not yet available on main.

## restore
Restore a deleted archive back to its original location (or managed path).

Usage:

```
7zarch-go restore <id> [--force] [--dry-run]
```

- --force: overwrite destination if it exists
- --dry-run: show actions without making changes

## trash list
List deleted archives with retention countdown.

Usage:
```
7zarch-go trash list [--within-days N] [--before YYYY-MM-DD] [--json]
```

- --within-days: only show items purging within N days
- --before: only show items deleted before date
- --json: machine-readable output

## trash purge
Permanently delete trashed archives that are eligible for purge.

Usage:
```
7zarch-go trash purge [--all] [--within-days N] [--force] [--dry-run]
```

- --all: ignore retention and purge everything
- --within-days: only purge items purging within N days
- --force: skip confirmation prompts
- --dry-run: show actions without making changes

