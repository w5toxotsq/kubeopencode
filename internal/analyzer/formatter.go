package analyzer

import (
	"fmt"
	"strings"
)

// Severity represents the severity level of an analysis finding.
type Severity string

const (
	SeverityCritical Severity = "CRITICAL"
	SeverityWarning  Severity = "WARNING"
	SeverityInfo     Severity = "INFO"
)

// Finding represents a single issue or recommendation from the analysis.
type Finding struct {
	Severity    Severity
	Title       string
	Description string
	Suggestion  string
}

// AnalysisResult holds the structured output of a resource analysis.
type AnalysisResult struct {
	ResourceName string
	ResourceKind string
	Findings     []Finding
	Summary      string
}

// FormatResult formats an AnalysisResult for human-readable terminal output.
func FormatResult(result *AnalysisResult) string {
	var sb strings.Builder

	// Use a clearer separator with dashes for easier reading in the terminal
	sb.WriteString(fmt.Sprintf("\n--- Analysis: %s/%s ---\n", result.ResourceKind, result.ResourceName))

	if len(result.Findings) == 0 {
		sb.WriteString("No issues found.\n")
		return sb.String()
	}

	counts := map[Severity]int{}
	for _, f := range result.Findings {
		counts[f.Severity]++
	}

	sb.WriteString(fmt.Sprintf("Found %d finding(s): %d critical, %d warnings, %d info\n\n",
		len(result.Findings),
		counts[SeverityCritical],
		counts[SeverityWarning],
		counts[SeverityInfo],
	))

	for i, f := range result.Findings {
		sb.WriteString(fmt.Sprintf("[%d] [%s] %s\n", i+1, f.Severity, f.Title))
		if f.Description != "" {
			sb.WriteString(fmt.Sprintf("    Description: %s\n", f.Description))
		}
		if f.Suggestion != "" {
			sb.WriteString(fmt.Sprintf("    Suggestion:  %s\n", f.Suggestion))
		}
		sb.WriteString("\n")
	}

	if result.Summary != "" {
		sb.WriteString(fmt.Sprintf("Summary: %s\n", result.Summary))
	}

	return sb.String()
}

// ParseRawAnalysis attempts to extract structured findings from a raw LLM response.
// Falls back to wrapping the entire response as a single INFO finding.
func ParseRawAnalysis(resourceName, resourceKind, raw string) *AnalysisResult {
	result := &AnalysisResult{
		ResourceName: resourceName,
		ResourceKind: resourceKind,
	}

	lines := strings.Split(strings.TrimSpace(raw), "\n")
	var current *Finding

	for _, line := range lines {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "[CRITICAL]"):
			if current != nil {
				result.Findings = append(result.Findings, *current)
			}
			current = &Finding{Severity: SeverityCritical, Title: strings.TrimSpace(strings.TrimPrefix(line, "[CRITICAL]"))}
		case strings.HasPrefix(line, "[WARNING]"):
			if current != nil {
				result.Findings = append(result.Findings, *current)
			}
			current = &Finding{Severity: SeverityWarning, Title: strings.TrimSpace(strings.TrimPrefix(line, "[WARNING]"))}
		case strings.HasPrefix(line, "[INFO]"):
			if current != nil {
				result.Findings = append(result.Findings, *current)
			}
			current = &Finding{Severity: SeverityInfo, Title: strings.TrimSpace(strings.TrimPrefix(line, "[INFO]"))}
		case strings.HasPrefix(line, "Suggestion:") && current != nil:
			current.Suggestion = strings.TrimSpace(strings.TrimPrefix(line, "Suggestion:"))
		}
	}

	// Don't forget to append the last finding if present
	if current != nil {
		result.Findings = append(result.Findings, *current)
	}

	return result
}
