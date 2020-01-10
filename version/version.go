// Tideland Go Data Structures and Algorithms - Version
//
// Copyright (C) 2014-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package version

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"strconv"
	"strings"

	"tideland.dev/go/trace/failure"
)

//--------------------
// CONST
//--------------------

// Precedence describes if a version is newer, equal, or older.
type Precedence int

// Level describes the level, on which a version differentiates
// from an other.
type Level string

// Separator, precedences, and part identifiers.
const (
	Metadata = "+"

	Newer Precedence = 1
	Equal Precedence = 0
	Older Precedence = -1

	Major      Level = "major"
	Minor      Level = "minor"
	Patch      Level = "patch"
	PreRelease Level = "pre-release"
	All        Level = "all"
)

//--------------------
// VERSION
//--------------------

// Version implements the version interface.
type Version struct {
	major      int
	minor      int
	patch      int
	preRelease []string
	metadata   []string
}

// New returns a simple version instance. Parts of pre-release
// and metadata are passed as optional strings separated by
// version.Metadata ("+").
func New(major, minor, patch int, prmds ...string) Version {
	if major < 0 {
		major = 0
	}
	if minor < 0 {
		minor = 0
	}
	if patch < 0 {
		patch = 0
	}
	v := Version{
		major: major,
		minor: minor,
		patch: patch,
	}
	isPR := true
	for _, prmd := range prmds {
		if isPR {
			if prmd == Metadata {
				isPR = false
				continue
			}
			v.preRelease = append(v.preRelease, validID(prmd, true))
		} else {
			v.metadata = append(v.metadata, validID(prmd, false))
		}
	}
	return v
}

// Parse retrieves a version out of a string.
func Parse(vsnstr string) (Version, error) {
	// Split version, pre-release, and metadata.
	npmstrs, err := splitVersionString(vsnstr)
	if err != nil {
		return New(0, 0, 0), err
	}
	// Parse these parts.
	nums, err := parseNumberString(npmstrs[0])
	if err != nil {
		return New(0, 0, 0), err
	}
	prmds := []string{}
	if npmstrs[1] != "" {
		prmds = strings.Split(npmstrs[1], ".")
	}
	if npmstrs[2] != "" {
		prmds = append(prmds, Metadata)
		prmds = append(prmds, strings.Split(npmstrs[2], ".")...)
	}
	// Done.
	return New(nums[0], nums[1], nums[2], prmds...), nil
}

// Major returns the major version.
func (v Version) Major() int {
	return v.major
}

// Minor returns the minor version.
func (v Version) Minor() int {
	return v.minor
}

// Patch return the patch version.
func (v Version) Patch() int {
	return v.patch
}

// PreRelease returns a possible pre-release of the version.
func (v Version) PreRelease() string {
	return strings.Join(v.preRelease, ".")
}

// Metadata returns a possible build metadata of the version.
func (v Version) Metadata() string {
	return strings.Join(v.metadata, ".")
}

// Compare compares this version to the passed one. The result
// is from the perspective of this one.
func (v Version) Compare(cv Version) (Precedence, Level) {
	// Standard version parts.
	switch {
	case v.major < cv.Major():
		return Older, Major
	case v.major > cv.Major():
		return Newer, Major
	case v.minor < cv.Minor():
		return Older, Minor
	case v.minor > cv.Minor():
		return Newer, Minor
	case v.patch < cv.Patch():
		return Older, Patch
	case v.patch > cv.Patch():
		return Newer, Patch
	}
	// Now the parts of the pre-release.
	cvpr := []string{}
	for _, cvprPart := range strings.Split(cv.PreRelease(), ".") {
		if cvprPart != "" {
			cvpr = append(cvpr, cvprPart)
		}
	}
	vlen := len(v.preRelease)
	cvlen := len(cvpr)
	count := vlen
	if cvlen < vlen {
		count = cvlen
	}
	for i := 0; i < count; i++ {
		vn, verr := strconv.Atoi(v.preRelease[i])
		cvn, cverr := strconv.Atoi(cvpr[i])
		if verr == nil && cverr == nil {
			// Numerical comparison.
			switch {
			case vn < cvn:
				return Older, PreRelease
			case vn > cvn:
				return Newer, PreRelease
			}
			continue
		}
		// Alphanumerical comparison.
		switch {
		case v.preRelease[i] < cvpr[i]:
			return Older, PreRelease
		case v.preRelease[i] > cvpr[i]:
			return Newer, PreRelease
		}
	}
	// Still no clean result, so the shorter
	// pre-relese is older.
	switch {
	case vlen < cvlen:
		return Newer, PreRelease
	case vlen > cvlen:
		return Older, PreRelease
	}
	// Last but not least: we are equal.
	return Equal, All
}

// Less returns true if this version is less than the passed one.
// This means this version is older.
func (v Version) Less(cv Version) bool {
	precedence, _ := v.Compare(cv)
	return precedence == Older
}

// String implements the fmt.Stringer interface.
func (v Version) String() string {
	vs := fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
	if len(v.preRelease) > 0 {
		vs += "-" + v.PreRelease()
	}
	if len(v.metadata) > 0 {
		vs += Metadata + v.Metadata()
	}
	return vs
}

//--------------------
// TOOLS
//--------------------

// validID reduces the passed identifier to a valid one. If we care
// for numeric identifiers leading zeros will be removed.
func validID(id string, numeric bool) string {
	out := []rune{}
	letter := false
	digit := false
	hyphen := false
	for _, r := range id {
		switch {
		case r >= 'a' && r <= 'z':
			letter = true
			out = append(out, r)
		case r >= 'A' && r <= 'Z':
			letter = true
			out = append(out, r)
		case r >= '0' && r <= '9':
			digit = true
			out = append(out, r)
		case r == '-':
			hyphen = true
			out = append(out, r)
		}
	}
	if numeric && digit && !letter && !hyphen {
		// Digits only, and we care for it.
		// Remove leading zeros.
		for len(out) > 0 && out[0] == '0' {
			out = out[1:]
		}
		if len(out) == 0 {
			out = []rune{'0'}
		}
	}
	return string(out)
}

// splitVersionString separates the version string into numbers,
// pre-release, and metadata strings.
func splitVersionString(vsnstr string) ([]string, error) {
	npXm := strings.SplitN(vsnstr, Metadata, 2)
	switch len(npXm) {
	case 1:
		nXp := strings.SplitN(npXm[0], "-", 2)
		switch len(nXp) {
		case 1:
			return []string{nXp[0], "", ""}, nil
		case 2:
			return []string{nXp[0], nXp[1], ""}, nil
		}
	case 2:
		nXp := strings.SplitN(npXm[0], "-", 2)
		switch len(nXp) {
		case 1:
			return []string{nXp[0], "", npXm[1]}, nil
		case 2:
			return []string{nXp[0], nXp[1], npXm[1]}, nil
		}
	}
	return nil, failure.New("version is malformed: %v", vsnstr)
}

// parseNumberString retrieves major, minor, and patch number
// of the passed string.
func parseNumberString(nstr string) ([]int, error) {
	nstrs := strings.Split(nstr, ".")
	if len(nstrs) < 1 || len(nstrs) > 3 {
		return nil, failure.New("version is malformed: %v", nstr)
	}
	vsn := []int{1, 0, 0}
	for i, nstr := range nstrs {
		num, err := strconv.Atoi(nstr)
		if err != nil {
			return nil, failure.New("version is malformed: %v", err)
		}
		if num < 0 {
			return nil, failure.New("version is malformed: %v", nstr)
		}
		vsn[i] = num
	}
	return vsn, nil
}

// EOF
