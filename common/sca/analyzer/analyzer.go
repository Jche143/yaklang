package analyzer

import (
	"bufio"
	"io"
	"io/fs"
	"os"
	"strings"
	"sync"

	godeptypes "github.com/aquasecurity/go-dep-parser/pkg/types"

	"github.com/samber/lo"
	"github.com/yaklang/yaklang/common/sca/types"
	"github.com/yaklang/yaklang/common/utils"
)

const (
	headerSize = 4

	AllMode ScanMode = 0
	PkgMode          = 1 << (iota - 1)
	LanguageMode
)

var (
	analyzers = make(map[TypAnalyzer]Analyzer, 0)
)

type TypAnalyzer string
type ScanMode int
type Analyzer interface {
	Analyze(AnalyzeFileInfo) ([]types.Package, error)
	Match(MatchInfo) int
}

type AnalyzeFileInfo struct {
	path        string
	f           *os.File
	matchStatus int
}

type MatchInfo struct {
	path   string
	fi     fs.FileInfo
	header []byte
}

type Task struct {
	fileInfo AnalyzeFileInfo
	a        Analyzer
}

type AnalyzerGroup struct {
	analyzers []Analyzer

	// consume
	ch         chan Task
	numWorkers int

	// return
	pkgs []types.Package
	err  error

	// scanned file
	scannedFiles map[string]struct{}
}

func RegisterAnalyzer(typ TypAnalyzer, a Analyzer) {
	if _, ok := analyzers[typ]; ok {
		return
	}
	analyzers[typ] = a
}

func FilterAnalyzer(mode ScanMode) []Analyzer {
	ret := make([]Analyzer, 0, len(analyzers))
	if mode == AllMode {
		return lo.MapToSlice(analyzers, func(_ TypAnalyzer, a Analyzer) Analyzer {
			return a
		})
	}

	for analyzerName, a := range analyzers {
		// filter by ScanMode
		if mode&PkgMode == PkgMode {
			if strings.HasSuffix(string(analyzerName), "-pkg") {
				ret = append(ret, a)
				continue
			}
		}
		if mode&LanguageMode == LanguageMode {
			if strings.HasSuffix(string(analyzerName), "-lang") {
				ret = append(ret, a)
				continue
			}
		}

	}
	return ret
}

func NewAnalyzerGroup(numWorkers int, scanMode ScanMode) *AnalyzerGroup {
	return &AnalyzerGroup{
		ch:           make(chan Task),
		numWorkers:   numWorkers,
		scannedFiles: make(map[string]struct{}),
		analyzers:    FilterAnalyzer(scanMode),
	}
}

func (ag *AnalyzerGroup) Error() error {
	return ag.err
}

func (ag *AnalyzerGroup) Packages() []types.Package {
	return lo.Uniq(ag.pkgs)
}

func (ag *AnalyzerGroup) Append(a ...Analyzer) {
	ag.analyzers = append(ag.analyzers, a...)
}

func (ag *AnalyzerGroup) Consume(wg *sync.WaitGroup) {
	wg.Add(ag.numWorkers)

	for i := 0; i < ag.numWorkers; i++ {
		go func() {
			defer wg.Done()
			for task := range ag.ch {
				defer func() {
					name := task.fileInfo.f.Name()
					task.fileInfo.f.Close()
					os.Remove(name)
				}()
				pkgs, err := task.a.Analyze(task.fileInfo)
				if err != nil {
					ag.err = err
					return
				}
				ag.pkgs = append(ag.pkgs, pkgs...)
			}
		}()
	}
}

func (ag *AnalyzerGroup) Close() {
	close(ag.ch)
}

// write
func (ag *AnalyzerGroup) Analyze(path string, fi fs.FileInfo, r io.Reader) error {
	var (
		header []byte
		err    error
	)
	br := bufio.NewReader(r)

	for _, a := range ag.analyzers {
		// if scanned, skip
		if _, ok := ag.scannedFiles[path]; ok {
			continue
		}

		if fi.Mode().IsRegular() {
			header, err = br.Peek(headerSize)
			if err != nil && err != io.EOF {
				return utils.Errorf("read file header error: %v", err)
			}
		}

		matchStatus := a.Match(MatchInfo{
			path:   path,
			fi:     fi,
			header: header,
		})

		if matchStatus == 0 {
			continue
		}
		// match type > 0 mean matched and need to analyze

		// save
		f, err := os.CreateTemp("", "fanal-file-*")
		if err != nil {
			return utils.Errorf("failed to create a temporary file for analyzer")
		}

		if _, err := io.Copy(f, br); err != nil {
			return utils.Errorf("failed to copy the file: %v", err)
		}
		f.Seek(0, 0)

		// send
		task := Task{
			fileInfo: AnalyzeFileInfo{
				path:        path,
				f:           f,
				matchStatus: matchStatus,
			},
			a: a,
		}
		ag.ch <- task

		// add to scanned files
		ag.scannedFiles[path] = struct{}{}
	}
	return nil
}

func ParseLanguageConfiguration(fi AnalyzeFileInfo, parser godeptypes.Parser) ([]types.Package, error) {
	parsedLibs, _, err := parser.Parse(fi.f)
	if err != nil {
		return nil, err
	}

	pkgs := lo.Map(parsedLibs, func(lib godeptypes.Library, index int) types.Package {
		return types.Package{
			Name:    lib.Name,
			Version: lib.Version,
		}
	})
	return pkgs, nil
}
