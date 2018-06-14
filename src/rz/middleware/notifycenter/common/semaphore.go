package common
//
//func NewSemaphore() (*Semaphore) {
//	semaphore := &Semaphore{}
//	go semaphore.release()
//
//	return semaphore
//}
//
//type Semaphore struct {
//	flag        int
//	blocked     bool
//	blockChan   chan bool
//	releaseFunc func()
//}
//
//func (myself *Semaphore) WaitOne() {
//	myself.flag++
//
//	if 0 < myself.flag {
//		<-myself.blockChan
//	}
//
//}
//
//func (myself *Semaphore) Release() {
//	myself.flag--
//
//	if myself.blocked {
//
//	}
//	go func() {
//		myself.blockChan <- true
//
//	}()
//}
//
//func (myself *Semaphore) release() {
//		myself.blockChan <- true
//
//}
