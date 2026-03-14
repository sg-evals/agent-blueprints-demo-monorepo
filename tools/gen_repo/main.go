package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

var (
	seed      = flag.Int64("seed", 42, "Random seed for deterministic generation")
	services  = flag.Int("services", 10, "Number of service stubs to generate")
	libs      = flag.Int("libs", 18, "Number of library stubs to generate")
	outputDir = flag.String("output-dir", ".", "Output directory for generated files")
)

var serviceAdjectives = []string{
	"fast", "smart", "core", "edge", "cloud", "data", "stream", "batch",
	"async", "sync", "multi", "micro", "hyper", "auto", "meta", "prime",
	"next", "base", "open", "flex", "deep", "safe", "live", "real",
}

var serviceNouns = []string{
	"processor", "handler", "manager", "controller", "dispatcher", "scheduler",
	"aggregator", "transformer", "validator", "resolver", "analyzer", "monitor",
	"collector", "publisher", "consumer", "indexer", "exporter", "importer",
}

var libPrefixes = []string{
	"go", "lib", "pkg", "util", "common", "shared", "core", "base",
}

var libSuffixes = []string{
	"cache", "queue", "pool", "codec", "store", "config", "crypto", "hash",
	"limiter", "mapper", "parser", "render", "router", "schema", "serial",
	"socket", "stream", "timer", "token", "tracer", "worker", "writer",
}

func main() {
	flag.Parse()

	rng := rand.New(rand.NewSource(*seed))

	svcNames := generateServiceNames(rng, *services)
	libNames := generateLibNames(rng, *libs)

	fmt.Printf("Generating with seed=%d: %d services, %d libs\n", *seed, *services, *libs)
	fmt.Printf("Output directory: %s\n\n", *outputDir)

	for _, name := range svcNames {
		if err := generateService(rng, *outputDir, name, libNames); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating service %s: %v\n", name, err)
			os.Exit(1)
		}
		fmt.Printf("  Created service: apps/%s\n", name)
	}

	for _, name := range libNames {
		if err := generateLib(*outputDir, name); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating lib %s: %v\n", name, err)
			os.Exit(1)
		}
		fmt.Printf("  Created lib:     libs/%s\n", name)
	}

	fmt.Printf("\nDone. Generated %d services and %d libraries.\n", len(svcNames), len(libNames))
}

func generateServiceNames(rng *rand.Rand, count int) []string {
	used := make(map[string]bool)
	var names []string
	for len(names) < count {
		adj := serviceAdjectives[rng.Intn(len(serviceAdjectives))]
		noun := serviceNouns[rng.Intn(len(serviceNouns))]
		name := adj + "-" + noun
		if !used[name] {
			used[name] = true
			names = append(names, name)
		}
	}
	return names
}

func generateLibNames(rng *rand.Rand, count int) []string {
	used := make(map[string]bool)
	var names []string
	for len(names) < count {
		prefix := libPrefixes[rng.Intn(len(libPrefixes))]
		suffix := libSuffixes[rng.Intn(len(libSuffixes))]
		name := prefix + suffix
		if !used[name] {
			used[name] = true
			names = append(names, name)
		}
	}
	return names
}

