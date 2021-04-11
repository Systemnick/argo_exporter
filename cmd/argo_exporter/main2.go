package main

import (
	"errors"
	"time"
)

var ErrSome = errors.New("some error")
var ErrStop = errors.New("stop error")

type RoutineResult struct {
	WorkTime time.Duration
	Err      error
}

const (
	RoutineCountMin         = 50
	RoutineCountMax         = 80
	RoutineSleepMin         = 0
	RoutineSleepMax         = 3
	ErrProbabilityPercent   = 30
	CriticalErrorPercentage = 50
)

// Программа запустит рандомное количество горутин от 50 до 80 штук.
// Каждая горутина проспит рандомно от 0 до 3 секунд и с вероятностью 30% вернёт ошибку.
// Посчитать суммарное время, которое проработали горутины, прервать программу, если ошибок > 50%.
// Вывести статистику по ошибкам.
// func main() {
// 	rand.Seed(time.Now().Unix())
//
// 	routineCount := RoutineCountMin + rand.Intn(RoutineCountMax-RoutineCountMin+1)
//
// 	ch := make(chan RoutineResult, routineCount)
// 	done := make(chan struct{})
//
// 	for i := 0; i < routineCount; i++ {
// 		go SleepWithError(ch, done)
// 	}
//
// 	var (
// 		resultCount      int
// 		interruptedCount int
// 		errCount         int
// 		workTimeSum      time.Duration
// 		stopFlag         bool
// 	)
//
// 	for {
// 		result := <-ch
// 		resultCount++
// 		workTimeSum += result.WorkTime
//
// 		if errors.Is(result.Err, ErrSome) {
// 			errCount++
// 		}
//
// 		if errors.Is(result.Err, ErrStop) {
// 			interruptedCount++
// 		}
//
// 		// TODO: Уточнить постановку задачи:
// 		//  - считать процент ошибок от отработавших горутин, тогда нужно сохранять значение в переменную для вывода;
// 		//  - считать процент ошибок от общей массы горутин, тогда очень сложно добиться соблюдения условия остановки.
// 		if errCount*100/routineCount > CriticalErrorPercentage && !stopFlag {
// 			close(done)
// 			stopFlag = true
// 		}
//
// 		if resultCount >= routineCount {
// 			break
// 		}
// 	}
//
// 	fmt.Println("---=== Summary ===---")
// 	fmt.Printf("Goroutine spent time: %s\n", workTimeSum)
// 	fmt.Printf("Goroutine count: %d\n", routineCount)
// 	fmt.Printf("Goroutine interrupted count: %d\n", interruptedCount)
// 	fmt.Printf("Goroutine some error count: %d (~%d%%)\n", errCount, errCount*100/routineCount)
// }
//
// func SleepWithError(out chan RoutineResult, cancel chan struct{}) {
// 	var err error
//
// 	startTime := time.Now()
// 	defer func() {
// 		endTime := time.Now()
//
// 		workTime := endTime.Sub(startTime)
//
// 		out <- RoutineResult{
// 			WorkTime: workTime,
// 			Err:      err,
// 		}
// 	}()
//
// 	sleepSeconds := RoutineSleepMin + rand.Intn(RoutineSleepMax-RoutineSleepMin+1)
// 	sleepTime := time.Duration(sleepSeconds) * time.Second
// 	timer := time.NewTimer(sleepTime)
// 	success := timer.C
//
// 	select {
// 	case <-success:
// 	case <-cancel:
// 		timer.Stop()
// 		err = ErrStop
// 		return
// 	}
//
// 	errProbability := rand.Intn(100)
// 	if errProbability < ErrProbabilityPercent {
// 		err = ErrSome
// 	}
// }
//
// func PollRegistrar() {
//
// }
