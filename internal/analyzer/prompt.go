package analyzer

import (
	"fmt"
	"strings"
)

// systemPrompt is the base system prompt for the AI analyzer.
const systemPrompt = `You are an expert Kubernetes administrator and security engineer.
Your task is to analyze Kubernetes resource configurations and provide actionable insights.

For each resource, you should:
1. Identify potential security vulnerabilities or misconfigurations
2. Highlight performance or reliability concerns
3. Check for missing best practices (resource limits, health checks, etc.)
4. Suggest specific improvements with corrected YAML examples where applicable

Be concise but thorough. Format your response with clear sections:
- **Summary**: Brief overview of findings
- **Issues Found**: List of problems with severity (Critical/High/Medium/Low/Info)
- **Recommendations**: Specific actionable steps to remediate
- **Best Practices**: Any additional suggestions for improvement`

// buildAnalysisPrompt constructs the prompt for analyzing a Kubernetes resource.
func buildAnalysisPrompt(resourceType, resourceName, namespace, resourceJSON string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Analyze the following Kubernetes %s resource:\n\n", resourceType))

	if namespace != "" {
		sb.WriteString(fmt.Sprintf("Name: %s\nNamespace: %s\n\n", resourceName, namespace))
	} else {
		sb.WriteString(fmt.Sprintf("Name: %s\n\n", resourceName))
	}

	sb.WriteString("Resource Configuration (JSON):\n```json\n")
	sb.WriteString(resourceJSON)
	sb.WriteString("\n```\n\n")
	// Personal preference: also ask for a one-line TL;DR at the top so I can
	// skim results quickly when reviewing a lot of resources.
	sb.WriteString("Start with a single TL;DR sentence, then provide a detailed analysis of this resource configuration, focusing on security, reliability, and best practices.")

	return sb.String()
}

// buildMultiResourcePrompt constructs a prompt for analyzing multiple resources at once.
func buildMultiResourcePrompt(resources []resourceEntry) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Analyze the following %d Kubernetes resources for issues and improvements:\n\n", len(resources)))

	for i, r := range resources {
		sb.WriteString(fmt.Sprintf("--- Resource %d: %s/%s ---\n", i+1, r.Kind, r.Name))
		if r.Namespace != "" {
			sb.WriteString(fmt.Sprintf("Namespace: %s\n", r.Namespace))
		}
		sb.WriteString("```json\n")
		sb.WriteString(r.JSON)
		sb.WriteString("\n```\n\n")
	}

	// Personal preference: request a per-resource TL;DR table before the deep-dive
	// so it's easy to spot the worst offenders at a glance.
	sb.WriteString("Begin with a markdown table summarising each resource (name, kind, worst severity, one-line finding), ")
	sb.WriteString("then provide a consolidated analysis covering cross-resource issues, dependencies, and overall cluster health implications.")

	return sb.String()
}

// resourceEntry holds the data needed to build a multi-resource prompt.
type resourceEntry struct {
	Kind      string
	Name      string
	Namespace string
	JSON      string
}

// severityLevels defines the recognized severity levels for issues.
// Includes "Info" in addition to the standard Critical/High/Medium/Low so that
// purely informational observations can be surfaced without implying a problem.
var severityLevels = []string{"Critical", "High", "Medium", "Low", "Info"}
