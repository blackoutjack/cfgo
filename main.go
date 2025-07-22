package main

import (
    "flag"
    "fmt"
    "os"

    tree_sitter "github.com/tree-sitter/go-tree-sitter"
    tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"

    "golze/cfg"
    log "golze/util"
)

func parse(parser *tree_sitter.Parser, code []byte) *tree_sitter.Tree {

    tree := parser.Parse(code, nil)

    return tree
}

func getParser() *tree_sitter.Parser {
    parser := tree_sitter.NewParser()
    parser.SetLanguage(tree_sitter.NewLanguage(tree_sitter_go.Language()))
    return parser
}

func main() {

    noCFG := flag.Bool("no-cfg", false,
        "print the AST but do not generate CFGs")
    noAST := flag.Bool("no-ast", false, "print the CFG but not AST")

    flag.Parse()

    if flag.NArg() < 1 {
        log.PrintErrAndDie("please provide a source file to analyze")
    }
    filename := flag.Arg(0)
    code, err := os.ReadFile(filename)
    if err != nil {
        log.PrintErrAndDie("unable to read file %s: %w", filename, err)
    }

    parser := getParser()

    sourceTree := parse(parser, code)
    defer sourceTree.Close()

    if !*noAST {
        fmt.Println(sourceTree.RootNode().ToSexp())
    }

    if *noCFG { return }

    fileCFG, err := cfg.NewCFG(sourceTree.RootNode(), code)
    if err != nil {
        log.PrintErrAndDie("failed to create CFG: %w", err)
    }

    fnCFGs := []*cfg.CFG{}

    funcsToDo := []cfg.FuncDef{}
    funcsToDo = append(funcsToDo, fileCFG.FuncDefs...)
    for _, funcDef := range funcsToDo {
        funcCFG, err := cfg.NewCFG(funcDef.AST, code)
        if err != nil {
            funcName := "<anonymous>"
            if funcDef.Id != nil {
                funcName = funcDef.Id.Name
            }
            fnCFGs = append(fnCFGs, funcCFG)
            log.PrintErr("failed to create CFG for function %s", funcName)
            continue
        }
        funcsToDo = append(funcsToDo, funcCFG.FuncDefs...)
    }

    fmt.Println(fileCFG)
    for fnCFG := range fnCFGs {
        fmt.Println(fnCFG)
    }
}
