package influxdb

import (
	"time"

	bulkQuerygen "github.com/antondavidsen/influxdb-comparisons/bulk_query_gen"
)

// InfluxQL query for "Group Window" on the standard cardinality IoT dataset
func NewInfluxQLGroupWindowTransposeMin(dbConfig bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	return NewInfluxGroupWindowTransposeQuery(Min, LowCardinality, InfluxQL, dbConfig, queriesFullRange, queryInterval, scaleVar)
}

// Flux Query query for "Group Window" on the standard cardinality IoT dataset
func NewFluxGroupWindowTransposeMin(dbConfig bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	return NewInfluxGroupWindowTransposeQuery(Min, LowCardinality, Flux, dbConfig, queriesFullRange, queryInterval, scaleVar)
}

// InfluxQL query for "Group Window" on the high cardinality Metaquery dataset
func NewInfluxQLGroupWindowTransposeMinCardinality(dbConfig bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	return NewInfluxGroupWindowTransposeQuery(Min, HighCardinality, InfluxQL, dbConfig, queriesFullRange, queryInterval, scaleVar)
}

// Flux query for "Group Window" on the high cardinality Metaquery dataset
func NewFluxGroupWindowTransposeMinCardinality(dbConfig bulkQuerygen.DatabaseConfig, queriesFullRange bulkQuerygen.TimeInterval, queryInterval time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	return NewInfluxGroupWindowTransposeQuery(Min, HighCardinality, Flux, dbConfig, queriesFullRange, queryInterval, scaleVar)
}
