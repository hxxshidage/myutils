package utime

import (
	"fmt"
	"strconv"
	"time"
)

const (
	ProgramFmt = "2006-01-02 15:04:05.000"

	DefaultFmt = "2006-01-02 15:04:05"

	YMDFmt = "2006-01-02"

	YMNoSepFmt = "200601"

	YMDNoSepFmt = "20060102"

	TimeFmt = "15:04:05"
)

var cnLoc, _ = time.LoadLocation("Asia/Shanghai")

func UtcNow() time.Time {
	return time.Now().UTC()
}

func TsSec() int64 {
	return time.Now().Unix()
}

func TsMill() int64 {
	return time.Now().UnixMilli()
}

func FromTs(millTs int64) time.Time {
	return time.UnixMilli(millTs)
}

func FromTsAsCn(millTs int64) time.Time {
	return time.UnixMilli(millTs).In(cnLoc)
}

func ToTs(t time.Time) int64 {
	return t.UnixMilli()
}

func Now() time.Time {
	return time.Now()
}

func CnNow() time.Time {
	return time.Now().In(cnLoc)
}

func Ts() int64 {
	return time.Now().Unix()
}

func Utc2Cn(t time.Time) time.Time {
	return t.In(cnLoc)
}

func Fmt(t time.Time, patten string) string {
	return t.Format(patten)
}

func FmtProgram() string {
	return FmtWith(ProgramFmt)
}

func FmtDefault() string {
	return FmtWith(DefaultFmt)
}

func FmtWith(patten string) string {
	return Now().Format(patten)
}

func FmtVal(t time.Time) string {
	return Fmt(t, DefaultFmt)
}

func FmtAsYmdNum(t *time.Time) uint32 {
	if t == nil || t.IsZero() {
		return 0
	}

	ymdNum, _ := strconv.ParseUint(fmt.Sprintf("%d%d%d", t.Year(), t.Month(), t.Day()), 10, 32)

	return uint32(ymdNum)
}

func FmtYmdWithVal(t time.Time) string {
	return Fmt(t, YMDFmt)
}

func FmtYmd(t time.Time) string {
	return FmtWith(YMDFmt)
}

func PastDays(days int) time.Time {
	return PastDaysWithVal(UtcNow(), days)
}

func PastDaysWithVal(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, -days)
}

func YesterdayWithVal(t time.Time) time.Time {
	return PastDaysWithVal(t, 1)
}

func Yesterday() time.Time {
	return PastDays(1)
}

func MinMaxTimeInDayWithVal(t time.Time) []time.Time {
	mi := t.Truncate(time.Hour * 24)
	ma := mi.Add(time.Hour*24 - time.Nanosecond)
	return []time.Time{
		mi, ma,
	}
}

func MinMaxTimeInDay() []time.Time {
	mi := UtcNow().Truncate(time.Hour * 24)
	ma := mi.Add(time.Hour*24 - time.Nanosecond)
	return []time.Time{
		mi, ma,
	}
}

func MinTimeInDay() time.Time {
	return UtcNow().Truncate(time.Hour * 24)
}

func MinTimeInDayWithVal(t time.Time) time.Time {
	return t.Truncate(time.Hour * 24)
}
