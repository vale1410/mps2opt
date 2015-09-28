package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

var ver = flag.Bool("ver", false, "Show version info.")
var filename = flag.String("f", "", "Path of the file.")
var gringo = flag.Bool("gringo", false, "Ouput in Potasscos Gringo Format.")
var lp = flag.Bool("lp", false, "Ouput in LP Format (Gurobi/CPlex).")
var minizinc = flag.Bool("minizinc", true, "Ouput in Minizinc Format.")
var maxDom = int64(1000)

func main() {

	flag.Parse()

	var pbs []Linear
	var vars map[string]Bound
	var err error

	if *ver {
		fmt.Println(`MIPLIB .mps converter: Tag 0.2
Copyright (C) Data61 and Valentin Mayer-Eichberger
License GPLv2+: GNU GPL version 2 or later <http://gnu.org/licenses/gpl.html>
There is NO WARRANTY, to the extent permitted by law.`)
		return
	}

	f := ""
	if strings.HasSuffix(flag.Arg(0), "mps") {
		f = flag.Arg(0)
	} else if strings.HasSuffix(*filename, "mps") {
		f = *filename
	}

	if f != "" {
		pbs, vars, err = ParseMPS(f)
		if err != nil {
			fmt.Println(err)
			os.Exit(100)
			return
		}

		if *gringo {
			PrintGringo(pbs, vars)
		} else if *lp {
			PrintLP(pbs, vars)
		} else if *minizinc {
			PrintMinizinc(pbs, vars)
		}

	} else {
		os.Exit(100)
	}

	return
}

type Constraint struct {
	name string
}

type EquationType int

const (
	LE  EquationType = iota //"<="
	GE                      //">="
	EQ                      //"=="
	OPT                     //"MIN"
)

func (e EquationType) String() string {
	switch e {
	case LE:
		return "<="
	case GE:
		return ">="
	case EQ:
		return "=="
	case OPT:
		return "MIN"
	}
	return ""
}

type Bound struct {
	lb, ub int64
}

type Entry struct {
	id     string
	Weight int64
}

type Linear struct {
	Desc    string
	Entries []Entry
	K       int64
	Typ     EquationType
}

func getEquationType(s string) EquationType {
	if s == "E" {
		return EQ
	} else if s == "L" {
		return LE
	} else if s == "G" {
		return GE
	}
	if s != "N" {
		panic("Unknown equation type " + s)
	}
	return OPT
}

func PrintMinizinc(pbs []Linear, vars map[string]Bound) {

	for x, b := range vars {
		fmt.Println("var", b.lb, "..", b.ub, ":", x, ";")
	}

	for _, t := range pbs {
		if len(t.Entries) > 0 {

			switch t.Typ {
			case LE:
				fmt.Print("constraint ")
			case GE:
				fmt.Print("constraint ")
			case EQ:
				fmt.Print("constraint ")
			case OPT:
				fmt.Print("solve minimize ")
			}

			for i, e := range t.Entries {
				if i == 0 {
					if e.Weight > 0 {
						if e.Weight == 1 {
							fmt.Print(e.id)
						} else {
							fmt.Print(e.Weight, "*", e.id)
						}
					} else {
						fmt.Print(e.Weight, "*", e.id)
					}
				} else {
					if e.Weight > 0 {
						if e.Weight == 1 {
							fmt.Print(" + ", e.id)
						} else {
							fmt.Print(" + ", e.Weight, "*", e.id)
						}
					} else {
						if e.Weight == -1 {
							fmt.Print(" - ", e.id)
						} else {
							fmt.Print(" - ", -e.Weight, "*", e.id)
						}
					}
				}
			}

			switch t.Typ {
			case LE:
				fmt.Print(" <= ", t.K, ";\n")
			case GE:
				fmt.Print(" >= ", t.K, ";\n")
			case EQ:
				fmt.Print(" == ", t.K, ";\n")
			case OPT:
				fmt.Print(" ;\n")
			}
		}
	}

	//fmt.Print("output [")
	//first := true
	//for x, _ := range vars {
	//	if first {
	//		first = false
	//	} else {
	//		fmt.Print(",")
	//	}
	//	fmt.Print("\"", x, "\",\"=\",", "show(", x, "),\"\\n\"\n")
	//}
	//fmt.Println("];")
}

func PrintEntries(entries []Entry) {
	if len(entries) > 0 {

		for i, e := range entries {
			if i == 0 {
				if e.Weight > 0 {
					if e.Weight == 1 {
						fmt.Print(e.id)
					} else {
						fmt.Print(e.Weight, " ", e.id)
					}
				} else {
					fmt.Print(e.Weight, " ", e.id)
				}
			} else {
				if e.Weight > 0 {
					if e.Weight == 1 {
						fmt.Print(" + ", e.id)
					} else {
						fmt.Print(" + ", e.Weight, " ", e.id)
					}
				} else {
					if e.Weight == -1 {
						fmt.Print(" - ", e.id)
					} else {
						fmt.Print(" - ", -e.Weight, " ", e.id)
					}
				}
			}
		}
	}
}

