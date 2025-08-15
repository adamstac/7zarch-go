package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// GitPatternValidator validates git history and state against DDD patterns
type GitPatternValidator struct {
	SessionLogPattern    *regexp.Regexp
	CoordinationCommits  *regexp.Regexp
	BranchNamingPattern  *regexp.Regexp
	SessionCommitPattern *regexp.Regexp
}

func NewGitPatternValidator() *GitPatternValidator {
	return &GitPatternValidator{
		SessionLogPattern:    regexp.MustCompile(`(?i)^# Session Log - .+`),
		CoordinationCommits:  regexp.MustCompile(`(?i)(coordination:|feat:.+coordination|session:)`),
		BranchNamingPattern:  regexp.MustCompile(`^(feature|feat|fix|docs|chore|refactor|hotfix)/.+$|^main$|^feat/7ep-\d+-.*$`),
		SessionCommitPattern: regexp.MustCompile(`(?i)session: (start|end)`),
	}
}

type GitValidationResult struct {
	SessionLogs      []SessionLogValidation
	CommitPatterns   []CommitValidation
	BranchCompliance BranchValidation
	Errors          []string
	Warnings        []string
}

type SessionLogValidation struct {
	FilePath   string
	HasStart   bool
	HasEnd     bool
	ValidTiming bool
	Errors     []string
}

type CommitValidation struct {
	Hash       string
	Message    string
	Compliant  bool
	HasCoord   bool
	Error      string
}

type BranchValidation struct {
	Current     string
	Compliant   bool
	Error       string
}

func (gpv *GitPatternValidator) ValidateRepository(baseDir string) GitValidationResult {
	result := GitValidationResult{
		SessionLogs:    []SessionLogValidation{},
		CommitPatterns: []CommitValidation{},
		Errors:        []string{},
		Warnings:      []string{},
	}

	// Validate session logs
	sessionLogs := gpv.findSessionLogs(baseDir)
	for _, logPath := range sessionLogs {
		logValidation := gpv.validateSessionLog(logPath)
		result.SessionLogs = append(result.SessionLogs, logValidation)
	}

	// Validate recent commit patterns
	commits := gpv.getRecentCommits(baseDir, 20)
	for _, commit := range commits {
		commitValidation := gpv.validateCommitMessage(commit)
		result.CommitPatterns = append(result.CommitPatterns, commitValidation)
	}

	// Validate current branch naming
	result.BranchCompliance = gpv.validateCurrentBranch(baseDir)

	return result
}

func (gpv *GitPatternValidator) findSessionLogs(baseDir string) []string {
	var sessionLogs []string
	
	logsDir := filepath.Join(baseDir, "docs/logs")
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		return sessionLogs
	}

	filepath.Walk(logsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if strings.HasSuffix(path, ".md") && strings.Contains(path, "session-") {
			sessionLogs = append(sessionLogs, path)
		}
		
		return nil
	})

	return sessionLogs
}

func (gpv *GitPatternValidator) validateSessionLog(logPath string) SessionLogValidation {
	validation := SessionLogValidation{
		FilePath: logPath,
		Errors:   []string{},
	}

	file, err := os.Open(logPath)
	if err != nil {
		validation.Errors = append(validation.Errors, fmt.Sprintf("Cannot read session log: %v", err))
		return validation
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var content strings.Builder
	
	for scanner.Scan() {
		line := scanner.Text()
		content.WriteString(line + "\n")
		
		// Check for session start indicators
		if strings.Contains(line, "Session started by DDD Framework") {
			validation.HasStart = true
		}
		
		// Check for session end indicators
		if strings.Contains(line, "Session completed by DDD Framework") {
			validation.HasEnd = true
		}
	}

	contentStr := content.String()

	// Validate session log header format
	if !gpv.SessionLogPattern.MatchString(contentStr) {
		validation.Errors = append(validation.Errors, "Session log missing standard header format")
	}

	// Check for required timing sections
	if !strings.Contains(contentStr, "Start Time:") {
		validation.Errors = append(validation.Errors, "Missing session start time")
	}

	// Validate timing format if both start and end present
	if validation.HasStart && validation.HasEnd {
		validation.ValidTiming = strings.Contains(contentStr, "Duration:")
		if !validation.ValidTiming {
			validation.Errors = append(validation.Errors, "Missing or invalid session duration calculation")
		}
	}

	return validation
}

func (gpv *GitPatternValidator) getRecentCommits(baseDir string, count int) []GitCommit {
	cmd := exec.Command("git", "log", fmt.Sprintf("-%d", count), "--oneline", "--format=%H|%s")
	cmd.Dir = baseDir
	
	output, err := cmd.Output()
	if err != nil {
		// Properly handle git command failures
		return []GitCommit{}
	}

	var commits []GitCommit
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "|", 2)
		if len(parts) == 2 {
			commits = append(commits, GitCommit{
				Hash:    parts[0],
				Message: parts[1],
			})
		}
	}

	return commits
}

