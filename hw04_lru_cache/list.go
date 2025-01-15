package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(value interface{}) *ListItem
	PushBack(value interface{}) *ListItem
	Remove(item *ListItem)
	MoveToFront(item *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	firstItem *ListItem
	lastItem  *ListItem
	len       int
}

func NewList() List {
	return new(list)
}

func (thisList *list) Len() int {
	return thisList.len
}

func (thisList *list) Front() *ListItem {
	return thisList.firstItem
}

func (thisList *list) Back() *ListItem {
	return thisList.lastItem
}

func (thisList *list) PushFront(value interface{}) *ListItem {
	newItem := &ListItem{
		Value: value,
		Next:  thisList.Front(),
		Prev:  nil,
	}

	if thisList.Front() == nil {
		thisList.lastItem = newItem
		thisList.firstItem = newItem
	} else {
		thisList.firstItem.Prev = newItem
		thisList.firstItem = newItem
	}

	thisList.len++

	return newItem
}

func (thisList *list) PushBack(value interface{}) *ListItem {
	newItem := &ListItem{
		Value: value,
		Next:  nil,
		Prev:  thisList.Back(),
	}

	if thisList.Back() == nil {
		thisList.lastItem = newItem
		thisList.firstItem = newItem
	} else {
		thisList.lastItem.Next = newItem
		thisList.lastItem = newItem
	}

	thisList.len++

	return newItem
}

func (thisList *list) Remove(item *ListItem) {
	switch {
	case item.Prev == nil:
		item.Next.Prev = nil
		thisList.firstItem = item.Next
	case item.Next == nil:
		item.Prev.Next = nil
		thisList.lastItem = item.Prev
	default:
		item.Prev.Next = item.Next
		item.Next.Prev = item.Prev
	}

	thisList.len--
}

func (thisList *list) MoveToFront(item *ListItem) {
	if item != nil && item.Prev != nil {
		thisList.Remove(item)
		// Т.к. Remove(item) уменьшает len - прибавляем назад
		thisList.len++

		item.Prev = nil
		item.Next = thisList.Front()
		thisList.firstItem.Prev = item
		thisList.firstItem = item
	}
}
