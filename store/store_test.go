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
package store

import (
	"os"
	"testing"

	"github.com/svaishnavy/nano/blocks"
)

func TestInit(t *testing.T) {
	Init(TestConfigLive)

	os.RemoveAll(TestConfigLive.Path)
}

func TestGenesisBalance(t *testing.T) {
	Init(TestConfigLive)

	block := FetchBlock(blocks.LiveGenesisBlockHash)

	if GetBalance(block).String() != "ffffffffffffffffffffffffffffffff" {
		t.Errorf("Genesis block has invalid initial balance")
	}
	os.RemoveAll(TestConfigLive.Path)
}

func TestMissingBlock(t *testing.T) {
	Init(TestConfig)

	block := FetchBlock(blocks.LiveGenesisBlockHash)

	if block != nil {
		t.Errorf("Found live genesis on test config")
	}
	os.RemoveAll(TestConfig.Path)
}
