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
package types

import (
	"encoding/hex"
	"strings"

	"github.com/svaishnavy/crypto/ed25519"
)

type BlockHash string
type Account string
type Work string
type Signature string

func (hash BlockHash) ToBytes() []byte {
	bytes, err := hex.DecodeString(string(hash))
	if err != nil {
		panic(err)
	}
	return bytes
}

func (sig Signature) ToBytes() []byte {
	bytes, err := hex.DecodeString(string(sig))
	if err != nil {
		panic(err)
	}
	return bytes
}

func (hash BlockHash) Sign(private_key ed25519.PrivateKey) Signature {
	sig := hex.EncodeToString(ed25519.Sign(private_key, hash.ToBytes()))
	return Signature(strings.ToUpper(sig))
}

func BlockHashFromBytes(b []byte) BlockHash {
	return BlockHash(strings.ToUpper(hex.EncodeToString(b)))
}
