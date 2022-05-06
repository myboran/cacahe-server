package tcp

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"mycache/httprest/cache"
	"net"
	"strconv"
	"strings"
)

type Server struct {
	cache.Cache
}

func New(c cache.Cache) *Server {
	return &Server{c}
}

func (s *Server) Listen() {
	l, e := net.Listen("tcp", ":12346")
	if e != nil {
		panic(e)
	}
	for {
		c, e := l.Accept()
		if e != nil {
			panic(e)
		}
		go s.process(c)
	}
}

func (s *Server) process(conn net.Conn) {
	defer conn.Close()
	// 读取一个字节
	r := bufio.NewReader(conn)
	for {
		op, e := r.ReadByte()
		if e != nil {
			log.Println("close conn due to error: ", e)
		}
		switch op {
		case 'S':
			e = s.set(conn, r)
		case 'G':
			e = s.get(conn, r)
		case 'D':
			e = s.del(conn, r)
		default:
			log.Println("close conn due to invalid operation:", op)
			return
		}
		if e != nil {
			log.Println("close conn due to error: ", e)
			return
		}
	}
}

func (s *Server) set(conn net.Conn, r *bufio.Reader) error {
	// S<klen><SP><vlen><SP><key><value>
	k, v, err := s.readKeyAndValue(r)
	if err != nil {
		return err
	}
	return s.sendResponse(v, s.Set(k, v), conn)
}

func (s *Server) get(conn net.Conn, r *bufio.Reader) error {
	k, err := s.readKey(r)
	if err != nil {
		return err
	}
	v, err := s.Get(k)
	return s.sendResponse(v, err, conn)
}

func (s *Server) del(conn net.Conn, r *bufio.Reader) error {
	k, err := s.readKey(r)
	if err != nil {
		return err
	}
	return s.sendResponse(nil, s.Del(k), conn)
}

func (s *Server) readKeyAndValue(r *bufio.Reader) (string, []byte, error) {
	// S<klen><SP><vlen><SP><key><value>
	klen, err := s.readLen(r)
	if err != nil {
		return "", nil, err
	}
	vlen, err := s.readLen(r)
	if err != nil {
		return "", nil, err
	}
	k := make([]byte, klen)
	// ReadFull 将 r 中的 len(buf) 个字节准确地读取到 buf 中
	_, err = io.ReadFull(r, k)
	if err != nil {
		return "", nil, err
	}
	v := make([]byte, vlen)
	_, err = io.ReadFull(r, v)
	if err != nil {
		return "", nil, err
	}
	return string(k), v, nil
}

func (s *Server) readLen(r *bufio.Reader) (int, error) {
	// ReadString 读取直到输入中第一次出现 delim，返回一个字符串，其中包含直到并包括分隔符的数据
	tmp, err := r.ReadString(' ')
	if err != nil {
		return 0, err
	}
	// TrimSpace 返回字符串 s 的切片，删除所有前导和尾随空格
	l, err := strconv.Atoi(strings.TrimSpace(tmp))
	if err != nil {
		return 0, err
	}
	return l, nil
}

func (s *Server) sendResponse(value []byte, err error, conn net.Conn) error {
	if err != nil {
		errString := err.Error()
		tmp := fmt.Sprintf("-%d", len(errString)) + errString
		_, e := conn.Write([]byte(tmp))
		return e
	}
	vlen := fmt.Sprintf("%d", len(value))
	_, e := conn.Write(append([]byte(vlen), value...))
	return e
}

func (s *Server) readKey(r *bufio.Reader) (string, error) {
	klen, err := s.readLen(r)
	if err != nil {
		return "", err
	}
	k := make([]byte, klen)
	_, err = io.ReadFull(r, k)
	if err != nil {
		return "", err
	}
	return string(k), nil
}
