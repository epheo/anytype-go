package tests

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
)

type ApiDefinition struct {
	Paths      map[string]map[string]interface{} `json:"paths"`
	Components struct {
		Schemas map[string]struct {
			Enum []interface{} `json:"enum"`
		} `json:"schemas"`
	} `json:"components"`
	Info map[string]interface{} `json:"info"`
}

type EndpointInfo struct {
	Path        string
	Method      string
	Tag         string
	Summary     string
	OperationID string
	Implemented bool
}

type methodPath struct {
	Method string
	Path   string
}

// Endpoints the SDK consciously does not implement, keyed by operationId.
// Anything unimplemented and absent here fails the test.
var allowedUnimplemented = map[string]bool{}

// SDK typed string enums checked against the spec's top-level schema enums.
// Names differ on purpose (FilterFormat vs PropertyFormat), so the link is explicit.
var enumMappings = []struct{ SDKType, SpecSchema string }{
	{"FilterCondition", "FilterCondition"},
	{"FilterOperator", "FilterOperator"},
	{"FilterFormat", "PropertyFormat"},
	{"PropertyFormat", "PropertyFormat"},
	{"TypeLayout", "TypeLayout"},
	{"Color", "Color"},
	{"SortDirection", "SortDirection"},
	{"SortProperty", "SortProperty"},
}

func TestApiCoverage(t *testing.T) {
	apiDef, err := loadApiDefinition()
	if err != nil {
		t.Fatalf("load API definition: %v", err)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	root := filepath.Dir(currentDir)

	endpoints := extractEndpoints(apiDef)
	clientPaths, err := extractClientPaths(filepath.Join(root, "client"))
	if err != nil {
		t.Fatalf("extract client paths: %v", err)
	}

	markImplemented(endpoints, clientPaths)
	displayCoverageStats(t, endpoints)

	// Gate 1: every unimplemented endpoint must be explicitly allowlisted.
	for _, e := range endpoints {
		if !e.Implemented && !allowedUnimplemented[e.OperationID] {
			t.Errorf("unimplemented endpoint not in allowlist: %s %s [%s]", e.Method, e.Path, e.OperationID)
		}
	}

	// Gate 2: every client route must map to a real endpoint; an orphan means a
	// typo or a path the extractor mis-read (the class of bug this test once hid).
	known := endpointSet(endpoints)
	for _, cp := range clientPaths {
		if !known[cp] {
			t.Errorf("orphan client route matches no API endpoint: %s %s", cp.Method, cp.Path)
		}
	}

	// Gate 3: SDK typed enums must cover their spec counterparts.
	checkEnumDrift(t, apiDef, root)
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
			if tags, ok := methodMap["tags"].([]interface{}); ok && len(tags) > 0 {
				tag, _ = tags[0].(string)
			}
			summary, _ = methodMap["summary"].(string)
			operationId, _ = methodMap["operationId"].(string)

			endpoints = append(endpoints, EndpointInfo{
				Path:        path,
				Method:      strings.ToUpper(method),
				Tag:         tag,
				Summary:     summary,
				OperationID: operationId,
			})
		}
	}

	return endpoints
}

// extractClientPaths parses the client package and folds each newRequest call's
// method and path argument into a normalized (method, path) pair. Parsing the AST
// rather than scraping text handles inline concatenation, fmt.Sprintf, and path
// variables uniformly, which regex matching did not.
func extractClientPaths(dir string) ([]methodPath, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	seen := make(map[methodPath]bool)
	var paths []methodPath

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".go") || strings.HasSuffix(f.Name(), "_test.go") {
			continue
		}

		node, err := parser.ParseFile(fset, filepath.Join(dir, f.Name()), nil, 0)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", f.Name(), err)
		}

		for _, decl := range node.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				continue
			}
			assigns := collectStringAssigns(fn.Body)

			ast.Inspect(fn.Body, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}
				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok || sel.Sel.Name != "newRequest" || len(call.Args) < 3 {
					return true
				}

				method := methodFromExpr(call.Args[1])
				raw := evalPathExpr(call.Args[2], assigns)
				if method == "" || !strings.HasPrefix(raw, "/") {
					return true
				}

				mp := methodPath{Method: method, Path: normalizePath(raw)}
				if !seen[mp] {
					seen[mp] = true
					paths = append(paths, mp)
				}
				return true
			})
		}
	}

	return paths, nil
}

