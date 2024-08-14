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
	"github.com/berachain/beacon-kit/mod/config"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/payload/pkg/attributes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
)

type AttributesFactoryInput[
	LoggerT log.Logger[any],
] struct {
	depinject.In

	ChainSpec common.ChainSpec
	Config    *config.Config
	Logger    LoggerT
}

// ProvideAttributesFactory provides an AttributesFactory for the client.
func ProvideAttributesFactory[
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT BeaconState[
		BeaconStateT, BeaconBlockHeaderT, BeaconStateMarshallableT,
		Eth1DataT, ExecutionPayloadHeaderT, ForkT, KVStoreT,
		ValidatorT, ValidatorsT, WithdrawalT,
	],
	BeaconStateMarshallableT any,
	Eth1DataT any,
	ExecutionPayloadHeaderT any,
	ForkT any,
	KVStoreT any,
	LoggerT log.Logger[any],
	PayloadAttributesT PayloadAttributes[PayloadAttributesT, WithdrawalT],
	ValidatorT any,
	ValidatorsT any,
	WithdrawalT any,
](
	in AttributesFactoryInput[LoggerT],
) (*attributes.Factory[
	BeaconStateT,
	PayloadAttributesT,
	WithdrawalT,
], error) {
	return attributes.NewAttributesFactory[
		BeaconStateT,
		PayloadAttributesT,
		WithdrawalT,
	](
		in.ChainSpec,
		in.Logger,
		in.Config.PayloadBuilder.SuggestedFeeRecipient,
	), nil
}
