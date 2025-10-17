package etcdinspectioncontroller

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	klog "k8s.io/klog/v2"

	"github.com/etcd-monitor/taskmaster/pkg/controllers/etcdinspection"
	"github.com/etcd-monitor/taskmaster/pkg/controllers/util"
	"github.com/etcd-monitor/taskmaster/pkg/k8s"
	"github.com/etcd-monitor/taskmaster/pkg/signals"
)

type EtcdInspectionCommand struct {
	out                io.Writer
	kubeconfig         string
	masterURL          string
	labelSelector      string
	leaseLockName      string
	leaseLockNamespace string
}

// NewEtcdInspectionControllerCommand creates a *cobra.Command object with default parameters
func NewEtcdInspectionControllerCommand(out io.Writer) *cobra.Command {
	cc := &EtcdInspectionCommand{out: out}
	cmd := &cobra.Command{
		Use:   "inspection",
		Short: "run inspection controller",
		Long: `The inspection controller is a daemon, it will watches the changes of etcdinspection resources
 through the apiserver and makes changes attempting to move the current state towards the desired state.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().VisitAll(func(flag *pflag.Flag) {
				klog.V(1).Infof("FLAG: --%s=%q", flag.Name, flag.Value)
			})
			if err := cc.run(); err != nil {
				return err
			}
			return nil
		},
	}
	// Add command flags
	cmd.Flags().StringVar(&cc.kubeconfig, "kubeconfig", "", "Path to a kubeconfig file")
	cmd.Flags().StringVar(&cc.masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig.")
	cmd.Flags().StringVar(&cc.labelSelector, "label-selector", "", "Label selector for filtering resources")
	cmd.Flags().StringVar(&cc.leaseLockName, "lease-lock-name", "", "Name of the lease lock")
	cmd.Flags().StringVar(&cc.leaseLockNamespace, "lease-lock-namespace", "", "Namespace of the lease lock")
	return cmd
}

func (c *EtcdInspectionCommand) run() error {
	// Set up signals to gracefully shutdown
	stopCh := signals.SetupSignalHandler()

	cfg, err := k8s.GetClientConfig(c.kubeconfig)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, etcdClient, kubeInformerFactory, etcdInformerFactory, err := k8s.GenerateInformer(cfg, c.labelSelector)
	if err != nil {
		klog.Fatalf("Error building clientset: %s", err.Error())
	}

	controller := etcdinspection.NewController(
		kubeClient,
		etcdClient,
		kubeInformerFactory,
		etcdInformerFactory,
		util.NewSimpleClientBuilder(c.kubeconfig),
	)

	// Start informer factories
	kubeInformerFactory.Start(stopCh)
	etcdInformerFactory.Start(stopCh)

	// Leader election setup
	if c.leaseLockName != "" && c.leaseLockNamespace != "" {
		id := string(uuid.NewUUID())
		lock := &resourcelock.LeaseLock{
			LeaseMeta: metav1.ObjectMeta{
				Name:      c.leaseLockName,
				Namespace: c.leaseLockNamespace,
			},
			Client: kubeClient.CoordinationV1(),
			LockConfig: resourcelock.ResourceLockConfig{
				Identity: id,
			},
		}

		leaderelection.RunOrDie(context.TODO(), leaderelection.LeaderElectionConfig{
			Lock:            lock,
			ReleaseOnCancel: true,
			LeaseDuration:   60 * time.Second,
			RenewDeadline:   15 * time.Second,
			RetryPeriod:     5 * time.Second,
			Callbacks: leaderelection.LeaderCallbacks{
				OnStartedLeading: func(ctx context.Context) {
					if err = controller.Run(2, stopCh); err != nil {
						klog.Fatalf("Error running controller: %s", err.Error())
					}
				},
				OnStoppedLeading: func() {
					klog.Infof("leader lost: %s", id)
					os.Exit(0)
				},
				OnNewLeader: func(identity string) {
					if identity == id {
						return
					}
					klog.Infof("new leader elected: %s", identity)
				},
			},
		})
	} else {
		if err = controller.Run(2, stopCh); err != nil {
			klog.Fatalf("Error running controller: %s", err.Error())
		}
	}

	<-stopCh
	return nil
}
