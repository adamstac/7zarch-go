// 7zarch Static Blog Generator
// Converts markdown posts to beautiful, fast static website
package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
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
	highlighting "github.com/yuin/goldmark-highlighting"
	"gopkg.in/yaml.v3"
)

// Command line flags for deployment strategy support
var (
	sourceDir = flag.String("source", "posts", "Source directory for markdown files")
	outputDir = flag.String("output", "public", "Output directory for generated files")
	watch     = flag.Bool("watch", false, "Watch for file changes and rebuild")
	verbose   = flag.Bool("verbose", false, "Verbose output")
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
	URL         string
}

// BlogConfig holds site-wide configuration
type BlogConfig struct {
	Title       string
	Description string
	BaseURL     string
	Author      string
	GitHubURL   string
}

func main() {
	flag.Parse()

	config := BlogConfig{
		Title:       "7zarch Blog",
		Description: "7-Zip with a brain deserves a blog with soul",
		BaseURL:     "https://adamstac.github.io/7zarch-go",
		Author:      "7zarch Team",
		GitHubURL:   "https://github.com/adamstac/7zarch-go",
	}

	fmt.Printf("🚀 Building blog from %s/ to %s/\n", *sourceDir, *outputDir)

	if *verbose {
		fmt.Printf("📁 Source: %s/\n", *sourceDir)
		fmt.Printf("📁 Output: %s/\n", *outputDir)
		fmt.Printf("⚙️  Watch mode: %t\n", *watch)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create output directory: %v", err))
	}

	// Build the blog
	buildBlog(config)

	if *watch {
		fmt.Println("👀 Watch mode not implemented yet - run manually for now")
	}

	fmt.Println("✅ Blog built successfully!")
}

func buildBlog(config BlogConfig) {
	start := time.Now()

	// 1. Load all markdown posts from source directory
	posts := loadPosts(*sourceDir + "/")
	fmt.Printf("📝 Found %d posts\n", len(posts))

	if len(posts) == 0 {
		fmt.Printf("⚠️  No posts found in %s/ - create some .md files to get started\n", *sourceDir)
		return
	}

	// 2. Sort by date (newest first)
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	// 3. Generate individual post pages
	for _, post := range posts {
		generatePostPage(post, config)
		if *verbose {
			fmt.Printf("  ✓ Generated %s.html\n", post.Slug)
		}
	}

	// 4. Generate index page
	generateIndexPage(posts, config)
	if *verbose {
		fmt.Println("  ✓ Generated index.html")
	}

	// 5. Generate RSS feed
	generateRSSFeed(posts, config)
	if *verbose {
		fmt.Println("  ✓ Generated feed.xml")
	}

	// 6. Copy static assets
	copyStaticAssets()

	duration := time.Since(start)
	fmt.Printf("⚡ Built %d posts in %v\n", len(posts), duration)
}

// loadPosts reads all .md files from the specified directory
func loadPosts(dir string) []Post {
	var posts []Post

	// Check if source directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("⚠️  Source directory %s doesn't exist - creating it\n", dir)
		if err := os.MkdirAll(dir, 0755); err != nil {
			panic(fmt.Sprintf("Failed to create source directory: %v", err))
		}
		return posts
	}

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %v", path, err)
		}

		post := parsePost(content, path)
		posts = append(posts, post)

		return nil
	})

	if err != nil {
		panic(fmt.Sprintf("Failed to load posts: %v", err))
	}

	return posts
}

// parsePost extracts frontmatter and converts markdown to HTML
func parsePost(content []byte, filename string) Post {
	// Split frontmatter from content
	parts := bytes.SplitN(content, []byte("---"), 3)
	if len(parts) < 3 {
		panic(fmt.Sprintf("Invalid frontmatter in %s - posts need YAML frontmatter", filename))
	}

	// Parse YAML frontmatter
	var post Post
	if err := yaml.Unmarshal(parts[1], &post); err != nil {
		panic(fmt.Sprintf("Failed to parse frontmatter in %s: %v", filename, err))
	}

	// Validate required fields
	if post.Title == "" {
		panic(fmt.Sprintf("Missing title in %s", filename))
	}
	if post.Slug == "" {
		panic(fmt.Sprintf("Missing slug in %s", filename))
	}
	if post.Date.IsZero() {
		panic(fmt.Sprintf("Missing or invalid date in %s", filename))
	}

	// Convert markdown content to HTML
	post.Content = processMarkdown(parts[2])
	post.ReadingTime = estimateReadingTime(string(parts[2]))
	post.Filename = filename
	post.URL = fmt.Sprintf("/%s.html", post.Slug)

	return post
}

