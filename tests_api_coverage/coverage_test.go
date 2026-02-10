package tests

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

type ApiDefinition struct {
	Paths map[string]map[string]interface{} `json:"paths"`
	Info  map[string]interface{}            `json:"info"`
}

type EndpointInfo struct {
	Path        string
	Method      string
	Tag         string
	Summary     string
	OperationID string
	Implemented bool
}

var (
	// Matches fmt.Sprintf("/path/%s/...", ...) or inline "/path" strings
	reSprintfPath = regexp.MustCompile(`fmt\.Sprintf\("(/[^"]+)"`)
	// Matches urlPath/endpoint/path := "/literal/path"
	reLiteralPath = regexp.MustCompile(`(?:endpoint|urlPath|path)\s*:?=\s*"(/[^"]+)"`)
	// Matches string concatenation assignments: endpoint := "/spaces/" + var + "/members"
	reConcatAssign = regexp.MustCompile(`(?:endpoint|urlPath|path)\s*:?=\s*(".+)$`)
	// Matches inline literal paths in newRequest calls
	reInlinePath = regexp.MustCompile(`newRequest\([^,]+,\s*[^,]+,\s*"(/[^"]+)"`)
	// For normalizing path parameters
	reParamPlaceholder = regexp.MustCompile(`\{[^}]*\}|%[a-zA-Z]`)
)

func TestApiCoverage(t *testing.T) {
	apiDef, err := loadApiDefinition()
	if err != nil {
		t.Fatalf("Failed to load API definition: %v", err)
	}

	endpoints := extractEndpoints(apiDef)
	checkSDKImplementation(endpoints)
	displayCoverageStats(t, endpoints)
}

func loadApiDefinition() (*ApiDefinition, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}

	data, err := os.ReadFile(filepath.Join(currentDir, "api_definition.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to read API definition file: %w", err)
	}

	var apiDef ApiDefinition
	if err := json.Unmarshal(data, &apiDef); err != nil {
		return nil, fmt.Errorf("failed to parse API definition JSON: %w", err)
	}

	return &apiDef, nil
}

func extractEndpoints(apiDef *ApiDefinition) []EndpointInfo {
	var endpoints []EndpointInfo

	for path, pathData := range apiDef.Paths {
		for method, methodData := range pathData {
			if method == "parameters" || method == "description" || method == "summary" {
				continue
			}

			methodMap, ok := methodData.(map[string]interface{})
			if !ok {
				continue
			}

			var tag, summary, operationId string
			if tags, exists := methodMap["tags"]; exists {
				if tagArr, ok := tags.([]interface{}); ok && len(tagArr) > 0 {
					tag, _ = tagArr[0].(string)
				}
			}
			if sum, exists := methodMap["summary"].(string); exists {
				summary = sum
			}
			if opid, exists := methodMap["operationId"].(string); exists {
				operationId = opid
			}

			endpoints = append(endpoints, EndpointInfo{
				Path:        path,
				Method:      strings.ToUpper(method),
				Tag:         tag,
				Summary:     summary,
				OperationID: operationId,
				Implemented: false,
			})
		}
	}

	return endpoints
}

func checkSDKImplementation(endpoints []EndpointInfo) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Warning: couldn't get current directory: %v\n", err)
		return
	}

	clientDir := filepath.Join(filepath.Dir(currentDir), "client")
	clientPaths, err := extractClientPaths(clientDir)
	if err != nil {
		fmt.Printf("Warning: couldn't extract client paths: %v\n", err)
		return
	}

	for i := range endpoints {
		apiNorm := normalizeAPIPath(endpoints[i].Path)
		for _, cp := range clientPaths {
			if apiNorm == cp {
				endpoints[i].Implemented = true
				break
			}
		}
	}
}

func normalizeAPIPath(p string) string {
	p = strings.TrimPrefix(p, "/v1")
	p = reParamPlaceholder.ReplaceAllString(p, "{}")
	return strings.TrimSuffix(p, "/")
}

