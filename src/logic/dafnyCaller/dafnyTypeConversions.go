// Formally verified E-Voting using Dafny
// Copyright (C) 2025 Authors Gruppe EinS
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package votingSystemsCaller

import (
	"dafny"
	"errors"
)

func sliceInt_to_DafnySet(x []int) (res dafny.Set) {
	// convert []int to set<int>
	xInterface := make([]interface{}, len(x))
	for i, element := range x {
		xInterface[i] = dafny.IntOf(element)
	}
	xSet := dafny.SetOf(xInterface...)
	return xSet
}

func sliceSliceInt_to_dafnySeqSetInt(x [][]int) (res dafny.Sequence) {
	// converts [][]int into dafny Type seq<set<int>>
	xInterface := make([]interface{}, len(x))
	for i, innerSlice := range x {
		innerSliceInterface := make([]interface{}, len(innerSlice))
		for j, element := range innerSlice {
			innerSliceInterface[j] = dafny.IntOf(element)
		}

		xInterface[i] = dafny.SetOf(innerSliceInterface...)
	}
	isString := false
	xSeq := dafny.SeqFromArray(xInterface, isString)
	return xSeq
}

func sliceSliceInt_to_dafnySeqSeqInt(x [][]int) (res dafny.Sequence) {
	// converts [][]int into dafny Type seq<seq<int>>
	xInterface := make([]interface{}, len(x))
	isString := false
	for i, innerSlice := range x {
		innerSliceInterface := make([]interface{}, len(innerSlice))
		for j, element := range innerSlice {
			innerSliceInterface[j] = dafny.IntOf(element)
		}

		xInterface[i] = dafny.SeqFromArray(innerSliceInterface, isString)
	}
	xSeq := dafny.SeqFromArray(xInterface, isString)
	return xSeq
}

func dafnyMapInt2Int_to_mapInt2Int(dafnyMap dafny.Map) (map[int]int, error) {
	// convert dafny Type map<int, int> into map[int]int
	resultMap := make(map[int]int)
	var keyValueTupleSet dafny.Set = dafnyMap.Items()
	mapIter := keyValueTupleSet.Iterator()
	// Iterieren über (key, value) Tuple
	for {
		keyValueTupleUntyped, hasNext := mapIter()
		if !hasNext {
			break
		}

		// cast as Tuple
		var keyValueTuple dafny.Tuple
		keyValueTuple, ok := keyValueTupleUntyped.(dafny.Tuple)
		if !ok {
			err := errors.New("cant cast value to dafny.Tuple")
			return make(map[int]int), err
		}

		untypedKey := *(keyValueTuple.IndexInt(0)) // IndexInt returnt *interface{}
		key, ok := untypedKey.(dafny.Int)          // cast als dafny.Int
		if !ok {
			err := errors.New("cant cast key from DafnyMapvalue to dafny.Int")
			return make(map[int]int), err
		}

		untypedValue := *(keyValueTuple.IndexInt(1)) // IndexInt returnt *interface{}
		value, ok := untypedValue.(dafny.Int)        // cast als dafny.Int
		if !ok {
			err := errors.New("cant cast votesOfCandidate from DafnyMap to dafny.Int")
			return make(map[int]int), err
		}
		resultMap[int(key.Int32())] = int(value.Int32())
	}
	return resultMap, nil
}

func dafnySetInt_to_sliceInt(s dafny.Set) ([]int, error) {
	//convert dafny Type set<int> into []int
	resultSlice := make([]int, 0, s.CardinalityInt())

	setIter := s.Iterator()
	for {
		untypedElement, hasNext := setIter()
		if !hasNext {
			break
		}
		element, ok := untypedElement.(dafny.Int)
		if !ok {
			return make([]int, 0), errors.New("cant cast element from dafny Set to dafny.Int")
		}
		resultSlice = append(resultSlice, int(element.Int32()))
	}
	return resultSlice, nil
}

func sliceOfmapInt2Int_to_DafnySeqOfMapInt2Int(sliceOfMaps [](map[int]int)) dafny.Sequence {
	//convert [](map[int][int]) into dafny Type seq<map<int,int>>

	resSeqAsInterface := make([]interface{}, len(sliceOfMaps))
	for i, mapInSlice := range sliceOfMaps {

		dafnyMap := mapInt2Int_to_DafnyMapInt2Int(mapInSlice)
		resSeqAsInterface[i] = dafnyMap
	}
	isString := false
	resSeq := dafny.SeqFromArray(resSeqAsInterface, isString)
	return resSeq
}

func mapInt2Int_to_DafnyMapInt2Int(m map[int]int) (res dafny.Map) {
	// convert map[int]int into dafny map<int,int>

	mapBuilder := dafny.NewMapBuilder()
	for key, value := range m {
		dafnyIntKey := dafny.IntOf(key)
		dafnyIntValue := dafny.IntOf(value)
		mapBuilder.Add(dafnyIntKey, dafnyIntValue)
	}
	return mapBuilder.ToMap()
}
