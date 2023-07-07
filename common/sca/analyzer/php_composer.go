package analyzer

import (
	"encoding/json"
	"io"
	"path/filepath"
	"strings"

	"github.com/aquasecurity/go-dep-parser/pkg/php/composer"
	"github.com/yaklang/yaklang/common/sca/types"
	"github.com/yaklang/yaklang/common/utils"
	"golang.org/x/exp/slices"
)

const (
	TypComposer TypAnalyzer = "composer-lang"

	phpLockFile = "composer.lock"
	phpJsonFile = "composer.json"

	statusComposerLock int = 1
	statusComposerJson int = 2
)

func init() {
	RegisterAnalyzer(TypComposer, NewPHPComposerAnalyzer())
}

type composerAnalyzer struct{}

func NewPHPComposerAnalyzer() *composerAnalyzer {
	return &composerAnalyzer{}
}

type composerJson struct {
	Require map[string]string `json:"require"`
}

func (a composerAnalyzer) Analyze(afi AnalyzeFileInfo) ([]types.Package, error) {
	fi := afi.self
	switch fi.matchStatus {
	case statusComposerLock:
		// parse composer lock file
		lockParser := composer.NewParser()
		pkgs, err := ParseLanguageConfiguration(fi, lockParser)
		if err != nil {
			return nil, err
		}

		// parse composer json file
		var p map[string]string
		jsonPath := filepath.Join(filepath.Dir(fi.path), "composer.json")
		if jsonFi, ok := afi.matchedFileInfos[jsonPath]; ok {
			p, err = a.parseComposerJson(jsonFi.f)
			if err != nil {
				p = nil
			}
		}
		if p != nil {
			for i, pkg := range pkgs {
				if _, ok := p[pkg.Name]; !ok {
					pkgs[i].Indirect = true
				}
			}
		}
		return pkgs, nil
	}
	return nil, nil
}
func (a composerAnalyzer) parseComposerJson(f io.Reader) (map[string]string, error) {
	jsonFile := composerJson{}
	err := json.NewDecoder(f).Decode(&jsonFile)
	if err != nil {
		return nil, utils.Errorf("json decode error: %v", err)
	}
	return jsonFile.Require, nil
}

func (a composerAnalyzer) Match(info MatchInfo) int {
	fileName := filepath.Base(info.path)
	// Skip `composer.lock` inside `vendor` folder
	if slices.Contains(strings.Split(info.path, "/"), "vendor") {
		return 0
	}
	if fileName == phpJsonFile {
		return statusComposerJson
	}
	if fileName == phpLockFile {
		return statusComposerLock
	}
	return 0
}