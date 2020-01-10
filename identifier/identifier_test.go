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

// TestTypeAsIdentifierPart tests the creation of identifiers based on types.
func TestTypeAsIdentifierPart(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	// Type as identifier.
	var tai TypeToSplitForIdentifier

	id := identifier.TypeAsIdentifierPart(tai)
	assert.Equal(id, "type-to-split-for-identifier")

	id = identifier.TypeAsIdentifierPart(identifier.NewUUID())
	assert.Equal(id, "u-u-i-d")
}

// TestIdentifier tests the creation of identifiers based on parts.
func TestIdentifier(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	// Identifier.
	id := identifier.Identifier("One", 2, "three four")
	assert.Equal(id, "one:2:three-four")

	id = identifier.Identifier(2011, 6, 22, "One, two, or  three things.")
	assert.Equal(id, "2011:6:22:one-two-or-three-things")
}

// TestSepIdentifier tests the creation of identifiers based on parts
// with defined seperators.
func TestSepIdentifier(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	id := identifier.SepIdentifier("+", 1, "oNe", 2, "TWO", "3", "ÄÖÜ")
	assert.Equal(id, "1+one+2+two+3")

	id = identifier.LimitedSepIdentifier("+", false, 1, "oNe", 2, "TWO", "3", "ÄÖÜ")
	assert.Equal(id, "1+one+2+two+3+äöü")

	id = identifier.LimitedSepIdentifier("+", true, "     ", 1, "oNe", 2, "TWO", "3", "ÄÖÜ", "Four", "+#-:,")
	assert.Equal(id, "1+one+2+two+3+four")
}

//--------------------
// HELPER
//--------------------

// Type as part of an identifier.
type TypeToSplitForIdentifier bool

// EOF
