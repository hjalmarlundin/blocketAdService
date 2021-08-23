package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var a App

func Test_DELETE_ErrorIfNotExist(t *testing.T) {

	// Arrange
	a.initialize("testDb.db")
	clearTable()
	addTestData()

	// Act
	req, _ := http.NewRequest("DELETE", "/ad/23423940234", nil)
	response := executeRequest(req)

	// Assert
	checkResponsecode(t, http.StatusNotFound, response.Code)
}

func Test_DELETE_DeletesCorrectUser(t *testing.T) {

	// Arrange
	a.initialize("testDb.db")
	clearTable()
	addTestData()

	ads := getAdverts(t)
	if len(ads) != 4 {
		t.Errorf("Expected 4 elements in list. Received %v", len(ads))
	}

	// Act
	req, _ := http.NewRequest("DELETE", "/ad/1", nil)
	response := executeRequest(req)
	checkResponsecode(t, http.StatusOK, response.Code)

	// Assert
	adsAfterDelete := getAdverts(t)
	if len(adsAfterDelete) != 3 {
		t.Errorf("Expected 3 elements in the after delete list. Received %v", len(adsAfterDelete))
	}
	for _, v := range adsAfterDelete {
		if v.ID == "1" {
			t.Errorf("Expected ID 1 to be deleted. Received %v", len(adsAfterDelete))
		}
	}

}

func Test_POST_AddsToDb(t *testing.T) {

	// Arrange
	a.initialize("testDb.db")
	clearTable()

	ads := getAdverts(t)
	if len(ads) != 0 {
		t.Errorf("Expected database to be empty. Received %v", len(ads))
	}
	newAd := &Ad{ID: "1", Subject: "Add1", Body: "This is what Iam selling", Email: "myEmal@hotmail.com", Created: time.Now()}
	asJSON, _ := json.Marshal(newAd)

	// Act
	req, _ := http.NewRequest("POST", "/ad", bytes.NewBuffer(asJSON))
	response := executeRequest(req)

	// Assert
	checkResponsecode(t, http.StatusOK, response.Code)
	adsAfter := getAdverts(t)
	if len(adsAfter) != 1 {
		t.Errorf("Expected database to contain one entry. Received %v", len(adsAfter))
	}

	if newAd.Body != adsAfter[0].Body {
		t.Errorf("Expected: %v, Received: %v", newAd.Body, adsAfter[0].Body)
	}
}

func Test_GET_ReturnsAllUsers(t *testing.T) {

	// Arrange
	a.initialize("testDb.db")
	clearTable()
	addTestData()

	// Act
	ads := getAdverts(t)

	// Assert
	if len(ads) != 4 {
		t.Errorf("Expected 4 elements in list. Received %v", len(ads))
	}

}

func Test_GET_PriceIsOptional(t *testing.T) {

	// Arrange
	a.initialize("testDb.db")
	clearTable()

	entry := addASingleTestData(nil)

	// Act
	ads := getAdverts(t)

	// Assert
	if len(ads) != 1 {
		t.Errorf("Expected 1 elements in list. Received %v", len(ads))
	}

	if ads[0].Price != nil {
		t.Errorf("Expected price to be null: %v. Received %v", entry.Price, ads[0].Price)
	}

}

func Test_GET_ReturnsCorrectData(t *testing.T) {
	// Arrange
	a.initialize("testDb.db")
	clearTable()

	price := 40.43
	entry := addASingleTestData(&price)

	// Act
	ads := getAdverts(t)

	// Assert
	if len(ads) != 1 {
		t.Errorf("Expected 1 elements in list. Received %v", len(ads))
	}

	for _, v := range ads {

		if v.Body != entry.Body {
			t.Errorf("Expected body to be the same: %v. Received %v", entry.Body, v.Body)
		}
		if v.Subject != entry.Subject {
			t.Errorf("Expected subject to be the same: %v. Received %v", entry.Subject, v.Subject)
		}
		if v.Email != entry.Email {
			t.Errorf("Expected email to be the same: %v. Received %v", entry.Email, v.Email)
		}
		if *v.Price != *entry.Price {
			t.Errorf("Expected price to be the same: %v. Received %v", *entry.Price, *v.Price)
		}

	}

}

func Test_GET_SortByPriceDesc(t *testing.T) {
	// Arrange
	a.initialize("testDb.db")
	clearTable()
	addTestData()

	// Act
	allAdsVerts := get(t, "/ads?sort_by=price.desc")

	// Assert
	if len(allAdsVerts) != 4 {
		t.Errorf("Expected four results from get request. Got %v\n", len(allAdsVerts))
	}

	expectedPrices := []float64{900.5, 134.13, 32.42}
	AssertPrice(allAdsVerts[0], expectedPrices[0], t)
	AssertPrice(allAdsVerts[1], expectedPrices[1], t)
	AssertPrice(allAdsVerts[2], expectedPrices[2], t)
	AssertPriceIsNull(allAdsVerts[3], t)
}

