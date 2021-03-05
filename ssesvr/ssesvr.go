package ssesvr

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
)

type RespStream struct {
	Stream *types.HijackedResponse
}

func (resp *RespStream) SseServer(r *gin.Engine) {
	r.GET("/stream", func(c *gin.Context) {
		/*
		   we don't want the stream lasts forever, set the timeout
		*/
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		chanStream := make(chan string) // to consume lines read from docker
		done := make(chan bool)         // to indicate when the work is done
		/*
		   this is where we handle the request context
		*/
		go func() {
			for {
				select {
				case <-c.Request.Context().Done():
					// client gave up
					done <- true
					return
				case <-ctx.Done():
					// timeout
					switch ctx.Err() {
					case context.DeadlineExceeded:
						fmt.Println("timeout")
					}
					done <- true
					return
				}
			}
		}()
		/*
		   send log lines to channel
		*/
		rd := bufio.NewReader(resp.Stream.Reader)
		var mu sync.RWMutex
		go func() {
			for {
				mu.Lock()
				// read lines from the reader
				str, _, err := rd.ReadLine()
				if err != nil {
					log.Println("Read Error:", err.Error())
					done <- true
					return
				}
				// send the lines to channel
				chanStream <- string(str)
				mu.Unlock()
			}
		}()
		count := 0 // to indicate the message id
		isStreaming := c.Stream(func(w io.Writer) bool {
			for {
				select {
				case <-done:
					// when deadline is reached, send 'end' event
					c.SSEvent("end", "end")
					return false
				case msg := <-chanStream:
					// send events to client
					c.Render(-1, sse.Event{
						Id:    strconv.Itoa(count),
						Event: "message",
						Data:  msg,
					})
					count++
					return true
				}
			}
		})
		if !isStreaming {
			fmt.Println("stream closed")
		}
	})
}
