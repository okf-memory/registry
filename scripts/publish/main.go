package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type IndexItem struct {
	ID          string `json:"id"`
	Version     string `json:"version"`
	Tier        string `json:"tier,omitempty"`
	License     string `json:"license,omitempty"`
	Hash        string `json:"hash"`
	ManifestURL string `json:"manifest_url"`
	DownloadURL string `json:"download_url"`
}

type IndexManifest struct {
	Version int         `json:"version"`
	Bundles []IndexItem `json:"bundles"`
}

func main() {
	idxData, err := os.ReadFile("public/index.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading public/index.json: %v\n", err)
		os.Exit(1)
	}

	var manifest IndexManifest
	if err := json.Unmarshal(idxData, &manifest); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing public/index.json: %v\n", err)
		os.Exit(1)
	}

	for _, b := range manifest.Bundles {
		if b.Tier == "community" {
			// Community bundles are decentralized and hosted on the author's own GitHub repository.
			continue
		}

		safeID := strings.ReplaceAll(b.ID, "/", "-")
		tag := fmt.Sprintf("%s-v%s", safeID, b.Version)
		fileName := fmt.Sprintf("%s-%s.tgz", filepath.Base(b.ID), b.Version)

		var localPath string
		parent := filepath.Dir(b.ID)
		if parent == "." {
			localPath = filepath.Join("public", "downloads", fileName)
		} else {
			localPath = filepath.Join("public", "downloads", parent, fileName)
		}

		if _, err := os.Stat(localPath); os.IsNotExist(err) {
			fmt.Printf("⚠️  Archive not found locally: %s, skipping\n", localPath)
			continue
		}

		// Check if release tag already exists
		viewCmd := exec.Command("gh", "release", "view", tag)
		if err := viewCmd.Run(); err == nil {
			fmt.Printf("✓ Release %s already exists\n", tag)
			continue
		}

		// Create release and upload asset
		title := fmt.Sprintf("%s v%s", b.ID, b.Version)
		licenseNote := ""
		if b.License != "" {
			licenseNote = fmt.Sprintf("* **License:** `%s`\n", b.License)
		}
		notes := fmt.Sprintf("OKF Memory Bundle `%s` version `%s`.\n\n%s* **SHA-256:** `%s`\n* **Manifest:** `/bundles/%s-%s.json`", b.ID, b.Version, licenseNote, b.Hash, b.ID, b.Version)
		createCmd := exec.Command("gh", "release", "create", tag, localPath, "--title", title, "--notes", notes)
		createCmd.Stdout = os.Stdout
		createCmd.Stderr = os.Stderr

		fmt.Printf("📦 Creating release %s with asset %s...\n", tag, localPath)
		if err := createCmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create release %s: %v\n", tag, err)
		} else {
			fmt.Printf("✓ Successfully created release %s\n", tag)
		}
	}
}
