package cache_test

import (
	"testing"

	"github.com/bagmangood/cache"
)

func TestLRU(t *testing.T) {
	size := 5
	values := [][]string{
		[]string{"zero", "a"},
		[]string{"one", "b"},
		[]string{"two", "c"},
		[]string{"three", "d"},
		[]string{"four", "e"},
		[]string{"five", "f"},
		[]string{"six", "g"},
	}

	testCache := cache.NewLRU(size)

	result, err := testCache.Read("doesn't exist")
	if result != nil {
		t.Errorf("not found did not return nil")
	}
	if err != cache.NotFound {
		t.Errorf("expected a return of notfound, got %v instead", err)
	}

	for _, pair := range values {
		k, v := pair[0], pair[1]

		testCache.Write(k, v)

		result, err = testCache.Read(k)

		if err != nil {
			t.Errorf("errored immediately after write with %v", err.Error())
		}

		if result != v {
			t.Errorf("read key %v and expected %v, got %v", k, v, result)
		}
	}

	if testCache.Size() != size {
		t.Errorf("error in how big the cache is, expected %v but was %v", size, testCache.Size())
	}

	result, err = testCache.Read("zero")

	if err != cache.NotFound {
		t.Errorf(
			"expected bumped value %v to return NotFound, instead got value %v with error %v",
			"zero",
			result,
			err,
		)
	}

	result, err = testCache.Read("one")

	if err != cache.NotFound {
		t.Errorf(
			"expected bumped value %v to return NotFound, instead got value %v with error %v",
			"one",
			result,
			err,
		)
	}
}
