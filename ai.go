package main

import (
	"time"
	"sync"
	"math/rand"
)

// reproduce copies all quens from the left half of the mother chessboard
// to the left half of the child chessboard. It also copies all queens from
// the right half of the father chessboard to the right half of the child
// chessboard.
func reproduce(mother, father *chessboard) (first, second *chessboard) {
	first, second = newChessboard(), newChessboard()

	if ChessMaxX % 2 != 0 {
		panic("The X axis of the chessboard must be even!")
	}

	for _, queen := range mother.queens {
		if queen.x <= ChessMaxX / 2 {
			queenCopy := *queen
			err := first.placeQueen(&queenCopy)
			if err != nil {
				panic(err)
			}
		} else if queen.x >=(ChessMaxX / 2) + 1 {
			queenCopy := *queen
			err := second.placeQueen(&queenCopy)
			if err != nil {
				panic(err)
			}
		}

	}

	for _, queen := range father.queens {
		if queen.x >= (ChessMaxX / 2) + 1 {
			queenCopy := *queen
			err := first.placeQueen(&queenCopy)
			if err != nil {
				panic(err)
			}
		} else if queen.x <= ChessMaxX / 2 {
			queenCopy := *queen
			err := second.placeQueen(&queenCopy)
			if err != nil {
				panic(err)
			}
		}
	}

	first.fitness = first.numAttacks()
	second.fitness = second.numAttacks()

	return
}

func mutate(populationMember *chessboard) {
	columnOfQueenToMutate := rand.Intn(len(populationMember.board)) + 1

	// remove old queen
	for row := 1 ; row <= len(populationMember.board[0]) ; row++ {
		tmpQueen := newQueen(columnOfQueenToMutate, row)
		if _, ok := (populationMember.isTaken(tmpQueen)).(*boardLocationTakenError); ok {
			populationMember.removeQueen(tmpQueen)
		}
	}

	rowOfQueenToMutate := rand.Intn(len(populationMember.board[columnOfQueenToMutate])) + 1
	newQueen := newQueen(columnOfQueenToMutate, rowOfQueenToMutate)
	if err := populationMember.placeQueen(newQueen); err != nil {
		panic(err)
	}
}

type population struct {
	population []*chessboard
	exit chan bool

	fitnessLock sync.Mutex
	fitnessRequirement int
}

func newPopulation(numMembers int) *population {
	newPop := &population{}
	newPop.population = make([]*chessboard, numMembers)
	newPop.exit = make(chan bool)
	newPop.fitnessRequirement = 64

	for index := range newPop.population {
		newPop.population[index] = newRandomChessboard()
	}

	return newPop
}

func (pop *population) shouldMutate() bool {
	pop.fitnessLock.Lock()
	if rand.Intn(65) <= pop.fitnessRequirement {
		pop.fitnessLock.Unlock()
		return true
	}
	pop.fitnessLock.Unlock()

	return false
}

func (pop *population) getRandomPop() *chessboard {
	candidates := make([]*chessboard, 0)

	pop.fitnessLock.Lock()
	fitnessReq := pop.fitnessRequirement
	pop.fitnessLock.Unlock()

	for _, chessboard := range pop.population {
		if chessboard.fitness <= fitnessReq {
			candidates = append(candidates, chessboard)
		}
	}

	// if not enough candidates are fit to reproduce, use the population to reproduce
	numCandidates := len(candidates)
	if numCandidates < 2 {
		index := rand.Intn(len(pop.population))
		return pop.population[index]
	} else {
		index := rand.Intn(len(candidates))
		return candidates[index]
	}
}

func (pop *population) genticSearch() (solution *chessboard) {
	time.AfterFunc(8 * time.Second, func() {
		l.Println("TIME UP. Signalling shutdown")
		pop.exit <- true
		pop.exit <- true
		l.Println("Shut down")
	})

	go func() {
		ch := time.Tick(1 * time.Second)
		for {
			select {
			case <-pop.exit:
				return
			case <-ch:
				pop.fitnessLock.Lock()
				pop.fitnessRequirement /= 2
				pop.fitnessLock.Unlock()
				l.Println("Reduced fitnessRequirement (simmulated annealing)")
			}
		}
	}()

	l.Println("Starting search")
	for {
		select {
		case <- pop.exit:
			return nil
		default:
			newPops := make([]*chessboard, 0)
			for i := 0 ; i < len(pop.population); i++ {
				// Use simmulated annealing inside of getRandomPop to control fitness requirements
				mother := pop.getRandomPop()
				father := pop.getRandomPop()
				child1, child2 := reproduce(mother, father)

				if child1.fitness <= 1 {
					return child1
				}

				if child2.fitness <= 1 {
					return child2
				}

				// Use simulated annealing inside shouldMutate to control mutation
				if pop.shouldMutate() {
					mutate(child1)
					if child1.fitness <= 1 {
						return child1
					}
				}

				if pop.shouldMutate() {
					mutate(child2)
					if child2.fitness <= 1 {
						return child2
					}
				}

				newPops = append(newPops, child1, child2)
			}

			// Out with the old, in with the new
			pop.population = newPops
		}
	}
}
