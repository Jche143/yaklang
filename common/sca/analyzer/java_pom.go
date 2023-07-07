package analyzer

import (
	"path/filepath"

	"github.com/aquasecurity/go-dep-parser/pkg/java/pom"
	"github.com/yaklang/yaklang/common/sca/types"
)

const (
	TypPom TypAnalyzer = "pom-lang"

	MavenPom = "pom.xml"

	statusPom int = 1
)

func init() {
	RegisterAnalyzer(TypPom, NewJavaPomAnalyzer())
}

type pomAnalyzer struct{}

func NewJavaPomAnalyzer() *pomAnalyzer {
	return &pomAnalyzer{}
}

func (a pomAnalyzer) Analyze(afi AnalyzeFileInfo) ([]types.Package, error) {
	fi := afi.self
	switch fi.matchStatus {
	case statusPom:
		p := pom.NewParser(fi.path, pom.WithOffline(false))
		pkgs, err := ParseLanguageConfiguration(fi, p)
		if err != nil {
			return nil, err
		}
		return pkgs, nil
	}
	return nil, nil
}

func (a pomAnalyzer) Match(info MatchInfo) int {
	if filepath.Base(info.path) == MavenPom {
		return statusPom
	}
	return 0
}