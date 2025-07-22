package cfg

import (
    "fmt"
    "maps"
    "slices"
    "strings"
    "strconv"

    tree_sitter "github.com/tree-sitter/go-tree-sitter"

    log "golze/util"
)

var nodeIndex uint64 = 0

type Node struct {
    id uint64
}
func NewNode() *Node {
    node := Node{nodeIndex}
    nodeIndex++
    return &node
}
func (n *Node) String() string {
    if n == nil { return "nil" }
    return strconv.FormatUint(n.id, 10)
}
func nodeSort(a, b *Node) int {
    if a == nil && b == nil { return 0 }
    if a == nil { return -1 }
    if b == nil { return 1 }
    if a.id < b.id { return -1 }
    if a.id > b.id { return 1 }
    return 0
}

type Edge struct {
    Value *tree_sitter.Node
    Dest *Node
}

type Graph map[*Node][]Edge

type Decl struct {
    Name string
}

type FuncDef struct {
    Id *Decl
    AST *tree_sitter.Node
}
func (fd *FuncDef) String() string {
    if fd == nil { return "nil" }
    if fd.Id == nil {
        return "<anonymous>"
    }
    return fd.Id.Name
}

type CFG struct {
    Decls []Decl
    FuncDefs []FuncDef
    Entry *Node
    Exit *Node
    Graph Graph
}
func (g *CFG) String() string {
    builder := strings.Builder{}

    builder.WriteString("Declarations:")
    for _, decl := range g.Decls {
        builder.WriteString(fmt.Sprintf(" %s", decl.Name))
    }
    builder.WriteString("\n")

    builder.WriteString("Function definitions:")
    for _, fdef := range g.FuncDefs {
        builder.WriteString(fmt.Sprintf(" %s", fdef.String()))
    }
    builder.WriteString("\n")

    builder.WriteString("Entry: ");
    builder.WriteString(g.Entry.String())
    builder.WriteString("\n")

    builder.WriteString("Exit: ")
    builder.WriteString(g.Exit.String())
    builder.WriteString("\n")

    builder.WriteString("Graph:\n")

    // sort for testing stability
    nodes := maps.Keys(g.Graph)
    sortedNodes := slices.SortedStableFunc(nodes, nodeSort)
    for _, n := range sortedNodes {
        edges := g.Graph[n]
        srcStr := "nil"
        if n != nil {
            srcStr = n.String()
        }
        for _, e := range edges {
            valStr := "nil"
            if e.Value != nil {
                valStr = e.Value.ToSexp()
            }
            destStr := "nil"
            if e.Dest != nil {
                destStr = e.Dest.String()
            }
            builder.WriteString(fmt.Sprintf("  %s => %s: %s\n",
                srcStr, destStr, valStr))
        }
    }

    return builder.String()
}

func extractImportNameFromPath(
    code []byte,
    cursor *tree_sitter.TreeCursor) (string, error) {

    text := cursor.Node().Utf8Text(code)
    text = strings.Trim(text, "\"`")
    if len(text) == 0 {
        return "", fmt.Errorf("empty string found for import_spec path")
    }
    parts := strings.Split(text, "/")
    last := parts[len(parts)-1]
    // %%% check valid identifier

    return last, nil
}

func visitImportSpec(
    cfg *CFG,
    code []byte,
    source *Node,
    dest *Node,
    cursor *tree_sitter.TreeCursor) error {

    cfg.Graph[source] = append(cfg.Graph[source], Edge{cursor.Node(), dest})

    if !cursor.GotoFirstChild() {
        return fmt.Errorf("import_spec has no children")
    }
    defer cursor.GotoParent()

    isSplat := false
    for {
        fieldName := cursor.FieldName()
        switch fieldName {
        case "name":
            if !cursor.GotoFirstChild() {
                importName := cursor.Node().Utf8Text(code)
                if importName != "_" {
                    cfg.Decls = append(cfg.Decls, Decl{importName})
                }
                return nil
            }

            switch cursor.Node().Kind() {
            case ".":
                isSplat = true
                cursor.GotoParent()
            default:
                return fmt.Errorf(
                    "unexpected child of import_spec 'name': %s",
                    cursor.Node().Kind())
            }
        case "path":
            importName, err := extractImportNameFromPath(code, cursor)
            if err != nil {
                return fmt.Errorf(
                    "failed to extract import name from path: %w", err)
            }
            if isSplat {
                importName = fmt.Sprintf("%s.*", importName)
            }
            cfg.Decls = append(cfg.Decls, Decl{importName})
            return nil
        default:
            // skip over parens, other cruft
            log.PrintErr("skipping import_spec field: %s", fieldName)
            continue
        }
        if !cursor.GotoNextSibling() {
            break
        }
    }

    return fmt.Errorf("unable to determine declaration from import_spec")
}

