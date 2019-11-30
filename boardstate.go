package wqLogic

type PosiGroup struct {
	memb     *Set
	surr     *Set
	color    PosiColor
	board    *BoardState
	movesLen int // last updated at

	pinColor PosiColor // appointed by players
	//	estColor PosiColor // estimated
}

func (pg *PosiGroup) init(b *BoardState, color PosiColor) {
	pg.board = b
	pg.movesLen = b.movesLen
	pg.memb = NewSet()
	pg.surr = NewSet()
	pg.color = color
}

func (pg *PosiGroup) add(posi int, color PosiColor) bool {
	if pg.color != color ||
		(pg.memb.Len() > 0 && !pg.surr.Has(posi)) {
		return false
	}
	pg._add(posi)
	pg.board.groupSurr(posi)
	return true // added!
}

func (pg *PosiGroup) _add(posi int) {
	pg.memb.Insert(posi)
	pg.surr.Remove(posi)
	for _, s := range posiSurr(posi) {
		if !pg.memb.Has(s) {
			pg.surr.Insert(s)
		}
	}
	pg.board.groups[posi] = pg
}

func (pg *PosiGroup) dead() bool {
	r := true
	pg.surr.Do(func(s interface{}) bool {
		if pg.board.groups[s.(int)].color == Unput {
			r = false
			return true
		}
		return false
	})
	return r
}

type BoardState struct {
	groups   map[int]*PosiGroup
	movesLen int
	preKO    int // posi that just made a 1-1 kill

	deadLen [2]int // prisoners: [0] for black, [1] for white
	moyo    [2]int // territory: [0] for black, [1] for white

	ids map[int]interface{} //*Move
}

func (bs *BoardState) init() {
	bs.groups = make(map[int]*PosiGroup)
	bs.movesLen = 0
	bs.preKO = -1
	bs.ids = make(map[int]interface{})
}

func (bs *BoardState) setupUnput() {
	g := bs.newGroup(Unput)
	for p := 0; p <= 360; p++ {
		g._add(p)
	}
}

func (bs *BoardState) newGroup(color PosiColor) *PosiGroup {
	g := new(PosiGroup)
	g.init(bs, color)
	return g
}

func (bs *BoardState) groupSurr(posi int) {
	g := bs.groups[posi]
	for _, s := range posiSurr(posi) {
		if g.movesLen > bs.groups[s].movesLen &&
			g.add(s, bs.groups[s].color) {
			for _, ss := range posiSurr(s) {
				bs.groupSurr(ss)
			}
		}
	}
}

// return: isValidMove bool, deadCount int
func (bs *BoardState) libertyCheck(posi int) (bool, int) {
	g := bs.groups[posi]
	gDead := g.dead()
	sDead := 0
	ko := -1
	for _, s := range posiSurr(posi) {
		sg := bs.groups[s]
		if sg.color == -g.color && sg.dead() {
			sg.color = Unput
			ko = s
			sDead += sg.memb.Len()
		}
	}

	if !gDead {
		return true, sDead // valid state
	}
	if sDead == 0 {
		return false, 0
	}
	if g.memb.Len() > 1 || sDead > 1 {
		return true, sDead
	}
	if bs.preKO == -1 {
		bs.preKO = posi
		return true, sDead
	}
	if bs.preKO == ko {
		return false, 0 //KO!
	}

	return true, sDead
}

func (bs *BoardState) copy(from *BoardState) {
	gs := NewSet()
	for _, fg := range from.groups {
		gs.Insert(fg)
	}
	gs.Do(func(ug interface{}) bool {
		fg := ug.(*PosiGroup)
		g := bs.newGroup(fg.color)
		g.movesLen = fg.movesLen
		g.pinColor = fg.pinColor
		fg.memb.Do(func(p interface{}) bool {
			g.memb.Insert(p.(int))
			bs.groups[p.(int)] = g
			return false
		})
		fg.surr.Do(func(p interface{}) bool {
			g.surr.Insert(p.(int))
			return false
		})
		return false
	})
	bs.movesLen = from.movesLen
	bs.preKO = from.preKO
	bs.deadLen[0] = from.deadLen[0]
	bs.deadLen[1] = from.deadLen[1]
	bs.moyo[0] = from.moyo[0]
	bs.moyo[1] = from.moyo[1]
	for p, id := range from.ids {
		bs.ids[p] = id
	}
}

func (bs *BoardState) NewStone(posi int, color PosiColor, id ...interface{}) bool {
	if bs.movesLen == 0 {
		bs.setupUnput()
	}
	if bs.movesLen > 0 && bs.groups[posi].color != Unput {
		return false
	}

	backup := new(BoardState)
	backup.init()
	backup.copy(bs)

	bs.movesLen++
	g := bs.newGroup(color)
	g.add(posi, color)
	for _, s := range posiSurr(posi) {
		if bs.groups[s].movesLen < bs.movesLen {
			sg := bs.newGroup(bs.groups[s].color)
			sg.add(s, bs.groups[s].color)
		}
	}

	ok, dead := bs.libertyCheck(posi)
	if ok {
		if bs.preKO == backup.preKO {
			bs.preKO = -1
		}
		if len(id) == 1 {
			bs.ids[posi] = id[0]
		}
		if color == Black {
			bs.deadLen[1] += dead
		} else {
			bs.deadLen[0] += dead
		}
	} else {
		bs.copy(backup)
	}
	return ok
}

func (bs *BoardState) TurnColor(posi int) {
	g := bs.groups[posi]
	if g.color == Unput {
		return
	}
	if g.pinColor == Unput {
		g.pinColor = -g.color
	} else {
		g.pinColor = -g.pinColor
	}
}

func (bs *BoardState) CalcMoyo() {
	bs.moyo[0] = 0
	bs.moyo[1] = 0
	bwgs := NewSet()
	gs := NewSet()
	for _, g := range bs.groups {
		if g.color == Unput {
			gs.Insert(g)
		} else {
			bwgs.Insert(g)
		}
	}
	bwgs.Do(func(v interface{}) bool {
		g := v.(*PosiGroup)
		if g.pinColor == Unput {
			g.pinColor = g.color
		}
		if g.color == -g.pinColor {
			if g.color == Black {
				bs.moyo[1] += g.memb.Len() * 2
			} else {
				bs.moyo[0] += g.memb.Len() * 2
			}
		}
		return false
	})
	gs.Do(func(v interface{}) bool {
		g := v.(*PosiGroup)
		surrColor := Unput
		g.surr.Do(func(p interface{}) bool {
			sg := bs.groups[p.(int)]
			if surrColor == Unput {
				surrColor = sg.pinColor
				return false
			}
			if surrColor == sg.pinColor {
				return false
			}
			surrColor = Unput
			return true
		})
		if surrColor == Unput {
			return false
		}
		g.pinColor = surrColor
		if surrColor == Black {
			bs.moyo[0] += g.memb.Len()
		} else {
			bs.moyo[1] += g.memb.Len()
		}
		return false
	})
	bs.moyo[0] += bs.deadLen[1]
	bs.moyo[1] += bs.deadLen[0]
}
