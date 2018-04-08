// Copyright 2018 rootkiwi
//
// screen_share_remote_go is licensed under GNU General Public License 3 or later.
//
// See LICENSE for more details.

package remote

import (
	"crypto/tls"
	"encoding/binary"
	"io"
	"log"
	"net"
	"strconv"

	"github.com/rootkiwi/screen_share_remote_go/conf"
	"github.com/rootkiwi/screen_share_remote_go/password"
	"github.com/rootkiwi/screen_share_remote_go/remote/internal/webserver"
)

// Listen waits for connection from screen_share.
func Listen(conf *conf.Config, done chan struct{}) {
	config := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		CipherSuites: []uint16{tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384},
		Certificates: []tls.Certificate{*conf.Cert},
		ClientAuth:   tls.NoClientCert,
	}
	ln, err := tls.Listen("tcp", ":"+strconv.Itoa(conf.Port), config)
	if err != nil {
		log.Fatalf("error listen: %v\n", err)
	}
	defer ln.Close()

	type accepted struct {
		conn net.Conn
		err  error
	}
	newConn := make(chan accepted, 1)

Loop:
	for {
		go func() {
			conn, err := ln.Accept()
			newConn <- accepted{conn, err}
		}()
		select {
		case conn := <-newConn:
			if conn.err != nil {
				continue
			}
			if exit := handleConnection(conf, conn.conn, done); exit {
				break Loop
			}
		case <-done:
			break Loop
		}
	}
	defer func() { done <- struct{}{} }()
}

// handleConnection handles connection from screen_share.
func handleConnection(conf *conf.Config, conn net.Conn, done chan struct{}) (exit bool) {
	defer conn.Close()
	pageTitle, err := doInitialTransfers(conf.PasswordHash, conn)
	if err != nil {
		return false
	}
	entering := make(chan webserver.FrameQueue)
	leaving := make(chan webserver.FrameQueue, 10)
	go webserver.Start(conf.WebPort, pageTitle, entering, leaving)
	defer webserver.Stop()
	frameQueues := make(map[webserver.FrameQueue]struct{})
	closeFrameQueues := func() {
		for queue, _ := range frameQueues {
			close(queue)
		}
	}
	frames := make(chan []byte, 10)
	go readFrames(conn, frames)
	for {
		select {
		case queue := <-entering:
			if err := handleEntering(conn, frameQueues, queue); err != nil {
				return false
			}
		case queue := <-leaving:
			if err := handleLeaving(conn, frameQueues, queue); err != nil {
				return false
			}
		case f, ok := <-frames:
			if !ok {
				closeFrameQueues()
				return false
			}
			for queue, _ := range frameQueues {
				isFull := len(queue) == cap(queue)
				if isFull {
					close(queue)
					delete(frameQueues, queue)
					continue
				}
				queue <- f
			}
		case <-done:
			closeFrameQueues()
			return true
		}
	}
}

const (
	firstNewConnection = 0
	newConnection      = 1
	zeroConnections    = 2
)

// handleEntering handles new WebSocket client connected,
// Notifying screen_share about new connection and puts client frame queue in frameQueues set.
func handleEntering(conn net.Conn, frameQueues map[webserver.FrameQueue]struct{}, queue webserver.FrameQueue) error {
	var message byte
	if len(frameQueues) == 0 {
		message = firstNewConnection
	} else {
		message = newConnection
	}
	if _, err := conn.Write([]byte{message}); err != nil {
		return err
	}
	frameQueues[queue] = struct{}{}
	return nil
}

// handleLeaving handles disconnected WebSocket client connected,
// Notifying screen_share if zeroConnections and removes client frame queue from frameQueues set.
func handleLeaving(conn net.Conn, frameQueues map[webserver.FrameQueue]struct{}, queue webserver.FrameQueue) error {
	if len(frameQueues) == 1 {
		if _, err := conn.Write([]byte{zeroConnections}); err != nil {
			return err
		}
	}
	delete(frameQueues, queue)
	return nil
}

// readFrames read h264 frames from screen_share.
func readFrames(conn net.Conn, frames chan<- []byte) {
	frameSize := make([]byte, 4)
	for {
		_, err := io.ReadFull(conn, frameSize)
		if err != nil {
			close(frames)
			return
		}
		frame := make([]byte, byteArrToInt(frameSize))
		_, err = io.ReadFull(conn, frame)
		if err != nil {
			close(frames)
			return
		}
		frames <- frame
	}
}

// doInitialTransfers validates password and reads pageTitle from screen_share.
func doInitialTransfers(passwordHash string, conn net.Conn) (pageTitle string, err error) {
	const (
		wrongPassword   = 0
		correctPassword = 1
	)
	conn.Write([]byte("im_a_screen_share_remote_server_i_promise"))

	clientPassword, err := readPassword(conn)
	if err != nil {
		return "", err
	}
	if !password.Validate(passwordHash, clientPassword) {
		conn.Write([]byte{wrongPassword})
		log.Printf("wrong password attempt from: %s\n", conn.RemoteAddr())
		return "", err
	}
	conn.Write([]byte{correctPassword})
	log.Printf("correct password by: %s\n", conn.RemoteAddr())

	pageTitle, err = readPageTitle(conn)
	if err != nil {
		return "", err
	}
	return pageTitle, nil
}

func readPassword(conn net.Conn) ([]byte, error) {
	passwordSizeSlice := make([]byte, 4)
	_, err := io.ReadFull(conn, passwordSizeSlice)
	if err != nil {
		return nil, err
	}
	passwordSize := byteArrToInt(passwordSizeSlice)
	clientPassword := make([]byte, passwordSize)
	_, err = io.ReadFull(conn, clientPassword)
	if err != nil {
		return nil, err
	}
	return clientPassword, nil
}

func readPageTitle(conn net.Conn) (string, error) {
	titleSizeSlice := make([]byte, 4)
	_, err := io.ReadFull(conn, titleSizeSlice)
	if err != nil {
		return "", err
	}
	titleSize := byteArrToInt(titleSizeSlice)
	pageTitle := make([]byte, titleSize)
	_, err = io.ReadFull(conn, pageTitle)
	if err != nil {
		return "", err
	}
	return string(pageTitle), nil
}

func byteArrToInt(arr []byte) uint32 {
	return binary.BigEndian.Uint32(arr)
}
