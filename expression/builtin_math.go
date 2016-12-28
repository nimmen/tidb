// Copyright 2013 The ql Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSES/QL-LICENSE file.

// Copyright 2015 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package expression

import (
	"math"
	"math/rand"

	"github.com/juju/errors"
	"github.com/pingcap/tidb/context"
	"github.com/pingcap/tidb/util/types"
)

type absFuncClass struct {
	baseFuncClass
}

type builtinAbs struct {
	baseBuiltinFunc
}

func (b *absFuncClass) getFunction(args []Expression, ctx context.Context) (builtinFunc, error) {
	err := b.checkValid(args)
	if err != nil {
		return nil, errors.Trace(err)
	}
	f := &builtinAbs{baseBuiltinFunc: newBaseBuiltinFunc(args, true, ctx)}
	f.self = f
	return f, nil
}

// See http://dev.mysql.com/doc/refman/5.7/en/mathematical-functions.html#function_abs
func (b *builtinAbs) eval(args []types.Datum) (d types.Datum, err error) {
	if args, err = b.evalArgs(args); err != nil {
		return d, errors.Trace(err)
	}
	d = args[0]
	switch d.Kind() {
	case types.KindNull:
		return d, nil
	case types.KindUint64:
		return d, nil
	case types.KindInt64:
		iv := d.GetInt64()
		if iv >= 0 {
			d.SetInt64(iv)
			return d, nil
		}
		d.SetInt64(-iv)
		return d, nil
	default:
		// we will try to convert other types to float
		// TODO: if time has no precision, it will be a integer
		f, err := d.ToFloat64(b.ctx.GetSessionVars().StmtCtx)
		d.SetFloat64(math.Abs(f))
		return d, errors.Trace(err)
	}
}

type ceilFuncClass struct {
	baseFuncClass
}

type builtinCeil struct {
	baseBuiltinFunc
}

func (b *ceilFuncClass) getFunction(args []Expression, ctx context.Context) (builtinFunc, error) {
	err := b.checkValid(args)
	if err != nil {
		return nil, errors.Trace(err)
	}
	f := &builtinCeil{baseBuiltinFunc: newBaseBuiltinFunc(args, true, ctx)}
	f.self = f
	return f, nil
}

// See http://dev.mysql.com/doc/refman/5.7/en/mathematical-functions.html#function_ceiling
func (b *builtinCeil) eval(args []types.Datum) (d types.Datum, err error) {
	if args, err = b.evalArgs(args); err != nil {
		return d, errors.Trace(err)
	}
	if args[0].IsNull() ||
		args[0].Kind() == types.KindUint64 || args[0].Kind() == types.KindInt64 {
		return args[0], nil
	}

	f, err := args[0].ToFloat64(b.ctx.GetSessionVars().StmtCtx)
	if err != nil {
		return d, errors.Trace(err)
	}
	d.SetFloat64(math.Ceil(f))
	return
}

type logFuncClass struct {
	baseFuncClass
}

type builtinLog struct {
	baseBuiltinFunc
}

func (b *logFuncClass) getFunction(args []Expression, ctx context.Context) (builtinFunc, error) {
	err := b.checkValid(args)
	if err != nil {
		return nil, errors.Trace(err)
	}
	f := &builtinLog{baseBuiltinFunc: newBaseBuiltinFunc(args, true, ctx)}
	f.self = f
	return f, nil
}

// See http://dev.mysql.com/doc/refman/5.7/en/mathematical-functions.html#function_log
func (b *builtinLog) eval(args []types.Datum) (d types.Datum, err error) {
	if args, err = b.evalArgs(args); err != nil {
		return d, errors.Trace(err)
	}
	sc := b.ctx.GetSessionVars().StmtCtx
	switch len(args) {
	case 1:
		x, err := args[0].ToFloat64(sc)
		if err != nil {
			return d, errors.Trace(err)
		}

		if x <= 0 {
			return d, nil
		}

		d.SetFloat64(math.Log(x))
		return d, nil
	case 2:
		b, err := args[0].ToFloat64(sc)
		if err != nil {
			return d, errors.Trace(err)
		}

		x, err := args[1].ToFloat64(sc)
		if err != nil {
			return d, errors.Trace(err)
		}

		if b <= 1 || x <= 0 {
			return d, nil
		}

		d.SetFloat64(math.Log(x) / math.Log(b))
		return d, nil
	}
	return
}

type log2FuncClass struct {
	baseFuncClass
}

type builtinLog2 struct {
	baseBuiltinFunc
}

func (b *log2FuncClass) getFunction(args []Expression, ctx context.Context) (builtinFunc, error) {
	err := b.checkValid(args)
	if err != nil {
		return nil, errors.Trace(err)
	}
	f := &builtinLog2{baseBuiltinFunc: newBaseBuiltinFunc(args, true, ctx)}
	f.self = f
	return f, nil
}

