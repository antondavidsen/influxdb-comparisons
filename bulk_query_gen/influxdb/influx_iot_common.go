package influxdb

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	bulkDataGenIot "github.com/antondavidsen/influxdb-comparisons/bulk_data_gen/iot"
	bulkQuerygen "github.com/antondavidsen/influxdb-comparisons/bulk_query_gen"
)

// InfluxIot produces Influx-specific queries for all the devops query types.
type InfluxIot struct {
	InfluxCommon
}

// NewInfluxIotCommon makes an InfluxIot object ready to generate Queries.
func NewInfluxIotCommon(lang Language, dbConfig bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	if _, ok := dbConfig[bulkQuerygen.DatabaseName]; !ok {
		panic("need influx database name")
	}

	return &InfluxIot{
		InfluxCommon: *newInfluxCommon(lang, dbConfig[bulkQuerygen.DatabaseName], queriesFullRange, scaleVar),
	}
}

// Dispatch fulfills the QueryGenerator interface.
func (d *InfluxIot) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	bulkQuerygen.IotDispatchAll(d, i, q, d.ScaleVar)
	return q
}

func (d *InfluxIot) AverageTemperatureDayByHourOneHome(q bulkQuerygen.Query) {
	d.averageTemperatureDayByHourNHomes(q.(*bulkQuerygen.HTTPQuery), 1, time.Hour*12)
}

// averageTemperatureDayByHourNHomes populates a Query with a query that looks like:
// SELECT avg(temperature) from air_condition_room where (home_id = '$HHOME_ID_1' or ... or hostname = '$HOSTNAME_N') and time >= '$HOUR_START' and time < '$HOUR_END' group by time(1h)
func (d *InfluxIot) averageTemperatureDayByHourNHomes(qi bulkQuerygen.Query, nHomes int, timeRange time.Duration) {
	interval := d.AllInterval.RandWindow(timeRange)
	nn := rand.Perm(d.ScaleVar)[:nHomes]

	homes := []string{}
	for _, n := range nn {
		homes = append(homes, fmt.Sprintf(bulkDataGenIot.SmartHomeIdFormat, n))
	}

	homeClauses := []string{}
	for _, s := range homes {
		if d.language == InfluxQL {
			homeClauses = append(homeClauses, fmt.Sprintf("home_id = '%s'", s))
		} else {
			homeClauses = append(homeClauses, fmt.Sprintf(`r.home_id == "%s"`, s))
		}
	}

	combinedHomesClause := strings.Join(homeClauses, " or ")

	var query string
	if d.language == InfluxQL {
		query = fmt.Sprintf("SELECT mean(temperature) from air_condition_room where (%s) and time >= '%s' and time < '%s' group by time(1h)", combinedHomesClause, interval.StartString(), interval.EndString())
	} else {
		query = fmt.Sprintf(`from(bucket:"%s") `+
			`|> range(start:%s, stop:%s) `+
			`|> filter(fn:(r) => r._measurement == "air_condition_room" and r._field == "temperature" and (%s)) `+
			`|> aggregateWindow(every:1h, fn:mean) `+
			`|> yield()`,
			d.DatabaseName,
			interval.StartString(), interval.EndString(),
			combinedHomesClause)
	}

	humanLabel := fmt.Sprintf("InfluxDB (%s) mean temperature, rand %4d homes, rand %s by 1h", d.language.String(), nHomes, timeRange)
	q := qi.(*bulkQuerygen.HTTPQuery)
	d.getHttpQuery(humanLabel, interval.StartString(), query, q)
}

