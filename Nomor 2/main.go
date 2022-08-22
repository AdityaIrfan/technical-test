package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func main() {
	var players int
	fmt.Print("Pemain : ")
	_, err := fmt.Scanf("%d", &players)
	if err != nil {
		fmt.Println("Hanya boleh menginputkan angka")
	} else {
		var dices int
		fmt.Print("Dadu : ")
		_, err := fmt.Scanf("%d", &dices)
		if err != nil {
			fmt.Println("Hanya boleh menginputkan angka")
		} else {
			fmt.Println("====================")
			process(players, dices)
		}
	}
}

func process(players int, dices int) {
	playerDice := definePlayerDice(players, dices)
	start(1, players, playerDice, nil)
}

func start(turn int, players int, playerDice [][]int, results [][]int) {

	fmt.Println("Giliran " + strconv.Itoa(turn) + " lempar dadu: ")
	rollValue := rollDice(players, playerDice)
	printMe(rollValue, results)

	if turn == 1 {
		results = definePlayerResult(players) // [0] is a lot of value 1, [1] is score
	}

	rollValue, results = evaluate(rollValue, results)
	rollValue = moveValue(rollValue, results)
	fmt.Println("Setelah Evaluasi:")
	if isDone, last := printMe(rollValue, results); isDone {
		fmt.Println("Game berakhir karena hanya pemain #" + strconv.Itoa(last) + " yang memiliki dadu.")
		if isDefault, isWinnerMoreThanOne, indexes := getTheWinner(results); !isDefault {
			if !isWinnerMoreThanOne {
				fmt.Println("Game dimenangkan oleh pemain #" + indexes + " karena memiliki poin lebih banyak dari pemain lainnya")
			} else {
				fmt.Println("Ada lebih dari satu pemain yang memiliki nilai tinggi, yaitu " + indexes)
			}
		} else {
			fmt.Println("Tidak ada yang memenangkan game karena tidak ada yang memiliki poin")
		}
		return
	}
	turn++
	start(turn, players, rollValue, results)
}

func moveValue(rollValue [][]int, results [][]int) [][]int {
	for i, result := range results {
		if result[0] > 0 {
			for j := 0; j < result[0]; j++ {
				if i != len(rollValue)-1 {
					rollValue[i+1] = append(rollValue[i+1], 1)
					result[0]--
				} else {
					rollValue[0] = append(rollValue[0], 1)
					result[0]--
				}
			}
		}
	}

	return rollValue
}

func definePlayerDice(players int, dices int) [][]int {
	var playerDice [][]int

	for i := 0; i < players; i++ {
		var temp []int
		for j := 0; j < dices; j++ {
			temp = append(temp, j)
		}
		playerDice = append(playerDice, temp)
	}

	return playerDice
}

func rollDice(players int, playerDice [][]int) [][]int {
	var rollValue [][]int
	rand.Seed(time.Now().UTC().UnixNano())

	for i := 0; i < players; i++ {
		var temp []int = nil
		for j := 0; j < len(playerDice[i]); j++ {
			temp = append(temp, rand.Intn(7-1)+1)
		}
		rollValue = append(rollValue, temp)
	}

	return rollValue
}

/** return bool means isDone */
func evaluate(rollValue [][]int, results [][]int) ([][]int, [][]int) {
	var (
		check = func(value []int, i int) {
			for _, each := range value {
				if each == 1 {
					rollValue[i] = RemoveIndex(rollValue[i], getIndex(value, each))
					value = rollValue[i]
					results[i][0]++
				} else if each == 6 {
					rollValue[i] = RemoveIndex(rollValue[i], getIndex(value, each))
					value = rollValue[i]
					results[i][1]++
				}
			}
		}
	)

	for i, value := range rollValue {
		check(value, i)
	}

	return rollValue, results
}

func definePlayerResult(players int) [][]int {
	var results [][]int

	for i := 0; i < players; i++ {
		results = append(results, []int{0, 0})
	}

	return results
}

func getIndex(values []int, value int) int {
	for i, v := range values {
		if v == value {
			return i
		}
	}

	return -1
}

func RemoveIndex(s []int, index int) []int {
	ret := make([]int, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

func printMe(rollValue [][]int, result [][]int) (bool, int) {
	var (
		remainingPlayer int
		last            int
		isResult        bool
	)

	if result != nil {
		isResult = true
	}

	for i, value := range rollValue {
		var point int = 0

		if isResult {
			point = result[i][1]
		}

		fmt.Print("\tPemain #" + strconv.Itoa(i+1) + " (" + strconv.Itoa(point) + "): ")
		if len(value) != 0 {
			for j, v := range value {
				if j != len(value)-1 {
					fmt.Print(strconv.Itoa(v) + ", ")
				} else {
					fmt.Print(v)
				}
			}
			last = i + 1
			remainingPlayer++
		} else {
			fmt.Print("_ (Berhenti bermain karena tidak memiliki dadu)")
		}
		fmt.Println()
	}

	if remainingPlayer <= 1 {
		return true, last
	}

	return false, last
}

/** return isDefault (there's no winner), isDraw (more than one player get highest score), index as string*/
func getTheWinner(result [][]int) (bool, bool, string) {
	var (
		indexes   []string
		isDefault bool = true
		highest   int
	)
	for i, value := range result {

		if value[1] > highest {
			indexes = nil
			highest = value[1]
			indexes = append(indexes, strconv.Itoa(i+1))
			isDefault = false
		}

		if highest != 0 {
			if value[1] == highest {
				indexes = append(indexes, strconv.Itoa(i+1))
				isDefault = false
			}
		}
	}

	return isDefault, len(removeDuplicateStr(indexes)) > 1, strings.Join(removeDuplicateStr(indexes)[:], ", ")
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