// collectStringAssigns maps each locally assigned variable to its value expression
// so evalPathExpr can resolve a path passed to newRequest by name. Compound
// assignments (endpoint += "?"+query) are skipped: the query is dropped anyway.
func collectStringAssigns(body *ast.BlockStmt) map[string]ast.Expr {
	assigns := make(map[string]ast.Expr)
	ast.Inspect(body, func(n ast.Node) bool {
		as, ok := n.(*ast.AssignStmt)
		if !ok || (as.Tok != token.DEFINE && as.Tok != token.ASSIGN) {
			return true
		}
		if len(as.Lhs) != 1 || len(as.Rhs) != 1 {
			return true
		}
		if id, ok := as.Lhs[0].(*ast.Ident); ok {
			assigns[id.Name] = as.Rhs[0]
		}
		return true
	})
	return assigns
}

// evalPathExpr folds a path expression into a template, rendering every dynamic
// segment as "{}" so it aligns with a spec path parameter after normalization.
func evalPathExpr(expr ast.Expr, assigns map[string]ast.Expr) string {
	switch e := expr.(type) {
	case *ast.BasicLit:
		return litString(e)
	case *ast.BinaryExpr:
		if e.Op == token.ADD {
			return evalPathExpr(e.X, assigns) + evalPathExpr(e.Y, assigns)
		}
	case *ast.Ident:
		if v, ok := assigns[e.Name]; ok {
			return evalPathExpr(v, assigns)
		}
	case *ast.CallExpr:
		if isSprintf(e) && len(e.Args) > 0 {
			if lit, ok := e.Args[0].(*ast.BasicLit); ok {
				return litString(lit)
			}
		}
		// withListParams(path, opts) wraps the real path in its first argument.
		if fn, ok := e.Fun.(*ast.Ident); ok && fn.Name == "withListParams" && len(e.Args) > 0 {
			return evalPathExpr(e.Args[0], assigns)
		}
	}
	return "{}"
}

func methodFromExpr(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.SelectorExpr: // http.MethodGet
		if x, ok := e.X.(*ast.Ident); ok && x.Name == "http" {
			return strings.ToUpper(strings.TrimPrefix(e.Sel.Name, "Method"))
		}
	case *ast.BasicLit:
		return strings.ToUpper(litString(e))
	}
	return ""
}

func isSprintf(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	return ok && pkg.Name == "fmt" && sel.Sel.Name == "Sprintf"
}

func litString(lit *ast.BasicLit) string {
	if lit.Kind != token.STRING {
		return ""
	}
	if s, err := strconv.Unquote(lit.Value); err == nil {
		return s
	}
	return strings.Trim(lit.Value, "`\"")
}

// reParamPlaceholder collapses spec params ({id}) and format verbs (%s) to "{}".
var reParamPlaceholder = regexp.MustCompile(`\{[^}]*\}|%[a-zA-Z]`)

// normalizePath renders spec and client paths into one comparable form: no /v1
// prefix, no query string, dynamic segments as "{}", no trailing slash.
func normalizePath(p string) string {
	if i := strings.Index(p, "?"); i >= 0 {
		p = p[:i]
	}
	p = strings.TrimPrefix(p, "/v1")
	p = reParamPlaceholder.ReplaceAllString(p, "{}")
	return strings.TrimSuffix(p, "/")
}

func markImplemented(endpoints []EndpointInfo, clientPaths []methodPath) {
	have := make(map[methodPath]bool, len(clientPaths))
	for _, cp := range clientPaths {
		have[cp] = true
	}
	for i := range endpoints {
		mp := methodPath{Method: endpoints[i].Method, Path: normalizePath(endpoints[i].Path)}
		endpoints[i].Implemented = have[mp]
	}
}

func endpointSet(endpoints []EndpointInfo) map[methodPath]bool {
	set := make(map[methodPath]bool, len(endpoints))
	for _, e := range endpoints {
		set[methodPath{Method: e.Method, Path: normalizePath(e.Path)}] = true
	}
	return set
}

