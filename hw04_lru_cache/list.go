package hw04lrucache

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

type list struct {
	head   *ListItem
	tail   *ListItem
	length int
}

func NewList() List {
	return new(list)
}

// Len длина списка.
func (l *list) Len() int {
	return l.length
}

// Front первый элемент списка.
func (l *list) Front() *ListItem {
	return l.head
}

// Back последний элемент списка.
func (l *list) Back() *ListItem {
	return l.tail
}

// PushFront добавить значение в начало.
func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{Value: v}
	if l.head != nil {
		newItem.Next = l.head
		l.head.Prev = newItem
	} else {
		l.tail = newItem
	}
	l.head = newItem
	l.length++
	return newItem
}

// PushBack добавить значение в конец.
func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{Value: v}
	if l.tail != nil {
		newItem.Prev = l.tail
		l.tail.Next = newItem
	} else {
		l.head = newItem
	}
	l.tail = newItem
	l.length++
	return newItem
}

// Remove удалить элемент.
func (l *list) Remove(i *ListItem) {
	if i.Prev == nil {
		l.head = i.Next
		i.Next.Prev = nil
	} else {
		i.Prev.Next = i.Next
	}

	if i.Next == nil {
		l.tail = i.Prev
		i.Prev.Next = nil
	} else {
		i.Next.Prev = i.Prev
	}

	l.length--
}

// MoveToFront переместить элемент в начало.
func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}
