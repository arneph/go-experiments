package main

import (
    "fmt"
//    "image"
    "image/color"
//    "imgage/gif"
//    "io"
//    "math"
    "math/rand"
    "time"
//    "os"
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
        elements[i] = i
    }

    for i := range elements {
        j := rand.Intn(len(elements))

        elements[i], elements[j] = elements[j], elements[i]
    }

    swaps := 0

    for {
        res, i, j := bubbleSortStep(elements)

        if (!res) {
            swaps++
            fmt.Printf("Swap %03d: %02d with %02d\n", swaps, i, j);

        }else{
            break
        }
    }
    
    fmt.Printf("%d swaps\n", swaps);
}

func printElements(elements []int) {
    for i, v := range elements {
        fmt.Printf("%2d: %2d\n", i, v);
    }
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
