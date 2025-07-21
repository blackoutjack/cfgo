package cfg

import (
    "fmt" 
    "strings"
    "strconv"

    tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

var nodeIndex = 0

type Node struct {
    id int
}
func NewNode() *Node {
    node := Node{nodeIndex}
    nodeIndex++
    return &node
}
func (n *Node) String() string {
    if n == nil { return "nil" }
    return strconv.Itoa(n.id)
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

    builder.WriteString("Declarations: ")
    for _, decl := range g.Decls {
        builder.WriteString(decl.Name)
    }
    builder.WriteString("\n")

    builder.WriteString("Function definitions: ")
    for _, fdef := range g.FuncDefs {
        builder.WriteString(fdef.String())
    }
    builder.WriteString("\n")

    builder.WriteString("Entry: ");
    builder.WriteString(g.Entry.String())
    builder.WriteString("\n")

    builder.WriteString("Exit: ")
    builder.WriteString(g.Exit.String())
    builder.WriteString("\n")

    builder.WriteString("Graph:\n")
    for n, edges := range g.Graph {
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

func visitNode(cfg *CFG, source *Node, dest *Node, cursor *tree_sitter.TreeCursor) error {
    switch cursor.Node().Kind() {
    case "source_file":
        /*
        for _, child := range node.Children(node.Walk()) {
            subdest := NewNode()
            err := visitNode(cfg, source, subdest, &child)
            if err != nil {
                fmt.Println(err.Error())
            } else {
                source = subdest
            }
        }
        */
        if cursor.GotoFirstChild() {
            newDest := NewNode()
            err := visitNode(cfg, source, newDest, cursor)
            if err != nil {
                fmt.Println(err.Error())
            } else {
                source = newDest
            }

            for cursor.GotoNextSibling() {
                newDest := NewNode()
                err = visitNode(cfg, source, newDest, cursor)
                if err != nil {
                    fmt.Println(err.Error())
                } else {
                    source = newDest
                }
            }
        }
        
        // implicit return
        cfg.Graph[source] = append(cfg.Graph[source], Edge{nil, dest})

    case "package_clause":
        cfg.Graph[source] = append(cfg.Graph[source], Edge{cursor.Node(), dest})
    case "import_declaration":
        return fmt.Errorf("unsupported AST node: %s", cursor.Node().Kind())
    default:
        return fmt.Errorf("unsupported AST node: %s", cursor.Node().Kind())
    }
    return nil
}

func NewCFG(ast *tree_sitter.Node) (*CFG, error) {
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

    err := visitNode(&cfg, entryNode, exitNode, ast.Walk())

    return &cfg, err
}
