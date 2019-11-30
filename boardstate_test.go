package wqLogic

import "testing"

func TestPosiGroup(t *testing.T) {
	bs := new(BoardState)
	bs.init()
	bs.setupUnput()

	g := bs.newGroup(White)

	g._add(0)

	if bs.groups[0].memb.Len() != 1 ||
		bs.groups[0].surr.Len() != 2 {
		t.Errorf("posiGroup memb: %v surr: %v", bs.groups[0].memb, bs.groups[0].surr)
	}
}

func TestSetupUnput(t *testing.T) {
	bs := new(BoardState)
	bs.init()
	bs.setupUnput()

	if bs.groups[0].memb.Len() != 361 ||
		bs.groups[0].surr.Len() != 0 {
		t.Errorf("setupUnput memb: %v surr: %v", bs.groups[0].memb, bs.groups[0].surr)
	}

	if bs.groups[360].memb.Len() != 361 ||
		bs.groups[360].surr.Len() != 0 {
		t.Errorf("setupUnput memb: %v surr: %v", bs.groups[360].memb, bs.groups[360].surr)
	}
}

func TestNewStone(t *testing.T) {

	bs := new(BoardState)
	bs.init()

	r := bs.NewStone(0, White)

	if !r {
		t.Errorf("NewStone() failed")
	}

	if bs.movesLen != 1 {
		t.Errorf("movesLen wrong")
	}

	if bs.preKO != -1 {
		t.Errorf("preKO wrong")
	}

	if bs.groups[0].color != White {
		t.Errorf("posi 0 group color wrong")
	}

	if bs.groups[0].memb.Len() != 1 {
		t.Errorf("posi 0 group memb wrong %v", bs.groups[0].memb)
	}

	if bs.groups[0].surr.Len() != 2 {
		t.Errorf("posi 0 group surr wrong: %v != 2", bs.groups[0].surr) //.Len())
	}

	if bs.groups[360].color != Unput {
		t.Errorf("posi 360 group color wrong")
	}

	if bs.groups[360].memb.Len() != 360 {
		t.Errorf("posi 360 group memb wrong: %v != 360", bs.groups[360].memb.Len())
	}

	if bs.groups[360].surr.Len() != 1 {
		t.Errorf("posi 360 group surr wrong: %v != 1", bs.groups[360].surr) //.Len())
	}

	r = bs.NewStone(1, Black)
	r = bs.NewStone(19, Black)
	if !r {
		t.Errorf("posi 1 Black: [0].color: %v [0].memb: %v surr: %v [1].memb: %v surr: %v [2].memb: %v surr: %v [19].memb: %v surr: %v != 1",
			bs.groups[0].color, bs.groups[0].memb, bs.groups[0].surr,
			bs.groups[1].memb, bs.groups[1].surr,
			bs.groups[2].memb, bs.groups[2].surr,
			bs.groups[19].memb, bs.groups[19].surr,
		)
	}

	r = bs.NewStone(19, White)
	if r {
		t.Errorf("cannot put on a stone")
	}

	r = bs.NewStone(20, White)
	r = bs.NewStone(38, White)
	r = bs.NewStone(0, White)
	if bs.groups[19].color != Unput {
		t.Errorf("19 shall be dead")
	}
	if bs.preKO == -1 {
		t.Errorf("preKO %v != -1", bs.preKO)
	}

	r = bs.NewStone(19, Black)
	if r {
		t.Errorf("KO not allowed")
	}

	r = bs.NewStone(2, Black)
	if bs.preKO != -1 {
		t.Errorf("preKO %v != -1", bs.preKO)
	}

	r = bs.NewStone(19, Black)
	if !r {
		t.Errorf("not KO")
	}
	if bs.preKO == -1 {
		t.Errorf("19 preKO %v != -1", bs.preKO)
	}

	r = bs.NewStone(0, White)
	if r {
		t.Errorf("KO! not allowed")
	}

	r = bs.NewStone(3, White)
	r = bs.NewStone(21, White)
	r = bs.NewStone(0, Black)
	if r {
		t.Error("out of liberty")
	}

	r = bs.NewStone(0, White)
	if !r || bs.groups[1].color != Unput ||
		bs.groups[2].color != Unput ||
		bs.groups[19].color != Unput || bs.preKO != -1 {
		t.Error("1 2 19 should be dead")
	}

	/*
		t.Errorf("movesLen: %v  preKO: %v", bs.movesLen, bs.preKO)

		for p := 0; p <= 360; p++ {
			if bs.groups[p].memb.Len() < 100 {//|| p == 360 {
				t.Errorf("[%v] memb: %v surr: %v color: %v", p, bs.groups[p].memb, bs.groups[p].surr, bs.groups[p].color)
			}
		}
	*/
}
