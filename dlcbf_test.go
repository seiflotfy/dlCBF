package dlcbf

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestDlcbf(t *testing.T) {
	dlcbf, _ := NewDlcbf(256, 256)
	fd, err := os.Open("/usr/share/dict/web2")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	scanner := bufio.NewScanner(fd)

	for scanner.Scan() {
		s := []byte(scanner.Text())
		dlcbf.Add(s)
	}

	fmt.Println(dlcbf.GetCount())
	fmt.Println("---------")
	fd, err = os.Open("/usr/share/dict/web2")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	scanner = bufio.NewScanner(fd)

	for scanner.Scan() {
		s := []byte(scanner.Text())
		dlcbf.Delete(s)
	}

	fmt.Println(dlcbf.GetCount())
	/*
		for i, table := range dlcbf.tables {
			for j, bucket := range table {
				if bucket.count > 0 {
					fmt.Println(i, j, bucket)
				}
			}
		}

		fd, err = os.Open("/usr/share/dict/web2")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		scanner = bufio.NewScanner(fd)

		for scanner.Scan() {
			s := []byte(scanner.Text())
			dlcbf.Delete(s)
		}

		fmt.Println(dlcbf.GetCount())
		for i, table := range dlcbf.tables {
			for j, bucket := range table {
				if bucket.count > 0 {
					fmt.Println(i, j, bucket)
				}
			}
		}
	*/
}
