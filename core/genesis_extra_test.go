// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// This file is a derived work, based on the go-ethereum library whose original
// notices appear below.
//
// It is distributed under a license compatible with the licensing terms of the
// original code from which it is derived.
//
// Much love to the original authors for their work.
// **********
// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"math/big"
	"testing"
	"time"

	"github.com/ava-labs/coreth/params"
	"github.com/ava-labs/coreth/params/extras"
	"github.com/ava-labs/coreth/utils"
	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/libevm/core/rawdb"
	"github.com/ava-labs/libevm/core/types"
	"github.com/ava-labs/libevm/triedb"
	"github.com/stretchr/testify/require"
)

func TestGenesisEthUpgrades(t *testing.T) {
	db := rawdb.NewMemoryDatabase()
	preEthUpgrades := params.WithExtra(
		&params.ChainConfig{
			ChainID:        big.NewInt(43114), // Specifically refers to mainnet for this UT
			HomesteadBlock: big.NewInt(0),
			// For this test to be a proper regression test, DAOForkBlock and
			// DAOForkSupport should be set to match the values in
			// [params.SetEthUpgrades]. Otherwise, in case of a regression, the test
			// would pass as there would be a mismatch at genesis, which is
			// incorrectly considered a success.
			DAOForkBlock:        big.NewInt(0),
			DAOForkSupport:      true,
			EIP150Block:         big.NewInt(0),
			EIP155Block:         big.NewInt(0),
			EIP158Block:         big.NewInt(0),
			ByzantiumBlock:      big.NewInt(0),
			ConstantinopleBlock: big.NewInt(0),
			PetersburgBlock:     big.NewInt(0),
			IstanbulBlock:       big.NewInt(0),
			MuirGlacierBlock:    big.NewInt(0),
		},
		&extras.ChainConfig{
			NetworkUpgrades: extras.NetworkUpgrades{
				ApricotPhase1BlockTimestamp: utils.NewUint64(0),
				ApricotPhase2BlockTimestamp: utils.NewUint64(0),
			},
		},
	)

	tdb := triedb.NewDatabase(db, triedb.HashDefaults)
	config := *preEthUpgrades
	// Set this up once, just to get the genesis hash
	_, genHash, err := SetupGenesisBlock(db, tdb, &Genesis{Config: &config}, common.Hash{}, false)
	require.NoError(t, err)

	// Write the configuration back to the db as it would be in prior versions
	rawdb.WriteChainConfig(db, genHash, preEthUpgrades)

	// Make some other block
	block := types.NewBlock(
		&types.Header{
			Number:     big.NewInt(1640340), // Berlin activation on mainnet
			Difficulty: big.NewInt(1),
			ParentHash: genHash,
			Time:       uint64(time.Now().Unix()),
		},
		nil, nil, nil, nil,
	)
	rawdb.WriteBlock(db, block)

	// We should still be able to re-initialize
	config = *preEthUpgrades
	require.NoError(t, params.SetEthUpgrades(&config)) // New versions will set additional fields eg, LondonBlock
	_, _, err = SetupGenesisBlock(db, tdb, &Genesis{Config: &config}, block.Hash(), false)
	require.NoError(t, err)
}
