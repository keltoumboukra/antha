package enzymes

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

var data0 = `
[
        {
          "jm": "RHAMNOSE",
          "seq": "AAAAGGATCTCAAGAAGATCCTTTGATCTTTTCTACGGGGTCTGACGCTCAGTGGAACGACGCGCGCGTAACTCACGTTAAGGGATTTTGGTCATGAGCTTGCGCCGTCCCGTCAAGTCAGCGTAATGCTCTGCTTTTAGAAAAACTCATCGAGCATCAAATGAAACTGCAATTTATTCATATCAGGATTATCAATACCATATTTTTGAAAAAGCCGTTTCTGTAATGAAGGAGAAAACTCACCGAGGCAGTTCCATAGGATGGCAAGATCCTGGTATCGGTCTGCGATTCCGACTCGTCCAACATCAATACAACCTATTAATTTCCCCTCGTCAAAAATAAGGTTATCAAGTGAGAAATCACCATGAGTGACGACTGAATCCGGTGAGAATGGCAAAAGTTTATGCATTTCTTTCCAGACTTGTTCAACAGGCCAGCCATTACGCTCGTCATCAAAATCACTCGCATCAACCAAACCGTTATTCATTCGTGATTGCGCCTGAGCGAGGCGAAATACGCGATCGCTGTTAAAAGGACAATTACAAACAGGAATCGAGTGCAACCGGCGCAGGAACACTGCCAGCGCATCAACAATATTTTCACCTGAATCAGGATATTCTTCTAATACCTGGAACGCTGTTTTTCCGGGGATCGCAGTGGTGAGTAACCATGCATCATCAGGAGTACGGATAAAATGCTTGATGGTCGGAAGTGGCATAAATTCCGTCAGCCAGTTTAGTCTGACCATCTCATCTGTAACATCATTGGCAACGCTACCTTTGCCATGTTTCAGAAACAACTCTGGCGCATCGGGCTTCCCATACAAGCGATAGATTGTCGCACCTGATTGCCCGACATTATCGCGAGCCCATTTATACCCATATAAATCAGCATCCATGTTGGAATTTAATCGCGGCCTCGACGTTTCCCGTTGAATATGGCTCATATTCTTCCTTTTTCAATATTATTGAAGCATTTATCAGGGTTATTGTCTCATGAGCGGATACATATTTGAATGTATTTAGAAAAATAAACAAATAGGGGTCAGTGTTACAACCAATTAACCAATTCTGAACATTATCGCGAGCCCATTTATACCTGAATATGGCTCATAACACCCCTTGTTTGCCTGGCGGCAGTAGCGCGGTGGTCCCACCTGACCCCATGCCGAACTCAGAAGTGAAACGCCGTAGCGCCGATGGTAGTGTGGGGACTCCCCATGCGAGAGTAGGGAACTGCCAGGCATCAAATAAAACGAAAGGCTCAGTCGAAAGACTGGGCCTTTCGCCCGGGCTAATTAGGGGGTGTCGCCCTTTACACGTACTTAGTCGCTGAAGCTCTTCAGCGGGTCTCAGGCACACCACAATTCAGCAAATTGTGAACATCATCACGTTCATCTTTCCCTGGTTGCCAATGGCCCATTTTCCTGTCAGTAACGAGAAGGTCGCGAATTCAGGCGCTTTTTAGACTGGTCGTAATGAAATTCTTTTTAAGGAGGTAAAAAATGCATCACCACCATCACCACATGAGAAGAGCCGTCAATCGAGTTCGTACCTAAGGGCGACACAAAATTTATTCTAAATGCATAATAAATACTGATAACATCTTATAGTTTGTATTATATTTTGTATTATCGTTGACATGTATAATTTTGATATCAAAAACTGATTTTCCCTTTATTATTTTCGAGATTTATTTTCTTAATTCTCTTTAACAAACTAGAAATATTGTATATACAAAAAATCATAAATAATAGATGAATAGTTTAATTATAGGTGTTCATCAATCGAAAAAGCAACGTATCTTATTTAAAGTGCGTTGCTTTTTTCTCATTTATAAGGTTAAATAATTCTCATATATCAAGCAAAGTGACAGGCGCCCTTAAATATTCTGACAAATGCTCTTTCCCTAAACTCCCCCCATAAAAAAACCCGCCGAAGCGGGTTTTTACGTTATTTGCGGATTAACGATTACTCGTTATCAGAACCGCCCAGGGGGCCCGAGCTTAAGACTGGCCGTCGTTTTACAACACAGAAAGAGTTTGTAGAAACGCAAAAAGGCCATCCGTCAGGGGCCTTCTGCTTAGTTTGATGCCTGGCAGTTCCCTACTCTCGCCTTCCGCTTCCTCGCTCACTGACTCGCTGCGCTCGGTCGTTCGGCTGCGGCGAGCGGTATCAGCTCACTCAAAGGCGGTAATACGGTTATCCACAGAATCAGGGGATAACGCAGGAAAGAACATGTGAGCAAAAGGCCAGCAAAAGGCCAGGAACCGTAAAAAGGCCGCGTTGCTGGCGTTTTTCCATAGGCTCCGCCCCCCTGACGAGCATCACAAAAATCGACGCTCAAGTCAGAGGTGGCGAAACCCGACAGGACTATAAAGATACCAGGCGTTTCCCCCTGGAAGCTCCCTCGTGCGCTCTCCTGTTCCGACCCTGCCGCTTACCGGATACCTGTCCGCCTTTCTCCCTTCGGGAAGCGTGGCGCTTTCTCATAGCTCACGCTGTAGGTATCTCAGTTCGGTGTAGGTCGTTCGCTCCAAGCTGGGCTGTGTGCACGAACCCCCCGTTCAGCCCGACCGCTGCGCCTTATCCGGTAACTATCGTCTTGAGTCCAACCCGGTAAGACACGACTTATCGCCACTGGCAGCAGCCACTGGTAACAGGATTAGCAGAGCGAGGTATGTAGGCGGTGCTACAGAGTTCTTGAAGTGGTGGGCTAACTACGGCTACACTAGAAGAACAGTATTTGGTATCTGCGCTCTGCTGAAGCCAGTTACCTTCGGAAAAAGAGTTGGTAGCTCTTGATCCGGCAAACAAACCACCGCTGGTAGCGGTGGTTTTTTTGTTTGCAAGCAGCAGATTACGCGCAGAAA",
          "plasmid": true
        },
        {
          "jm": "COMET_P1",
          "seq": "GCTCTTCTATGACGGCATTGACGGAAGGCGCAAAATTGTTCGAAAAAGAAATCCCATACATCACCGAACTGGAAGGCGATGTTGAAGGTATGAAGTTCATCATTAAGGGTGAGGGCACCGGCGATGCAACTACGGGCACCATTAAAGCGAAGTATATCTGCACCACCGGTGACGTTCCGGTGCCGTGGAGCACGCTGGTCACCACCCTGACCTATGGCGCGCAGTGTTTCGCGAAGTACGGTCCGGAACTGGAAGAGC"
        },
        {
          "jm": "COMET_P2",
          "seq": "GCTCTTCGACTGAAGGACTTCTATAAGAGCTGTATGCCTGAGGGCTATGTTCAGGAGCGTACCATTACCTTTGAGGGTGATGGTGTCTTTAAGACGCGTGCTGAGGTGACCTTTGAGAATGGTTCCGTGTACAATCGCGTGAAACTGAATGGTCAAGGTTTTAAGAAAGATGGTCACGTGCTGGGCAAAAACCTGGAGTTTAACTTTACTCCGCATTGCCTGTGCATTTGGGGCGACCAAGCGAACCACGGTCTTGAAGAGC"
        },
        {
          "jm": "COMET_P3",
          "seq": "GCTCTTCGTCTGAAAAGCGCGTTCAAGATTATGCACGAGATTACGGGTAGCAAAGAGGACTTCATCGTGGCCGACCACACGCAGATGAACACCCCGATCGGTGGCGGTCCGGTCCATGTCCCGGAGTACCACCACTTGACCGTTTGGACCTCTTTCGGTAAAGACCCGGATGATGACGAAACGGATCATCTGAATATTGTTGAGGTTATCAAAGCCGTCGACCTGGAAACTTACCGTTAATGATAATGAGGTAGAAGAGC"
        },
        {
          "jm": "Vector",
          "seq": "CGGAAAAAGAGTTGGTAGCTCTTGATCCGGCAAACAAACCACCGCTGGTAGCGGTGGTTTTTTTGTTTGCAAGCAGCAGATTACGCGCAGAAAAAAAGGATCTCAAGAAGATCCTTTGATCTTTTCTACGGGGTCTGACGCTCAGTGGAACGACGCGCGCGTAACTCACGTTAAGGGATTTTGGTCATGAGCTTGCGCCGTCCCGTCAAGTCAGCGTATTTTCGAGACGTTACGCCCCGCCCTGCCACTCATCGCAGTACTGTTGTAATTCATTAAGCATTCTGCCGACATGGAAGCCATCACAAACGGCATGATGAACCTGAATCGCCAGCGGCATCAGCACCTTGTCGCCTTGCGTATAATATTTGCCCATGGTGAAAACGGGGGCGAAGAAGTTGTCCATATTGGCCACGTTTAAATCAAAACTGGTGAAACTCACCCAGGGATTGGCTGACACGAAAAACATATTCTCAATAAATCCTTTAGGGAAATAGGCCAGGTTTTCACCGTAACACGCCACATCTTGCGAATATATGTGTAGAAACTGCCGGAAATCGTCGTGGTATTCACTCCAGAGCGATGAAAACGTTTCAGTTTGCTCATGGAAAACGGTGTAACATGGGTGAACACTATCCCATATCACCAGCTCACCGTCTTTCATTGCCATACGGAATTCTGGATGAGCATTCATCAGGCGGGCAAGAATGTGAATAAAGGCCGGATAAAACTTGTGCTTATTTTTCTTTACGGTTTTTAAAAAGGCCGTAATATCCAGCTGAACGGTCTGGTTATAGGTACATTGAGCAACTGACTGAAATGCCTCAAAATGTTCTTTACGATGCCATTGGGATATATCAACGGTGGTATATCCAGTGATTTTTTTCTCCATATTCTTCCTTTTTCAATATTATTGAAGCATTTATCAGGGTTATTGTCTCATGAGCGGATACATATTTGAATGTATTTAGAAAAATAAACAAATAGGGGTCAGTGTTACAACCAATTAACCAATTCTGATGCGCGTCTCTCCCCTTTGCCTGGCGGCAGTAGCGCGGTGGTCCCACCTGACCCCATGCCGAACTCAGAAGTGAAACGCCGTAGCGCCGATGGTAGTGTGGGGACTCCCCATGCGAGAGTAGGGAACTGCCAGGCATCAAATAAAACGAAAGGCTCAGTCGAAAGACTGGGCCTTTCGCCCGGGCTAATTAGGGGGTGTCGCCCTTCGCTGAATCACTGCCCGCTTTCCAGTCGGGAAACCTGTCGTGCCAGCTGCATTAATGAATCGGCCAACGCGCGGGGAGAGGCGGTTTGCGTATTGGGCGCCAGGGTGGTTTTTCTTTTCACCAGTGAGACTGGCAACAGCTGATTGCCCTTCACCGCCTGGCCCTGAGAGAGTTGCAGCAAGCGGTCCACGCTGGTTTGCCCCAGCAGGCGAAAATCCTGTTTGATGGTGGTTAACGGCGGGATATAACATGAGCTATCTTCGGTATCGTCGTATCCCACTACCGAGATATCCGCACCAACGCGCAGCCCGGACTCGGTAATGGCGCGCATTGCGCCCAGCGCCATCTGATCGTTGGCAACCAGCATCGCAGTGGGAACGATGCCCTCATTCAGCATTTGCATGGTTTGTTGAAAACCGGACATGGCACTCCAGTCGCCTTCCCGTTCCGCTATCGGCTGAATTTGATTGCGAGTGAGATATTTATGCCAGCCAGCCAGACGCAGACGCGCCGAGACAGAACTTAATGGGCCCGCTAACAGCGCGATTTGCTGGTGACCCAATGCGACCAGATGCTCCACGCCCAGTCGCGTACCGTCCTCATGGGAGAAAATAATACTGTTGATGGGTGTCTGGTCAGAGACATCAAGAAATAACGCCGGAACATTAGTGCAGGCAGCTTCCACAGCAATGGCATCCTGGTCATCCAGCGGATAGTTAATGATCAGCCCACTGACGCGTTGCGCGAGAAGATTGTGCACCGCCGCTTTACAGGCTTCGACGCCGCTTCGTTCTACCATCGACACCACCACGCTGGCACCCAGTTGATCGGCGCGAGATTTAATCGCCGCGACAATTTGCGACGGCGCGTGCAGGGCCAGACTGGAGGTGGCAACGCCAATCAGCAACGACTGTTTGCCCGCCAGTTGTTGTGCCACGCGGTTGGGAATGTAATTCAGCTCCGCCATCGCCGCTTCCACTTTTTCCCGCGTTTTCGCAGAAACGTGGCTGGCCTGGTTCACCACGCGGGAAACGGTCTGATAAGAGACACCGGCATACTCTGCGACATCGTATAACGTTACTGGTTTCATTAAGTGGTGGGACTTAACTGAGAGAAGGCCCGCGATAGTACGGTATGGCGCTTTCCAAAACGCGGTAAGCTACCGCGCGGCGCAATTTGTTTTAATAAAGATCTCCCTTTGGCAGCGAGAAGAGCGACGTCCACATATACCTGCCGTTCACTATTATTTAGTGAAATGAGATATTATGATATTTTCTGAATTGTGATTAAAAAGGCAACTTTATGCCCATGCAACAGAAACTATAAAAAATACAGAGAATGAAAAGAAACAGATAGATTTTTTAGTTCTTTAGGCCCGTAGTCTGCAAATCCTTTTATGATTTTCTATCAAACAAAAGAGGAAAATAGACCAGTTGCAATCCAAACGAGAGTCTAATAGAATGAGGTCGAAAAGTAAATCGCGCGGGTTTGTTACTGATAAAGCAGGCAAGACCTAAAATGTGTAAAGGGCAAAGTGTATACTTTGGCGTCACCCCTTACATATTTTAGGTCTTTTTTTATTGTGCGTAACTAACTTGCCATCTTCAAACAGGAGGGCTGGAAGAAGCAGACCGCTAACACAGTACATAAAAAAGGAGACATGAACGATGAACATCAAAAAGTTTGCAAAACAAGCAACAGTATTAACCTTTACTACCGCACTGCTGGCAGGAGGCGCAACTCAAGCGTTTGCGAAAGAAACGAACCAAAAGCCATATAAGGAAACATACGGCATTTCCCATATTACACGCCATGATATGCTGCAAATCCCTGAACAGCAAAAAAATGAAAAATATCAAGTTCCTGAATTCGATTCGTCCACAATTAAAAATATCTCTTCTGCAAAAGGCCTGGACGTTTGGGACAGCTGGCCATTACAAAACGCTGACGGCACTGTCGCAAACTATCACGGCTACCACATCGTCTTTGCATTAGCCGGAGATCCTAAAAATGCGGATGACACATCGATTTACATGTTCTATCAAAAAGTCGGCGAAACTTCTATTGACAGCTGGAAAAACGCTGGCCGCGTCTTTAAAGACAGCGACAAATTCGATGCAAATGATTCTATCCTAAAAGACCAAACACAAGAATGGTCAGGTTCAGCCACATTTACATCTGACGGAAAAATCCGTTTATTCTACACTGATTTCTCCGGTAAACATTACGGCAAACAAACACTGACAACTGCACAAGTTAACGTATCAGCATCAGACAGCTCTTTGAACATCAACGGTGTAGAGGATTATAAATCAATCTTTGACGGTGACGGAAAAACGTATCAAAATGTACAGCAGTTCATCGATGAAGGCAACTACAGCTCAGGCGACAACCATACGCTGAGAGATCCTCACTACGTAGAAGATAAAGGCCACAAATACTTAGTATTTGAAGCAAACACTGGAACTGAAGATGGCTACCAAGGCGAAGAATCTTTATTTAACAAAGCATACTATGGCAAAAGCACATCATTCTTCCGTCAAGAAAGTCAAAAACTTCTGCAAAGCGATAAAAAACGCACGGCTGAGTTAGCAAACGGCGCTCTCGGTATGATTGAGCTAAACGATGATTACACACTGAAAAAAGTGATGAAACCGCTGATTGCATCTAACACAGTAACAGATGAAATTGAACGCGCGAACGTCTTTAAAATGAACGGCAAATGGTACCTGTTCACTGACTCCCGCGGATCAAAAATGACGATTGACGGCATTACGTCTAACGATATTTACATGCTTGGTTATGTTTCTAATTCTTTAACTGGCCCATACAAGCCGCTGAACAAAACTGGCCTTGTGTTAAAAATGGATCTTGATCCTAACGATGTAACCTTTACTTACTCACACTTCGCTGTACCTCAAGCGAAAGGAAACAATGTCGTGATTACAAGCTATATGACAAACAGAGGATTCTACGCAGACAAACAATCAACGTTTGCGCCAAGCTTCCTGCTGAACATCAAAGGCAAGAAAACATCTGTTGTCAAAGACAGCATCCTTGAACAAGGACAATTAACAGTTAACAAATAAAAACGCAAAAGAAAATGCCGATATCCTATTGGCATTGACGCTCTTCAGGTTAAAAAGCAAGCTGATAAACCGATACAATTAAAGGCTCCTTTTGGAGCCTTTTTTTTTGGAGATTTTCAACATGAAAAAATTATTATTTGATGATCAGATAGCGGCGGGGAACTGCCAGACATCAAATAAAACAAAAGGCTCAGTCGGAAGACTGGGCCTTTTGTTTTATCTGTTGTTTGTCGGTGAACACTCTCCCGGCGGTGAGACCCGTCAAAAGGGCGACACAAAATTTATTCTAAATGCATAATAAATACTGATAACATCTTATAGTTTGTATTATATTTTGTATTATCGTTGACATGTATAATTTTGATATCAAAAACTGATTTTCCCTTTATTATTTTCGAGATTTATTTTCTTAATTCTCTTTAACAAACTAGAAATATTGTATATACAAAAAATCATAAATAATAGATGAATAGTTTAATTATAGGTGTTCATCAATCGAAAAAGCAACGTATCTTATTTAAAGTGCGTTGCTTTTTTCTCATTTATAAGGTTAAATAATTCTCATATATCAAGCAAAGTGACAGGCGCCCTTAAATATTCTGACAAATGCTCTTTCCCTAAACTCCCCCCATAAAAAAACCCGCCGAAGCGGGTTTTTACGTTATTTGCGGATTAACGATTACTCGTTATCAGAACCGCCCAGGGGGCCCGAGCTTAAGACTGGCCGTCGTTTTACAACACAGAAAGAGTTTGTAGAAACGCAAAAAGGCCATCCGTCAGGGGCCTTCTGCTTAGTTTGATGCCTGGCAGTTCCCTACTCTCGCCTTCCGCTTCCTCGCTCACTGACTCGCTGCGCTCGGTCGTTCGGCTGCGGCGAGCGGTATCAGCTCACTCAAAGGCGGTAATACGGTTATCCACAGAATCAGGGGATAACGCAGGAAAGAACATGTGAGCAAAAGGCCAGCAAAAGGCCAGGAACCGTAAAAAGGCCGCGTTGCTGGCGTTTTTCCATAGGCTCCGCCCCCCTGACGAGCATCACAAAAATCGACGCTCAAGTCAGAGGTGGCGAAACCCGACAGGACTATAAAGATACCAGGCGTTTCCCCCTGGAAGCTCCCTCGTGCGCTCTCCTGTTCCGACCCTGCCGCTTACCGGATACCTGTCCGCCTTTCTCCCTTCGGGAAGCGTGGCGCTTTCTCATAGCTCACGCTGTAGGTATCTCAGTTCGGTGTAGGTCGTTCGCTCCAAGCTGGGCTGTGTGCACGAACCCCCCGTTCAGCCCGACCGCTGCGCCTTATCCGGTAACTATCGTCTTGAGTCCAACCCGGTAAGACACGACTTATCGCCACTGGCAGCAGCCACTGGTAACAGGATTAGCAGAGCGAGGTATGTAGGCGGTGCTACAGAGTTCTTGAAGTGGTGGGCTAACTACGGCTACACTAGAAGAACAGTATTTGGTATCTGCGCTCTGCTGAAGCCAGTTACCTT"
        }
]
`

