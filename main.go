package main

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func main() {
	// dayOne()
	// dayTwo()
	// dayThree()
	// dayFour()
	// dayFive()
	// daySix()
	// daySeven()
	// dayEight()
}

func dayOne() {
	position := 50
	reachedZero := 0
	clickedOnZero := 0

	fmt.Println(position)
	line := make(chan string)
	go readFileLineByLine("day1Input.txt", line)
	for v := range line {
		zeroClick := 0
		dir := v[0:1]
		click, _ := strconv.Atoi(v[1:])
		if position == 0 {
			if dir == "L" {
				position -= (click - 100)
				for position < 0 {
					zeroClick++
					position += 100
				}
			} else {
				position += click
				for position > 100 {
					zeroClick++
					position -= 100
				}
			}
		} else {
			if dir == "L" {
				position -= click
				for position < 0 {
					zeroClick++
					position += 100
				}
			} else {
				position += click
				for position > 100 {
					zeroClick++
					position -= 100
				}
			}
		}
		if position == 100 {
			position = 0
		}
		clickedOnZero += zeroClick
		if position == 0 {
			reachedZero++
			clickedOnZero++
		}
	}

	fmt.Printf("first answer: %d\n", reachedZero)
	fmt.Printf("second answer: %d\n", clickedOnZero)
}

func dayTwo() {
	sumRepeatingIdsTwice, sumRepeatingIds := 0, 0

	input, err := os.ReadFile("day2Input.txt")
	if err != nil {
		fmt.Println("Error during reading file")
	}
	ranges := strings.Split(string(input), ",")
	for _, idRange := range ranges {
		idStartAndStop := strings.SplitN(idRange, "-", 2)
		start, _ := strconv.Atoi(idStartAndStop[0])
		stop, _ := strconv.Atoi(idStartAndStop[1])
		for i := start; i <= stop; i++ {
			numAsString := strconv.Itoa(i)
			if checkIfRepeatingTwice(numAsString) {
				sumRepeatingIdsTwice += i
			}
			if checkIfRepeating(numAsString) {
				sumRepeatingIds += i
			}
		}
	}
	fmt.Printf("first answer %d \n", sumRepeatingIdsTwice)
	fmt.Printf("second answer %d \n", sumRepeatingIds)
}

func checkIfRepeatingTwice(text string) bool {
	len := len(text)
	if len%2 > 0 {
		return false
	}
	return text[0:len/2] == text[len/2:]
}

func checkIfRepeating(text string) bool {
	len := len(text)
	for l := 1; l <= len/2; l++ {
		if len%l > 0 {
			continue
		}
		if strings.Repeat(text[0:l], int(len/l)) == text {
			return true
		}
	}
	return false
}

func dayThree() {
	banks := make(chan string)
	go readFileLineByLine("day3Input.txt", banks)
	maxJoltage := 0
	maxJoltageOf12 := 0

	for bank := range banks {
		maxJoltage += maxBankJolatge(bank)
		maxJoltageOf12 += maxBankJoltageOf12(bank)
	}

	fmt.Printf("third day first answer: %d\n", maxJoltage)
	fmt.Printf("third day second answer: %d\n", maxJoltageOf12)
}

func maxBankJolatge(bank string) int {
	maxFirst, maxIndex := getMaxInString(bank[:(len(bank) - 1)])
	maxSecond, _ := getMaxInString(bank[maxIndex+1:])

	maxJoltage, _ := strconv.Atoi(fmt.Sprintf("%d%d", maxFirst, maxSecond))
	return maxJoltage
}

func maxBankJoltageOf12(bank string) int {
	var max int
	numbers := ""
	index := 0
	for i := 0; i < 12; i++ {
		minLen := 11 - i
		prevIndex := index
		max, index = getMaxInString(bank[index:(len(bank) - minLen)])
		index = prevIndex + index + 1
		numbers = strings.Join([]string{numbers, strconv.Itoa(max)}, "")
	}
	maxJoltage, _ := strconv.Atoi(numbers)
	return maxJoltage
}

