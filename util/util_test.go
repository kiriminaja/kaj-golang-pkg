package util

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInArray(t *testing.T) {
	t.Parallel()
	testCase := []struct {
		Case       int
		Value      interface{}
		Collection interface{}
		Valid      bool
	}{
		{
			Case:  1,
			Value: "test",
			Collection: []string{
				"test",
				"test123",
			},
			Valid: true,
		},
		{
			Case:  2,
			Value: 2,
			Collection: []string{
				"test",
				"test123",
			},
			Valid: false,
		},
		{
			Case:       3,
			Value:      int64(3),
			Collection: []int64{1, 2, 3},
			Valid:      true,
		},

		{
			Case:  4,
			Value: "case-4",
			Collection: map[string]string{
				"case-4": "case-4",
			},
			Valid: false,
		},
	}

	for _, v := range testCase {

		valid := InArray(v.Value, v.Collection)

		if valid == v.Valid {
			t.Logf("scenario #%v exptected %v, got %v", v.Case, v.Valid, valid)
		}

		if valid != v.Valid {
			t.Errorf("scenario #%v exptected %v, got %v", v.Case, v.Valid, valid)
		}
	}
}

func TestDumpToString(t *testing.T) {
	t.Parallel()
	if result := DumpToString("test"); result == "test" {
		t.Logf("expected: %v, got: %v", "test", result)
	} else {
		t.Fatalf("expected: %v, got: %v", "test", result)
	}

	if result := DumpToString(1); strings.Trim(result, "\n") == "1" {
		t.Logf("expected: %v, got: %s", "1", result)
	} else {
		t.Fatalf("expected: %v, got: %v", "1", result)
	}
}

func TestDebugPrint(t *testing.T) {
	DebugPrint("test print")
}

func TestEnvironmentTransform(t *testing.T) {
	t.Parallel()
	testCase := map[string]string{
		"production":  "prod",
		"staging":     "stg",
		"development": "dev",
		"prod":        "prod",
		"stg":         "stg",
		"dev":         "dev",
		"devs":        "",
	}

	for k, v := range testCase {
		assert.Equal(t, v, EnvironmentTransform(k))
	}
}

func TestStringToDate(t *testing.T) {
	t.Parallel()
	t.Run("parse string to date", func(t *testing.T) {
		expected := `19/09/2019 06:24:33`
		tm := StringToDate(expected)
		assert.Equal(t, expected, tm.Format(`02/01/2006 15:04:05`))
	})
}

func TestStringToDateE(t *testing.T) {

	t.Parallel()
	// ts := `19/09/2019 06:24:33`
	// fmt.Errorf("unable to parse date: %s", s)
	// layout := `2006-01-02 15:04:05.000`
	scenarios := []struct {
		StringDate    string
		FormatDate    string
		ExpectedError bool
	}{
		{StringDate: `19/09/2019 06:24:33 000`, ExpectedError: true, FormatDate: "02/01/2006 15:04:05"},
		{StringDate: "2019-09-19T19:19:40+07:00", ExpectedError: false, FormatDate: "2006-01-02T15:04:05Z07:00"},
		{StringDate: "2019-09-19T19:19:40", ExpectedError: false, FormatDate: "2006-01-02T15:04:05"},
		{StringDate: "Thu, 19 Sep 2019 19:19:40 +0700", ExpectedError: false, FormatDate: "Mon, 02 Jan 2006 15:04:05 -0700"},
		{StringDate: "Thu, 19 Sep 2019 19:19:40 WIB", ExpectedError: false, FormatDate: "Mon, 02 Jan 2006 15:04:05 MST"},
		{StringDate: "19 Sep 19 19:19 +0700", ExpectedError: false, FormatDate: "02 Jan 06 15:04 -0700"},
		{StringDate: "19 Sep 19 19:19 WIB", ExpectedError: false, FormatDate: "02 Jan 06 15:04 MST"},
		{StringDate: "Thursday, 19-Sep-19 19:19:40 WIB", ExpectedError: false, FormatDate: "Monday, 02-Jan-06 15:04:05 MST"},
		{StringDate: "Thu Sep 19 19:19:40 2019", ExpectedError: false, FormatDate: "Mon Jan _2 15:04:05 2006"},
		{StringDate: "Thu Sep 19 19:19:40 WIB 2019", ExpectedError: false, FormatDate: "Mon Jan _2 15:04:05 MST 2006"},
		{StringDate: "Thu Sep 19 19:19:40 +0700 2019", ExpectedError: false, FormatDate: "Mon Jan 02 15:04:05 -0700 2006"},
		{StringDate: "2019-09-19 19:19:40.376382 +0700 WIB", ExpectedError: false, FormatDate: "2006-01-02 15:04:05.999999999 -0700 MST"},
		{StringDate: "2019-09-19", ExpectedError: false, FormatDate: "2006-01-02"},
		{StringDate: "19 Sep 2019", ExpectedError: false, FormatDate: "02 Jan 2006"},
		{StringDate: "2019-09-19T19:19:40+0700", ExpectedError: false, FormatDate: "2006-01-02T15:04:05-0700"},
		{StringDate: "2019-09-19 19:19:40 +07:00", ExpectedError: false, FormatDate: "2006-01-02 15:04:05 -07:00"},
		{StringDate: "2019-09-19 19:19:40 +0700", ExpectedError: false, FormatDate: "2006-01-02 15:04:05 -0700"},
		{StringDate: "2019-09-19 19:19:40+07:00", ExpectedError: false, FormatDate: "2006-01-02 15:04:05Z07:00"},
		{StringDate: "2019-09-19 19:19:40+0700", ExpectedError: false, FormatDate: "2006-01-02 15:04:05Z0700"},
		{StringDate: "2019-09-19 19:19:40", ExpectedError: false, FormatDate: "2006-01-02 15:04:05"},
		{StringDate: "7:19PM", ExpectedError: false, FormatDate: "3:04PM"},
		{StringDate: "Sep 19 19:19:40", ExpectedError: false, FormatDate: "Jan _2 15:04:05"},
		{StringDate: "Sep 19 19:19:40.376", ExpectedError: false, FormatDate: "Jan _2 15:04:05.000"},
		{StringDate: "Sep 19 19:19:40.376382", ExpectedError: false, FormatDate: "Jan _2 15:04:05.000000"},
		{StringDate: "Sep 19 19:19:40.376382000", ExpectedError: false, FormatDate: "Jan _2 15:04:05.000000000"},
		{StringDate: "19/09/2019 19:19:40", ExpectedError: false, FormatDate: "02/01/2006 15:04:05"},
		{StringDate: "19/09/2019 19:19:40.376", ExpectedError: false, FormatDate: "02/01/2006 15:04:05.000"},
		{StringDate: "19/09/2019", ExpectedError: false, FormatDate: "02/01/2006"},
	}

	for _, sc := range scenarios {

		r, e := StringToDateE(sc.StringDate)

		t.Log(r.Format(sc.FormatDate), "-", e)

		if sc.ExpectedError {
			assert.Error(t, e)
		} else {
			assert.NoError(t, e)
		}

	}

}

