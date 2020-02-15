package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/klogr"
	klog "k8s.io/klog/v2"
	controllers "sigs.k8s.io/controller-runtime"
)

// var tm []TagMapping = []TagMapping{
// 	TagMapping{
// 		Key: "nodeType", Value: "kinvey",
// 		TagKey: "workload_type", TagValue: "trusted",
// 	},
// 	TagMapping{
// 		Key: "nodeType", Value: "fsr",
// 		TagKey: "workload_type", TagValue: "untrusted",
// 	},
// }

var log = klogr.New().WithName("ec2-tag-controller")

var cfg Config

func main() {
	klog.InitFlags(nil)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	viper.BindPFlags(pflag.CommandLine)
	pflag.Parse()

	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	viper.AddConfigPath(".")    // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	viper.Unmarshal(&cfg)

	manager, err := controllers.NewManager(controllers.GetConfigOrDie(), controllers.Options{})
	if err != nil {
		log.Error(err, "could not create manager")
		os.Exit(1)
	}

	sess := session.Must(session.NewSession())
	ec2client := ec2.New(sess, &aws.Config{})

	err = controllers.NewControllerManagedBy(manager).For(&corev1.Node{}).Complete(&reconcileNode{
		client:    manager.GetClient(),
		tm:        cfg.Mappings,
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
