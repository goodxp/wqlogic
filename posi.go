package wqLogic

//x, y: 0..18
//posi: 0..360

func xy2posi(x, y int) int {
	return y*19 + x
}

func posi2xy(p int) (x int, y int) {
	return p % 19, (p - p%19) / 19
}

func posiSurr(p int) []int {
	s := []int{}
	if left := p - 1; p%19-1 >= 0 {
		s = append(s, left)
	}
	if right := p + 1; p%19+1 <= 18 {
		s = append(s, right)
	}
	if above := p - 19; above >= 0 {
		s = append(s, above)
	}
	if below := p + 19; below <= 360 {
		s = append(s, below)
	}
	return s
}

const (
	Black PosiColor = 1
	White PosiColor = -1
	Unput PosiColor = 0
)

type PosiColor int
