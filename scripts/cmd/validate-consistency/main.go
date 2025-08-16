package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ConsistencyChecker validates cross-document relationships in DDD framework
type ConsistencyChecker struct {
	RoleAssignments   map[string][]string // agent -> assignments
	NextCoordination  map[string]string   // agent -> status in NEXT.md
	SevenEPStatuses   map[string]string   // 7EP number -> status
	TeamContextRefs   []string            // files referencing team context
	ContentDuplicates map[string][]string // content hash -> file locations
}

func NewConsistencyChecker() *ConsistencyChecker {
	return &ConsistencyChecker{
		RoleAssignments:   make(map[string][]string),
		NextCoordination:  make(map[string]string),
		SevenEPStatuses:   make(map[string]string),
		TeamContextRefs:   []string{},
		ContentDuplicates: make(map[string][]string),
	}
}

func (cc *ConsistencyChecker) LoadFrameworkState(baseDir string) error {
	// Load role file assignments
	if err := cc.loadRoleAssignments(filepath.Join(baseDir, "docs/development/roles")); err != nil {
		return fmt.Errorf("loading role assignments: %w", err)
	}

	// Load NEXT.md coordination
	if err := cc.loadNextCoordination(filepath.Join(baseDir, "docs/development/NEXT.md")); err != nil {
		return fmt.Errorf("loading NEXT.md coordination: %w", err)
	}

	// Load 7EP statuses
	if err := cc.loadSevenEPStatuses(filepath.Join(baseDir, "docs/7eps")); err != nil {
		return fmt.Errorf("loading 7EP statuses: %w", err)
	}

	// Check team context references
	if err := cc.checkTeamContextReferences(baseDir); err != nil {
		return fmt.Errorf("checking team context references: %w", err)
	}

	return nil
}

func (cc *ConsistencyChecker) loadRoleAssignments(roleDir string) error {
	return filepath.Walk(roleDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || !strings.HasSuffix(path, ".md") {
			return err
		}

		filename := filepath.Base(path)
		if filename == "ROLE-TEMPLATE.md" || filename == "README.md" {
			return nil
		}

		agentName := strings.ToUpper(strings.TrimSuffix(filename, ".md"))
		
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		inActiveWork := false
		
		for scanner.Scan() {
			line := scanner.Text()
			
			if strings.Contains(line, "### Active Work") {
				inActiveWork = true
				continue
			}
			
			if inActiveWork && strings.HasPrefix(line, "###") {
				inActiveWork = false
			}
			
			if inActiveWork && strings.HasPrefix(line, "- **") {
				// Extract assignment name
				re := regexp.MustCompile(`- \*\*([^*]+)\*\* - ([A-Z]+)`)
				matches := re.FindStringSubmatch(line)
				if len(matches) >= 3 {
					assignment := matches[1]
					status := matches[2]
					cc.RoleAssignments[agentName] = append(cc.RoleAssignments[agentName], 
						fmt.Sprintf("%s (%s)", assignment, status))
				}
			}
		}

		return scanner.Err()
	})
}

func (cc *ConsistencyChecker) loadNextCoordination(nextPath string) error {
	file, err := os.Open(nextPath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inActiveWork := false
	
	for scanner.Scan() {
		line := scanner.Text()
		
		if strings.Contains(line, "Current Active Work") {
			inActiveWork = true
			continue
		}
		
		if inActiveWork && strings.HasPrefix(line, "##") {
			inActiveWork = false
		}
		
		if inActiveWork && strings.HasPrefix(line, "**") {
			// Extract agent status from NEXT.md - handle both formats
			re := regexp.MustCompile(`\*\*([A-Za-z-]+):\*\* (.+)`)
			matches := re.FindStringSubmatch(line)
			if len(matches) >= 3 {
				agent := strings.ToUpper(strings.Replace(matches[1], "-", "", -1))
				// Convert Amp-s to AMP, CC to CLAUDE, AC to AUGMENT, Adam to ADAM
				switch agent {
				case "AMPS":
					agent = "AMP"
				case "CC":
					agent = "CLAUDE" 
				case "AC":
					agent = "AUGMENT"
				}
				status := matches[2]
				cc.NextCoordination[agent] = status
			}
		}
	}

	return scanner.Err()
}

func (cc *ConsistencyChecker) loadSevenEPStatuses(sevenEPDir string) error {
	return filepath.Walk(sevenEPDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || !strings.HasSuffix(path, ".md") || info.Name() == "index.md" || info.Name() == "template.md" {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var epNumber, status string
		
		// Extract 7EP number from filename
		re := regexp.MustCompile(`7ep-(\d+)-`)
		matches := re.FindStringSubmatch(info.Name())
		if len(matches) >= 2 {
			epNumber = matches[1]
		}

		// Find status in file
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "**Status:**") {
				status = strings.TrimSpace(strings.TrimPrefix(line, "**Status:**"))
				break
			}
		}

		if epNumber != "" && status != "" {
			cc.SevenEPStatuses[epNumber] = status
		}

		return scanner.Err()
	})
}

