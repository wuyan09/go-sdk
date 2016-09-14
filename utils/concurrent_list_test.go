package utils

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

func aTestRemoveArray(t *testing.T) {
	a := []interface{}{"abc1", "abc2", "abc3", "abc4", "abc5", "abc6", "abc7"}

	a1 := removeArray(a, 0)

	if len(a1) != len(a)-1 {
		t.Errorf("Invaliate len:expect:%d,actual:%d", len(a)-1, len(a1))
	}

	if !reflect.DeepEqual(a[1:], a1) {
		t.Errorf("Not equsls, expect:%v,actual:%v", a[1:], a1)
	}
}

func aTestSyncList1(t *testing.T) {
	list := NewConcurrentList()

	if !list.Empty() {
		t.Error("Not emtpy")
	}
	n := 1000
	for i := 0; i < n; i++ {
		list.Add(fmt.Sprintf("abc_%d", i))
	}

	if list.Size() != n {
		t.Errorf("Invalidate size:%d", list.Size())
	}

	for i := list.Size() - 1; i >= 0; i-- {
		v := list.Remove(i)

		if v != fmt.Sprintf("abc_%d", i) {
			t.Errorf("Not equals.%v", v)
		}
	}

	if !list.Empty() || list.Size() != 0 {
		t.Error("Not emtpy")
	}
}

func aTestRemoveSyncList(t *testing.T) {
	list := NewConcurrentList()

	if !list.Empty() {
		t.Error("Not emtpy")
	}
	n := 1000
	for i := 0; i < n; i++ {
		list.Add(fmt.Sprintf("abc_%d", i))
	}

	result := list.RemoveFun(func(i int, v interface{}) bool {
		if str, ok := v.(string); ok && strings.HasPrefix(str, "abc_") {
			return true
		}

		return false
	})

	if result == nil || len(result) != n {
		t.Errorf("Invalidate remove fun result.%v", result)
	} else {
		for i, v := range result {
			expect := fmt.Sprintf("abc_%d", i)

			if v != expect {
				t.Errorf("Invalidate result,expect:%v,actual:%v", expect, v)
			}
		}
	}
}

func aTestSyncListCon(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU())

	list := NewConcurrentList()

	f := func(group *sync.WaitGroup) {

		for i := 0; i < 10; i++ {
			n1 := rand.Intn(1000)
			n2 := rand.Intn(1000)
			str := fmt.Sprintf("%d_eee_%d", n1, n2)

			list.Add(str)

			t := rand.Intn(50)
			time.Sleep(time.Duration(t+1) * time.Millisecond)
		}

		group.Done()
	}

	group := new(sync.WaitGroup)

	n := 10000

	group.Add(n)

	start := time.Now()

	for i := 0; i < n; i++ {
		go f(group)
	}

	group.Wait()

	fmt.Printf("Used time:%d ms,list size=%d\n", time.Now().Sub(start)/1e6, list.Size())

	time.Sleep(5 * time.Minute)

}
