package main

import (
	"fmt"
	"os"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

func main() {
	fmt.Println("Press any key to play sound: Press Esc to QUIT")

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer keyboard.Close()

	soundFile, err := os.Open("./output.wav")
	if err != nil {
		panic(fmt.Sprintf("Failed to open sound file: %v", err))
	}
	defer soundFile.Close()

	_, format, err := wav.Decode(soundFile)
	if err != nil {
		panic(fmt.Sprintf("Failed to decode sound file: %v", err))
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		fmt.Printf("Key pressed: %q (Special Key: %v)\n", char, key)

		if key == keyboard.KeyEsc {
			break
		}

		soundFile, err := os.Open("./output.wav")
		if err != nil {
			fmt.Printf("Failed to open sound file: %v\n", err)
			continue
		}

		streamer, _, err := wav.Decode(soundFile)
		if err != nil {
			fmt.Printf("Failed to decode sound file: %v\n", err)
			soundFile.Close()
			continue
		}

		speaker.Play(beep.Seq(streamer, beep.Callback(func() {
			fmt.Println("Sound playback finished")
			streamer.Close()
			soundFile.Close()
		})))
	}
}
