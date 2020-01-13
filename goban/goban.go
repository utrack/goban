package goban

import (
	"bufio"
	"flag"
	"go/ast"
	"go/types"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/types/typeutil"
)

func Analyzer() *analysis.Analyzer {
	ret := &analysis.Analyzer{
		Name:     "goban",
		Doc:      "prohibits usage of certain symbols",
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
	flg := flag.NewFlagSet("", flag.ExitOnError)
	cfgPath := flg.String("cfg", ".goban.cfg", "path to newline-delimited list of banned symbols")
	ret.Flags = *flg
	ret.Run = run(cfgPath)
	return ret
}

var (
	bannedPatterns map[string]string // map (symbol name)->(comment)
	bpMtx          sync.Mutex
)

func getBannedPtsMap(path string) map[string]string {
	bpMtx.Lock()
	defer bpMtx.Unlock()
	if bannedPatterns == nil {
		err := loadTrie(path)
		if err != nil {
			log.Fatal(err)
		}
	}

	return bannedPatterns
}

func run(cfgPath *string) func(*analysis.Pass) (interface{}, error) {
	return func(pass *analysis.Pass) (interface{}, error) {
		patterns := getBannedPtsMap(*cfgPath)
		inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
		nodeFilter := []ast.Node{
			(*ast.CallExpr)(nil),
		}

		inspect.Preorder(nodeFilter, func(n ast.Node) {
			call := n.(*ast.CallExpr)
			fn, _ := typeutil.Callee(pass.TypesInfo, call).(*types.Func)
			if fn == nil {
				return
			}
			if comment, ok := patterns[fn.FullName()]; ok {
				if comment != "" {
					comment = " - " + comment
				}
				pass.Reportf(call.Pos(), "func %v is banned%v", fn.FullName(), comment)
			}
		})
		return nil, nil
	}
}

// TODO: create a real trie to match symbols by wildcards
func loadTrie(path string) error {
	bannedPatterns = map[string]string{}
	f, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "when opening goban config file")
	}

	bio := bufio.NewReader(f)
	for {
		line, err := bio.ReadString('\n')
		if err == io.EOF {
			return nil
		}

		var comment string
		commentPointIdx := strings.Index(line, "#")
		if commentPointIdx > -1 {
			comment = line[commentPointIdx+1:]
			comment = strings.Trim(comment, " \n\t\r")
			line = line[:commentPointIdx]
		}

		line = strings.Split(line, "#")[0]
		line = strings.Trim(line, " \n\t\r")
		if line == "" {
			continue
		}
		bannedPatterns[line] = comment
	}
}
