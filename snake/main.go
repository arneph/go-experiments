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

    f, err := os.Create("./snake/Player 1.gif")

    if err != nil {
        fmt.Printf("%v\n", err)

        return
    }

    defer f.Close()

    w := bufio.NewWriter(f)

    err = makeGIF(w, b, simplePlayer)

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
    BoardW = 72
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
