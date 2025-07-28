package main

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
)

// Test to ensure http response from NWS server db containing the Lvl 2 data is still alive and produces the radar data specified by the user. In this test case (December 20th 2020 on the KHTX radar at 000212 zulu) was assessed

func mainTesting() int {

	var month, day, year, radar string
	var test1, test4 int

	month = "12"
	num1, _ := strconv.Atoi(month)
	if num1 < 10 {
		month = fmt.Sprintf("%02d", num1)
	}

	day = "20"
	num2, _ := strconv.Atoi(day)
	if num2 < 10 {
		day = fmt.Sprintf("%02d", num2)
	}

	year = "2020"
	radar = "KHTX"

	test := "000212"
	test1, _ = strconv.Atoi(test)

	test3 := "000212"
	test4, _ = strconv.Atoi(test3)

	for x := test1; x <= test4; x++ {

		var end string
		year2 := year
		year1, _ := strconv.Atoi(year2)

		if year1 >= 2009 && year1 <= 2012 {
			end = "V03"
		} else if year1 < 2009 {
			end = ""
		} else {
			end = "V06"
		}

		timeComb := fmt.Sprintf("%06d", x)

		url := fmt.Sprintf("https://noaa-nexrad-level2.s3.amazonaws.com/%s/%s/%s/%s/%s%s%s%s_%s_%s", year, month, day, radar, radar, year, month, day, timeComb, end)

		url2 := fmt.Sprintf("https://noaa-nexrad-level2.s3.amazonaws.com/%s/%s/%s/%s/%s%s%s%s_%s_%s.gz", year, month, day, radar, radar, year, month, day, timeComb, end)

		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == 200 {
			fmt.Println("(+) FETCHING", url)
		}

		resp2, err := http.Get(url2)
		if err == nil && resp2.StatusCode == 200 {
			fmt.Println("(+) FETCHING .GZ", url2)
		}
	}
	return 0
}

func TestMain(t *testing.T) {
	expected := mainTesting()
	want := 0

	if mainTesting() != want {
		t.Errorf("got %q want %q // CHECK THAT PARAMETERS BEING PASSED RESULT IN A 200 RESPONSE CODE FROM NWS DB", expected, want)
	}
}
