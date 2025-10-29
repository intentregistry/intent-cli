package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// PackageRegistry stores published packages in memory
type PackageRegistry struct {
	mu       sync.RWMutex
	packages map[string]*Package
	files    map[string][]byte // filename -> content
}

type Package struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Tarball   string `json:"tarball"`
	SHA256    string `json:"sha256"`
	PublishedAt time.Time `json:"published_at"`
}

var registry = &PackageRegistry{
	packages: make(map[string]*Package),
	files:    make(map[string][]byte),
}

// POST /v1/packages/publish - Publish a package
func publishHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(50 << 20) // 50MB
	if err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	// Get uploaded file
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "no file uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file content
	fileContent := make([]byte, fileHeader.Size)
	if _, err := file.Read(fileContent); err != nil {
		http.Error(w, "failed to read file", http.StatusInternalServerError)
		return
	}

	// Parse form data
	sha256 := r.FormValue("sha256")
	if sha256 == "" {
		http.Error(w, "sha256 required", http.StatusBadRequest)
		return
	}

	// Extract package name and version from filename
	// Expected format: @scope/name-version.itpkg or name-version.itpkg
	filename := fileHeader.Filename
	name, version, err := parsePackageName(filename)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid package name: %v", err), http.StatusBadRequest)
		return
	}

	pkgKey := fmt.Sprintf("%s@%s", name, version)

	registry.mu.Lock()
	defer registry.mu.Unlock()

	// Check if version already exists
	if _, exists := registry.packages[pkgKey]; exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": "version already exists"})
		return
	}

	// Store package
	tarballPath := fmt.Sprintf("/v1/packages/tarball/%s/%s.itpkg", name, version)
	registry.packages[pkgKey] = &Package{
		Name:        name,
		Version:     version,
		Tarball:     fmt.Sprintf("http://localhost:8080%s", tarballPath),
		SHA256:      sha256,
		PublishedAt: time.Now(),
	}
	registry.files[tarballPath] = fileContent

	log.Printf("✅ Published: %s@%s (sha256: %s)", name, version, sha256)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "package published successfully",
		"name":    name,
		"version": version,
	})
}

// GET /v1/packages/resolve?spec=@scope/name[@version] - Resolve package
func resolveHandler(w http.ResponseWriter, r *http.Request) {
	spec := r.URL.Query().Get("spec")
	if spec == "" {
		http.Error(w, "spec parameter required", http.StatusBadRequest)
		return
	}

	// Parse spec: @scope/name or @scope/name@version
	name, version := parseSpec(spec)

	registry.mu.RLock()
	defer registry.mu.RUnlock()

	// If no version specified, find latest
	if version == "" {
		latest := findLatestVersion(name)
		if latest == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "package not found"})
			return
		}
		version = latest
	}

	pkgKey := fmt.Sprintf("%s@%s", name, version)
	pkg, exists := registry.packages[pkgKey]
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "package not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pkg)
}

// GET /v1/packages/tarball/:name/:version.itpkg - Download tarball
func tarballHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	// Extract name and version from path: /v1/packages/tarball/@scope/name/version.itpkg
	parts := strings.Split(path, "/")
	if len(parts) < 6 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	registry.mu.RLock()
	defer registry.mu.RUnlock()

	content, exists := registry.files[path]
	if !exists {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/gzip")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))
	w.Write(content)
}

// GET /v1/packages/search?q=query - Search packages
func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	registry.mu.RLock()
	defer registry.mu.RUnlock()

	var results []*Package
	for _, pkg := range registry.packages {
		if query == "" || strings.Contains(pkg.Name, query) {
			results = append(results, pkg)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"count":    len(results),
		"packages": results,
	})
}

// GET /health - Health check
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"time":   time.Now(),
	})
}

// Helper: Parse package name from filename
func parsePackageName(filename string) (string, string, error) {
	// Remove extension
	filename = strings.TrimSuffix(filename, ".itpkg")

	// Expected format: scope-name-version or name-version
	// For now, assume last part is version
	parts := strings.Split(filename, "-")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid package name format")
	}

	version := parts[len(parts)-1]
	name := strings.Join(parts[:len(parts)-1], "-")

	return name, version, nil
}

// Helper: Parse spec like @scope/name or @scope/name@version
func parseSpec(spec string) (string, string) {
	// Remove @ prefix if present
	spec = strings.TrimPrefix(spec, "@")

	// Split by @
	parts := strings.Split(spec, "@")
	if len(parts) == 1 {
		// No version specified
		return parts[0], ""
	}
	return parts[0], parts[1]
}

// Helper: Find latest version of a package
func findLatestVersion(name string) string {
	var latest string
	var latestTime time.Time

	for key, pkg := range registry.packages {
		if strings.HasPrefix(key, name+"@") {
			if pkg.PublishedAt.After(latestTime) {
				latestTime = pkg.PublishedAt
				latest = pkg.Version
			}
		}
	}
	return latest
}

func main() {
	port := "8080"
	if env := os.Getenv("PORT"); env != "" {
		port = env
	}

	// Routes
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/v1/packages/publish", publishHandler)
	http.HandleFunc("/v1/packages/resolve", resolveHandler)
	http.HandleFunc("/v1/packages/search", searchHandler)
	http.HandleFunc("/v1/packages/tarball/", tarballHandler)

	log.Printf("🚀 Mock IntentRegistry API starting on http://localhost:%s", port)
	log.Printf("   Health: http://localhost:%s/health", port)
	log.Printf("   Publish: POST /v1/packages/publish", port)
	log.Printf("   Resolve: GET /v1/packages/resolve?spec=@scope/name[@version]", port)
	log.Printf("   Search: GET /v1/packages/search?q=query", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
