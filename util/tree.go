package util

import (
    "fmt"

    tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

func PrintTree(
    code []byte,
    cursor *tree_sitter.TreeCursor,
    levels int,
    indent string,
    locs bool) {

    if levels == 0 {
        return
    }

    n := cursor.Node()
    fieldName := cursor.FieldName()
    if len(fieldName) > 0 {
        fieldName = fmt.Sprintf("%s:", fieldName)
    }
    if locs {
        r := n.Range()
        fmt.Printf("%s%s%s @ %d:%d-%d:%d\n",
            indent, fieldName, n.Kind(),
            r.StartPoint.Row + 1, r.StartPoint.Column + 1,
            r.EndPoint.Row + 1, r.EndPoint.Column + 1)
    } else {
        fmt.Printf("%s%s%s\n",
            indent, fieldName, n.Kind())
    }

    if !cursor.GotoFirstChild() {
        return
    }
    defer cursor.GotoParent()
    PrintTree(code, cursor, levels - 1, indent + "  ", locs)

    for cursor.GotoNextSibling() {
        PrintTree(code, cursor, levels - 1, indent + "  ", locs)
    }
}

func PrintFullTree(code []byte, cursor *tree_sitter.TreeCursor, locs bool) {
    PrintTree(code, cursor, -1, "", locs)
}

func PrintNodeWithChildren(code []byte, cursor *tree_sitter.TreeCursor) {
    PrintTree(code, cursor, 2, "", true)
}

func ExpectNodeKind(cursor *tree_sitter.TreeCursor, kind string) (string, error) {
    if cursor.Node().Kind() != kind {
        return "", fmt.Errorf("expected %s, got %s", kind, cursor.Node().Kind())
    }
    return kind, nil
}
