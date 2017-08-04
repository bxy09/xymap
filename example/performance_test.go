package example_test

import (
	"fmt"
	"github.com/bxy09/xymap/example"
	"github.com/docker/docker/pkg/testutil/assert"
	"math/rand"
	"runtime"
	"runtime/debug"
	"testing"
)

var length = 1000

func doNothing(a int){}

func BenchmarkMappingIter(b *testing.B) {
	debug.SetGCPercent(-1)
	mapping := make(map[int]int)
	for i := 0; i < length; i++ {
		mapping[i] = rand.Int()
	}
	runtime.GC()
	b.Run("pair", func(b *testing.B) {
		sum := 0
		for j := 0; j < b.N; j++ {
			for i, k := range mapping {
				sum += i + k
			}
		}
		doNothing(sum)
	})
	runtime.GC()
	b.Run("key", func(b *testing.B) {
		sum := 0
		for j := 0; j < b.N; j++ {
			for i := range mapping {
				sum += i + mapping[i]
			}
		}
		doNothing(sum)
	})
	runtime.GC()
	b.Run("value", func(b *testing.B) {
		sum := 0
		for j := 0; j < b.N; j++ {
			for i := 0; i < length; i++ {
				sum += i + mapping[i]
			}
		}
		doNothing(sum)
	})
	for i := 0; i < length; i++ {
		mapping[i + length] = rand.Int()
		delete(mapping, i * 2)
	}
	b.Run("pair-half-empty", func(b *testing.B) {
		sum := 0
		for j := 0; j < b.N; j++ {
			for i, k := range mapping {
				sum += i + k
			}
		}
		doNothing(sum)
	})
}

func BenchmarkXYMapIter(b *testing.B) {
	debug.SetGCPercent(-1)
	mapping := example.NewXYMapIntInt()
	for i := 0; i < length; i++ {
		mapping.Set(i, rand.Int())
	}
	runtime.GC()
	b.Run("pair", func(b *testing.B) {
		sum := 0
		for j := 0; j < b.N; j++ {
			mapping.Iterate(func(i int, k int) (broken bool) {
				sum += i + k
				return
			})
		}
		doNothing(sum)
	})
	b.Run("value", func(b *testing.B) {
		sum := 0
		for j := 0; j < b.N; j++ {
			for i := 0; i < length; i++ {
				k, _ := mapping.Get(i)
				sum += i + k
			}
		}
		doNothing(sum)
	})
	for i := 0; i < length; i++ {
		mapping.Set(i + length, rand.Int())
		mapping.Delete(i * 2)
	}
	b.Run("pair-half-empty", func(b *testing.B) {
		sum := 0
		for j := 0; j < b.N; j++ {
			mapping.Iterate(func(i int, k int) bool {
				sum += i + k
				return false
			})
		}
		doNothing(sum)
	})
}

func BenchmarkSliceIter(b *testing.B) {
	debug.SetGCPercent(-1)
	slice := make([]int, length)
	for i := 0; i < length; i++ {
		slice[i] = rand.Int()
	}
	runtime.GC()
	b.Run("pair", func(b *testing.B) {
		sum := 0
		for j := 0; j < b.N; j++ {
			for i, k := range slice {
				sum += i + k
			}
		}
		doNothing(sum)
	})
	runtime.GC()
	b.Run("key", func(b *testing.B) {
		sum := 0
		for j := 0; j < b.N; j++ {
			for i := range slice {
				sum += i + slice[i]
			}
		}
		doNothing(sum)
	})
}

func BenchmarkSetDeleteIter(b *testing.B) {
	debug.SetGCPercent(-1)
	mapping := make(map[int]int)
	mapxy := example.NewXYMapIntInt()
	for i := 0; i < length; i++ {
		value := rand.Int()
		mapxy.Set(i, value)
	}
	runtime.GC()
	b.Run("map", func(b *testing.B) {
		for j := 0; j < b.N; j++ {
			for i := 0; i < length; i++ {
				mapping[i] = i
			}
			for i := 0; i < length; i++ {
				delete(mapping, i)
			}
		}
	})
	runtime.GC()
	b.Run("mapxy", func(b *testing.B) {
		for j := 0; j < b.N; j++ {
			for i := 0; i < length; i++ {
				mapxy.Set(i, i)
			}
			for i := 0; i < length; i++ {
				mapxy.Delete(i)
			}
		}
	})
}

func TestIdentify(t *testing.T) {
	tt := func(t *testing.T) {
		mapping := make(map[int]int)
		mapxy := example.NewXYMapIntInt()
		for i := 0; i < length; i++ {
			key := rand.Int()
			value := rand.Int()
			mapping[key] = value
			_, exist := mapxy.Set(key, value)
			assert.Equal(t, exist, false)
		}
		sumM := 0
		sumXY := 0
		for i, v := range mapping {
			sumM += i + v
		}
		mapxy.Iterate(func(k, v int) bool {
			sumXY += k + v
			return false
		})
		assert.Equal(t, sumM, sumXY)
		// test for several empty slot
		for k, v := range mapping {
			vv, exist := mapxy.Get(k)
			assert.Equal(t, exist, true)
			assert.Equal(t, v, vv)
			if rand.Int() % 20 > 18 {
				delete(mapping, k)
				vv, exist := mapxy.Delete(k)
				assert.Equal(t, exist, true)
				assert.Equal(t, v, vv)
			}
		}
		t.Logf("Delete to %d/%d", len(mapping), length)
		sumM = 0
		sumXY = 0
		for i, v := range mapping {
			sumM += i + v
		}
		mapxy.Iterate(func(k, v int) bool {
			sumXY += k + v
			return false
		})
		assert.Equal(t, sumM, sumXY)
		// test for compress
		for k, v := range mapping {
			vv, exist := mapxy.Get(k)
			assert.Equal(t, exist, true)
			assert.Equal(t, v, vv)
			if rand.Int() % 20 > 1 {
				delete(mapping, k)
				vv, exist := mapxy.Delete(k)
				assert.Equal(t, exist, true)
				assert.Equal(t, v, vv)
			}
		}
		t.Logf("Delete to %d/%d", len(mapping), length)
		sumM = 0
		sumXY = 0
		for i, v := range mapping {
			sumM += i + v
		}
		mapxy.Iterate(func(k, v int) (broken bool) {
			sumXY += k + v
			assert.Equal(t, v, mapping[k])
			return
		})
		assert.Equal(t, sumM, sumXY)
		for i := 0; i < 100; i++ {
			v, vExist := mapping[i]
			vv, vvExist := mapxy.Get(i)
			assert.Equal(t, v, vv)
			assert.Equal(t, vExist, vvExist)
		}
	}
	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("run %d", i), tt)
	}
}
