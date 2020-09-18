package service

type SecLimit struct {
	count int
	curTime int64
}

func (this *SecLimit) Count(nowTime int64) {
	if this.curTime != nowTime {
		this.count = 1
		this.curTime = nowTime
	} else {
		this.count++
	}
}
func (this *SecLimit) Check(nowTime int64) int {
	if this.curTime != nowTime {
		return 1
	}
	return this.count
}