const (
	success = "\u2713"
	failed  = "\u2717"
)

func TestGenerateReferenceID(t *testing.T) {
	t.Parallel()
	t.Log("Test generation of reference id")
	{
		var tester []string
		for i := 0; i < 20000; i++ {
			ref := GenerateReferenceID("")
			tester = append(tester, ref)
		}

		var isSame bool

		for i := 0; i < len(tester); i++ {

			for j := 0; j < len(tester); j++ {

				if i != j {
					isSame = tester[i] == tester[j]
				}

			}

		}

		if isSame {
			t.Errorf("%s expected no two numbers have the same value", failed)
		} else {
			t.Logf("%s expected no two numbers have the same value", success)
		}
	}
}

func TestGenerateRandomNumberString(t *testing.T) {
	t.Parallel()
	t.Log("Test generation of string number")
	{
		var tester []string
		for i := 0; i < 20000; i++ {
			ref := GenerateRandomNumberString(6)
			tester = append(tester, ref)
		}

		var isSame bool

		for i := 0; i < len(tester); i++ {

			for j := 0; j < len(tester); j++ {

				if i != j {
					isSame = tester[i] == tester[j]
				}

			}

		}

		if isSame {
			t.Errorf("%s expected no two numbers have the same value", failed)
		} else {
			t.Logf("%s expected no two numbers have the same value", success)
		}
	}
}

func TestGenerateAppID(t *testing.T) {
	t.Parallel()
	t.Log("Test generation of string number")
	{
		var tester []string
		for i := 0; i < 20000; i++ {
			ref := GenerateAppID("")
			tester = append(tester, ref)
		}

		var isSame bool

		for i := 0; i < len(tester); i++ {

			for j := 0; j < len(tester); j++ {

				if i != j {
					isSame = tester[i] == tester[j]
				}

			}

		}

		if isSame {
			t.Errorf("%s expected no two numbers have the same value", failed)
		} else {
			t.Logf("%s expected no two numbers have the same value", success)
		}
	}
}

func TestReplacer(t *testing.T) {
	t.Parallel()
	testCase := map[string]string{
		"testCase#1": "testCase_1",
		"testCase@1": "testCase-1",
	}

	rule := map[string]string{
		"#": "_",
		"@": "-",
	}

	for k, v := range testCase {
		assert.Equal(t, v, Replacer(rule, k))
	}
}

func TestIsSameType(t *testing.T) {
	assert.Equal(t, true, IsSameType(struct{}{}, struct{}{}))
	assert.Equal(t, false, IsSameType(struct{}{}, int64(1)))
}

func TestToString(t *testing.T) {
	t.Parallel()
	scenarios := []struct {
		Input    interface{}
		Expected string
	}{
		{
			Input:    int(1),
			Expected: "1",
		},
		{
			Input:    int8(2),
			Expected: "2",
		},
		{
			Input:    int16(3),
			Expected: "3",
		},
		{
			Input:    int32(4),
			Expected: "4",
		},
		{
			Input:    int64(5),
			Expected: "5",
		},
		{
			Input:    uint(6),
			Expected: "6",
		},
		{
			Input:    uint8(7),
			Expected: "7",
		},
		{
			Input:    uint16(8),
			Expected: "8",
		},
		{
			Input:    uint32(9),
			Expected: "9",
		},
		{
			Input:    uint64(10),
			Expected: "10",
		},
		{
			Input:    float32(11),
			Expected: "11",
		},
		{
			Input:    float64(12),
			Expected: "12",
		},
		{
			Input:    bool(true),
			Expected: "true",
		},
		{
			Input: struct {
			}{},
			Expected: "{}",
		},

		{
			Input:    string(`test`),
			Expected: "test",
		},
	}

	for _, sc := range scenarios {
		r := ToString(sc.Input)

		assert.Equal(t, sc.Expected, r)
	}
}
