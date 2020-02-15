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
)

const cloudTagsAnnotationKey = "kinvey.com/cloudTags"

// reconcileNode reconciles Nodes
type reconcileNode struct {
	// client can be used to retrieve objects from the APIServer.
	client    client.Client
	log       logr.Logger
	tm        []TagMapping
	ec2client ec2iface.EC2API
}

// Implement reconcile.Reconciler so the controller can reconcile objects
var _ reconcile.Reconciler = &reconcileNode{}

func (r *reconcileNode) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// set up a convenient log object so we don't have to type request over and over again
	log := r.log.WithValues("request", request)

	// Fetch the ReplicaSet from the cache
	node := &corev1.Node{}
	err := r.client.Get(context.TODO(), request.NamespacedName, node)
	if errors.IsNotFound(err) {
		log.Error(nil, "Could not find Node")
		return reconcile.Result{}, nil
	}

	if err != nil {
		log.Error(err, "Could not fetch Node")
		return reconcile.Result{}, err
	}

	log.Info("Reconciling Node", "node name", node.Name)

	tags := materializeTagsForNodeFromMapping(node, r.tm)

	// Set the label if it is missing

	desiredJson, err := json.Marshal(tags)

	if err != nil {
		return reconcile.Result{}, err
	}

	if currentJson, ok := node.Annotations[cloudTagsAnnotationKey]; !ok || !bytes.Equal(desiredJson, []byte(currentJson)) {
		r.log.Info("Node needs tagging")
		var ec2Tags []*ec2.Tag
		for _, t := range tags {
			ec2Tags = append(ec2Tags, &ec2.Tag{Key: aws.String(t.Key), Value: aws.String(t.Value)})
		}

		instanceID := ExtractInstanceIDFromNode(node)

		cti := &ec2.CreateTagsInput{
			Resources: []*string{aws.String(instanceID)},
			Tags:      ec2Tags,
		}

		if _, err := r.ec2client.CreateTags(cti); err != nil {
			r.log.WithValues("instance_id", instanceID).Error(err, "Failed creating tags")
			return reconcile.Result{}, err
		}
		node.Annotations[cloudTagsAnnotationKey] = string(desiredJson)
	}

	err = r.client.Update(context.TODO(), node)
	if err != nil {
		log.Error(err, "Could not write Node")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func materializeTagsForNodeFromMapping(n *corev1.Node, tm []TagMapping) []Tag {
	var tl []Tag
	for _, t := range tm {
		if l, ok := n.Labels[t.Key]; ok {
			if l == t.Value {
				tl = append(tl, Tag{Key: t.TagKey, Value: t.TagValue})
			}
		}
	}
	return tl
}

func ExtractInstanceIDFromNode(n *corev1.Node) string {
	parts := strings.Split(n.Spec.ProviderID, "/")
	return parts[len(parts)-1]
}
