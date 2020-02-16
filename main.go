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

var log = klogr.New().WithName("ec2-tag-controller")

var cfg Config

func main() {
	setupConfigAndFlags()
	err := viper.ReadInConfig()
	viper.Unmarshal(&cfg)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	manager, err := controllers.NewManager(controllers.GetConfigOrDie(), controllers.Options{})
	if err != nil {
		log.Error(err, "could not create manager")
		os.Exit(1)
	}

	viper.SetDefault("annotation_key", "kinvey.com/cloudTags")

	sess := session.Must(session.NewSession())
	ec2client := ec2.New(sess, &aws.Config{})

	err = controllers.NewControllerManagedBy(manager).For(&corev1.Node{}).Complete(&reconcileNode{
		client:                 manager.GetClient(),
		tm:                     cfg.Mappings,
		log:                    log,
		ec2client:              ec2client,
		cloudTagsAnnotationKey: viper.GetString("annotation_key"),
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

func setupConfigAndFlags() {

	klog.InitFlags(nil)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	viper.BindPFlags(pflag.CommandLine)
	pflag.Parse()

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
}
