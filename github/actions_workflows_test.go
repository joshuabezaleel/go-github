// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestActionsService_ListWorkflows(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/workflows", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":4,"workflows":[{"name":"workflow1","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"},{"name":"workflow2","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}]}`)
	})

	opt := &ListOptions{Page: 2, PerPage: 2}
	workflows, _, err := client.Actions.ListWorkflows(context.Background(), "o", "r", opt)
	if err != nil {
		t.Errorf("Actions.ListWorkflows returned error: %v", err)
	}

	want := &Workflows{
		TotalCount: 4,
		Workflows: []*Workflow{
			{Name: "workflow1", CreatedAt: Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)}, UpdatedAt: Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)}},
			{Name: "workflow", CreatedAt: Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)}, UpdatedAt: Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)}},
		},
	}
	if !reflect.DeepEqual(workflows, want) {
		t.Errorf("Actions.ListWorkflows returned %+v, want %+v", workflows, want)
	}
}

func TestActionsService_GetWorkflow(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/actions/workflows/72844", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":72844,"created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}`)
	})

	workflow, _, err := client.Actions.GetWorkflow(context.Background(), "o", "r", 72844)
	if err != nil {
		t.Errorf("Actions.GetWorkflow returned error: %v", err)
	}

	want := &Workflow{
		ID:        72844,
		CreatedAt: Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)},
		UpdatedAt: Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)},
	}
	if !reflect.DeepEqual(workflow, want) {
		t.Errorf("Actions.GetWorkflow returned %+v, want %+v", workflow, want)
	}
}
