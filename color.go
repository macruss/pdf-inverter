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

// cmykToInvertedRGB converts a CMYK color to its visual complement in RGB.
// Direct per-channel CMYK inversion is wrong (subtractive model); this goes through RGB.
func cmykToInvertedRGB(c, m, y, k float64) (r, g, b float64) {
	r = 1 - (1-c)*(1-k)
	g = 1 - (1-m)*(1-k)
	b = 1 - (1-y)*(1-k)
	return
}