var data1 = `
[
        {
          "jm": "promoter_1",
          "seq": "GCTCTTCATAGTATGGGGAAAACCTTTTAACTCCATGATAGGGGCATTACGGAGATGGGATGAGATCACTTATACTCGTGTAATAGTAAGTTCCCGCAGTCCCTCTAACTCACAGCGGTTGTACGTCGACGGACACACGCGGACTGAAGAGC"
        },
        {
          "jm": "rbs_1",
          "seq": "GCTCTTCAGACTTCACCTAGTCGACATGACTAGGCCTCTGAGCTTTCTCTAGAAATTGGTGTTTGGCAGTCACAAAACCGCAATAGTATGACTCCACGAGGAAAATGTTACTCCCACGTTTCACTTGCACATTGTCGCTGAAGAGC"
        },
        {
          "jm": "cds_1",
          "seq": "GCTCTTCACGCGATTCTGACTGACCCGGGAACTCACCTCCACACCCGTTGCTTACGGGACAGCCAGAATATTTAAGGATGTCGAGCACTACCCCAGCCAAAATACTTGCTCTAATCAATTGGACGGCGGCGCGAAGGAGGAATTTTGAAGAGC"
        },
        {
          "jm": "Vector",
          "seq": "GCTCTTCATTTGTACTCACCTTCAGTCTCAGTGCAGCGGGGCTAATCGTAATTTTCATCCTTGATCGGTACCTAGGAACGGTATACGGGAAATACGTTAGGGCTCGCGGACAAAACATCCCTGGGCATCTAAGTGGATCCACCCCTCGTGGATTAGCTTTCGATGTGACAACCTACCACCACCTTCATCCTCTCGCTGGTGCTGGCACAGCATATGTCGGTGATTCAGTTGTGGCGACTGATGGATCTACACGCTCCAGGGGGGGTGGGCTTCCATAAAAAAGGAGCGGGGGAATTTATCCTGGAACACCGCACAGGGCCTGCGCGACGGAGATCGAGTTTTTTCAAGCCGGAATGTTAAACGCCCACCCACGTGTATAAATTAATGAGCAAGATCTCGCTCTGTTAACGGTGAAAATGATTAGCAGGAAGGCATAAATGTGTAGTGGAAGCACGGGGTGTATATAGGTACAGATTACATGCCAACCTCTGGTTGGAACTGTTTGAAGATGGGTTGTAAGATACATGGGGGCTTGCGGACGTGGGAGCCATCAGGAATACCGCGATCACTACGACCACAGTAGTGAAGAGC"
        }
]
`

