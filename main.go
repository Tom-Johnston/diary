package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	targetStart := time.Date(2025, time.October, 1, 0, 0, 0, 0, time.UTC)
	targetEnd := targetStart.AddDate(1, 0, 0)
	outputFile := "diary.html"
	numLines := 8

	f, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}

	currDate := targetStart
	for currDate.Weekday() != time.Monday {
		currDate = currDate.AddDate(0, 0, -1)
	}
	fmt.Printf("Target start: %v\n", targetStart)
	fmt.Printf("Actual start: %v\n", currDate)

	// Preamble
	preamble := `
<!DOCTYPE html>
<html lang="en">

<head>
  <meta charset="UTF-8">
  <title>2025 Academic Diary</title>
  <meta name="Description" content="2025 Academic Diary">
  <link rel="stylesheet" type="text/css" href="horizontal.css">
</head>

<body>`

	fmt.Fprintln(f, preamble)

	// Add the initial page
	// We will start with the month of the targetStart even if we could fit it into the start of the next month.

	fmt.Fprintln(f, `<div class="page all">`)

	fmt.Fprintf(f, "<h1 id=\"title\" class=\"title\">%v %v &ndash; %v %v</h1>", targetStart.Month(), targetStart.Year(), targetEnd.AddDate(0, 0, -1).Month(), targetEnd.AddDate(0, 0, -1).Year())

	allCalDate := time.Date(targetStart.Year(), targetStart.Month(), 1, 0, 0, 0, 0, time.UTC)
	for allCalDate.Before(targetEnd) {

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
			for i := 0; i < 7; i++ {
				class := "cal-entry"
				if currCalDate.Month() != calMonth {
					class += " cal-other-month"
				}
				fmt.Fprintf(f, "<span class=\"%v\">%v</span>\n", class, currCalDate.Day())
				currCalDate = currCalDate.AddDate(0, 0, 1)
			}
		}
		fmt.Fprintln(f, "</div>")

		allCalDate = allCalDate.AddDate(0, 1, 0)
	}

	fmt.Fprintln(f, "</div>")

	for currDate.Before(targetEnd) {
		fmt.Fprintln(f, `<div class="page week">`)

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
			calMonth = currDate.AddDate(0, 1, 0).Month()
			calYear = currDate.AddDate(0, 1, 0).Year()
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
				fmt.Fprintf(f, "<span class=\"%v\">%v</span>\n", class, currCalDate.Day())
				currCalDate = currCalDate.AddDate(0, 0, 1)
			}
		}
		fmt.Fprintln(f, "</div>")

		// Now add each of the boxes
		for i := 0; i < 7; i++ {
			fmt.Fprintf(f, "<div class=\"day %v\">\n", strings.ToLower(currDate.Weekday().String()))
			fmt.Fprintln(f, "<div class=\"day-title\">")
			fmt.Fprintf(f, "<h2 class=\"day-num\">%v</h2>\n", currDate.Day())
			fmt.Fprintf(f, "<span class=\"day-name\">%v</span>\n", currDate.Weekday().String())
			fmt.Fprintf(f, "<span class=\"day-month\">%v</span>\n", currDate.Month())
			fmt.Fprintln(f, "</div>")
			if currDate.Weekday() == time.Monday {
				fmt.Fprintln(f, "<a href=\"#title\" class=\"home-anchor\"><img class=\"home\" src=\"calendar-tight.svg\" alt=\"View entire year\"/></a>")
			}
			for i := 0; i < numLines; i++ {
				fmt.Fprintln(f, "<div class=\"dotted\"></div>")
			}
			if currDate.Weekday() != time.Saturday && currDate.Weekday() != time.Sunday {
				for i := 0; i < numLines; i++ {
					fmt.Fprintln(f, "<div class=\"dotted\"></div>")
				}
			}
			fmt.Fprintln(f, "</div>")
			currDate = currDate.AddDate(0, 0, 1)
		}
		// <div class="day mon">
		//     <div class="day-title">
		//         <h2 class="day-num">8</h2>
		//         <span class="day-name">Monday</span>
		//         <span class="day-month">October</span>
		//     </div>
		//     <img class="home" src="calendar-tight.svg" alt="View entire year"/>
		//     <div class="dotted"></div>
		//     <div class="dotted"></div>
		//     <div class="dotted"></div>
		//     <div class="dotted"></div>
		//     <div class="dotted"></div>
		//     <div class="dotted"></div>
		//     <div class="dotted"></div>
		//     <div class="dotted"></div>
		//     <div class="dotted"></div>
		//     <div class="dotted"></div>
		//     <div class="dotted"></div>
		//     <div class="dotted"></div>
		//     <div class="dotted"></div>
		//     <div class="dotted"></div>
		//     <div class="dotted"></div>
		//     <div class="dotted"></div>
		// </div>

		fmt.Fprintln(f, "</div>")

	}

	// End
	end := `
</body>

</html>`

	fmt.Fprintln(f, end)
}
