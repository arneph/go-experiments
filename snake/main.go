package main

import (
    "bufio"
    "fmt"
    "image"
    "image/color"
    "image/gif"
    "io"
    "math/rand"
    "time"
    "os"
)

func init() {
    rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
    b := makeBoard()

    f, err := os.Create("./snake/GoodPlayer.gif")

    if err != nil {
        fmt.Printf("%v\n", err)

        return
    }

    defer f.Close()

    w := bufio.NewWriter(f)

    err = makeGIF(w, b, goodPlayer)

    if err != nil {
        fmt.Printf("%v\n", err)

        return
    }

    w.Flush()
}

var palette = []color.Color{color.RGBA{0xdf, 0xdf, 0xdf, 0xff},
                            color.RGBA{0xef, 0xef, 0xef, 0xff},
                            color.RGBA{0xff, 0xff, 0x00, 0xff},
                            color.RGBA{0xff, 0x00, 0x00, 0xff}}

const (
    grayIndex = 0
    lightGrayIndex = 1
    yellowIndex = 2
    redIndex = 3
)

func makeGIF(out io.Writer,
             b *Board,
             player func(*Board)) error {
    anim := gif.GIF{}

    for b.Lost == false {
        addGIFFrame(&anim, b)

        player(b)

        b.Tick()
    }

    addGIFFrame(&anim, b)

    anim.LoopCount = len(anim.Image)

    return gif.EncodeAll(out, &anim)
}

func addGIFFrame(anim *gif.GIF, b *Board) {
    const fieldW = 16

    rect := image.Rect(0, 0, BoardW * fieldW + 1, BoardH * fieldW + 1)
    frame := image.NewPaletted(rect, palette)

    for x := 0; x < BoardW * fieldW + 1; x++ {
        for y := 0; y < BoardH * fieldW + 1; y++ {
            frame.SetColorIndex(x, y, grayIndex)
        }
    }

    for x := 0; x < BoardW; x++ {
        for y := 0; y < BoardH; y++ {
            var colorIndex uint8

            if b.Item.X == x && b.Item.Y == y {
                colorIndex = redIndex
            }else if b.Fields[x][y] == 0 {
                colorIndex = lightGrayIndex
            }else{
                colorIndex = yellowIndex
            }

            for i := 1; i < fieldW; i++ {
                for j := 1; j < fieldW; j++ {
                    frame.SetColorIndex(x * fieldW + i, y * fieldW + j, colorIndex)
                }
            }
        }
    }

    anim.Delay = append(anim.Delay, 5)
    anim.Image = append(anim.Image, frame)
}

const (
    BoardW = 24
    BoardH = 24
)

type Board struct {
    Fields    [BoardW][BoardH]BoardField
    Head      BoardLocation
    Direction BoardDirection
    Length    int
    Item      BoardLocation
    Lost      bool
}

type BoardField int

type BoardLocation struct {
    X, Y int
}

type BoardDirection int

const (
    Left  BoardDirection = iota
    Right BoardDirection = iota
    Up    BoardDirection = iota
    Down  BoardDirection = iota
)

func makeBoard() *Board {
    b := Board{}

    b.Fields[1][1] = 1
    b.Head = BoardLocation{1, 1}
    b.Direction = Right
    b.Length = 1
    b.Item = b.Head

    for b.Fields[b.Item.X][b.Item.Y] != 0 {
        b.Item.X, b.Item.Y = rand.Intn(BoardW), rand.Intn(BoardH)
    }

    b.Lost = false

    return &b
}

func (b *Board) Tick() {
    if b.Lost {
        panic("Board.Tick: Game already lost")
    }

    //Update head location:
    switch b.Direction {
    case Left:
        b.Head.X--
    case Right:
        b.Head.X++
    case Up:
        b.Head.Y--
    case Down:
        b.Head.Y++
    default:
        panic("Board.Tick: Invalid direction")
    }

    //Check if snake leaves board
    if b.Head.X < 0 || b.Head.X >= BoardW ||
       b.Head.Y < 0 || b.Head.Y >= BoardH {
        b.Lost = true

        return
    }

    //Check if snake is biting its own tail
    if b.Fields[b.Head.X][b.Head.Y] != 0 {
        b.Lost = true

        return
    }

    //Check if snake has found item
    foundItem := b.Head == b.Item

    //Update new field with head:
    b.Fields[b.Head.X][b.Head.Y] = BoardField(b.Length + 1)

    for x := 0; x < BoardW; x++ {
        for y := 0; y < BoardH; y++ {
            if b.Fields[x][y] == 0 {
                continue
            }

            if foundItem {
                b.Fields[x][y] += 1
            }else{
                b.Fields[x][y] -= 1
            }
        }
    }

    if foundItem {
        b.Length += 2

        //Place new item:
        for b.Fields[b.Item.X][b.Item.Y] != 0 {
            b.Item.X, b.Item.Y = rand.Intn(BoardW), rand.Intn(BoardH)
        }
    }
}

func simplePlayer(b *Board) {
    var h, v int

    if b.Head.X < b.Item.X {
        h = 1
    }else if b.Head.X == b.Item.X {
        h = 0
    }else{
        h = -1
    }

    if b.Head.Y < b.Item.Y {
        v = 1
    }else if b.Head.Y == b.Item.Y {
        v = 0
    }else{
        v = -1
    }

    if h == 1 {
        if b.Direction != Left {
            b.Direction = Right
        }else{
            b.Direction = Up
        }
    }else if h == -1 {
        if b.Direction != Right {
            b.Direction = Left
        }else{
            b.Direction = Up
        }
    }else if v == 1 {
        if b.Direction != Up {
            b.Direction = Down
        }else{
            b.Direction = Left
        }
    }else if v == -1 {
        if b.Direction != Down {
            b.Direction = Up
        }else{
            b.Direction = Left
        }
    }
}

