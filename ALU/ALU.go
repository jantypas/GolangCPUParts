package Onyx1ALU

import "math"

const (
	ALU_FLAGS_ERROR        = 0x0000_0000_0000_0001
	ALU_FLAGS_ZERO         = 0x0000_0000_0000_0002
	ALU_FLAGS_NEGATIVE     = 0x0000_0000_0000_0004
	ALU_FLAGS_CARRY        = 0x0000_0000_0000_0008
	ALU_FLAGS_DIVIDEBYZERO = 0x0000_0000_0000_0010
	ALU_FLAGS_INVALIDOP    = 0x0000_0000_0000_0020
)

const (
	ALU_OP_ADDINT64  = 0x0000_0000_0000_0001
	ALU_OP_SUBINT64  = 0x0000_0000_0000_0002
	ALU_OP_MULTINT64 = 0x0000_0000_0000_0003
	ALU_OP_DIVINT64  = 0x0000_0000_0000_0004
	ALU_OP_ANDINT64  = 0x0000_0000_0000_0005
	ALU_OP_NOTINT64  = 0x0000_0000_0000_0006
	ALU_OP_ORINT64   = 0x0000_0000_0000_0007
	ALU_OP_XORINT64  = 0x0000_0000_0000_0008
	ALU_OP_SHLINT64  = 0x0000_0000_0000_0009
	ALU_OP_SHRINT64  = 0x0000_0000_0000_000A
	ALU_OP_FADD64    = 0x0000_0000_0000_000B
	ALU_OP_FSUB64    = 0x0000_0000_0000_000C
	ALU_OP_FMULT64   = 0x0000_0000_0000_000D
	ALU_OP_FDIV64    = 0x0000_0000_0000_000E
	ALU_OP_FSIN64    = 0x0000_0000_0000_000F
	ALU_OP_FCOS64    = 0x0000_0000_0000_0010
	ALU_OP_FTAN64    = 0x0000_0000_0000_0011
	ALU_OP_FLN64     = 0x0000_0000_0000_0012
	ALU_OP_FEXP64    = 0x0000_0000_0000_0013
	ALU_OP_FSQRT64   = 0x0000_0000_0000_0014
)

func ALUInt64(op int, parmA int64, parmB int64) (outA int64, outB int64, flags uint64) {
	switch op {
	case ALU_OP_ADDINT64:
		outA = parmA + parmB
		if outA < 0 {
			flags |= ALU_FLAGS_NEGATIVE
		}
		if outA == 0 {
			flags |= ALU_FLAGS_ZERO
		}
		return outA, parmB, flags
	case ALU_OP_SUBINT64:
		outA = parmA - parmB
		if outA < 0 {
			flags |= ALU_FLAGS_NEGATIVE
		}
		if outA == 0 {
			flags |= ALU_FLAGS_ZERO
		}
		return outA, outB, flags
	case ALU_OP_MULTINT64:
		outA = parmA * parmB
		if outA == 0 {
			flags |= ALU_FLAGS_ZERO
		}
		if outA < 0 {
			flags |= ALU_FLAGS_NEGATIVE
		}
		return outA, outB, flags
	case ALU_OP_DIVINT64:
		if parmB == 0 {
			flags |= ALU_FLAGS_ERROR
			flags |= ALU_FLAGS_DIVIDEBYZERO
			return 0, 0, flags
		}
		outA = parmA / parmB
		outB = parmA % parmB
		if outA < 0 {
			flags |= ALU_FLAGS_NEGATIVE
		}
		if outA == 0 {
			flags |= ALU_FLAGS_ZERO
		}
		return outA, outB, flags
	case ALU_OP_ANDINT64:
		outA := parmA & parmB
		if outA < 0 {
			flags |= ALU_FLAGS_NEGATIVE
		}
		if outA == 0 {
			flags |= ALU_FLAGS_ZERO
		}
		return outA, outB, flags
	case ALU_OP_NOTINT64:
		outA = ^parmA
		if outA < 0 {
			flags |= ALU_FLAGS_NEGATIVE
		}
		if outA == 0 {
			flags |= ALU_FLAGS_ZERO
		}
		return outA, outB, flags
	case ALU_OP_ORINT64:
		outA = parmA | parmB
		if outA < 0 {
			flags |= ALU_FLAGS_NEGATIVE
		}
		if outA == 0 {
			flags |= ALU_FLAGS_ZERO
		}
		return outA, outB, flags
	case ALU_OP_XORINT64:
		outA = parmA ^ parmB
		if outA < 0 {
			flags |= ALU_FLAGS_NEGATIVE
		}
		if outA == 0 {
			flags |= ALU_FLAGS_ZERO
		}
		return outA, outB, flags
	case ALU_OP_SHLINT64:
		if parmB >= 64 {
			flags |= ALU_FLAGS_ERROR
			flags |= ALU_FLAGS_INVALIDOP
			return 0, 0, flags
		}
		outA = parmA << parmB
		if outA < 0 {
			flags |= ALU_FLAGS_NEGATIVE
		}
		if outA == 0 {
			flags |= ALU_FLAGS_ZERO
		}
		return outA, outB, flags
	case ALU_OP_SHRINT64:
		if parmB >= 64 {
			flags |= ALU_FLAGS_ERROR
			flags |= ALU_FLAGS_INVALIDOP
			return 0, 0, flags
		}
		outA := parmA >> parmB
		if outA < 0 {
			flags |= ALU_FLAGS_NEGATIVE
		}
		if outA == 0 {
			flags |= ALU_FLAGS_ZERO
		}
		return outA, outB, flags
	}
	flags |= ALU_FLAGS_INVALIDOP
	flags |= ALU_FLAGS_ERROR
	return 0, 0, flags
}

