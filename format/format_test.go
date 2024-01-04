package format

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/vmihailenco/msgpack/v5"
)

type formater struct {
	Marshal   func(v interface{}) ([]byte, error)
	Unmarshal func(data []byte, v interface{}) error
}

var formatters = map[string]formater{
	"JSON":    {Marshal: json.Marshal, Unmarshal: json.Unmarshal},
	"MsgPack": {Marshal: msgpack.Marshal, Unmarshal: json.Unmarshal},
}

type testPage struct {
	Length int
	Items  []testItem
	Next   string
}

type testItem struct {
	IntA   uint8
	IntB   uint16
	IntC   uint32
	IntD   uint64
	FloatA float32
	FloatB float64
	Str    string
	StrArr []string
	Blob   []byte
}

var result []byte

var sizes = []int{5, 10, 100, 1000, 10000}

func genInput(size int) testPage {
	items := []testItem{}
	for i := 0; i < size; i++ {
		items = append(items, testItem{0xFF, 0xFFFF, 0xFFFFFFFF, 0xFFFFFFFFFFFFFFFF, (1 << 32), (1 << 64), "abcdeffdaslfualnxoiudfljasdfl", []string{"hello", "world", "abcdeffdaslfualnxoiudfljasdfl"}, []byte("sfldfasdfasd;lfkhdlkfjas")})
	}
	return testPage{
		Length: size,
		Items:  items,
		Next:   "some-string",
	}
}

func BenchmarkFormatters(b *testing.B) {
	for name, fns := range formatters {
		for _, size := range sizes {
			input := genInput(size)
			b.Run(fmt.Sprintf("%s--size_%d", name, size), func(b *testing.B) {
				var bs []byte
				bsTotal := 0
				for n := 0; n < b.N; n++ {
					bs, _ = fns.Marshal(input)
					bsTotal += len(bs)
					ti := testPage{}
					fns.Unmarshal(bs, &ti)
				}
				result = bs
				b.ReportMetric(float64(bsTotal)/float64(b.N)/(1<<10), "KB")
			})
		}
	}
}
