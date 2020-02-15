package main

import (
	"flag"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	corev1 "k8s.io/api/core/v1"
	controllers "sigs.k8s.io/controller-runtime"

	"k8s.io/klog/klogr"
	klog "k8s.io/klog/v2"
)

var tm []tagMapping = []tagMapping{
	tagMapping{
		Key: "nodeType", Value: "kinvey",
		TagKey: "workload_type", TagValue: "trusted",
	},
	tagMapping{
		Key: "nodeType", Value: "fsr",
		TagKey: "workload_type", TagValue: "untrusted",
	},
}

var log = klogr.New().WithName("ec2-tag-controller")

func main() {
	klog.InitFlags(nil)
	flag.Parse()
	manager, err := controllers.NewManager(controllers.GetConfigOrDie(), controllers.Options{})
	if err != nil {
		log.Error(err, "could not create manager")
		os.Exit(1)
	}

	sess := session.Must(session.NewSession())
	ec2client := ec2.New(sess, &aws.Config{})

	err = controllers.NewControllerManagedBy(manager).For(&corev1.Node{}).Complete(&reconcileNode{
		client:    manager.GetClient(),
		tm:        tm,
		log:       log,
		ec2client: ec2client,
	})

	if err != nil {
		log.Error(err, "could not create controller")
		os.Exit(1)
	}

	if err := manager.Start(controllers.SetupSignalHandler()); err != nil {
		log.Error(err, "could not start manager")
		os.Exit(1)
	}

}
