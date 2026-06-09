// Package nextdate provides an internal scheduling engine to parse custom recurrence rules
// and calculate subsequent task execution milestones across daily, weekly, monthly, and yearly intervals.
package nextdate

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/antonlearn/go-final-project/internal/model"
	"github.com/antonlearn/go-final-project/pkg/format"
	"github.com/antonlearn/go-final-project/pkg/logger"
)

var errEmpty = errors.New("line with repetition rule: empty line")

// CheckDate evaluates a task's target execution date. If the date is omitted, it defaults to the current day.
// If the target execution date falls in the past, the schedule is advanced to its next logical occurrence
// based on the task's designated repetition rule string.
func CheckDate(task *model.Task, appLogger *logger.Logger) error {
	nowDateStr := time.Now().Format(format.YYYYMMDD)

	if task.Date == "" {
		task.Date = nowDateStr
		appLogger.Infof("Task date was empty, set to current date: %s", nowDateStr)
	}

	taskDate, err := time.Parse(format.YYYYMMDD, task.Date)
	if err != nil {
		appLogger.Errorf("Failed to parse task date %s: %v", task.Date, err)
		return err
	}

	nowDate, err := time.Parse(format.YYYYMMDD, nowDateStr)
	if err != nil {
		appLogger.Errorf("Failed to parse current date: %v", err)
		return err
	}

	nextDateStr, err := NextDate(nowDate, task.Date, task.Repeat, appLogger)

	if firstDateMoreSecondDate(nowDate, taskDate) {
		if errors.Is(err, errEmpty) {
			task.Date = nowDateStr
			appLogger.Info("Empty repeat rule detected, using current date")
			return nil
		} else {
			if err != nil {
				return err
			}
			task.Date = nextDateStr
			appLogger.Infof("Task date updated to next occurrence: %s", nextDateStr)
			return nil
		}
	}

	appLogger.Infof("Task date %s is in the future, no update needed", task.Date)
	return nil
}

// firstDateMoreSecondDate returns true if the first timestamp chronologically succeeds the second timestamp.
func firstDateMoreSecondDate(firstDate, secondDate time.Time) bool {
	return firstDate.After(secondDate)
}

// daysInMonth calculates the precise total number of calendar days contained within the target date's month.
func daysInMonth(date time.Time) int {
	return 32 - time.Date(date.Year(), date.Month(), 32, 0, 0, 0, 0, time.UTC).Day()
}

