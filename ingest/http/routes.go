package http

import (
	"fmt"
	ias_pg "ias/automation/db/pg"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	ias_pg.NewPostgresStorage(nil).QueryData("select device_name from ppj_tree_sensor")
	fmt.Fprint(w, "Hello World!")
}

func getAllTreeSensorHandler(w http.ResponseWriter, r *http.Request) {

	sensors, _ := ias_pg.NewPostgresStorage(nil).QueryData("select * from ppj_tree_sensor")
	for _, sensor := range sensors {
		fmt.Fprintf(w, "Sensor: %+v\n", sensor)
	}
}
