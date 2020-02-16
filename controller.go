package main

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"

	"github.com/go-logr/logr"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

// reconcileNode reconciles Nodes
type reconcileNode struct {
	// client can be used to retrieve objects from the APIServer.
	client                 client.Client
	log                    logr.Logger
	tm                     []Mapping
	ec2client              ec2iface.EC2API
	cloudTagsAnnotationKey string
}

// Implement reconcile.Reconciler so the controller can reconcile objects
var _ reconcile.Reconciler = &reconcileNode{}

func (r *reconcileNode) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// set up a convenient log object so we don't have to type request over and over again
	log := r.log.WithValues("node_name", request.NamespacedName.Name)

	node, err := r.GetNodeByName(request.NamespacedName)

	if errors.IsNotFound(err) {
		log.Error(nil, "Could not find Node")
		return reconcile.Result{}, nil
	}

	if err != nil {
		log.Error(err, "Could not fetch Node")
		return reconcile.Result{}, err
	}

	log.Info("Reconciling Node")

	tags := materializeTagsForNodeFromMapping(node, r.tm)
	desiredJSON, err := json.Marshal(tags)

	if err != nil {
		return reconcile.Result{}, err
	}

	if currentJSON, ok := node.Annotations[r.cloudTagsAnnotationKey]; ok && bytes.Equal(desiredJSON, []byte(currentJSON)) {
		return reconcile.Result{}, nil
	}

	log.Info("Node needs tagging")

	err = r.CreateTagsForNode(tags, node)
	if err != nil {
		log.Error(err, "Could not create tags")
		return reconcile.Result{}, err
	}

	node.Annotations[r.cloudTagsAnnotationKey] = string(desiredJSON)
	err = r.client.Update(context.TODO(), node)
	if err != nil {
		log.Error(err, "Could not write Node")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *reconcileNode) GetNodeByName(name types.NamespacedName) (*corev1.Node, error) {
	// Fetch the ReplicaSet from the cache
	node := &corev1.Node{}
	err := r.client.Get(context.TODO(), name, node)
	return node, err
}

func (r *reconcileNode) CreateTagsForNode(tags []*ec2.Tag, node *corev1.Node) error {
	instanceID := extractInstanceIDFromNode(node)

	cti := &ec2.CreateTagsInput{
		Resources: []*string{aws.String(instanceID)},
		Tags:      tags,
	}

	_, err := r.ec2client.CreateTags(cti)
	return err
}

func materializeTagsForNodeFromMapping(n *corev1.Node, tm []Mapping) []*ec2.Tag {
	var tl []*ec2.Tag
	for _, t := range tm {
		if l, ok := n.Labels[t.Key]; ok {
			if l == t.Value {
				tl = append(tl, &ec2.Tag{Key: aws.String(t.TagKey), Value: aws.String(t.TagValue)})
			}
		}
	}
	return tl
}

// extractInstanceIDFromNode ...
func extractInstanceIDFromNode(n *corev1.Node) string {
	parts := strings.Split(n.Spec.ProviderID, "/")
	return parts[len(parts)-1]
}
