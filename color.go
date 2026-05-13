package main

func invertValue(v float64) float64 {
	return 1.0 - v
}

func invertRGB(r, g, b float64) (float64, float64, float64) {
	return invertValue(r), invertValue(g), invertValue(b)
}

func invertGray(g float64) float64 {
	return invertValue(g)
}

func invertCMYK(c, m, y, k float64) (float64, float64, float64, float64) {
	return invertValue(c), invertValue(m), invertValue(y), invertValue(k)
}
