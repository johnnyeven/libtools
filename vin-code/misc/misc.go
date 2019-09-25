package misc

// save VIN years
var vinYear map[rune][2]int

func init() {
	// init VIN years
	// VINs will not include I (i), O (o), Q (q), U, Z or the number 0
	vinYear = make(map[rune][2]int)
	yearStart := 1980
	for i := 'A'; i <= 'Z'; i++ {
		if 'I' == i || 'O' == i || 'Q' == i || 'U' == i || 'Z' == i {
			continue
		}
		year1 := yearStart
		year2 := yearStart + 30
		years := [...]int{year1, year2}
		vinYear[i] = years
		yearStart += 1
	}
	for i := '1'; i <= '9'; i++ {
		year1 := yearStart
		year2 := yearStart + 30
		years := [...]int{year1, year2}
		vinYear[i] = years
		yearStart += 1
	}
}