func mediocrePlayer(b *Board) {
    //Determine possible directions:
    p := make(map[BoardDirection]bool, 4)

    p[Left] = true
    p[Right] = true
    p[Up] = true
    p[Down] = true

    //Avoid leaving board:
    if b.Head.X == 0 {
        p[Left] = false
    }else if b.Head.X == BoardW - 1 {
        p[Right] = false
    }

    if b.Head.Y == 0 {
        p[Up] = false
    }else if b.Head.Y == BoardH - 1 {
        p[Down] = false
    }

    //Avoid biting own tail:
    if p[Left] && b.Fields[b.Head.X - 1][b.Head.Y] != 0 {
        p[Left] = false
    }
    if p[Right] && b.Fields[b.Head.X + 1][b.Head.Y] != 0 {
        p[Right] = false
    }
    if p[Up] && b.Fields[b.Head.X][b.Head.Y - 1] != 0 {
        p[Up] = false
    }
    if p[Down] && b.Fields[b.Head.X][b.Head.Y + 1] != 0 {
        p[Down] = false
    }

    //Determine preferred direction:
    var h, v int

    if b.Head.X < b.Item.X {
        h = 1
    }else if b.Head.X == b.Item.X {
        h = 0
    }else{
        h = -1
    }

    if b.Head.Y < b.Item.Y {
        v = 1
    }else if b.Head.Y == b.Item.Y {
        v = 0
    }else{
        v = -1
    }

    if h == 1 && p[Right] {
        b.Direction = Right

    }else if h == -1 && p[Left] {
        b.Direction = Left

    }else if v == 1 && p[Down] {
        b.Direction = Down

    }else if v == -1 && p[Up] {
        b.Direction = Up

    }else if p[b.Direction] == false {
        for d, q := range p {
            if q {
                b.Direction = d
                break
            }
        }
    }
}

func goodPlayer(b *Board) {
    //Determine possible directions:
    p := make(map[BoardDirection]bool, 4)

    p[Left] = true
    p[Right] = true
    p[Up] = true
    p[Down] = true

    //Avoid leaving board:
    if b.Head.X == 0 {
        p[Left] = false
    }else if b.Head.X == BoardW - 1 {
        p[Right] = false
    }

    if b.Head.Y == 0 {
        p[Up] = false
    }else if b.Head.Y == BoardH - 1 {
        p[Down] = false
    }

    //Avoid biting own tail:
    if p[Left] && b.Fields[b.Head.X - 1][b.Head.Y] != 0 {
        p[Left] = false
    }
    if p[Right] && b.Fields[b.Head.X + 1][b.Head.Y] != 0 {
        p[Right] = false
    }
    if p[Up] && b.Fields[b.Head.X][b.Head.Y - 1] != 0 {
        p[Up] = false
    }
    if p[Down] && b.Fields[b.Head.X][b.Head.Y + 1] != 0 {
        p[Down] = false
    }

    //Avoid getting trapped:
    sizeFunc := func(l BoardLocation) (n int) {
        var queue []BoardLocation
        seen := make(map[BoardLocation]bool)

        n = 0
        queue = append(queue, l)

        for len(queue) > 0 {
            l = queue[0]

            queue = queue[1:]

            if seen[l] {
                continue
            }

            seen[l] = true

            if l.X < 0 || l.X >= BoardW ||
               l.Y < 0 || l.Y >= BoardH {
                continue
            }

            if b.Fields[l.X][l.Y] != 0 {
                continue
            }

            n++
            queue = append(queue, BoardLocation{l.X - 1, l.Y},
                                  BoardLocation{l.X + 1, l.Y},
                                  BoardLocation{l.X, l.Y - 1},
                                  BoardLocation{l.X, l.Y + 1})
        }

        return
    }

    if p[Left] && sizeFunc(BoardLocation{b.Head.X - 1, b.Head.Y}) < b.Length + 2 {
        p[Left] = false
    }
    if p[Right] && sizeFunc(BoardLocation{b.Head.X + 1, b.Head.Y}) < b.Length + 2 {
        p[Right] = false
    }
    if p[Up] && sizeFunc(BoardLocation{b.Head.X, b.Head.Y - 1}) < b.Length + 2 {
        p[Up] = false
    }
    if p[Down] && sizeFunc(BoardLocation{b.Head.X, b.Head.Y + 1}) < b.Length + 2 {
        p[Down] = false
    }

    //Determine preferred direction:
    var h, v int

    if b.Head.X < b.Item.X {
        h = 1
    }else if b.Head.X == b.Item.X {
        h = 0
    }else{
        h = -1
    }

    if b.Head.Y < b.Item.Y {
        v = 1
    }else if b.Head.Y == b.Item.Y {
        v = 0
    }else{
        v = -1
    }

    //Determine new direction:
    if h == 1 && p[Right] {
        b.Direction = Right

    }else if h == -1 && p[Left] {
        b.Direction = Left

    }else if v == 1 && p[Down] {
        b.Direction = Down

    }else if v == -1 && p[Up] {
        b.Direction = Up

    }else{
        for d, q := range p {
            if q {
                b.Direction = d
                break
            }
        }
    }
}