// See http://dev.mysql.com/doc/refman/5.7/en/mathematical-functions.html#function_log2
func (b *builtinLog2) eval(args []types.Datum) (d types.Datum, err error) {
	if args, err = b.evalArgs(args); err != nil {
		return d, errors.Trace(err)
	}
	sc := b.ctx.GetSessionVars().StmtCtx
	x, err := args[0].ToFloat64(sc)
	if err != nil {
		return d, errors.Trace(err)
	}

	if x <= 0 {
		return
	}

	d.SetFloat64(math.Log2(x))
	return
}

type log10FuncClass struct {
	baseFuncClass
}

type builtinLog10 struct {
	baseBuiltinFunc
}

func (b *log10FuncClass) getFunction(args []Expression, ctx context.Context) (builtinFunc, error) {
	err := b.checkValid(args)
	if err != nil {
		return nil, errors.Trace(err)
	}
	f := &builtinLog10{baseBuiltinFunc: newBaseBuiltinFunc(args, true, ctx)}
	f.self = f
	return f, nil
}

// See http://dev.mysql.com/doc/refman/5.7/en/mathematical-functions.html#function_log10
func (b *builtinLog10) eval(args []types.Datum) (d types.Datum, err error) {
	if args, err = b.evalArgs(args); err != nil {
		return d, errors.Trace(err)
	}
	sc := b.ctx.GetSessionVars().StmtCtx
	x, err := args[0].ToFloat64(sc)
	if err != nil {
		return d, errors.Trace(err)
	}

	if x <= 0 {
		return
	}

	d.SetFloat64(math.Log10(x))
	return

}

type randFuncClass struct {
	baseFuncClass
}

type builtinRand struct {
	baseBuiltinFunc
}

func (b *randFuncClass) getFunction(args []Expression, ctx context.Context) (builtinFunc, error) {
	err := b.checkValid(args)
	if err != nil {
		return nil, errors.Trace(err)
	}
	f := &builtinRand{baseBuiltinFunc: newBaseBuiltinFunc(args, false, ctx)}
	f.self = f
	return f, nil
}

// See http://dev.mysql.com/doc/refman/5.7/en/mathematical-functions.html#function_rand
func (b *builtinRand) eval(args []types.Datum) (d types.Datum, err error) {
	if args, err = b.evalArgs(args); err != nil {
		return d, errors.Trace(err)
	}
	if len(args) == 1 && !args[0].IsNull() {
		seed, err := args[0].ToInt64(b.ctx.GetSessionVars().StmtCtx)
		if err != nil {
			return d, errors.Trace(err)
		}
		rand.Seed(seed)
	}
	d.SetFloat64(rand.Float64())
	return d, nil
}

type powFuncClass struct {
	baseFuncClass
}

type builtinPow struct {
	baseBuiltinFunc
}

func (b *powFuncClass) getFunction(args []Expression, ctx context.Context) (builtinFunc, error) {
	err := b.checkValid(args)
	if err != nil {
		return nil, errors.Trace(err)
	}
	f := &builtinPow{baseBuiltinFunc: newBaseBuiltinFunc(args, true, ctx)}
	f.self = f
	return f, nil
}

// See http://dev.mysql.com/doc/refman/5.7/en/mathematical-functions.html#function_pow
func (b *builtinPow) eval(args []types.Datum) (d types.Datum, err error) {
	if args, err = b.evalArgs(args); err != nil {
		return d, errors.Trace(err)
	}
	sc := b.ctx.GetSessionVars().StmtCtx
	x, err := args[0].ToFloat64(sc)
	if err != nil {
		return d, errors.Trace(err)
	}

	y, err := args[1].ToFloat64(sc)
	if err != nil {
		return d, errors.Trace(err)
	}
	d.SetFloat64(math.Pow(x, y))
	return d, nil
}

type roundFuncClass struct {
	baseFuncClass
}

type builtinRound struct {
	baseBuiltinFunc
}

func (b *roundFuncClass) getFunction(args []Expression, ctx context.Context) (builtinFunc, error) {
	err := b.checkValid(args)
	if err != nil {
		return nil, errors.Trace(err)
	}
	f := &builtinRound{baseBuiltinFunc: newBaseBuiltinFunc(args, true, ctx)}
	f.self = f
	return f, nil
}

// See http://dev.mysql.com/doc/refman/5.7/en/mathematical-functions.html#function_round
func (b *builtinRound) eval(args []types.Datum) (d types.Datum, err error) {
	if args, err = b.evalArgs(args); err != nil {
		return d, errors.Trace(err)
	}
	sc := b.ctx.GetSessionVars().StmtCtx
	x, err := args[0].ToFloat64(sc)
	if err != nil {
		return d, errors.Trace(err)
	}

	dec := 0
	if len(args) == 2 {
		y, err1 := args[1].ToInt64(sc)
		if err1 != nil {
			return d, errors.Trace(err1)
		}
		dec = int(y)
	}
	d.SetFloat64(types.Round(x, dec))
	return d, nil
}
