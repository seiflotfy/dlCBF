# dlCBF
[![GoDoc](https://godoc.org/github.com/seiflotfy/dlCBF?status.svg)](https://godoc.org/github.com/seiflotfy/dlCBF)

A d-left Counting Bloom Filter (dlCBF) is simple hashing-based Counting Bloom Filter (CBF) alternative based on d-left hashing that offers the same functionality as a CBF, but uses less space, generally saving a factor of two or more.

For details about the algorithm and citations please use this article for now

["Bloom Filters via d-Left Hashing and Dynamic Bit Reassignment Extended Abstract" by Flavio Bonomi, Michael Mitzenmacher, Rina Panigrahy, Sushil Singh and George Varghese](http://www.eecs.harvard.edu/~michaelm/postscripts/aller2006.pdf)

##### Note: 
* This implmentation currently does not implement the counter feature. (coming soon)
* Insertions of items that already exists in the filter will fail (item can only exist once). 
* fingerprint size set to a static 16 Bit (also to be changed soon)

### Usage
```go
dlcbf, err := NewDlcbfForCapacity(1000000)
	
// Add Item
dlcbf.Add("jon snow is alive")

count := dlcbf.GetCount() // count >> 1

member := dlcbf.IsMember("jon snow is alive") // member = true
	
// Remove Item
dlcbf.Delete("jon snow is alive") // returns true

// Remove Item again
dlcbf.Delete("jon snow is alive") // returns false
	
count := dlcbf.GetCount() // count >> 0
```
