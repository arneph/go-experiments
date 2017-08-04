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

    elements := make([]int, 50)

    for i := range elements {
        elements[i] = i + 1
    }

    for i := range elements {
        j := rand.Intn(len(elements))

        elements[i], elements[j] = elements[j], elements[i]
    }

    f, err := os.Create("./OddEvenSort" + strconv.Itoa(len(elements)) + ".gif")

    if err != nil {
        fmt.Printf("%v\n", err)

        return
    }

    defer f.Close()

    w := bufio.NewWriter(f)

    err = makeGIF(w, elements, oddEvenSort)

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
             sortFunc func([]int, func(int, int))) error {
    anim := gif.GIF{}

    sortFunc(elements, func(i, j int) {
        addGIFFrame(&anim, elements, i, j)

        elements[i], elements[j] = elements[j], elements[i]
    })

    addGIFFrame(&anim, elements, -1, -1)

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

func bubbleSort(elements []int, swapFunc func(i, j int)) {
    for i := 1; i < len(elements); i++ {
        j := i -1

        if (elements[j] <= elements[i]) {
            continue
        }

        swapFunc(i, j)

        i = 0
    }
}

func oddEvenSort(elements []int, swapFunc func(i, j int)) {
    sorted := false

    for !sorted {
        sorted = true

        for i := 1; i < len(elements) - 1; i += 2 {
            j := i + 1

            if (elements[i] > elements[j]) {
                swapFunc(i, j);
                sorted = false
            }
        }
        
        for i := 0; i < len(elements) - 1; i += 2 {
            j := i + 1

            if (elements[i] > elements[j]) {
                swapFunc(i, j);
                sorted = false
            }
        }
    }
}

func insertionSort(elements []int, swapFunc func(i, j int)) {
    for i := 1; i < len(elements); i++ {
        for j, k := i - 1, i; j >= 0 && elements[j] > elements[k]; j, k = j - 1, k - 1 {
            swapFunc(j, k)
        }
    }
}

func selectionSort(elements []int, swapFunc func(i, j int)) {
    for i := 0; i < len(elements); i++ {
        minIndex := i

        for j := i + 1; j < len(elements); j++ {
            if elements[j] < elements[minIndex] {
                minIndex = j
            }
        }

        swapFunc(i, minIndex)
    }
}

func quickSort(elements []int, swapFunc func(i, j int)) {
    n := len(elements)

    if n < 2 {
        return
    }

    //Determine random pivot:
    pivotIndex := rand.Intn(n)
    pivot := elements[pivotIndex]

    //Partition slice:
    swapFunc(pivotIndex, n - 1) //Pivot is last element while partitioning

    j := -1 //Index of last element of first partition

    for i := 0; i < n - 1; i++ {
        if (elements[i] <= pivot) {
            j++

            swapFunc(i, j)
        }
    }

    swapFunc(j + 1, n - 1) //Pivot goes between the two partitions

    //Sort paritions:
    if (j > -1) {
        quickSort(elements[0:j + 1], swapFunc)
    }
    if (j + 2 < n) {
        quickSort(elements[j + 2:n], func(a, b int) {
            swapFunc(a + j + 2, b + j + 2) //Indices need to be converted
        })
    }
}