func getMaxInString(numbers string) (int, int) {
	max := 0
	maxIndex := 0
	for i, v := range numbers {
		value, _ := strconv.Atoi(string(v))
		if value > max {
			max = value
			maxIndex = i
		}
	}
	return max, maxIndex
}

func dayFour() {
	mapLines := make(chan string)
	go readFileLineByLine("day4Input.txt", mapLines)

	var rollMap []string
	for line := range mapLines {
		rollMap = append(rollMap, line)
	}

	availableFields := getAvailableRollFields(rollMap)
	fmt.Printf("fourth day first answer: %d\n", len(availableFields))

	allAvailableCnt := len(availableFields)
	for len(availableFields) != 0 {
		rollMap = removeAvailable(rollMap, availableFields)
		availableFields = getAvailableRollFields(rollMap)
		allAvailableCnt += len(availableFields)
	}
	fmt.Printf("fourth day second answer: %d\n", allAvailableCnt)
}

type Field struct {
	Row          int
	Column       int
	HasRoll      bool
	NeighbourCnt int
}

func (f Field) IsAvailable() bool {
	return f.HasRoll && f.NeighbourCnt < 4
}

func getAvailableRollFields(rollMap []string) []Field {
	fields := make(chan Field)
	defer close(fields)

	spotCnt := 0
	for i, line := range rollMap {
		for j := range line {
			spotCnt++
			go countNeighbours(rollMap, i, j, fields)
		}
	}

	var availableFields []Field
	for k := 0; k < spotCnt; k++ {
		field := <-fields
		if field.IsAvailable() {
			availableFields = append(availableFields, field)
		}
	}
	return availableFields
}

func countNeighbours(rollMap []string, row int, col int, out chan Field) {
	neighbourCnt := 0
	if rollMap[row][col] != '@' {
		out <- Field{row, col, false, 0}
		return
	}

	if chkRow := row - 1; chkRow >= 0 {
		if chkCol := col - 1; chkCol >= 0 {
			if rollMap[chkRow][chkCol] == '@' {
				neighbourCnt++
			}
		}
	}
	if chkRow := row - 1; chkRow >= 0 {
		if rollMap[chkRow][col] == '@' {
			neighbourCnt++
		}
	}
	if chkRow := row - 1; chkRow >= 0 {
		if chkCol := col + 1; chkCol < len(rollMap[chkRow]) {
			if rollMap[chkRow][chkCol] == '@' {
				neighbourCnt++
			}
		}
	}
	if chkCol := col - 1; chkCol >= 0 {
		if rollMap[row][chkCol] == '@' {
			neighbourCnt++
		}
	}
	if chkCol := col + 1; chkCol < len(rollMap[row]) {
		if rollMap[row][chkCol] == '@' {
			neighbourCnt++
		}
	}
	if chkRow := row + 1; chkRow < len(rollMap) {
		if chkCol := col - 1; chkCol >= 0 {
			if rollMap[chkRow][chkCol] == '@' {
				neighbourCnt++
			}
		}
	}
	if chkRow := row + 1; chkRow < len(rollMap) {
		if rollMap[chkRow][col] == '@' {
			neighbourCnt++
		}
	}
	if chkRow := row + 1; chkRow < len(rollMap) {
		if chkCol := col + 1; chkCol < len(rollMap[chkRow]) {
			if rollMap[chkRow][chkCol] == '@' {
				neighbourCnt++
			}
		}
	}
	out <- Field{row, col, true, neighbourCnt}
}

func removeAvailable(rollMap []string, available []Field) []string {
	for _, field := range available {
		row := []byte(rollMap[field.Row])
		row[field.Column] = '.'
		rollMap[field.Row] = string(row)
	}
	return rollMap
}

