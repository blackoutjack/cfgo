package main

import (
    "flag"
    "fmt"
    "os"

    tree_sitter "github.com/tree-sitter/go-tree-sitter"
    tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"

    "cfgo/cfg"
    log "cfgo/util"
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
        "print the parse tree but do not generate CFGs")
    noAST := flag.Bool("no-ast", false, "print the CFG but not the parse tree")

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
        log.PrintFullTree(code, sourceTree.RootNode().Walk(), false)
    }

    if *noCFG { return }

    fileCFG, err := cfg.NewCFG(sourceTree.RootNode(), code)
    if err != nil {
        log.PrintErr("error(s) while creating CFG: %w", err)
    }

    fmt.Println(fileCFG)
    for _, funcDef := range fileCFG.FuncDefs {
        fmt.Printf("Function: %s\n", funcDef.ID.Name) 
        fmt.Println(funcDef.CFG)
    }
}
