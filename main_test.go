package main

import (
	"math/rand"
	"runtime"
	"sync"
	"testing"
)

var globalResultChan chan int
var globalResult int
var workerCount = 4

func NRand(n int) []int {
	i := make([]int, n)
	for ind := range i {
		i[ind] = rand.Int()
	}
	return i
}

func BenchmarkStoreRegular(b *testing.B) {
	nums := NRand(b.N)
	rm := NewRegularIntMap()
	b.ResetTimer()
	for _, v := range nums {
		rm.Store(v, v)
	}
}

func BenchmarkStoreSync(b *testing.B) {
	nums := NRand(b.N)
	var sm sync.Map
	b.ResetTimer()
	for _, v := range nums {
		sm.Store(v, v)
	}
}

func BenchmarkDeleteRegular(b *testing.B) {
	nums := NRand(b.N)
	rm := NewRegularIntMap()
	for _, v := range nums {
		rm.Store(v, v)
	}
	
	b.ResetTimer()
	for _, v := range nums {
		rm.Delete(v)
	}
}

func BenchmarkDeleteSync(b *testing.B) {
	nums := NRand(b.N)
	var sm sync.Map
	for _, v := range nums {
		sm.Store(v, v)
	}
	
	b.ResetTimer()
	for _, v := range nums {
		sm.Delete(v)
	}
}

func BenchmarkLoadRegularFound(b *testing.B) {
	nums := NRand(b.N)
	rm := NewRegularIntMap()
	for _, v := range nums {
		rm.Store(v, v)
	}
	
	currentResult := 0
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		currentResult, _ = rm.Load(nums[i])
	}
	globalResult = currentResult
}

func BenchmarkLoadRegularNotFound(b *testing.B) {
	nums := NRand(b.N)
	rm := NewRegularIntMap()
	for _, v := range nums {
		rm.Store(v, v)
	}
	currentResult := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		currentResult, _ = rm.Load(i)
	}
	globalResult = currentResult
}

func BenchmarkLoadSyncFound(b *testing.B) {
	nums := NRand(b.N)
	var sm sync.Map
	for _, v := range nums {
		sm.Store(v, v)
	}
	currentResult := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r, ok := sm.Load(nums[i])
		if ok {
			currentResult = r.(int)
		}
	}
	globalResult = currentResult
}

func BenchmarkLoadSyncNotFound(b *testing.B) {
	nums := NRand(b.N)
	var sm sync.Map
	for _, v := range nums {
		sm.Store(v, v)
	}
	currentResult := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r, ok := sm.Load(i)
		if ok {
			currentResult = r.(int)
		}
	}
	globalResult = currentResult
}

func BenchmarkRegularStableKeys(b *testing.B) {
	runtime.GOMAXPROCS(workerCount)
	
	rm := NewRegularIntMap()
	populateMap(b.N, rm)
	
	var wg sync.WaitGroup
	wg.Add(workerCount)
	
	// Holds our final results, to prevent compiler optimizations.
	globalResultChan = make(chan int, workerCount)
	
	b.ResetTimer()
	
	for wc := 0; wc < workerCount; wc++ {
		go func(n int) {
			currentResult := 0
			for i := 0; i < n; i++ {
				currentResult, _ = rm.Load(5)
			}
			globalResultChan <- currentResult
			wg.Done()
		}(b.N)
	}
	
	wg.Wait()
}

func BenchmarkSyncStableKeys(b *testing.B) {
	runtime.GOMAXPROCS(workerCount)
	
	var sm sync.Map
	populateSyncMap(b.N, &sm)
	
	var wg sync.WaitGroup
	wg.Add(workerCount)
	
	// Holds our final results, to prevent compiler optimizations.
	globalResultChan = make(chan int, workerCount)
	
	b.ResetTimer()
	
	for wc := 0; wc < workerCount; wc++ {
		go func(n int) {
			currentResult := 0
			for i := 0; i < n; i++ {
				r, ok := sm.Load(5)
				if ok {
					currentResult = r.(int)
				}
			}
			globalResultChan <- currentResult
			wg.Done()
		}(b.N)
	}
	
	wg.Wait()
}

func BenchmarkRegularStableKeysFound(b *testing.B) {
	runtime.GOMAXPROCS(workerCount)
	
	rm := NewRegularIntMap()
	values := populateMap(b.N, rm)
	
	var wg sync.WaitGroup
	wg.Add(workerCount)
	
	// Holds our final results, to prevent compiler optimizations.
	globalResultChan = make(chan int, workerCount)
	
	b.ResetTimer()
	
	for wc := 0; wc < workerCount; wc++ {
		go func(n int) {
			currentResult := 0
			for i := 0; i < n; i++ {
				currentResult, _ = rm.Load(values[i])
			}
			globalResultChan <- currentResult
			wg.Done()
		}(b.N)
	}
	
	wg.Wait()
}

func BenchmarkSyncStableKeysFound(b *testing.B) {
	runtime.GOMAXPROCS(workerCount)
	
	var sm sync.Map
	values := populateSyncMap(b.N, &sm)
	
	var wg sync.WaitGroup
	wg.Add(workerCount)
	
	// Holds our final results, to prevent compiler optimizations.
	globalResultChan = make(chan int, workerCount)
	
	b.ResetTimer()
	
	for wc := 0; wc < workerCount; wc++ {
		go func(n int) {
			currentResult := 0
			for i := 0; i < n; i++ {
				r, ok := sm.Load(values[i])
				if ok {
					currentResult = r.(int)
				}
			}
			globalResultChan <- currentResult
			wg.Done()
		}(b.N)
	}
	
	wg.Wait()
}

func populateMap(n int, intMap *RegularIntMap) []int {
	i := make([]int, n)
	for ind := range i {
		i[ind] = rand.Int()
	}
	return i
}

// interface{}
func populateSyncMap(n int, i2 *sync.Map) []int {
	i := make([]int, n)
	for ind := range i {
		i[ind] = rand.Int()
	}
	return i
}
