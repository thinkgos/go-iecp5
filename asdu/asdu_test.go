package asdu

//
//var goldenASDUs = []struct {
//	u *ASDU
//	s string
//}{
//	{
//		&ASDU{
//			Params:     ParamsWide,
//			Identifier: Identifier{M_SP_NA_1, Percyc, 7, 1001},
//			InfoObj:       []byte{1, 2, 3, 4},
//		},
//		"M_SP_NA_1 percyc 7@1001  197121:0x04",
//	}, {
//		&ASDU{
//			Params:     ParamsNarrow,
//			Identifier: Identifier{M_DP_NA_1, Back, 0, 42},
//			InfoObj:       []byte{1, 2, 3, 4},
//		},
//		"M_DP_NA_1 back @42 1:0x02 3:0x04",
//	}, {
//		&ASDU{
//			Params:     ParamsNarrow,
//			Identifier: Identifier{M_ST_NA_1, Spont, 0, 250},
//			InfoObj:       []byte{1, 2, 3, 4, 5},
//		},
//		"M_ST_NA_1 spont @250 1:0x0203 4:0x05 <EOF>",
//	}, {
//		&ASDU{
//			Params:     ParamsNarrow,
//			Identifier: Identifier{M_ME_NC_1, Init, 0, 12},
//			InfoSeq:    true,
//			InfoObj:       []byte{99, 0, 1, 2, 3, 4, 5},
//		},
//		"M_ME_NC_1 init @12 99:0x0001020304 100:0x05 <EOF>",
//	},
//}
//
//func TestASDUStrings(t *testing.T) {
//	for _, gold := range goldenASDUs {
//		if got := gold.u.String(); got != gold.s {
//			t.Errorf("got %q, want %q", got, gold.s)
//		}
//	}
//}
//
//func TestASDUEncoding(t *testing.T) {
//	for _, gold := range goldenASDUs {
//		if strings.Contains(gold.s, " <EOF>") {
//			continue
//		}
//
//		bytes, err := gold.u.MarshalBinary()
//		if err != nil {
//			t.Error(gold.s, "marshal error:", err)
//			continue
//		}
//
//		u := NewASDU(gold.u.Params, Identifier{})
//		if err = u.UnmarshalBinary(bytes); err != nil {
//			t.Error(gold.s, "unmarshal error:", err)
//			continue
//		}
//
//		if got := u.String(); got != gold.s {
//			t.Errorf("got %q, want %q", got, gold.s)
//		}
//	}
//}