func generateService(rng *rand.Rand, outDir, name string, availableLibs []string) error {
	dir := filepath.Join(outDir, "apps", name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	modulePath := "github.com/sg-evals/agent-blueprints-demo-monorepo/apps/" + name
	goMod := fmt.Sprintf("module %s\n\ngo 1.22\n", modulePath)
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0o644); err != nil {
		return err
	}

	pkgName := toIdentifier(name)
	port := 8080 + rng.Intn(920)

	// Pick 1-3 random endpoints
	endpoints := []string{"/healthz"}
	possibleEndpoints := []string{"/status", "/metrics", "/ready", "/info", "/version", "/ping"}
	numExtra := 1 + rng.Intn(3)
	rng.Shuffle(len(possibleEndpoints), func(i, j int) {
		possibleEndpoints[i], possibleEndpoints[j] = possibleEndpoints[j], possibleEndpoints[i]
	})
	for i := 0; i < numExtra && i < len(possibleEndpoints); i++ {
		endpoints = append(endpoints, possibleEndpoints[i])
	}

	var handlerRegistrations strings.Builder
	var handlerFuncs strings.Builder

	for _, ep := range endpoints {
		handlerName := "handle" + toCamelCase(ep)
		handlerRegistrations.WriteString(fmt.Sprintf("\thttp.HandleFunc(%q, %s)\n", ep, handlerName))
		handlerFuncs.WriteString(fmt.Sprintf(`
func %s(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "{\"service\":%q,\"endpoint\":%q,\"status\":\"ok\"}\n")
}
`, handlerName, pkgName, ep))
	}

	mainGo := fmt.Sprintf(`package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Printf("Starting %s on :%d")
%s
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", %d), nil))
}
%s`, pkgName, port, handlerRegistrations.String(), "%d", port, handlerFuncs.String())

	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(mainGo), 0o644); err != nil {
		return err
	}

	// Generate test file
	testGo := fmt.Sprintf(`package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleHealthz(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	handleHealthz(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %%d, got %%d", http.StatusOK, w.Code)
	}

	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %%s", ct)
	}
}
`)

	if err := os.WriteFile(filepath.Join(dir, "main_test.go"), []byte(testGo), 0o644); err != nil {
		return err
	}

	return nil
}

func generateLib(outDir, name string) error {
	dir := filepath.Join(outDir, "libs", name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	modulePath := "github.com/sg-evals/agent-blueprints-demo-monorepo/libs/" + name
	goMod := fmt.Sprintf("module %s\n\ngo 1.22\n", modulePath)
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0o644); err != nil {
		return err
	}

	pkgName := toIdentifier(name)

	srcGo := fmt.Sprintf(`// Package %s provides %s functionality.
package %s

// Config holds configuration for %s.
type Config struct {
	Enabled bool
	Name    string
}

// New creates a new %s instance with the given config.
func New(cfg Config) *Instance {
	return &Instance{config: cfg}
}

// Instance is the main type for %s.
type Instance struct {
	config Config
}

// IsEnabled returns whether the instance is enabled.
func (inst *Instance) IsEnabled() bool {
	return inst.config.Enabled
}

// GetName returns the configured name.
func (inst *Instance) GetName() string {
	return inst.config.Name
}
`, pkgName, name, pkgName, name, name, name)

	if err := os.WriteFile(filepath.Join(dir, pkgName+".go"), []byte(srcGo), 0o644); err != nil {
		return err
	}

	testGo := fmt.Sprintf(`package %s

import "testing"

func TestNew(t *testing.T) {
	inst := New(Config{Enabled: true, Name: %q})
	if !inst.IsEnabled() {
		t.Error("expected instance to be enabled")
	}
	if inst.GetName() != %q {
		t.Errorf("expected name %%q, got %%q", %q, inst.GetName())
	}
}

func TestDisabled(t *testing.T) {
	inst := New(Config{Enabled: false, Name: ""})
	if inst.IsEnabled() {
		t.Error("expected instance to be disabled")
	}
}
`, pkgName, name, name, name)

	if err := os.WriteFile(filepath.Join(dir, pkgName+"_test.go"), []byte(testGo), 0o644); err != nil {
		return err
	}

	return nil
}

func toIdentifier(name string) string {
	return strings.ReplaceAll(strings.ReplaceAll(name, "-", ""), "_", "")
}

func toCamelCase(path string) string {
	path = strings.TrimPrefix(path, "/")
	parts := strings.Split(path, "-")
	var result strings.Builder
	for _, p := range parts {
		if len(p) > 0 {
			result.WriteString(strings.ToUpper(p[:1]) + p[1:])
		}
	}
	return result.String()
}