func dayFive() {
	dataLine := make(chan string)
	go readFileLineByLine("day5Input.txt", dataLine)

	var ranges []Range
	for data := range dataLine {
		if data == "" {
			break
		}
		parts := strings.Split(data, "-")
		start, _ := strconv.Atoi(parts[0])
		end, _ := strconv.Atoi(parts[1])
		ranges = append(ranges, Range{start, end})
	}

	var ingredients []string
	for data := range dataLine {
		ingredients = append(ingredients, data)
	}

	unifiedRanges := removeOverlaps(ranges)
	isFresh := make(chan bool)
	ingredinetCnt := 0
	for _, ingredient := range ingredients {
		id, _ := strconv.Atoi(ingredient)
		go checkIngredientInAnyRange(unifiedRanges, id, isFresh)
		ingredinetCnt++
	}

	freshIngredientsCnt := 0
	for ; ingredinetCnt > 0; ingredinetCnt-- {
		if <-isFresh {
			freshIngredientsCnt++
		}
	}
	close(isFresh)

	var idCnt int = 0
	for _, idRange := range unifiedRanges {
		idCnt += idRange.Length()
	}

	fmt.Printf("fifth day first answer: %d\n", freshIngredientsCnt)
	fmt.Printf("fifth day second answer: %d\n", idCnt)
}

type Range struct {
	Start int
	End   int
}

func (r Range) IsIn(number int) bool {
	return r.Start <= number && r.End >= number
}

func (r Range) IsValid() bool {
	return r.Start <= r.End
}

func (r Range) Length() int {
	return r.End - r.Start + 1
}

func checkIngredientInAnyRange(ranges []Range, number int, isFresh chan bool) {
	for _, idRange := range ranges {
		if idRange.IsIn(number) {
			isFresh <- true
			return
		}
	}
	isFresh <- false
}

func removeOverlaps(ranges []Range) []Range {
	var uniqueRanges []Range
	for _, idRange := range ranges {
		shouldAdd := true
		for i, uniqueRange := range uniqueRanges {
			if uniqueRange.IsIn(idRange.Start) {
				idRange = Range{uniqueRange.End + 1, idRange.End}
			}
			if uniqueRange.IsIn(idRange.End) {
				idRange = Range{idRange.Start, uniqueRange.Start - 1}
			}
			if idRange.IsIn(uniqueRange.Start) && idRange.IsIn(uniqueRange.End) {
				uniqueRanges[i] = idRange
				shouldAdd = false
				break
			}
			if !idRange.IsValid() {
				break
			}
		}
		if idRange.IsValid() && shouldAdd {
			uniqueRanges = append(uniqueRanges, idRange)
		}
	}
	return uniqueRanges
}

func daySix() {
	inputData := make(chan string)
	go readFileLineByLine("day6Input.txt", inputData)

	inputChannels := broadcastChannel(inputData, 2)

	answers := make(chan int)
	go func() {
		problems := make(chan Problem)
		go firstPartProblems(inputChannels[0], problems)
		go evaluateProblems(problems, answers)
	}()

	go func() {
		problems := make(chan Problem)
		go secondPartProblems(inputChannels[1], problems)
		go evaluateProblems(problems, answers)
	}()

	for range 2 {
		fmt.Printf("sixth day answer: %d\n", <-answers)
	}
}

type Problem struct {
	Operator  string
	Arguments []int
}

func evaluateProblems(problems chan Problem, out chan int) {
	noOfProblems := 0
	solutions := make(chan int)
	for problem := range problems {
		switch problem.Operator {
		case "+":
			go sum(problem.Arguments, solutions)
		case "*":
			go product(problem.Arguments, solutions)
		}
		noOfProblems++
	}

	grandTotal := 0
	for i := 0; i < noOfProblems; i++ {
		grandTotal += <-solutions
	}
	out <- grandTotal
}

