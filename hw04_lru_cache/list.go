package hw04lrucache

import "sync"

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

// Структура списка.
type list struct {
	length int          // Длина списка
	front  *ListItem    // Указатель на первый элемент списка
	back   *ListItem    // Указатель на последний элемент списка
	mutex  sync.RWMutex // Для горутино-безопасности
}

// Создание списка.
func NewList() List {
	return new(list)
}

// Создание нового элемента списка.
func (l *list) newItem(v interface{}) (*ListItem, bool) {
	item := new(ListItem)
	item.Value = v

	// Если элементы отсутствуют, то создаем первый.
	if l.length <= 0 {
		l.front = item
		l.back = item
		l.length = 1
		return item, true
	}

	return item, false
}

// Возврат длины списка.
func (l *list) Len() int {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.length
}

// Возвращение первого элемента списка.
func (l *list) Front() *ListItem {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.front
}

// Возвращение последнего элемента списка.
func (l *list) Back() *ListItem {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.back
}

// Запись данных в начало списка и возврат нового элемента.
func (l *list) PushFront(v interface{}) *ListItem {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	// Создание нового элемента списка.
	item, ok := l.newItem(v)
	if ok {
		return item
	}

	// Обновление полей в персом элементе и запись нового в начало.
	item.Next = l.front
	l.front.Prev = item
	l.front = item
	l.length++
	return item
}

// Запись данных в конец списка и возврат нового элемента.
func (l *list) PushBack(v interface{}) *ListItem {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	// Создание нового элемента списка.
	item, ok := l.newItem(v)
	if ok {
		return item
	}

	// Обновление полей в последнем элементе и запись нового в конец.
	l.back.Next = item
	item.Prev = l.back
	l.back = item
	l.length++
	return item
}

// Удаление элемента из списка.
func (l *list) Remove(i *ListItem) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	// Проверка на пустой входной параметр
	if i == nil {
		return
	}

	// Проверка если элемент единственный в списке.
	// Если это так то очищаем его.
	if (i.Prev == nil) && (i.Next == nil) {
		l.front = nil
		l.back = nil
		l.length = 0
		return
	}

	// Проверка если в начале списка.
	// Если это так то обновляем поля и очищаем его.
	if i.Prev == nil {
		front := l.front
		l.front = front.Next
		l.front.Prev = nil
		front = nil
		l.length--
		return
	}

	// Проверка если в конце списка.
	// Если это так то обновляем поля и очищаем его.
	if i.Next == nil {
		back := l.back
		l.back = back.Prev
		l.back.Next = nil
		back = nil
		l.length--
		return
	}

	// Иначе если элемент посредине, обновляем поля и очищаем его.
	nextItem := i.Next
	nextItem.Prev = i.Prev
	i.Prev.Next = nextItem
	l.length--
}

// Перемещение элемента списка в начало.
func (l *list) MoveToFront(i *ListItem) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	// Проверка на пустой входной параметр
	if i == nil {
		return
	}

	// Проверка если элемент уже в начале списка.
	if i.Prev == nil {
		return
	}

	// Обновление промежуточных полей списка.
	// Если он находится в конце, иначе посредине.
	if i.Next == nil {
		itemPrev := i.Prev
		itemPrev.Next = nil
		l.back = itemPrev
	} else {
		prevItem := i.Prev
		nextItem := i.Next
		prevItem.Next = nextItem
		nextItem.Prev = prevItem
	}

	// Обнволение полей первого элемента списка.
	i.Prev = nil
	i.Next = l.front
	l.front.Prev = i
	l.front = i
}
