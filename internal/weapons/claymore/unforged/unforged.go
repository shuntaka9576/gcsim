package unforged

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("the unforged", weapon)
	core.RegisterWeaponFunc("theunforged", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {

	shd := .15 + float64(r)*.05
	c.Shields.AddBonus(func() float64 {
		return shd
	})

	stacks := 0
	icd := 0
	duration := 0

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if icd > c.F {
			return false
		}
		if duration < c.F {
			stacks = 0
		}
		stacks++
		if stacks > 5 {
			stacks = 0
		}
		icd = c.F + 18
		return false
	}, fmt.Sprintf("memory-dust-%v", char.Name()))

	val := make([]float64, core.EndStatType)
	atk := 0.03 + 0.01*float64(r)
	char.AddMod(core.CharStatMod{
		Key:    "memory",
		Expiry: -1,
		Amount: func() ([]float64, bool) {
			if duration > c.F {
				val[core.ATKP] = atk * float64(stacks)
				if c.Shields.IsShielded(char.CharIndex()) {
					val[core.ATKP] *= 2
				}
				return val, true
			}
			stacks = 0
			return nil, false
		},
	})
	return "theunforged"
}
