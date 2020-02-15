package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("Controller", func() {
	var node *corev1.Node
	Describe("Extracting InstanceID from Node", func() {
		BeforeEach(func() {
			node = &corev1.Node{Spec: corev1.NodeSpec{ProviderID: "aws:///us-east-1a/i-053d6b928c4e33e44"}}
		})
		It("Should extract InstanceID from a Node", func() {
			instanceID := ExtractInstanceIDFromNode(node)
			expected := "i-053d6b928c4e33e44"
			Expect(instanceID).To(Equal(expected))
		})
	})
})

// func TestMaterializeTagsForNodeFromMapping(t *testing.T) {
// 	node := &corev1.Node{ObjectMeta: metav1.ObjectMeta{
// 		Labels: map[string]string{"nodeType": "bla"},
// 	}}

// 	tm := []tagMapping{
// 		tagMapping{Key: "nodeType", Value: "bla", TagKey: "1", TagValue: "2"},
// 	}

// 	tags := materializeTagsForNodeFromMapping(node, tm)

// 	if taglen := len(tags); taglen != 1 {
// 		t.Errorf("Expected tag length of 1, got: %v", taglen)
// 	}

// 	if len(tags) == 1 && (tags[0].Key != "1" || tags[0].Value != "2") {
// 		t.Errorf("Expected tag Tag{Key: 1, Value:2, got %v", tags[0])
// 	}
// }k8s.io/apimachinery/pkg/api/errors"

// })
