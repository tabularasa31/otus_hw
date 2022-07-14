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
	newItem := &ListItem{Value: v, Next: l.head, Prev: nil}
	if l.head != nil {
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
	switch {
	case i == l.head:
		i.Next.Prev = nil
		l.head = i.Next
	case i == l.tail:
		i.Prev.Next = nil
		l.tail = i.Prev
	default:
		i.Next.Prev = i.Prev
		i.Prev.Next = i.Next
	}
	//	i = nil
	l.length--
}

// MoveToFront переместить элемент в начало.
func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)

	l.head.Prev, i.Next, l.head = i, l.head, i
	l.length++
}
