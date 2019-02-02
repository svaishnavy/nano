/*
Copyright (c) 2018 Frank Hamand
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package wallet

import (
	"encoding/hex"
	"testing"

	"github.com/svaishnavy/nano/address"
	"github.com/svaishnavy/nano/blocks"
	"github.com/svaishnavy/nano/store"
	"github.com/svaishnavy/nano/uint128"
)

func TestNew(t *testing.T) {
	store.Init(store.TestConfig)

	w := New(blocks.TestPrivateKey)
	if w.GetBalance() != blocks.GenesisAmount {
		t.Errorf("Genesis block doesn't have correct balance")
	}
}

func TestPoW(t *testing.T) {
	blocks.WorkThreshold = 0xff00000000000000
	store.Init(store.TestConfig)
	w := New(blocks.TestPrivateKey)

	if w.GeneratePoWAsync() != nil || !w.WaitingForPoW() {
		t.Errorf("Failed to start PoW generation")
	}

	if w.GeneratePoWAsync() == nil {
		t.Errorf("Started PoW while already in progress")
	}

	_, err := w.Send(blocks.TestGenesisBlock.Account, uint128.FromInts(0, 1))

	if err == nil {
		t.Errorf("Created send block without PoW")
	}

	w.WaitPoW()

	send, _ := w.Send(blocks.TestGenesisBlock.Account, uint128.FromInts(0, 1))

	if !blocks.ValidateBlockWork(send) {
		t.Errorf("Invalid work")
	}

}

func TestSend(t *testing.T) {
	blocks.WorkThreshold = 0xff00000000000000
	store.Init(store.TestConfig)
	w := New(blocks.TestPrivateKey)

	w.GeneratePowSync()
	amount := uint128.FromInts(1, 1)

	send, _ := w.Send(blocks.TestGenesisBlock.Account, amount)

	if w.GetBalance() != blocks.GenesisAmount.Sub(amount) {
		t.Errorf("Balance unchanged after send")
	}

	_, err := w.Send(blocks.TestGenesisBlock.Account, blocks.GenesisAmount)
	if err == nil {
		t.Errorf("Sent more than account balance")
	}

	w.GeneratePowSync()
	store.StoreBlock(send)
	receive, _ := w.Receive(send.Hash())
	store.StoreBlock(receive)

	if w.GetBalance() != blocks.GenesisAmount {
		t.Errorf("Balance not updated after receive, %x != %x", w.GetBalance().GetBytes(), blocks.GenesisAmount.GetBytes())
	}

}

func TestOpen(t *testing.T) {
	blocks.WorkThreshold = 0xff00000000000000
	store.Init(store.TestConfig)
	amount := uint128.FromInts(1, 1)

	sendW := New(blocks.TestPrivateKey)
	sendW.GeneratePowSync()

	_, priv := address.GenerateKey()
	openW := New(hex.EncodeToString(priv))
	send, _ := sendW.Send(openW.Address(), amount)
	openW.GeneratePowSync()

	_, err := openW.Open(send.Hash(), openW.Address())
	if err == nil {
		t.Errorf("Expected error for referencing unstored send")
	}

	if openW.GetBalance() != uint128.FromInts(0, 0) {
		t.Errorf("Open should start at zero balance")
	}

	store.StoreBlock(send)
	_, err = openW.Open(send.Hash(), openW.Address())
	if err != nil {
		t.Errorf("Open block failed: %s", err)
	}

	if openW.GetBalance() != amount {
		t.Errorf("Open balance didn't equal send amount")
	}

	_, err = openW.Open(send.Hash(), openW.Address())
	if err == nil {
		t.Errorf("Expected error for creating duplicate open block")
	}

}
