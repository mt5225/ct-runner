package runner

import (
	"bufio"
	"context"
	"io"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
)

// Streaming conttainer output
func (c *Command) Streaming(r *gin.Engine) {
	r.GET("/stream/"+c.ContainerInstance.ID, func(r *gin.Context) {
		/*
		  inspect output in container
		*/
		docker, err := client.NewEnvClient()

		if err != nil {
			panic(err)
		}

		resp, err := docker.ContainerExecAttach(
			c.ContainerInstance.Context,
			c.ContainerInstance.RunID,
			types.ExecStartCheck{})

		if err != nil {
			panic(err)
		}

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
				case <-r.Request.Context().Done():
					// client gave up
					done <- true
					return
				case <-ctx.Done():
					// timeout
					switch ctx.Err() {
					case context.DeadlineExceeded:
						log.Printf("timeout")
					}
					done <- true
					return
				}
			}
		}()

		/*
		   send log lines to channel
		*/
		rd := bufio.NewReader(resp.Reader)
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
		isStreaming := r.Stream(func(w io.Writer) bool {
			for {
				select {
				case <-done:
					// when deadline is reached, send 'end' event
					r.SSEvent("end", "end")
					return false
				case msg := <-chanStream:
					// send events to client
					r.Render(-1, sse.Event{
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
			log.Printf("stream closed, kill container %s", c.ContainerInstance.ID)
			c.stopAndRemoveContainer(docker)
		}
	})
}
