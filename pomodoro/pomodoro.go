package pomodoro

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

type pomodoro struct {
	Type       string
	Duration   time.Duration
	StartSound string
	StopSound  string
}

func New(pomoType string, pomoDuration int, startSoundPath string, stopSoundPath string) pomodoro {

	duration := time.Duration(pomoDuration) * time.Second
	p := pomodoro{pomoType, duration, startSoundPath, stopSoundPath}

	return p
}

func (p pomodoro) Repeat(repetitions int) {
	numPomodoros := repetitions
	breakPomo := New("break", 1, "assets/sounds/start_beep.wav", "assets/sounds/stop_beep.wav")

	for i := 0; i < numPomodoros; i++ {
		executePomodoro(p)
		executePomodoro(breakPomo)
	}

	fmt.Println("All done! Good job!")
}

func executePomodoro(p pomodoro) {
	playSound(p.StartSound)
	fmt.Printf("Starting %s session...\n", p.Type)
	startTimer(time.Duration(p.Duration), p.Type)
	playSound(p.StopSound)
}

func startTimer(duration time.Duration, label string) {
	timer := time.NewTimer(duration)

	<-timer.C
	fmt.Printf("Ending %s session: Time's up!\n", label)
}

func playSound(soundFile string) {
	f, err := os.Open(soundFile)
	if err != nil {
		panic(err)
	}

	s, format, err := wav.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		panic(err)
	}

	done := make(chan bool)
	speaker.Play(beep.Seq(s, beep.Callback(func() {
		done <- true
	})))
	<-done
}
