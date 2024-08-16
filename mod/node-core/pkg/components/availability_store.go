// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package components

import (
	"os"

	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/async/pkg/broker"
	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/storage/pkg/filedb"
	"github.com/berachain/beacon-kit/mod/storage/pkg/manager"
	"github.com/berachain/beacon-kit/mod/storage/pkg/pruner"
	"github.com/cosmos/cosmos-sdk/client/flags"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cast"
)

// AvailabilityStoreInput is the input for the ProviderAvailabilityStore
// function for the depinject framework.
type AvailabilityStoreInput[
	AttestationDataT any,
	BeaconBlockBodyT BeaconBlockBody[
		BeaconBlockBodyT, AttestationDataT, DepositT,
		Eth1DataT, ExecutionPayloadT, SlashingInfoT,
	],
	DepositT any,
	Eth1DataT any,
	ExecutionPayloadT any,
	LoggerT log.AdvancedLogger[any, LoggerT],
	SlashingInfoT any,
] struct {
	depinject.In
	AppOpts   servertypes.AppOptions
	ChainSpec common.ChainSpec
	Logger    LoggerT
}

// ProvideAvailibilityStore provides the availability store.
func ProvideAvailibilityStore[
	AttestationDataT any,
	BeaconBlockBodyT BeaconBlockBody[
		BeaconBlockBodyT, AttestationDataT, DepositT,
		Eth1DataT, ExecutionPayloadT, SlashingInfoT,
	],
	DepositT any,
	Eth1DataT any,
	ExecutionPayloadT any,
	LoggerT log.AdvancedLogger[any, LoggerT],
	SlashingInfoT any,
](
	in AvailabilityStoreInput[
		AttestationDataT, BeaconBlockBodyT, DepositT, Eth1DataT,
		ExecutionPayloadT, LoggerT, SlashingInfoT,
	],
) (*dastore.Store[BeaconBlockBodyT], error) {
	return dastore.New[BeaconBlockBodyT](
		filedb.NewRangeDB(
			filedb.NewDB(
				filedb.WithRootDirectory(
					cast.ToString(
						in.AppOpts.Get(flags.FlagHome),
					)+"/data/blobs",
				),
				filedb.WithFileExtension("ssz"),
				filedb.WithDirectoryPermissions(os.ModePerm),
				filedb.WithLogger(in.Logger),
			),
		),
		in.Logger.With("service", "da-store"),
		in.ChainSpec,
	), nil
}

// AvailabilityPrunerInput is the input for the ProviderAvailabilityPruner
// function for the depinject framework.
type AvailabilityPrunerInput[
	AvailabilityStoreT AvailabilityStore[
		BeaconBlockBodyT, BlobSidecarsT,
	],
	BeaconBlockT any,
	BeaconBlockBodyT any,
	BlobSidecarsT any,
	LoggerT log.AdvancedLogger[any, LoggerT],
] struct {
	depinject.In
	AvailabilityStore AvailabilityStoreT
	BlockBroker       *broker.Broker[*asynctypes.Event[BeaconBlockT]]
	ChainSpec         common.ChainSpec
	Logger            LoggerT
}

// ProvideAvailabilityPruner provides a availability pruner for the depinject
// framework.
func ProvideAvailabilityPruner[
	AvailabilityStoreT AvailabilityStore[
		BeaconBlockBodyT, BlobSidecarsT,
	],
	BeaconBlockT BeaconBlock[
		BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
	],
	BeaconBlockBodyT any,
	BeaconBlockHeaderT any,
	BlobSidecarsT any,
	LoggerT log.AdvancedLogger[any, LoggerT],
](
	in AvailabilityPrunerInput[
		AvailabilityStoreT, BeaconBlockT,
		BeaconBlockBodyT, BlobSidecarsT, LoggerT,
	],
) (DAPruner, error) {
	subCh, err := in.BlockBroker.Subscribe()
	if err != nil {
		in.Logger.Error("failed to subscribe to block feed", "err", err)
		return nil, err
	}

	// build the availability pruner if IndexDB is available.
	return pruner.NewPruner[
		BeaconBlockT,
		*asynctypes.Event[BeaconBlockT],
		AvailabilityStoreT,
	](
		in.Logger.With("service", manager.AvailabilityPrunerName),
		in.AvailabilityStore,
		manager.AvailabilityPrunerName,
		subCh,
		dastore.BuildPruneRangeFn[
			BeaconBlockT,
			*asynctypes.Event[BeaconBlockT],
		](in.ChainSpec),
	), nil
}