func visitImportSpecList(cfg *CFG,
    code []byte,
    source *Node,
    dest *Node,
    cursor *tree_sitter.TreeCursor) error {

    importCount := cursor.Node().ChildCount()

    if !cursor.GotoFirstChild() {
        return fmt.Errorf("import_spec_list has no children")
    }
    defer cursor.GotoParent()

    // ignore the first child, it is a '('

    var cumulErr error = nil
    var childIndex uint = 1
    lastImport := importCount - 2
    for cursor.GotoNextSibling() {
        switch cursor.Node().Kind() {
        case "import_spec":
            var theDest *Node
            if childIndex == lastImport {
                theDest = dest
            } else {
                theDest = NewNode()
            }
            err := visitImportSpec(cfg, code, source, theDest, cursor)
            cumulErr = log.CombineErrors(cumulErr, err)
            source = theDest
        default:
            cumulErr = log.CombineErrors(cumulErr, fmt.Errorf(
                "unexpected child of import_spec_list: %s",
                cursor.Node().Kind()))
        }
        if childIndex == lastImport {
            break
        }
        childIndex += 1
    }

    return cumulErr
}

func visitImportDeclaration(cfg *CFG,
    code []byte,
    source *Node,
    dest *Node,
    cursor *tree_sitter.TreeCursor) error {

    // first child is the "import" keyword
    if !cursor.GotoFirstChild() {
        return fmt.Errorf("no children of import_declaration")
    }
    defer cursor.GotoParent()

    if !cursor.GotoNextSibling() {
        return fmt.Errorf("no interesting children of import_declaration")
    }

    switch cursor.Node().Kind() {
    case "import_spec_list":
        return visitImportSpecList(cfg, code, source, dest, cursor)
    case "import_spec":
        return visitImportSpec(cfg, code, source, dest, cursor)
    default:
        err := fmt.Errorf("unexpected child of import_declaration: %s", cursor.Node().Kind())
        cursor.GotoParent()
        return err
    }
}

func visitNode(cfg *CFG,
    code []byte,
    source *Node,
    dest *Node,
    cursor *tree_sitter.TreeCursor) error {

    switch cursor.Node().Kind() {
    case "source_file":

        if cursor.GotoFirstChild() {
            newDest := NewNode()
            err := visitNode(cfg, code, source, newDest, cursor)
            if err != nil {
                fmt.Println(err.Error())
            }
            source = newDest

            for cursor.GotoNextSibling() {
                newDest := NewNode()
                err = visitNode(cfg, code, source, newDest, cursor)
                if err != nil {
                    fmt.Println(err.Error())
                }
                source = newDest
            }
            cursor.GotoParent()  // just for good measure
        }

        // implicit return
        cfg.Graph[source] = append(cfg.Graph[source], Edge{nil, dest})

    case "package_clause":
        cfg.Graph[source] = append(
            cfg.Graph[source], Edge{cursor.Node(), dest})

    case "import_declaration":
        err := visitImportDeclaration(cfg, code, source, dest, cursor)
        if err != nil {
            return fmt.Errorf("failed to process import_declaration: %w", err)
        }

    default:
        return fmt.Errorf("unsupported AST node: %s", cursor.Node().Kind())

    }
    return nil
}

func NewCFG(ast *tree_sitter.Node, code []byte) (*CFG, error) {
    entryNode := NewNode()
    exitNode := NewNode()
    graph := Graph{
        entryNode: []Edge{},
    }

    cfg := CFG{
        Decls: []Decl{},
        FuncDefs: []FuncDef{},
        Entry: entryNode,
        Exit: exitNode,
        Graph: graph,
    }

    err := visitNode(&cfg, code, entryNode, exitNode, ast.Walk())

    return &cfg, err
}
