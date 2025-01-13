package main

import (
	"cloak/pkg/subs"
	"fmt"
	"path/filepath"
)


type ElasticModule struct {
	ignoreList   []string
	renamePaths  bool
	replacePairs []SearchReplacePair
}

func NewElasticModule() *ElasticModule {
	return &ElasticModule{
		ignoreList:  []string{".git", ".github", "docs", "vendor"},
		renamePaths: true,
		replacePairs: []SearchReplacePair{
			{search: "IfconfigReq", replace: "Frank"},
			{search: "ImpersonateReq", replace: "Steve"},
			{search: "InvokeMigrateReq", replace: "Paul"},
			{search: "RevToSelfReq", replace: "Gerald"},
			{search: "ScreenshotReq", replace: "Smith"},
			{search: "SideloadReq", replace: "Alex"},
			{search: "InvokeSpawnDllReq", replace: "Derek"},
			{search: "NetstatReq", replace: "Grant"},
			{search: "httpSessionInit", replace: "Robert"},
			{search: "screenshotRequested", replace: "Wayne"},
			{search: "RegistryReadReq", replace: "Roberto"},
			{search: "RequestResend", replace: "Frankie"},
			{search: "GetPrivInfo", replace: "Wallace"},
			{search: "-NoExit", replace: "-nOExIt"},
		},
	}
}

func (m *ElasticModule) Name() string {
	return "Elastic"
}

func (m *ElasticModule) Run(config *Config, verbose bool) error {
	startPath := filepath.Join(config.RunDir, "sliver")

	// Start the recursive search and replace
	var err error
	for _, pair := range m.replacePairs {

		err = subs.SearchAndReplace(startPath, pair.search, pair.replace, m.ignoreList, verbose)
		if err != nil {
			return fmt.Errorf("[Elastic] [SearchAndReplace] error during execution: %v", err)
		}

		err = subs.SearchAndRenameFiles(startPath, pair.search, pair.replace, m.ignoreList, verbose)
		if err != nil {
			return fmt.Errorf("[Elastic] [SearchAndRenameFiles] error during execution: %v", err)
		}

		err = subs.SearchAndRenameDirectories(startPath, pair.search, pair.replace, m.ignoreList, verbose)
		if err != nil {
			return fmt.Errorf("[Elastic] [SearchAndRenameDirectories] error during execution: %v", err)
		}
	}

	return nil
}
