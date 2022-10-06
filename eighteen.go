package main

import (
	"fmt"
)

/**
--- Day 18: Snailfish ---

You descend into the ocean trench and encounter some snailfish. They say they saw the sleigh keys! They'll even tell you
which direction the keys went if you help one of the smaller snailfish with his math homework. Snailfish numbers aren't
like regular numbers. Instead, every snailfish number is a pair - an ordered list of two elements. Each element of the
pair can be either a regular number or another pair. Pairs are written as [x,y], where x and y are the elements within
the pair. Here are some example snailfish numbers, one snailfish number per line:

[1,2]
[[1,2],3]
[9,[8,7]]
[[1,9],[8,5]]
[[[[1,2],[3,4]],[[5,6],[7,8]]],9]
[[[9,[3,8]],[[0,9],6]],[[[3,7],[4,9]],3]]
[[[[1,3],[5,3]],[[1,3],[8,7]]],[[[4,9],[6,9]],[[8,2],[7,3]]]]

This snailfish homework is about addition. To add two snailfish numbers, form a pair from the left and right parameters
of the addition operator. For example, [1,2] + [[3,4],5] becomes [[1,2],[[3,4],5]]. There's only one problem: snailfish
numbers must always be reduced, and the process of adding two snailfish numbers can result in snailfish numbers that
need to be reduced. To reduce a snailfish number, you must repeatedly do the first action in this list that applies to
the snailfish number:

    If any pair is nested inside four pairs, the leftmost such pair explodes.
    If any regular number is 10 or greater, the leftmost such regular number splits.

Once no action in the above list applies, the snailfish number is reduced. During reduction, at most one action applies,
after which the process returns to the top of the list of actions. For example, if split produces a pair that meets the
explode criteria, that pair explodes before other splits occur. To explode a pair, the pair's left value is added to the
first regular number to the left of the exploding pair (if any), and the pair's right value is added to the first
regular number to the right of the exploding pair (if any). Exploding pairs will always consist of two regular numbers.
Then, the entire exploding pair is replaced with the regular number 0. Here are some examples of a single explode
action:

    [[[[[9,8],1],2],3],4] becomes [[[[0,9],2],3],4] (the 9 has no regular number to its left, so it is not added to any
		regular number).
    [7,[6,[5,[4,[3,2]]]]] becomes [7,[6,[5,[7,0]]]] (the 2 has no regular number to its right, and so it is not added
		to any regular number).
    [[6,[5,[4,[3,2]]]],1] becomes [[6,[5,[7,0]]],3].
    [[3,[2,[1,[7,3]]]],[6,[5,[4,[3,2]]]]] becomes [[3,[2,[8,0]]],[9,[5,[4,[3,2]]]]] (the pair [3,2] is unaffected
		because the pair [7,3] is further to the left; [3,2] would explode on the next action).
    [[3,[2,[8,0]]],[9,[5,[4,[3,2]]]]] becomes [[3,[2,[8,0]]],[9,[5,[7,0]]]].

To split a regular number, replace it with a pair; the left element of the pair should be the regular number divided by
two and rounded down, while the right element of the pair should be the regular number divided by two and rounded up.
For example, 10 becomes [5,5], 11 becomes [5,6], 12 becomes [6,6], and so on. Here is the process of finding the reduced
result of [[[[4,3],4],4],[7,[[8,4],9]]] + [1,1]:

after addition: [[[[[4,3],4],4],[7,[[8,4],9]]],[1,1]]
after explode:  [[[[0,7],4],[7,[[8,4],9]]],[1,1]]
after explode:  [[[[0,7],4],[15,[0,13]]],[1,1]]
after split:    [[[[0,7],4],[[7,8],[0,13]]],[1,1]]
after split:    [[[[0,7],4],[[7,8],[0,[6,7]]]],[1,1]]
after explode:  [[[[0,7],4],[[7,8],[6,0]]],[8,1]]

Once no reduce actions apply, the snailfish number that remains is the actual result of the addition operation:
[[[[0,7],4],[[7,8],[6,0]]],[8,1]]. The homework assignment involves adding up a list of snailfish numbers (your puzzle
input). The snailfish numbers are each listed on a separate line. Add the first snailfish number and the second, then
add that result and the third, then add that result and the fourth, and so on until all numbers in the list have been
used once. For example, the final sum of this list is [[[[1,1],[2,2]],[3,3]],[4,4]]:

[1,1]
[2,2]
[3,3]
[4,4]

The final sum of this list is [[[[3,0],[5,3]],[4,4]],[5,5]]:

[1,1]
[2,2]
[3,3]
[4,4]
[5,5]

The final sum of this list is [[[[5,0],[7,4]],[5,5]],[6,6]]:

[1,1]
[2,2]
[3,3]
[4,4]
[5,5]
[6,6]

Here's a slightly larger example:

[[[0,[4,5]],[0,0]],[[[4,5],[2,6]],[9,5]]]
[7,[[[3,7],[4,3]],[[6,3],[8,8]]]]
[[2,[[0,8],[3,4]]],[[[6,7],1],[7,[1,6]]]]
[[[[2,4],7],[6,[0,5]]],[[[6,8],[2,8]],[[2,1],[4,5]]]]
[7,[5,[[3,8],[1,4]]]]
[[2,[2,2]],[8,[8,1]]]
[2,9]
[1,[[[9,3],9],[[9,0],[0,7]]]]
[[[5,[7,4]],7],1]
[[[[4,2],2],6],[8,7]]

The final sum [[[[8,7],[7,7]],[[8,6],[7,7]]],[[[0,7],[6,6]],[8,7]]] is found after adding up the above snailfish numbers:

  [[[0,[4,5]],[0,0]],[[[4,5],[2,6]],[9,5]]]
+ [7,[[[3,7],[4,3]],[[6,3],[8,8]]]]
= [[[[4,0],[5,4]],[[7,7],[6,0]]],[[8,[7,7]],[[7,9],[5,0]]]]

  [[[[4,0],[5,4]],[[7,7],[6,0]]],[[8,[7,7]],[[7,9],[5,0]]]]
+ [[2,[[0,8],[3,4]]],[[[6,7],1],[7,[1,6]]]]
= [[[[6,7],[6,7]],[[7,7],[0,7]]],[[[8,7],[7,7]],[[8,8],[8,0]]]]

  [[[[6,7],[6,7]],[[7,7],[0,7]]],[[[8,7],[7,7]],[[8,8],[8,0]]]]
+ [[[[2,4],7],[6,[0,5]]],[[[6,8],[2,8]],[[2,1],[4,5]]]]
= [[[[7,0],[7,7]],[[7,7],[7,8]]],[[[7,7],[8,8]],[[7,7],[8,7]]]]

  [[[[7,0],[7,7]],[[7,7],[7,8]]],[[[7,7],[8,8]],[[7,7],[8,7]]]]
+ [7,[5,[[3,8],[1,4]]]]
= [[[[7,7],[7,8]],[[9,5],[8,7]]],[[[6,8],[0,8]],[[9,9],[9,0]]]]

  [[[[7,7],[7,8]],[[9,5],[8,7]]],[[[6,8],[0,8]],[[9,9],[9,0]]]]
+ [[2,[2,2]],[8,[8,1]]]
= [[[[6,6],[6,6]],[[6,0],[6,7]]],[[[7,7],[8,9]],[8,[8,1]]]]

  [[[[6,6],[6,6]],[[6,0],[6,7]]],[[[7,7],[8,9]],[8,[8,1]]]]
+ [2,9]
= [[[[6,6],[7,7]],[[0,7],[7,7]]],[[[5,5],[5,6]],9]]

  [[[[6,6],[7,7]],[[0,7],[7,7]]],[[[5,5],[5,6]],9]]
+ [1,[[[9,3],9],[[9,0],[0,7]]]]
= [[[[7,8],[6,7]],[[6,8],[0,8]]],[[[7,7],[5,0]],[[5,5],[5,6]]]]

  [[[[7,8],[6,7]],[[6,8],[0,8]]],[[[7,7],[5,0]],[[5,5],[5,6]]]]
+ [[[5,[7,4]],7],1]
= [[[[7,7],[7,7]],[[8,7],[8,7]]],[[[7,0],[7,7]],9]]

  [[[[7,7],[7,7]],[[8,7],[8,7]]],[[[7,0],[7,7]],9]]
+ [[[[4,2],2],6],[8,7]]
= [[[[8,7],[7,7]],[[8,6],[7,7]]],[[[0,7],[6,6]],[8,7]]]

To check whether it's the right answer, the snailfish teacher only checks the magnitude of the final sum. The magnitude
of a pair is 3 times the magnitude of its left element plus 2 times the magnitude of its right element. The magnitude of
a regular number is just that number. For example, the magnitude of [9,1] is 3*9 + 2*1 = 29; the magnitude of [1,9] is
3*1 + 2*9 = 21. Magnitude calculations are recursive: the magnitude of [[9,1],[1,9]] is 3*29 + 2*21 = 129. Here are a
few more magnitude examples:

    [[1,2],[[3,4],5]] becomes 143.
    [[[[0,7],4],[[7,8],[6,0]]],[8,1]] becomes 1384.
    [[[[1,1],[2,2]],[3,3]],[4,4]] becomes 445.
    [[[[3,0],[5,3]],[4,4]],[5,5]] becomes 791.
    [[[[5,0],[7,4]],[5,5]],[6,6]] becomes 1137.
    [[[[8,7],[7,7]],[[8,6],[7,7]]],[[[0,7],[6,6]],[8,7]]] becomes 3488.

So, given this example homework assignment:

[[[0,[5,8]],[[1,7],[9,6]]],[[4,[1,2]],[[1,4],2]]]
[[[5,[2,8]],4],[5,[[9,9],0]]]
[6,[[[6,2],[5,6]],[[7,6],[4,7]]]]
[[[6,[0,7]],[0,9]],[4,[9,[9,0]]]]
[[[7,[6,4]],[3,[1,3]]],[[[5,5],1],9]]
[[6,[[7,3],[3,2]]],[[[3,8],[5,7]],4]]
[[[[5,4],[7,7]],8],[[8,3],8]]
[[9,3],[[9,9],[6,[4,9]]]]
[[2,[[7,7],7]],[[5,8],[[9,3],[0,2]]]]
[[[[5,2],5],[8,[3,7]]],[[5,[7,5]],[4,4]]]

The final sum is:

[[[[6,6],[7,6]],[[7,7],[7,0]]],[[[7,7],[7,7]],[[7,8],[9,9]]]]

The magnitude of this final sum is 4140. Add up all of the snailfish numbers from the homework assignment in the order
they appear. What is the magnitude of the final sum?

-- Soln --

    			[[6,[5,[4,[3,2]]]],1]  -->  (Becomes [[6,[5,[7,0]]],3])
					/         \
			[6,[5,[4,[3,2]]]]  1
				/      \
				6	[5,[4,[3,2]]]
						/      \
						5	[4,[3,2]]
							  /   \
							  4  [3,2]
								  / \
								 3   2
Considering [[[[[9,8],1],2],3],4] in tree form

			[[[[[9,8],1],2],3],4]
				/      \
	[[[[9,8],1],2],3]   4
		/          \
	[[[9,8],1],2]   3
		/      \
	[[9,8],1]   2
	/      \
  [9,8]     1
   / \
  9   8

Then we can understand operations of SnailFish arithmetics with tree semantics
 - "If any pair is nested inside four pairs, the leftmost such pair explodes."
	--> Perform in-order traversal and find first pair, if present, at depth=4. Add left to the next left at
		depth=3 (if no right branch, d=2, ...). When pair found, add to most right element. (i.e. perform post order
		traversal on d=3 (or 2, 1) pair). Similarly, for right, find parent with right node, add to left most child of
		parent node.

 - "If any regular number is 10 or greater, the leftmost such regular number splits."
	-->  Perform in-order traversal and find first leaf node with value>=10. If found, convert into pair of
		[v//2, (v+1)//2]

 - If neither of these traversals provide an action, then it is reduced.

 - To add to Snailfish numbers, construct new tree with summands as leafs: i.e. for  a+b, evaluate the tree
	a+b
    / \
   a   b

  - "The magnitude of a pair is 3 times the magnitude of its left element plus 2 times the magnitude of its right
	element"
	--> Perform post order traversal. For a node |v| = 3*|left| + 2*|right|
*/

