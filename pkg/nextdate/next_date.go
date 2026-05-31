// Package nextdate
package nextdate

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/antonlearn/go-final-project/internal/db"
	"github.com/antonlearn/go-final-project/pkg/format"
)

var errEmpty = errors.New("line with repetition rule: empty line")

func CheckDate(task *db.Task) error {

	nowDateStr := time.Now().Format(format.DateFormatTemplateYYYYMMDD)
	if task.Date == "" {
		task.Date = nowDateStr
	}
	taskDate, err := time.Parse(format.DateFormatTemplateYYYYMMDD, task.Date)
	if err != nil {
		return err
	}
	nowDate, err := time.Parse(format.DateFormatTemplateYYYYMMDD, nowDateStr)
	if err != nil {
		return err
	}
	var nextDateStr string
	nextDateStr, err = NextDate(nowDate, task.Date, task.Repeat)
	if firstDateMoreSecondDate(nowDate, taskDate) {
		if errors.Is(err, errEmpty) {
			task.Date = nowDateStr
			return nil
		} else {
			if err != nil {
				return err
			}
			task.Date = nextDateStr
			return nil
		}
	}
	return nil
}

func firstDateMoreSecondDate(firstDate, secondDate time.Time) bool {
	return firstDate.After(secondDate)
}

