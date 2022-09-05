package Benchmark

import (
	"fmt"
	"strconv"
	"time"
)

var Kilo float64 = 1000
var Mega float64 = 1000 * Kilo
var Giga float64 = 1000 * Mega
var Tera float64 = 1000 * Giga
var Peta float64 = 1000 * Tera
var ThroughputMessage chan string
var verbose = false

func CalculateThroughput(startTime time.Time, endTime time.Time, MessageSize int, NbMessages int) string {

	dataBytes := float64(MessageSize * NbMessages)
	dataBits := dataBytes * 8

	if verbose {

		PrintDataAmount(float64(dataBytes))

		fmt.Println("Starting Time: ", startTime)

		fmt.Println("Finish Time: ", endTime)
	}

	elapsedTime := endTime.Sub(startTime)

	if verbose {

		fmt.Println("Elapsed Time: ", elapsedTime)
	}

	Throughput := dataBits / elapsedTime.Seconds()

	return getThroughputString(Throughput)
}

func CalculateResponseTime(startTime time.Time, endTime time.Time) time.Duration {

	if verbose {

		fmt.Println("Starting Time: ", startTime)

		fmt.Println("Finish Time: ", endTime)
	}

	return endTime.Sub(startTime)
}

func getThroughputString(Throughput float64) string {

	return getMetric(Throughput) + "b/s"
}

func PrintDataAmount(dataAmount float64) {

	fmt.Printf("Data Size: %sB\n", getMetric(dataAmount))
}

func getMetric(size float64) string {

	if size >= Peta {
		return strconv.FormatFloat(size/Peta, 'f', 2, 32) + "P"
	} else if size >= Tera {
		return strconv.FormatFloat(size/Tera, 'f', 2, 32) + "T"
	} else if size >= Giga {
		return strconv.FormatFloat(size/Giga, 'f', 2, 32) + "G"
	} else if size >= Mega {
		return strconv.FormatFloat(size/Mega, 'f', 2, 32) + "M"
	} else if size >= Kilo {
		return strconv.FormatFloat(size/Kilo, 'f', 2, 32) + "K"
	} else {
		return strconv.FormatFloat(size, 'f', 2, 32)
	}
}
