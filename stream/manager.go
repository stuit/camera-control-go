package stream

import (
	"bufio"
	"camera-control/mq"
	"fmt"
	"os/exec"
	"sync"
)

type StreamJob struct {
	StopChan chan struct{}
}

var jobs = make(map[string]*StreamJob)
var mu sync.Mutex

func StartStream(cameraID, rtspURL string, everyN int) {
	mu.Lock()
	if _, ok := jobs[cameraID]; ok {
		mu.Unlock()
		fmt.Println("Stream already running")
		return
	}
	job := &StreamJob{StopChan: make(chan struct{})}
	jobs[cameraID] = job
	mu.Unlock()

	go func() {
		cmd := exec.Command("ffmpeg", "-i", rtspURL, "-f", "image2pipe", "-vcodec", "mjpeg", "-")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Println("Failed to get stdout:", err)
			return
		}

		if err := cmd.Start(); err != nil {
			fmt.Println("Failed to start FFmpeg:", err)
			return
		}

		reader := bufio.NewReader(stdout)
		done := make(chan error, 1)

		go func() {
			done <- cmd.Wait() // waits for FFmpeg to exit
		}()

		frameCount := 0
		for {
			select {
			case <-job.StopChan:
				_ = cmd.Process.Kill()
				<-done // wait for cleanup
				return

			case err := <-done:
				fmt.Println("FFmpeg exited:", err)
				return

			default:
				frame, err := reader.ReadBytes(0xD9) // end of JPEG
				if err != nil {
					fmt.Println("Read error:", err)
					return
				}

				frameCount++
				if frameCount%everyN == 0 {
					mq.SendToRabbit(frame, cameraID)
				}
			}
		}
	}()

}

func StopStream(cameraID string) {
	mu.Lock()
	if job, ok := jobs[cameraID]; ok {
		close(job.StopChan)
	}
	mu.Unlock()
}
