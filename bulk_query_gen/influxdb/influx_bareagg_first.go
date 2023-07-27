package influxdb

import (
	"time"

	bulkQuerygen "github.com/antondavidsen/influxdb-comparisons/bulk_query_gen"
)

func NewInfluxQLBareAggregateFirst(dbConfig bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, _ time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	return NewInfluxBareAggregateQuery(First, InfluxQL, dbConfig, queriesFullRange, scaleVar)
}

func NewFluxBareAggregateFirst(dbConfig bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, _ time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	return NewInfluxBareAggregateQuery(First, Flux, dbConfig, queriesFullRange, scaleVar)
}
