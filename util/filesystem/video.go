package filesystem

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func GetVideoSize(videoPath string) (width int, height int, err error) {
	cmd := exec.Command("ffprobe", "-v", "quiet", "-select_streams", "v", "-show_entries", "stream=width,height", videoPath)
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}
	widthRe := regexp.MustCompile(`width=([0-9]+)`)
	heightRe := regexp.MustCompile(`height=([0-9]+)`)
	widthMatches := widthRe.FindSubmatch(output)
	heightMatches := heightRe.FindSubmatch(output)
	width, _ = strconv.Atoi(string(widthMatches[1]))
	height, _ = strconv.Atoi(string(heightMatches[1]))
	return width, height, nil
}

func DeleteVideoLogo(videoPath string, newVideo string, x int, y int, width int, height int) error {
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-b:v", "1000k", "-vf", fmt.Sprintf("delogo=x=%d:y=%d:w=%d:h=%d:show=0", x, y, width, height), newVideo)
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}

func SliceVideo(videoPath string, outputM3u8 string, timeSegment int) error {
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-c:v", "libx264", "-c:a", "copy", "-f", "hls", "-threads", "8", "-hls_time", strconv.Itoa(timeSegment), "-hls_list_size", "0", outputM3u8)
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}

func MergeVideoByM3u8(inputVideos []string, outputVideo string) error {
	dir := path.Dir(inputVideos[0])
	filename := dir + "/" + "input.txt"
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	for _, input := range inputVideos {
		file.Write([]byte("file " + path.Base(input) + "\n"))
	}
	cmd := exec.Command("ffmpeg","-f", "concat", "-safe", "0", "-i", filename, "-c", "copy", outputVideo)
	_, err = cmd.Output()
	if err != nil {
		return err
	}
	return nil
}

func MergeVideoByFileName(outputVideo string, inputVideos ...string) error {
	dir := path.Dir(inputVideos[0])
	filename := dir + "/" + outputVideo
	cmd := exec.Command("ffmpeg", "-i", fmt.Sprintf("concat:%s", strings.Join(inputVideos, "|")), "-c", "copy", filename)
	_, err :=  cmd.Output()
	if err != nil {
		return err
	}
	return nil
}

func GetDuration(filename string) (float64, error) {
	cmd := exec.Command("ffprobe", "-v", "quiet", "-select_streams", "v", "-show_entries", "stream=duration", filename)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return -1, fmt.Errorf(fmt.Sprint(err) + ": " + stderr.String())
	}

	re := regexp.MustCompile(`duration=([^\n]+)`)
	res := strings.TrimRight(strings.TrimLeft(re.FindString(out.String()), "duration="), "\r")

	return strconv.ParseFloat(res, 32)
}

func CreatePreviewVideo(inputVideos []string, outputVideo string, size string, removeAudio bool) error {
	dir := path.Dir(inputVideos[0])
	filename := dir + "/" + outputVideo
	var cmd *exec.Cmd
	if removeAudio {
		cmd = exec.Command("ffmpeg", "-i", fmt.Sprintf("concat:%s", strings.Join(inputVideos, "|")), "-vcodec", "copy", "-an", "-fs", size, filename)
	} else {
		cmd = exec.Command("ffmpeg", "-i", fmt.Sprintf("concat:%s", strings.Join(inputVideos, "|")), "-c", "copy", "-fs", size, filename)
	}
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}

func MergeEncryptionTsToVideo(m3u8File string, outputVideo string) error {
	// ffmpeg -allowed_extensions ALL -i index.m3u8 -c copy new.mp4
	dir := path.Dir(m3u8File)
	filename := dir + "/" + outputVideo
	var cmd *exec.Cmd
	cmd = exec.Command("ffmpeg", "-allowed_extensions", "ALL", "-i", m3u8File, "-c", "copy", filename)
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}

func CreatePreviewVideoByEncryptionTs(m3u8File string, outputVideo string, size string, removeAudio bool) error {
	dir := path.Dir(m3u8File)
	filename := dir + "/" + outputVideo
	var cmd *exec.Cmd
	if removeAudio {
		fmt.Println("ffmpeg", "-allowed_extensions", "ALL", "-i", m3u8File, "-vcodec", "copy", "-an", "-fs", size, filename)
		cmd = exec.Command("ffmpeg", "-allowed_extensions", "ALL", "-i", m3u8File, "-vcodec", "copy", "-an", "-fs", size, filename)
	} else {
		cmd = exec.Command("ffmpeg", "-allowed_extensions", "ALL", "-i", m3u8File, "-c", "copy", "-fs", size, filename)
	}
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}