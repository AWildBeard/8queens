package main

import (
	"io/ioutil"
	"testing"

	"gopkg.in/yaml.v3"
)

type AttackConfig struct {
	Cases map[string]map[string][]string
}

type attackCase struct {
	attackingQueen *queen
	attackedQueens []*queen
	missedQueens   []*queen
}

var config = &AttackConfig{}

func init() {
	configContent, err := ioutil.ReadFile("chessboard_queen_attack_test_cases.yaml")
	if err != nil {
		panic(err)
	}

	if err = yaml.Unmarshal(configContent, config); err != nil {
		panic(err)
	}
}

func TestAttackDetectsCorrectQueens(t *testing.T) {
	for attacker, defenders := range config.Cases {
		t.Logf("attacker: %s\n", attacker)
		for typ, queens := range defenders {
			t.Logf("type: %s\n", typ)
			for _, queen := range queens {
				t.Logf("\t%s\n", queen)
			}
		}
	}
}
