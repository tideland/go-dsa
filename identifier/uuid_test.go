// Tideland Go Data Structures and Algorithms - Identifier - Unit Tests
//
// Copyright (C) 2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package identifier_test

//--------------------
// IMPORTS
//--------------------

import (
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/dsa/identifier"
)

//--------------------
// TESTS
//--------------------

// TestStandardUUID tests the standard UUID.
func TestStandardUUID(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	// Asserts.
	uuid := identifier.NewUUID()
	assert.Equal(uuid.Version(), identifier.UUIDv4)
	uuidShortStr := uuid.ShortString()
	uuidStr := uuid.String()
	assert.Equal(len(uuid), 16)
	assert.Match(uuidShortStr, "[0-9a-f]{32}")
	assert.Match(uuidStr, "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")
	// Check for copy.
	uuidA := identifier.NewUUID()
	uuidB := uuidA.Copy()
	for i := 0; i < len(uuidA); i++ {
		uuidA[i] = 0
	}
	assert.Different(uuidA, uuidB)
}

// TestUUIDVersions tests the creation of different UUID versions.
func TestUUIDVersions(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	ns := identifier.UUIDNamespaceOID()
	// Asserts.
	uuidV1, err := identifier.NewUUIDv1()
	assert.Nil(err)
	assert.Equal(uuidV1.Version(), identifier.UUIDv1)
	assert.Equal(uuidV1.Variant(), identifier.UUIDVariantRFC4122)
	assert.Logf("UUID V1: %v", uuidV1)
	uuidV3, err := identifier.NewUUIDv3(ns, []byte{4, 7, 1, 1})
	assert.Nil(err)
	assert.Equal(uuidV3.Version(), identifier.UUIDv3)
	assert.Equal(uuidV3.Variant(), identifier.UUIDVariantRFC4122)
	assert.Logf("UUID V3: %v", uuidV3)
	uuidV4, err := identifier.NewUUIDv4()
	assert.Nil(err)
	assert.Equal(uuidV4.Version(), identifier.UUIDv4)
	assert.Equal(uuidV4.Variant(), identifier.UUIDVariantRFC4122)
	assert.Logf("UUID V4: %v", uuidV4)
	uuidV5, err := identifier.NewUUIDv5(ns, []byte{4, 7, 1, 1})
	assert.Nil(err)
	assert.Equal(uuidV5.Version(), identifier.UUIDv5)
	assert.Equal(uuidV5.Variant(), identifier.UUIDVariantRFC4122)
	assert.Logf("UUID V5: %v", uuidV5)
}

// TestUUIDByHex tests creating UUIDs from hex strings.
func TestUUIDByHex(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	// Asserts.
	_, err := identifier.NewUUIDByHex("ffff")
	assert.ErrorMatch(err, `.* source length is not 32`)
	_, err = identifier.NewUUIDByHex("012345678901234567890123456789zz")
	assert.ErrorMatch(err, `.* source is no hex value: .*`)
	_, err = identifier.NewUUIDByHex("012345678901234567890123456789ab")
	assert.Nil(err)
}

// EOF
