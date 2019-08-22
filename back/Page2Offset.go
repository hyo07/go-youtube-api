package back

import "strconv"

func Page2Offset(page string) (offset int) {
	intPage, err := strconv.Atoi(page)
	if err != nil {
		intPage = 1
	}
	offset = (intPage - 1) * 10
	return
}
