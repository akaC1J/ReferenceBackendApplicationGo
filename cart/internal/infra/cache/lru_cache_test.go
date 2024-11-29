package cache

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLruCache_BasicOperations(t *testing.T) {
	c := NewLruCache[int, string](10)

	// Добавляем элементы
	c.Put(1, "one")
	c.Put(2, "two")

	// Проверяем, что элементы корректно добавлены
	val, ok := c.Get(1)
	assert.True(t, ok, "1 should exist in cache")
	assert.Equal(t, "one", val)

	val, ok = c.Get(2)
	assert.True(t, ok, "2 should exist in cache")
	assert.Equal(t, "two", val)
}

func TestLruCache_UpdateExisting(t *testing.T) {
	c := NewLruCache[int, string](10)

	c.Put(1, "one")
	c.Put(1, "uno") // Обновляем значение

	val, ok := c.Get(1)
	assert.True(t, ok, "1 should exist in cache")
	assert.Equal(t, "uno", val, "Value for key 1 should be updated to 'uno'")
}

func TestLruCache_Overflow(t *testing.T) {
	c := NewLruCache[int, string](2)

	c.Put(1, "one")
	c.Put(2, "two")
	c.Put(3, "three") // Добавляем, что превышает емкость

	// Проверяем, что наименее недавно использованный элемент (1) был удален
	_, ok := c.Get(1)
	assert.False(t, ok, "1 should be evicted from cache")

	// Проверяем, что остальные элементы остались
	val, ok := c.Get(2)
	assert.True(t, ok, "2 should still exist in cache")
	assert.Equal(t, "two", val)

	val, ok = c.Get(3)
	assert.True(t, ok, "3 should still exist in cache")
	assert.Equal(t, "three", val)
}

func TestLruCache_MoveToHead(t *testing.T) {
	c := NewLruCache[int, string](2)

	c.Put(1, "one")
	c.Put(2, "two")
	c.Get(1)          // Делаем 1 самым недавно использованным
	c.Put(3, "three") // Превышаем емкость

	// Проверяем, что элемент 2 был удален, так как он стал самым старым
	_, ok := c.Get(2)
	assert.False(t, ok, "2 should be evicted from cache")

	// Проверяем, что 1 и 3 остались
	val, ok := c.Get(1)
	assert.True(t, ok, "1 should still exist in cache")
	assert.Equal(t, "one", val)

	val, ok = c.Get(3)
	assert.True(t, ok, "3 should still exist in cache")
	assert.Equal(t, "three", val)
}

func TestLruCache_EmptyCache(t *testing.T) {
	c := NewLruCache[int, string](10)

	// Проверяем, что в пустом кэше ничего нет
	_, ok := c.Get(1)
	assert.False(t, ok, "Empty cache should not have any elements")
}

func TestLruCache_ZeroCapacity(t *testing.T) {
	c := NewLruCache[int, string](0)

	c.Put(1, "one")
	_, ok := c.Get(1)
	assert.False(t, ok, "Cache with zero capacity should not store elements")
}

func TestLruCache_SingleElement(t *testing.T) {
	c := NewLruCache[int, string](1)

	c.Put(1, "one")
	val, ok := c.Get(1)
	assert.True(t, ok, "1 should exist in cache")
	assert.Equal(t, "one", val)

	// Добавляем другой элемент, удаляя первый
	c.Put(2, "two")
	_, ok = c.Get(1)
	assert.False(t, ok, "1 should be evicted from cache")
	val, ok = c.Get(2)
	assert.True(t, ok, "2 should exist in cache")
	assert.Equal(t, "two", val)
}

func TestLruCache_CustomKeyType(t *testing.T) {
	type Key struct {
		ID   int
		Name string
	}

	c := NewLruCache[Key, string](10)

	key1 := Key{ID: 1, Name: "one"}
	key2 := Key{ID: 2, Name: "two"}

	c.Put(key1, "value1")
	c.Put(key2, "value2")

	val, ok := c.Get(key1)
	assert.True(t, ok, "Custom key 1 should exist in cache")
	assert.Equal(t, "value1", val)

	val, ok = c.Get(key2)
	assert.True(t, ok, "Custom key 2 should exist in cache")
	assert.Equal(t, "value2", val)
}

func TestLruCache_FrequentAccess(t *testing.T) {
	c := NewLruCache[int, string](3)

	c.Put(1, "one")
	c.Put(2, "two")
	c.Put(3, "three")
	c.Get(1)         // Делаем 1 самым недавно использованным
	c.Put(4, "four") // Превышаем емкость

	// Проверяем, что 2 был удален
	_, ok := c.Get(2)
	assert.False(t, ok, "2 should be evicted from cache")

	// Проверяем, что остальные элементы остались
	val, ok := c.Get(1)
	assert.True(t, ok, "1 should still exist in cache")
	assert.Equal(t, "one", val)

	val, ok = c.Get(3)
	assert.True(t, ok, "3 should still exist in cache")
	assert.Equal(t, "three", val)

	val, ok = c.Get(4)
	assert.True(t, ok, "4 should still exist in cache")
	assert.Equal(t, "four", val)
}

func TestLruCache_ConcurrentAccess(t *testing.T) {
	c := NewLruCache[int, string](10)
	var wg sync.WaitGroup

	// Количество горутин
	numRoutines := 50
	numOperations := 1000

	// Добавление данных в кэш
	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				c.Put(id*numOperations+j, fmt.Sprintf("value-%d", id*numOperations+j))
			}
		}(i)
	}

	// Одновременное чтение данных из кэша
	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				c.Get(id * numOperations)
			}
		}(i)
	}

	wg.Wait()
	t.Log("Concurrent test completed without deadlocks or crashes")
}

func TestLruCache_ConcurrentPutAndGet(t *testing.T) {
	c := NewLruCache[int, string](10)

	// Использ ожидания завершения всех горутин
	var wg sync.WaitGroup

	numRoutines := 10
	numOperations := 100

	// Одновременное добавление и чтение
	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				c.Put(j, fmt.Sprintf("value-%d", j))
				c.Get(j)
			}
		}(i)
	}

	wg.Wait()
	t.Log("Concurrent put and get test completed without deadlocks or crashes")
}

func TestLruCache_DataIntegrity(t *testing.T) {
	c := NewLruCache[int, string](10)

	// Проверка на одновременн данных
	var wg sync.WaitGroup

	// Записываем данные
	for i := 0; i < 10; i++ {
		c.Put(i, fmt.Sprintf("value-%d", i))
	}

	numRoutines := 10
	numOperations := 100

	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				for k := 0; k < 10; k++ {
					val, ok := c.Get(k)
					if !ok {
						t.Errorf("Key %d not found", k)
					} else if val != fmt.Sprintf("value-%d", k) {
						t.Errorf("Data corruption detected: expected value-%d, got %s", k, val)
					}
				}
			}
		}()
	}

	wg.Wait()
	t.Log("Data integrity test completed")
}
