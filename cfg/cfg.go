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
    Type string
}

type FuncDef struct {
    Id *Decl
    Tree *tree_sitter.Node
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
        builder.WriteString(fmt.Sprintf(" %s:%s", decl.Name, decl.Type))
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

func printNodeWithChildren(code []byte, cursor *tree_sitter.TreeCursor) {
    fmt.Printf("Children of %s: %s\n", cursor.Node().Kind(),
        cursor.Node().Utf8Text(code))
    i := 0
    if cursor.GotoFirstChild() {
        defer cursor.GotoParent()

        name := cursor.FieldName()
        if len(name) > 0 {
            name = fmt.Sprintf("/%s", name)
        }
        fmt.Printf("%d%s %s: %s\n", i, name, cursor.Node().Kind(),
            cursor.Node().Utf8Text(code))

        for cursor.GotoNextSibling() {
            i += 1
            name := cursor.FieldName()
            if len(name) > 0 {
                name = fmt.Sprintf("/%s", name)
            }
            fmt.Printf("%d%s %s: %s\n", i, name, cursor.Node().Kind(),
                cursor.Node().Utf8Text(code))
        }
    }
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
                    cfg.Decls = append(cfg.Decls, Decl{importName, "module"})
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
                cfg.Decls = append(cfg.Decls, Decl{importName, "any"})
            } else {
                cfg.Decls = append(cfg.Decls, Decl{importName, "module"})
            }
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

func visitImportSpecList(
    cfg *CFG,
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

func visitImportDeclaration(
    cfg *CFG,
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

func extractTypeFromExpression(
    cursor *tree_sitter.TreeCursor) (string, error) {

    switch cursor.Node().Kind() {
    case "int_literal":
        return "int", nil
    case "interpreted_string_literal":
        return "string", nil
    default:
        return "<unknown>", fmt.Errorf("unsupported expression type: %s", cursor.Node().Kind())
    }
}

func extractTypesFromExpressionList(
    code []byte,
    cursor *tree_sitter.TreeCursor) ([]string, error) {

    typeList := []string{}

    if !cursor.GotoFirstChild() {
        return []string{}, fmt.Errorf("expression_list has no children")
    }
    defer cursor.GotoParent()

    var cumulErr error = nil

    typeName, err := extractTypeFromExpression(cursor)
    if err != nil {
        typeName = "<unknown>"
        cumulErr = log.CombineErrors(cumulErr, err)
    }
    typeList = append(typeList, typeName)

    for cursor.GotoNextSibling() {
        if cursor.Node().Kind() == "," {
            continue
        }
        typeName, err := extractTypeFromExpression(cursor)
        if err != nil {
            typeName = "<unknown>"
            cumulErr = log.CombineErrors(cumulErr, err)
        }
        typeList = append(typeList, typeName)
    }

    return typeList, cumulErr
}

func visitVarSpec(
    cfg *CFG,
    code []byte,
    cursor *tree_sitter.TreeCursor) error {

    var cumulErr error = nil
    names := []string{}
    for _, nameNode := range cursor.Node().ChildrenByFieldName(
        "name", cursor.Node().Walk()) {
        
        if nameNode.Kind() != "identifier" {
            cumulErr = log.CombineErrors(cumulErr, fmt.Errorf(
                "expected var_spec name to be an identifier, got %s",
                nameNode.Kind()))
            continue
        }
        varName := nameNode.Utf8Text(code)
        names = append(names, varName)
    }

    explicitTypeName := ""
    typeNode := cursor.Node().ChildByFieldName("type")
    if typeNode != nil {
        explicitTypeName = typeNode.Utf8Text(code)
    }

    valueNode := cursor.Node().ChildByFieldName("value")
    var valueTypeNames []string = nil
    if valueNode != nil {
        var err error = nil
        valueTypeNames, err = extractTypesFromExpressionList(code, valueNode.Walk())
        if err != nil {
            cumulErr = log.CombineErrors(cumulErr, fmt.Errorf(
                "failed to determine types from var_spec value expressions"))
        } else {
            // sanity checks
            if len(valueTypeNames) != len(names) {
                cumulErr = log.CombineErrors(cumulErr, fmt.Errorf(
                    "%d types extracted from const_spec value " +
                    "expressions does not match %d declared names",
                    len(valueTypeNames), len(names)))
            } else if explicitTypeName != "" {
                for _, valueType := range valueTypeNames {
                    if valueType != explicitTypeName {
                        cumulErr = log.CombineErrors(cumulErr, fmt.Errorf(
                            "explicit type '%s' and inferred value type '%s' do not match", 
                            explicitTypeName, valueType))
                    }
                }
            }
        }
    }

    for i, declName := range names {
        typeName := ""
        if explicitTypeName != "" {
            typeName = explicitTypeName
        } else if valueTypeNames != nil && i < len(valueTypeNames) {
            typeName = valueTypeNames[i]
        } else {
            typeName = "<unknown>"
        }
        cfg.Decls = append(cfg.Decls, Decl{declName, typeName})
    }

    return cumulErr
}

func visitVarDeclaration(
    cfg *CFG,
    code []byte,
    source *Node,
    dest *Node,
    cursor *tree_sitter.TreeCursor) error {

    cfg.Graph[source] = append(cfg.Graph[source], Edge{cursor.Node(), dest})

    if !cursor.GotoFirstChild() {
        return fmt.Errorf("no children for var declaration")
    }
    defer cursor.GotoParent()
    
    if cursor.Node().Kind() != "var" {
        return fmt.Errorf("unexpected first child for var declaration: %s", cursor.Node().Kind())
    }

    if !cursor.GotoNextSibling() || cursor.Node().Kind() != "var_spec" {
        return fmt.Errorf("no var_spec for var_declaration")
    }

    return visitVarSpec(cfg, code, cursor)
}

func visitConstSpec(
    cfg *CFG,
    code []byte,
    cursor *tree_sitter.TreeCursor) error {

    var cumulErr error = nil
    names := []string{}
    for _, nameNode := range cursor.Node().ChildrenByFieldName(
        "name", cursor.Node().Walk()) {

        if nameNode.Kind() == "," {
            continue
        }
        
        if nameNode.Kind() != "identifier" {
            cumulErr = log.CombineErrors(cumulErr, fmt.Errorf(
                "expected const_spec name to be an identifier, got %s",
                nameNode.Kind()))
            continue
        }
        constName := nameNode.Utf8Text(code)
        names = append(names, constName)
    }

    explicitTypeName := ""
    typeNode := cursor.Node().ChildByFieldName("type")
    if typeNode != nil {
        explicitTypeName = typeNode.Utf8Text(code)
    }

    valueNode := cursor.Node().ChildByFieldName("value")
    var valueTypeNames []string = nil
    if valueNode != nil {
        var err error = nil
        valueTypeNames, err = extractTypesFromExpressionList(code, valueNode.Walk())
        if err != nil {
            cumulErr = log.CombineErrors(cumulErr, fmt.Errorf(
                "failed to determine types from const_spec value expressions: %w", err))
        } else {
            // sanity checks
            if len(valueTypeNames) != len(names) {
                cumulErr = log.CombineErrors(cumulErr, fmt.Errorf(
                    "%d types extracted from const_spec value " +
                    "expressions does not match %d declared names",
                    len(valueTypeNames), len(names)))
            } else if explicitTypeName != "" {
                for _, valueType := range valueTypeNames {
                    if valueType != explicitTypeName {
                        cumulErr = log.CombineErrors(cumulErr, fmt.Errorf(
                            "explicit type '%s' and inferred value type '%s' do not match", 
                            explicitTypeName, valueType))
                    }
                }
            }
        }
    } else {
        cumulErr = log.CombineErrors(cumulErr, fmt.Errorf(
            "const_spec has no value expression"))
    }

    for i, declName := range names {
        typeName := ""
        if explicitTypeName != "" {
            typeName = explicitTypeName
        } else if valueTypeNames != nil && i < len(valueTypeNames) {
            typeName = valueTypeNames[i]
        } else {
            typeName = "<unknown>"
        }
        cfg.Decls = append(cfg.Decls, Decl{declName, typeName})
    }

    return cumulErr
}

func visitConstDeclaration(
    cfg *CFG,
    code []byte,
    source *Node,
    dest *Node,
    cursor *tree_sitter.TreeCursor) error {

    cfg.Graph[source] = append(cfg.Graph[source], Edge{cursor.Node(), dest})

    if !cursor.GotoFirstChild() {
        return fmt.Errorf("no children for const declaration")
    }
    defer cursor.GotoParent()
    
    if cursor.Node().Kind() != "const" {
        return fmt.Errorf("unexpected first child for const declaration: %s",
            cursor.Node().Kind())
    }

    if !cursor.GotoNextSibling() || cursor.Node().Kind() != "const_spec" {
        return fmt.Errorf("no const_spec for const_declaration")
    }

    return visitConstSpec(cfg, code, cursor)
}

func extractTypeParameterDeclarationName(
    code []byte,
    cursor *tree_sitter.TreeCursor) (string, error) {

    paramText := cursor.Node().Utf8Text(code)

    // %%% shady, just to avoid having spaces in type string representations
    paramText = strings.Replace(paramText, " ", "<", 1)
    
    return paramText, nil
}

func extractTypeParameterListText(
    code []byte,
    cursor *tree_sitter.TreeCursor) (string, error) {

    if cursor.Node().Kind() != "type_parameter_list" {
        return "", fmt.Errorf("expected type_parameter_list, got %s",
            cursor.Node().Kind())
    }

    if !cursor.GotoFirstChild() {
        return "", fmt.Errorf("no children for type_parameter_list")
    }
    defer cursor.GotoParent()

    typeParamBuf := strings.Builder{}
    var cumulErr error = nil

    if cursor.Node().Kind() != "[" {
        cumulErr = log.CombineErrors(cumulErr, fmt.Errorf(
            "expected opening bracket in type_parameter_list, got %s",
            cursor.Node().Kind()))
    }

    typeParamBuf.WriteString("[")

    if cursor.Node().Kind() == "type_parameter_declaration" {
        tpdName, err := extractTypeParameterDeclarationName(code, cursor)
        if err != nil {
            cumulErr = log.CombineErrors(cumulErr, fmt.Errorf(
                "failed to extract type parameter declaration name: %w", err))
            typeParamBuf.WriteString("<error>")
        } else {
            typeParamBuf.WriteString(tpdName)
        }
    }

    for cursor.GotoNextSibling() {
        if cursor.Node().Kind() == "," {
            typeParamBuf.WriteString(",")
        } else if cursor.Node().Kind() == "]" {
            typeParamBuf.WriteString("]")
        } else if cursor.Node().Kind() == "type_parameter_declaration" {
            tpdName, err := extractTypeParameterDeclarationName(code, cursor)
            if err != nil {
                cumulErr = log.CombineErrors(cumulErr, fmt.Errorf(
                    "failed to extract type parameter declaration name: %w", err))
                typeParamBuf.WriteString("<error>")
            } else {
                typeParamBuf.WriteString(tpdName)
            }
        } else {
            cumulErr = log.CombineErrors(cumulErr, fmt.Errorf(
                "expected type_parameter_declaration, got %s",
                cursor.Node().Kind()))
            typeParamBuf.WriteString("<error>")
        }
    }

    return typeParamBuf.String(), cumulErr
}

func extractParameterDeclarationTypeName(
    code []byte,
    cursor *tree_sitter.TreeCursor) (string, error) {

    if cursor.Node().Kind() != "parameter_declaration" {
        return "", fmt.Errorf("expected parameter_declaration, got %s",
            cursor.Node().Kind())
    }

    nameNodes := cursor.Node().ChildrenByFieldName("name", cursor.Copy())
    numParams := len(nameNodes)

    typeNode := cursor.Node().ChildByFieldName("type")
    if typeNode == nil {
        return "<error>", fmt.Errorf("parameter_declaration has no type")
    }
    typeText := typeNode.Utf8Text(code)
    totalTypeText := typeText

    for i := 1; i<numParams; i++ {
        totalTypeText += "," + typeText
    }

    return totalTypeText, nil
}

func extractParameterListText(
    code []byte,
    cursor *tree_sitter.TreeCursor) (string, error) {

    if cursor.Node().Kind() != "parameter_list" {
        return "", fmt.Errorf("expected parameter_list, got %s",
            cursor.Node().Kind())
    }

    if !cursor.GotoFirstChild() {
        return "", fmt.Errorf("no children for parameter_list")
    }
    defer cursor.GotoParent()

    paramBuf := strings.Builder{}

    var cumulErr error = nil

    if cursor.Node().Kind() != "(" {
        cumulErr = log.CombineErrors(cumulErr, fmt.Errorf(
            "expected opening paren in parameters, got %s",
            cursor.Node().Kind()))
    }

    paramBuf.WriteString("(")

    for cursor.GotoNextSibling() {
        if cursor.Node().Kind() == "," {
            paramBuf.WriteString(",")
        } else if cursor.Node().Kind() == ")" {
            break
        } else if cursor.Node().Kind() == "parameter_declaration" {
            tpdName, err := extractParameterDeclarationTypeName(code, cursor)
            if err != nil {
                cumulErr = log.CombineErrors(cumulErr, fmt.Errorf(
                    "failed to extract parameter declaration name: %w", err))
                paramBuf.WriteString("<error>")
            } else {
                paramBuf.WriteString(tpdName)
            }
        } else {
            cumulErr = log.CombineErrors(cumulErr, fmt.Errorf(
                "expected parameter_declaration, got %s",
                cursor.Node().Kind()))
            paramBuf.WriteString("<error>")
        }
    }

    if cursor.GotoNextSibling() {
        cumulErr = log.CombineErrors(cumulErr, fmt.Errorf(
            "unexpected additional child in parameter_list: %s",
            cursor.Node().Kind()))
    }

    paramBuf.WriteString(")")

    return paramBuf.String(), cumulErr
}

func visitFunctionDeclaration(
    cfg *CFG,
    code []byte,
    source *Node,
    dest *Node,
    cursor *tree_sitter.TreeCursor) error {

    cfg.Graph[source] = append(cfg.Graph[source], Edge{cursor.Node(), dest})

    var cumulErr error = nil

    funcName := ""
    nameNode := cursor.Node().ChildByFieldName("name")
    if nameNode == nil {
        cumulErr = log.CombineErrors(cumulErr, fmt.Errorf(
            "no name for function declaration"))
    } else {
        funcName = nameNode.Utf8Text(code)
    }

    funcTypeBuf := strings.Builder{}
    funcTypeBuf.WriteString("func")

    typeParamNode := cursor.Node().ChildByFieldName("type_parameters")
    if typeParamNode != nil {
        typeParamText, err := extractTypeParameterListText(code, typeParamNode.Walk())
        if err != nil {
            cumulErr = log.CombineErrors(cumulErr, err)
        }
        funcTypeBuf.WriteString(typeParamText)
    }

    paramNode := cursor.Node().ChildByFieldName("parameters")
    if paramNode == nil {
        cumulErr = log.CombineErrors(cumulErr, fmt.Errorf(
            "no parameter list for function declaration"))
    } else {
        paramText, err := extractParameterListText(code, paramNode.Walk())
        if err != nil {
            cumulErr = log.CombineErrors(cumulErr, err)
        }
        funcTypeBuf.WriteString(paramText)
    }

    resultNode := cursor.Node().ChildByFieldName("result")
    if resultNode != nil {
        resultText, err := extractParameterListText(code, resultNode.Walk())
        if err != nil {
            cumulErr = log.CombineErrors(cumulErr, err)
        }
        funcTypeBuf.WriteString(resultText)
    }

    cfg.Decls = append(cfg.Decls, Decl{funcName, funcTypeBuf.String()})
    cfg.FuncDefs = append(cfg.FuncDefs, FuncDef{&cfg.Decls[len(cfg.Decls)-1], cursor.Node()})

    return cumulErr
}

func visitNode(
    cfg *CFG,
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

    case "var_declaration":
        err := visitVarDeclaration(cfg, code, source, dest, cursor)
        if err != nil {
            return fmt.Errorf("failed to process var_declaration: %w", err)
        }

    case "const_declaration":
        err := visitConstDeclaration(cfg, code, source, dest, cursor)
        if err != nil {
            return fmt.Errorf("failed to process const_declaration: %w", err)
        }

    case "function_declaration":
        err := visitFunctionDeclaration(cfg, code, source, dest, cursor)
        if err != nil {
            return fmt.Errorf("failed to process func_declaration: %w", err)
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