func checkEnumDrift(t *testing.T, apiDef *ApiDefinition, root string) {
	spec := specEnums(apiDef)
	sdk, err := sdkEnums(root)
	if err != nil {
		t.Errorf("extract SDK enums: %v", err)
		return
	}

	t.Logf("=== Enum Drift (%d mappings) ===", len(enumMappings))
	for _, m := range enumMappings {
		specVals, ok := spec[m.SpecSchema]
		if !ok {
			t.Errorf("enum mapping %q: spec schema %q not found", m.SDKType, m.SpecSchema)
			continue
		}
		sdkVals, ok := sdk[m.SDKType]
		if !ok {
			t.Errorf("enum mapping %q: SDK type not found in package", m.SDKType)
			continue
		}

		sdkSet := toSet(sdkVals)
		var missing []string
		for _, v := range specVals {
			if !sdkSet[v] {
				missing = append(missing, v)
			}
		}
		if len(missing) > 0 {
			sort.Strings(missing)
			t.Errorf("SDK %s is missing spec %s values: %v", m.SDKType, m.SpecSchema, missing)
			continue
		}
		t.Logf("%s: OK (%d values)", m.SDKType, len(specVals))
	}
	t.Logf("")
}

func specEnums(apiDef *ApiDefinition) map[string][]string {
	out := make(map[string][]string)
	for name, sch := range apiDef.Components.Schemas {
		if len(sch.Enum) == 0 {
			continue
		}
		vals := make([]string, 0, len(sch.Enum))
		for _, v := range sch.Enum {
			if s, ok := v.(string); ok {
				vals = append(vals, s)
			}
		}
		out[name] = vals
	}
	return out
}

// sdkEnums collects the values of every typed string constant in the root package,
// keyed by type name, so they can be checked against the spec without hand-copying.
func sdkEnums(dir string) (map[string][]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	out := make(map[string][]string)

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".go") || strings.HasSuffix(f.Name(), "_test.go") {
			continue
		}
		node, err := parser.ParseFile(fset, filepath.Join(dir, f.Name()), nil, 0)
		if err != nil {
			continue
		}
		for _, decl := range node.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.CONST {
				continue
			}
			for _, spec := range gd.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok || len(vs.Values) == 0 {
					continue
				}
				typeName, ok := vs.Type.(*ast.Ident)
				if !ok {
					continue
				}
				for _, val := range vs.Values {
					if lit, ok := val.(*ast.BasicLit); ok && lit.Kind == token.STRING {
						out[typeName.Name] = append(out[typeName.Name], litString(lit))
					}
				}
			}
		}
	}
	return out, nil
}

func toSet(vals []string) map[string]bool {
	set := make(map[string]bool, len(vals))
	for _, v := range vals {
		set[v] = true
	}
	return set
}

func displayCoverageStats(t *testing.T, endpoints []EndpointInfo) {
	total := len(endpoints)
	implemented := 0

	type counter struct{ Total, Implemented int }
	tagStats := make(map[string]*counter)
	methodStats := make(map[string]*counter)

	bump := func(m map[string]*counter, key string, done bool) {
		c := m[key]
		if c == nil {
			c = &counter{}
			m[key] = c
		}
		c.Total++
		if done {
			c.Implemented++
		}
	}

	for _, e := range endpoints {
		if e.Implemented {
			implemented++
		}
		bump(tagStats, e.Tag, e.Implemented)
		bump(methodStats, e.Method, e.Implemented)
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
	for _, tag := range sortedKeys(tagStats) {
		s := tagStats[tag]
		name := tag
		if name == "" {
			name = "(no tag)"
		}
		t.Logf("%s: %.1f%% (%d/%d)", name, pct(s.Implemented, s.Total), s.Implemented, s.Total)
	}
	t.Logf("")

	t.Logf("=== Coverage by HTTP Method ===")
	for _, method := range sortedKeys(methodStats) {
		s := methodStats[method]
		t.Logf("%s: %.1f%% (%d/%d)", method, pct(s.Implemented, s.Total), s.Implemented, s.Total)
	}
	t.Logf("")

	if total-implemented == 0 {
		return
	}

	t.Logf("=== Unimplemented Endpoints ===")
	byTag := make(map[string][]EndpointInfo)
	for _, e := range endpoints {
		if e.Implemented {
			continue
		}
		tag := e.Tag
		if tag == "" {
			tag = "(no tag)"
		}
		byTag[tag] = append(byTag[tag], e)
	}
	for _, tag := range sortedKeys(byTag) {
		t.Logf("\n%s:", tag)
		for _, e := range byTag[tag] {
			opInfo := ""
			if e.OperationID != "" {
				opInfo = fmt.Sprintf(" [%s]", e.OperationID)
			}
			t.Logf("  %s %s%s", e.Method, e.Path, opInfo)
			if e.Summary != "" {
				t.Logf("    %s", e.Summary)
			}
		}
	}
	t.Logf("")
}

func pct(done, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(done) * 100 / float64(total)
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
