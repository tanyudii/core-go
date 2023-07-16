package errutil_test

import (
	"errors"
	"github.com/tanyudii/core-go/errutil"
	"testing"
)

func TestWrapper(t *testing.T) {
	var (
		ErrLvl1 = errors.New("lvl1")
		ErrLvl2 = errors.New("lvl2")
	)

	testCases := []struct {
		name   string
		err    error
		expect func(err error) bool
	}{
		{
			name: "Func_A",
			err:  ErrLvl1,
			expect: func(err error) bool {
				return errors.Is(err, ErrLvl1)
			},
		},
		{
			name: "Func_B",
			err:  ErrLvl2,
			expect: func(err error) bool {
				return errors.Is(err, ErrLvl2)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			wrapper := errutil.NewWrapper(tc.name)
			err := wrapper.Wrap(tc.err)
			if !tc.expect(err) {
				t.Errorf("expect %v, got %v", tc.err, err)
			}
		})
	}
}
