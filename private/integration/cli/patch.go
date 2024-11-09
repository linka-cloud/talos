// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

//go:build integration_cli

package cli

import (
	"encoding/json"
	"fmt"

	"github.com/siderolabs/talos/private/integration/base"
	"github.com/siderolabs/talos/pkg/machinery/config/machine"
	"github.com/siderolabs/talos/pkg/machinery/constants"
)

// PatchSuite verifies dmesg command.
type PatchSuite struct {
	base.CLISuite
}

// SuiteName ...
func (suite *PatchSuite) SuiteName() string {
	return "cli.PatchSuite"
}

// TestSuccess successful run.
func (suite *PatchSuite) TestSuccess() {
	node := suite.RandomDiscoveredNodeInternalIP(machine.TypeControlPlane)

	patch := map[string]any{
		"cluster": map[string]any{
			"proxy": map[string]any{
				"image": fmt.Sprintf("%s:v%s", constants.KubeProxyImage, constants.DefaultKubernetesVersion),
			},
		},
	}

	data, err := json.Marshal(patch)
	suite.Require().NoError(err)

	suite.RunCLI([]string{"patch", "--nodes", node, "--patch", string(data), "machineconfig", "--mode=no-reboot"},
		base.StdoutEmpty(),
		base.StderrNotEmpty(),
	)
	suite.RunCLI([]string{"patch", "--nodes", node, "--patch", string(data), "machineconfig", "--mode=no-reboot", "--dry-run"},
		base.StdoutEmpty(),
		base.StderrNotEmpty(),
	)
}

// TestError runs comand with error.
func (suite *PatchSuite) TestError() {
	node := suite.RandomDiscoveredNodeInternalIP(machine.TypeControlPlane)

	patch := []map[string]any{
		{
			"op":   "crash",
			"path": "/cluster/proxy",
			"value": map[string]any{
				"image": fmt.Sprintf("%s:v%s", constants.KubeProxyImage, constants.DefaultKubernetesVersion),
			},
		},
	}

	data, err := json.Marshal(patch)
	suite.Require().NoError(err)

	suite.RunCLI([]string{"patch", "--nodes", node, "--patch", string(data), "machineconfig"},
		base.StdoutEmpty(), base.ShouldFail())
	suite.RunCLI([]string{"patch", "--nodes", node, "--patch", string(data), "machineconfig", "v1alpha2"},
		base.StdoutEmpty(), base.ShouldFail())
	suite.RunCLI([]string{"patch", "--nodes", node, "--patch-file", "/nnnope", "machineconfig"},
		base.StdoutEmpty(), base.ShouldFail())
	suite.RunCLI([]string{"patch", "--nodes", node, "--patch", "it's not even a json", "machineconfig"},
		base.StdoutEmpty(), base.ShouldFail())
}

func init() {
	allSuites = append(allSuites, new(PatchSuite))
}
