package main

import (
    "fmt"
    "os"

    tree_sitter "github.com/tree-sitter/go-tree-sitter"
    tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"

    "golze/cfg"
)

func printErr(msg string, args... any) {
    err := fmt.Errorf(msg, args...)
    fmt.Println(err.Error())
}

func printErrAndDie(msg string, args... any) {
    printErr(msg, args...)
    os.Exit(1)
}

func parse(parser *tree_sitter.Parser, code []byte) *tree_sitter.Tree {

    tree := parser.Parse([]byte(code), nil)

    return tree
}

func getParser() *tree_sitter.Parser {
    parser := tree_sitter.NewParser()
    parser.SetLanguage(tree_sitter.NewLanguage(tree_sitter_go.Language()))
    return parser
}

func main() {
    if len(os.Args) < 2 {
        printErrAndDie("please provide a source file to analyze")
    }
    filename := os.Args[1]
    code, err := os.ReadFile(filename)
    if err != nil {
        printErrAndDie("unable to read file %s: %w", filename, err)
    }

    parser := getParser()

    //parse("package test\n\nimport \"fmt\"\n\nfunc main() { fmt.Println(\"HELLO\") }")
    sourceTree := parse(parser, code)
    defer sourceTree.Close()

    fmt.Println(sourceTree.RootNode().ToSexp())

    fileCFG, err := cfg.NewCFG(sourceTree.RootNode())
    if err != nil {
        printErrAndDie("failed to create CFG: %w", err)
    }

    fnCFGs := []*cfg.CFG{}

    funcsToDo := []cfg.FuncDef{}
    funcsToDo = append(funcsToDo, fileCFG.FuncDefs...)
    for _, funcDef := range funcsToDo {
        funcCFG, err := cfg.NewCFG(funcDef.AST)
        if err != nil {
            funcName := "<anonymous>"
            if funcDef.Id != nil {
                funcName = funcDef.Id.Name
            }
            fnCFGs = append(fnCFGs, funcCFG)
            printErr("failed to create CFG for function %s", funcName)
            continue
        }
        funcsToDo = append(funcsToDo, funcCFG.FuncDefs...)
    }

    fmt.Println(fileCFG)
    for fnCFG := range fnCFGs {
        fmt.Println(fnCFG)
    }
}