// NextDate computes the subsequent string-formatted milestone execution date matching a repetition sequence
// based on a baseline window timestamp, starting anchor date, and evaluation rules.
func NextDate(now time.Time, dstart string, repeat string, appLogger *logger.Logger) (string, error) {
	if repeat == "" {
		appLogger.Info("Empty repeat rule provided")
		return "", errEmpty
	}

	resultDate, err := time.Parse(format.YYYYMMDD, dstart)
	if err != nil {
		appLogger.Errorf("Failed to parse start date %s: %v", dstart, err)
		return "", err
	}

	repeats := strings.Split(repeat, " ")
	appLogger.Infof("Processing repeat rule: %s", repeat)

	errTooManyParameters := fmt.Errorf("line with repetition rule %s: too many parameters", repeat)
	errTooFewParameters := fmt.Errorf("line with repetition rule %s: too few parameters", repeat)
	errUnsupportedFormat := fmt.Errorf("line with repetition rule %s: unsupported format", repeat)

	var errRuleParameter error
	if len(repeats) == 2 {
		errRuleParameter = fmt.Errorf("line with repetition rule %s and parameter %s", repeats[0], repeats[1])
	} else if len(repeats) == 3 {
		errRuleParameter = fmt.Errorf("line with repetition rule %s and parameter %s", repeats[0], repeats[2])
	}

	switch repeats[0] {
	case "y":
		if len(repeats) > 1 {
			return "", errTooManyParameters
		}
		appLogger.Info("Applying yearly repetition rule")

		if resultDate.After(now) {
			return resultDate.Format(format.YYYYMMDD), nil
		}
		for {
			resultDate = resultDate.AddDate(1, 0, 0)
			if resultDate.After(now) {
				break
			}
		}
		return resultDate.Format(format.YYYYMMDD), nil

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
			appLogger.Errorf("Invalid number format in 'd' rule: %s", repeats[1])
			return "", fmt.Errorf("%s: %w", errRuleParameter, err)
		}
		if num > maxD {
			return "", fmt.Errorf("%s: %w", errRuleParameter, fmt.Errorf("maximum %d allowed interval in days has been exceeded", maxD))
		}
		if num < 1 {
			return "", fmt.Errorf("%s: %w", errRuleParameter, errors.New("invalid interval in days"))
		}

		appLogger.Infof("Applying daily repetition rule with interval: %d days", num)

		if resultDate.After(now) {
			return resultDate.Format(format.YYYYMMDD), nil
		}
		for {
			resultDate = resultDate.AddDate(0, 0, num)
			if resultDate.After(now) {
				break
			}
		}
		return resultDate.Format(format.YYYYMMDD), nil

	case "w":
		if len(repeats) < 2 {
			return "", errTooFewParameters
		}
		if len(repeats) > 2 {
			return "", errTooManyParameters
		}

		appLogger.Info("Applying weekly repetition rule")

		weekdays := map[int]time.Weekday{
			1: time.Monday, 2: time.Tuesday, 3: time.Wednesday, 4: time.Thursday,
			5: time.Friday, 6: time.Saturday, 7: time.Sunday,
		}
		var weekdaysName []time.Weekday
		const daysInWeek = 7

		for _, weekday := range strings.Split(repeats[1], ",") {
			weekdayNumber, err := strconv.Atoi(weekday)
			if err != nil {
				appLogger.Errorf("Invalid weekday number: %s", weekday)
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

		if !resultDate.After(now) {
			resultDate = now
		} else {
			if slices.Contains(weekdaysName, resultDate.Weekday()) {
				return resultDate.Format(format.YYYYMMDD), nil
			}
		}

		for {
			resultDate = resultDate.AddDate(0, 0, 1)
			if slices.Contains(weekdaysName, resultDate.Weekday()) {
				return resultDate.Format(format.YYYYMMDD), nil
			}
		}

	case "m":
		if len(repeats) < 2 {
			return "", errTooFewParameters
		}
		if len(repeats) > 3 {
			return "", errTooManyParameters
		}

		appLogger.Info("Applying monthly repetition rule")

		const (
			minMonthday, invalidMonthday, maxMonthday = -2, 0, 31
			minMonth, maxMonth                        = 1, 12
		)
		var (
			monthdays []int
			months    []time.Month
		)

		for _, monthday := range strings.Split(repeats[1], ",") {
			monthdayNumber, err := strconv.Atoi(monthday)
			if err != nil {
				appLogger.Errorf("Invalid monthday format: %s", monthday)
				return "", fmt.Errorf("%s: %w", errRuleParameter, err)
			}
			if monthdayNumber < minMonthday || monthdayNumber == invalidMonthday || monthdayNumber > maxMonthday {
				return "", fmt.Errorf("%s: %w", errRuleParameter, errors.New("invalid day of the month ordinal number"))
			}
			monthdays = append(monthdays, monthdayNumber)
		}

		hasMonthFilter := false
		if len(repeats) > 2 {
			hasMonthFilter = true
			for _, month := range strings.Split(repeats[2], ",") {
				monthNumber, err := strconv.Atoi(month)
				if err != nil {
					appLogger.Errorf("Invalid month format: %s", month)
					return "", fmt.Errorf("%s: %w", errRuleParameter, err)
				}
				if monthNumber < minMonth || monthNumber > maxMonth {
					return "", fmt.Errorf("%s: %w", errRuleParameter, errors.New("invalid month ordinal number"))
				}
				months = append(months, time.Month(monthNumber))
			}
		}

		if !resultDate.After(now) {
			resultDate = now
		} else {
			if matchMonthRule(resultDate, monthdays, months, hasMonthFilter) {
				return resultDate.Format(format.YYYYMMDD), nil
			}
		}

		for {
			resultDate = resultDate.AddDate(0, 0, 1)
			if matchMonthRule(resultDate, monthdays, months, hasMonthFilter) {
				return resultDate.Format(format.YYYYMMDD), nil
			}
		}

	default:
		appLogger.Errorf("Unsupported repeat rule: %s", repeats[0])
		return "", errUnsupportedFormat
	}
}

// matchMonthRule evaluates if a given date satisfies the combined month-of-year
// boundaries and relative day-of-month indices.
func matchMonthRule(date time.Time, monthdays []int, months []time.Month, hasMonthFilter bool) bool {
	if hasMonthFilter && !slices.Contains(months, date.Month()) {
		return false
	}

	dOfM := daysInMonth(date)
	for _, d := range monthdays {
		if d > 0 && date.Day() == d {
			return true
		}
		if d == -1 && date.Day() == dOfM {
			return true
		}
		if d == -2 && date.Day() == dOfM-1 {
			return true
		}
	}
	return false
}
