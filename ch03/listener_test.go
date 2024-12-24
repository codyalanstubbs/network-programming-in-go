package ch03

import (
	"io"
	"net"
	"syscall"
	"testing"
	"time"
)

func TestListener(t *testing.T) {
  listener, err := net.Listen("tcp", "127.0.0.1:0")
  if err != nil {
    t.Fatal(err)
  }

  defer func() { _ = listener.Close() }()

  t.Logf("bound to %q", listener.Addr())
}

func TestDialer(t *testing.T) {
  listener, err := net.Listen("tcp", "127.0.0.1:")
  if err != nil {
    t.Fatal(err)
  }

  done := make(chan struct{})

  go func()  {
    defer func() { done <- struct{}{} }()

    for {
      conn, err := listener.Accept()

      if err != nil {
        t.Log(err)
        return
      }

      go func(c net.Conn) {
        defer func() {
          c.Close()
          done <- struct{}{}
        }()

        buf := make([]byte, 1024)
        for {
          n, err := c.Read(buf)
          if err != nil {
            if err != io.EOF {
              t.Error(err)
            }
            return
          }

          t.Logf("received: %q", buf[:n])
        }

      }(conn)
    }
    
  }()

  conn, err := net.Dial("tcp", listener.Addr().String())

  if err != nil {
    t.Fatal(err)
  }

  conn.Close()
  <-done
  listener.Close()
  <-done

}

func DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
  d := net.Dialer{
      Control: func(_, addr string, _ syscall.RawConn) error {
          return &net.DNSError{
        Err: "connection timed out",
        Name: addr,
        Server: "127.0.0.1",
        IsTimeout: true,
        IsTemporary: true,
      }
      },
    Timeout: timeout,
  }

  return d.Dial(network, address)
}