func Eighteen() interface{} {
	sfns := GetAllSnailfishNumbers()

	// Reduce all SnailFishNumbers
	for _, sfn := range sfns {
		start := sfn.Print()
		(&sfn).Reduce()
		fmt.Println(fmt.Sprintf("DONE WITH: %s --> %s", start, (&sfn).Print()))
		fmt.Println()
	}

	// Add leading two SnailfishNumbers, then reduce. Combine all SnailfishNumbers like this.
	initSfn := &sfns[0]
	for _, sfn := range sfns[1:] {
		fmt.Println(fmt.Sprintf("Starting. %s + %s", initSfn.Print(), sfn.Print()))

		addSfn := Add(initSfn, &sfn)
		fmt.Println("Added.", addSfn.Print())
		fmt.Println("l r .", addSfn.left.Print(), addSfn.right.Print())
		(&addSfn).Reduce()

		fmt.Println(fmt.Sprintf("Result.    = %s", addSfn.Print()))
		initSfn = &addSfn
	}

	// Compute magnitude of final SnailfishNumber
	return initSfn.Magnitude()
}

func GetAllSnailfishNumbers() []SnailfishNumber {
	var sfns []SnailfishNumber
	s := GetScanner(18)
	for s.Scan() {
		sfn := ConstructSnailFishNumber(s.Text(), nil)
		sfns = append(sfns, sfn)
	}
	return sfns
}

