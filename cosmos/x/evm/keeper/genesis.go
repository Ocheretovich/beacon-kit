// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2023, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
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

package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/itsdevbear/bolaris/cosmos/x/evm/store"
	"github.com/itsdevbear/bolaris/cosmos/x/evm/types"
)

func (k *Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {
	genesisStore := store.NewGenesis(ctx.KVStore(k.storeKey))
	if err := genesisStore.Store(data.Eth1GenesisHash); err != nil {
		panic(err)
	}
}

func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	genesisStore := store.NewGenesis(ctx.KVStore(k.storeKey))
	return &types.GenesisState{
		Eth1GenesisHash: genesisStore.Retrieve().Hex(),
	}
}
