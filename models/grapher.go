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
	Blocking              bool    `json:"blocking"`
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
			continue
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
		if st.Blocking {
			c.Blocking = true
		}
		c.SkillTests = append(c.SkillTests, st.IDStat)
	}

	return locGraphOut, nil
}
