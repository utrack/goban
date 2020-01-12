package goban

import (
	"bufio"
	"flag"
	"go/ast"
	"go/types"
	"io"
	"log"
	"os"

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
	bannedPatterns map[string]struct{}
)

func run(cfgPath *string) func(*analysis.Pass) (interface{}, error) {
	return func(pass *analysis.Pass) (interface{}, error) {
		if bannedPatterns == nil {
			err := loadTrie(*cfgPath)
			if err != nil {
				log.Fatal(err)
			}
		}
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
			if _, ok := bannedPatterns[fn.FullName()]; ok {
				pass.Reportf(fn.Pos(), "func %v is banned", fn.FullName())
			}
		})
		return nil, nil
	}
}

func loadTrie(path string) error {
	bannedPatterns = map[string]struct{}{}
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
		bannedPatterns[line] = struct{}{}
	}
}
