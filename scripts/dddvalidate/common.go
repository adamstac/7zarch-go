package dddvalidate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// ValidationError represents a validation failure
type ValidationError struct {
	File       string
	LineNumber int
	Message    string
	Severity   string
}

// ValidationWarning represents a validation warning
type ValidationWarning struct {
	File       string
	LineNumber int
	Message    string
}

// ValidationResult represents the outcome of validating a document
type ValidationResult struct {
	FilePath   string
	DocType    DocumentType
	Errors     []ValidationError
	Warnings   []ValidationWarning
	Compliance float64
}

// DocumentType represents different types of DDD framework documents
type DocumentType int

const (
	RoleFile DocumentType = iota
	WorkflowAction
	SevenEP
	TeamCoordination
	FrameworkDoc
)

func (dt DocumentType) String() string {
	switch dt {
	case RoleFile:
		return "Role File"
	case WorkflowAction:
		return "Workflow Action"
	case SevenEP:
		return "7EP Specification"
	case TeamCoordination:
		return "Team Coordination"
	case FrameworkDoc:
		return "Framework Documentation"
	default:
		return "Unknown"
	}
}

// DocumentParser provides common markdown parsing functionality
type DocumentParser struct {
	parser goldmark.Markdown
}

func NewDocumentParser() *DocumentParser {
	return &DocumentParser{
		parser: goldmark.New(),
	}
}

// ParseDocument parses markdown content into AST
func (dp *DocumentParser) ParseDocument(content []byte) ast.Node {
	reader := text.NewReader(content)
	return dp.parser.Parser().Parse(reader)
}

// ExtractHeaders extracts all headers from document with line numbers
func (dp *DocumentParser) ExtractHeaders(doc ast.Node, content []byte) map[string]HeaderInfo {
	headers := make(map[string]HeaderInfo)
	lines := strings.Split(string(content), "\n")
	
	ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		
		if heading, ok := node.(*ast.Heading); ok {
			line := heading.Segment.Start
			lineNum := 1
			for i, l := range lines {
				if i >= line {
					lineNum = i + 1
					break
				}
			}
			
			text := string(heading.Text(content))
			headers[text] = HeaderInfo{
				Level:      heading.Level,
				LineNumber: lineNum,
				Text:       text,
			}
		}
		
		return ast.WalkContinue, nil
	})
	
	return headers
}

type HeaderInfo struct {
	Level      int
	LineNumber int
	Text       string
}

// ExtractMetadataHeaders extracts YAML-style metadata from document header
func (dp *DocumentParser) ExtractMetadataHeaders(content []byte) map[string]string {
	metadata := make(map[string]string)
	lines := strings.Split(string(content), "\n")
	
	// Look for metadata in first 20 lines before first ## heading
	for i, line := range lines {
		if i > 20 || strings.HasPrefix(line, "## ") {
			break
		}
		
		// Match **Key:** value pattern
		metaPattern := regexp.MustCompile(`^\*\*([^*]+):\*\*\s*(.*)`)
		if matches := metaPattern.FindStringSubmatch(line); len(matches) >= 3 {
			key := strings.TrimSpace(matches[1])
			value := strings.TrimSpace(matches[2])
			metadata[key] = value
		}
	}
	
	return metadata
}

// ValidateDirectory walks directory and validates matching files
func ValidateDirectory(dir string, filePattern string, validator func(string, []byte) ValidationResult) ([]ValidationResult, error) {
	var results []ValidationResult
	
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Log but continue on file system errors
			fmt.Printf("Warning: cannot access %s: %v\n", path, err)
			return nil
		}
		
		if info.IsDir() {
			return nil
		}
		
		matched, err := filepath.Match(filePattern, filepath.Base(path))
		if err != nil {
			return err
		}
		
		if matched {
			content, err := os.ReadFile(path)
			if err != nil {
				fmt.Printf("Warning: cannot read %s: %v\n", path, err)
				return nil
			}
			
			result := validator(path, content)
			results = append(results, result)
		}
		
		return nil
	})
	
	return results, err
}

// StandardRegexPatterns provides common validation patterns
type StandardRegexPatterns struct {
	RequiredHeaders      []*regexp.Regexp
	ForbiddenContent     []*regexp.Regexp
	BranchNaming         *regexp.Regexp
	ConventionalCommit   *regexp.Regexp
	CoordinationCommit   *regexp.Regexp
	SessionLogFormat     *regexp.Regexp
}

