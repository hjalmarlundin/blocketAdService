package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Ad-service blocket code exercise")

	a := App{}
	a.initialize("adServiceDb.db")

	// a.DB.Create(&Ad{ID: "1", Subject: "Add1", Body: "This is what Iam selling", Price: 32.42, Email: "myEmal@hotmail.com", Created: time.Now()})
	// a.DB.Create(&Ad{ID: "2", Subject: "Add2", Body: "This is another ad", Email: "anotherEmail@msn.com", Price: 134.13, Created: time.Date(2021, time.January, 3, 0, 0, 0, 0, time.Local)})
	// a.DB.Create(&Ad{ID: "3", Subject: "Add3", Body: "This is another ad", Email: "anotherEmail@msn.com", Price: 900.5, Created: time.Date(2020, time.April, 3, 0, 0, 0, 0, time.Local)})
	a.DB.Create(&Ad{ID: "4", Subject: "Add4", Body: "This is another ad", Email: "anotherEmail@msn.com", Created: time.Date(1988, time.December, 3, 0, 0, 0, 0, time.Local)})
	a.run(":10000")

	// // price1 := 32.42
	// // price2 := 134.42
	// // price3 := 4.2312
	// // Create som test data

}

func hej() {}
