package main

import (
    "fmt"
    "math/rand"
    "time"
)

func main() {
    rand.Seed(time.Now().UTC().UnixNano())

    elements := make([]int, 10000000)

    for i := range elements {
        elements[i] = i + 1
    }

    for i := range elements {
        j := rand.Intn(len(elements))

        elements[i], elements[j] = elements[j], elements[i]
    }

    elements1 := make([]int, len(elements))
    elements2 := make([]int, len(elements))
    copy(elements1, elements)
    copy(elements2, elements)

    fmt.Printf("Number of items:      %d\n", len(elements))

    start := time.Now()

    quickSort(elements1, 0)

    fmt.Printf("QuickSort sequential: %05.03fs\n", time.Since(start).Seconds())

    start = time.Now()

    quickSort(elements2, 2)

    fmt.Printf("QuickSort parallel:   %05.03fs\n", time.Since(start).Seconds())
}

func quickSort(elements []int, parallel int) {
    n := len(elements)

    if n < 2 {
        return
    }

    //Determine random pivot:
    pivotIndex := rand.Intn(n)
    pivot := elements[pivotIndex]

    //Partition slice:
    elements[pivotIndex], elements[n - 1] = elements[n - 1], elements[pivotIndex] //Pivot is last element while partitioning

    j := -1 //Index of last element of first partition

    for i := 0; i < n - 1; i++ {
        if (elements[i] <= pivot) {
            j++

            elements[i], elements[j] = elements[j], elements[i]
        }
    }

    elements[j + 1], elements[n - 1] = elements[n - 1], elements[j + 1] //Pivot goes between the two partitions

    //Sort paritions:
    if (parallel == 0) {
        if (j > -1) {
            quickSort(elements[0:j + 1], 0)
        }
        if (j + 2 < n) {
            quickSort(elements[j + 2:n], 0)
        }

    }else{
        type Result struct{
            index int
            partition []int
        }

        results := make(chan Result, 2)

        if (j > -1) {
            partition := make([]int, j + 1)
            copy(partition, elements[0:j + 1])

            go func() {
                quickSort(partition, parallel - 1)
                results <- Result{1, partition}
            } ()

        }else{
            results <- Result{0, nil}
        }
        if (j + 2 < n) {
            partition := make([]int, n - j - 2)
            copy(partition, elements[j + 2:n])

            go func() {
                quickSort(partition, parallel - 1)
                results <- Result{2, partition}
            } ()

        }else{
            results <- Result{0, nil}
        }

        for i := 0; i < 2; i++ {
            res := <- results

            if (res.index == 1) {
                for i, e := range res.partition {
                    elements[i] = e
                }
            }else if (res.index == 2) {
                for i, e := range res.partition {
                    elements[j + 2 + i] = e
                }
            }
        }

        close(results)
    }
}
