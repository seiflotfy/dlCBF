package dlcbf

const bucketHeight = 8
const fingerprintBits = 13

type counter uint8
type fingerprint uint16

type f struct {
	count counter
	fp    fingerprint
}

type field struct {
	all uint16
	f   f
}

type bucketfp struct {
	bucketI     uint16
	fingerprint fingerprint
}

type bucket struct {
	fields [fingerprintBits]field
	count  uint8
}

type table struct {
	buckets []bucket
}

/*
Dlcbf is a struct representing a d-left Counting Bloom Filter
*/
type Dlcbf struct {
	tables []table
	d      uint8
	b      uint8
	bits   uint8
}

type loc struct {
	field field
}