func NewStandardPatterns() *StandardRegexPatterns {
	return &StandardRegexPatterns{
		RequiredHeaders: []*regexp.Regexp{
			regexp.MustCompile(`(?i)\*\*\s*Last\s+Updated\s*:\*\*`),
			regexp.MustCompile(`(?i)\*\*\s*Status\s*:\*\*`),
			regexp.MustCompile(`(?i)\*\*\s*Current\s+Focus\s*:\*\*`),
		},
		ForbiddenContent: []*regexp.Regexp{
			regexp.MustCompile(`(?i)Adam\s+Stacoviak.*@adamstac.*Project\s+owner`),
			regexp.MustCompile(`(?i)Human\s+Team.*Adam\s+Stacoviak`),
			regexp.MustCompile(`(?i)AI\s+Team.*AC.*Augment\s+Code.*You!`),
		},
		BranchNaming: regexp.MustCompile(`^(feature|feat|fix|docs|chore|refactor|hotfix)/.+$|^main$|^feat/7ep-\d+-.*$`),
		ConventionalCommit: regexp.MustCompile(`^(feat|fix|docs|style|refactor|test|chore)(\(.+\))?\s*:\s*.+`),
		CoordinationCommit: regexp.MustCompile(`(?i)(coordination:|session:|feat:.+coordination)`),
		SessionLogFormat: regexp.MustCompile(`(?i)^#\s+Session\s+Log\s+-\s+.+`),
	}
}

// ValidateHeaderFormats validates metadata header formats
func ValidateHeaderFormats(metadata map[string]string) []ValidationError {
	var errors []ValidationError
	
	// Validate Last Updated format
	if lastUpdated, exists := metadata["Last Updated"]; exists {
		// Accept YYYY-MM-DD or YYYY-MM-DD HH:MM
		datePattern := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}(\s+\d{2}:\d{2})?$`)
		if !datePattern.MatchString(lastUpdated) {
			errors = append(errors, ValidationError{
				LineNumber: 0,
				Message:    "Last Updated must be in YYYY-MM-DD or YYYY-MM-DD HH:MM format",
				Severity:   "error",
			})
		}
	}
	
	// Validate Status field values
	if status, exists := metadata["Status"]; exists {
		validStatuses := []string{"Available", "Active", "Blocked"}
		isValid := false
		for _, valid := range validStatuses {
			if strings.EqualFold(status, valid) {
				isValid = true
				break
			}
		}
		if !isValid {
			errors = append(errors, ValidationError{
				LineNumber: 0,
				Message:    "Status must be one of: Available, Active, Blocked",
				Severity:   "error",
			})
		}
	}
	
	return errors
}

// FilePathUtils provides safe file path operations
type FilePathUtils struct{}

func (fpu *FilePathUtils) IsRoleFile(path string) bool {
	basename := filepath.Base(path)
	return strings.HasSuffix(basename, ".md") &&
		!strings.Contains(basename, "ROLE-TEMPLATE") &&
		!strings.Contains(basename, "README") &&
		filepath.Dir(path) == "docs/development/roles"
}

func (fpu *FilePathUtils) SanitizePath(path string) (string, error) {
	// Clean and validate path to prevent injection
	cleaned := filepath.Clean(path)
	
	// Ensure path is within expected bounds (basic safety)
	if strings.Contains(cleaned, "..") {
		return "", fmt.Errorf("path traversal not allowed: %s", path)
	}
	
	return cleaned, nil
}

// OutputFormatter provides consistent formatting across tools
type OutputFormatter struct {
	UseColors bool
	JSONMode  bool
}

func NewOutputFormatter() *OutputFormatter {
	return &OutputFormatter{
		UseColors: true,
		JSONMode:  false,
	}
}

func (of *OutputFormatter) FormatError(err ValidationError) string {
	if of.JSONMode {
		return fmt.Sprintf(`{"type":"error","file":"%s","line":%d,"message":"%s"}`,
			err.File, err.LineNumber, err.Message)
	}
	
	if of.UseColors {
		if err.LineNumber > 0 {
			return fmt.Sprintf("  ❌ Line %d: %s", err.LineNumber, err.Message)
		}
		return fmt.Sprintf("  ❌ %s", err.Message)
	}
	
	return fmt.Sprintf("ERROR: %s", err.Message)
}

func (of *OutputFormatter) FormatWarning(warning ValidationWarning) string {
	if of.JSONMode {
		return fmt.Sprintf(`{"type":"warning","file":"%s","line":%d,"message":"%s"}`,
			warning.File, warning.LineNumber, warning.Message)
	}
	
	if of.UseColors {
		if warning.LineNumber > 0 {
			return fmt.Sprintf("  ⚠️  Line %d: %s", warning.LineNumber, warning.Message)
		}
		return fmt.Sprintf("  ⚠️  %s", warning.Message)
	}
	
	return fmt.Sprintf("WARNING: %s", warning.Message)
}

func (of *OutputFormatter) FormatSuccess(message string) string {
	if of.JSONMode {
		return fmt.Sprintf(`{"type":"success","message":"%s"}`, message)
	}
	
	if of.UseColors {
		return fmt.Sprintf("  ✅ %s", message)
	}
	
	return fmt.Sprintf("SUCCESS: %s", message)
}

// TimeUtils provides portable time operations
type TimeUtils struct{}

func (tu *TimeUtils) ParseSessionTime(timeStr string) (time.Time, error) {
	// Try common formats
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
	}
	
	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return t, nil
		}
	}
	
	return time.Time{}, fmt.Errorf("cannot parse time format: %s", timeStr)
}

func (tu *TimeUtils) FormatDuration(duration time.Duration) string {
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}
