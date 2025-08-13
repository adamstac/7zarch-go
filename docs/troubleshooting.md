# 7zarch-go Troubleshooting Guide

## Common Issues and Solutions

### Archive Not Found Errors

**Problem**: `Archive 'xyz' not found` when trying to show, delete, or move an archive.

**Solutions**:
1. Check available archives: `7zarch-go list`
2. Use a longer ID prefix if multiple archives match
3. Try using the full ULID instead of a partial match
4. Verify the archive hasn't been deleted: `7zarch-go list --deleted`

### Multiple Archives Match

**Problem**: `Archive ID 'abc' matches multiple archives`

**Solutions**:
1. Use a longer prefix (e.g., 'abcd' instead of 'abc')
2. Use the full ULID shown in the error message
3. Use `7zarch-go list` to see all archives and their full IDs

### Permission Denied Errors

**Problem**: Cannot create or access archives due to permission errors.

**Solutions**:
1. Check file permissions: `ls -la <path>`
2. Ensure you have write access to the managed storage directory
3. For external archives, verify read permissions on the source
4. Try using `sudo` if appropriate (not recommended for regular use)

### Missing Archives

**Problem**: Archives show as "missing" in the list.

**Solutions**:
1. Check if the file was moved: `7zarch-go list --missing`
2. Verify the original path still exists
3. Use `7zarch-go move` to update the path if relocated
4. Consider using managed storage to avoid path issues

### Build/Installation Issues

**Problem**: `make dev` or build commands fail.

**Solutions**:
1. Ensure Go 1.24+ is installed: `go version`
2. Install Goreleaser: `brew install goreleaser`
3. Clean and rebuild: `make clean && make dev`
4. Check for missing dependencies: `make deps`

### Database Errors

**Problem**: Database-related errors when running commands.

**Solutions**:
1. Check database status: `7zarch-go db status`
2. Run migrations if needed: `7zarch-go db migrate`
3. Create a backup: `7zarch-go db backup`
4. Reset database (last resort): `rm ~/.7zarch-go/registry.db`

### Compression Issues

**Problem**: Archives are too large or compression is slow.

**Solutions**:
1. Use appropriate profile: `--profile media` for photos/videos
2. Use `--profile documents` for text/office files
3. Check available disk space before creating large archives
4. Use `--dry-run` to preview without creating

### Shell Completion Not Working

**Problem**: Tab completion doesn't work in terminal.

**Solutions**:
1. Generate completion script: `7zarch-go completion <shell>`
2. Source the completion file in your shell config
3. Restart your terminal or reload shell config
4. Verify completion is enabled in your shell

## Debug Mode

For any issue, enable debug mode to see detailed information:

```bash
# Show performance metrics and detailed output
7zarch-go list --debug
7zarch-go show <id> --debug
```

Debug output includes:
- Query execution time
- Memory usage
- Database size
- Result counts
- Error stack traces

## Getting Help

1. **Check command help**: `7zarch-go help <command>`
2. **View examples**: `7zarch-go <command> --help`
3. **Read documentation**: `/docs/reference/`
4. **Report issues**: https://github.com/adamstac/7zarch-go/issues

## Performance Tips

1. **Use managed storage** for better performance and tracking
2. **Filter large lists** with `--pattern` or `--status` flags
3. **Use appropriate display modes** for your terminal size
4. **Enable debug mode** to identify bottlenecks
5. **Keep database optimized** with regular maintenance

## Best Practices

1. **Regular backups**: Use `7zarch-go db backup` periodically
2. **Consistent naming**: Use descriptive archive names
3. **Profile selection**: Choose the right compression profile
4. **Managed storage**: Prefer managed over external storage
5. **Version control**: Track important archives with metadata