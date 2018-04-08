// Copyright 2018 rootkiwi
//
// screen_share_remote_go is licensed under GNU General Public License 3 or later.
//
// See LICENSE for more details.

package webserver

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/websocket"
)

var (
	loadFilesOnce sync.Once
	js            packr.Box
	css           packr.Box
	indexTmpl     *template.Template
)

var server *http.Server

type FrameQueue chan<- []byte

func Start(port int, pageTitle string, entering, leaving chan<- FrameQueue) {
	loadFilesOnce.Do(loadStaticFiles)
	mux := http.NewServeMux()
	mux.Handle("/", serveIndex(pageTitle))
	mux.Handle("/js/", serveStaticJs())
	mux.Handle("/css/", serveStaticCSS())
	mux.Handle("/ws/", serveWs(entering, leaving))
	server = &http.Server{Addr: ":" + strconv.Itoa(port), Handler: mux}
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatalf("failed starting web server: %v\n", err)
	}
}

func loadStaticFiles() {
	js = packr.NewBox("../../../web/static/js")
	css = packr.NewBox("../../../web/static/css")
	index := packr.NewBox("../../../web/dynamic").String("index.html")
	indexTmpl = template.Must(template.New("index").Parse(index))
}

func Stop() {
	server.Close()
}

func serveIndex(pageTitle string) http.Handler {
	buf := new(bytes.Buffer)
	indexTmpl.Execute(buf, pageTitle)
	indexHTML := buf.String()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "", "/", "index.html":
			w.Write([]byte(indexHTML))
		default:
			http.NotFound(w, r)
		}
	})
}

func serveStaticJs() http.Handler {
	return http.StripPrefix("/js/", noDirListing(http.FileServer(js)))
}

func serveStaticCSS() http.Handler {
	return http.StripPrefix("/css/", noDirListing(http.FileServer(css)))
}

func noDirListing(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}

var upgrader = websocket.Upgrader{
	WriteBufferSize: 1024 * 8,
}

func serveWs(entering, leaving chan<- FrameQueue) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("fail upgrade websocket: %v\n", err)
			return
		}
		defer conn.Close()
		frameQueue := make(chan []byte, 80)
		entering <- frameQueue
		defer func() { leaving <- frameQueue }()
		go readLoop(conn)
		for frame := range frameQueue {
			if err := conn.WriteMessage(websocket.BinaryMessage, frame); err != nil {
				return
			}
		}
	})
}

func readLoop(c *websocket.Conn) {
	for {
		if _, _, err := c.NextReader(); err != nil {
			c.Close()
			break
		}
	}
}
