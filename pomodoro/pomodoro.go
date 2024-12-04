package employee

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

type Pomodoro struct {
	Type     string
	Duration int
}

type Employee struct {
	FirstName   string
	LastName    string
	TotalLeaves int
	LeavesTaken int
}

func (e Employee) LeavesRemaining() {
	fmt.Printf("%s %s has %d leaves remaining\n", e.FirstName, e.LastName, (e.TotalLeaves - e.LeavesTaken))
}

func (p Pomodoro) TimeRemaining() {
	workDuration := 1 * time.Second
	breakDuration := 1 * time.Second
	numPomodoros := 2

	for i := 0; i < numPomodoros; i++ {
		fmt.Println("Starting work session...")
		playSound("assets/sounds/start_beep.wav")
		startTimer(workDuration)
		playSound("assets/sounds/stop_beep.wav")

		fmt.Println("Starting break session...")
		playSound("assets/sounds/start_beep.wav")
		startTimer(breakDuration)
		playSound("assets/sounds/stop_beep.wav")
	}

	fmt.Println("All done! Good job!")
}

func startTimer(duration time.Duration) {
	timer := time.NewTimer(duration)

	<-timer.C
	fmt.Println("Time's up!")
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
