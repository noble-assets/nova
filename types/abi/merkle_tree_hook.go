// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, NASD Inc. All rights reserved.
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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abi

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// MerkleTreeHookMetaData contains all meta data concerning the MerkleTreeHook contract.
var MerkleTreeHookMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"root\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// MerkleTreeHookABI is the input ABI used to generate the binding from.
// Deprecated: Use MerkleTreeHookMetaData.ABI instead.
var MerkleTreeHookABI = MerkleTreeHookMetaData.ABI

// MerkleTreeHook is an auto generated Go binding around an Ethereum contract.
type MerkleTreeHook struct {
	MerkleTreeHookCaller     // Read-only binding to the contract
	MerkleTreeHookTransactor // Write-only binding to the contract
	MerkleTreeHookFilterer   // Log filterer for contract events
}

// MerkleTreeHookCaller is an auto generated read-only Go binding around an Ethereum contract.
type MerkleTreeHookCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MerkleTreeHookTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MerkleTreeHookTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MerkleTreeHookFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MerkleTreeHookFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MerkleTreeHookSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MerkleTreeHookSession struct {
	Contract     *MerkleTreeHook   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MerkleTreeHookCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MerkleTreeHookCallerSession struct {
	Contract *MerkleTreeHookCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// MerkleTreeHookTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MerkleTreeHookTransactorSession struct {
	Contract     *MerkleTreeHookTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// MerkleTreeHookRaw is an auto generated low-level Go binding around an Ethereum contract.
type MerkleTreeHookRaw struct {
	Contract *MerkleTreeHook // Generic contract binding to access the raw methods on
}

// MerkleTreeHookCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MerkleTreeHookCallerRaw struct {
	Contract *MerkleTreeHookCaller // Generic read-only contract binding to access the raw methods on
}

// MerkleTreeHookTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MerkleTreeHookTransactorRaw struct {
	Contract *MerkleTreeHookTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMerkleTreeHook creates a new instance of MerkleTreeHook, bound to a specific deployed contract.
func NewMerkleTreeHook(address common.Address, backend bind.ContractBackend) (*MerkleTreeHook, error) {
	contract, err := bindMerkleTreeHook(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MerkleTreeHook{MerkleTreeHookCaller: MerkleTreeHookCaller{contract: contract}, MerkleTreeHookTransactor: MerkleTreeHookTransactor{contract: contract}, MerkleTreeHookFilterer: MerkleTreeHookFilterer{contract: contract}}, nil
}

// NewMerkleTreeHookCaller creates a new read-only instance of MerkleTreeHook, bound to a specific deployed contract.
func NewMerkleTreeHookCaller(address common.Address, caller bind.ContractCaller) (*MerkleTreeHookCaller, error) {
	contract, err := bindMerkleTreeHook(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MerkleTreeHookCaller{contract: contract}, nil
}

// NewMerkleTreeHookTransactor creates a new write-only instance of MerkleTreeHook, bound to a specific deployed contract.
func NewMerkleTreeHookTransactor(address common.Address, transactor bind.ContractTransactor) (*MerkleTreeHookTransactor, error) {
	contract, err := bindMerkleTreeHook(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MerkleTreeHookTransactor{contract: contract}, nil
}

// NewMerkleTreeHookFilterer creates a new log filterer instance of MerkleTreeHook, bound to a specific deployed contract.
func NewMerkleTreeHookFilterer(address common.Address, filterer bind.ContractFilterer) (*MerkleTreeHookFilterer, error) {
	contract, err := bindMerkleTreeHook(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MerkleTreeHookFilterer{contract: contract}, nil
}

// bindMerkleTreeHook binds a generic wrapper to an already deployed contract.
func bindMerkleTreeHook(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MerkleTreeHookMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MerkleTreeHook *MerkleTreeHookRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MerkleTreeHook.Contract.MerkleTreeHookCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MerkleTreeHook *MerkleTreeHookRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MerkleTreeHook.Contract.MerkleTreeHookTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MerkleTreeHook *MerkleTreeHookRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MerkleTreeHook.Contract.MerkleTreeHookTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MerkleTreeHook *MerkleTreeHookCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MerkleTreeHook.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MerkleTreeHook *MerkleTreeHookTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MerkleTreeHook.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MerkleTreeHook *MerkleTreeHookTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MerkleTreeHook.Contract.contract.Transact(opts, method, params...)
}

// Root is a free data retrieval call binding the contract method 0xebf0c717.
//
// Solidity: function root() view returns(bytes32)
func (_MerkleTreeHook *MerkleTreeHookCaller) Root(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _MerkleTreeHook.contract.Call(opts, &out, "root")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Root is a free data retrieval call binding the contract method 0xebf0c717.
//
// Solidity: function root() view returns(bytes32)
func (_MerkleTreeHook *MerkleTreeHookSession) Root() ([32]byte, error) {
	return _MerkleTreeHook.Contract.Root(&_MerkleTreeHook.CallOpts)
}

// Root is a free data retrieval call binding the contract method 0xebf0c717.
//
// Solidity: function root() view returns(bytes32)
func (_MerkleTreeHook *MerkleTreeHookCallerSession) Root() ([32]byte, error) {
	return _MerkleTreeHook.Contract.Root(&_MerkleTreeHook.CallOpts)
}
