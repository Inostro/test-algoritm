package main

// go test -v
import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

type SafeMap struct {
	mu  sync.Mutex
	m   map[int]int
	cnt int /// счётчик обращений
	add int /// счётчик добавлений
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		m: make(map[int]int),
	}
}

func (s *SafeMap) Inc(key int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cnt++

	if _, ok := s.m[key]; !ok {
		s.add++
		s.m[key] = 0
	}

	s.m[key]++
}

func TestAlgoritm(t *testing.T) {
	year := 1955

	sm := NewSafeMap()

	keys := make([]int, year)
	for i := 1; i <= year; i++ {
		keys[i-1] = i
	}

	rand.Seed(time.Now().UnixNano())
	for i := range keys {
		j := rand.Intn(i + 1)
		keys[i], keys[j] = keys[j], keys[i]
	}
	var wg sync.WaitGroup
	for g := 0; g < 4; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, k := range keys {
				sm.Inc(k)
			}
		}()
	}
	wg.Wait()

	sm.mu.Lock()
	defer sm.mu.Unlock()

	for k, v := range sm.m {
		if v != 4 {
			t.Errorf("ключ %d = %d, а надо 4", k, v)
		}
	}

	if sm.cnt != year*4 {
		t.Errorf("обращений %d, надо %d", sm.cnt, year*4)
	}
	if sm.add != year {
		t.Errorf("добавлений %d, надо %d", sm.add, year)
	}
}
