// Copyright 2019 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package main

import "testing"

func next() []Datum {
	return nil
}

type Operator interface {
	next() []Datum
}

type Datum interface {
	Mul(Datum) Datum
}

type Int struct {
	int
}

func (i Int) Mul(datum Datum) Datum {
	arg := datum.(Int).int
	return Int{int: i.int * arg}
}

type Float64 struct {
	float64
}

func (f Float64) Mul(datum Datum) Datum {
	arg := datum.(Float64).float64
	return Float64{float64: f.float64 * arg}
}

type mulOperator struct {
	input             Operator
	arg               Datum
	columnsToMultiply []int
}

func (m mulOperator) next() []Datum {
	row := m.input.next()
	if row == nil {
		return nil
	}
	for _, c := range m.columnsToMultiply {
		row[c] = row[c].Mul(m.arg)
	}
	return row
}

type tableReader struct {
	curIdx int
	rows   [][]Datum
}

func (t *tableReader) next() []Datum {
	if len(t.rows) == 0 {
		return nil
	}
	return t.rows[t.curIdx]
}

func (t *tableReader) reset() {
	t.curIdx = 0
}

// makeInput creates numRows rows of numCols each of the given type. For each
// row, all of its columns will be its index (zero-indexed).
func makeInput(numRows int, numCols int, t Datum) [][]Datum {
	result := make([][]Datum, numRows)
	for i := range result {
		result[i] = make([]Datum, numCols)
	}
	switch t.(type) {
	case Int:
		for i := 0; i < numRows; i++ {
			for j := 0; j < numCols; j++ {
				result[i][j] = Int{int: i}
			}
		}
	case Float64:
		for i := 0; i < numRows; i++ {
			for j := 0; j < numCols; j++ {
				result[i][j] = Float64{float64: float64(i)}
			}
		}
	default:
		panic("unhandled type")
	}
	return result
}

func BenchmarkInterface(b *testing.B) {

}
