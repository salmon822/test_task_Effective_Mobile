package time_manager

import "time"

type TimeManager interface {
	Now() time.Time
	MillisecondsToTime(milliseconds int64) time.Time
	Locale() *time.Location
	LocaleOffsetMilli() int64
}

type timeManager struct {
	localeOffsetMilli int64
	locale            *time.Location
}

func New(locale int64) TimeManager {
	return &timeManager{
		localeOffsetMilli: locale * time.Hour.Milliseconds(),
		locale:            time.FixedZone("MSC", int(locale)*3600),
	}
}

func (t *timeManager) Locale() *time.Location {
	return t.locale
}

func (t *timeManager) LocaleOffsetMilli() int64 {
	return t.localeOffsetMilli
}

func (t timeManager) Now() time.Time {
	return time.Now().In(t.locale)
}

func (t *timeManager) MillisecondsToTime(milliseconds int64) time.Time {
	return time.UnixMilli(milliseconds).In(t.locale)
}
