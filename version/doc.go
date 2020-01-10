// Tideland Go Data Structures and Algorithms - Version
//
// Copyright (C) 2014-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package version helps other packages to provide information about their
// current version and compare it to others. It follows the idea of semantic
// versioning (see http://semver.org/).
//
// Version instances can be created via New() with explicit passed
// field values or via Parse() and a passed sting. Beside accessing
// the individual fields two versions can be compared with Compare()
// and Less().
package version

// EOF
