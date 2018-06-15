package nursetags

import (
	"fmt"
	"testing"
)

const ErrorMismatch = `Result mismatch
expect: %s
result: %s

expect - result: %s
result - result: %s`

func init() {
	for i := 1; i <= 100; i++ {
		tags := make([]string, 0)
		for j := 1; j <= 10; j++ {
			if i%j == 0 {
				tags = append(tags, fmt.Sprintf("divisor:%01d", j))
			}
		}
		post := Post{
			Tags:      tags,
			Namespace: "test",
			External:  true,
			ID:        fmt.Sprintf("number:%03d", i),
			Value:     fmt.Sprintf("%02d", i),
		}

		databasePostCreate(post)
	}
}

func TestParseSimpleAndNot(t *testing.T) {
	expect := NewSetFromSlice([]string{
		"number:006",
		"number:018",
		"number:030",
		"number:042",
		"number:054",
		"number:066",
		"number:078",
		"number:090",
	})
	query := map[string]interface{}{
		"and": []interface{}{"divisor:2", "divisor:3"},
		"not": []interface{}{"divisor:4"},
	}
	result := parseQuery(query)
	if !expect.Equal(result) {
		t.Errorf(
			ErrorMismatch,
			expect,
			result,
			expect.Difference(result),
			result.Difference(expect),
		)
	}
}

func TestParseSimpleAndOr(t *testing.T) {
	expect := NewSetFromSlice([]string{
		"number:012",
		"number:024",
		"number:036",
		"number:048",
		"number:060",
		"number:072",
		"number:084",
		"number:096",
	})
	query := map[string]interface{}{
		"and": []interface{}{"divisor:3", "divisor:4"},
		"or":  []interface{}{"divisor:2"},
	}
	result := parseQuery(query)
	if !expect.Equal(result) {
		t.Errorf(
			ErrorMismatch,
			expect,
			result,
			expect.Difference(result),
			result.Difference(expect),
		)
	}
}

func TestParseSimpleAnd(t *testing.T) {
	expect := NewSetFromSlice([]string{
		"number:010",
		"number:020",
		"number:030",
		"number:040",
		"number:050",
		"number:060",
		"number:070",
		"number:080",
		"number:090",
		"number:100",
	})
	query := map[string]interface{}{
		"and": []interface{}{"divisor:2", "divisor:10"},
	}
	result := parseQuery(query)
	if !expect.Equal(result) {
		t.Errorf(
			ErrorMismatch,
			expect,
			result,
			expect.Difference(result),
			result.Difference(expect),
		)
	}
}

func TestParseSimpleOr(t *testing.T) {
	expect := NewSetFromSlice([]string{
		"number:002",
		"number:004",
		"number:005",
		"number:006",
		"number:008",
		"number:010",
		"number:012",
		"number:014",
		"number:015",
		"number:016",
		"number:018",
		"number:020",
		"number:022",
		"number:024",
		"number:025",
		"number:026",
		"number:028",
		"number:030",
		"number:032",
		"number:034",
		"number:035",
		"number:036",
		"number:038",
		"number:040",
		"number:042",
		"number:044",
		"number:045",
		"number:046",
		"number:048",
		"number:050",
		"number:052",
		"number:054",
		"number:055",
		"number:056",
		"number:058",
		"number:060",
		"number:062",
		"number:064",
		"number:065",
		"number:066",
		"number:068",
		"number:070",
		"number:072",
		"number:074",
		"number:075",
		"number:076",
		"number:078",
		"number:080",
		"number:082",
		"number:084",
		"number:085",
		"number:086",
		"number:088",
		"number:090",
		"number:092",
		"number:094",
		"number:095",
		"number:096",
		"number:098",
		"number:100",
	})
	query := map[string]interface{}{
		"or": []interface{}{"divisor:2", "divisor:5"},
	}
	result := parseQuery(query)
	if !expect.Equal(result) {
		t.Errorf(
			ErrorMismatch,
			expect,
			result,
			expect.Difference(result),
			result.Difference(expect),
		)
	}
}

func TestParseSimpleString(t *testing.T) {
	expect := NewSetFromSlice([]string{
		"number:010",
		"number:020",
		"number:030",
		"number:040",
		"number:050",
		"number:060",
		"number:070",
		"number:080",
		"number:090",
		"number:100",
	})
	query := "divisor:10"
	result := parseQuery(query)
	if !expect.Equal(result) {
		t.Errorf(
			ErrorMismatch,
			expect,
			result,
			expect.Difference(result),
			result.Difference(expect),
		)
	}
}
