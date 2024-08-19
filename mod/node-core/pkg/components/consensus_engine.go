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
	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft"
)

// ConsensusEngineInput is the input for the consensus engine.
type ConsensusEngineInput[
	AvailabilityStoreT any,
	BlockStoreT any,
	BeaconStateT any,
	StorageBackendT StorageBackend[
		AvailabilityStoreT,
		BeaconStateT,
		BlockStoreT,
		*DepositStore,
	],
] struct {
	depinject.In
	ConsensusMiddleware *ABCIMiddleware
	StorageBackend      StorageBackendT
}

// ProvideConsensusEngine is a depinject provider for the consensus engine.
func ProvideConsensusEngine[
	AvailabilityStoreT any,
	BeaconBlockHeaderT any,
	BeaconStateT BeaconState[
		BeaconStateT, BeaconBlockHeaderT, BeaconStateMarshallableT,
		*Eth1Data, *ExecutionPayloadHeader, *Fork, KVStoreT,
		*Validator, Validators, *Withdrawal,
	],
	BeaconStateMarshallableT any,
	BlockStoreT any,
	KVStoreT any,
	StorageBackendT StorageBackend[
		AvailabilityStoreT,
		BeaconStateT,
		BlockStoreT,
		*DepositStore,
	],
](
	in ConsensusEngineInput[
		AvailabilityStoreT,
		BlockStoreT,
		BeaconStateT,
		StorageBackendT,
	],
) (*cometbft.ConsensusEngine[
	*AttestationData, BeaconStateT, *SlashingInfo,
	*SlotData, StorageBackendT, *ValidatorUpdate,
], error) {
	return cometbft.NewConsensusEngine[
		*AttestationData,
		BeaconStateT,
		*SlashingInfo,
		*SlotData,
		StorageBackendT,
		*ValidatorUpdate,
	](
		in.ConsensusMiddleware,
		in.StorageBackend,
	), nil
}
