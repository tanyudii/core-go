package errutil_test

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tanyudii/core-go/errutil"
	"testing"
)

func TestWrap(t *testing.T) {
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

func TestWrapMsg(t *testing.T) {
	testCases := []struct {
		name      string
		fn        string
		err       error
		msg       []string
		expectMsg string
	}{
		{
			name:      "EmptyMsg",
			err:       errors.New("lvl1"),
			msg:       nil,
			expectMsg: "EmptyMsg: lvl1",
		},
		{
			name:      "WithMsg1",
			err:       errors.New("lvl1"),
			msg:       []string{"msg1"},
			expectMsg: "WithMsg1: msg1: lvl1",
		},
		{
			name:      "WithMsg2",
			err:       errors.New("lvl1"),
			msg:       []string{"msg1", "msg2"},
			expectMsg: "WithMsg2: msg1: msg2: lvl1",
		},
		{
			name:      "WithMsgAndFn",
			fn:        "Fn Name",
			err:       errors.New("lvl1"),
			msg:       []string{"msg1", "msg2"},
			expectMsg: "WithMsgAndFn: Fn Name: msg1: msg2: lvl1",
		},
		{
			name:      "WithMsgAndFn And Wrapped",
			fn:        "Fn Name",
			err:       fmt.Errorf("%w: %w", errors.New("lvl1"), errors.New("lvl1.1")),
			msg:       []string{"msg1", "msg2"},
			expectMsg: "WithMsgAndFn And Wrapped: Fn Name: msg1: msg2: lvl1: lvl1.1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			wrapper := errutil.NewWrapper(tc.name).SetFn(tc.fn)
			err := wrapper.Wrap(tc.err, tc.msg...)
			assert.Equal(t, tc.expectMsg, err.Error())
		})
	}
}