func firstPartProblems(input chan string, out chan Problem) {
	var rows [][]string
	for line := range input {
		re := regexp.MustCompile(`( +)`)
		line = strings.Trim(re.ReplaceAllString(line, " "), " ")
		rows = append(rows, strings.Split(line, " "))
	}

	noOfProblems := 0
	noOfArguments := 0
	if len(rows) > 0 && len(rows[0]) > 0 {
		noOfProblems = len(rows[0])
		noOfArguments = len(rows) - 1
	}

	for column := 0; column < noOfProblems; column++ {
		operator := rows[noOfArguments][column]
		arguments := make([]int, noOfArguments)
		for i := 0; i < noOfArguments; i++ {
			arguments[i], _ = strconv.Atoi(rows[i][column])
		}
		out <- Problem{operator, arguments}
	}
	close(out)
}

func secondPartProblems(input chan string, out chan Problem) {
	var rows []string
	for line := range input {
		rows = append(rows, line)
	}

	var problemArgumentLengths []int
	var arguments [][]string
	var operators []string
	for _, char := range rows[len(rows)-1] {
		if char != ' ' {
			if len(problemArgumentLengths) > 0 {
				problemArgumentLengths[len(operators)-1]--
			}
			operators = append(operators, string(char))
			problemArgumentLengths = append(problemArgumentLengths, 1)
			arguments = append(arguments, make([]string, len(rows)-1))
		} else {
			problemArgumentLengths[len(operators)-1]++
		}
	}
	for arg, line := range rows[0 : len(rows)-1] {
		for problem, length := range problemArgumentLengths {
			arguments[problem][arg] = line[0:length]
			if len(line) > length {
				line = line[length+1:]
			}
		}
	}

	for problem, args := range arguments {
		ansembledArgs := make([]int, len(args[0]))
		for _, argRow := range args {
			for i, argChar := range argRow {
				if argChar != ' ' {
					num, _ := strconv.Atoi(string(argChar))
					ansembledArgs[i] = 10*ansembledArgs[i] + num
				}
			}
		}
		out <- Problem{operators[problem], ansembledArgs}
	}
	close(out)
}

func sum(args []int, solution chan int) {
	sum := 0
	for _, num := range args {
		sum += num
	}
	solution <- sum
}

func product(args []int, solution chan int) {
	prod := 1
	for _, num := range args {
		prod *= num
	}
	solution <- prod
}

func daySeven() {
	inputData := make(chan string)
	go readFileLineByLine("day7Input.txt", inputData)

	splitter := regexp.MustCompile(`\^`)
	beams, _ := findAllStringStartIndex(regexp.MustCompile(`[S]`), <-inputData, -1)
	hitCnt := 0

	for line := range inputData {
		splitterLocations, cnt := findAllStringStartIndex(splitter, line, -1)
		if cnt == 0 {
			continue
		}
		for pos, active := range beams {
			if active > 0 && splitterLocations[pos] > 0 {
				beams[pos-1] = active + beams[pos-1]
				beams[pos+1] = active + beams[pos+1]
				beams[pos] = 0
				hitCnt++
			}
		}
	}
	sumOfTimelines := 0
	for _, timeLines := range beams {
		sumOfTimelines += timeLines
	}
	fmt.Printf("seventh day first answer: %d\n", hitCnt)
	fmt.Printf("seventh day second answer: %d\n", sumOfTimelines)
}

func findAllStringStartIndex(re *regexp.Regexp, input string, n int) (map[int]int, int) {
	indices := make(map[int]int)
	for i := range input {
		indices[i] = 0
	}
	locations := re.FindAllStringIndex(input, n)

	activeLocations := 0
	for _, location := range locations {
		indices[location[0]] = 1
		activeLocations++
	}
	return indices, activeLocations
}