func ConstructSnailFishNumber(s string, parent *SnailfishNumber) SnailfishNumber {
	var left SnailfishNumber
	var rightIndex int
	result := SnailfishNumber{
		parent: parent,
	}

	if s[1] >= 48 && s[1] <= 57 {
		left = SnailfishNumber{
			v:      string(s[1]),
			left:   nil,
			right:  nil,
			parent: &result,
		}
		rightIndex = 3

	} else {
		left = ConstructSnailFishNumber(s[1:], &result)
		rightIndex = len(left.v) + 2
	}

	var right SnailfishNumber
	if s[rightIndex] >= 48 && s[rightIndex] <= 57 {
		right = SnailfishNumber{
			v:      string(s[rightIndex]),
			left:   nil,
			right:  nil,
			parent: &result,
		}
	} else {
		right = ConstructSnailFishNumber(s[rightIndex:], &result)
	}
	result.v = fmt.Sprintf("[%s,%s]", left.v, right.v)
	result.left = &left
	result.right = &right
	return result

}

type SnailfishNumber struct {
	v      string
	left   *SnailfishNumber
	right  *SnailfishNumber
	parent *SnailfishNumber
}

func (sfn SnailfishNumber) IsLeaf() bool {
	return (sfn.left == nil) && (sfn.right == nil)
}

