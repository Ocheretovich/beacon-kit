// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package cmd

import (
	"os"
	"time"

	"github.com/cometbft/cometbft/libs/log"
	// "github.com/itsdevbear/bolaris/light/runner"

	"github.com/spf13/cobra"

	"github.com/berachain/beacon-kit/light/app"
	"github.com/berachain/beacon-kit/light/client/provider/comet"
)

// beacond light beacond-2061 --hash 587FD10EF595CB799E36F0C21A51861C6D2C81C7452868FA2B8178C7C1689710 --height 6461

const (
	shortDescription = "Run a light node."
	longDescription  = `Run a light node.

	All calls that can be tracked back to a block header by a proof
	will be verified before passing them back to the caller.
	
	Furthermore to the chainID, a fresh instance of a light client will
	need a primary RPC address, a trusted hash and height and witness RPC addresses
	(if not using sequential verification). To restart the node, thereafter
	only the chainID is required.
	`
	example = `light beacond-2061 -primary-addr http://52.57.29.196:26657
	-witness-addr http://public-seed-node.cosmoshub.certus.one:26657
	--height 6461 --hash 587FD10EF595CB799E36F0C21A51861C6D2C81C7452868FA2B8178C7C1689710`
)

func LightClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "light [chainID]",
		Short:   shortDescription,
		Long:    longDescription,
		RunE:    runLightClient,
		Args:    cobra.ExactArgs(1),
		Example: example,
	}

	cmd.Flags().String(listenAddr, defaultListeningAddress, listenAddrDesc)
	cmd.Flags().String(engineURL, defaultEngineURL, engineURLDesc)
	cmd.Flags().String(primaryAddr, defaultPrimaryAddress, primaryAddrDesc)
	cmd.Flags().String(primaryEthAddr, defaultPrimaryEthAddr, primaryEthAddrDesc)
	cmd.Flags().String(witnessAddrsJoined, defaultWitnessAddresses, witnessAddrsJoinedDesc)
	cmd.Flags().String(dir, defaultDir, dirDesc)
	cmd.Flags().Int(maxOpenConnections, defaultMaxOpenConn, maxOpenConnectionsDesc)
	cmd.Flags().Duration(trustingPeriod, defaultTrustPeriod, trustingPeriodDesc)
	cmd.Flags().Int64(trustedHeight, defaultTrustedHeight, trustedHeightDesc)
	cmd.Flags().BytesHex(trustedHash, []byte{}, trustedHashDesc)
	cmd.Flags().String(logLevel, defaultLogLevel, logLevelDesc)
	cmd.Flags().String(trustLevel, defaultTrustLevel, trustLevelDesc)
	cmd.Flags().Bool(sequential, defaultSequential, sequentialDesc)

	return cmd
}

func runLightClient(cmd *cobra.Command, args []string) error {
	// Initialize logger.
	logLvl, err := cmd.Flags().GetString(logLevel)
	if err != nil {
		return err
	}
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	var option log.Option
	if logLvl == "info" {
		option, _ = log.AllowLevel("info")
	} else {
		option, _ = log.AllowLevel("debug")
	}
	logger = log.NewFilter(logger, option)

	config, err := ConfigFromCmd(logger, args[0], cmd)
	if err != nil {
		return err
	}

	// Start the light client.
	err = comet.StartProxy(*config.Comet)
	if err != nil {
		return err
	}

	time.Sleep(1 * time.Second)
	app.RunLightNode(cmd.Context(), config)
	return nil
}
