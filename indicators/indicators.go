/*
	import "github.com/thetruetrade/gotrade"

	Package indicators provides a range of technical trading indicators.
	All indicators follow the basic structure of:
		- receiving price data,	processing this price data and storing the transformed result.
		- maximum and minimum bounds of the transformed results are calculated automatically.
		- a lookback period indicating the lag between source data and the transformed result.
		- the source data bar from which the indicator is valid

 	Functions are provided for each indicator that provide indicator creation
 	for the following scenarios:

 	Online Usage
		- the data stream length is not known ahead of time, e.g. real time data streams
 	Offline Usage
		- the data stream length is known ahead of time, e.g. historical data streams

	Both scenarios provide the following indicator creation functions
 		* Indicator with default parameters
 		* Indicator with default parameters for attachment to a data stream
 		* Indicator with specified parameters
 		* Indicator with specified parameters for attachment to a data stream
 		* Indicator without storage with specified parameters
			- for use inside other indicators, has no storage of results which is instead
			- provided via a callback when it becomes available for use in the parent indicator.
*/
package indicators

import (
	"errors"
	"github.com/thetruetrade/gotrade"
	"math"
)

var (
	ErrSourceDataEmpty                      = errors.New("Source data is empty")
	ErrNotEnoughSourceDataForLookbackPeriod = errors.New("Source data does not contain enough data for the specfied lookback period")
	ErrLookbackPeriodMustBeGreaterThanZero  = errors.New("Lookback period must be greater than 0")
	ErrValueAvailableActionIsNil            = errors.New("ValueAvailableAction cannot be empty")

	// lookback minimum
	MinimumLookbackPeriod int = 0
	// lookback maximum
	MaximumLookbackPeriod int = 100000
)

type Indicator interface {
	// the source data bar number from which this indicator is valid, starts at bar 1.
	ValidFromBar() int
	// the lookback period, if applicable, the amount of lag the indicator displays with regards to the source data.
	GetLookbackPeriod() int
	// the length of the transformed data generated by the indicator.
	Length() int
}

type IndicatorWithTimePeriod interface {
	GetTimePeriod() int
}

type IndicatorWithFloatBounds interface {
	// the minimum bound of the transformed data generated by the indicator.
	MinValue() float64
	// the maximum bound of the transformed data generated by the indicator.
	MaxValue() float64
}

type IndicatorWithIntBounds interface {
	// the minimum bound of the transformed data generated by the indicator.
	MinValue() int64
	// the maximum bound of the transformed data generated by the indicator.
	MaxValue() int64
}

type baseFloatBounds struct {
	minValue float64
	maxValue float64
}

func newBaseFloatBounds() *baseFloatBounds {
	ind := baseFloatBounds{minValue: math.MaxFloat64, maxValue: math.SmallestNonzeroFloat64}
	return &ind
}

func (ind *baseFloatBounds) MinValue() float64 {
	return ind.minValue
}

func (ind *baseFloatBounds) MaxValue() float64 {
	return ind.maxValue
}

type baseIntBounds struct {
	minValue int64
	maxValue int64
}

func newBaseIntBounds() *baseIntBounds {
	ind := baseIntBounds{minValue: math.MaxInt64, maxValue: math.MinInt64}
	return &ind
}

func (ind *baseIntBounds) MinValue() int64 {
	return ind.minValue
}

func (ind *baseIntBounds) MaxValue() int64 {
	return ind.maxValue
}

type baseIndicator struct {
	validFromBar   int
	dataLength     int
	selectData     gotrade.DataSelectionFunc
	lookbackPeriod int
}

func newBaseIndicator(lookbackPeriod int) *baseIndicator {
	ind := baseIndicator{lookbackPeriod: lookbackPeriod, validFromBar: -1}
	return &ind
}

func (ind *baseIndicator) ValidFromBar() int {
	return ind.validFromBar
}

func (ind *baseIndicator) GetLookbackPeriod() int {
	return ind.lookbackPeriod
}

func (ind *baseIndicator) Length() int {
	return ind.dataLength
}

type baseIndicatorWithTimePeriod struct {
	timePeriod int
}

type baseIndicatorWithFloatBounds struct {
	*baseIndicator
	*baseFloatBounds
}

func newBaseIndicatorWithFloatBounds(lookbackPeriod int) *baseIndicatorWithFloatBounds {
	ind := baseIndicatorWithFloatBounds{
		baseIndicator:   newBaseIndicator(lookbackPeriod),
		baseFloatBounds: newBaseFloatBounds()}
	return &ind
}

type baseIndicatorWithIntBounds struct {
	*baseIndicator
	*baseIntBounds
}

func newBaseIndicatorWithIntBounds(lookbackPeriod int) *baseIndicatorWithIntBounds {
	ind := baseIndicatorWithIntBounds{
		baseIndicator: newBaseIndicator(lookbackPeriod),
		baseIntBounds: newBaseIntBounds()}
	return &ind
}

func newBaseIndicatorWithTimePeriod(timePeriod int) *baseIndicatorWithTimePeriod {
	ind := baseIndicatorWithTimePeriod{timePeriod: timePeriod}
	return &ind
}

func (ind *baseIndicatorWithTimePeriod) GetTimePeriod() int {
	return ind.timePeriod
}

type ValueAvailableActionFloat func(dataItem float64, streamBarIndex int)
type ValueAvailableActionInt func(dataItem int64, streamBarIndex int)
type ValueAvailableActionDOHLCV func(dataItem gotrade.DOHLCV, streamBarIndex int)
type ValueAvailableActionBollinger func(dataItemUpperBand float64, dataItemMiddleBand float64, dataItemLowerBand float64, streamBarIndex int)
type ValueAvailableActionMACD func(dataItemMACD float64, dataItemSignal float64, dataItemHistogram float64, streamBarIndex int)
type ValueAvailableActionAroon func(dataItemAroonUp float64, dataItemAroonDown float64, streamBarIndex int)
type ValueAvailableActionStoch func(dataItemK float64, dataItemD float64, streamBarIndex int)
type ValueAvailableActionLinearReg func(dataItem float64, slope float64, intercept float64, streamBarIndex int)
