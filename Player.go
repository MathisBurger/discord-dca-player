package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func Play(file string, vc *discordgo.VoiceConnection) error {
	if strings.HasSuffix(file, ".dca") {
		return streamToVC(file, vc)
	} else {
		ext := filepath.Ext(file)
		dcaFilePath := strings.ReplaceAll(file, ext, ".dca")
		if _, err := os.Stat(dcaFilePath); err == nil {
			return streamToVC(dcaFilePath, vc)
		}
		// ffmpeg -i test.mp3 -f s16le -ar 48000 -ac 2 pipe:1 | dca > test.dca
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd := exec.Command("ffmpeg", "-i", file, "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1", "|", "dca", ">", dcaFilePath)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			return err
		}
		return streamToVC(dcaFilePath, vc)
	}
}

func streamToVC(fileString string, vc *discordgo.VoiceConnection) error {
	file, err := os.Open(fileString)
	if err != nil {
		fmt.Println("Error opening dca file :", err)
		return err
	}

	var buffer [][]byte = make([][]byte, 0)
	var opuslen int16

	for {
		err = binary.Read(file, binary.LittleEndian, &opuslen)
		if err == io.EOF || errors.Is(err, io.ErrUnexpectedEOF) {
			err := file.Close()
			if err != nil {
				return err
			}
			break
		}
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}
		buffer = append(buffer, InBuf)
	}
	time.Sleep(250 * time.Millisecond)
	err = vc.Speaking(true)
	if err != nil {
		return err
	}
	for _, buff := range buffer {
		vc.OpusSend <- buff
	}
	err = vc.Speaking(false)
	if err != nil {
		return err
	}
	time.Sleep(250 * time.Millisecond)
	err = vc.Disconnect()
	if err != nil {
		return err
	}

	return nil
}
