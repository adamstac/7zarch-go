# E2E TUI Actions Smoke Test

Date: 2025-08-13
Related PR: #21

Scope
- Validate TUI actions overlay, confirm dialog, delete/restore operations
- Verify error overlay and status hints display

Preconditions
- Repo built using safe 3-line build (or new build system)
- Sample managed archives exist (at least one present, one deleted in trash)

Steps
1) Launch TUI
   - ./dist/7zarch-go tui
2) Select a present archive
   - Move cursor to an item
   - Press Space to select (✓ indicator appears)
3) Delete (soft)
   - Press 'a' → select "Delete" → Enter
   - Confirm with Enter (or 'y')
   - Expect: item status changes to deleted after refresh; status shows "delete: N item(s)"
4) Restore
   - Filter to deleted (f until deleted)
   - Select one or more deleted items
   - 'a' → "Restore" → Enter → confirm
   - Expect: items move back to present; status shows "restore: N item(s)"
5) Errors overlay
   - Induce an error (e.g., lock permissions on file or move across restricted path) if feasible
   - After action, check status shows "errors=N (e)"
   - Press 'e' to open errors overlay; Esc to close

Notes
- Actions apply to selected rows; if none selected, to current row
- 'r' refreshes; '?' shows help; 'q' quits

Results
- [ ] Delete succeeded
- [ ] Restore succeeded
- [ ] Errors collected and shown
- [ ] Status hints displayed correctly

Follow-ups
- Consider Clear errors key/action
- Wire Mark Uploaded, Move actions with same patterns