type GitCommit struct {
	Hash    string
	Message string
}

func (gpv *GitPatternValidator) validateCommitMessage(commit GitCommit) CommitValidation {
	validation := CommitValidation{
		Hash:    commit.Hash,
		Message: commit.Message,
	}

	// Check for coordination patterns
	validation.HasCoord = gpv.CoordinationCommits.MatchString(commit.Message)

	// Validate commit message format (conventional commits style) - improved pattern
	conventionalPattern := regexp.MustCompile(`^(feat|fix|docs|style|refactor|test|chore)(\(.+\))?\s*:\s*.+`)
	validation.Compliant = conventionalPattern.MatchString(commit.Message)

	if !validation.Compliant {
		validation.Error = "Commit message doesn't follow conventional format: type(scope): description"
	}

	// Additional checks for coordination commits
	if validation.HasCoord {
		if !strings.Contains(commit.Message, "coordination:") && 
		   !strings.Contains(commit.Message, "session:") {
			validation.Error = "Coordination commit missing coordination context in message"
		}
	}

	return validation
}

func (gpv *GitPatternValidator) validateCurrentBranch(baseDir string) BranchValidation {
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = baseDir
	
	output, err := cmd.Output()
	if err != nil {
		return BranchValidation{
			Current:   "unknown",
			Compliant: false,
			Error:     fmt.Sprintf("Cannot determine current branch: %v", err),
		}
	}

	currentBranch := strings.TrimSpace(string(output))
	
	validation := BranchValidation{
		Current: currentBranch,
	}

	// Allow main branch or properly named feature branches
	if currentBranch == "main" {
		validation.Compliant = true
	} else if gpv.BranchNamingPattern.MatchString(currentBranch) {
		validation.Compliant = true
	} else {
		validation.Compliant = false
		validation.Error = "Branch name doesn't follow pattern: feature/description or feat/7ep-XXXX-name"
	}

	return validation
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: validate-git-patterns <base-directory>")
		os.Exit(1)
	}

	baseDir := os.Args[1]
	validator := NewGitPatternValidator()

	fmt.Println("🔍 DDD Git Pattern Validation")
	fmt.Println("=============================")
	fmt.Printf("Generated: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	result := validator.ValidateRepository(baseDir)

	// Report session log validation
	fmt.Println("📋 Session Log Validation")
	if len(result.SessionLogs) == 0 {
		fmt.Println("  ⚠️  No session logs found in docs/logs/")
	} else {
		for _, sessionLog := range result.SessionLogs {
			filename := filepath.Base(sessionLog.FilePath)
			fmt.Printf("  📁 %s\n", filename)
			
			if len(sessionLog.Errors) == 0 {
				fmt.Printf("    ✅ Valid session log format\n")
			} else {
				for _, err := range sessionLog.Errors {
					fmt.Printf("    ❌ %s\n", err)
				}
			}
		}
	}
	fmt.Println()

	// Report commit pattern validation
	fmt.Println("📋 Recent Commit Pattern Validation")
	coordinationCommits := 0
	compliantCommits := 0
	
	for _, commit := range result.CommitPatterns {
		if commit.HasCoord {
			coordinationCommits++
		}
		if commit.Compliant {
			compliantCommits++
		}
		
		if commit.Error != "" {
			fmt.Printf("  ❌ %s: %s\n", commit.Hash[:8], commit.Error)
		}
	}
	
	fmt.Printf("  📊 Compliant commits: %d/%d (%.1f%%)\n", 
		compliantCommits, len(result.CommitPatterns), 
		float64(compliantCommits)/float64(len(result.CommitPatterns))*100)
	fmt.Printf("  🤝 Coordination commits: %d (%.1f%%)\n", 
		coordinationCommits,
		float64(coordinationCommits)/float64(len(result.CommitPatterns))*100)
	fmt.Println()

	// Report branch validation
	fmt.Println("📋 Branch Naming Validation")
	if result.BranchCompliance.Compliant {
		fmt.Printf("  ✅ Current branch '%s' follows naming convention\n", result.BranchCompliance.Current)
	} else {
		fmt.Printf("  ❌ Current branch '%s': %s\n", result.BranchCompliance.Current, result.BranchCompliance.Error)
	}
	fmt.Println()

	// Overall summary
	totalErrors := len(result.Errors)
	for _, sessionLog := range result.SessionLogs {
		totalErrors += len(sessionLog.Errors)
	}
	for _, commit := range result.CommitPatterns {
		if commit.Error != "" {
			totalErrors++
		}
	}
	if !result.BranchCompliance.Compliant {
		totalErrors++
	}

	fmt.Println("📊 Git Pattern Summary")
	fmt.Println("=====================")
	fmt.Printf("Total errors: %d\n", totalErrors)

	if totalErrors == 0 {
		fmt.Println("✅ All git patterns comply with DDD framework!")
		os.Exit(0)
	} else {
		fmt.Printf("❌ Git pattern validation failed\n")
		os.Exit(1)
	}
}