func ALUFloat64(op int, parmA float64, parmB float64) (outA float64, outB float64, flags uint64) {
	switch op {
	case ALU_OP_FADD64:
		outA = parmA + parmB
		if outA < 0 {
			flags |= ALU_FLAGS_NEGATIVE
		}
		if outA == 0 {
			flags |= ALU_FLAGS_ZERO
		}
		return outA, outB, flags
	case ALU_OP_FSUB64:
		outA = parmA - parmB
		if outA < 0 {
			flags |= ALU_FLAGS_NEGATIVE
		}
		if outA == 0 {
			flags |= ALU_FLAGS_ZERO
		}
		return outA, outB, flags
	case ALU_OP_FMULT64:
		outA = parmA * parmB
		if outA < 0 {
			flags |= ALU_FLAGS_NEGATIVE
		}
		if outA == 0 {
			flags |= ALU_FLAGS_ZERO
		}
		return outA, outB, flags
	case ALU_OP_FDIV64:
		if parmB == 0 {
			flags |= ALU_FLAGS_ERROR
			flags |= ALU_FLAGS_DIVIDEBYZERO
			return 0, 0, flags
		}
		outA = parmA / parmB
		if outA < 0 {
			flags |= ALU_FLAGS_NEGATIVE
		}
		if outA == 0 {
			flags |= ALU_FLAGS_ZERO
		}
		return outA, outB, flags
	case ALU_OP_FSIN64:
		outA = math.Sin(parmA)
		if outA < 0 {
			flags |= ALU_FLAGS_NEGATIVE
		}
		if outA < 0 {
			flags |= ALU_FLAGS_NEGATIVE
		}
		return outA, outB, flags
	case ALU_OP_FCOS64:
		outA = math.Cos(parmA)
		if outA < 0 {
			flags |= ALU_FLAGS_NEGATIVE
		}
		if outA < 0 {
			flags |= ALU_FLAGS_NEGATIVE
		}
		return outA, outB, flags
	case ALU_OP_FTAN64:
		outA = math.Tan(parmA)
		if outA < 0 {
			flags |= ALU_FLAGS_NEGATIVE
		}
		if outA < 0 {
			flags |= ALU_FLAGS_NEGATIVE
		}
		return outA, outB, flags
	case ALU_OP_FEXP64:
		outA := math.Exp(parmA)
		if outA < 0 {
			flags |= ALU_FLAGS_NEGATIVE
		}
		if outA < 0 {
			flags |= ALU_FLAGS_NEGATIVE
		}
		return outA, outB, flags
	case ALU_OP_FLN64:
		outA := math.Log2(parmA)
		if outA < 0 {
			flags |= ALU_FLAGS_NEGATIVE
		}
		if outA < 0 {
			flags |= ALU_FLAGS_NEGATIVE
		}
		return outA, outB, flags
	case ALU_OP_FSQRT64:
		outA = math.Sqrt(parmA)
		if outA < 0 {
			flags |= ALU_FLAGS_NEGATIVE
		}
		if outA < 0 {
			flags |= ALU_FLAGS_NEGATIVE
		}
		return outA, outB, flags
	}
	return 0, 0, ALU_FLAGS_INVALIDOP
}
