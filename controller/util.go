package controller

import "time"

func retry(inv time.Duration, fn func() error) {
	for {
		err := fn()
		if err != nil {
			time.Sleep(inv)
			continue
		}
		break
	}
}
