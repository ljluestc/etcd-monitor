package etcdinspectioncontroller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEtcdInspectionController_Comprehensive(t *testing.T) {
	t.Run("NewEtcdInspectionControllerCommand", func(t *testing.T) {
		cmd := NewEtcdInspectionControllerCommand(nil)
		assert.NotNil(t, cmd)
		assert.Equal(t, "inspection", cmd.Use)
		assert.Equal(t, "run inspection controller", cmd.Short)
	})

	t.Run("EtcdInspectionCommand structure", func(t *testing.T) {
		cc := &EtcdInspectionCommand{}
		assert.NotNil(t, cc)
		assert.Empty(t, cc.out)
		assert.Empty(t, cc.kubeconfig)
		assert.Empty(t, cc.masterURL)
		assert.Empty(t, cc.labelSelector)
		assert.Empty(t, cc.leaseLockName)
		assert.Empty(t, cc.leaseLockNamespace)
	})

	t.Run("Command flags", func(t *testing.T) {
		cmd := NewEtcdInspectionControllerCommand(nil)
		
		// Check that all expected flags exist
		assert.NotNil(t, cmd.Flag("kubeconfig"))
		assert.NotNil(t, cmd.Flag("master"))
		assert.NotNil(t, cmd.Flag("label-selector"))
		assert.NotNil(t, cmd.Flag("lease-lock-name"))
		assert.NotNil(t, cmd.Flag("lease-lock-namespace"))
	})

	t.Run("Command description", func(t *testing.T) {
		cmd := NewEtcdInspectionControllerCommand(nil)
		assert.Contains(t, cmd.Long, "inspection controller")
		assert.Contains(t, cmd.Long, "daemon")
		assert.Contains(t, cmd.Long, "etcdinspection")
	})
}