func (sfn *SnailfishNumber) isLeftChild() bool {
	parent := sfn.parent
	defer func() {
		// recover from panic if one occured. Set err to nil otherwise.
		if recover() != nil {
			fmt.Println("PANICKED", parent.Print(), sfn.Print())
			panic("")
		}
	}()
	return *parent.left == *sfn

}

func (sfn *SnailfishNumber) Print() string {
	if sfn == nil {
		return ""
	}
	if sfn.IsLeaf() {
		if len(sfn.v) > 1 {
			return fmt.Sprintf("%d%d", sfn.v[0]-48, sfn.v[1]-48)
		}
		return fmt.Sprintf("%d", sfn.v[0]-48)
	} else {
		return fmt.Sprintf("[%s,%s]", sfn.left.Print(), sfn.right.Print())
	}
}

func (sfn *SnailfishNumber) Reduce() {
	// Precedence: Explode, Split

	for reductionOccurred := true; reductionOccurred; reductionOccurred = sfn.Explode() || sfn.Split() {
		//fmt.Println()
		//fmt.Println("START OF ROUND:", sfn.Print())
	}
}

func Add(x *SnailfishNumber, y *SnailfishNumber) SnailfishNumber {
	root := SnailfishNumber{
		v:      fmt.Sprintf("[%s,%s]", x.Print(), y.Print()),
		left:   x,
		right:  y,
		parent: nil,
	}
	x.parent = &root
	y.parent = &root
	return root
}

func (sfn *SnailfishNumber) Magnitude() int {
	if sfn.IsLeaf() {
		return unsafeIntParse(sfn.v)
	}
	return sfn.left.Magnitude() + sfn.right.Magnitude()
}

//- "If any pair is nested inside four pairs, the leftmost such pair explodes."
//--> Perform in-order traversal and find first pair, if present, at depth=4. Add left to the next left at
//depth=3 (if no right branch, d=2, ...). When pair found, add to most right element. (i.e. perform post order
//traversal on d=3 (or 2, 1) pair). Similarly, for right, find parent with right node, add to left most child of
//parent node.

