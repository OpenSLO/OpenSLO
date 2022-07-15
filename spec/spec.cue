package v1alpha1

import "time"

#spec: {
	// list?
	service: string
	// optional
	// TODO: 1050 exceeds go regex limit (1000)
	description?: string
	// object oneOf
	timeWindow:      #timeWindowSecond | #timeWindowCalendar
	budgetingMethod: "Occurrences" | "Timeslices"
    // objectives:
    // indicator:
}

#indicator: {

}

#thresholdMetric: {
	source:    string
	queryType: string
	query:     string
}

#timeWindow: {
	unit:      string
	count:     string
	isRolling: bool
}

#timeWindowSecond: {
	#timeWindow & {
		unit:  "Second"
		count: "numeric"
	}
}

#timeWindowCalendar: {
	#timeWindow & {
		unit:  "Year" | "Quarter" | "Month" | "Week" | "Day"
		count: "numeric"
	}
	calendar: {
		// note: leap-seconds not supported
		startTime: time.Format(time.RFC3339)
	}
}
