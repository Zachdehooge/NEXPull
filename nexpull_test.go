package main

import (
	"net/http"
	"testing"
)

// Test to ensure http response from NWS server db containing the Lvl 2 data is still alive and produces the radar data specified by the user. In this test case (December 20th 2020 on the KHTX radar at 000212 zulu) was assessed

func MainTesting() int {

	// TODO: Rework tests to mimic download and to test if the SW3 bucket is still working
	// ! URL's to evaluate are below.

	url := "https://unidata-nexrad-level2.s3.amazonaws.com/2008/10/09/KHTX/KHTX20081009_000833_V03.gz"

	url2 := "https://unidata-nexrad-level2.s3.amazonaws.com/2026/03/15/KHTX/KHTX20260315_024508_V06"

	resp, _ := http.Get(url)
	resp2, _ := http.Get(url2)

	if resp.StatusCode != 200 {
		return 1
	}

	if resp2.StatusCode != 200 {
		return 1
	}

	resp.Body.Close()
	resp2.Body.Close()

	return 0
}

func TestMainTesting(t *testing.T) {
	want := 0
	got := MainTesting()

	if want != got {
		t.Errorf("MainTesting() = %d; wanted %d", got, want)
	}
}
