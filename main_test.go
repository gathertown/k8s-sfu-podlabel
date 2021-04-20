package main

import "testing"

func TestLabelGroups(t *testing.T) {
	type test struct {
		pods  int
		index int
		label string
	}
	tests := []test{
		{pods: 10, index: 1, label: "group1"},
		{pods: 10, index: 4, label: "group2"},
		{pods: 10, index: 6, label: "group3"},
		{pods: 10, index: 9, label: "group4"},
		{pods: 75, index: 7, label: "group1"},
		{pods: 75, index: 10, label: "group2"},
		{pods: 75, index: 30, label: "group3"},
		{pods: 75, index: 60, label: "group4"},
		{pods: 250, index: 1, label: "group1"},
		{pods: 250, index: 24, label: "group1"},
		{pods: 250, index: 30, label: "group2"},
		{pods: 250, index: 150, label: "group3"},
		{pods: 250, index: 201, label: "group4"},
	}

	for _, test := range tests {
		label := getLabelGroup(test.pods, test.index)

		if label != test.label {
			t.Fatalf("For index=%d,pods=%d expected: %v, got: %v", test.index, test.pods, label, test.label)
		}

	}

}