func PrintLP(pbs []Linear, vars map[string]Bound) {

	// find optimization
	fmt.Println("Minimize")
	for _, t := range pbs {
		if t.Typ == OPT {
			PrintEntries(t.Entries)
			fmt.Println()
			break
		}
	}

	// iterate over the others
	fmt.Println("Subject To")
	for _, t := range pbs {
		if t.Typ == OPT {
			continue
		}
		PrintEntries(t.Entries)

		switch t.Typ {
		case LE:
			fmt.Print(" <= ", t.K, "\n")
		case GE:
			fmt.Print(" >= ", t.K, "\n")
		case EQ:
			fmt.Print(" = ", t.K, "\n")
		}
	}

	// Bounds
	fmt.Println("Bounds")

	for x, b := range vars {
		if b.lb != 0 || b.ub != 1 {
			if b.ub == math.MaxInt32 {
				fmt.Println(b.lb, " <= ", x)
			} else {
				fmt.Println(b.lb, " <= ", x, "<=", b.ub)
			}
		}
	}

	fmt.Println("Generals")
	c := 0

	for x, b := range vars {
		if b.lb != 0 || b.ub != 1 {
			fmt.Print(x, " ")
			c++
			if c == 10 {
				c = 0
				fmt.Println()
			}
		}
	}
	if c != 0 {
		c = 0
		fmt.Println()
	}

	fmt.Println("Binary")

	for x, b := range vars {
		if b.lb == 0 && b.ub == 1 {
			fmt.Print(x, " ")
			c++
			if c == 10 {
				c = 0
				fmt.Println()
			}
		}
	}
	if c != 0 {
		c = 0
		fmt.Println()
	}
	fmt.Println("End")
}

func PrintGringo(pbs []Linear, vars map[string]Bound) {

	fmt.Println("#hide.")

	for x, b := range vars {

		if b.ub-b.lb > maxDom {
			fmt.Println("Variable domain exeeds maxDom: ", b.ub, b.lb)
			panic("Variable domain too big")
		}

		if b.lb == 0 && b.ub == 1 {
			fmt.Println("{", x, "}.")
		} else {
			fmt.Println(x+"_ge_"+strconv.FormatInt(b.lb, 10), ".")
			for d := b.lb + 1; d <= b.ub; d++ {
				fmt.Println("{", x+"_ge_"+strconv.FormatInt(d, 10), "}.")
				if d != b.ub {
					fmt.Println(":- not", x+"_ge_"+strconv.FormatInt(d, 10), ",", x+"_ge_"+strconv.FormatInt(d+1, 10), ".")
				}
			}
		}
	}
	for _, t := range pbs {
		if len(t.Entries) > 0 {

			switch t.Typ {
			case LE:
				fmt.Print(":- ", t.K+1, " [ ")
			case GE:
				fmt.Print(":- [ ")
			case EQ:
				fmt.Print(":- not ", t.K, " [ ")
			case OPT:
				fmt.Print("#minimize[")
			}

			for i, e := range t.Entries {
				if i != 0 {
					fmt.Print(" , ")
				}
				b := vars[e.id]
				if b.lb == 0 && b.ub == 1 {
					fmt.Print(e.id, "=", e.Weight)
				} else {
					//if b.lb != 1 {
					//	fmt.Print(" true=", b.lb-1, " , ")
					//}
					for d := b.lb + 1; d <= b.ub; d++ {
						fmt.Print(e.id+"_ge_"+strconv.FormatInt(d, 10), "=", e.Weight)
						if d != b.ub {
							fmt.Print(" , ")
						}
					}
				}
			}

			switch t.Typ {
			case LE:
				fmt.Print(" ]")
			case GE:
				fmt.Print(" ] ", t.K-1)
			case EQ:
				fmt.Print(" ] ", t.K)
			case OPT:
				fmt.Print("]")
			}
			fmt.Println(". %", t.Desc)
		}
	}
}

func reformatIdentifier(s string) string {
	s = strings.ToLower(strings.Replace(s, "#", "_", -1))
	s = strings.Replace(s, ",", "_", -1)
	s = strings.Replace(s, "(", "_", -1)
	s = strings.Replace(s, ")", "_", -1)
	return s
}

