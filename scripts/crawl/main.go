package main

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type GitHubSearchResponse struct {
	TotalCount        int          `json:"total_count"`
	IncompleteResults bool         `json:"incomplete_results"`
	Items             []GitHubRepo `json:"items"`
}

type GitHubRepo struct {
	FullName      string `json:"full_name"`
	HTMLURL       string `json:"html_url"`
	DefaultBranch string `json:"default_branch"`
	Description   string `json:"description"`
	UpdatedAt     string `json:"updated_at"`
}

type CommunityBundle struct {
	ID           string `json:"id"`
	Owner        string `json:"owner"`
	Repo         string `json:"repo"`
	Version      string `json:"version"`
	OKFVersion   string `json:"okf_version,omitempty"`
	License      string `json:"license,omitempty"`
	Status       string `json:"status"` // "stable", "draft", "deprecated" (OKF v0.2 §5.3)
	IsStale      bool   `json:"is_stale,omitempty"`
	Tier         string `json:"tier"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	RepoURL      string `json:"repo_url"`
	ConceptCount int    `json:"concept_count"`
	UpdatedAt    string `json:"updated_at"`
}

type ETagCache struct {
	SearchETag string `json:"search_etag"`
}

type OKFValidationResult struct {
	DeclaredVersion string   `json:"declared_version"`
	ConceptCount    int      `json:"concept_count"`
	Errors          []string `json:"errors"`
	GateFindings    []string `json:"gate_findings"`
	BrokenLinks     []any    `json:"broken_links"`
	Orphans         []string `json:"orphans"`
	StaleCount      int      `json:"stale_count"`
	IsConformant    bool     `json:"is_conformant"`
	GatePassed      bool     `json:"gate_passed"`
}

func main() {
	dataFile := filepath.Join("data", "community.json")
	etagFile := filepath.Join("data", "cache_etags.json")

	if err := os.MkdirAll("data", 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating data dir: %v\n", err)
		os.Exit(1)
	}

	okfBinary, err := exec.LookPath("okf")
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠️ 'okf' CLI not found in PATH. Install with: go install github.com/okf-memory/okf-agent-memory/cmd/okf@latest\n")
		if _, statErr := os.Stat(dataFile); statErr == nil {
			fmt.Println("Preserving existing community.json data.")
			return
		}
		os.Exit(1)
	}
	fmt.Printf("✓ Using OKF validator binary: %s\n", okfBinary)

	var etagCache ETagCache
	if raw, err := os.ReadFile(etagFile); err == nil {
		_ = json.Unmarshal(raw, &etagCache)
	}

	token := os.Getenv("GH_TOKEN")
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}

	client := &http.Client{Timeout: 30 * time.Second}
	apiURL := "https://api.github.com/search/repositories?q=topic:okf-memory-bundle+topic:okf-memory&sort=updated&order=desc&per_page=100"

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating request: %v\n", err)
		os.Exit(1)
	}

	req.Header.Set("User-Agent", "okf-memory-crawler/0.2")
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if etagCache.SearchETag != "" {
		req.Header.Set("If-None-Match", etagCache.SearchETag)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("⚠️ Network issue querying GitHub API: %v (continuing with cached community data)\n", err)
		setGitHubOutput("changed", "false")
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusNotModified {
		fmt.Println("✓ Community repositories unchanged (304 Not Modified, 0 rate-limit used)")
		setGitHubOutput("changed", "false")
		return
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("⚠️ GitHub API returned HTTP %d: %s (continuing with cached data)\n", resp.StatusCode, string(body))
		setGitHubOutput("changed", "false")
		return
	}

	newETag := resp.Header.Get("ETag")
	if newETag != "" {
		etagCache.SearchETag = newETag
		if etagBytes, err := json.MarshalIndent(etagCache, "", "  "); err == nil {
			_ = os.WriteFile(etagFile, etagBytes, 0o644)
		}
	}

	var searchResp GitHubSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding GitHub search response: %v\n", err)
		return
	}

	var communityBundles []CommunityBundle

	for _, repo := range searchResp.Items {
		// Exclude official core repository
		if strings.HasPrefix(strings.ToLower(repo.FullName), "okf-memory/") {
			continue
		}

		branch := repo.DefaultBranch
		if branch == "" {
			branch = "main"
		}

		parts := strings.Split(repo.FullName, "/")
		owner := repo.FullName
		repoName := repo.FullName
		if len(parts) == 2 {
			owner = parts[0]
			repoName = parts[1]
		}

		bundle, err := downloadAndValidateBundle(client, token, repo, branch, owner, repoName)
		if err != nil {
			fmt.Printf("❌ Rejected %s: %v\n", repo.FullName, err)
			continue
		}
		communityBundles = append(communityBundles, *bundle)
		fmt.Printf("✓ Validated & Indexed: %s (v%s, %d concepts)\n", bundle.ID, bundle.Version, bundle.ConceptCount)
	}

	outBytes, err := json.MarshalIndent(communityBundles, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling community bundles: %v\n", err)
		setGitHubOutput("changed", "false")
		return
	}

	existingBytes, _ := os.ReadFile(dataFile)
	hasChanged := !bytes.Equal(bytes.TrimSpace(existingBytes), bytes.TrimSpace(outBytes))

	if hasChanged {
		if err := os.WriteFile(dataFile, outBytes, 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", dataFile, err)
			return
		}
		fmt.Printf("✓ Finished crawl: %d verified community bundle(s) indexed (new or updated)\n", len(communityBundles))
		setGitHubOutput("changed", "true")
	} else {
		fmt.Printf("✓ Finished crawl: %d verified community bundle(s) (unchanged)\n", len(communityBundles))
		setGitHubOutput("changed", "false")
	}
}

func setGitHubOutput(name, value string) {
	if ghOut := os.Getenv("GITHUB_OUTPUT"); ghOut != "" {
		f, err := os.OpenFile(ghOut, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644)
		if err == nil {
			_, _ = fmt.Fprintf(f, "%s=%s\n", name, value)
			_ = f.Close()
		}
	}
}

func downloadAndValidateBundle(client *http.Client, token string, repo GitHubRepo, branch, owner, repoName string) (*CommunityBundle, error) {
	tmpDir, err := os.MkdirTemp("", "okf-crawl-*")
	if err != nil {
		return nil, fmt.Errorf("internal tempdir error: %w", err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	tarURL := fmt.Sprintf("https://api.github.com/repos/%s/tarball/%s", repo.FullName, branch)
	req, err := http.NewRequest(http.MethodGet, tarURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "okf-memory-crawler/0.2")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		status := 0
		if resp != nil {
			status = resp.StatusCode
			_ = resp.Body.Close()
		}
		return nil, fmt.Errorf("failed to fetch repository tarball (HTTP %d)", status)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if err := extractTarGz(resp.Body, tmpDir); err != nil {
		return nil, fmt.Errorf("failed to extract repository archive: %w", err)
	}

	bundlePath := findBundleDir(tmpDir)
	if bundlePath == "" {
		return nil, fmt.Errorf("no valid index.md found in repository root or standard paths")
	}

	cmd := exec.Command("okf", "validate", bundlePath, "--strict", "--stale", "--json")
	out, _ := cmd.CombinedOutput()

	var valResult OKFValidationResult
	if err := json.Unmarshal(out, &valResult); err != nil {
		return nil, fmt.Errorf("okf validator returned malformed output: %w", err)
	}

	if !valResult.GatePassed || !valResult.IsConformant || len(valResult.Errors) > 0 {
		return nil, fmt.Errorf("strict OKF conformance validation failed (%d error(s))", len(valResult.Errors))
	}

	indexPath := filepath.Join(bundlePath, "index.md")
	f, err := os.Open(indexPath)
	if err != nil {
		return nil, fmt.Errorf("unable to open index.md: %w", err)
	}
	defer func() {
		_ = f.Close()
	}()

	fm, ok := parseFrontmatter(f)
	if !ok || fm["okf_version"] != "0.2" {
		return nil, fmt.Errorf("index.md frontmatter must declare okf_version: \"0.2\"")
	}

	version := fm["bundle_version"]
	if version == "" {
		version = "1.0.0"
	}
	title := fm["title"]
	if title == "" {
		title = repo.FullName
	}
	desc := fm["description"]
	if desc == "" {
		desc = repo.Description
	}
	license := fm["license"]
	if license == "" {
		license = "MIT"
	}
	status := "stable"
	if s := strings.ToLower(fm["status"]); s == "draft" || s == "stable" || s == "deprecated" {
		status = s
	}
	isStale := valResult.StaleCount > 0

	return &CommunityBundle{
		ID:           "github.com/" + repo.FullName,
		Owner:        owner,
		Repo:         repoName,
		Version:      version,
		OKFVersion:   fm["okf_version"],
		License:      license,
		Status:       status,
		IsStale:      isStale,
		Tier:         "community",
		Title:        title,
		Description:  desc,
		RepoURL:      repo.HTMLURL,
		ConceptCount: valResult.ConceptCount,
		UpdatedAt:    repo.UpdatedAt,
	}, nil
}

func extractTarGz(r io.Reader, destDir string) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	grPre, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer func() {
		_ = grPre.Close()
	}()

	hasKnowledgeDir := false
	rootWrapper := ""
	firstRootChecked := false

	trPre := tar.NewReader(grPre)
	for {
		hdr, err := trPre.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		cleanName := filepath.ToSlash(filepath.Clean(hdr.Name))
		if cleanName == "." || cleanName == ".." || strings.HasPrefix(cleanName, "../") {
			continue
		}

		parts := strings.Split(cleanName, "/")
		if !firstRootChecked {
			if len(parts) > 1 || hdr.Typeflag == tar.TypeDir {
				rootWrapper = parts[0]
			}
			firstRootChecked = true
		} else if rootWrapper != "" && parts[0] != rootWrapper {
			rootWrapper = ""
		}

		if cleanName == "knowledge/index.md" || strings.HasSuffix(cleanName, "/knowledge/index.md") {
			hasKnowledgeDir = true
		}
	}

	gzr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer func() {
		_ = gzr.Close()
	}()

	tr := tar.NewReader(gzr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		cleanName := filepath.ToSlash(filepath.Clean(hdr.Name))
		if cleanName == "." || cleanName == ".." || strings.HasPrefix(cleanName, "../") {
			continue
		}

		var relDest string
		if hasKnowledgeDir {
			idx := strings.Index(cleanName, "/knowledge/")
			if idx != -1 {
				relDest = cleanName[idx+len("/knowledge/"):]
			} else if strings.HasPrefix(cleanName, "knowledge/") {
				relDest = cleanName[len("knowledge/"):]
			} else {
				continue
			}
		} else {
			if rootWrapper != "" {
				if cleanName == rootWrapper {
					continue
				}
				if strings.HasPrefix(cleanName, rootWrapper+"/") {
					relDest = cleanName[len(rootWrapper)+1:]
				} else {
					relDest = cleanName
				}
			} else {
				relDest = cleanName
			}
		}

		relDest = filepath.Clean(relDest)
		if relDest == "" || relDest == "." || relDest == ".." || strings.HasPrefix(relDest, "..") || filepath.IsAbs(relDest) {
			continue
		}

		target := filepath.Join(destDir, filepath.FromSlash(relDest))
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, hdr.FileInfo().Mode())
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				_ = f.Close()
				return err
			}
			if err := f.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

func findBundleDir(root string) string {
	if _, err := os.Stat(filepath.Join(root, "index.md")); err == nil {
		return root
	}
	var found string
	_ = filepath.Walk(root, func(p string, fi os.FileInfo, err error) error {
		if err != nil || !fi.IsDir() {
			return nil
		}
		if _, err := os.Stat(filepath.Join(p, "index.md")); err == nil {
			found = p
			return filepath.SkipDir
		}
		return nil
	})
	return found
}

func parseFrontmatter(r io.Reader) (map[string]string, bool) {
	scanner := bufio.NewScanner(r)
	if !scanner.Scan() {
		return nil, false
	}
	if strings.TrimSpace(scanner.Text()) != "---" {
		return nil, false
	}

	result := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "---" {
			return result, true
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			k := strings.TrimSpace(parts[0])
			v := strings.Trim(strings.TrimSpace(parts[1]), "\"'")
			result[k] = v
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, false
	}
	return result, false
}
