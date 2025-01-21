package IOSupport

type TimerQueue struct {
	CountdownTimer uint64
	Interrupt      uint64
}

var TimerService map[uint]TimerQueue

func Timer_Initialize() {
	TimerService = make(map[uint]TimerQueue)
}

func Timer_Terminate() {

}

func Timer_Add(timer uint, countdown uint64) {
}

func Timer_Remove(id uint) error {
	return nil
}
