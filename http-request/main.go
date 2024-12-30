package main

import (
	"github.com/pkg/errors"
	"github.com/soypat/cyw43439/examples/common"
	"github.com/soypat/seqs"
	"github.com/soypat/seqs/httpx"
	"github.com/soypat/seqs/stacks"
	"math/rand"
	"net/netip"
	"time"

	"log/slog"

	"machine"
)

const connTimeout = 5 * time.Second
const tcpbufsize = 2030 // MTU - ethhdr - iphdr - tcphdr
const ourHostname = "tinygo-http-client"

type Connection struct {
	serverAddrStr string
	conn          *stacks.TCPConn
	closeConn     func(err error)
	svAddr        netip.AddrPort
	rng           *rand.Rand
	logger        *slog.Logger
	clientAddr    netip.AddrPort
	routerhw      [6]byte
}

type RetryConfig struct {
	delay      time.Duration
	retryCount int
}

func main() {
	connection := NewConnection("192.168.1.201:8080")

	response, err := connection.sendRequestWithRetry("GET", RetryConfig{
		delay:      5 * time.Second,
		retryCount: 10,
	})
	if err != nil {
		println(err.Error())
	} else {
		println(response)
	}
}

func NewConnection(targetAddress string) *Connection {
	logger := slog.New(slog.NewTextHandler(machine.Serial, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	_, stack, _, err := common.SetupWithDHCP(common.SetupConfig{
		Hostname: ourHostname,
		Logger:   logger,
		TCPPorts: 1, // For HTTP over TCP.
		UDPPorts: 1, // For DNS.
	})
	start := time.Now()
	if err != nil {
		panic("setup DHCP:" + err.Error())
	}
	svAddr, err := netip.ParseAddrPort(targetAddress)
	if err != nil {
		panic("parsing server address:" + err.Error())
	}
	// Resolver router's hardware address to dial outside our network to internet.
	routerhw, err := common.ResolveHardwareAddr(stack, svAddr.Addr())
	if err != nil {
		panic("router hwaddr resolving:" + err.Error())
	}

	rng := rand.New(rand.NewSource(int64(time.Now().Sub(start))))
	// Start TCP server.
	clientAddr := netip.AddrPortFrom(stack.Addr(), uint16(rng.Intn(65535-1024)+1024))
	conn, err := stacks.NewTCPConn(stack, stacks.TCPConnConfig{
		TxBufSize: tcpbufsize,
		RxBufSize: tcpbufsize,
	})

	if err != nil {
		panic("conn create:" + err.Error())
	}

	closeConn := func(err error) {
		if err != nil {
			slog.Error("tcpconn:closing", slog.String("err", err.Error()))
		}
		conn.Close()
		for !conn.State().IsClosed() {
			slog.Info("tcpconn:waiting", slog.String("state", conn.State().String()))
			time.Sleep(1000 * time.Millisecond)
		}
	}

	return &Connection{
		serverAddrStr: targetAddress,
		rng:           rng,
		conn:          conn,
		closeConn:     closeConn,
		svAddr:        svAddr,
		clientAddr:    clientAddr,
		routerhw:      routerhw,
		logger:        logger,
	}
}

func (c *Connection) sendRequestWithRetry(method string, retry RetryConfig) (response string, err error) {
	defer func() {
		c.closeConn(err)
	}()
	for i := 0; i < retry.retryCount; i++ {
		c.logger.Info("dialing", slog.String("serveraddr", c.serverAddrStr))
		response, err = c.sendRequest(method)
		if err == nil {
			return response, nil
		}
		c.logger.Info("send request failed. Got %v", err)
		c.conn.SetDeadline(time.Time{}) // Disable the deadline.
		time.Sleep(retry.delay)
	}
	return "", errors.New("retry limit reached")
}

func (c *Connection) sendRequest(method string) (string, error) {
	// Here we create the HTTP request and generate the bytes. The Header method
	// returns the raw header bytes as should be sent over the wire.
	var req httpx.RequestHeader
	req.SetRequestURI("/")
	// If you need a Post request change "GET" to "POST" and then add the
	// post data to reqbytes: `postReq := append(reqbytes, postData...)` and send postReq over TCP.
	req.SetMethod(method)
	req.SetHost(c.svAddr.Addr().String())
	req.SetUserAgent(ourHostname)
	reqbytes := req.Header()
	c.logger.Info("tcp:ready",
		slog.String("clientaddr", c.clientAddr.String()),
		slog.String("serveraddr", c.serverAddrStr),
	)

	// Make sure to timeout the connection if it takes too long.
	c.conn.SetDeadline(time.Now().Add(connTimeout))
	err := c.conn.OpenDialTCP(c.clientAddr.Port(), c.routerhw, c.svAddr, seqs.Value(c.rng.Intn(65535-1024)+1024))
	if err != nil {
		return "", errors.Wrap(err, "opening TCP")
	}

	retries := 50
	for c.conn.State() != seqs.StateEstablished && retries > 0 {
		time.Sleep(100 * time.Millisecond)
		retries--
	}
	c.logger.Info("tcp connection state: ", slog.Uint64("state", uint64(c.conn.State())))
	c.conn.SetDeadline(time.Time{}) // Disable the deadline.
	if retries == 0 {
		return "", errors.New("wait for connection established exceeded")
	}

	return c.do(reqbytes)
}

func (c *Connection) do(reqBytes []byte) (string, error) {
	rxBuf := make([]byte, 4096)
	// Send the request.
	_, err := c.conn.Write(reqBytes)
	if err != nil {
		return "", errors.Wrap(err, "writing to TCP")
	}
	time.Sleep(500 * time.Millisecond)
	c.conn.SetDeadline(time.Now().Add(connTimeout))
	n, err := c.conn.Read(rxBuf)
	if n == 0 && err != nil {
		return "", errors.Wrap(err, "reading response")
	} else if n == 0 {
		return "", errors.Wrap(err, "no response")
	}
	println("got HTTP response!")
	c.closeConn(nil)
	return string(rxBuf[:n]), nil
}
