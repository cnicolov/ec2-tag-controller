package main

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestExtractInstanceIDFromNode(t *testing.T) {
	node := &corev1.Node{Spec: corev1.NodeSpec{ProviderID: "aws:///us-east-1a/i-053d6b928c4e33e44"}}
	instanceID := extractInstanceIDFromNode(node)
	expected := "i-053d6b928c4e33e44"
	if instanceID != expected {
		t.Errorf("Expected %v to equal %v", instanceID, expected)
	}
}

func TestMaterializeTagsForNodeFromMapping(t *testing.T) {
	node := &corev1.Node{ObjectMeta: metav1.ObjectMeta{
		Labels: map[string]string{"nodeType": "bla"},
	}}

	tm := []tagMapping{
		tagMapping{Key: "nodeType", Value: "bla", TagKey: "1", TagValue: "2"},
	}

	tags := materializeTagsForNodeFromMapping(node, tm)

	if taglen := len(tags); taglen != 1 {
		t.Errorf("Expected tag length of 1, got: %v", taglen)
	}

	if len(tags) == 1 && (tags[0].Key != "1" || tags[0].Value != "2") {
		t.Errorf("Expected tag Tag{Key: 1, Value:2, got %v", tags[0])
	}
}
