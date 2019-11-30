package wqLogic

import "testing"

func TestGameInfo(t *testing.T) {
	k := NewKifu()
	k.Komi = 7.5
	k.Handicap = 2
	k.AddMove(NewMove(xy2posi(3, 3), Black), nil)

	if k.Komi != 7.5 || k.Handicap != 2 {
		t.Errorf("GameInfo error 1 %v %v", k.Komi, k.Handicap)
	}

	s := k.ToSGF()
	k1 := NewKifu()
	k1.LoadSGF(s)

	if k1.Komi != 7.5 || k1.Handicap != 2 {
		t.Errorf("GameInfo error 2 %v %v", k1.Komi, k1.Handicap)
	}

	k1.Komi = 6.5
	k1.Handicap = 7
	k2 := NewKifu()
	k2.LoadSGF(k1.ToSGF())

	if k2.Komi != 6.5 || k2.Handicap != 7 {
		t.Errorf("GameInfo error 3 %v %v", k2.Komi, k2.Handicap)
	}
}

func TestState(t *testing.T) {
	k := NewKifu()
	var move *Move
	for i := 0; i < 19; i++ {
		move = k.AddMove(NewMove(xy2posi(i, 9), Black), move)
		move = k.AddMove(NewMove(xy2posi(i, 10), White), move)
	}

	k.TurnColor(move, xy2posi(0, 8))
	color, pinColor, forkNum, moveNum := k.GetPosiState(move, xy2posi(0, 8))
	if moveNum != 0 || color != Unput || pinColor != Black || forkNum != 0 {
		t.Errorf("TurnColor() error 1: %v, %v, %v, %v", color, pinColor, forkNum, moveNum)
	}
	color, pinColor, _, moveNum = k.GetPosiState(move, xy2posi(0, 10))
	if moveNum != 2 || color != White || pinColor != White {
		t.Errorf("GetPosiState() error 1")
	}

	result, moyo, deadLen := k.Forecast(move)
	if result != 19 || moyo[0] != 9*19 || deadLen[1] != 0 {
		t.Errorf("Forecast() error 1")
	}

	k.TurnColor(move, xy2posi(9, 9))
	result, moyo, deadLen = k.Forecast(move)
	if result != -361 || moyo[0] != 0 || deadLen[0] != 0 {
		t.Errorf("Forecast() error 2: %v, %v, %v", result, moyo, deadLen)
	}

	move = k.AddMove(NewMove(xy2posi(0, 0), White), move)
	move = k.AddMove(NewMove(xy2posi(1, 0), Black), move)
	move = k.AddMove(NewMove(xy2posi(0, 1), Black), move)

	color, pinColor, forkNum, moveNum = k.GetPosiState(move, xy2posi(0, 0))
	if pinColor != Black || moveNum != 0 || forkNum != 0 || color != Unput {
		t.Errorf("GetPosiState() error 2:  %v, %v, %v %v", color, pinColor, forkNum, moveNum)
	}

	result, moyo, deadLen = k.Forecast(move)
	if result != 18 || deadLen[1] != 1 || moyo[0] != 9*19-1 {
		t.Errorf("Forecast() error 3  %v, %v, %v", result, moyo, deadLen)
	}
}
