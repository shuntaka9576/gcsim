package keqing

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Keqing, NewChar)
}

type char struct {
	*character.Tmpl
	eStartFrame int
	c2ICD       int
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Electro

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 40
	}
	c.Energy = float64(e)
	c.EnergyMax = 40
	c.Weapon.Class = core.WeaponClassSword
	c.NormalHitNum = 5
	c.BurstCon = 3
	c.SkillCon = 5
	c.CharZone = core.ZoneLiyue

	if c.Base.Cons >= 2 {
		c.c2()
	}

	if c.Base.Cons >= 4 {
		c.c4()
	}

	return &c, nil
}

var delay = [][]int{{8}, {20}, {25}, {25, 35}, {34}}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		return 25
	default:
		c.Core.Log.NewEvent("ActionStam not implemented", core.LogActionEvent, c.Index, "action", a.String())
		return 0
	}
}

func (c *char) c4() {

	cb := func(args ...interface{}) bool {

		atk := args[1].(*core.AttackEvent)
		if atk.Info.ActorIndex != c.Index {
			return false
		}

		c.AddMod(core.CharStatMod{
			Key: "c4",
			Amount: func() ([]float64, bool) {

				val := make([]float64, core.EndStatType)
				val[core.ATKP] = 0.25
				return val, true
			},
			Expiry: c.Core.F + 600,
		})
		return false
	}

	c.Core.Events.Subscribe(core.OnOverload, cb, "keqing-c4")
	c.Core.Events.Subscribe(core.OnElectroCharged, cb, "keqing-c4")
	c.Core.Events.Subscribe(core.OnSuperconduct, cb, "keqing-c4")
	c.Core.Events.Subscribe(core.OnSwirlElectro, cb, "keqing-c4")
	c.Core.Events.Subscribe(core.OnCrystallizeElectro, cb, "keqing-c4")
}

func (c *char) c2() {
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if c.Core.F < c.c2ICD {
			return false
		}
		if c.Core.Rand.Float64() < 0.5 {
			c.c2ICD = c.Core.F + 300
			c.QueueParticle("keqing", 1, core.Electro, 100)
			c.Core.Log.NewEvent("keqing c2 proc'd", core.LogCharacterEvent, c.Index, "next ready", c.c2ICD)
		}
		return false
	}, "keqingc2")
}

func (c *char) activateC6(src string) {
	val := make([]float64, core.EndStatType)
	val[core.ElectroP] = 0.06
	c.AddMod(core.CharStatMod{
		Key:    src,
		Amount: func() ([]float64, bool) { return val, true },
		Expiry: c.Core.F + 480,
	})
}
