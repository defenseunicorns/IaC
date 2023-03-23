package test_test

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/terraform"
	teststructure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os/exec"
	"testing"
)

// To run this test, we first have to apply with the EKS public endpoint on, then apply again to turn the endpoint off. When destroying we need to do the opposite.
func TestExamplesCompleteSecure(t *testing.T) {
	t.Parallel()
	tempFolder := teststructure.CopyTerraformFolderToTemp(t, "../..", "examples/complete")
	terraformOptionsNoTargets := &terraform.Options{
		TerraformDir: tempFolder,
		Upgrade:      false,
		VarFiles: []string{
			"fixtures.common.tfvars",
			"fixtures.secure.tfvars",
		},
	}
	terraformOptionsWithVPCAndBastionTargets := &terraform.Options{
		TerraformDir: tempFolder,
		Upgrade:      false,
		VarFiles: []string{
			"fixtures.common.tfvars",
			"fixtures.secure.tfvars",
		},
		Targets: []string{
			"module.vpc",
			"module.bastion",
		},
	}
	terraformOptionsWithEKSTarget := &terraform.Options{
		TerraformDir: tempFolder,
		Upgrade:      false,
		VarFiles: []string{
			"fixtures.common.tfvars",
			"fixtures.secure.tfvars",
		},
		Targets: []string{
			"module.eks",
		},
	}
	defer teardownTestExamplesCompleteSecure(t, terraformOptionsNoTargets, terraformOptionsWithEKSTarget)
	setupTestExamplesCompleteSecure(t, terraformOptionsNoTargets, terraformOptionsWithVPCAndBastionTargets)
}

func setupTestExamplesCompleteSecure(t *testing.T, terraformOptionsNoTargets *terraform.Options, terraformOptionsWithVPCAndBastionTargets *terraform.Options) {
	t.Helper()
	teststructure.RunTestStage(t, "SETUP", func() {
		terraform.InitAndApply(t, terraformOptionsWithVPCAndBastionTargets)
		bastionInstanceID := terraform.Output(t, terraformOptionsWithVPCAndBastionTargets, "bastion_instance_id")
		//nolint:godox
		// TODO: Figure out how to parse the input variables to get the bastion password rather than having to hardcode it
		bastionPassword := "my-password"
		vpcCidr := terraform.Output(t, terraformOptionsWithVPCAndBastionTargets, "vpc_cidr")
		err := applyWithSshuttle(t, bastionInstanceID, bastionPassword, vpcCidr, terraformOptionsNoTargets)
		require.NoError(t, err)
	})
}

func teardownTestExamplesCompleteSecure(t *testing.T, terraformOptionsNoTargets *terraform.Options, terraformOptionsWithEKSTarget *terraform.Options) {
	t.Helper()
	teststructure.RunTestStage(t, "TEARDOWN", func() {
		bastionInstanceID := terraform.Output(t, terraformOptionsWithEKSTarget, "bastion_instance_id")
		//nolint:godox
		// TODO: Figure out how to parse the input variables to get the bastion password rather than having to hardcode it
		bastionPassword := "my-password"
		vpcCidr := terraform.Output(t, terraformOptionsWithEKSTarget, "vpc_cidr")
		err := destroyWithSshuttle(t, bastionInstanceID, bastionPassword, vpcCidr, terraformOptionsWithEKSTarget)
		assert.NoError(t, err)
		terraform.Destroy(t, terraformOptionsNoTargets)
	})
}

func applyWithSshuttle(t *testing.T, bastionInstanceID string, bastionPassword string, vpcCidr string, terraformOptions *terraform.Options) error {
	t.Helper()
	cmd, err := runSshuttleInBackground(t, bastionInstanceID, bastionPassword, vpcCidr)
	if err != nil {
		return err
	}
	defer func(t *testing.T, cmd *exec.Cmd) {
		t.Helper()
		err := stopSshuttle(t, cmd)
		require.NoError(t, err)
	}(t, cmd)
	terraform.Apply(t, terraformOptions)
	return nil
}

func destroyWithSshuttle(t *testing.T, bastionInstanceID string, bastionPassword string, vpcCidr string, terraformOptions *terraform.Options) error {
	t.Helper()
	cmd, err := runSshuttleInBackground(t, bastionInstanceID, bastionPassword, vpcCidr)
	if err != nil {
		return err
	}
	defer func(t *testing.T, cmd *exec.Cmd) {
		t.Helper()
		err := stopSshuttle(t, cmd)
		require.NoError(t, err)
	}(t, cmd)
	terraform.Destroy(t, terraformOptions)
	return nil
}

func runSshuttleInBackground(t *testing.T, bastionInstanceID string, bastionPassword string, vpcCidr string) (*exec.Cmd, error) {
	t.Helper()
	cmd := exec.Command("sshuttle", "-e", fmt.Sprintf(`sshpass -p "%s" ssh -q -o CheckHostIP=no -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null`, bastionPassword), "--dns", "--disable-ipv6", "-vr", fmt.Sprintf("ec2-user@%s", bastionInstanceID), vpcCidr) //nolint:gosec
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start sshuttle: %w", err)
	}
	return cmd, nil
}

func stopSshuttle(t *testing.T, cmd *exec.Cmd) error {
	t.Helper()
	if err := cmd.Process.Kill(); err != nil {
		return fmt.Errorf("failed to stop sshuttle: %w", err)
	}
	return nil
}