// processMarkdown converts markdown to HTML with syntax highlighting
func processMarkdown(content []byte) template.HTML {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,        // GitHub Flavored Markdown
			extension.Typographer, // Smart quotes, dashes
			highlighting.NewHighlighting(
				highlighting.WithStyle("github"),
				highlighting.WithFormatOptions(
					html.WithLineNumbers(true),
					html.WithClasses(true),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(), // Auto-generate heading IDs
		),
	)

	var buf bytes.Buffer
	if err := md.Convert(content, &buf); err != nil {
		panic(fmt.Sprintf("Failed to convert markdown: %v", err))
	}

	return template.HTML(buf.String())
}

// joinURL safely joins base URL and path
func joinURL(baseURL, path string) string {
	if strings.HasSuffix(baseURL, "/") && strings.HasPrefix(path, "/") {
		return baseURL + path[1:]
	}
	if !strings.HasSuffix(baseURL, "/") && !strings.HasPrefix(path, "/") {
		return baseURL + "/" + path
	}
	return baseURL + path
}

// estimateReadingTime calculates reading time based on word count
func estimateReadingTime(content string) int {
	words := strings.Fields(content)
	minutes := len(words) / 200 // 200 words per minute
	if minutes == 0 {
		return 1
	}
	return minutes
}

// generatePostPage creates an individual HTML page for a post
func generatePostPage(post Post, config BlogConfig) {
	tmpl, err := template.ParseFiles("templates/post.html")
	if err != nil {
		panic(fmt.Sprintf("Failed to parse post template: %v", err))
	}

	filename := filepath.Join(*outputDir, fmt.Sprintf("%s.html", post.Slug))
	file, err := os.Create(filename)
	if err != nil {
		panic(fmt.Sprintf("Failed to create %s: %v", filename, err))
	}
	defer file.Close()

	data := struct {
		Post   Post
		Config BlogConfig
	}{
		Post:   post,
		Config: config,
	}

	if err := tmpl.Execute(file, data); err != nil {
		panic(fmt.Sprintf("Failed to execute post template: %v", err))
	}
}

// generateIndexPage creates the main blog index
func generateIndexPage(posts []Post, config BlogConfig) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		panic(fmt.Sprintf("Failed to parse index template: %v", err))
	}

	filename := filepath.Join(*outputDir, "index.html")
	file, err := os.Create(filename)
	if err != nil {
		panic(fmt.Sprintf("Failed to create %s: %v", filename, err))
	}
	defer file.Close()

	data := struct {
		Posts  []Post
		Config BlogConfig
	}{
		Posts:  posts,
		Config: config,
	}

	if err := tmpl.Execute(file, data); err != nil {
		panic(fmt.Sprintf("Failed to execute index template: %v", err))
	}
}

// generateRSSFeed creates an RSS feed for the blog
func generateRSSFeed(posts []Post, config BlogConfig) {
	// Create template with custom functions
	tmpl := template.New("rss.xml").Funcs(template.FuncMap{
		"joinURL": joinURL,
	})
	
	tmpl, err := tmpl.ParseFiles("templates/rss.xml")
	if err != nil {
		panic(fmt.Sprintf("Failed to parse RSS template: %v", err))
	}

	filename := filepath.Join(*outputDir, "feed.xml")
	file, err := os.Create(filename)
	if err != nil {
		panic(fmt.Sprintf("Failed to create %s: %v", filename, err))
	}
	defer file.Close()

	data := struct {
		Posts     []Post
		Config    BlogConfig
		BuildDate string
	}{
		Posts:     posts,
		Config:    config,
		BuildDate: time.Now().Format(time.RFC1123),
	}

	if err := tmpl.Execute(file, data); err != nil {
		panic(fmt.Sprintf("Failed to execute RSS template: %v", err))
	}
}

// copyStaticAssets copies CSS, fonts, images to output directory
func copyStaticAssets() {
	staticDir := "static"
	outputStaticDir := filepath.Join(*outputDir, "static")

	// Check if static directory exists
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		if *verbose {
			fmt.Println("📁 No static/ directory found - skipping asset copy")
		}
		return
	}

	// Ensure output static directory exists
	if err := os.MkdirAll(outputStaticDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create static output directory: %v", err))
	}

	// Copy all files from static/ to output/static/
	err := filepath.WalkDir(staticDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(staticDir, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(outputStaticDir, relPath)

		// Ensure destination directory exists
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}

		// Copy file
		src, err := os.Open(path)
		if err != nil {
			return err
		}
		defer src.Close()

		dst, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return err
		}

		if *verbose {
			fmt.Printf("  ✓ Copied %s\n", relPath)
		}

		return nil
	})

	if err != nil {
		panic(fmt.Sprintf("Failed to copy static assets: %v", err))
	}

	if *verbose {
		fmt.Println("  ✓ Static assets copied")
	}
}