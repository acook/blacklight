package main

import (
	"os"
)

const (
	IO_READ = iota
	IO_WRITE
)

type IO struct {
	Name  T
	Queue *Queue
	FD    uint
	File  *os.File
	Mode  uint // 01 read, 10 write
}

var FDtable map[uint]*IO = make(map[uint]*IO)

func initFDtable() {
	stdin := new(IO)
	stdin.Name = "stdin"
	stdin.FD = 0
	stdin.File = os.Stdin
	stdin.Mode = IO_READ

	FDtable[0] = stdin

	stdout := new(IO)
	stdout.Name = "stdout"
	stdout.FD = 1
	stdout.File = os.Stdout
	stdout.Mode = IO_WRITE

	FDtable[1] = stdout

	stderr := new(IO)
	stderr.Name = "stderr"
	stderr.FD = 2
	stderr.File = os.Stderr
	stderr.Mode = IO_WRITE

	FDtable[2] = stderr
}

func ReadFD(i int, q *Queue) *IO {
	fd := FDtable[uint(i)]
	fd.Queue = q

	threads.Add(1)
	go func(fd *IO, q *Queue) {
		defer threads.Done()

		for b := make([]byte, 1); ; {
			l, _ := fd.File.Read(b)
			if l > 0 {
				q.Enqueue(R(string(b)[0]))
			} else {
				fd.File.Close()
				q.Enqueue(NewNil("EOF"))
				return
			}
		}
	}(fd, q)

	return fd
}

func WriteFD(i int, q *Queue) *IO {
	fd := FDtable[uint(i)]
	fd.Queue = q

	threads.Add(1)
	go func(fd *IO, q *Queue) {
		defer threads.Done()
		var b []byte

		for {
			b = q.Dequeue().(byter).Bytes()

			if b == nil {
				fd.File.Close()
				return
			} else {
				l, _ := fd.File.Write(b)
				if l < len(b) {
					panic("WriteFile: Write Error!")
				}
			}
		}
	}(fd, q)

	return fd
}

func ReadFile(filename T, q *Queue) *Tag {
	fd := new(IO)
	fd.Mode = IO_READ
	fd.Queue = q
	fd.File, _ = os.Open(string(filename))
	fd.FD = uint(fd.File.Fd())

	FDtable[fd.FD] = fd

	threads.Add(1)
	go func(fd *IO, q *Queue) {
		defer threads.Done()
		b := make([]byte, 1)

		for {
			l, _ := fd.File.Read(b)
			if l > 0 {
				q.Enqueue(C(b[0]))
			} else {
				fd.File.Close()
				q.Enqueue(NewNil("EOF"))
				return
			}
		}
	}(fd, q)

	return NewFileTag("File#"+filename.Print(), fd)
}

func WriteFile(filename T, q *Queue) *Tag {
	fd := new(IO)
	fd.Mode = IO_WRITE
	fd.Queue = q
	fd.File, _ = os.Create(string(filename))
	fd.FD = uint(fd.File.Fd())

	FDtable[fd.FD] = fd

	threads.Add(1)
	go func(fd *IO, q *Queue) {
		defer threads.Done()
		var b []byte

		for {
			b = q.Dequeue().(byter).Bytes()

			if b == nil {
				fd.File.Close()
				return
			} else {
				l, _ := fd.File.Write(b)
				if l < len(b) {
					panic("WriteFile: Write Error!")
				}
			}
		}
	}(fd, q)

	return NewFileTag("File#"+filename.Print(), fd)
}
