package dlcbf

import (
	"encoding/binary"
	"errors"
	"hash/fnv"
	"math"
)

const bucketHeight = 8

type fingerprint uint16

type target struct {
	bucketIndex uint
	fingerprint fingerprint
}

type bucket struct {
	entries [bucketHeight]fingerprint
	count   uint8
}

type table []bucket

/*
Dlcbf is a struct representing a d-left Counting Bloom Filter
*/
type Dlcbf struct {
	tables     []table
	numTables  uint
	numBuckets uint
}

/*
NewDlcbf returns a newly created Dlcbf
*/
func NewDlcbf(numTables uint, numBuckets uint) (*Dlcbf, error) {

	if numBuckets < numTables {
		return nil, errors.New("numBuckets has to be greater than numTables")
	}

	dlcbf := &Dlcbf{
		numTables:  numTables,
		numBuckets: numBuckets,
		tables:     make([]table, numTables, numTables),
	}

	for i := range dlcbf.tables {
		dlcbf.tables[i] = make(table, numBuckets, numBuckets)
	}

	return dlcbf, nil
}

/*
NewDlcbfForCapacity returns a newly created Dlcbf for a given max Capacity
*/
func NewDlcbfForCapacity(capacity uint) (*Dlcbf, error) {
	t := capacity / (4096 * bucketHeight)
	return NewDlcbf(t, 4096)
}

func (dlcbf *Dlcbf) getTargets(data []byte) []target {
	hasher := fnv.New64a()
	hasher.Write(data)
	fp := hasher.Sum(nil)
	hsum := hasher.Sum64()

	h1 := uint32(hsum & 0xffffffff)
	h2 := uint32((hsum >> 32) & 0xffffffff)

	indices := make([]uint, dlcbf.numTables, dlcbf.numTables)
	for i := uint(0); i < dlcbf.numTables; i++ {
		saltedHash := uint((h1 + uint32(i)*h2))
		indices[i] = (saltedHash % dlcbf.numBuckets)
	}

	targets := make([]target, dlcbf.numTables, dlcbf.numTables)
	for i := uint(0); i < dlcbf.numTables; i++ {
		targets[i] = target{
			bucketIndex: uint(indices[i]),
			fingerprint: fingerprint(binary.LittleEndian.Uint16(fp)),
		}
	}
	return targets
}

/*
Add data to filter return true if insertion was successful,
return false if data already in filter or size limit was exceeeded
*/
func (dlcbf *Dlcbf) Add(data []byte) bool {
	targets := dlcbf.getTargets(data)

	_, _, target := dlcbf.lookup(targets)
	if target != nil {
		return false
	}

	minCount := uint8(math.MaxUint8)
	tableI := uint(0)

	for i, target := range targets {
		tmpCount := dlcbf.tables[i][target.bucketIndex].count
		if tmpCount < minCount && tmpCount < bucketHeight {
			minCount = dlcbf.tables[i][target.bucketIndex].count
			tableI = uint(i)
		}
	}

	if minCount == uint8(math.MaxUint8) {
		return false
	}
	bucket := &dlcbf.tables[tableI][targets[tableI].bucketIndex]
	bucket.entries[minCount] = targets[tableI].fingerprint
	bucket.count++
	return true
}

/*
Delete data to filter return true if deletion was successful,
return false if data not in filter
*/
func (dlcbf *Dlcbf) Delete(data []byte) bool {
	deleted := false
	targets := dlcbf.getTargets(data)
	for i, target := range targets {
		for j, fp := range dlcbf.tables[i][target.bucketIndex].entries {
			if fp == target.fingerprint {
				if dlcbf.tables[i][target.bucketIndex].count == 0 {
					continue
				}
				dlcbf.tables[i][target.bucketIndex].count--
				k := 0
				for l, fp := range dlcbf.tables[i][target.bucketIndex].entries {
					if j == l {
						continue
					}
					dlcbf.tables[i][target.bucketIndex].entries[k] = fp
					k++
				}
				lastindex := dlcbf.tables[i][target.bucketIndex].count
				dlcbf.tables[i][target.bucketIndex].entries[lastindex] = 0
				deleted = true
			}
		}
	}
	return deleted
}

func (dlcbf *Dlcbf) lookup(targets []target) (uint, uint, *target) {
	for i, target := range targets {
		for j, fp := range dlcbf.tables[i][target.bucketIndex].entries {
			if fp == target.fingerprint {
				return uint(i), uint(j), &target
			}
		}
	}
	return 0, 0, nil
}

/*
IsMember return true if data is in filter
*/
func (dlcbf *Dlcbf) IsMember(data []byte) bool {
	targets := dlcbf.getTargets(data)
	_, _, bfp := dlcbf.lookup(targets)
	return bfp != nil
}

/*
GetCount return cardinlaity count of current filter
*/
func (dlcbf *Dlcbf) GetCount() uint {
	count := uint(0)
	for _, table := range dlcbf.tables {
		for _, bucket := range table {
			count += uint(bucket.count)
		}
	}
	return count
}
