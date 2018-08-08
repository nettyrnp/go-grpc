package util

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
	"strings"
)

const (
	UK = "(+44)"
)

type Person struct {
	Id           int32
	Name         string
	Email        string
	MobileNumber string
}

func (p Person) ToString() string {
	return fmt.Sprintf("%v,%v,%v,%v", p.Id, p.Name, p.Email, p.MobileNumber)
}

func Ingest(reader *csv.Reader, offset, limit int) []Person {
	peopleMap := map[Person]bool{}
	var people = []Person{}
	var count = offset
	for count <= offset+limit {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		count++
		if count == 1 {
			continue
		}
		id, err := strconv.Atoi(line[0])
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		p := Person{
			Id:           int32(id),
			Name:         line[1],
			Email:        line[2],
			MobileNumber: normalizeAndAddPrefix(line[3]),
		}
		if !peopleMap[p] {
			peopleMap[p] = true
			people = append(people, p)
		}
	}
	sort.SliceStable(people, func(i, j int) bool {
		return people[i].Id < people[j].Id
	})
	return people
}

func normalizeAndAddPrefix(s string) string {
	s = strings.Replace(s, " ", "", -1)
	s = strings.Replace(s, "(", "", -1)
	s = strings.Replace(s, ")", "", -1)
	s = UK + s
	return s
}
