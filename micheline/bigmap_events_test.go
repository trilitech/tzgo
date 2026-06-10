package micheline

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/trilitech/tzgo/tezos"
)

func TestBigmapEvent_JSON_UpdateWithEmptyValueBecomesRemove(t *testing.T) {
	var e BigmapEvent
	// value omitted -> invalid Prim -> translated to remove
	if err := json.Unmarshal([]byte(`{"action":"update","big_map":"1"}`), &e); err != nil {
		t.Fatal(err)
	}
	if e.Action != DiffActionRemove {
		t.Fatalf("action=%s want=%s", e.Action, DiffActionRemove)
	}
}

func TestBigmapEvent_MarshalJSON_Actions(t *testing.T) {
	key := NewString("k")
	val := NewString("v")
	kh := tezos.MustParseExprHash("exprtrNvsCL6MY7rfibVC6t8uqVJgAXfjGDyBzxonZeUpiunP5P9KC")

	// update
	eu := BigmapEvent{
		Action:  DiffActionUpdate,
		Id:      7,
		KeyHash: kh,
		Key:     key,
		Value:   val,
	}
	ubuf, err := json.Marshal(eu)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(ubuf, []byte(`"action":"update"`)) {
		t.Fatalf("missing action in update json: %s", string(ubuf))
	}

	// alloc
	ea := BigmapEvent{
		Action:    DiffActionAlloc,
		Id:        7,
		KeyType:   Prim{Type: PrimNullary, OpCode: T_STRING},
		ValueType: Prim{Type: PrimNullary, OpCode: T_BYTES},
	}
	abuf, err := json.Marshal(ea)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(abuf, []byte(`"action":"alloc"`)) {
		t.Fatalf("missing action in alloc json: %s", string(abuf))
	}

	// copy
	ec := BigmapEvent{
		Action:   DiffActionCopy,
		SourceId: 1,
		DestId:   2,
	}
	cbuf, err := json.Marshal(ec)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(cbuf, []byte(`"action":"copy"`)) {
		t.Fatalf("missing action in copy json: %s", string(cbuf))
	}
}

func TestBigmapEvents_BinaryRoundtrip_And_KeyHelpers(t *testing.T) {
	// Use a real mainnet bigmap fixture to get a realistic key/value encoding.
	buf, err := os.ReadFile(filepath.Join("testdata-mainnet", "bigmap", "KT1AFA2mwNUMNd4SsujE1YYp29vd8BZejyKW-515.json"))
	if err != nil {
		t.Fatal(err)
	}
	var tests []testcase
	if err := json.Unmarshal(buf, &tests); err != nil {
		t.Fatal(err)
	}
	if len(tests) == 0 {
		t.Fatal("fixture is empty")
	}
	tc := tests[0]

	// derive expr hash from testcase name suffix (contains "...-expr<hash>")
	var kh tezos.ExprHash
	if i := strings.LastIndex(tc.Name, "-expr"); i >= 0 {
		// include the "expr" prefix (hash string starts at i+1)
		h := tc.Name[i+1:]
		kh = tezos.MustParseExprHash(h)
	} else {
		t.Fatalf("unexpected testcase name format: %q", tc.Name)
	}

	// build key type from big_map type: args[0] is key type
	bmTyp := checkTypeEncoding(t, tc)
	if bmTyp.OpCode != T_BIG_MAP || len(bmTyp.Args) != 2 {
		t.Fatalf("unexpected fixture big_map type: %s (%d args)", bmTyp.OpCode, len(bmTyp.Args))
	}
	keyTyp := Type{bmTyp.Args[0]}

	key := checkKeyEncoding(t, tc)
	val := checkValueEncoding(t, tc)

	ev := BigmapEvents{
		{
			Action:  DiffActionUpdate,
			Id:      515,
			KeyHash: kh,
			Key:     key,
			Value:   val,
		},
		{
			Action:    DiffActionAlloc,
			Id:        516,
			KeyType:   bmTyp.Args[0],
			ValueType: bmTyp.Args[1],
		},
		{
			Action:   DiffActionCopy,
			SourceId: 1,
			DestId:   2,
		},
		{
			Action:  DiffActionRemove,
			Id:      515,
			KeyHash: kh,
			Key:     key,
		},
	}

	// Filter
	if got := ev.Filter(515); len(got) != 2 {
		t.Fatalf("Filter(515) len=%d want=2", len(got))
	}

	// Encoding for each event
	if ev[0].Encoding() != ev[0].Key.OpCode.PrimType() {
		t.Fatalf("Encoding(update) mismatch")
	}
	if ev[1].Encoding() != ev[1].KeyType.OpCode.PrimType() {
		t.Fatalf("Encoding(alloc) mismatch")
	}

	// GetKey/GetKeyPtr should not panic and should return a key for a valid type.
	k := ev[0].GetKey(keyTyp)
	if !k.Type.IsValid() {
		t.Fatalf("GetKey returned invalid key type")
	}
	kp := ev[0].GetKeyPtr(keyTyp)
	if kp == nil || !kp.Type.IsValid() {
		t.Fatalf("GetKeyPtr returned nil/invalid key")
	}

	// Binary roundtrip
	bin, err := ev.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	var out BigmapEvents
	if err := out.UnmarshalBinary(bin); err != nil {
		t.Fatal(err)
	}
	if len(out) != len(ev) {
		t.Fatalf("roundtrip len=%d want=%d", len(out), len(ev))
	}
}

func TestBigmapEvents_UnmarshalBinary_Errors(t *testing.T) {
	// Build a binary blob with invalid prim payload (not a pair).
	buf := bytes.NewBuffer(nil)
	_ = binary.Write(buf, binary.BigEndian, uint32(1))                  // id
	_ = buf.WriteByte(byte(DiffActionUpdate))                           // action
	_ = Prim{Type: PrimInt, Int: tezos.NewZ(1).Big()}.EncodeBuffer(buf) // invalid for decoder expectations

	var out BigmapEvents
	if err := out.UnmarshalBinary(buf.Bytes()); err == nil {
		t.Fatalf("expected error for invalid keypair prim")
	}
}
