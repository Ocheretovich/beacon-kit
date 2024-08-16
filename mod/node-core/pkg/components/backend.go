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
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/storage"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
)

// StorageBackendInput is the input for the ProvideStorageBackend function.
type StorageBackendInput[
	AvailabilityStoreT AvailabilityStore[
		BeaconBlockBodyT, BlobSidecarsT,
	],
	BeaconBlockT any,
	BeaconBlockBodyT any,
	BeaconBlockStoreT BlockStore[BeaconBlockT],
	BlobSidecarsT any,
] struct {
	depinject.In
	AvailabilityStore AvailabilityStoreT
	BlockStore        BeaconBlockStoreT
	ChainSpec         common.ChainSpec
	DepositStore      *DepositStore
	KVStore           *KVStore
}

// ProvideStorageBackend is the depinject provider that returns a beacon storage
// backend.
func ProvideStorageBackend[
	AvailabilityStoreT AvailabilityStore[
		BeaconBlockBodyT, BlobSidecarsT,
	],
	BeaconBlockT any,
	BeaconBlockBodyT constraints.SSZMarshallable,
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconBlockStoreT BlockStore[BeaconBlockT],
	BlobSidecarsT any,
](
	in StorageBackendInput[
		AvailabilityStoreT, BeaconBlockT,
		BeaconBlockBodyT, BeaconBlockStoreT, BlobSidecarsT,
	],
) *StorageBackend {
	return storage.NewBackend[
		AvailabilityStoreT,
		BeaconBlockT,
		BeaconBlockBodyT,
		BeaconBlockHeaderT,
		*BeaconState,
		*BeaconStateMarshallable,
		BlobSidecarsT,
		BeaconBlockStoreT,
		*Deposit,
		*DepositStore,
		*Eth1Data,
		*ExecutionPayloadHeader,
		*Fork,
		*KVStore,
		*Validator,
		Validators,
		*Withdrawal,
		WithdrawalCredentials,
	](
		in.ChainSpec,
		in.AvailabilityStore,
		in.KVStore,
		in.DepositStore,
		in.BlockStore,
	)
}
