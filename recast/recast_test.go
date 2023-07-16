package recast_test

import (
	"encoding/json"
	"github.com/tanyudii/core-go/recast"
	"testing"
)

func TestNormal(t *testing.T) {
	type Entity struct {
		Serial string
		Name   string
		Email  string
	}

	type CreateRequest struct {
		Name  string
		Email string
	}

	type UpdateRequest struct {
		Serial string `json:"-"`
		Name   string
		Email  string
	}

	testCases := []struct {
		name   string
		entity interface{}
		req    interface{}
		expect interface{}
	}{
		{
			name:   "create test",
			entity: &Entity{},
			req: &CreateRequest{
				Name:  "name",
				Email: "email",
			},
			expect: &Entity{
				Name:  "name",
				Email: "email",
			},
		},
		{
			name: "update test",
			entity: &Entity{
				Serial: "current_serial",
				Name:   "current_name",
				Email:  "current_email",
			},
			req: &UpdateRequest{
				Serial: "serial",
				Name:   "name",
				Email:  "email",
			},
			expect: &Entity{
				Serial: "current_serial",
				Name:   "name",
				Email:  "email",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := recast.Recast(tc.req, tc.entity)
			if err != nil {
				t.Error(err)
			}

			entityMarshal, _ := json.Marshal(tc.entity)
			expectMarshal, _ := json.Marshal(tc.expect)
			if string(entityMarshal) != string(expectMarshal) {
				t.Errorf("expect %v, got %v", string(expectMarshal), string(entityMarshal))
			}
		})
	}
}

func TestDeep(t *testing.T) {
	type Entity struct {
		Serial   string
		Name     string
		Email    string
		Parent   *Entity
		Children []*Entity
	}

	type CreateRequest struct {
		Name     string
		Email    string
		Parent   *CreateRequest
		Children []*CreateRequest
	}

	type UpdateRequest struct {
		Serial   string `json:"-"`
		Name     string
		Email    string
		Parent   *UpdateRequest
		Children []*UpdateRequest
	}

	testCases := []struct {
		name   string
		entity interface{}
		req    interface{}
		expect interface{}
	}{
		{
			name:   "create test",
			entity: &Entity{},
			req: &CreateRequest{
				Name:  "name",
				Email: "email",
				Parent: &CreateRequest{
					Name:  "parent_name",
					Email: "parent_email",
				},
				Children: []*CreateRequest{
					{
						Name:  "child_name",
						Email: "child_email",
					},
				},
			},
			expect: &Entity{
				Name:  "name",
				Email: "email",
				Parent: &Entity{
					Name:  "parent_name",
					Email: "parent_email",
				},
				Children: []*Entity{
					{
						Name:  "child_name",
						Email: "child_email",
					},
				},
			},
		},
		{
			name: "update test",
			entity: &Entity{
				Serial: "current_serial",
				Name:   "current_name",
				Email:  "current_email",
				Parent: &Entity{
					Serial: "current_parent_serial",
					Name:   "current_parent_name",
					Email:  "current_parent_email",
				},
				Children: []*Entity{
					{
						Serial: "current_child_serial",
						Name:   "current_child_name",
						Email:  "current_child_email",
					},
				},
			},
			req: &UpdateRequest{
				Serial: "serial",
				Name:   "name",
				Email:  "email",
				Parent: &UpdateRequest{
					Serial: "parent_serial",
					Name:   "parent_name",
					Email:  "parent_email",
				},
				Children: []*UpdateRequest{
					{
						Serial: "child_serial",
						Name:   "child_name",
						Email:  "child_email",
					},
				},
			},
			expect: &Entity{
				Serial: "current_serial",
				Name:   "name",
				Email:  "email",
				Parent: &Entity{
					Serial: "current_parent_serial",
					Name:   "parent_name",
					Email:  "parent_email",
				},
				Children: []*Entity{
					{
						Serial: "current_child_serial",
						Name:   "child_name",
						Email:  "child_email",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := recast.Recast(tc.req, tc.entity)
			if err != nil {
				t.Error(err)
			}

			entityMarshal, _ := json.Marshal(tc.entity)
			expectMarshal, _ := json.Marshal(tc.expect)
			if string(entityMarshal) != string(expectMarshal) {
				t.Errorf("expect %v, got %v", string(expectMarshal), string(entityMarshal))
			}
		})
	}
}
