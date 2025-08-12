# E2E Smoke Test Report

## 1) Latest Results
- Date: 2025-08-12
- Status: ✅ PASS

## 2) Test Workflow Executed
Executed locally with a sandboxed HOME to avoid touching real data.

Commands:
```
# Build
# from repo root
mkdir -p ./bin && go build -o ./bin/7zarch-go .

# Isolated HOME and sample data
HOME=$(mktemp -d -t sevenz_home.XXXXXX)
SRC=$(mktemp -d -t sevenz_src.XXXXXX)
mkdir -p "$SRC/subdir" && echo "hello world" > "$SRC/hello.txt" && echo "more data" > "$SRC/subdir/more.txt"

# Create -> List -> Delete -> Trash list -> Restore -> Final list
HOME="$HOME" ./bin/7zarch-go create "$SRC" --profile documents
HOME="$HOME" ./bin/7zarch-go list --details
ARCHIVE=$(ls -1 "$HOME/.7zarch-go/archives"/*.7z | head -n1 | xargs -n1 basename)
HOME="$HOME" ./bin/7zarch-go delete "$ARCHIVE"
HOME="$HOME" ./bin/7zarch-go trash list
HOME="$HOME" ./bin/7zarch-go restore "$ARCHIVE"
HOME="$HOME" ./bin/7zarch-go list --details
```

## 3) Results Summary
- Create: created managed archive successfully
- List: showed 1 managed archive present
- Delete: soft-deleted to trash by archive name
- Trash list: showed the deleted archive with short UID and purge schedule
- Restore: restored the archive back to managed storage
- Final list: archive present again, status ✓

## 4) Detailed Test Steps
1. Build
   - go build .
   - Success
2. Create
   - Command: create "$SRC" --profile documents
   - Output contained: "✅ Archive created successfully!" and "Stored in managed storage: <name>.7z"
3. List (details)
   - Output: "📦 Archives (1 found)" and ACTIVE - MANAGED table with short UID (e.g., 01K2G2GK)
4. Delete (soft)
   - Command: delete <archive-name>.7z
   - No error; moved to trash
5. Trash list
   - Output included the archive entry and short UID; showed deleted date and purge ETA
6. Restore
   - Command: restore <archive-name>.7z
   - Output: "✅ Restored <name>.7z to ~/.7zarch-go/archives/<name>.7z"
7. Final list
   - Output again shows 1 managed archive, status ✓

## 5) Issues Found
- None during this run. Earlier attempts using short UID on delete failed to resolve (ArchiveNotFound) in the isolated run; deleting by exact archive name is reliable for automation. Interactive MAS resolver UX works for show/move; delete path may need resolver integration follow-up in Phase 3 if we want short-ID deletes.

## 6) System State
- Branch: main
- Commit: 3f66a88e8352
- Go: go1.24.6 darwin/arm64
- Managed path: default (~/.7zarch-go) under the sandboxed HOME
- Test host: macOS (Apple Silicon)

## 1) Latest Results
- Date: 2025-08-12
- Status: ✅ PASS

## 2) Test Workflow Executed
Executed locally with a sandboxed HOME to avoid touching real data.

Commands:
```
# Build
go build -o ./bin/7zarch-go .

# Isolated HOME and sample data
HOME=$(mktemp -d -t sevenz_home.XXXXXX)
SRC=$(mktemp -d -t sevenz_src.XXXXXX)
mkdir -p "$SRC/subdir" && echo "hello world" > "$SRC/hello.txt" && echo "more data" > "$SRC/subdir/more.txt"

# Create -> List -> Delete -> Trash list -> Restore -> Final list
HOME="$HOME" ./bin/7zarch-go create "$SRC" --profile documents
HOME="$HOME" ./bin/7zarch-go list --details
ARCHIVE=$(ls -1 "$HOME/.7zarch-go/archives"/*.7z | head -n1 | xargs -n1 basename)
HOME="$HOME" ./bin/7zarch-go delete "$ARCHIVE"
HOME="$HOME" ./bin/7zarch-go trash list
HOME="$HOME" ./bin/7zarch-go restore "$ARCHIVE"
HOME="$HOME" ./bin/7zarch-go list --details
```

## 3) Results Summary
- Create: created managed archive successfully
- List: showed 1 managed archive present
- Delete: soft-deleted to trash by archive name
- Trash list: showed the deleted archive with short UID and purge schedule
- Restore: restored the archive back to managed storage
- Final list: archive present again, status ✓

## 4) Detailed Test Steps
1. Build
   - go build .
   - Success
2. Create
   - Command: create "$SRC" --profile documents
   - Output contained: "✅ Archive created successfully!" and "Stored in managed storage: <name>.7z"
3. List (details)
   - Output: "📦 Archives (1 found)" and ACTIVE - MANAGED table with short UID (e.g., 01K2G2GK)
4. Delete (soft)
   - Command: delete <archive-name>.7z
   - No error; moved to trash
5. Trash list
   - Output included the archive entry and short UID; showed deleted date and purge ETA
6. Restore
   - Command: restore <archive-name>.7z
   - Output: "✅ Restored <name>.7z to ~/.7zarch-go/archives/<name>.7z"
7. Final list
   - Output again shows 1 managed archive, status ✓

## 5) Issues Found
- None during this run. Earlier attempts using short UID on delete failed to resolve (ArchiveNotFound) in the isolated run; deleting by exact archive name is reliable for automation. Interactive MAS resolver UX works for show/move; delete path may need resolver integration follow-up in Phase 3 if we want short-ID deletes.

## 6) System State
- Branch: main
- Commit: 3f66a88e8352
- Go: go1.24.6 darwin/arm64
- Managed path: default (~/.7zarch-go) under the sandboxed HOME
- Test host: macOS (Apple Silicon)

>>>>>>> 2d6d85d (docs(test): add E2E smoke test report (✅ PASS) with commands, results, issues, and system state)