func daysInMonth(date time.Time) int {
	return 32 - time.Date(date.Year(), date.Month(), 32, 0, 0, 0, 0, time.UTC).Day()
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	var (
		errTooManyParameters = fmt.Errorf("line with repetition rule %s: too many parameters", repeat)
		errTooFewParameters  = fmt.Errorf("line with repetition rule %s: too few parameters", repeat)
		errUnsupportedFormat = fmt.Errorf("line with repetition rule %s: unsupported format", repeat)
	)
	if repeat == "" {
		return "", errEmpty
	}
	resultDate, err := time.Parse(format.DateFormatTemplateYYYYMMDD, dstart)
	if err != nil {
		return "", err
	}
	repeats := strings.Split(repeat, " ")
	var errRuleParameter error
	switch len(repeats) {
	case 2:
		errRuleParameter = fmt.Errorf("line with repetition rule %s and parameter %s", repeats[0], repeats[1])
	case 3:
		errRuleParameter = fmt.Errorf("line with repetition rule %s and parameter %s", repeats[0], repeats[2])
	}
	switch repeats[0] {
	case "y":
		if len(repeats) > 1 {
			return "", errTooManyParameters
		}
		for {
			resultDate = resultDate.AddDate(1, 0, 0)
			if firstDateMoreSecondDate(resultDate, now) {
				break
			}
		}
		return resultDate.Format(format.DateFormatTemplateYYYYMMDD), nil
	case "d":
		const maxD = 400
		if len(repeats) < 2 {
			return "", errTooFewParameters
		}
		if len(repeats) > 2 {
			return "", errTooManyParameters
		}
		num, err := strconv.Atoi(repeats[1])
		if err != nil {
			return "", fmt.Errorf("%s: %w", errRuleParameter, err)
		}
		if num > maxD {
			return "", fmt.Errorf("%s: %w", errRuleParameter, fmt.Errorf("maximum %d allowed interval in days has been exceeded", maxD))
		}
		if num < 1 {
			return "", fmt.Errorf("%s: %w", errRuleParameter, errors.New("invalid interval in days"))
		}
		for {
			resultDate = resultDate.AddDate(0, 0, num)
			if firstDateMoreSecondDate(resultDate, now) {
				break
			}
		}
		return resultDate.Format(format.DateFormatTemplateYYYYMMDD), nil
	case "w":
		if firstDateMoreSecondDate(now, resultDate) {
			resultDate = now
		}
		if len(repeats) < 2 {
			return "", errTooFewParameters
		}
		if len(repeats) > 2 {
			return "", errTooManyParameters
		}
		weekdays := map[int]time.Weekday{
			1: time.Monday, 2: time.Tuesday, 3: time.Wednesday, 4: time.Thursday,
			5: time.Friday, 6: time.Saturday, 7: time.Sunday,
		}
		var weekdaysName []time.Weekday
		const daysInWeek = 7
		for weekday := range strings.SplitSeq(repeats[1], ",") {
			weekdayNumber, err := strconv.Atoi(weekday)
			if err != nil {
				return "", fmt.Errorf("%s: %w", errRuleParameter, err)
			}
			if weekdayNumber > daysInWeek {
				return "", fmt.Errorf("%s: %w", errRuleParameter, fmt.Errorf("maximum %d number of days in a week has been exceeded", daysInWeek))
			}
			if weekdayNumber < 1 {
				return "", fmt.Errorf("%s: %w", errRuleParameter, errors.New("invalid day of the week ordinal number"))
			}
			weekdaysName = append(weekdaysName, weekdays[weekdayNumber])
		}
		for {
			resultDate = resultDate.AddDate(0, 0, 1)
			if slices.Contains(weekdaysName, resultDate.Weekday()) {
				return resultDate.Format(format.DateFormatTemplateYYYYMMDD), nil
			}
		}
	case "m":
		if firstDateMoreSecondDate(now, resultDate) {
			resultDate = now
		}
		if len(repeats) < 2 {
			return "", errTooFewParameters
		}
		if len(repeats) > 3 {
			return "", errTooManyParameters
		}
		const (
			minMonthday, invalidMonthday, maxMonthday = -2, 0, 31
			minMonth, maxMonth                        = 1, 12
		)
		var (
			monthdays []int
			months    []time.Month
		)
		for monthday := range strings.SplitSeq(repeats[1], ",") {
			monthdayNumber, err := strconv.Atoi(monthday)
			if err != nil {
				return "", fmt.Errorf("%s: %w", errRuleParameter, err)
			}
			if monthdayNumber < minMonthday || monthdayNumber == invalidMonthday || monthdayNumber > maxMonthday {
				return "", fmt.Errorf("%s: %w", errRuleParameter, errors.New("invalid day of the month ordinal number"))
			}
			monthdays = append(monthdays, monthdayNumber)
		}
		var monthRepeatsNotApply bool
		if len(repeats) > 2 {
			monthRepeatsNotApply = true
			for month := range strings.SplitSeq(repeats[2], ",") {
				monthNumber, err := strconv.Atoi(month)
				if err != nil {
					return "", fmt.Errorf("%s: %w", errRuleParameter, err)
				}
				if monthNumber < minMonth || monthNumber > maxMonth {
					return "", fmt.Errorf("%s: %w", errRuleParameter, errors.New("invalid month ordinal number"))
				}
				months = append(months, time.Month(monthNumber))
			}
		}
		for {
			resultDate = resultDate.AddDate(0, 0, 1)
			dOfM := daysInMonth(resultDate)
			var isPenultimateDay, isLastDay bool
			const penultimateDayOfMonth, lastDayOfMonth = -2, -1
			if resultDate.Day() == dOfM {
				isLastDay = true
			}
			if resultDate.Day() == dOfM-1 {
				isPenultimateDay = true
			}
			if slices.Contains(monthdays, penultimateDayOfMonth) && isPenultimateDay {
				return resultDate.Format(format.DateFormatTemplateYYYYMMDD), nil
			}
			if slices.Contains(monthdays, lastDayOfMonth) && isLastDay {
				return resultDate.Format(format.DateFormatTemplateYYYYMMDD), nil
			}
			if slices.Contains(monthdays, resultDate.Day()) && !monthRepeatsNotApply {
				return resultDate.Format(format.DateFormatTemplateYYYYMMDD), nil
			}
			if slices.Contains(monthdays, resultDate.Day()) && slices.Contains(months, resultDate.Month()) {
				return resultDate.Format(format.DateFormatTemplateYYYYMMDD), nil
			}
		}
	default:
		return "", errUnsupportedFormat
	}
}
