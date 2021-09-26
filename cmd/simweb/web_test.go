package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/wowsims/tbc/api"
	"google.golang.org/protobuf/proto"
)

var basicSpec = &api.PlayerOptions_ElementalShaman{
	ElementalShaman: &api.ElementalShaman{
		Agent: &api.ElementalShaman_ElementalShamanAgent{
			Type: api.ElementalShaman_ElementalShamanAgent_Adaptive,
		},
		Talents: &api.ShamanTalents{
			// ElementalDevastation
			ElementalFury:      true,
			Convection:         5,
			Concussion:         5,
			ElementalFocus:     true,
			CallOfThunder:      5,
			UnrelentingStorm:   3,
			ElementalPrecision: 3,
			LightningMastery:   5,
			ElementalMastery:   true,
			LightningOverload:  5,
		},
		Options: &api.ElementalShaman_ElementalShamanOptions{
			WaterShield: true,
		},
	},
}

var basicConsumes = &api.Consumes{
	FlaskOfBlindingLight: true,
	BlackenedBasilisk:    true,
	BrilliantWizardOil:   true,
	SuperManaPotion:      true,
	DarkRune:             true,
}

var basicBuffs = &api.Buffs{
	ArcaneBrilliance: true,
	BlessingOfKings:  true,
	Bloodlust:        1,
	MoonkinAura:      api.TristateEffect_TristateEffectRegular,
	ManaSpringTotem:  api.TristateEffect_TristateEffectRegular,
	TotemOfWrath:     1,
	WrathOfAirTotem:  api.TristateEffect_TristateEffectRegular,
}

var p1Equip = &api.EquipmentSpec{
	Items: []*api.ItemSpec{
		{Id: 29035, Gems: []int32{34220, 24059}, Enchant: 29191},
		{Id: 28762},
		{Id: 29037, Gems: []int32{24059, 24059}, Enchant: 28909},
		{Id: 28766},
		{Id: 29519},
		{Id: 29521},
		{Id: 28780},
		{Id: 29520},
		{Id: 30541},
		{Id: 28810},
		{Id: 30667},
		{Id: 28753},
		{Id: 28785},
		{Id: 29370},
		{Id: 28248},
		{Id: 28770, Enchant: 22555},
		{Id: 29268},
	},
}

// TestIndividualSim is just a smoke test to make sure the http server works as expected.
//   Don't modify this test unless the proto defintions change and this no longer compiles.
func TestIndividualSim(t *testing.T) {
	req := &api.IndividualSimRequest{
		Player: &api.Player{
			Options: &api.PlayerOptions{
				Race:     api.Race_RaceTroll10,
				Spec:     basicSpec,
				Consumes: basicConsumes,
			},
			Equipment: p1Equip,
		},
		Buffs: basicBuffs,
		Encounter: &api.Encounter{
			Duration:   120,
			NumTargets: 1,
		},
		Iterations: 5000,
		RandomSeed: 1,
		Debug:      false,
	}

	msgBytes, err := proto.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to encode request: %s", err.Error())
	}

	r, err := http.Post("http://localhost:3333/individualSim", "application/x-protobuf", bytes.NewReader(msgBytes))
	if err != nil {
		t.Fatalf("Failed to POST request: %s", err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("Failed to read result body: %s", err.Error())
		return
	}

	isr := &api.IndividualSimResult{}
	if err := proto.Unmarshal(body, isr); err != nil {
		t.Fatalf("Failed to parse request: %s", err.Error())
		return
	}

	log.Printf("RESULT: %#v", isr)
}

func TestCalcStatWeight(t *testing.T) {
	req := &api.IndividualSimRequest{
		Player: &api.Player{
			Options: &api.PlayerOptions{
				Race:     api.Race_RaceTroll10,
				Spec:     basicSpec,
				Consumes: basicConsumes,
			},
			Equipment: p1Equip,
		},
		Buffs: basicBuffs,
		Encounter: &api.Encounter{
			Duration:   120,
			NumTargets: 1,
		},
		Iterations: 5000,
		RandomSeed: 1,
		Debug:      false,
	}

	msgBytes, err := proto.Marshal(&api.StatWeightsRequest{
		Options:         req,
		StatsToWeigh:    []api.Stat{api.Stat_StatSpellPower, api.Stat_StatSpellCrit, api.Stat_StatSpellHit},
		EpReferenceStat: api.Stat_StatSpellPower,
	})
	if err != nil {
		t.Fatalf("Failed to encode request: %s", err.Error())
	}

	r, err := http.Post("http://localhost:3333/statWeights", "application/x-protobuf", bytes.NewReader(msgBytes))
	if err != nil {
		t.Fatalf("Failed to POST request: %s", err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("Failed to read result body: %s", err.Error())
		return
	}

	isr := &api.IndividualSimResult{}
	if err := proto.Unmarshal(body, isr); err != nil {
		t.Fatalf("Failed to parse request: %s", err.Error())
		return
	}

	log.Printf("RESULT: %#v", isr)
}