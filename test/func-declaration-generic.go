package test

func generic[E any, F ~string](a E, b F) (F, E) {
    return b, a
}