func Test_GET_SortByPriceAsc(t *testing.T) {
	// Arrange
	a.initialize("testDb.db")
	clearTable()
	addTestData()

	// Act
	allAdsVerts := get(t, "/ads?sort_by=price.asc")

	// Assert
	if len(allAdsVerts) != 4 {
		t.Errorf("Expected four results from get request. Got %v\n", len(allAdsVerts))
	}

	expectedPrices := []float64{900.5, 134.13, 32.42}
	AssertPrice(allAdsVerts[3], expectedPrices[0], t)
	AssertPrice(allAdsVerts[2], expectedPrices[1], t)
	AssertPrice(allAdsVerts[1], expectedPrices[2], t)
	AssertPriceIsNull(allAdsVerts[0], t)
}

func Test_GET_SortByCreatedAsc(t *testing.T) {
	// Arra ge
	a.initialize("testDb.db")
	clearTable()
	addTestData()

	// Act
	allAdsVerts := get(t, "/ads?sort_by=date.asc")

	// Assert
	if len(allAdsVerts) != 4 {
		t.Errorf("Expected four results from get request. Got %v\n", len(allAdsVerts))
	}

	expectedDates := []time.Time{time.Date(2021, time.August, 0, 0, 0, 0, 0, time.Local), time.Date(2021, time.January, 0, 0, 0, 0, 0, time.Local),
		time.Date(2020, time.April, 0, 0, 0, 0, 0, time.Local), time.Date(1988, time.December, 0, 0, 0, 0, 0, time.Local)}
	AssertTime(allAdsVerts[0], expectedDates[3], t)
	AssertTime(allAdsVerts[1], expectedDates[2], t)
	AssertTime(allAdsVerts[2], expectedDates[1], t)
	AssertTime(allAdsVerts[3], expectedDates[0], t)

}

func Test_GET_SortByCreatedDesc(t *testing.T) {
	// Arrange
	a.initialize("testDb.db")
	clearTable()
	addTestData()

	// Act
	allAdsVerts := get(t, "/ads?sort_by=date.desc")

	// Assert
	if len(allAdsVerts) != 4 {
		t.Errorf("Expected four results from get request. Got %v\n", len(allAdsVerts))
	}

	expectedDates := []time.Time{time.Date(2021, time.August, 0, 0, 0, 0, 0, time.Local), time.Date(2021, time.January, 0, 0, 0, 0, 0, time.Local), time.Date(2020, time.April, 0, 0, 0, 0, 0, time.Local), time.Date(1988, time.December, 0, 0, 0, 0, 0, time.Local)}
	AssertTime(allAdsVerts[0], expectedDates[0], t)
	AssertTime(allAdsVerts[1], expectedDates[1], t)
	AssertTime(allAdsVerts[2], expectedDates[2], t)
	AssertTime(allAdsVerts[3], expectedDates[3], t)

}

// Helper functions

func AssertTime(add1 Ad, t time.Time, test *testing.T) {
	if add1.Created != t {
		test.Errorf("Expected %v. Got %v\n", t, add1.Created)
	}
}

func AssertPrice(add1 Ad, expectedPrice float64, t *testing.T) {
	if *add1.Price != expectedPrice {
		t.Errorf("Expected %v. Got %v\n", expectedPrice, *add1.Price)
	}
}

func AssertPriceIsNull(add1 Ad, t *testing.T) {
	if add1.Price != nil {
		t.Errorf("Expected price to be null. Got %v\n", *add1.Price)
	}
}

func get(t *testing.T, query string) []Ad {
	request, _ := http.NewRequest("GET", query, nil)
	response := executeRequest(request)
	checkResponsecode(t, http.StatusOK, response.Code)
	reqBody, _ := ioutil.ReadAll(response.Body)

	var ads []Ad
	json.Unmarshal(reqBody, &ads)
	return ads
}

func getAdverts(t *testing.T) []Ad {
	return get(t, "/ads")
}

func checkResponsecode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func addTestData() {
	price1 := 32.42
	price2 := 134.13
	price3 := 900.5

	a.DB.Create(&Ad{ID: "1", Subject: "Add1", Body: "This is what Iam selling", Price: &price1, Email: "myEmal@hotmail.com", Created: time.Date(2021, time.August, 0, 0, 0, 0, 0, time.Local)})
	a.DB.Create(&Ad{ID: "2", Subject: "Add2", Body: "This is another ad", Email: "test@msn.com", Price: &price2, Created: time.Date(2021, time.January, 0, 0, 0, 0, 0, time.Local)})
	a.DB.Create(&Ad{ID: "3", Subject: "Add3", Body: "This is another ad", Email: "anotherEmail@spray.com", Price: &price3, Created: time.Date(2020, time.April, 0, 0, 0, 0, 0, time.Local)})
	a.DB.Create(&Ad{ID: "4", Subject: "Add4", Body: "This is another ad", Email: "emalj@Hoppla.com", Created: time.Date(1988, time.December, 0, 0, 0, 0, 0, time.Local)})
}

func addASingleTestData(price *float64) Ad {
	x := &Ad{ID: "1", Subject: "Add1", Body: "This is what Iam selling", Price: price, Email: "myEmal@hotmail.com", Created: time.Now()}
	a.DB.Create(x)
	return *x
}

func clearTable() {
	a.DB.Where("1 = 1").Delete(&Ad{})
}
