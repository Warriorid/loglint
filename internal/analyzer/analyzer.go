package analyzer

import (
	"go/ast"
	"go/token"
	"strings"
	"unicode/utf8"

	"github.com/loglint/internal/rules"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const analyzerName = "loglint"
const analyzerDoc = `loglint checks log messages for style and security issues.

Rules:
  1. Log messages must start with a lowercase letter.
  2. Log messages must be in English only.
  3. Log messages must not contain special characters or emojis.
  4. Log messages must not contain potentially sensitive data.
`

func New(cfg *Config) *analysis.Analyzer {
	if cfg == nil {
		cfg = DefaultConfig()
	}
	a := &loglintAnalyzer{cfg: cfg}
	return &analysis.Analyzer{
		Name:     analyzerName,
		Doc:      analyzerDoc,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run:      a.run,
	}
}

var Analyzer = New(DefaultConfig())

type loglintAnalyzer struct {
	cfg *Config
}

var loggerPackages = map[string]bool{
	"log":             true,
	"log/slog":        true,
	"go.uber.org/zap": true,
}

var logMethods = map[string]bool{
	"Print":    true,
	"Printf":   true,
	"Println":  true,
	"Fatal":    true,
	"Fatalln":  true,
	"Panic":    true,
	"Panicln":  true,
	"Debug":    true,
	"Info":     true,
	"Warn":     true,
	"Error":    true,
	"Log":      true,
	"DebugCtx": true,
	"InfoCtx":  true,
	"WarnCtx":  true,
	"ErrorCtx": true,
	"DPanic":   true,
	"Debugf":   true,
	"Infof":    true,
	"Warnf":    true,
	"Errorf":   true,
	"Fatalf":   true,
	"Panicf":   true,
	"DPanicf":  true,
	"Debugw":   true,
	"Infow":    true,
	"Warnw":    true,
	"Errorw":   true,
	"Fatalw":   true,
	"Panicw":   true,
	"DPanicw":  true,
}

func (a *loglintAnalyzer) run(pass *analysis.Pass) (interface{}, error) {
	loggerLocalNames := collectLoggerImports(pass)
	if len(loggerLocalNames) == 0 {
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return
		}

		msgArg, methodName, isLogger := extractLogCall(call, loggerLocalNames)
		if !isLogger {
			return
		}
		_ = methodName

		msg, isDynamic := extractStringValue(msgArg)
		if isDynamic {
			if a.cfg.CheckSensitiveData {
				if hasSensitiveInExpr(msgArg, a.cfg.SensitiveKeywords) {
					pass.Reportf(call.Pos(), "log message may contain sensitive data")
				}
			}
			return
		}

		a.checkMessage(pass, call.Pos(), msgArg, msg)
	})

	return nil, nil
}

func (a *loglintAnalyzer) checkMessage(pass *analysis.Pass, pos token.Pos, msgArg ast.Expr, msg string) {
	if a.cfg.CheckLowercase && !rules.CheckLowercase(msg) {
		r, size := utf8.DecodeRuneInString(msg)
		fix := strings.ToLower(string(r)) + msg[size:]
		pass.Report(analysis.Diagnostic{
			Pos:     msgArg.Pos(),
			End:     msgArg.End(),
			Message: "log message should start with a lowercase letter",
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: "Convert first letter to lowercase",
					TextEdits: []analysis.TextEdit{
						{
							Pos:     msgArg.Pos(),
							End:     msgArg.End(),
							NewText: []byte(`"` + fix + `"`),
						},
					},
				},
			},
		})
	}

	if a.cfg.CheckEnglishOnly && !rules.CheckEnglishOnly(msg) {
		pass.Reportf(pos, "log message should be in English only")
	}

	if a.cfg.CheckSpecialChars && !rules.CheckNoSpecialChars(msg) {
		pass.Reportf(pos, "log message should not contain special characters or emojis")
	}

	if a.cfg.CheckSensitiveData && !rules.CheckNoSensitiveData(msg, a.cfg.SensitiveKeywords) {
		pass.Reportf(pos, "log message may contain sensitive data")
	}
}

func collectLoggerImports(pass *analysis.Pass) map[string]string {
	result := make(map[string]string)
	for _, file := range pass.Files {
		for _, imp := range file.Imports {
			path := strings.Trim(imp.Path.Value, `"`)
			if !loggerPackages[path] {
				continue
			}
			var localName string
			if imp.Name != nil {
				localName = imp.Name.Name
			} else {
				parts := strings.Split(path, "/")
				localName = parts[len(parts)-1]
			}
			result[localName] = path
		}
	}
	return result
}

func extractLogCall(call *ast.CallExpr, loggerLocalNames map[string]string) (ast.Expr, string, bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, "", false
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return nil, "", false
	}

	methodName := sel.Sel.Name
	if !logMethods[methodName] {
		return nil, "", false
	}

	importPath, isLoggerPkg := loggerLocalNames[ident.Name]
	if !isLoggerPkg {
		return nil, "", false
	}

	if len(call.Args) == 0 {
		return nil, "", false
	}

	msgIndex := 0
	if importPath == "log/slog" && methodName == "Log" {
		msgIndex = 2
	} else if importPath == "log/slog" && strings.HasSuffix(methodName, "Ctx") {
		msgIndex = 1
	}

	if msgIndex >= len(call.Args) {
		return nil, "", false
	}

	return call.Args[msgIndex], methodName, true
}

func extractStringValue(arg ast.Expr) (string, bool) {
	lit, ok := arg.(*ast.BasicLit)
	if !ok {
		return "", true
	}
	if lit.Kind.String() != "STRING" {
		return "", true
	}
	val := lit.Value
	if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
		return val[1 : len(val)-1], false
	}
	if len(val) >= 2 && val[0] == '`' && val[len(val)-1] == '`' {
		return val[1 : len(val)-1], false
	}
	return val, false
}

func hasSensitiveInExpr(expr ast.Expr, keywords []string) bool {
	found := false
	ast.Inspect(expr, func(n ast.Node) bool {
		if found {
			return false
		}
		switch v := n.(type) {
		case *ast.BasicLit:
			val := strings.Trim(v.Value, "\"`")
			if !rules.CheckNoSensitiveData(val, keywords) {
				found = true
			}
		case *ast.Ident:
			if !rules.CheckNoSensitiveData(v.Name, keywords) {
				found = true
			}
		}
		return !found
	})
	return found
}
