---
title: "Test Post: 7zarch Blog Generator"
date: 2025-08-15
author: Claude Code (CC)
slug: test-post
summary: Testing the newly implemented static blog generator with proper YAML frontmatter and markdown content.
---

# Testing Our Blog Generator

This is a test post to verify that our **7zarch blog generator** is working correctly.

## Features Being Tested

- YAML frontmatter parsing
- Markdown to HTML conversion
- Syntax highlighting for code blocks
- Template rendering

## Code Example

Here's some Go code to test syntax highlighting:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, 7zarch blog!")
    
    // Test different syntax elements
    posts := []string{"post1", "post2", "post3"}
    for i, post := range posts {
        fmt.Printf("%d: %s\n", i+1, post)
    }
}
```

## Lists and Formatting

- **Bold text** works
- *Italic text* works  
- `Inline code` works
- Links to [7zarch GitHub](https://github.com/adamstac/7zarch-go) work

1. Ordered lists
2. Also work properly
3. With proper styling

## Conclusion

If you can read this formatted properly, our blog generator is working! 🎉

This test verifies:
- ✅ Frontmatter parsing
- ✅ Markdown processing  
- ✅ Template rendering
- ✅ Static file generation