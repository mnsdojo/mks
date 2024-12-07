package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

type SoundEffect struct {
	Name string
	File string
}

var (
	sounds = []SoundEffect{
		{"Sound 1", "./output.wav"},
		{"Sound 2", "./output.wav"},
		{"Sound 3", "./output.wav"},
	}
	selectedSound *SoundEffect
	playMutex     sync.Mutex
	isPlaying     bool
)

func main() {
	// Initialize speaker once
	err := initializeSpeaker()
	if err != nil {
		log.Fatalf("Failed to initialize speaker: %v", err)
	}
	defer speaker.Close()

	// Validate sound files exist
	for i := range sounds {
		if !fileExists(sounds[i].File) {
			log.Printf("Warning: Sound file %s does not exist", sounds[i].File)
		}
	}

	selectSound()

	if err := keyboard.Open(); err != nil {
		log.Fatalf("Failed to open keyboard: %v", err)
	}
	defer keyboard.Close()

	fmt.Println("MKS")
	fmt.Println("Commands:")
	fmt.Println("- Press any key to play current sound")
	fmt.Println("- Press 'n' to change sound")
	fmt.Println("- Press 'ESC' to exit")

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			log.Printf("Keyboard error: %v", err)
			continue
		}

		switch {
		case key == keyboard.KeyEsc:
			fmt.Println("Exiting program")
			return

		case char == 'n':
			// Change sound
			selectSound()

		case key == keyboard.KeySpace || (char != 0 && char != 'n'):
			// Play sound on any key press
			go playCurrentSound()
		}
	}
}

func selectSound() {
	fmt.Println("\nAvailable Sound Effects:")
	for i, sound := range sounds {
		fmt.Printf("%d. %s\n", i+1, sound.Name)
	}

	var input string
	fmt.Print("Select sound effect (1-3): ")
	fmt.Scanln(&input)

	index, err := strconv.Atoi(input)
	if err != nil || index < 1 || index > len(sounds) {
		fmt.Println("Invalid selection. Defaulting to sound 1.")
		selectedSound = &sounds[0]
	} else {
		selectedSound = &sounds[index-1]
	}

	fmt.Printf("Selected: %s\n", selectedSound.Name)
}

func playCurrentSound() {
	playMutex.Lock()
	if isPlaying {
		playMutex.Unlock()
		fmt.Println("Sound is already playing")
		return
	}
	isPlaying = true

	playMutex.Unlock()
	if selectedSound == nil {
		fmt.Println("No sound selected")
		playMutex.Lock()
		isPlaying = false
		playMutex.Unlock()
		return
	}
	err := playSound(selectedSound.File)
	if err != nil {
		fmt.Printf("failed to play sound :%v\n", err)

	}
	playMutex.Lock()
	isPlaying = false
	playMutex.Unlock()
}

func initializeSpeaker() error {
	// Determine a reasonable sample rate (44.1 kHz is standard)
	sampleRate := beep.SampleRate(44100)
	return speaker.Init(sampleRate, sampleRate.N(time.Second/10))
}

func playSound(filePath string) error {
	// Open sound file
	soundFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open sound file: %w", err)
	}
	defer soundFile.Close()

	// Decode WAV file
	streamer, _, err := wav.Decode(soundFile)
	if err != nil {
		return fmt.Errorf("failed to decode sound file: %w", err)
	}
	defer streamer.Close()

	// Create a channel to wait for sound to finish
	done := make(chan bool)

	// Play the sound
	speaker.Play(beep.Seq(
		streamer,
		beep.Callback(func() {
			done <- true
		}),
	))

	// Wait for playback to complete
	<-done

	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