func extractClientPaths(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var paths []string

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".go") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, f.Name()))
		if err != nil {
			continue
		}

		for _, line := range strings.Split(string(data), "\n") {
			trimmed := strings.TrimSpace(line)
			var raw string

			switch {
			case reSprintfPath.MatchString(trimmed):
				raw = reSprintfPath.FindStringSubmatch(trimmed)[1]
			case reInlinePath.MatchString(trimmed):
				raw = reInlinePath.FindStringSubmatch(trimmed)[1]
			case reLiteralPath.MatchString(trimmed) && !strings.Contains(trimmed, "+"):
				raw = reLiteralPath.FindStringSubmatch(trimmed)[1]
			case reConcatAssign.MatchString(trimmed) && strings.Contains(trimmed, "+"):
				raw = extractConcatPath(trimmed)
			}

			if raw == "" {
				continue
			}

			normalized := normalizePath(raw)
			if !seen[normalized] {
				seen[normalized] = true
				paths = append(paths, normalized)
			}
		}
	}
	return paths, nil
}

func extractConcatPath(line string) string {
	idx := strings.Index(line, ":=")
	if idx < 0 {
		idx = strings.Index(line, "=")
	}
	if idx < 0 {
		return ""
	}
	rhs := line[idx+2:]

	reQuoted := regexp.MustCompile(`"([^"]*)"`)
	parts := strings.Split(rhs, "+")
	var result strings.Builder
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if m := reQuoted.FindStringSubmatch(part); m != nil {
			result.WriteString(m[1])
		} else {
			result.WriteString("{}")
		}
	}
	path := result.String()
	if strings.HasPrefix(path, "/") {
		return path
	}
	return ""
}

func normalizePath(path string) string {
	path = reParamPlaceholder.ReplaceAllString(path, "{}")
	if idx := strings.Index(path, "?"); idx > -1 {
		path = path[:idx]
	}
	return path
}

func displayCoverageStats(t *testing.T, endpoints []EndpointInfo) {
	total := len(endpoints)
	implemented := 0

	tagStats := make(map[string]struct{ Total, Implemented int })
	methodStats := make(map[string]struct{ Total, Implemented int })

	for _, e := range endpoints {
		if e.Implemented {
			implemented++
		}

		ts := tagStats[e.Tag]
		ts.Total++
		if e.Implemented {
			ts.Implemented++
		}
		tagStats[e.Tag] = ts

		ms := methodStats[e.Method]
		ms.Total++
		if e.Implemented {
			ms.Implemented++
		}
		methodStats[e.Method] = ms
	}

	var coverage float64
	if total > 0 {
		coverage = float64(implemented) * 100 / float64(total)
	}

	t.Logf("\n=== API COVERAGE REPORT ===\n")
	t.Logf("Total Endpoints: %d", total)
	t.Logf("Implemented: %d", implemented)
	t.Logf("Not Implemented: %d", total-implemented)
	t.Logf("Coverage: %.1f%%\n", coverage)

	t.Logf("=== Coverage by Tag ===")
	for tag, stats := range tagStats {
		tagName := tag
		if tagName == "" {
			tagName = "(no tag)"
		}
		t.Logf("%s: %.1f%% (%d/%d)", tagName, float64(stats.Implemented)*100/float64(stats.Total), stats.Implemented, stats.Total)
	}
	t.Logf("")

	t.Logf("=== Coverage by HTTP Method ===")
	for method, stats := range methodStats {
		t.Logf("%s: %.1f%% (%d/%d)", method, float64(stats.Implemented)*100/float64(stats.Total), stats.Implemented, stats.Total)
	}
	t.Logf("")

	if total-implemented > 0 {
		t.Logf("=== Unimplemented Endpoints ===")
		unimplementedByTag := make(map[string][]EndpointInfo)
		for _, e := range endpoints {
			if !e.Implemented {
				tag := e.Tag
				if tag == "" {
					tag = "(no tag)"
				}
				unimplementedByTag[tag] = append(unimplementedByTag[tag], e)
			}
		}
		for tag, endpointList := range unimplementedByTag {
			t.Logf("\n%s:", tag)
			for _, e := range endpointList {
				opInfo := ""
				if e.OperationID != "" {
					opInfo = fmt.Sprintf(" [%s]", e.OperationID)
				}
				t.Logf("  %s %s%s", e.Method, e.Path, opInfo)
				if e.Summary != "" {
					t.Logf("    → %s", e.Summary)
				}
			}
		}
		t.Logf("")
	}
}