func TestAssembly(t *testing.T) {
	for idx, data := range []string{data0, data1} {
		var parts []wtype.DNASequence
		if err := json.Unmarshal([]byte(data), &parts); err != nil {
			t.Fatal(err)
		}

		last := len(parts) - 1
		output, count, _, seq, err := Assemblysimulator(Assemblyparameters{
			Enzymename:   "sapi",
			Vector:       parts[last],
			Partsinorder: parts[:last],
		})

		if err != nil {
			t.Fatalf("input %d: %s: %s", idx, err, output)
		}
		if count != 1 {
			t.Fatalf("input %d: expected successful assembly but none found", idx)
		}

		bps := 0
		expected_len := 0
		for _, part := range parts {
			bps += len(part.Seq)
			tp := sisplitfwd(part.Seq)
			tp = revcmp(sisplitfwd(revcmp(tp)))
			tp2 := ""
			if tp == "" {
				tp = sisplitfwd(part.Seq)
				tp2 = revcmp(sisplitfwd(revcmp(part.Seq)))
				expected_len += len(tp) + len(tp2) - 3
			} else {
				expected_len += len(tp) - 3
			}
		}

		if len(seq.Seq) != expected_len {
			t.Fatal(fmt.Sprintf("Data %d error: length %d is not equal to expected %d", idx, len(seq.Seq), expected_len))
		}
	}
}

func sisplitfwd(s string) string {
	sa := strings.SplitAfter(s, "GCTCTTC")
	if len(sa) != 2 {
		return ""
	} else {
		return sa[1][1:]
	}
}

func cmp(s string) string {
	switch s {
	case "A":
		return "T"
	case "C":
		return "G"
	case "G":
		return "C"
	case "T":
		return "A"
	}
	return "N"
}

func revcmp(s string) string {
	s2 := ""

	for k := len(s) - 1; k >= 0; k-- {
		s2 += string(cmp(string(s[k])))
	}

	return s2
}

func siprint(s string) {
	s = strings.Replace(s, "GCTCTTC", "gctcttc", -1)
	s = strings.Replace(s, "GAAGAGC", "gaagagc", -1)
	fmt.Println(s)
}