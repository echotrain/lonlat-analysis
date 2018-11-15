package main

import (
	"log"
	"bufio"
	"encoding/csv"
	"fmt"
	"math"
	"os"

	"github.com/echotrain/lonlat-analysis/db"
)

// LonLatData holds the queried lonlat data
var LonLatData struct {
	StandardLon string `db:"s_lon"`
	StandardLat string `db:"s_lat"`
	ProviderLon string `db:"p_lon"`
	ProviderLat string `db:"p_lat"`
}

// Report stores a string record for each lonlat that's out of range.
type Report struct {
	Records       [][]string
	standardCount int
	providerCount int
	mismatchCount int
}

func main() {
	fmt.Println("Hello world!")
	conn, err := newConnection()
	if err != nil {
		return err
	}

	report := Report{}

	csvHeaders := []string{
		"standard_id", "hex_display_id", "provider", "external_id",
		"standard_name", "provider_name", "standard_address", "standard_city",
		"standard_state", "provider_address", "provider_city", "provider_state",
	}

	report.Records = append(report.Records, csvHeaders)

	csvFile, _ := os.Open("records.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))

	for {
		line, err := reader.Read()
		if err != nil {
			log.Fatal(error)
		}
		
		lonlatData := LonLatData{}
		conn.DB.Get(
			&lonlatData,
			lonlatQuery(line[2]),
			line[0],
			line[3],
		)


	}


}

// newConnection centralizes database connection. User responsible for closing their connection.
func newConnection() (conn *db.Connection, err error) {
	// setup database connections
	return db.NewConnection(
		os.Getenv("LOCAL_PGSQL_DB_HOST"),
		os.Getenv("LOCAL_PGSQL_DB_USERNAME"),
		os.Getenv("LOCAL_PGSQL_DB_PASSWORD"),
		os.Getenv("LOCAL_PGSQL_DB_NAME"),
		os.Getenv("LOCAL_PGSQL_DB_SSLMODE"),
		os.Getenv("LOCAL_PGSQL_DB_PORT"),
	)
}

fun lonlatQuery(table string) string {
	return fmt.Sprintf(
		`
		SELECT 
			s.longitude as s_lon, 
			s.latitude as s_lat, 
			p.longitude as p_lon, 
			p.latitude as p_lat 
		FROM standard_hotels s 
		LEFT OUTER JOIN %s_hotels p ON s.id=p.standard_hotel_id 
		WHERE s.id=$1 AND p.external_id=$2; 
		`,
		table,
	)
}

// haversin(theta) function
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// Distance function returns the distance (in meters) between two points
func Distance(lon1, lat1, lon2, lat2 float64) float64 {
	var lo1, la1, lo2, la2, r float64
	lo1 = lon1 * math.Pi / 180
	la1 = lat1 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180
	la2 = lat2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}
