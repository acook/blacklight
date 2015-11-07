// use runtime.SetFinalizer to close IOs when their queues are GC'd?
// Or when their FD tag is GC'd?
// https://golang.org/pkg/runtime/#SetFinalizer

package main

import (
	"os"
)

type IO struct {
	Name  string
	Queue *Queue
}

func NewIO(i datatypes, q *Queue) *Tag {
	switch i.(type) {
	case *Number:
		fd := NewFD(i.Value().(int), q)
		return NewTag("FD#"+i.String(), fd)
	case *CharVector:
		file := NewFile(i.String(), q)
		return NewTag("File#"+i.String(), file)
	default:
		panic("NewIO: unrecognized type for IO - " + i.String())
	}
	return nil
}

var FDtable map[uint]*os.File = make(map[uint]*os.File)
var FDtableinit bool

func initFDtable() {
	if !FDtableinit {
		FDtable[0] = os.Stdin
		FDtable[1] = os.Stdout
		FDtable[2] = os.Stderr
		FDtableinit = true
	}
}

type FD struct {
	IO
	FD   uint
	File *os.File
}

func NewFD(i int, q *Queue) *FD {
	initFDtable()
	fd := new(FD)
	fd.Queue = q
	fd.FD = uint(i)
	fd.File = FDtable[uint(i)]

	threads.Add(1)
	go func(fd *FD, q *Queue) {
		defer threads.Done()
		b := make([]byte, 1)
		for {
			l, _ := fd.File.Read(b)
			if l > 0 {
				q.Enqueue(NewCharFromString(string(b)))
			} else {
				fd.File.Close()
				q.Enqueue(NewNil("EOF"))
				return
			}
		}
	}(fd, q)

	return fd
}

func NewFile(filename string, q *Queue) *FD {
	fd := new(FD)
	fd.Queue = q
	fd.File, _ = os.Open(filename)
	fd.FD = uint(fd.File.Fd())

	FDtable[fd.FD] = fd.File

	threads.Add(1)
	go func(fd *FD, q *Queue) {
		defer threads.Done()
		b := make([]byte, 1)
		for {
			l, _ := fd.File.Read(b)
			if l > 0 {
				q.Enqueue(NewCharFromString(string(b)))
			} else {
				fd.File.Close()
				q.Enqueue(NewNil("EOF"))
				return
			}
		}
	}(fd, q)

	return fd
}