func (sfn *SnailfishNumber) Explode() bool {
	//fmt.Println("Exploding....", sfn.Print())
	explodeSfn := GetSfnToExplode(sfn, 0)
	if explodeSfn == nil {
		//fmt.Println("No explosion")
		return false
	}

	// Convert leaf value, E.g. "3" -> into uint8 3 before use
	// TODO: don't respect if explodeSfn.parent.left = explodeSfn
	//fmt.Println("exploding ")
	//PrintSnailfishNumber(sfn)
	//PrintSnailfishNumber(explodeSfn)

	explodeSfn.addToLeftParent(explodeSfn.left.v[0] - 48)
	explodeSfn.addToRightParent(explodeSfn.right.v[0] - 48)

	// Convert exploded Snailfish Number to a "0".
	if explodeSfn.isLeftChild() {
		*(explodeSfn.parent.left) = SnailfishNumber{
			v:      "0",
			left:   nil,
			right:  nil,
			parent: explodeSfn.parent,
		}
	} else {
		*(explodeSfn.parent.right) = SnailfishNumber{
			v:      "0",
			left:   nil,
			right:  nil,
			parent: explodeSfn.parent,
		}
	}

	return true
}

func PrintSnailfishNumber(sfn *SnailfishNumber) {
	if sfn.parent != nil {
		fmt.Printf("root: %s\n", sfn.parent.Print())
	}
	fmt.Printf("node: %s\n", sfn.Print())
	if sfn.IsLeaf() {
		return
	}
	fmt.Printf("left: %s\n", sfn.left.Print())
	fmt.Printf("right: %s\n", sfn.right.Print())
}

func GetSfnToExplode(sfn *SnailfishNumber, depth int) *SnailfishNumber {
	// Don't explode leaf node (i.e. a literal number), explode its parent (e.g. explode [3,4] not 3)
	if sfn.IsLeaf() && depth >= 3+1 {
		return sfn.parent
	}
	if sfn.IsLeaf() {
		return nil
	}

	l := GetSfnToExplode(sfn.left, depth+1)
	if l != nil {
		return l
	}

	r := GetSfnToExplode(sfn.right, depth+1)
	if r != nil {
		return r
	}

	return nil
}

func (sfn *SnailfishNumber) Split() bool {
	//fmt.Println("Splitting....", sfn.Print())
	splitSfn := GetSfnToSplit(sfn)
	if splitSfn == nil {
		//fmt.Println("No split")
		return false
	}

	splitV := splitSfn.v[0] - 48
	splitSfn.left = &SnailfishNumber{
		v:      string((splitV / 2) + 48),
		left:   nil,
		right:  nil,
		parent: splitSfn,
	}

	splitSfn.right = &SnailfishNumber{
		v:      string(((splitV + 1) / 2) + 48),
		left:   nil,
		right:  nil,
		parent: splitSfn,
	}
	return true
}

func GetSfnToSplit(sfn *SnailfishNumber) *SnailfishNumber {
	if sfn.IsLeaf() && sfn.v[0] > 57 {
		return sfn
	}
	if sfn.IsLeaf() {
		return nil
	}

	l := GetSfnToSplit(sfn.left)
	if l != nil {
		return l
	}

	r := GetSfnToSplit(sfn.right)
	if r != nil {
		return r
	}

	return nil

}

func (sfn *SnailfishNumber) addToRightChild(v uint8) {
	if sfn.IsLeaf() {
		sfn.v = string(sfn.v[0] + v)
	} else {
		sfn.right.addToRightChild(v)
	}
}

func (sfn *SnailfishNumber) addToLeftChild(v uint8) {
	if sfn.IsLeaf() {
		sfn.v = string(sfn.v[0] + v)
	} else {
		sfn.left.addToLeftChild(v)
	}
}

func (sfn *SnailfishNumber) addToLeftParent(v uint8) {
	if sfn.parent != nil && !sfn.isLeftChild() {
		sfn.parent.left.addToRightChild(v)
	} else if sfn.parent == nil {
		return
	} else {
		sfn.parent.addToLeftParent(v)
	}
}

func (sfn *SnailfishNumber) addToRightParent(v uint8) {
	if sfn.parent != nil && sfn.isLeftChild() {
		sfn.parent.right.addToRightParent(v)
	} else if sfn.parent == nil {
		return
	} else {
		sfn.parent.addToRightParent(v)
	}
}

/**
		[[[10,0],[15,9]],[[8,15],[0,[5,4]]]]
     		/ 				\
		[[10,0],[15,9]]    [[8,15],[0,[5,4]]]
		  /    \				/    \
		[10,0]  [15, 9]		[8, 15]  [0,[5,4]]
		 / \      /  \		  /  \      /  \
		10, 0    15   9      8   15    0   [5,4]
											/ \
										   5   4
*/
