// Copyright 2012 Aryan Naraghi (aryan.naraghi@gmail.com)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package difflib

import (
	"reflect"
	"strings"
	"testing"
)

var lcsTests = []struct {
	seq1 string
	seq2 string
	lcs  int
}{
	{"", "", 0},
	{"abc", "abc", 3},
	{"mzjawxu", "xmjyauz", 4},
	{"human", "chimpanzee", 4},
	{"Hello, world!", "Hello, world!", 13},
	{"Hello, world!", "H     e    l  l o ,   w  o r l  d   !", 13},
}

func TestLongestCommonSubsequenceMatrix(t *testing.T) {
	for i, test := range lcsTests {
		seq1 := strings.Split(test.seq1, "")
		seq2 := strings.Split(test.seq2, "")
		matrix := longestCommonSubsequenceMatrix(seq1, seq2)
		lcs := matrix[len(matrix)-1][len(matrix[0])-1] // Grabs the lower, right value.
		if lcs != test.lcs {
			t.Errorf("%d. longestCommonSubsequence(%v, %v)[last][last] => %d, expected %d",
				i, seq1, seq2, lcs, test.lcs)
		}
	}
}

var numEqualStartAndEndElementsTests = []struct {
	seq1  string
	seq2  string
	start int
	end   int
}{
	{"", "", 0, 0},
	{"abc", "", 0, 0},
	{"", "abc", 0, 0},
	{"abc", "abc", 3, 0},
	{"abhelloc", "abbyec", 2, 1},
	{"abchello", "abcbye", 3, 0},
	{"helloabc", "byeabc", 0, 3},
}

func TestNumEqualStartAndEndElements(t *testing.T) {
	for i, test := range numEqualStartAndEndElementsTests {
		seq1 := strings.Split(test.seq1, "")
		seq2 := strings.Split(test.seq2, "")
		start, end := numEqualStartAndEndElements(seq1, seq2)
		if start != test.start || end != test.end {
			t.Errorf("%d. numEqualStartAndEndElements(%v, %v) => (%d, %d), expected (%d, %d)",
				i, seq1, seq2, start, end, test.start, test.end)
		}
	}
}

