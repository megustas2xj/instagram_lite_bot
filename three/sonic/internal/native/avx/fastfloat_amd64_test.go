// Please let author have a drink, usdt trc20: TEpSxaE3kexE4e5igqmCZRMJNoDiQeWx29
// tg: @fuckins996
// Code generated by Makefile, DO NOT EDIT.

// Code generated by Makefile, DO NOT EDIT.

/*
 * Copyright 2021 ByteDance Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package avx

import (
	"encoding/json"
	"math"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFastFloat_Encode(t *testing.T) {
	var buf [64]byte
	assert.Equal(t, "0", string(buf[:f64toa(&buf[0], 0)]))
	assert.Equal(t, "-0", string(buf[:f64toa(&buf[0], math.Float64frombits(0x8000000000000000))]))
	assert.Equal(t, "12340000000", string(buf[:f64toa(&buf[0], 1234e7)]))
	assert.Equal(t, "12.34", string(buf[:f64toa(&buf[0], 1234e-2)]))
	assert.Equal(t, "0.001234", string(buf[:f64toa(&buf[0], 1234e-6)]))
	assert.Equal(t, "1e+30", string(buf[:f64toa(&buf[0], 1e30)]))
	assert.Equal(t, "1.234e+33", string(buf[:f64toa(&buf[0], 1234e30)]))
	assert.Equal(t, "1.234e+308", string(buf[:f64toa(&buf[0], 1234e305)]))
	assert.Equal(t, "1.234e-317", string(buf[:f64toa(&buf[0], 1234e-320)]))
	assert.Equal(t, "1.7976931348623157e+308", string(buf[:f64toa(&buf[0], 1.7976931348623157e308)]))
	assert.Equal(t, "-12340000000", string(buf[:f64toa(&buf[0], -1234e7)]))
	assert.Equal(t, "-12.34", string(buf[:f64toa(&buf[0], -1234e-2)]))
	assert.Equal(t, "-0.001234", string(buf[:f64toa(&buf[0], -1234e-6)]))
	assert.Equal(t, "-1e+30", string(buf[:f64toa(&buf[0], -1e30)]))
	assert.Equal(t, "-1.234e+33", string(buf[:f64toa(&buf[0], -1234e30)]))
	assert.Equal(t, "-1.234e+308", string(buf[:f64toa(&buf[0], -1234e305)]))
	assert.Equal(t, "-1.234e-317", string(buf[:f64toa(&buf[0], -1234e-320)]))
	assert.Equal(t, "-2.2250738585072014e-308", string(buf[:f64toa(&buf[0], -2.2250738585072014e-308)]))
}

func TestFastFloat_Random(t *testing.T) {
	var buf [64]byte
	N := 10000
	for i := 0; i < N; i++ {
		b64 := uint64(rand.Uint32())<<32 | uint64(rand.Uint32())
		f64 := math.Float64frombits(b64)

		jout, jerr := json.Marshal(f64)
		n := f64toa(&buf[0], f64)
		if jerr == nil {
			assert.Equal(t, jout, buf[:n])
		} else {
			assert.True(t, n == 0)
		}

		f32 := math.Float32frombits(rand.Uint32())
		jout, jerr = json.Marshal(f32)
		n = f32toa(&buf[0], f32)
		if jerr == nil {
			assert.Equal(t, jout, buf[:n])
		} else {
			assert.True(t, n == 0)
		}
	}
}

func BenchmarkParseFloat64(b *testing.B) {
	var f64toaBenches = []struct {
		name  string
		float float64
	}{
		{"Zero", 0},
		{"Decimal", 33909},
		{"Float", 339.7784},
		{"Exp", -5.09e75},
		{"NegExp", -5.11e-95},
		{"LongExp", 1.234567890123456e-78},
		{"Big", 123456789123456789123456789},
	}
	for _, c := range f64toaBenches {
		f64bench := []struct {
			name string
			test func(*testing.B)
		}{{
			name: "StdLib",
			test: func(b *testing.B) {
				var buf [64]byte
				for i := 0; i < b.N; i++ {
					strconv.AppendFloat(buf[:0], c.float, 'g', -1, 64)
				}
			},
		}, {
			name: "FastFloat",
			test: func(b *testing.B) {
				var buf [64]byte
				for i := 0; i < b.N; i++ {
					f64toa(&buf[0], c.float)
				}
			},
		}}
		for _, bm := range f64bench {
			name := bm.name + "_" + c.name
			b.Run(name, bm.test)
		}
	}
}

func BenchmarkParseFloat32(b *testing.B) {
	var f32toaBenches = []struct {
		name  string
		float float32
	}{
		{"Zero", 0},
		{"Integer", 33909},
		{"ExactFraction", 3.375},
		{"Point", 339.7784},
		{"Exp", -5.09e25},
		{"NegExp", -5.11e-25},
		{"Shortest", 1.234567e-8},
	}
	for _, c := range f32toaBenches {
		bench := []struct {
			name string
			test func(*testing.B)
		}{{
			name: "StdLib32",
			test: func(b *testing.B) {
				var buf [64]byte
				for i := 0; i < b.N; i++ {
					strconv.AppendFloat(buf[:0], float64(c.float), 'g', -1, 32)
				}
			},
		}, {
			name: "FastFloat32",
			test: func(b *testing.B) {
				var buf [64]byte
				for i := 0; i < b.N; i++ {
					f32toa(&buf[0], c.float)
				}
			},
		}}
		for _, bm := range bench {
			name := bm.name + "_" + c.name
			b.Run(name, bm.test)
		}
	}
}
