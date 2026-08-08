package main

import (
	"log"
	"sync"
	"time"
)

func makeThumbnails(filenames []string) {
	ch := make(chan struct{})
	for _, f := range filenames {
		// if _, err := ImageFile(f); err != nil {
		// 	log.Println(err)
		// }
		//
		go func(f string) {
			ImageFile(f)
			ch <- struct{}{}
		}(f)
	}

	for range filenames {
		<-ch
		println("done")
	}
}

func makeThumbnails2(filenames []string) {
	type item struct {
		thumbfile string
		err       error
	}

	ch := make(chan item, len(filenames))

	for _, f := range filenames {
		go func(f string) {
			var it item
			it.thumbfile, it.err = ImageFile(f)
			ch <- it
		}(f)
	}

	for range filenames {
		it := <-ch
		if it.err != nil {
			return
		}

		println(it.thumbfile)
	}
}

func ImageFile(f string) (string, error) {
	time.Sleep(1 * time.Second)
	i := 0
	for i < 1000000 {
		i++
	}
	time.Sleep(1 * time.Second)

	return f, nil
}

func makeThumbnails3(filenames <-chan string) int64 {
	sizes := make(chan int64)
	var wg sync.WaitGroup

	for f := range filenames {
		wg.Add(1)

		go func(f string) {
			defer wg.Done()
			thumb, err := ImageFile2(f)
			if err != nil {
				log.Println(err)
				return
			}

			sizes <- thumb
		}(f)
	}

	go func() {
		wg.Wait()
		close(sizes)
	}()

	var total int64
	for size := range sizes {
		total += size
	}

	return total
}

func ImageFile2(f string) (int64, error) {
	time.Sleep(1 * time.Second)
	i := 0
	for i < 1000000 {
		i++
	}
	time.Sleep(1 * time.Second)
	return int64(len(f)), nil
}

func main() {
	// filenames := []string{"abc", "abc", "abc", "abc", "abc", "abc"}

	// makeThumbnails(filenames)
	// makeThumbnails2(filenames)
	//
	filenames := make(chan string)

	go func() {
		defer close(filenames)
		for i := 0; i < 10; i++ {
			filenames <- "abc"
		}
	}()

	println(makeThumbnails3(filenames))
}
