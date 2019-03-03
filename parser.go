package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type Subject struct {
	Name   string
	Lector string
	Room   string
}

type Department struct {
	Number  string
	Lessons []Subject
}

type Interval struct {
	Start int
	End   int
}

var re = regexp.MustCompile(`[a-zA-z]([0-9]+)`)

func main() {
	res, err := http.Get("http://ras.phys.msu.ru/table/4/1.html")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatal("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if link, ok := s.Attr("href"); ok {
			fmt.Println(link)
			text := s.Text()
			fmt.Println(text)
		}
	})

	course := "4"
	var reGrp = regexp.MustCompile(course + `\d{2}`)
	//	var reInterval = regexp.MustCompile(`(` +  course + `\d{2})\s*\-\s*` + `(` +  course + `\d{2})`)

	grpbegin := "ГРУППЫ >>"
	grpEnd := "<< ГРУППЫ"
	var grpsFound int

	var isGroups bool
	groups := make(map[string]string)
	departments := make([]Department, 0, 5)
	Clss := make(map[string]string)

	eachColumn := make(map[int][]string)

	indx := 0
	doc.Find("td").Each(func(i int, std *goquery.Selection) {
		//		fmt.Println("TD")
		if class, ok := std.Attr("class"); ok {
			Clss[class] = ""
		}

		if grpsFound > 1 {
			return
		}
		//if set []Subject, 0, 5, program will panic. WHY?

		text := std.Text()
		if isGroups && text != grpEnd {
			resFromReg := reGrp.FindAllString(text, -1)
			eachColumn[indx] = resFromReg
			indx++
			for _, val := range resFromReg {
				depart := Department{Lessons: make([]Subject, 5, 5)}
				depart.Number = val
				departments = append(departments, depart)
			}
			groups[text] = ""
		}
		if text == grpbegin {
			grpsFound++
			isGroups = true
		} else if text == grpEnd {
			isGroups = false
		}
	})

	for key, val := range eachColumn {
		fmt.Println(key, val)
	}

	var time string
	var nextStr bool

	var ind int
	tditem := "tditem"
	tdsmall := "tdsmall"
	tdtime := "tdtime"
	t := "9:00- - -  10:35"
	var tmp int

	var classBeforeSmall0 string
	var numberBeforeSmall0 int
	var countSmall0 int
	var n int
	var is2Weeks bool
	var Spans = make([]Interval, 10, 10)
	var insertedGroups = make([]string, 5)

	var num int
	doc.Find("td").Each(func(i int, std *goquery.Selection) {
		text := std.Text()

		if class, ok := std.Attr("class"); ok {
			//			For debugging. To show only Monday.
			if text == t {
				tmp++
			}
			//			if tmp > 6 || tmp < 5 {
			if tmp > 2 {
				return
			}

			if strings.Contains(class, tdtime) {
				if time == text {
					num = 0
					nextStr = true
					ind = 0
					numberBeforeSmall0 = 0
				} else if time == "" {
					time = text
					nextStr = false
				} else {
					num = 0
					Spans = make([]Interval, 10, 10)
					n++
					time = text
					nextStr = false
					ind = 0
				}
			}

			std.Find("td").Each(func(i int, sel *goquery.Selection) {
				if small, ok := sel.Attr("class"); ok {
					if strings.Contains(small, "tdsmall0") {
						countSmall0++
					}
				}
			})

			if countSmall0 <= 0 {
				insertedGroups = make([]string, 5)
			}

			if countSmall0 > 0 && class != tdsmall+"0" {
				number := fromStringToInt(class)
				numberBeforeSmall0 = number
				classBeforeSmall0 = class

				return
			} else if countSmall0 == 0 {
				classBeforeSmall0 = class
			}

			var allGr = make([]string, 0, 5)
			var room string
			std.Find("nobr").Each(func(i int, sel *goquery.Selection) {
				room = sel.Text()
			})

			if strings.Contains(classBeforeSmall0, tditem) {
				is2Weeks = true
			} else {
				is2Weeks = false
			}

			if strings.Contains(class, tditem) {
				fmt.Println(class)
				number := fromStringToInt(class)

				subject := parseGroups(text, room)
				fmt.Printf("Name: %v\nRoom: %v\nLector: %v\n", subject.Name, subject.Room, subject.Lector)
				resFromReg := reGrp.FindAllString(text, -1)

				for i := ind; i < ind+number; i++ {
					allGr = append(allGr, eachColumn[i]...)
				}

				departments, insertedGroups = parseLine(departments, allGr, resFromReg, insertedGroups, subject, text, countSmall0-1, n, nextStr, is2Weeks)
				ind = ind + number

			} else if strings.Contains(class, tdsmall) {
				fmt.Println(class)
				number := fromStringToInt(class)

				subject := parseGroups(text, room)
				fmt.Printf("Name: %v\nRoom: %v\nLector: %v\n", subject.Name, subject.Room, subject.Lector)
				resFromReg := reGrp.FindAllString(text, -1)

				if numberBeforeSmall0 == 0 {
					numberBeforeSmall0 = number
				}

				if !nextStr {
					//					fmt.Println("==========================inserted groups and countSmall0 and is2Weeks================================")
					//					fmt.Println(insertedGroups, countSmall0, is2Weeks)
					//					fmt.Println("==========================================================")

					fmt.Println(class, ind, numberBeforeSmall0, classBeforeSmall0)
					if !strings.Contains(class, tdsmall+"0") || !strings.Contains(classBeforeSmall0, tditem) {
						if num == 0 || Spans[num].Start != ind && Spans[num].End != ind+numberBeforeSmall0 {
							span := Interval{Start: ind, End: ind + numberBeforeSmall0}
							Spans[num] = span
							fmt.Println("SPANS!!!!!", num, Spans[num])
							num++
						}
					}
					for i := ind; i < ind+numberBeforeSmall0; i++ {
						allGr = append(allGr, eachColumn[i]...)
					}
					//					departments = parseLine(departments, allGr, resFromReg, insertedGroups, subject, text, n, nextStr)

				} else { //NEXT STRING
					is2Weeks = false
					fmt.Println(Spans)
					for i := Spans[num].Start; i < Spans[num].End; i++ {
						allGr = append(allGr, eachColumn[i]...)
					}
					//					departments = parseLine(departments, allGr, resFromReg, insertedGroups, subject, text, countSmall0, n, nextStr)
					num++
				}

				departments, insertedGroups = parseLine(departments, allGr, resFromReg, insertedGroups, subject, text, countSmall0-1, n, nextStr, is2Weeks)

				if countSmall0 > 0 {
					countSmall0--
					if countSmall0 != 0 {
						return
					}
					ind = ind + numberBeforeSmall0
					numberBeforeSmall0 = 0
					return
				}
				ind = ind + number
			}
			fmt.Println(ind, time, class, text, "\n")
		}
	})
	for _, val := range departments {
		fmt.Println(val.Number)
		fmt.Println(val.Lessons, "\n")
	}
}
