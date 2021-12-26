package main

import (
	"fmt"
	"sync"
)

/**
--- Day 10: Syntax Scoring ---

You ask the submarine to determine the best route out of the deep-sea cave, but it only replies:

Syntax error in navigation subsystem on line: all of them

All of them?! The damage is worse than you thought. You bring up a copy of the navigation subsystem (your puzzle input).

The navigation subsystem syntax is made of several lines containing chunks. There are one or more chunks on each line, and chunks contain zero or more other chunks. Adjacent chunks are not separated by any delimiter; if one chunk stops, the next chunk (if any) can immediately start. Every chunk must open and close with one of four legal pairs of matching characters:

    If a chunk opens with (, it must close with ).
    If a chunk opens with [, it must close with ].
    If a chunk opens with {, it must close with }.
    If a chunk opens with <, it must close with >.

So, () is a legal chunk that contains no other chunks, as is []. More complex but valid chunks include ([]), {()()()}, <([{}])>, [<>({}){}[([])<>]], and even (((((((((()))))))))).

Some lines are incomplete, but others are corrupted. Find and discard the corrupted lines first.

A corrupted line is one where a chunk closes with the wrong character - that is, where the characters it opens and closes with do not form one of the four legal pairs listed above.

Examples of corrupted chunks include (], {()()()>, (((()))}, and <([]){()}[{}]). Such a chunk can appear anywhere within a line, and its presence causes the whole line to be considered corrupted.

For example, consider the following navigation subsystem:

[({(<(())[]>[[{[]{<()<>>
[(()[<>])]({[<{<<[]>>(
{([(<{}[<>[]}>{[]{[(<()>
(((({<>}<{<{<>}{[]{[]{}
[[<[([]))<([[{}[[()]]]
[{[{({}]{}}([{[{{{}}([]
{<[[]]>}<{[{[{[]{()[[[]
[<(<(<(<{}))><([]([]()
<{([([[(<>()){}]>(<<{{
<{([{{}}[<[[[<>{}]]]>[]]

Some of the lines aren't corrupted, just incomplete; you can ignore these lines for now. The remaining five lines are corrupted:

    {([(<{}[<>[]}>{[]{[(<()> - Expected ], but found } instead.
    [[<[([]))<([[{}[[()]]] - Expected ], but found ) instead.
    [{[{({}]{}}([{[{{{}}([] - Expected ), but found ] instead.
    [<(<(<(<{}))><([]([]() - Expected >, but found ) instead.
    <{([([[(<>()){}]>(<<{{ - Expected ], but found > instead.

Stop at the first incorrect closing character on each corrupted line.

Did you know that syntax checkers actually have contests to see who can get the high score for syntax errors in a file? It's true! To calculate the syntax error score for a line, take the first illegal character on the line and look it up in the following table:

    ): 3 points.
    ]: 57 points.
    }: 1197 points.
    >: 25137 points.

In the above example, an illegal ) was found twice (2*3 = 6 points), an illegal ] was found once (57 points), an illegal } was found once (1197 points), and an illegal > was found once (25137 points). So, the total syntax error score for this file is 6+57+1197+25137 = 26397 points!

Find the first illegal character in each corrupted line of the navigation subsystem. What is the total syntax error score for those errors?

*/
var characterScores = map[uint8]uint64{
	41:  3,     // )
	93:  57,    // ]
	125: 1197,  // }
	62:  25137, // >
}

var closedCharacters = map[uint8]bool{
	41:  true, // )
	62:  true, // >
	93:  true, // ]
	125: true, // }

	40:  false, // (
	60:  false, // <
	91:  false, // [
	123: false, // {
}

var leftRightPairs = map[uint8]uint8{
	40:  41,
	60:  62,
	91:  93,
	123: 125,
}

const errorScoreGoroutines = 1

func Ten() interface{} {

	inputs := make(chan string, 10)
	scores := make(chan uint64, 10)
	wg := sync.WaitGroup{}

	go sendLines(inputs)

	// Get error score
	wg.Add(errorScoreGoroutines)
	for i := 0; i < errorScoreGoroutines; i++ {
		go getErrorScores(inputs, scores, &wg)
	}
	go func(outputs chan uint64, wg *sync.WaitGroup) {
		wg.Wait()
		close(outputs)
	}(scores, &wg)

	// Sum Error Scores
	result := make(chan uint64, 10)
	go func(output chan uint64) {
		defer close(result)
		r := uint64(0)
		for s := range scores {
			r += s
		}
		output <- r
	}(result)

	// 981444
	return <-result
}

func getErrorScores(inputs chan string, outputs chan uint64, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := range inputs {
		s := getErrorScore(i)
		if s != 0 {
			outputs <- s
		}
	}
}

func getErrorScore(s string) uint64 {
	if len(s) == 0 {
		return 0
	}

	stack := []uint8{s[0]}

	for i := 1; i < len(s); i++ {
		// If it's a closing character, check its validity
		if closedCharacters[s[i]] {
			left := stack[len(stack)-1]

			// Remove left character from stack
			if isMatching(left, s[i]) {

				stack = append(stack[:len(stack)-1])
				//stack = stack[:maxInt(0, len(stack)-2)]

				// Then it's an error, get score
			} else {
				fmt.Println(characterScores[s[i]], s[:i+2])
				fmt.Println()
				return characterScores[s[i]]
			}

			// Add to stack
		} else {
			stack = append(stack, s[i])
		}
	}
	return 0
}

func isMatching(l uint8, r uint8) bool {
	fmt.Println(string(l), string(r))
	return leftRightPairs[l] == r
}

func sendLines(outputs chan string) {
	defer close(outputs)
	s := GetScanner(10)
	for s.Scan() {
		outputs <- s.Text()
	}
}
