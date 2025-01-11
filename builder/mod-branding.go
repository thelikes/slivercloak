package main

import (
	"cloak/pkg/subs"
	"fmt"
	"path/filepath"
)

type SearchReplacePair struct {
	search  string
	replace string
}

type BrandingModule struct {
	ignoreList   []string
	renamePaths  bool
	replacePairs []SearchReplacePair
}

func NewBrandingModule() *BrandingModule {
	return &BrandingModule{
		ignoreList:  []string{".git", ".github", "docs", "vendor"},
		renamePaths: true,
		replacePairs: []SearchReplacePair{
			{search: "sliver", replace: "gunner"},
			{search: "Sliver", replace: "Gunner"},
			{search: "SLIVER", replace: "GUNNER"},
			{search: "beacon", replace: "lazer"},
			{search: "Beacon", replace: "Lazer"},
			{search: "BEACON", replace: "LAZER"},
			{search: "bishopfox", replace: "knightbruce"},
			{search: "BishopFox", replace: "KnightBruce"},
		},
	}
}

func (m *BrandingModule) Name() string {
	return "branding"
}

func (m *BrandingModule) Run(config *Config, verbose bool) error {
	startPath := filepath.Join(config.RunDir, "sliver")

	// Start the recursive search and replace
	var err error
	for _, pair := range m.replacePairs {

		err = subs.SearchAndReplace(startPath, pair.search, pair.replace, m.ignoreList, verbose)
		if err != nil {
			return fmt.Errorf("[branding] [SearchAndReplace] error during execution: %v", err)
		}

		err = subs.SearchAndRenameFiles(startPath, pair.search, pair.replace, m.ignoreList, verbose)
		if err != nil {
			return fmt.Errorf("[branding] [SearchAndRenameFiles] error during execution: %v", err)
		}

		err = subs.SearchAndRenameDirectories(startPath, pair.search, pair.replace, m.ignoreList, verbose)
		if err != nil {
			return fmt.Errorf("[branding] [SearchAndRenameDirectories] error during execution: %v", err)
		}
	}

	return nil
}
