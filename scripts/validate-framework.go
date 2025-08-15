package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// DocumentType represents different types of DDD framework documents
type DocumentType int

const (
	RoleFile DocumentType = iota
	WorkflowAction
	SevenEP
	TeamCoordination
	FrameworkDoc
)

// ValidationResult represents the outcome of validating a document
type ValidationResult struct {
	FilePath   string
	DocType    DocumentType
	Errors     []ValidationError
	Warnings   []ValidationWarning
	Compliance float64
}

type ValidationError struct {
	LineNumber int
	Message    string
	Severity   string
}

type ValidationWarning struct {
	LineNumber int
	Message    string
}

// RoleFileValidator validates role files against 7EP-0019 standards
type RoleFileValidator struct {
	RequiredHeaders []string
	RequiredSections []string
	ForbiddenContent []string
}

func NewRoleFileValidator() *RoleFileValidator {
	return &RoleFileValidator{
		RequiredHeaders: []string{
			"Last Updated:",
			"Status:",
			"Current Focus:",
		},
		RequiredSections: []string{
			"🎯 Current Assignments",
			"🔗 Coordination Needed", 
			"✅ Recently Completed",
			"📝 Implementation Notes",
		},
		ForbiddenContent: []string{
			"Adam Stacoviak.*@adamstac.*Project owner",
			"Human Team",
			"AI Team.*AC.*Augment Code.*You!",
		},
	}
}

func (rv *RoleFileValidator) ValidateDocument(filePath string, content []byte) ValidationResult {
	result := ValidationResult{
		FilePath: filePath,
		DocType:  RoleFile,
		Errors:   []ValidationError{},
		Warnings: []ValidationWarning{},
	}

	// Parse markdown
	reader := text.NewReader(content)
	parser := goldmark.New()
	doc := parser.Parser().Parse(reader)

	// Check required headers
	headerSection := extractHeaderSection(doc, content)
	for _, header := range rv.RequiredHeaders {
		if !strings.Contains(headerSection, header) {
			result.Errors = append(result.Errors, ValidationError{
				LineNumber: 1,
				Message:    fmt.Sprintf("Missing required header field: %s", header),
				Severity:   "error",
			})
		}
	}

	// Validate header field formats
	if err := rv.validateHeaderFormats(headerSection); err != nil {
		result.Errors = append(result.Errors, ValidationError{
			LineNumber: 1,
			Message:    err.Error(),
			Severity:   "error",
		})
	}

	// Check required sections
	sections := extractSections(doc, content)
	for _, requiredSection := range rv.RequiredSections {
		if !containsSection(sections, requiredSection) {
			result.Errors = append(result.Errors, ValidationError{
				LineNumber: 0,
				Message:    fmt.Sprintf("Missing required section: %s", requiredSection),
				Severity:   "error",
			})
		}
	}

	// Check for forbidden content (content boundary violations)
	contentStr := string(content)
	for _, forbidden := range rv.ForbiddenContent {
		if matched, _ := regexp.MatchString(forbidden, contentStr); matched {
			result.Errors = append(result.Errors, ValidationError{
				LineNumber: 0,
				Message:    fmt.Sprintf("Content boundary violation: contains team context (should reference TEAM-CONTEXT.md)"),
				Severity:   "error",
			})
		}
	}

	// Check for TEAM-CONTEXT.md reference (except ADAM.md)
	if !strings.Contains(filepath.Base(filePath), "ADAM") {
		if !strings.Contains(contentStr, "TEAM-CONTEXT.md") {
			result.Warnings = append(result.Warnings, ValidationWarning{
				LineNumber: 0,
				Message:    "Recommended: Add reference to TEAM-CONTEXT.md",
			})
		}
	}

	// Calculate compliance score
	totalChecks := len(rv.RequiredHeaders) + len(rv.RequiredSections) + len(rv.ForbiddenContent)
	passedChecks := totalChecks - len(result.Errors)
	result.Compliance = float64(passedChecks) / float64(totalChecks) * 100

	return result
}

