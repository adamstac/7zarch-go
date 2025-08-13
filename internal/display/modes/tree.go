package modes

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/adamstac/7zarch-go/internal/display"
	"github.com/adamstac/7zarch-go/internal/storage"
)

// TreeDisplay provides hierarchical grouping of archives
type TreeDisplay struct{}

// NewTreeDisplay creates a new tree display mode
func NewTreeDisplay() *TreeDisplay {
	return &TreeDisplay{}
}

// Name returns the display mode name
func (td *TreeDisplay) Name() string {
	return "tree"
}

// MinWidth returns the minimum terminal width for this display
func (td *TreeDisplay) MinWidth() int {
	return 70
}

// Render displays archives in tree format
func (td *TreeDisplay) Render(archives []*storage.Archive, opts display.Options) error {
	if len(archives) == 0 {
		fmt.Printf("No archives found\n")
		fmt.Printf("Create archives with '7zarch-go create <path>' to see them here.\n")
		return nil
	}

	// Group archives by directory structure
	tree := td.buildDirectoryTree(archives)
	
	// Print summary header
	td.printSummary(archives)
	
	// Print the tree
	fmt.Printf("\nDirectory Structure:\n")
	td.printTree(tree, "", true, opts)

	return nil
}

// DirectoryNode represents a node in the directory tree
type DirectoryNode struct {
	Name     string
	Path     string
	Archives []*storage.Archive
	Children map[string]*DirectoryNode
	IsRoot   bool
}

// buildDirectoryTree creates a hierarchical structure from archives
func (td *TreeDisplay) buildDirectoryTree(archives []*storage.Archive) *DirectoryNode {
	root := &DirectoryNode{
		Name:     "Archives",
		Children: make(map[string]*DirectoryNode),
		IsRoot:   true,
	}

	// Group archives by their directory paths
	for _, archive := range archives {
		td.addToTree(root, archive)
	}

	return root
}

// addToTree adds an archive to the appropriate place in the tree
func (td *TreeDisplay) addToTree(root *DirectoryNode, archive *storage.Archive) {
	// Determine the grouping path
	var groupPath string
	
	if archive.Managed {
		// For managed archives, group by storage location
		groupPath = "Managed Storage"
	} else {
		// For external archives, group by directory
		dir := filepath.Dir(archive.Path)
		if dir == "." || dir == "/" {
			groupPath = "Root"
		} else {
			groupPath = dir
		}
	}

	// Navigate/create the path in the tree
	current := root
	parts := strings.Split(groupPath, string(filepath.Separator))
	
	for _, part := range parts {
		if part == "" || part == "." {
			continue
		}
		
		if _, exists := current.Children[part]; !exists {
			current.Children[part] = &DirectoryNode{
				Name:     part,
				Path:     part,
				Children: make(map[string]*DirectoryNode),
			}
		}
		current = current.Children[part]
	}
	
	// Add archive to the final node
	current.Archives = append(current.Archives, archive)
}

// printSummary prints the archive summary
func (td *TreeDisplay) printSummary(archives []*storage.Archive) {
	var managedCount, externalCount, missingCount, deletedCount int
	
	for _, a := range archives {
		if a.Status == "deleted" {
			deletedCount++
		} else if a.Managed {
			managedCount++
		} else {
			externalCount++
		}
		if a.Status == "missing" {
			missingCount++
		}
	}
	
	fmt.Printf("Archive Collection (%d archives found)\n", len(archives))
	fmt.Printf("Active: %d (Managed: %d, External: %d) | Missing: %d | Deleted: %d\n",
		managedCount+externalCount, managedCount, externalCount, missingCount, deletedCount)
}

// printTree recursively prints the directory tree
func (td *TreeDisplay) printTree(node *DirectoryNode, prefix string, isLast bool, opts display.Options) {
	if !node.IsRoot {
		// Print directory name
		connector := "├── "
		if isLast {
			connector = "└── "
		}
		
		fmt.Printf("%s%s%s", prefix, connector, node.Name)
		
		// Add archive count if any
		if len(node.Archives) > 0 {
			fmt.Printf(" (%d archives)", len(node.Archives))
		}
		fmt.Println()
		
		// Print archives in this directory
		if len(node.Archives) > 0 {
			td.printArchivesInDirectory(node.Archives, prefix, isLast, opts)
		}
	}
	
	// Sort children by name for consistent output
	var childNames []string
	for name := range node.Children {
		childNames = append(childNames, name)
	}
	sort.Strings(childNames)
	
	// Print child directories
	for i, childName := range childNames {
		child := node.Children[childName]
		childIsLast := i == len(childNames)-1
		
		var childPrefix string
		if node.IsRoot {
			childPrefix = ""
		} else if isLast {
			childPrefix = prefix + "    "
		} else {
			childPrefix = prefix + "│   "
		}
		
		td.printTree(child, childPrefix, childIsLast, opts)
	}
}

// printArchivesInDirectory prints archives within a directory node
func (td *TreeDisplay) printArchivesInDirectory(archives []*storage.Archive, prefix string, isLast bool, opts display.Options) {
	// Sort archives by name
	sort.Slice(archives, func(i, j int) bool {
		return archives[i].Name < archives[j].Name
	})
	
	for i, archive := range archives {
		archiveIsLast := i == len(archives)-1
		
		var archivePrefix string
		if isLast {
			archivePrefix = prefix + "    "
		} else {
			archivePrefix = prefix + "│   "
		}
		
		connector := "├── "
		if archiveIsLast {
			connector = "└── "
		}
		
		// Format archive entry
		status := td.formatTreeStatus(archive)
		size := display.FormatSize(archive.Size)
		age := td.formatTreeAge(archive.Created)
		
		fmt.Printf("%s%s📦 %s", archivePrefix, connector, archive.Name)
		
		if opts.Details {
			// Show detailed info
			id := archive.UID
			if len(id) > 12 {
				id = id[:12]
			}
			profile := archive.Profile
			if profile == "" {
				profile = "default"
			}
			
			fmt.Printf(" [%s] (%s, %s, %s, %s)", id, size, profile, age, status)
		} else {
			// Show basic info
			fmt.Printf(" (%s, %s, %s)", size, age, status)
		}
		
		fmt.Println()
	}
}

// formatTreeStatus returns a formatted status for tree display
func (td *TreeDisplay) formatTreeStatus(archive *storage.Archive) string {
	return display.FormatStatus(archive.Status, true) // Use icons for tree
}

// formatTreeAge formats duration since creation for tree display
func (td *TreeDisplay) formatTreeAge(created time.Time) string {
	age := time.Since(created)
	
	if age < time.Hour {
		mins := int(age.Minutes())
		return fmt.Sprintf("%dm", mins)
	}
	if age < 24*time.Hour {
		hours := int(age.Hours())
		return fmt.Sprintf("%dh", hours)
	}
	if age < 7*24*time.Hour {
		days := int(age.Hours() / 24)
		return fmt.Sprintf("%dd", days)
	}
	if age < 30*24*time.Hour {
		weeks := int(age.Hours() / (24 * 7))
		return fmt.Sprintf("%dw", weeks)
	}
	if age < 365*24*time.Hour {
		months := int(age.Hours() / (24 * 30))
		return fmt.Sprintf("%dmo", months)
	}
	years := int(age.Hours() / (24 * 365))
	return fmt.Sprintf("%dy", years)
}