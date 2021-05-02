package cloudburst

type Threshold struct {
	Upper int
	Lower int
}

func (t *Threshold) inRange(i int) bool {
	return i >= t.Lower && i <= t.Upper
}

func newThreshold(upper, lower int) Threshold {
	return Threshold{
		Upper: upper,
		Lower: lower,
	}
}
