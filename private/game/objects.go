package game

func Rect(x, y, w, h int, tags []int, p *[][]Pixel) {
    (*p)[y][x] = NewPixel(
        '+',
        tags,
    )
    (*p)[y+h][x] = (*p)[y][x]
    (*p)[y][x+w] = (*p)[y][x]
    (*p)[y+h][x+w] = (*p)[y][x+w]
    for i := 1; i < h; i++ {
        (*p)[y+i][x] = NewPixel(
            '|',
            tags,
        )
        (*p)[y+i][x+w] = (*p)[y+i][x]
    }
    for i := 1; i < w; i++ {
        (*p)[y][x+i] = NewPixel(
            '-',
            tags,
        )
        (*p)[y+h][x+i] = (*p)[y][x+i]
    }
}
