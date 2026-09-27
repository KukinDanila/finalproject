package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func afterNow(date, now time.Time) bool {
	dateStr := date.Format("20060102")
	nowStr := now.Format("20060102")
	return dateStr > nowStr
}

func toTaskWeekday(t time.Time) int {
	wd := int(t.Weekday())
	if wd == 0 {
		return 7
	}
	return wd
}

func lastDayOfMonth(t time.Time) time.Time {
	firstOfNextMonth := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
	return firstOfNextMonth.AddDate(0, 0, -1)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", err
	}

	if repeat == "" {
		return "", errors.New("не указано правило повторения")
	}

	parts := strings.Split(repeat, " ")

	if parts[0] == "d" {
		if len(parts) < 2 {
			return "", errors.New("не указан интервал в днях")
		}

		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", err
		}

		if days > 400 {
			return "", errors.New("превышен максимальный интервал в днях")
		}

		for {
			date = date.AddDate(0, 0, days)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format("20060102"), nil
	}

	if parts[0] == "y" {
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}

		return date.Format("20060102"), nil
	}

	if parts[0] == "w" {
		if len(parts) < 2 {
			return "", errors.New("не указаны дни недели")
		}

		weekdayStrs := strings.Split(parts[1], ",")
		var weekdays []int

		for _, s := range weekdayStrs {
			wd, err := strconv.Atoi(s)
			if err != nil {
				return "", err
			}
			if wd < 1 || wd > 7 {
				return "", errors.New("недопустимый день недели")
			}
			weekdays = append(weekdays, wd)
		}

		for {
			date = date.AddDate(0, 0, 1)

			if afterNow(date, now) {
				currentWeekday := toTaskWeekday(date)

				found := false
				for _, wd := range weekdays {
					if wd == currentWeekday {
						found = true
						break
					}
				}

				if found {
					break
				}
			}
		}

		return date.Format("20060102"), nil
	}

	if parts[0] == "m" {
		if len(parts) < 2 {
			return "", errors.New("не указаны дни месяца")
		}

		dayStrs := strings.Split(parts[1], ",")
		days := make(map[int]bool)

		for _, s := range dayStrs {
			d, err := strconv.Atoi(s)
			if err != nil {
				return "", err
			}
			if d == 0 || d < -2 || d > 31 {
				return "", errors.New("недопустимый день месяца")
			}
			days[d] = true
		}

		months := make(map[int]bool)
		if len(parts) >= 3 {
			monthStrs := strings.Split(parts[2], ",")
			for _, s := range monthStrs {
				m, err := strconv.Atoi(s)
				if err != nil {
					return "", err
				}
				if m < 1 || m > 12 {
					return "", errors.New("недопустимый месяц")
				}
				months[m] = true
			}
		}

		for {
			date = date.AddDate(0, 0, 1)

			if !afterNow(date, now) {
				continue
			}

			currentDay := date.Day()
			last := lastDayOfMonth(date)
			isLastDay := currentDay == last.Day()
			isSecondLastDay := currentDay == last.AddDate(0, 0, -1).Day()

			dayMatches := days[currentDay] ||
				(isLastDay && days[-1]) ||
				(isSecondLastDay && days[-2])

			if !dayMatches {
				continue
			}

			if len(months) > 0 && !months[int(date.Month())] {
				continue
			}

			break
		}

		return date.Format("20060102"), nil
	}

	return "", errors.New("неподдерживаемый формат правила повторения")
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowParam := r.FormValue("now")
	dateParam := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	if nowParam == "" {
		now = time.Now()
	} else {
		parsedNow, err := time.Parse("20060102", nowParam)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		now = parsedNow
	}

	next, err := NextDate(now, dateParam, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(next))
}
