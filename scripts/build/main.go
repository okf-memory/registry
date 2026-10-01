package main

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type BundleItem struct {
	ID          string `json:"id"`
	Version     string `json:"version"`
	OKFVersion  string `json:"okf_version,omitempty"`
	License     string `json:"license,omitempty"`
	Status      string `json:"status"` // "stable", "draft", "deprecated" (OKF v0.2 §5.3)
	IsStale     bool   `json:"is_stale,omitempty"`
	Tier        string `json:"tier"` // "official" or "community"
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Hash        string `json:"hash"`
	ManifestURL string `json:"manifest_url"`
	DownloadURL string `json:"download_url"`
	InstallCmd  string `json:"install_cmd,omitempty"`
}

type IndexManifest struct {
	Version   int          `json:"version"`
	UpdatedAt string       `json:"updated_at"`
	Bundles   []BundleItem `json:"bundles"`
}

type BundleManifest struct {
	ID          string `json:"id"`
	Version     string `json:"version"`
	OKFVersion  string `json:"okf_version,omitempty"`
	License     string `json:"license,omitempty"`
	Status      string `json:"status"` // "stable", "draft", "deprecated" (OKF v0.2 §5.3)
	IsStale     bool   `json:"is_stale,omitempty"`
	Hash        string `json:"hash"`
	DownloadURL string `json:"download_url"`
}

