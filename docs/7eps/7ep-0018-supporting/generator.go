// 7EP-0018 Supporting File: Blog Generator Implementation
// This file shows the core structure and key functions for the static blog generator

package main

import (
    "bytes"
    "fmt"
    "html/template"
    "io/fs"
    "os"
    "path/filepath"
    "sort"
    "strings"
    "time"
    
    "github.com/alecthomas/chroma/formatters/html"
    "github.com/yuin/goldmark"
    "github.com/yuin/goldmark/extension"
    "github.com/yuin/goldmark/parser"
    "github.com/yuin/goldmark-highlighting"
    "gopkg.in/yaml.v3"
)

// Post represents a blog post with frontmatter and content
type Post struct {
    Title       string    `yaml:"title"`
    Date        time.Time `yaml:"date"`
    Author      string    `yaml:"author"`
    Slug        string    `yaml:"slug"`
    Summary     string    `yaml:"summary"`
    Content     template.HTML
    ReadingTime int
    Filename    string
}

func main() {
    fmt.Println("Building 7zarch blog...")
    
    // 1. Load all markdown posts
    posts := loadPosts("posts/")
    fmt.Printf("Found %d posts\n", len(posts))
    
    // 2. Sort by date (newest first)
    sort.Slice(posts, func(i, j int) bool {
        return posts[i].Date.After(posts[j].Date)
    })
    
    // 3. Generate individual post pages
    for _, post := range posts {
        generatePostPage(post)
    }
    
    // 4. Generate index page
    generateIndexPage(posts)
    
    // 5. Generate RSS feed
    generateRSSFeed(posts)
    
    // 6. Copy static assets
    copyStaticAssets()
    
    fmt.Println("Blog built successfully!")
}

// loadPosts reads all .md files from the posts directory
func loadPosts(dir string) []Post {
    var posts []Post
    
    err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        
        if !strings.HasSuffix(path, ".md") {
            return nil
        }
        
        content, err := os.ReadFile(path)
        if err != nil {
            return err
        }
        
        post := parsePost(content, path)
        posts = append(posts, post)
        
        return nil
    })
    
    if err != nil {
        panic(err)
    }
    
    return posts
}

// parsePost extracts frontmatter and converts markdown to HTML
func parsePost(content []byte, filename string) Post {
    // Split frontmatter from content
    parts := bytes.SplitN(content, []byte("---"), 3)
    if len(parts) < 3 {
        panic("Invalid frontmatter in " + filename)
    }
    
    // Parse YAML frontmatter
    var post Post
    if err := yaml.Unmarshal(parts[1], &post); err != nil {
        panic(err)
    }
    
    // Convert markdown content to HTML
    post.Content = processMarkdown(parts[2])
    post.ReadingTime = estimateReadingTime(string(parts[2]))
    post.Filename = filename
    
    return post
}

// processMarkdown converts markdown to HTML with syntax highlighting
func processMarkdown(content []byte) template.HTML {
    md := goldmark.New(
        goldmark.WithExtensions(
            extension.GFM,                // GitHub Flavored Markdown
            extension.Typographer,       // Smart quotes, dashes
            highlighting.NewHighlighting(
                highlighting.WithStyle("github"),
                highlighting.WithFormatOptions(
                    html.WithLineNumbers(true),
                    html.WithClasses(true),
                ),
            ),
        ),
        goldmark.WithParserOptions(
            parser.WithAutoHeadingID(),  // Auto-generate heading IDs
        ),
    )
    
    var buf bytes.Buffer
    if err := md.Convert(content, &buf); err != nil {
        panic(err)
    }
    
    return template.HTML(buf.String())
}

// estimateReadingTime calculates reading time based on word count
func estimateReadingTime(content string) int {
    words := strings.Fields(content)
    return len(words) / 200 // 200 words per minute
}

// generatePostPage creates an individual HTML page for a post
func generatePostPage(post Post) {
    tmpl := template.Must(template.ParseFiles("templates/post.html"))
    
    filename := fmt.Sprintf("public/%s.html", post.Slug)
    file, err := os.Create(filename)
    if err != nil {
        panic(err)
    }
    defer file.Close()
    
    if err := tmpl.Execute(file, post); err != nil {
        panic(err)
    }
}

// generateIndexPage creates the main blog index
func generateIndexPage(posts []Post) {
    tmpl := template.Must(template.ParseFiles("templates/index.html"))
    
    file, err := os.Create("public/index.html")
    if err != nil {
        panic(err)
    }
    defer file.Close()
    
    data := struct {
        Posts []Post
        Title string
    }{
        Posts: posts,
        Title: "7zarch Blog",
    }
    
    if err := tmpl.Execute(file, data); err != nil {
        panic(err)
    }
}

// generateRSSFeed creates an RSS feed for the blog
func generateRSSFeed(posts []Post) {
    tmpl := template.Must(template.ParseFiles("templates/rss.xml"))
    
    file, err := os.Create("public/feed.xml")
    if err != nil {
        panic(err)
    }
    defer file.Close()
    
    data := struct {
        Posts []Post
        BuildDate string
    }{
        Posts: posts,
        BuildDate: time.Now().Format(time.RFC1123),
    }
    
    if err := tmpl.Execute(file, data); err != nil {
        panic(err)
    }
}

// copyStaticAssets copies CSS, fonts, images to public/
func copyStaticAssets() {
    // Copy static files (CSS, fonts, images)
    // Implementation would recursively copy static/ to public/static/
    fmt.Println("Copying static assets...")
}