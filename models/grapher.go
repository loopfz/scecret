package models

import "github.com/go-gorp/gorp"

type LocGraph struct {
	ID     int64        `json:"id"`
	Name   string       `json:"name"`
	Hidden bool         `json:"hidden"`
	Cards  []*CardGraph `json:"cards,omitempty"`
}

type CardGraph struct {
	ID                    int64   `json:"id"`
	Description           string  `json:"description"`
	Reveals               []int64 `json:"reveals,omitempty"`
	UnlockStateTokens     []int64 `json:"unlocks_state_tokens,omitempty"`
	IsUnlockedStateTokens []int64 `json:"is_unlocked_state_tokens,omitempty"`
	SkillTests            []int64 `json:"skill_tests,omitempty"`
}

func Graph(db *gorp.DbMap, scenar *Scenario) (interface{}, error) {

	locations, err := ListLocations(db, scenar)
	if err != nil {
		return nil, err
	}

	var locGraphOut []*LocGraph
	cards := make(map[int64]*CardGraph)

	for _, loc := range locations {
		loc_cards, err := loc.GetCards(db)
		if err != nil {
			return nil, err
		}
		locG := &LocGraph{
			ID:     loc.ID,
			Name:   loc.Name,
			Hidden: loc.Hidden,
		}
		locGraphOut = append(locGraphOut, locG)
		for _, c := range loc_cards {
			cG := &CardGraph{
				ID:          c.ID,
				Description: c.Description,
			}
			locG.Cards = append(locG.Cards, cG)
			cards[c.ID] = cG
		}
	}

	elemToCard := make(map[int64]int64)
	elemCards := make(map[int64]*Card)
	elems, err := GetElementCards(db, scenar)
	if err != nil {
		return nil, err
	}
	for _, elem := range elems {
		elemCards[elem.ID] = elem
	}
	elemLink, err := ListElementLinks(db, scenar, nil, nil)
	if err != nil {
		return nil, err
	}
	for _, el := range elemLink {
		if el.GivesUses {
			c, ok := elemCards[el.IDElement]
			if !ok {
				continue
			}
			// This allows backtracking:
			// card c.ID (element card) IS GIVEN BY el.IDCard
			// if el.IDCard represents an element card too,
			// its origin cab be recursively found too
			// until we reach a location card, which we will design as the actual origin
			// to abstract elements out of this graph
			elemToCard[c.ID] = el.IDCard
		}
	}

	stateTk, err := ListStateTokenLinks(db, scenar, nil)
	if err != nil {
		return nil, err
	}
	for _, tk := range stateTk {
		c, ok := cards[tk.IDCard]
		if !ok {
			continue
		}
		if tk.UnlocksUnlocked {
			c.UnlockStateTokens = append(c.UnlockStateTokens, tk.IDStateToken)
		} else {
			c.IsUnlockedStateTokens = append(c.IsUnlockedStateTokens, tk.IDStateToken)
		}
	}

	locLink, err := ListLocationLinks(db, scenar, nil, nil)
	if err != nil {
		return nil, err
	}
	for _, ll := range locLink {
		c, ok := cards[ll.IDCard]
		if !ok {
			initialLocationCard := recurseElementLinks(ll.IDCard, elemToCard)
			c, ok = cards[initialLocationCard]
			if !ok {
				// Could not backtrack to an initial location card
				continue
			}
		}
		c.Reveals = append(c.Reveals, ll.IDLocation)
	}

	skillTest, err := ListSkillTests(db, scenar, nil, nil)
	if err != nil {
		return nil, err
	}
	for _, st := range skillTest {
		c, ok := cards[st.IDCard]
		if !ok {
			continue
		}
		c.SkillTests = append(c.SkillTests, st.IDStat)
	}

	return locGraphOut, nil
}

func recurseElementLinks(IDCard int64, elemToCard map[int64]int64) int64 {
	var i int64
	ok := true
	retID := IDCard

	for ok {
		i, ok = elemToCard[retID]
		if ok {
			retID = i
		}
	}
	return retID
}