func (d *InfluxIot) IotAggregateKeep(qi bulkQuerygen.Query, timeRange time.Duration) {
	interval := d.AllInterval.RandWindow(timeRange)

	var query string
	if d.language == InfluxQL {
		query = fmt.Sprintf(`SELECT mean("co2_level") as "mean_value" FROM "air_quality_room" WHERE time > '%s' AND time < '%s' AND "room_id"='4' GROUP BY time(5m) FILL(null)`,
			interval.StartString(),
			interval.EndString())
	} else {
		query = fmt.Sprintf(`from(bucket: "%s") `+
			`|> range(start: %s, stop: %s) `+
			`|> filter(fn: (r) => r._measurement == "air_quality_room" and r._field == "co2_level")`+
			`|> filter(fn: (r) => r.room_id == "4")`+
			`|> keep(columns: ["_time", "_value"])`+
			`|> aggregateWindow(every: 5m, fn: mean)`,
			d.DatabaseName,
			interval.StartString(),
			interval.EndString())
	}

	humanLabel := fmt.Sprintf(`InfluxDB (%s) aggregate/keep`, d.language)
	q := qi.(*bulkQuerygen.HTTPQuery)
	d.getHttpQuery(humanLabel, interval.StartString(), query, q)
}

func (d *InfluxIot) IotAggregateDrop(qi bulkQuerygen.Query, timeRange time.Duration) {
	interval := d.AllInterval.RandWindow(timeRange)

	var query string
	if d.language == InfluxQL {
		query = fmt.Sprintf(`SELECT mean("co2_level") as "mean_value" FROM "air_quality_room" WHERE time > '%s' AND time < '%s' AND "room_id"='4' GROUP BY time(5m) FILL(null)`,
			interval.StartString(),
			interval.EndString())
	} else {
		query = fmt.Sprintf(`from(bucket: "%s") `+
			`|> range(start: %s, stop: %s) `+
			`|> filter(fn: (r) => r._measurement == "air_quality_room" and r._field == "co2_level")`+
			`|> filter(fn: (r) => r.room_id == "4")`+
			`|> drop(columns: ["_start", "_stop", "_measurement", "_field", "home_id", "room_id", "sensor_id"])`+
			`|> aggregateWindow(every: 5m, fn: mean)`,
			d.DatabaseName,
			interval.StartString(),
			interval.EndString())
	}

	humanLabel := fmt.Sprintf(`InfluxDB (%s) aggregate/drop`, d.language)
	q := qi.(*bulkQuerygen.HTTPQuery)
	d.getHttpQuery(humanLabel, interval.StartString(), query, q)
}

func (d *InfluxIot) StandAloneFilter(qi bulkQuerygen.Query) {
	interval := d.AllInterval.RandWindow(8 * time.Hour)

	var query string
	if d.language == InfluxQL {
		query = fmt.Sprintf(`SELECT level FROM light_level_room WHERE time > '%s' AND time < '%s'`, interval.StartString(), interval.EndString())
	} else {
		query = fmt.Sprintf(`from(bucket: "%s") `+
			`|> range(start: %s, stop: %s) `+
			`|> filter(fn: (r) => r["_measurement"] == "light_level_room" and r["_field"] == "level") `+
			`|> yield()`,
			d.DatabaseName,
			interval.StartString(),
			interval.EndString(),
		)
	}

	humanLabel := fmt.Sprintf(`InfluxDB (%s) 8 hrs Room Light Level (Raw Data)`, d.language)
	q := qi.(*bulkQuerygen.HTTPQuery)
	d.getHttpQuery(humanLabel, "n/a", query, q)
}

func (d *InfluxIot) IotSortedPivot(qi bulkQuerygen.Query, timeInterval time.Duration) {
	interval := d.AllInterval.RandWindow(timeInterval)
	var query string
	if d.language == InfluxQL {
		query = fmt.Sprintf(`SELECT * FROM air_quality_room WHERE time > '%s' AND time < '%s'`,
			interval.StartString(),
			interval.EndString())
	} else {
		query = fmt.Sprintf(`from(bucket: "%s") `+
			`|> range(start: %s, stop: %s) `+
			`|> filter(fn: (r) => r._measurement == "air_quality_room") `+
			`|> pivot(rowKey:["_time"], columnKey:["_field"], valueColumn:"_value")`,
			d.DatabaseName,
			interval.StartString(),
			interval.EndString())
	}

	humanLabel := fmt.Sprintf(`InfluxDB (%s) Sorted Pivot`, d.language)
	q := qi.(*bulkQuerygen.HTTPQuery)
	d.getHttpQuery(humanLabel, interval.StartString(), query, q)
}
