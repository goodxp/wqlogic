package wqLogic

import "testing"

func TestSGF(t *testing.T) {
	k := NewKifu()
	r := k.AddMove(NewMove(xy2posi(3, 3), Black), nil)
	k.AddMove(NewMove(xy2posi(15, 15), White), r)
	s := k.ToSGF()

	k1 := NewKifu()
	k1.LoadSGF(s)
	m1 := k1.GetMove(0, 1)
	m2 := k1.NextMove(m1)

	if m1.MoveNum != 1 || m1.Posi != xy2posi(3, 3) || m1.Color != Black ||
		m2.MoveNum != 2 || m2.Posi != xy2posi(15, 15) || m2.Color != White {

		t.Errorf("simple sgf save&load failed")
	}
}

func TestFork(t *testing.T) {
	moves := []*Move{
		NewMove(xy2posi(1, 1), Black),
		NewMove(xy2posi(2, 2), White),
		NewMove(xy2posi(3, 3), Black),
		NewMove(xy2posi(4, 4), White),
		NewMove(xy2posi(5, 5), Black),
		NewMove(xy2posi(6, 6), White),
		NewMove(xy2posi(7, 7), Black),
		NewMove(xy2posi(8, 8), Black),
	}

	k := NewKifu()

	m0 := k.AddMove(moves[0], nil)
	m1 := k.AddMove(moves[1], m0)
	m2 := k.AddMove(moves[2], m0)

	forks := k.GetForks(m0)
	if len(forks) != 1 || forks[0] != m2 {
		t.Errorf("GetForks() error: 1")
	}

	m3 := k.AddMove(moves[3], m1)
	m4 := k.AddMove(moves[4], m1)
	m5 := k.AddMove(moves[5], m2)
	m6 := k.AddMove(moves[6], m2)
	m7 := k.AddMove(moves[7], m2)

	forks = k.GetForks(m1)
	if len(forks) != 1 || forks[0] != m4 {
		t.Errorf("GetForks() error: 2")
	}

	forks = k.GetForks(m2)
	if len(forks) != 2 || forks[0] != m6 || forks[1] != m7 {
		t.Errorf("GetForks() error: 3")
	}

	if m4.ForkNum != 2 || m6.ForkNum != 3 || m7.ForkNum != 4 {
		t.Errorf("ForkNum error")
	}

	if m3.MoveNum != 3 || m7.MoveNum != 1 {
		t.Errorf("MoveNum error")
	}

	if k.PrevMove(m5) != m2 || k.PrevMove(m2) != m0 {
		t.Errorf("PrevMove() error")
	}

	k.UndoMove(m2)

	if k.NextMove(m0) != m1 || len(k.GetForks(m0)) != 0 {
		t.Errorf("UndoMove() error")
	}
}

func TestSGFRoot(t *testing.T) {
	k := NewKifu()
	s := k.ToSGF()
	moves := []*Move{
		NewMove(xy2posi(1, 1), Black),
		NewMove(xy2posi(2, 2), White),
		NewMove(xy2posi(3, 3), Black),
		NewMove(xy2posi(4, 4), White),
	}
	m0 := k.AddMove(moves[0], nil)
	m1 := k.AddMove(moves[1], m0)
	m2 := k.AddMove(moves[2], m0)

	forks := k.GetForks(m0)
	if len(forks) != 1 || forks[0] != m2 || k.NextMove(m0) != m1 {
		t.Errorf("GetForks() error: 4")
	}

	s = k.ToSGF()
	k1 := NewKifu()
	k1.LoadSGF(s)

	nm0 := k1.FirstMove
	if nm0.Posi != m0.Posi || k1.PrevMove(nm0) != nil {
		t.Errorf("SGF Root error: 1")
	}

	nm2 := k1.GetMove(1, 1)
	if nm2.Posi != m2.Posi {
		t.Errorf("GetMove error: 1")
	}

	m3 := k1.AddMove(moves[3], nm2)
	if k1.PrevMove(nm2) != nm0 {
		t.Errorf("SGF Root error: 2")
	}
	if k1.NextMove(nm2) != m3 {
		t.Errorf("SGF Root error: 3")
	}
}
