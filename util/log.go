package util

import (
    "fmt"
    "os"
)


func CombineErrors(first error, second error) error {
    if first == nil && second == nil { return nil }
    if second == nil { return first }
    if first == nil { return second }
    return fmt.Errorf("%w; %w", first, second)
}

func PrintErr(msg string, args... any) {
    err := fmt.Errorf(msg, args...)
    fmt.Fprintln(os.Stderr, err.Error())
}

func PrintErrAndDie(msg string, args... any) {
    PrintErr(msg, args...)
    os.Exit(1)
}
