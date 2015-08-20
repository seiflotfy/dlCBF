package dlcbf

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestDlcbf(t *testing.T) {
	dlcbf, _ := NewDlcbfForCapacity(1000000)
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

	count := dlcbf.GetCount()
	if float64(count)*100/235886 < 1 {
		t.Error("Expected error < 1 percent, got", float64(count)*100/235886)
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

	count = dlcbf.GetCount()
	if count != 0 {
		t.Error("Expected count == 0, got", count)
	}

}
