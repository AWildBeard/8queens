package main

import (
	"io/ioutil"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type attackCase struct {
	attackingQueen *queen
	attackedQueens []*queen
	missedQueens   []*queen
}

var (
	testCases []attackCase
)

func init() {

	type AttackConfig struct {
		Cases map[string]map[string][]string
	}

	config := &AttackConfig{}

	configContent, err := ioutil.ReadFile("chessboard_queen_attack_test_cases.yaml")
	if err != nil {
		panic(err)
	}

	if err = yaml.Unmarshal(configContent, config); err != nil {
		panic(err)
	}

	testCases = make([]attackCase, 0)

	for attacker, defenders := range config.Cases {
		newTestCase := attackCase{}
		newTestCase.attackingQueen = newQueen(int(attacker[0]-48), int(attacker[2]-48))
		newTestCase.attackedQueens = make([]*queen, 0)

		for typ, queens := range defenders {
			if typ == "victims" {
				for _, queen := range queens {
					newTestCase.attackedQueens = append(newTestCase.attackedQueens, newQueen(int(queen[0]-48), int(queen[2]-48)))
				}
			} else {
				for _, queen := range queens {
					newTestCase.missedQueens = append(newTestCase.missedQueens, newQueen(int(queen[0]-48), int(queen[2]-48)))
				}
			}
		}

		testCases = append(testCases, newTestCase)
	}
}

func TestAttackDetectsCorrectQueens(t *testing.T) {
	for _, testCase := range testCases {
		board := newChessboard()

		for _, victim := range testCase.attackedQueens {
			board.placeQueen(victim)
		}

		for _, survivor := range testCase.missedQueens {
			board.placeQueen(survivor)
		}

		board.placeQueen(testCase.attackingQueen)
		sut := board.attack(testCase.attackingQueen)

		builder := &strings.Builder{}
		board.print(builder)
		t.Log("------------------------------------")
		t.Log(builder.String())
		t.Logf("attacker: %v", *testCase.attackingQueen)
		for _, survivor := range testCase.missedQueens {
			for _, resultQueen := range sut {
				if *survivor == *resultQueen {
					t.Logf("survivor %v found in victim set", *survivor)
					t.Fail()
				}
			}
		}

		for _, victim := range testCase.attackedQueens {
			found := false
			for _, resultQueen := range sut {
				if *victim == *resultQueen {
					found = true
				}
			}

			if !found {
				t.Logf("victim %v not found in victim set", *victim)
				t.Fail()
			}
		}
	}
}
