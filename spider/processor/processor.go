package processor

import "sync"

type Processor interface {
	Run(url string)
}

func Run(urls []string, p Processor, poolSize int) {
	// 限制并发数
	sem := make(chan struct{}, poolSize)
	defer close(sem)
	var wg sync.WaitGroup

	for i, url := range urls {
		sem <- struct{}{} // 若 sem 满了，则发生阻塞，等待 <-sem
		wg.Add(1)
		go func(i int, url string) {
			defer func() {
				<-sem
				wg.Done()
			}()
			p.Run(url)
		}(i, url)
	}
	wg.Wait()
}