func (rv *RoleFileValidator) validateHeaderFormats(headerSection string) error {
	// Validate Last Updated format (YYYY-MM-DD or YYYY-MM-DD HH:MM)
	datePattern := `\*\*Last Updated:\*\*\s+\d{4}-\d{2}-\d{2}(\s+\d{2}:\d{2})?`
	if matched, _ := regexp.MatchString(datePattern, headerSection); !matched {
		return fmt.Errorf("Last Updated field must be in YYYY-MM-DD or YYYY-MM-DD HH:MM format")
	}

	// Validate Status field values
	statusPattern := `\*\*Status:\*\*\s+(Available|Active|Blocked)`
	if matched, _ := regexp.MatchString(statusPattern, headerSection); !matched {
		return fmt.Errorf("Status field must be one of: Available, Active, Blocked")
	}

	return nil
}

// Helper functions for markdown parsing
func extractHeaderSection(doc ast.Node, content []byte) string {
	// Extract everything before first ## heading
	lines := strings.Split(string(content), "\n")
	var headerLines []string
	
	for _, line := range lines {
		if strings.HasPrefix(line, "## ") {
			break
		}
		headerLines = append(headerLines, line)
	}
	
	return strings.Join(headerLines, "\n")
}

func extractSections(doc ast.Node, content []byte) []string {
	var sections []string
	lines := strings.Split(string(content), "\n")
	
	for _, line := range lines {
		if strings.HasPrefix(line, "## ") {
			sections = append(sections, strings.TrimSpace(line))
		}
	}
	
	return sections
}

func containsSection(sections []string, required string) bool {
	for _, section := range sections {
		if strings.Contains(section, required) {
			return true
		}
	}
	return false
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: validate-framework <directory>")
		os.Exit(1)
	}

	targetDir := filepath.Clean(os.Args[1])
	// Convert to absolute path for security
	absTargetDir, err := filepath.Abs(targetDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: cannot resolve absolute path: %v\n", err)
		os.Exit(1)
	}
	targetDir = absTargetDir
	
	validator := NewRoleFileValidator()
	
	var allResults []ValidationResult
	totalErrors := 0
	totalFiles := 0

	// Validate all role files
	roleDir := filepath.Join(targetDir, "docs/development/roles")
	walkErr := filepath.Walk(roleDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Warning: cannot access %s: %v\n", path, err)
			return nil // Continue on file system errors
		}

		if !strings.HasSuffix(path, ".md") || 
		   strings.Contains(path, "ROLE-TEMPLATE.md") || 
		   strings.Contains(path, "README.md") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Error reading %s: %v", path, err)
			return nil
		}

		result := validator.ValidateDocument(path, content)
		allResults = append(allResults, result)
		totalErrors += len(result.Errors)
		totalFiles++

		return nil
	})

	if walkErr != nil {
		log.Fatalf("Error walking directory: %v", walkErr)
	}

	// Output results
	fmt.Println("🔍 DDD Framework Validation Report")
	fmt.Println("==================================")
	fmt.Printf("Generated: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	for _, result := range allResults {
		filename := filepath.Base(result.FilePath)
		fmt.Printf("📋 %s (%.1f%% compliant)\n", filename, result.Compliance)
		
		for _, err := range result.Errors {
			if err.LineNumber > 0 {
				fmt.Printf("  ❌ Line %d: %s\n", err.LineNumber, err.Message)
			} else {
				fmt.Printf("  ❌ %s\n", err.Message)
			}
		}
		
		for _, warning := range result.Warnings {
			if warning.LineNumber > 0 {
				fmt.Printf("  ⚠️  Line %d: %s\n", warning.LineNumber, warning.Message)
			} else {
				fmt.Printf("  ⚠️  %s\n", warning.Message)
			}
		}
		
		if len(result.Errors) == 0 && len(result.Warnings) == 0 {
			fmt.Printf("  ✅ Fully compliant\n")
		}
		fmt.Println()
	}

	// Summary
	fmt.Println("📊 Validation Summary")
	fmt.Println("====================")
	fmt.Printf("Files validated: %d\n", totalFiles)
	fmt.Printf("Total errors: %d\n", totalErrors)

	if totalErrors == 0 {
		fmt.Println("✅ All documents pass DDD framework validation!")
		os.Exit(0)
	} else {
		fmt.Printf("❌ Framework validation failed with %d errors\n", totalErrors)
		os.Exit(1)
	}
}
