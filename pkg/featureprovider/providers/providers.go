package providers

import (
	// register consistency inspection feature/
	_ "github.com/etcd-monitor/taskmaster/pkg/featureprovider/providers/consistency"
	// register healthy inspection feature
	_ "github.com/etcd-monitor/taskmaster/pkg/featureprovider/providers/healthy"
	// register request inspection feature
	_ "github.com/etcd-monitor/taskmaster/pkg/featureprovider/providers/request"
	// register alarm inspection feature
	_ "github.com/etcd-monitor/taskmaster/pkg/featureprovider/providers/alarm"
)
