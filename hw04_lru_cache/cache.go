package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int               // Емкость кэша
	queue    List              // Очередь последних элементов
	items    map[Key]*ListItem // Словарь ключей с элементами
	mutex    sync.RWMutex      // Для горутино-безопасности
}

// Структура для хранения элемента в списке.
type cacheItem struct {
	key   Key
	value interface{}
}

// Создание кэша с установкой емкости.
func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

// Установка значения в кэш.
func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Если емкость кэша 0 или меньше то выходим
	if c.capacity <= 0 {
		return false
	}

	// Проверяем существующий элемент.
	// Обновляем значение и перемещаем в начало.
	if item, exists := c.items[key]; exists {
		item.Value = &cacheItem{key: key, value: value}
		c.queue.MoveToFront(item)
		return true
	}

	// Если достигнут лимит, удаляем последний элемент.
	if c.queue.Len() >= c.capacity {
		lastItem := c.queue.Back()
		if lastCacheItem, ok := lastItem.Value.(*cacheItem); ok {
			delete(c.items, lastCacheItem.key)
		}
		c.queue.Remove(lastItem)
	}

	// Добавляем новый элемент.
	newItem := c.queue.PushFront(&cacheItem{
		key:   key,
		value: value,
	})
	c.items[key] = newItem
	return false
}

// Чтение из кэша по ключу.
func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Проверка ключа на его наличие, и присваение значения если существует
	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	// Перемечаем найденое значение в начало списка
	c.queue.MoveToFront(item)
	// Возвращаем значение без ключа
	if cacheItem, ok := item.Value.(*cacheItem); ok {
		return cacheItem.value, true
	}

	return nil, false
}

// Очистка кэша.
func (c *lruCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
