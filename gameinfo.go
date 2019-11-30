package wqLogic

import (
	"actTree"
	"wqSGF"
)

type GameInfo struct {
	Komi     float32 // KM[7.5]
	Handicap int     // HA[2]AB[aa][bb]AW[cc]PL[W]  ???handicap stones

	*BoardState
	byForkNum int
	byMoveNum int
}

func (k *Kifu) GetPosiState(m *Move, posi int) (PosiColor, PosiColor, int, int) {
	k.updateBoardState(m)
	if k.groups[posi].color != Unput {
		return k.groups[posi].color, k.groups[posi].pinColor, k.ids[posi].(*Move).ForkNum, k.ids[posi].(*Move).MoveNum
	}
	return Unput, k.groups[posi].pinColor, 0, 0
}

func (k *Kifu) TurnColor(m *Move, posi int) {
	k.updateBoardState(m)
	k.BoardState.TurnColor(posi)
	k.CalcMoyo()
}

func (k *Kifu) Forecast(m *Move) (float32, [2]int, [2]int) {
	k.updateBoardState(m)
	result := float32(k.moyo[0]-k.moyo[1]) - k.Komi
	return result, k.moyo, k.deadLen
}

func (k *Kifu) updateBoardState(m *Move) {
	if m != nil && k.byForkNum == m.ForkNum && k.byMoveNum == m.MoveNum {
		return
	}
	var stack []*Move
	for pm := m; pm != nil; pm = k.PrevMove(pm) {
		stack = append(stack, pm)
	}
	k.BoardState = new(BoardState)
	k.BoardState.init()
	for n := len(stack) - 1; n >= 0; n-- {
		pm := stack[n]
		k.NewStone(pm.Posi, pm.Color, pm)
	}
	if m != nil {
		k.byForkNum = m.ForkNum
		k.byMoveNum = m.MoveNum
	}
	k.CalcMoyo()
}

func (k *Kifu) saveSgfRoot() {
	if k.tree.Root == nil {
		return
	}
	komi := -1
	handicap := -1
	if !isSgfRoot(k.tree.Root.Value.(*wqSGF.Node)) {
		t := new(actTree.Tree)
		t.AddNode(k.tree.Root, t.Add(newSgfRoot(), nil))
		k.tree = t
	} else {
		for i, prop := range k.tree.Root.Value.(*wqSGF.Node).Props {
			switch prop.Id {
			case "KM":
				komi = i
			case "HA":
				handicap = i
			}
		}
	}

	props := &k.tree.Root.Value.(*wqSGF.Node).Props
	if k.Komi != 0 {
		if komi != -1 {
			(*props)[komi].Vals[0] = wqSGF.V2Val(wqSGF.R2V(k.Komi))
		} else {
			*props = append(*props, wqSGF.Prop{Id: "KM", Vals: []string{wqSGF.V2Val(wqSGF.R2V(k.Komi))}})
		}
	}
	if k.Handicap != 0 {
		if handicap != -1 {
			(*props)[handicap].Vals[0] = wqSGF.V2Val(wqSGF.I2V(k.Handicap))
		} else {
			*props = append(*props, wqSGF.Prop{Id: "HA", Vals: []string{wqSGF.V2Val(wqSGF.I2V(k.Handicap))}})
		}
	}
}

func (k *Kifu) loadSgfRoot() {
	if k.tree.Root == nil {
		return
	}
	for _, prop := range k.tree.Root.Value.(*wqSGF.Node).Props {
		switch prop.Id {
		case "KM":
			k.Komi = wqSGF.V2R(wqSGF.Val2V(prop.Vals[0]))
		case "HA":
			k.Handicap = wqSGF.V2I(wqSGF.Val2V(prop.Vals[0]))
		}
	}
}

func newSgfRoot() *wqSGF.Node {
	return &wqSGF.Node{
		Props: []wqSGF.Prop{
			wqSGF.Prop{Id: "FF", Vals: []string{"[4]"}},
			wqSGF.Prop{Id: "GM", Vals: []string{"[1]"}},
			wqSGF.Prop{Id: "SZ", Vals: []string{"[19]"}},
			wqSGF.Prop{Id: "CA", Vals: []string{"[UTF-8]"}},
			// wqSGF.Prop{Id: "GX", Vals: []string{"[1]"}},
		},
	}
}

func isSgfRoot(root *wqSGF.Node) bool {
	if root == nil {
		return false
	}
	var propFF, propGM, propSZ, propCA bool
	for _, prop := range root.Props {
		if prop.Id == "FF" {
			propFF = true
		}
		if prop.Id == "GM" {
			propGM = true
		}
		if prop.Id == "SZ" {
			propSZ = true
		}
		if prop.Id == "CA" {
			propCA = true
		}
	}
	if !propFF && !propGM && !propSZ && !propCA {
		return false
	}
	return true
}