var diffTests = []struct {
	Seq1         string
	Seq2         string
	Diff         []DiffRecord
	HtmlDiff     string
	PPDiff       string
	AnchoredDiff []DiffRecord
}{
	{
		"",
		"",
		[]DiffRecord{
			{"", Common, 0, 0},
		},
		`<tr><td class="line-num">1</td><td><pre></pre></td><td><pre></pre></td><td class="line-num">1</td></tr>
`,
		`  0,  0   |
`,
		[]DiffRecord{
			{"", Common, 0, 0},
		},
	},

	{
		"same",
		"same",
		[]DiffRecord{
			{"same", Common, 0, 0},
		},
		`<tr><td class="line-num">1</td><td><pre>same</pre></td><td><pre>same</pre></td><td class="line-num">1</td></tr>
`,
		`  0,  0   |same
`,
		[]DiffRecord{
			{"same", Common, 0, 0},
		},
	},

	{
		`one
two
three
`,
		`one
two
three
`,
		[]DiffRecord{
			{"one", Common, 0, 0},
			{"two", Common, 1, 1},
			{"three", Common, 2, 2},
			{"", Common, 3, 3},
		},
		`<tr><td class="line-num">1</td><td><pre>one</pre></td><td><pre>one</pre></td><td class="line-num">1</td></tr>
<tr><td class="line-num">2</td><td><pre>two</pre></td><td><pre>two</pre></td><td class="line-num">2</td></tr>
<tr><td class="line-num">3</td><td><pre>three</pre></td><td><pre>three</pre></td><td class="line-num">3</td></tr>
<tr><td class="line-num">4</td><td><pre></pre></td><td><pre></pre></td><td class="line-num">4</td></tr>
`,
		`  0,  0   |one
  1,  1   |two
  2,  2   |three
  3,  3   |
`,
		[]DiffRecord{
			{"one", Common, 0, 0},
			{"two", Common, 1, 1},
			{"three", Common, 2, 2},
			{"", Common, 3, 3},
		},
	},

	{
		`one
two
three
`,
		`one
five
three
`,
		[]DiffRecord{
			{"one", Common, 0, 0},
			{"two", LeftOnly, 1, 1},
			{"five", RightOnly, 2, 1},
			{"three", Common, 2, 2},
			{"", Common, 3, 3},
		},
		`<tr><td class="line-num">1</td><td><pre>one</pre></td><td><pre>one</pre></td><td class="line-num">1</td></tr>
<tr><td class="line-num">2</td><td class="deleted"><pre>two</pre></td><td></td><td class="line-num"></td></tr>
<tr><td class="line-num"></td><td></td><td class="added"><pre>five</pre></td><td class="line-num">2</td></tr>
<tr><td class="line-num">3</td><td><pre>three</pre></td><td><pre>three</pre></td><td class="line-num">3</td></tr>
<tr><td class="line-num">4</td><td><pre></pre></td><td><pre></pre></td><td class="line-num">4</td></tr>
`,
		`  0,  0   |one
  1,  1 - |five
  2,  1 + |two
  2,  2   |three
  3,  3   |
`,
		[]DiffRecord{
			{"one", Common, 0, 0},
			{"two", LeftOnly, 1, 1},
			{"five", RightOnly, 2, 1},
			{"three", Common, 2, 2},
			{"", Common, 3, 3},
		},
	},

	{
		`Beethoven
Bach
Mozart
Chopin
`,
		`Beethoven
Bach
Brahms
Chopin
Liszt
Wagner
`,

		[]DiffRecord{
			{"Beethoven", Common, 0, 0},
			{"Bach", Common, 1, 1},
			{"Mozart", LeftOnly, 2, 2},
			{"Brahms", RightOnly, 3, 2},
			{"Chopin", Common, 3, 3},
			{"Liszt", RightOnly, 4, 4},
			{"Wagner", RightOnly, 4, 5},
			{"", Common, 4, 6},
		},
		`<tr><td class="line-num">1</td><td><pre>Beethoven</pre></td><td><pre>Beethoven</pre></td><td class="line-num">1</td></tr>
<tr><td class="line-num">2</td><td><pre>Bach</pre></td><td><pre>Bach</pre></td><td class="line-num">2</td></tr>
<tr><td class="line-num">3</td><td class="deleted"><pre>Mozart</pre></td><td></td><td class="line-num"></td></tr>
<tr><td class="line-num"></td><td></td><td class="added"><pre>Brahms</pre></td><td class="line-num">3</td></tr>
<tr><td class="line-num">4</td><td><pre>Chopin</pre></td><td><pre>Chopin</pre></td><td class="line-num">4</td></tr>
<tr><td class="line-num"></td><td></td><td class="added"><pre>Liszt</pre></td><td class="line-num">5</td></tr>
<tr><td class="line-num"></td><td></td><td class="added"><pre>Wagner</pre></td><td class="line-num">6</td></tr>
<tr><td class="line-num">5</td><td><pre></pre></td><td><pre></pre></td><td class="line-num">7</td></tr>
`,
		`  0,  0   |Beethoven
  1,  1   |Bach
  2,  2 - |Brahms
  3,  2 + |Mozart
  3,  3   |Chopin
  4,  4 - |Liszt
  5,  4 - |Wagner
  6,  4   |
`,
		[]DiffRecord{
			{"Beethoven", Common, 0, 0},
			{"Bach", Common, 1, 1},
			{"Mozart", LeftOnly, 2, 2},
			{"Brahms", RightOnly, 3, 2},
			{"Chopin", Common, 3, 3},
			{"Liszt", RightOnly, 4, 4},
			{"Wagner", RightOnly, 4, 5},
			{"", Common, 4, 6},
		},
	},

	{
		`adagio
vivace
staccato legato
presto
lento
`,
		`adagio adagio
staccato
staccato legato
staccato
legato
allegro
`,
		[]DiffRecord{
			{"adagio", LeftOnly, 0, 0},
			{"vivace", LeftOnly, 1, 0},
			{"adagio adagio", RightOnly, 2, 0},
			{"staccato", RightOnly, 2, 1},
			{"staccato legato", Common, 2, 2},
			{"presto", LeftOnly, 3, 3},
			{"lento", LeftOnly, 4, 3},
			{"staccato", RightOnly, 5, 3},
			{"legato", RightOnly, 5, 4},
			{"allegro", RightOnly, 5, 5},
			{"", Common, 5, 6},
		},
		`<tr><td class="line-num">1</td><td class="deleted"><pre>adagio</pre></td><td></td><td class="line-num"></td></tr>
<tr><td class="line-num">2</td><td class="deleted"><pre>vivace</pre></td><td></td><td class="line-num"></td></tr>
<tr><td class="line-num"></td><td></td><td class="added"><pre>adagio adagio</pre></td><td class="line-num">1</td></tr>
<tr><td class="line-num"></td><td></td><td class="added"><pre>staccato</pre></td><td class="line-num">2</td></tr>
<tr><td class="line-num">3</td><td><pre>staccato legato</pre></td><td><pre>staccato legato</pre></td><td class="line-num">3</td></tr>
<tr><td class="line-num">4</td><td class="deleted"><pre>presto</pre></td><td></td><td class="line-num"></td></tr>
<tr><td class="line-num">5</td><td class="deleted"><pre>lento</pre></td><td></td><td class="line-num"></td></tr>
<tr><td class="line-num"></td><td></td><td class="added"><pre>staccato</pre></td><td class="line-num">4</td></tr>
<tr><td class="line-num"></td><td></td><td class="added"><pre>legato</pre></td><td class="line-num">5</td></tr>
<tr><td class="line-num"></td><td></td><td class="added"><pre>allegro</pre></td><td class="line-num">6</td></tr>
<tr><td class="line-num">6</td><td><pre></pre></td><td><pre></pre></td><td class="line-num">7</td></tr>
`,
		`  0,  0 - |adagio adagio
  1,  0 - |staccato
  2,  0 + |adagio
  2,  1 + |vivace
  2,  2   |staccato legato
  3,  3 - |staccato
  4,  3 - |legato
  5,  3 - |allegro
  6,  3 + |presto
  6,  4 + |lento
  6,  5   |
`,
		[]DiffRecord{
			{"adagio", LeftOnly, 0, 0},
			{"vivace", LeftOnly, 1, 0},
			{"adagio adagio", RightOnly, 2, 0},
			{"staccato", RightOnly, 2, 1},
			{"staccato legato", Common, 2, 2},
			{"presto", LeftOnly, 3, 3},
			{"lento", LeftOnly, 4, 3},
			{"staccato", RightOnly, 5, 3},
			{"legato", RightOnly, 5, 4},
			{"allegro", RightOnly, 5, 5},
			{"", Common, 5, 6},
		},
	},

	{
		`alpha
beta
gama
delta
beta
pi
`,
		`solid
liquid
gas
beta
plasma
`,
		[]DiffRecord{
			{"alpha", LeftOnly, 0, 0},
			{"beta", LeftOnly, 1, 0},
			{"gama", LeftOnly, 2, 0},
			{"delta", LeftOnly, 3, 0},
			{"solid", RightOnly, 4, 0},
			{"liquid", RightOnly, 4, 1},
			{"gas", RightOnly, 4, 2},
			{"beta", Common, 4, 3},
			{"pi", LeftOnly, 5, 4},
			{"plasma", RightOnly, 6, 4},
			{"", Common, 6, 5},
		},
		`<tr><td class="line-num">1</td><td class="deleted"><pre>alpha</pre></td><td></td><td class="line-num"></td></tr>
<tr><td class="line-num">2</td><td class="deleted"><pre>beta</pre></td><td></td><td class="line-num"></td></tr>
<tr><td class="line-num">3</td><td class="deleted"><pre>gama</pre></td><td></td><td class="line-num"></td></tr>
<tr><td class="line-num">4</td><td class="deleted"><pre>delta</pre></td><td></td><td class="line-num"></td></tr>
<tr><td class="line-num"></td><td></td><td class="added"><pre>solid</pre></td><td class="line-num">1</td></tr>
<tr><td class="line-num"></td><td></td><td class="added"><pre>liquid</pre></td><td class="line-num">2</td></tr>
<tr><td class="line-num"></td><td></td><td class="added"><pre>gas</pre></td><td class="line-num">3</td></tr>
<tr><td class="line-num">5</td><td><pre>beta</pre></td><td><pre>beta</pre></td><td class="line-num">4</td></tr>
<tr><td class="line-num">6</td><td class="deleted"><pre>pi</pre></td><td></td><td class="line-num"></td></tr>
<tr><td class="line-num"></td><td></td><td class="added"><pre>plasma</pre></td><td class="line-num">5</td></tr>
<tr><td class="line-num">7</td><td><pre></pre></td><td><pre></pre></td><td class="line-num">6</td></tr>
`,
		`  0,  0 - |solid
  1,  0 - |liquid
  2,  0 - |gas
  3,  0 + |alpha
  3,  1 + |beta
  3,  2 + |gama
  3,  3 + |delta
  3,  4   |beta
  4,  5 - |plasma
  5,  5 + |pi
  5,  6   |
`,
		[]DiffRecord{
			{"alpha", LeftOnly, 0, 0},
			{"beta", LeftOnly, 1, 0},
			{"gama", LeftOnly, 2, 0},
			{"delta", LeftOnly, 3, 0},
			{"beta", LeftOnly, 4, 0},
			{"pi", LeftOnly, 5, 0},
			{"solid", RightOnly, 6, 0},
			{"liquid", RightOnly, 6, 1},
			{"gas", RightOnly, 6, 2},
			{"beta", RightOnly, 6, 3},
			{"plasma", RightOnly, 6, 4},
			{"", Common, 6, 5},
		},
	},
}

