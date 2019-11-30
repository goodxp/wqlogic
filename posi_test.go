package wqLogic

import "testing"

func TestPosiXY(t *testing.T) {
	p := xy2posi(0, 0)
	x, y := posi2xy(p)
	if x != 0 || y != 0 {
		t.Errorf("xy<=>posi: %v == 0x0", p)
	}

	p = xy2posi(18, 18)
	x, y = posi2xy(p)
	if x != 18 || y != 18 {
		t.Errorf("xy<=>posi: %v == 18x18", p)
	}

	p = xy2posi(17, 18)
	x, y = posi2xy(p)
	if x != 17 || y != 18 {
		t.Errorf("xy<=>posi: %v == 17x18", p)
	}

	p = xy2posi(18, 17)
	x, y = posi2xy(p)
	if x != 18 || y != 17 {
		t.Errorf("xy<=>posi: %v == 18x17", p)
	}

	p = xy2posi(18, 0)
	x, y = posi2xy(p)
	if x != 18 || y != 0 {
		t.Errorf("xy<=>posi: %v == 18x0", p)
	}

	p = xy2posi(0, 18)
	x, y = posi2xy(p)
	if x != 0 || y != 18 {
		t.Errorf("xy<=>posi: %v == 0x18", p)
	}
}

func TestPosiSurr(t *testing.T) {
	p := 0
	s := posiSurr(p)
	if len(s) != 2 {
		t.Errorf("posiSurr: %v == 2", s)
	}

	p = 360
	s = posiSurr(p)
	if len(s) != 2 {
		t.Errorf("posiSurr: %v == 2", s)
	}

	p = 18
	s = posiSurr(p)
	if len(s) != 2 {
		t.Errorf("posiSurr: %v == 2", s)
	}

	p = 342
	s = posiSurr(p)
	if len(s) != 2 {
		t.Errorf("posiSurr: %v == 2", s)
	}

	p = 17
	s = posiSurr(p)
	if len(s) != 3 {
		t.Errorf("posiSurr: %v == 3", s)
	}

	p = 341
	s = posiSurr(p)
	if len(s) != 3 {
		t.Errorf("posiSurr: %v == 3", s)
	}

	p = 343
	s = posiSurr(p)
	if len(s) != 3 {
		t.Errorf("posiSurr: %v == 3", s)
	}

	p = 19
	s = posiSurr(p)
	if len(s) != 3 {
		t.Errorf("posiSurr: %v == 3", s)
	}

	p = 20
	s = posiSurr(p)
	if len(s) != 4 {
		t.Errorf("posiSurr: %v == 4", s)
	}

	p = 340
	s = posiSurr(p)
	if len(s) != 4 {
		t.Errorf("posiSurr: %v == 4", s)
	}
}
