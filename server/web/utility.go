package web

import "time"

func nonZeroDuration(in, def time.Duration) time.Duration {
	if in == 0 {
		return def * time.Second
	}

	return in * time.Second
}