func TestDiff(t *testing.T) {
	for i, test := range diffTests {
		seq1 := strings.Split(test.Seq1, "\n")
		seq2 := strings.Split(test.Seq2, "\n")

		diff := Diff(seq1, seq2)
		if !reflect.DeepEqual(diff, test.Diff) {
			t.Errorf("%d. Diff(%v, %v) => %v, expected %v",
				i, seq1, seq2, diff, test.Diff)
		}

		htmlDiff := HTMLDiff(seq1, seq2)
		if htmlDiff != test.HtmlDiff {
			t.Errorf("%d. HtmlDiff(%v, %v) => %v, expected %v",
				i, seq1, seq2, htmlDiff, test.HtmlDiff)
		}

		ppDiff := PPDiff(seq1, seq2)
		if ppDiff != test.PPDiff {
			t.Errorf("%d. PPDiff(%v, %v) => '%v', expected '%v'",
				i, seq1, seq2, ppDiff, test.PPDiff)
		}

		anchoredDiff := AnchoredDiff(seq1, seq2)
		if !reflect.DeepEqual(anchoredDiff, test.AnchoredDiff) {
			t.Errorf("%d. AnchoredDiff(%v, %v) => %v, expected %v",
				i, seq1, seq2, anchoredDiff, test.Diff)
		}
	}
}
