package wqLogic

import (
	"actTree"
	"wqSGF"
)

type Move struct {
	Posi  int
	Color PosiColor

	MoveNum int
	ForkNum int

	node *actTree.Node
}

func NewMove(posi int, color PosiColor) *Move {
	return &Move{Posi: posi, Color: color, MoveNum: 1}
}

type Kifu struct {
	*GameInfo
	FirstMove *Move

	tree  *actTree.Tree
	moves map[*actTree.Node]*Move

	maxForkNum int
}

func NewKifu() *Kifu {
	k := new(Kifu)
	k.Init()
	return k
}

func (k *Kifu) Init() {
	k.tree = new(actTree.Tree)
	k.moves = make(map[*actTree.Node]*Move)
	k.GameInfo = new(GameInfo)
}

func (k *Kifu) AddMove(m *Move, base *Move) *Move {
	k.updateBoardState(base)
	if r := k.NewStone(m.Posi, m.Color, m); !r {
		return nil // against Go rules
	}
	k.byForkNum = m.ForkNum
	k.byMoveNum = m.MoveNum

	var baseNode *actTree.Node
	if base != nil {
		baseNode = base.node
	} else {
		if k.FirstMove != nil {
			baseNode, _ = k.FirstMove.node.Prev()
		}
	}

	m.node = k.tree.Add(new(wqSGF.Node), baseNode)
	if bm := k.moves[baseNode]; bm != nil {
		if kid, _ := baseNode.Next(); kid == m.node {
			m.ForkNum = bm.ForkNum
			m.MoveNum = bm.MoveNum + 1
		} else {
			k.maxForkNum++
			m.ForkNum = k.maxForkNum
			m.MoveNum = 1
		}
	}
	k.moves[m.node] = m
	if k.FirstMove == nil {
		k.FirstMove = m
	}
	return m
}

func (k *Kifu) UndoMove(m *Move) {
	if m == nil || m.node == nil {
		return
	}
	k.moves[m.node] = nil
	k.tree.Remove(m.node)
	m.node = nil
}

func (k *Kifu) GetMove(forkNum int, moveNum int) *Move {
	for _, m := range k.moves {
		if m.ForkNum == forkNum && m.MoveNum == moveNum {
			return m
		}
	}
	return nil
}

func (k *Kifu) PrevMove(m *Move) *Move {
	if m != nil && m.node != nil {
		p, _ := m.node.Prev()
		return k.moves[p]
	}
	return nil
}

func (k *Kifu) NextMove(m *Move) *Move {
	if m != nil && m.node != nil {
		kid, _ := m.node.Next()
		return k.moves[kid]
	}
	return nil
}

func (k *Kifu) GetForks(m *Move) []*Move {
	if m != nil && m.node != nil {
		if kid, _ := m.node.Next(); kid != nil {
			var ret []*Move
			for _, sib := kid.Next(); sib != nil; _, sib = sib.Next() {
				ret = append(ret, k.moves[sib])
			}
			return ret
		}
	}
	return nil
}

func (k *Kifu) LoadSGF(sgf string) {
	k.tree = wqSGF.Parse(sgf)
	k.loadSgfRoot()
	actTree.WalkThrough(k.tree.Root, func(n *actTree.Node) bool {
		var m *Move
		for _, prop := range n.Value.(*wqSGF.Node).Props {
			m = nil
			if prop.Id == "B" {
				x, y := wqSGF.V2P(wqSGF.Val2V(prop.Vals[0]))
				m = NewMove(xy2posi(x, y), Black)
				break
			}
			if prop.Id == "W" {
				x, y := wqSGF.V2P(wqSGF.Val2V(prop.Vals[0]))
				m = NewMove(xy2posi(x, y), White)
				break
			}
		}
		if m != nil {
			p, s := n.Prev()
			if s != nil && k.moves[s] != nil {
				k.maxForkNum++
				m.ForkNum = k.maxForkNum
				m.MoveNum = 1
			} else if p != nil && k.moves[p] != nil {
				m.ForkNum = k.moves[p].ForkNum
				m.MoveNum = k.moves[p].MoveNum + 1
			}
			m.node = n
			k.moves[n] = m
			if k.FirstMove == nil {
				k.FirstMove = m
			}
		}
		return false
	}, nil)
}

func (k *Kifu) ToSGF() string {
	for n, m := range k.moves {
		var pp *wqSGF.Prop
		for _, prop := range n.Value.(*wqSGF.Node).Props {
			pp = nil
			if prop.Id == "B" || prop.Id == "W" {
				pp = &prop
				break
			}
		}
		if pp == nil {
			pp = new(wqSGF.Prop)
			if m.Color == Black {
				pp.Id = "B"
			} else {
				pp.Id = "W"
			}
			x, y := posi2xy(m.Posi)
			pp.Vals = append(pp.Vals, wqSGF.V2Val(wqSGF.P2V(x, y)))
			n.Value.(*wqSGF.Node).Props = append(n.Value.(*wqSGF.Node).Props, *pp)
		}
	}
	k.saveSgfRoot()
	return wqSGF.ToSGF(k.tree)
}
