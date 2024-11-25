package rebalance_svc

func ema(values []float64, period int) (res []float64) {
	alpha := 2.0 / float64(period+1)
	emaValues := make([]float64, len(values))
	emaValues[0] = values[0] // Start with the first value

	for i := 1; i < len(values); i++ {
		emaValues[i] = alpha*values[i] + (1-alpha)*emaValues[i-1]
	}

	res = make([]float64, len(values))
	for i := range emaValues {
		if i >= (period - 1) { // First 32 values will be NaN since there's no full period
			res[i] = emaValues[i]
		} else {
			res[i] = 0 // Represent NaN as nil
		}
	}

	return res
}
