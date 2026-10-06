package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	title := "Diary 2027"
	// Inclusive.
	startDate := time.Date(2027, time.January, 1, 0, 0, 0, 0, time.UTC)
	// Not inclusive.
	endDate := startDate.AddDate(1, 0, 0)
	outputFile := "diary.html"
	numLines := 6

	f, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}

	currDate := startDate
	for currDate.Weekday() != time.Monday {
		currDate = currDate.AddDate(0, 0, -1)
	}
	fmt.Printf("Target start: %v\n", startDate)
	fmt.Printf("Actual start: %v\n", currDate)

	// Preamble
	fmt.Fprintln(f, `<!DOCTYPE html>`)
	fmt.Fprintln(f, `<html lang="en">`)
	fmt.Fprintln(f, `<head>`)
	fmt.Fprintln(f, `  <meta charset="UTF-8">`)
	fmt.Fprintf(f, "  <title>%v</title>\n", title)
	fmt.Fprintf(f, "  <meta name=\"Description\" content=\"%v\">\n", title)
	fmt.Fprintln(f, `  <link rel="stylesheet" type="text/css" href="horizontal.css">`)
	fmt.Fprintln(f, `</head>`)
	fmt.Fprintln(f, `<body>`)

	// Load the important dates.
	type occasion struct {
		name     string
		fontSize string
	}
	importantDates := map[time.Time][]occasion{}

	csvFile, err := os.Open("important-dates.csv")
	if err != nil {
		panic(err)
	}

	records, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		panic(err)
	}

	err = csvFile.Close()
	if err != nil {
		panic(err)
	}

	for _, record := range records[1:] {
		date, err := time.Parse("02/01/2006", record[0])
		if err != nil {
			panic(err)
		}
		importantDates[date] = append(importantDates[date], occasion{name: record[1], fontSize: strings.TrimSpace(record[2])})
	}
	fmt.Printf("Loaded %d important dates.\n", len(importantDates))

	// Add the initial page
	// We will start with the month of the targetStart even if we could fit it into the start of the next month.

	fmt.Fprintln(f, `<div class="page all">`)

	fmt.Fprintf(f, "<h1 id=\"title\" class=\"title\">%v %v &ndash; %v %v</h1>", startDate.Month(), startDate.Year(), endDate.AddDate(0, 0, -1).Month(), endDate.AddDate(0, 0, -1).Year())

	allCalDate := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, time.UTC)
	for allCalDate.Before(endDate) {

		calYear := allCalDate.Year()
		calMonth := allCalDate.Month()

		fmt.Fprintln(f, `<div class="calendar">`)
		fmt.Fprintf(f, "<span class=\"cal-month\">%v %v</span>\n", calMonth, calYear)
		fmt.Fprintln(f, `<span class="cal-day">M</span>`)
		fmt.Fprintln(f, `<span class="cal-day">T</span>`)
		fmt.Fprintln(f, `<span class="cal-day">W</span>`)
		fmt.Fprintln(f, `<span class="cal-day">T</span>`)
		fmt.Fprintln(f, `<span class="cal-day">F</span>`)
		fmt.Fprintln(f, `<span class="cal-day">S</span>`)
		fmt.Fprintln(f, `<span class="cal-day">S</span>`)

		// Find the first Monday

		currCalDate := allCalDate
		for currCalDate.Weekday() != time.Monday {
			currCalDate = currCalDate.AddDate(0, 0, -1)
		}

		endCal := time.Date(calYear, calMonth, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, 0)
		for currCalDate.Before(endCal) {
			linkTarget := currCalDate.Format("2006-01-02")
			for i := 0; i < 7; i++ {
				class := "cal-entry"
				if currCalDate.Month() != calMonth {
					class += " cal-other-month"
				}

				fmt.Fprintf(f, "<a href=\"#%v\"><span class=\"%v\">%v</span></a>\n", linkTarget, class, currCalDate.Day())
				currCalDate = currCalDate.AddDate(0, 0, 1)
			}
		}
		fmt.Fprintln(f, "</div>")

		allCalDate = allCalDate.AddDate(0, 1, 0)
	}

	fmt.Fprintln(f, "</div>")

	// Generate the gridTemplateRows
	gridRowTemplate := ""
	for _ = range numLines {
		gridRowTemplate += "1fr "
	}

	for currDate.Before(endDate) {
		fmt.Fprintf(f, `<div id="%v" class="page week">`, currDate.Format("2006-01-02"))

		// Add the calendar
		tmp := currDate
		inCurrMonth := 0
		for i := 0; i < 7; i++ {
			if tmp.Month() == currDate.Month() {
				inCurrMonth++
			}
			tmp = tmp.AddDate(0, 0, 1)
		}

		calMonth := currDate.Month()
		calYear := currDate.Year()

		if inCurrMonth <= 3 {
			nextMonth := time.Date(calYear, calMonth, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, 0)
			calMonth = nextMonth.Month()
			calYear = nextMonth.Year()
		}

		// Force the selected month to be within the months on the first page.

		// Cannot be too early.
		if time.Date(calYear, calMonth, 1, 0, 0, 0, 0, time.UTC).Before(time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, time.UTC)) {
			nextMonth := time.Date(calYear, calMonth, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, 0)
			calMonth = nextMonth.Month()
			calYear = nextMonth.Year()
		}

		// Cannot be too late.
		lastDay := endDate.AddDate(0, 0, -1)
		if time.Date(calYear, calMonth, 1, 0, 0, 0, 0, time.UTC).After(time.Date(lastDay.Year(), lastDay.Month(), 1, 0, 0, 0, 0, time.UTC)) {
			prevMonth := time.Date(calYear, calMonth, 1, 0, 0, 0, 0, time.UTC).AddDate(0, -1, 0)
			calMonth = prevMonth.Month()
			calYear = prevMonth.Year()
		}

		fmt.Fprintln(f, `<div class="calendar">`)
		fmt.Fprintf(f, "<span class=\"cal-month\">%v %v</span>\n", calMonth, calYear)
		fmt.Fprintln(f, `<span class="cal-day">M</span>`)
		fmt.Fprintln(f, `<span class="cal-day">T</span>`)
		fmt.Fprintln(f, `<span class="cal-day">W</span>`)
		fmt.Fprintln(f, `<span class="cal-day">T</span>`)
		fmt.Fprintln(f, `<span class="cal-day">F</span>`)
		fmt.Fprintln(f, `<span class="cal-day">S</span>`)
		fmt.Fprintln(f, `<span class="cal-day">S</span>`)

		// Find the first date in the calendar
		currCalDate := currDate
		for currCalDate.After(time.Date(calYear, calMonth, 1, 0, 0, 0, 0, time.UTC)) {
			currCalDate = currCalDate.AddDate(0, 0, -7)
		}
		endCal := time.Date(calYear, calMonth, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, 0)
		for currCalDate.Before(endCal) {
			linkTarget := currCalDate.Format("2006-01-02")
			currentWeek := false
			if currCalDate == currDate {
				currentWeek = true
			}
			for i := 0; i < 7; i++ {
				class := "cal-entry"
				if currentWeek {
					class += " cal-week"
				}
				if currCalDate.Month() != calMonth {
					class += " cal-other-month"
				}
				fmt.Fprintf(f, "<a href=\"#%v\"><span class=\"%v\">%v</span></a>\n", linkTarget, class, currCalDate.Day())
				currCalDate = currCalDate.AddDate(0, 0, 1)
			}
		}
		fmt.Fprintln(f, "</div>")

		// Now add each of the boxes
		for i := 0; i < 7; i++ {
			fmt.Fprintf(f, "<div class=\"day %v\" style=\"grid-template-rows:%v\">\n", strings.ToLower(currDate.Weekday().String()), gridRowTemplate)
			fmt.Fprintln(f, "<div class=\"day-title\">")
			fmt.Fprintf(f, "<h2 class=\"day-num\">%v</h2>\n", currDate.Day())
			fmt.Fprintf(f, "<span class=\"day-name\">%v</span>\n", currDate.Weekday().String())
			fmt.Fprintf(f, "<span class=\"day-month\">%v</span>\n", currDate.Month())
			fmt.Fprintln(f, "</div>")
			if currDate.Weekday() == time.Monday {
				fmt.Fprintln(f, "<a href=\"#title\" class=\"home-anchor\"><img class=\"home\" src=\"calendar-tight.svg\" alt=\"View entire year\"/></a>")
			}
			j := 0
			if occasions, ok := importantDates[currDate]; ok {
				for ; j < len(occasions); j++ {
					fmt.Fprintln(f, "<div class=\"lines\">")
					if occasions[j].fontSize != "" {
						fmt.Fprintf(f, "<span class=\"occasion\" style=\"font-size: %v\">%v</span>\n", occasions[j].fontSize, occasions[j].name)
					} else {
						fmt.Fprintf(f, "<span class=\"occasion\">%v</span>\n", occasions[j].name)
					}
					fmt.Fprintln(f, "</div>")
				}

			}

			for ; j < numLines; j++ {
				fmt.Fprintln(f, "<div class=\"lines\"></div>")
			}
			if currDate.Weekday() != time.Saturday && currDate.Weekday() != time.Sunday {
				for i := 0; i < numLines; i++ {
					fmt.Fprintln(f, "<div class=\"lines\"></div>")
				}
			}
			fmt.Fprintln(f, "</div>")
			currDate = currDate.AddDate(0, 0, 1)
		}

		fmt.Fprintln(f, "</div>")

	}

	// End
	end := `
</body>

</html>`

	fmt.Fprintln(f, end)
}
