package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	// Дополнительный тест на логику работы.
	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(3)

		c.Set("a", 1)
		c.Set("b", 2)
		c.Set("c", 3)
		// Должен вытеснить "a"
		c.Set("d", 4)

		// Проверка наличия "a"
		_, ok := c.Get("a")
		require.False(t, ok)

		// Считываю давний элемент и проверяю, что он существует.
		val, ok := c.Get("b")
		require.True(t, ok)
		require.Equal(t, 2, val)

		// Добавление нового элемента
		c.Set("e", 5)
		// Проверка на вытеснение "c", который был затронут наиболее давно.
		_, ok = c.Get("c")
		require.False(t, ok)
	})

	// Тест на кэш с нулевой емкостью.
	t.Run("null capacity", func(t *testing.T) {
		c := NewCache(0)

		c.Set("a", 1)
		c.Set("b", 2)
		c.Set("c", 3)

		_, ok := c.Get("a")
		require.False(t, ok)
	})

	// Дополнительный тест на очистку кэша.
	t.Run("clear", func(t *testing.T) {
		c := NewCache(5)

		c.Set("a", 1)
		c.Set("b", 2)

		c.Clear()

		_, ok := c.Get("a")
		require.False(t, ok)
		_, ok = c.Get("b")
		require.False(t, ok)
	})
}

func TestCacheMultithreading(_ *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