func main() {
	bundlesDir := "bundles"
	publicDir := "public"

	if err := os.RemoveAll(publicDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error cleaning public dir: %v\n", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(filepath.Join(publicDir, "downloads"), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating public/downloads: %v\n", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(filepath.Join(publicDir, "bundles"), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating public/bundles: %v\n", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(filepath.Join(publicDir, "badge"), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating public/badge: %v\n", err)
		os.Exit(1)
	}

	// Copy static canonical badge
	if badgeSvg, err := os.ReadFile(filepath.Join("src", "badge.svg")); err == nil {
		_ = os.WriteFile(filepath.Join(publicDir, "badge", "v0.2.svg"), badgeSvg, 0o644)
		_ = os.WriteFile(filepath.Join(publicDir, "badge", "okf-memory-bundle.svg"), badgeSvg, 0o644)
	}

	// Write .nojekyll to ensure GitHub Pages serves raw assets and doesn't run Jekyll
	_ = os.WriteFile(filepath.Join(publicDir, ".nojekyll"), []byte{}, 0o644)
	if cnameBytes, err := os.ReadFile("CNAME"); err == nil {
		_ = os.WriteFile(filepath.Join(publicDir, "CNAME"), cnameBytes, 0o644)
	}

	var indexItems []BundleItem

	err := filepath.Walk(bundlesDir, func(p string, fi os.FileInfo, err error) error {
		if err != nil || !fi.IsDir() {
			return nil
		}
		indexPath := filepath.Join(p, "index.md")
		if _, err := os.Stat(indexPath); err != nil {
			return nil
		}

		rel, err := filepath.Rel(bundlesDir, p)
		if err != nil || rel == "." {
			return nil
		}
		bundleID := filepath.ToSlash(rel)
		version := "1.0.0"
		okfVersion := "0.2"
		license := "MIT"
		status := "stable"
		isStale := false
		title := bundleID
		desc := ""

		// Read index.md to extract version, title, description, license, status if present
		if content, err := os.ReadFile(indexPath); err == nil {
			lines := strings.Split(string(content), "\n")
			for _, l := range lines {
				l = strings.TrimSpace(l)
				if strings.HasPrefix(l, "bundle_version:") || strings.HasPrefix(l, "version:") {
					parts := strings.SplitN(l, ":", 2)
					if len(parts) == 2 {
						v := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
						if v != "" {
							version = v
						}
					}
				} else if strings.HasPrefix(l, "okf_version:") {
					parts := strings.SplitN(l, ":", 2)
					if len(parts) == 2 {
						ov := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
						if ov != "" {
							okfVersion = ov
						}
					}
				} else if strings.HasPrefix(l, "license:") {
					parts := strings.SplitN(l, ":", 2)
					if len(parts) == 2 {
						lic := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
						if lic != "" {
							license = lic
						}
					}
				} else if strings.HasPrefix(l, "status:") {
					parts := strings.SplitN(l, ":", 2)
					if len(parts) == 2 {
						s := strings.ToLower(strings.Trim(strings.TrimSpace(parts[1]), `"'`))
						if s == "draft" || s == "stable" || s == "deprecated" {
							status = s
						}
					}
				} else if strings.HasPrefix(l, "title:") {
					parts := strings.SplitN(l, ":", 2)
					if len(parts) == 2 {
						t := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
						if t != "" {
							title = t
						}
					}
				} else if strings.HasPrefix(l, "description:") {
					parts := strings.SplitN(l, ":", 2)
					if len(parts) == 2 {
						d := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
						if d != "" {
							desc = d
						}
					}
				}
			}
		}

		// Automated staleness check via okf CLI
		if okfBin, err := exec.LookPath("okf"); err == nil {
			cmd := exec.Command(okfBin, "validate", p, "--json")
			if out, err := cmd.Output(); err == nil && len(out) > 0 {
				var vRes struct {
					StaleCount int `json:"stale_count"`
				}
				if err := json.Unmarshal(out, &vRes); err == nil && vRes.StaleCount > 0 {
					isStale = true
				}
			}
		}

		// Package into deterministic .tgz with BUNDLEID-VERSION.tgz
		baseName := filepath.Base(rel)
		parentDir := filepath.Dir(rel)
		fileNameWithVersion := fmt.Sprintf("%s-%s.tgz", baseName, version)

		var dlRelPath string
		if parentDir == "." {
			dlRelPath = filepath.Join("downloads", fileNameWithVersion)
		} else {
			dlRelPath = filepath.Join("downloads", parentDir, fileNameWithVersion)
		}
		tgzPath := filepath.Join(publicDir, dlRelPath)
		if err := os.MkdirAll(filepath.Dir(tgzPath), 0o755); err != nil {
			return err
		}

		hash, err := packageTgz(p, tgzPath)
		if err != nil {
			return fmt.Errorf("failed to package %s: %w", bundleID, err)
		}

		safeTagID := strings.ReplaceAll(bundleID, "/", "-")
		releaseTag := fmt.Sprintf("%s-v%s", safeTagID, version)

		downloadURL := "/" + filepath.ToSlash(dlRelPath)
		ghRepo := os.Getenv("GITHUB_REPOSITORY")
		if ghRepo == "" {
			ghRepo = "okf-memory/registry"
		}
		if os.Getenv("LOCAL_DOWNLOADS") != "1" {
			downloadURL = fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", ghRepo, releaseTag, fileNameWithVersion)
		}
		manifestURL := fmt.Sprintf("/bundles/%s.json", bundleID)

		// Create bundle manifest
		bManifest := BundleManifest{
			ID:          bundleID,
			Version:     version,
			OKFVersion:  okfVersion,
			License:     license,
			Status:      status,
			IsStale:     isStale,
			Hash:        hash,
			DownloadURL: downloadURL,
		}
		mBytes, _ := json.MarshalIndent(bManifest, "", "  ")

		// 1. Latest manifest: bundles/<bundle-id>.json
		manifestPathLatest := filepath.Join(publicDir, "bundles", rel+".json")
		if err := os.MkdirAll(filepath.Dir(manifestPathLatest), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(manifestPathLatest, mBytes, 0o644); err != nil {
			return err
		}

		// 2. Version-pinned manifest: bundles/<bundle-id>-<version>.json
		var manifestPathVersioned string
		if parentDir == "." {
			manifestPathVersioned = filepath.Join(publicDir, "bundles", fmt.Sprintf("%s-%s.json", baseName, version))
		} else {
			manifestPathVersioned = filepath.Join(publicDir, "bundles", parentDir, fmt.Sprintf("%s-%s.json", baseName, version))
		}
		if err := os.WriteFile(manifestPathVersioned, mBytes, 0o644); err != nil {
			return err
		}

		indexItems = append(indexItems, BundleItem{
			ID:          bundleID,
			Version:     version,
			OKFVersion:  okfVersion,
			License:     license,
			Status:      status,
			IsStale:     isStale,
			Tier:        "official",
			Title:       title,
			Description: desc,
			Hash:        hash,
			ManifestURL: manifestURL,
			DownloadURL: downloadURL,
			InstallCmd:  fmt.Sprintf("okf pull %s", bundleID),
		})

		fmt.Printf("✓ Packaged %-25s -> %s (%s)\n", bundleID, dlRelPath, hash[:17]+"...")
		return filepath.SkipDir
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Build error: %v\n", err)
		os.Exit(1)
	}

	officialCount := len(indexItems)
	communityCount := 0

	// Merge community bundles from data/community.json if available
	type CommunityItem struct {
		ID          string `json:"id"`
		Owner       string `json:"owner"`
		Repo        string `json:"repo"`
		Version     string `json:"version"`
		OKFVersion  string `json:"okf_version,omitempty"`
		License     string `json:"license,omitempty"`
		Status      string `json:"status,omitempty"`
		IsStale     bool   `json:"is_stale,omitempty"`
		Tier        string `json:"tier"`
		Title       string `json:"title"`
		Description string `json:"description"`
		RepoURL     string `json:"repo_url"`
		UpdatedAt   string `json:"updated_at"`
	}

	if commRaw, err := os.ReadFile(filepath.Join("data", "community.json")); err == nil {
		var commItems []CommunityItem
		if err := json.Unmarshal(commRaw, &commItems); err == nil {
			for _, c := range commItems {
				commStatus := c.Status
				if commStatus == "" {
					commStatus = "stable"
				}
				commTag := c.Version
				if commTag != "" && !strings.HasPrefix(commTag, "v") {
					commTag = "v" + commTag
				}
				commInstallCmd := fmt.Sprintf("okf pull %s", c.ID)
				if commTag != "" {
					commInstallCmd = fmt.Sprintf("okf pull %s@%s", c.ID, commTag)
				}

				indexItems = append(indexItems, BundleItem{
					ID:          c.ID,
					Version:     c.Version,
					OKFVersion:  c.OKFVersion,
					License:     c.License,
					Status:      commStatus,
					IsStale:     c.IsStale,
					Tier:        "community",
					Title:       c.Title,
					Description: c.Description,
					DownloadURL: c.RepoURL,
					InstallCmd:  commInstallCmd,
				})
			}
			communityCount = len(commItems)
		}
	}

	idxManifest := IndexManifest{
		Version:   1,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
		Bundles:   indexItems,
	}

	idxBytes, _ := json.MarshalIndent(idxManifest, "", "  ")
	if err := os.WriteFile(filepath.Join(publicDir, "index.json"), idxBytes, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing index.json: %v\n", err)
		os.Exit(1)
	}

	// Render web UI from src/index.html Go template
	tmpl, err := template.ParseFiles("src/index.html")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing template src/index.html: %v\n", err)
		os.Exit(1)
	}

	outHtml, err := os.Create(filepath.Join(publicDir, "index.html"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating public/index.html: %v\n", err)
		os.Exit(1)
	}

	templateData := struct {
		Bundles        []BundleItem
		OfficialCount  int
		CommunityCount int
		UpdatedAt      string
	}{
		Bundles:        indexItems,
		OfficialCount:  officialCount,
		CommunityCount: communityCount,
		UpdatedAt:      idxManifest.UpdatedAt,
	}

	if err := tmpl.Execute(outHtml, templateData); err != nil {
		_ = outHtml.Close()
		fmt.Fprintf(os.Stderr, "Error executing template: %v\n", err)
		os.Exit(1)
	}
	if err := outHtml.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Error closing public/index.html: %v\n", err)
		os.Exit(1)
	}

	if css, err := os.ReadFile("src/style.css"); err == nil {
		_ = os.WriteFile(filepath.Join(publicDir, "style.css"), css, 0o644)
	}
	if js, err := os.ReadFile("src/app.js"); err == nil {
		_ = os.WriteFile(filepath.Join(publicDir, "app.js"), js, 0o644)
	}
	fmt.Printf("Generated web UI at public/index.html with %d bundle card(s).\n", len(indexItems))

	fmt.Printf("\nGenerated registry index at public/index.json with %d bundle(s).\n", len(indexItems))
}

func packageTgz(srcDir, destTgz string) (string, error) {
	outFile, err := os.Create(destTgz)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = outFile.Close()
	}()

	hasher := sha256.New()
	mw := io.MultiWriter(outFile, hasher)

	gw := gzip.NewWriter(mw)
	gw.ModTime = time.Unix(0, 0).UTC()
	gw.OS = 255 // Neutral/unknown OS

	tw := tar.NewWriter(gw)

	// Collect all relative paths
	var relPaths []string
	err = filepath.Walk(srcDir, func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(srcDir, p)
		if err != nil || rel == "." {
			return nil
		}
		relPaths = append(relPaths, rel)
		return nil
	})
	if err != nil {
		return "", err
	}

	// Lexicographical sort for deterministic order across all environments
	sort.Strings(relPaths)

	for _, rel := range relPaths {
		p := filepath.Join(srcDir, rel)
		fi, err := os.Lstat(p)
		if err != nil {
			return "", err
		}

		hdr, err := tar.FileInfoHeader(fi, "")
		if err != nil {
			return "", err
		}

		hdr.Name = filepath.ToSlash(rel)
		hdr.ModTime = time.Unix(0, 0).UTC()
		hdr.Uid = 0
		hdr.Gid = 0
		hdr.Uname = ""
		hdr.Gname = ""
		if fi.IsDir() {
			hdr.Mode = 0o755
		} else {
			hdr.Mode = 0o644
		}

		if err := tw.WriteHeader(hdr); err != nil {
			return "", err
		}

		if fi.Mode().IsRegular() {
			f, err := os.Open(p)
			if err != nil {
				return "", err
			}
			if _, err := io.Copy(tw, f); err != nil {
				_ = f.Close()
				return "", err
			}
			if err := f.Close(); err != nil {
				return "", err
			}
		}
	}

	if err := tw.Close(); err != nil {
		return "", err
	}
	if err := gw.Close(); err != nil {
		return "", err
	}
	if err := outFile.Close(); err != nil {
		return "", err
	}

	return fmt.Sprintf("sha256:%x", hasher.Sum(nil)), nil
}
