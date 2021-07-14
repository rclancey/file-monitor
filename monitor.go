package monitor

import (
	"os"
	"time"
)

type FileChange struct {
	FileName string
	ModTime time.Time
}

type FileMonitor struct {
	FileNames []string
	Interval time.Duration
	IdleTime time.Duration
	C chan FileChange
	ticker *time.Ticker
	quit chan bool
}

func NewFileMonitor(interval, idle time.Duration, fns ...string) *FileMonitor {
	if interval <= 100 * time.Millisecond {
		interval = 100 * time.Millisecond
	}
	fm := &FileMonitor{
		FileNames: fns,
		Interval: interval,
		IdleTime: idle,
		C: make(chan FileChange, 10),
		ticker: nil,
	}
	fm.Start()
	return fm
}

func (fm *FileMonitor) Start() {
	fm.Stop()
	quit := make(chan bool, 2)
	fm.quit = quit
	lastMod := map[string]time.Time{}
	tick := time.NewTicker(fm.Interval)
	fm.ticker = tick
	go func() {
		for {
			select {
			case <-quit:
				return
			case <-tick.C:
				for _, fn := range fm.FileNames {
					st, err := os.Stat(fn)
					if err == nil {
						lm := st.ModTime()
						if time.Now().Sub(lm) >= fm.IdleTime {
							xlm, ok := lastMod[fn]
							if !ok {
								lastMod[fn] = lm
								fm.C <- FileChange{fn, lm}
							} else if lm.After(xlm) {
								lastMod[fn] = lm
								fm.C <- FileChange{fn, lm}
							}
						}
					} else {
						_, ok := lastMod[fn]
						if ok {
							delete(lastMod, fn)
							fm.C <- FileChange{fn, time.Now()}
						}
					}
				}
			}
		}
	}()
}

func (fm *FileMonitor) Stop() {
	if fm.quit != nil {
		fm.quit <- true
	}
	if fm.ticker != nil {
		fm.ticker.Stop()
		fm.ticker = nil
	}
}
