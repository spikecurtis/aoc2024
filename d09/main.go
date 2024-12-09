package main

import (
	"fmt"

	"github.com/spikecurtis/aoc2024/util"
)

type file struct {
	prev, next *file
	length     int
	id         int32
}

func (f *file) append(e *file) {
	l := f.prev
	f.prev = e
	l.next = e
	e.prev = l
	e.next = f
}

func (f *file) insertBefore(e *file) {
	f.append(e)
}

func (f *file) write(id int32, length int) {
	if f.id != -1 {
		panic("overwrite")
	}
	if f.length == length {
		f.id = id
		return
	}
	n := &file{
		length: length,
		id:     id,
	}
	f.insertBefore(n)
	f.length = f.length - length
}

func (f *file) remove() {
	f.id = -1
}

func main() {
	lines := util.GetInputLines()
	diskCompact := lines[0]
	//diskCompact = "2333133121414131402"
	diskRaw := make([]int32, 0)
	fID := int32(0)
	for i, b := range diskCompact {
		n := int(b - '0')
		w := fID
		if i%2 == 1 {
			w = -1
		} else {
			fID++
		}
		for range n {
			diskRaw = append(diskRaw, w)
		}
	}
	// compress pt 1
	s := 0
	e := len(diskRaw) - 1
	for s < e {
		if diskRaw[s] != -1 {
			s++
			continue
		}
		if diskRaw[e] == -1 {
			e--
			continue
		}
		diskRaw[s] = diskRaw[e]
		e--
		s++
	}
	diskRaw = diskRaw[:e+1]
	// checksum
	p1 := int64(0)
	for i, f := range diskRaw {
		if f == -1 {
			continue
		}
		p1 += int64(i) * int64(f)
	}
	fmt.Println("Part 1: ", p1)

	// part 2
	var first *file
	fID = 0
	for i, b := range diskCompact {
		n := int(b - '0')
		if n == 0 {
			continue
		}
		w := fID
		if i%2 == 1 {
			w = -1
		} else {
			fID++
		}
		f := &file{
			length: n,
			id:     w,
		}
		if w == 0 {
			first = f
			f.prev = first
			f.next = first
		} else {
			first.append(f)
		}
	}
	// compress
	cur := first.prev
files:
	for cur != first {
		if cur.id == -1 {
			cur = cur.prev
			continue
		}
		target := first
		for target != cur {
			if target.id == -1 && target.length >= cur.length {
				target.write(cur.id, cur.length)
				toRemove := cur
				cur = cur.prev
				toRemove.remove()
				continue files
			}
			target = target.next
		}
		cur = cur.prev
	}
	p2 := 0
	block := 0
	cur = first
	for {
		if cur.id != -1 {
			for b := block; b < block+cur.length; b++ {
				p2 += b * int(cur.id)
			}
		}
		block += cur.length
		cur = cur.next
		if cur == first {
			break
		}
	}
	fmt.Println("Part 2: ", p2)
}
