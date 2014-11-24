// antha/token/serialize_test.go: Part of the Antha language
// Copyright (C) 2014 The Antha authors. All rights reserved.
// 
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
// 
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// 
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
// 
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o 
// Synthace Ltd. The London Bioscience Innovation Centre
// 1 Royal College St, London NW1 0NH UK


package token

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"testing"
)

// equal returns nil if p and q describe the same file set;
// otherwise it returns an error describing the discrepancy.
func equal(p, q *FileSet) error {
	if p == q {
		// avoid deadlock if p == q
		return nil
	}

	// not strictly needed for the test
	p.mutex.Lock()
	q.mutex.Lock()
	defer q.mutex.Unlock()
	defer p.mutex.Unlock()

	if p.base != q.base {
		return fmt.Errorf("different bases: %d != %d", p.base, q.base)
	}

	if len(p.files) != len(q.files) {
		return fmt.Errorf("different number of files: %d != %d", len(p.files), len(q.files))
	}

	for i, f := range p.files {
		g := q.files[i]
		if f.set != p {
			return fmt.Errorf("wrong fileset for %q", f.name)
		}
		if g.set != q {
			return fmt.Errorf("wrong fileset for %q", g.name)
		}
		if f.name != g.name {
			return fmt.Errorf("different filenames: %q != %q", f.name, g.name)
		}
		if f.base != g.base {
			return fmt.Errorf("different base for %q: %d != %d", f.name, f.base, g.base)
		}
		if f.size != g.size {
			return fmt.Errorf("different size for %q: %d != %d", f.name, f.size, g.size)
		}
		for j, l := range f.lines {
			m := g.lines[j]
			if l != m {
				return fmt.Errorf("different offsets for %q", f.name)
			}
		}
		for j, l := range f.infos {
			m := g.infos[j]
			if l.Offset != m.Offset || l.Filename != m.Filename || l.Line != m.Line {
				return fmt.Errorf("different infos for %q", f.name)
			}
		}
	}

	// we don't care about .last - it's just a cache
	return nil
}

func checkSerialize(t *testing.T, p *FileSet) {
	var buf bytes.Buffer
	encode := func(x interface{}) error {
		return gob.NewEncoder(&buf).Encode(x)
	}
	if err := p.Write(encode); err != nil {
		t.Errorf("writing fileset failed: %s", err)
		return
	}
	q := NewFileSet()
	decode := func(x interface{}) error {
		return gob.NewDecoder(&buf).Decode(x)
	}
	if err := q.Read(decode); err != nil {
		t.Errorf("reading fileset failed: %s", err)
		return
	}
	if err := equal(p, q); err != nil {
		t.Errorf("filesets not identical: %s", err)
	}
}

func TestSerialization(t *testing.T) {
	p := NewFileSet()
	checkSerialize(t, p)
	// add some files
	for i := 0; i < 10; i++ {
		f := p.AddFile(fmt.Sprintf("file%d", i), p.Base()+i, i*100)
		checkSerialize(t, p)
		// add some lines and alternative file infos
		line := 1000
		for offs := 0; offs < f.Size(); offs += 40 + i {
			f.AddLine(offs)
			if offs%7 == 0 {
				f.AddLineInfo(offs, fmt.Sprintf("file%d", offs), line)
				line += 33
			}
		}
		checkSerialize(t, p)
	}
}