package shoesdb

import (
	"database/sql"
	"errors"
	"reflect"
	"strconv"
	"testing"
)

func Test_Conn(t *testing.T) {
	db := Conn()
	dbType := reflect.TypeOf(db)
	if dbType.String() != "*sql.DB" {
		t.Errorf("Expected type = *sql.DB, Actual type = %v\n", dbType.String())
	}
}

func Test_InsertShoes(t *testing.T) {
	type test struct {
		name        string
		expectedErr error
		expectedRes int
		names       []string
	}
	tests := []test{
		test{"Zero names.", ErrZeroShoes, 0, newShoes(0)},
		test{"One shoe.", nil, 1, newShoes(1)},
		test{"Five shoe.", nil, 5, newShoes(5)},
		test{"Fifty shoe.", nil, 50, newShoes(50)},
	}

	for _, te := range tests {
		t.Run(te.name, func(t *testing.T) {
			res, err := InsertShoes(new(mockExecer), te.names...)
			if res != te.expectedRes {
				t.Errorf("Expected Result = %v, Actual Result = %v\n", te.expectedRes, res)
			}
			if err != te.expectedErr {
				t.Errorf("Expected Error = %v, Actual Error = %v\n", te.expectedErr, err)
			}
		})
	}
}

func Test_InsertTrueToSizes(t *testing.T) {
	type test struct {
		name        string
		expectedErr error
		expectedRes int
		trueToSizes []int
	}
	tests := []test{
		test{"Zero trueToSizes.", ErrZeroTrueToSizes, 0, newTruetoSizes(0)},
		test{"One trueToSizes.", nil, 1, newTruetoSizes(1)},
		test{"Five trueToSizes.", nil, 5, newTruetoSizes(5)},
		test{"Fifty trueToSizes.", nil, 50, newTruetoSizes(50)},
	}

	for _, te := range tests {
		t.Run(te.name, func(t *testing.T) {
			res, err := InsertTrueToSizes(new(mockExecer), te.trueToSizes...)
			if res != te.expectedRes {
				t.Errorf("Expected Result = %v, Actual Result = %v\n", te.expectedRes, res)
			}
			if err != te.expectedErr {
				t.Errorf("Expected Error = %v, Actual Error = %v\n", te.expectedErr, err)
			}
		})
	}
}

func Test_SelectTrueToSize(t *testing.T) {
	type test struct {
		name        string
		expectedErr error
		identifier  interface{}
	}
	tests := []test{
		test{"Int identifier.", nil, 1},
		test{"String identifier.", nil, "shoe_1"},
		test{"Float64 identifier.", ErrIdentifierInvalidType, 1.1},
	}

	for _, te := range tests {
		t.Run(te.name, func(t *testing.T) {
			res, err := SelectTrueToSize(new(mockQueryer), te.identifier)
			if err != te.expectedErr {
				t.Errorf("Expected Error = %v, Actual Error = %v\n", te.expectedErr, err)
				t.SkipNow()
			}
			resValid := true
			for i, v := range res {
				if v != defTrueToSizeSet[i] {
					resValid = false
					break
				}
			}
			if !resValid {
				t.Errorf("Expected Res = %v, Actual Res = %v\n", defTrueToSizeSet, res)
			}
		})
	}
}

// newshoes creates a set of n shoes.
func newShoes(n int) []string {
	shoes := make([]string, n)
	for i, _ := range shoes {
		shoes[i] = "shoe_" + strconv.Itoa(i)
	}
	return shoes
}

// newTruetoSizes creates a set of n trueToSizes.
func newTruetoSizes(n int) []int {
	trueToSizes := make([]int, n)
	for i, _ := range trueToSizes {
		trueToSizes[i] = (i % 5) + 1
	}
	return trueToSizes
}

// A mock Execer for testing.
type mockExecer struct{}

func (m *mockExecer) Exec(query string, args ...interface{}) (sql.Result, error) {
	return mockSqlResult{0, int64(len(args))}, nil
}

// A mock sql.Result for testing.
type mockSqlResult struct {
	lastID  int64
	rowsAff int64
}

func (m mockSqlResult) LastInsertId() (int64, error) { return m.lastID, nil }
func (m mockSqlResult) RowsAffected() (int64, error) { return m.rowsAff, nil }

// a mock Queryer for testing.
type mockQueryer struct{}

func (m *mockQueryer) Query(query string, args ...interface{}) (Rower, error) {
	mr := new(mockRowerInt)
	mr.rows = make([]int, len(defTrueToSizeSet))
	copy(mr.rows, defTrueToSizeSet)
	return mr, nil
}

// defTrueToSizeSet is the default set of rows to be evaluated by mockRowerInt.
var defTrueToSizeSet = []int{1, 3, 4, 1}

// a mock Rower for testing.
type mockRowerInt struct {
	err     error
	rows    []int
	thisRow int
}

func (r *mockRowerInt) Err() error   { return r.err }
func (r *mockRowerInt) Close() error { return nil }
func (r *mockRowerInt) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	r.thisRow = r.rows[0]
	r.rows = r.rows[1:]
	return true
}

var (
	ErrDestLenNotEqualRowLen = errors.New("Number of destinations does not equal number of values.")
	ErrInvalidDestination    = errors.New("Destination type is invalid.")
)

func (r *mockRowerInt) Scan(iDest ...interface{}) error {
	if len(iDest) != 1 {
		r.err = ErrDestLenNotEqualRowLen
		return ErrDestLenNotEqualRowLen
	}

	dest, ok := iDest[0].(*int)
	if !ok {
		return ErrInvalidDestination
	}

	*dest = r.thisRow
	return nil
}