func dayEight() {
	inputData := make(chan string)
	go readFileLineByLine("day8Input.txt", inputData)
	maxSteps := 1000

	var points []Point
	for line := range inputData {
		points = append(points, point(line))
	}

	var circuits []*Circuit
	connectionSteps := 0
	for _, pair := range sortedPairs(points) {
		if connectionSteps == maxSteps {
			a, b, c, _ := largestThreeCircuits(circuits)
			fmt.Printf("eigth day first answer: %d\n", a.Size()*b.Size()*c.Size())
		}
		connectionSteps++

		needsNew := true
		fromId, toId, err := findCircuitsToBeMerged(pair, circuits)
		if err == nil {
			for conn := range circuits[fromId].Connections {
				circuits[toId].AddConnection(conn)
			}
			var newCircuits []*Circuit
			newCircuits = append(newCircuits, circuits[:fromId]...)
			newCircuits = append(newCircuits, circuits[fromId+1:]...)
			circuits = newCircuits
			continue
		}

		for cIndex, cValue := range circuits {
			if cValue.HasPoint(pair.Left) && cValue.HasPoint(pair.Right) {
				needsNew = false
				break
			}
			if cValue.HasPoint(pair.Left) || cValue.HasPoint(pair.Right) {
				circuits[cIndex].AddConnection(pair)
				needsNew = false
				break
			}
		}
		if needsNew {
			circuit := Circuit{make(map[Connection]bool)}
			circuit.AddConnection(pair)
			circuits = append(circuits, &circuit)
		}

		if len(circuits) == 1 && circuits[0].Size() == len(points) {
			fmt.Printf("eigth day second answer: %d\n", pair.Left.X*pair.Right.X)
			break
		}
	}
}

func point(line string) Point {
	coords := strings.Split(line, ",")
	x, _ := strconv.Atoi(coords[0])
	y, _ := strconv.Atoi(coords[1])
	z, _ := strconv.Atoi(coords[2])
	return Point{x, y, z}
}

type Point struct {
	X int
	Y int
	Z int
}

type Connection struct {
	Left     Point
	Right    Point
	Distance float64
}

type Circuit struct {
	Connections map[Connection]bool
}

func (c *Circuit) HasPoint(p Point) bool {
	for conn := range c.Connections {
		if conn.Left == p || conn.Right == p {
			return true
		}
	}
	return false
}

func (c *Circuit) AddConnection(conn Connection) {
	c.Connections[conn] = true
}

func (c *Circuit) Size() int {
	points := make(map[Point]bool)
	for c := range c.Connections {
		points[c.Left] = true
		points[c.Right] = true
	}
	return len(points)
}

func distance(p1, p2 Point) float64 {
	return math.Sqrt(math.Pow(float64(p1.X)-float64(p2.X), 2) + math.Pow(float64(p1.Y)-float64(p2.Y), 2) + math.Pow(float64(p1.Z)-float64(p2.Z), 2))
}

func largestThreeCircuits(circuits []*Circuit) (*Circuit, *Circuit, *Circuit, error) {
	if len(circuits) < 3 {
		return nil, nil, nil, fmt.Errorf("Too few input")
	}

	slices.SortStableFunc(circuits, func(c1, c2 *Circuit) int {
		return c2.Size() - c1.Size()
	})

	return circuits[0], circuits[1], circuits[2], nil
}

func sortedPairs(points []Point) []Connection {
	connections := make(map[Connection]bool)
	for _, basePoint := range points {
		for _, point := range points {
			if point == basePoint {
				continue
			}
			var connection Connection
			if basePoint.X < point.X {
				connection = Connection{basePoint, point, distance(basePoint, point)}
			} else {
				connection = Connection{point, basePoint, distance(basePoint, point)}
			}
			connections[connection] = true
		}
	}

	uniqueConnections := make([]Connection, 0, len(connections))
	for connection := range connections {
		uniqueConnections = append(uniqueConnections, connection)
	}
	slices.SortStableFunc(uniqueConnections, func(c1, c2 Connection) int {
		return int(c1.Distance - c2.Distance)
	})

	return uniqueConnections
}

func findCircuitsToBeMerged(connection Connection, circuits []*Circuit) (int, int, error) {
	left, right := -1, -1
	var err error = nil
	for key, c := range circuits {
		if c.HasPoint(connection.Left) && !c.HasPoint(connection.Right) {
			left = key
		}
		if c.HasPoint(connection.Right) && !c.HasPoint(connection.Left) {
			right = key
		}
	}
	if left == -1 || right == -1 {
		err = fmt.Errorf("No circuits to be merged")
	}
	return left, right, err
}