func (cc *ConsistencyChecker) checkTeamContextReferences(baseDir string) error {
	// Check for team context content in wrong locations
	teamContextPattern := "Adam Stacoviak.*@adamstac.*Project owner"
	
	return filepath.Walk(filepath.Join(baseDir, "docs"), func(path string, info os.FileInfo, err error) error {
		if err != nil || !strings.HasSuffix(path, ".md") {
			return err
		}

		// Skip TEAM-CONTEXT.md itself
		if strings.Contains(path, "TEAM-CONTEXT.md") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if matched, _ := regexp.MatchString(teamContextPattern, string(content)); matched {
			cc.TeamContextRefs = append(cc.TeamContextRefs, path)
		}

		return nil
	})
}

func (cc *ConsistencyChecker) ValidateConsistency() []ValidationIssue {
	var issues []ValidationIssue

	// Check role assignments vs NEXT.md consistency
	for agent, assignments := range cc.RoleAssignments {
		nextStatus, exists := cc.NextCoordination[agent]
		if !exists {
			issues = append(issues, ValidationIssue{
				Type:        "coordination-mismatch",
				Severity:    "error",
				Description: fmt.Sprintf("Agent %s has role assignments but no NEXT.md coordination entry", agent),
				Files:       []string{fmt.Sprintf("docs/development/roles/%s.md", agent), "docs/development/NEXT.md"},
			})
			continue
		}

		// Basic consistency check - if agent has ACTIVE assignments, should not be "Available" in NEXT.md
		hasActiveWork := false
		for _, assignment := range assignments {
			if strings.Contains(assignment, "(ACTIVE)") {
				hasActiveWork = true
				break
			}
		}

		if hasActiveWork && strings.Contains(nextStatus, "Available") {
			issues = append(issues, ValidationIssue{
				Type:        "coordination-mismatch", 
				Severity:    "warning",
				Description: fmt.Sprintf("Agent %s has ACTIVE work but NEXT.md shows 'Available'", agent),
				Files:       []string{fmt.Sprintf("docs/development/roles/%s.md", agent), "docs/development/NEXT.md"},
			})
		}
	}

	// Check for team context duplication
	if len(cc.TeamContextRefs) > 0 {
		issues = append(issues, ValidationIssue{
			Type:        "content-boundary-violation",
			Severity:    "error", 
			Description: "Team context found in files other than TEAM-CONTEXT.md",
			Files:       cc.TeamContextRefs,
		})
	}

	return issues
}

type ValidationIssue struct {
	Type        string
	Severity    string
	Description string
	Files       []string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: validate-consistency <base-directory>")
		os.Exit(1)
	}

	baseDir := filepath.Clean(os.Args[1])
	// Convert to absolute path for security
	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: cannot resolve absolute path: %v\n", err)
		os.Exit(1)
	}
	baseDir = absBaseDir
	
	checker := NewConsistencyChecker()

	fmt.Println("🔍 DDD Framework Consistency Check")
	fmt.Println("==================================")

	// Load framework state
	if err := checker.LoadFrameworkState(baseDir); err != nil {
		log.Fatalf("Error loading framework state: %v", err)
	}

	// Validate consistency
	issues := checker.ValidateConsistency()

	// Report results
	errorCount := 0
	warningCount := 0

	for _, issue := range issues {
		switch issue.Severity {
		case "error":
			fmt.Printf("❌ %s: %s\n", strings.ToUpper(issue.Type), issue.Description)
			errorCount++
		case "warning":
			fmt.Printf("⚠️  %s: %s\n", strings.ToUpper(issue.Type), issue.Description)
			warningCount++
		}
		
		for _, file := range issue.Files {
			fmt.Printf("   📁 %s\n", file)
		}
		fmt.Println()
	}

	// Summary
	fmt.Printf("📊 Consistency Summary\n")
	fmt.Printf("=====================\n")
	fmt.Printf("Errors: %d\n", errorCount)
	fmt.Printf("Warnings: %d\n", warningCount)

	if errorCount == 0 {
		fmt.Println("✅ Framework consistency validation passed!")
		os.Exit(0)
	} else {
		fmt.Printf("❌ Framework consistency validation failed\n")
		os.Exit(1)
	}
}
