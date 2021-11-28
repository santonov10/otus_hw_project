package lrucache

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

var _ List = &list{}

type list struct {
	len           int
	firstListItem *ListItem
	lastListItem  *ListItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.firstListItem
}

func (l *list) Back() *ListItem {
	return l.lastListItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	if l.Len() == 0 {
		return l.pushToEmptyList(v)
	}
	newListItem := NewListItem(v)
	secondListItem := l.firstListItem
	secondListItem.Prev = newListItem
	l.firstListItem = newListItem
	l.firstListItem.Next = secondListItem
	l.len++
	return newListItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	if l.Len() == 0 {
		return l.pushToEmptyList(v)
	}
	newListItem := NewListItem(v)
	preLastListItem := l.lastListItem
	preLastListItem.Next = newListItem
	l.lastListItem = newListItem
	l.lastListItem.Prev = preLastListItem
	l.len++
	return newListItem
}

func (l *list) pushToEmptyList(v interface{}) *ListItem {
	newListItem := NewListItem(v)
	l.firstListItem = newListItem
	l.lastListItem = newListItem
	l.len = 1
	return newListItem
}

func (l *list) Remove(i *ListItem) {
	if l.lastListItem == i {
		l.lastListItem = i.Prev
	}
	if l.firstListItem == i {
		l.firstListItem = i.Next
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.firstListItem != i {
		l.Remove(i)
		l.PushFront(i)
	}
}

func NewList() List {
	return new(list)
}

func NewListItem(v interface{}) *ListItem {
	newListItem := &ListItem{}
	switch t := v.(type) {
	case *ListItem:
		newListItem = v.(*ListItem)
	default:
		newListItem.Value = t
	}
	return newListItem
}
