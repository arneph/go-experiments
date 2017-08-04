package main

import (
    "bufio"
    "fmt"
    "image"
    "image/color"
    "image/gif"
    "io"
    "math/rand"
    "strconv"
    "time"
    "os"
)

var palette = []color.Color{color.Black,
                            color.RGBA{0x00, 0xff, 0x00, 0xff},
                            color.RGBA{0xff, 0x00, 0x00, 0xff}}

const (
    blackIndex = 0
    greenIndex = 1
    redIndex = 2
)

func main() {
    rand.Seed(time.Now().UTC().UnixNano())

    elements := make([]int, 100)

    for i := range elements {
        elements[i] = i + 1
    }

    for i := range elements {
        j := rand.Intn(len(elements))

        elements[i], elements[j] = elements[j], elements[i]
    }

    f, err := os.Create("./BubbleSort" + strconv.Itoa(len(elements)) + ".gif")

    if err != nil {
        fmt.Printf("%v\n", err)

        return
    }

    defer f.Close()

    w := bufio.NewWriter(f)

    err = makeGIF(w, elements, bubbleSortStep)

    if err != nil {
        fmt.Printf("%v\n", err)

        return
    }

    w.Flush()
}

func printElements(elements []int) {
    for i, v := range elements {
        fmt.Printf("%2d: %2d\n", i, v);
    }
}

func makeGIF(out io.Writer,
             elements []int,
             sortFunc func([]int) (bool, int, int)) error {
    anim := gif.GIF{}

    addGIFFrame(&anim, elements, -1, -1)

    for {
        res, i, j := sortFunc(elements)

        addGIFFrame(&anim, elements, i, j)

        if (res) {
            break
        }
    }

    anim.LoopCount = len(anim.Image)

    return gif.EncodeAll(out, &anim)
}

func addGIFFrame(anim *gif.GIF, elements []int, a, b int) {
    const ew = 5

    n := len(elements)
    rect := image.Rect(0, 0, n * ew, n * ew)
    frame := image.NewPaletted(rect, palette)

    for i, v := range elements {
        for x, y := 0, 0; x < ew && y < ew * v; x, y = (x + 1) % ew, y + ((x + 1) / ew) {
            if (i == a || i == b) {
                frame.SetColorIndex(i * ew + x, n * ew - y, redIndex)
            }else{
                frame.SetColorIndex(i * ew + x, n * ew - y, greenIndex)
            }
        }
    }

    anim.Delay = append(anim.Delay, 1)
    anim.Image = append(anim.Image, frame)
}

func bubbleSortStep(elements []int) (bool, int, int) {
    for i := 0; i < len(elements) - 1; i++ {
        j := i + 1

        if (elements[i] <= elements[j]) {
            continue
        }

        elements[i], elements[j] = elements[j], elements[i]

        return false, i, j
    }

    return true, -1, -1
}