func ParseMPS(f string) (pbs []Linear, vars map[string]Bound, err error) {
	input, err2 := os.Open(f)
	defer input.Close()
	if err2 != nil {
		err = errors.New("Please specifiy correct path to instance. Does not exist")
		return
	}
	scanner := bufio.NewScanner(input)

	pbs = make([]Linear, 0)
	vars = make(map[string]Bound)
	rowMap := make(map[string]int)

	state := 0
	i := -1

	for scanner.Scan() {
		i++
		l := strings.Trim(scanner.Text(), " ")
		if l == "" {
			continue
		}

		entries := strings.Fields(l)
		ts := Linear{}

		switch state {
		case 0:
			{
				if entries[0] == "NAME" {
					state = 1
				}
				//fmt.Println("name : ", entries[1])
			}
		case 1:
			{
				if entries[0] == "ROWS" {
					state = 2
				}
			}
		case 2: // rows
			if entries[0] == "COLUMNS" {
				state = 3
			} else {
				ts.Typ = getEquationType(entries[0])
				//ts.Desc = entries[1] strings.HasPrefix(entries[0], "LI")
				row := reformatIdentifier(entries[1])
				ts.Desc = row
				pbs = append(pbs, ts)
				rowMap[row] = len(pbs) - 1
			}
		case 3: // COLUMNS
			if entries[0] == "INTM" || entries[0] == "INT1" || entries[0] == "INT1END" || entries[0] == "INTEND" ||
				entries[0] == "INTSTART" || strings.HasPrefix(entries[0], "MARK") {
				// do nothing
				state = 3
			} else if entries[0] == "RHS" {
				state = 4
			} else {

				v := reformatIdentifier(entries[0])

				if _, ok := vars[v]; !ok {
					vars[v] = Bound{0, math.MaxInt32}
				}

				row := reformatIdentifier(entries[1])
				//fmt.Println(v, row)
				a, err2 := strconv.ParseInt(entries[2], 10, 64)
				// also parse Float64 and convert?
				if err2 != nil {
					fmt.Println("Problem in line", i, ":", l)
					err = errors.New("wrong number " + entries[2])
					return
				}

				e := Entry{v, a}
				pbs[rowMap[row]].Entries =
					append(pbs[rowMap[row]].Entries, e)

				if len(entries) == 5 {
					a, err = strconv.ParseInt(entries[4], 10, 64)
					// also parse Float64 and convert?
					if err != nil {
						fmt.Println("Problem in line", i, ":", l)
						err = errors.New("wrong number " + entries[4])
					}

					e = Entry{v, a}
					row = reformatIdentifier(entries[3])
					//fmt.Println(v, row)
					pbs[rowMap[row]].Entries =
						append(pbs[rowMap[row]].Entries, e)
				}
			}
		case 4:
			{
				if entries[0] == "BOUNDS" {
					state = 5
				} else if entries[0] == "RANGE" {
					fmt.Println("Problem in line", i, ":", l)
					err = errors.New("RANGE not supported yet in translator")
					return
				} else {
					a, err2 := strconv.ParseInt(entries[2], 10, 64)
					if err2 != nil {
						fmt.Println("Problem in line", i, ":", l)
						err = errors.New("wrong number " + entries[2])
						return
					}
					row := reformatIdentifier(entries[1])
					pbs[rowMap[row]].K = a
					if len(entries) == 5 {
						a, err2 = strconv.ParseInt(entries[4], 10, 64)
						// also parse Float64 and convert?
						if err2 != nil {
							fmt.Println("Problem in line", i, ":", l)
							err = errors.New("wrong number " + entries[4])
							return
						}

						row = reformatIdentifier(entries[3])
						//fmt.Println(v, row)
						pbs[rowMap[row]].K = a
					}
				}
			}
		case 5:
			{
				if entries[0] == "ENDATA" {
					state = 6
					continue
				}

				//LO    lower bound        b <= x
				//LI    integer variable   b <= x (< +inf)
				//UP    upper bound        x <= b
				//FX    fixed variable     x = b
				//FR    free variable
				//MI    lower bound -inf   -inf < x
				//BV    binary variable    x = 0 or 1

				if len(entries) <= 2 || len(entries) >= 5 {
					fmt.Println("Problem in line", i, ":", l)
					err = errors.New("wrong number of fields")
					return
				}

				id := reformatIdentifier(entries[2])
				b, ok := vars[id]
				if !ok {
					fmt.Println("Problem in line", i, ":", l)
					err = errors.New("variable not existing " + id)
					return
				}

				if strings.HasPrefix(entries[0], "BV") {
					b = Bound{0, 1}
				} else if strings.HasPrefix(entries[0], "FX") {
					value, err2 := strconv.ParseInt(entries[3], 10, 64)
					if err2 != nil {
						fmt.Println("Problem in line", i, ":", l)
						err = errors.New("wrong number " + entries[3])
						return
					}
					b = Bound{value, value}
				} else if strings.HasPrefix(entries[0], "LO") || strings.HasPrefix(entries[0], "LI") {
					value, err2 := strconv.ParseInt(entries[3], 10, 64)
					if err2 != nil {
						fmt.Println("Problem in line", i, ":", l)
						err = errors.New("wrong number " + entries[3])
						return
					}
					b.lb = value
				} else if strings.HasPrefix(entries[0], "UP") {
					value, err2 := strconv.ParseInt(entries[3], 10, 64)
					if err2 != nil {
						fmt.Println("Problem in line", i, ":", l)
						err = errors.New("wrong number " + entries[3])
						return
					}
					b.ub = value
				} else {
					fmt.Println("Problem in line", i, ":", l)
					err = errors.New("not supported yet in translator.")
					return
				}

				vars[id] = b

			}
		default:
			{
			}
		}
	}
	return
}
