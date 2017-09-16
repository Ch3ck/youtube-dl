// Copyright 2017 YTD Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// apiconv: Converts Decoded Video data to MP3, WEBM or MP4.
// NOTE: To reimplement using Go ffmpeg bindings.
package api

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//Converts Decoded Video file to mp3 by default with 123 bitrate or to
//flv if otherwise specified and downloads to system
func APIConvertVideo(file string, int bitrate, id string, decVideo []byte)  error {
	cmd := exec.Command("ffmpeg", "-i", "-", "-ab", fmt.Sprintf("%dk", bitrate), path)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		logrus.Fatalf(err)
	}
	if filepath.Ext(file) != ".mp3" && filepath.Ext(file) != ".flv" {
		file = file[:len(file)-4] + ".mp3"
	}

	logrus.Infof("Converting video to %q format", filepath.Ext(file))
	if filepath.Ext(file) == ".mp3" {
		// NOTE: To modify to use Go ffmpeg bindings or cgo
		_, err = exec.LookPath("ffmpeg")
		if err != nil {
			logrus.Errorf("ffmpeg not found on system")
		}

		cmd.Start()
		logrus.Infof("Downloading mp3 file to disk %s", path)
		cmd.Write(decVideo) //download file.

	} else {
		cmd, err = os.Create(path)
		if err != nil {
			logrus.Error("Unable to download video file.", err)
		}
		err = apiDownloadVideo(id, cmd)
		return err
	}

	return nil
}

//Downloads decoded video stream.
func APIDownloadVideo(videoUrl, cmd io.Writer) error {
	logrus.Infof("Downloading file stream")

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("requesting stream: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("reading answer: non 200 status code received: '%s'", err)
	}
	length, err := io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("saving file: %s (%d bytes copied)", err, length)
	}

	logrus.Infof("Downloaded %d bytes", length)

	return nil
}
