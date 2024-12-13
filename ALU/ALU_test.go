package Onyx1ALU

import (
	"math"
	"testing"
)

func TestALUFloat64(t *testing.T) {
	parmA := 30.0
	parmB := 20.0
	outA, outB, flags := ALUFloat64(ALU_OP_FADD64, parmA, parmB)
	if outA != 50.0 {
		t.Errorf("OP_FADDD64 Expected 50.0, got %f %f %b", outA, outB, flags)
	}
	outA, outB, flags = ALUFloat64(ALU_OP_FSUB64, parmA, parmB)
	if outA != 10.0 {
		t.Errorf("OP_FSUBD64 Expected 10.0, got %f %f %b", outA, outB, flags)
	}
	outA, outB, flags = ALUFloat64(ALU_OP_FMULT64, parmA, parmB)
	if outA != 600.0 {
		t.Errorf("OP_FMULTD64 Expected 300.0, got %f %f %b", outA, outB, flags)
	}
	outA, outB, flags = ALUFloat64(ALU_OP_FDIV64, parmA, parmB)
	if outA != 1.5 {
		t.Errorf("OP_FDIVD64 Expected 1.5, got %f %f %b", outA, outB, flags)
	}
	outA, outB, flags = ALUFloat64(ALU_OP_FSIN64, parmA, parmB)
	if outA != math.Sin(parmA) {
		t.Errorf("OP_FSIND64 Expected 0.9092974268256817, got %f %f %b", outA, outB, flags)
	}
	outA, outB, flags = ALUFloat64(ALU_OP_FCOS64, parmA, parmB)
	if outA != math.Cos(parmA) {
		t.Errorf("OP_FCOSD64 Expected 0.3826834323650898, got %f %f %b", outA, outB, flags)
	}
	outA, outB, flags = ALUFloat64(ALU_OP_FTAN64, parmA, parmB)
	if outA != math.Tan(parmA) {
		t.Errorf("OP_FTAND64 Expected %f, got %f %f %b", math.Sin(parmA), outA, outB, flags)
	}
	outA, outB, flags = ALUFloat64(ALU_OP_FEXP64, parmA, parmB)
	if outA != math.Exp(parmA) {
		t.Errorf("OP_FEXPD64 Expected %f, got %f %f %b", math.Exp(parmA), outA, outB, flags)
	}
	outA, outB, flags = ALUFloat64(ALU_OP_FLN64, parmA, parmB)
	if outA != math.Log2(parmA) {
		t.Errorf("OP_FLND64 Expected %f, got %f %f %b", math.Log2(parmA), outA, outB, flags)
	}
	outA, outB, flags = ALUFloat64(ALU_OP_FSQRT64, parmA, parmB)
	if outA != math.Sqrt(parmA) {
		t.Errorf("OP_FSQRTD64 Expected %f, got %f %f %b", math.Sqrt(parmA), outA, outB, flags)
	}
}

func TestALUInt64(t *testing.T) {
	var parmA int64 = 30
	var parmB int64 = 20
	outA, outB, flags := ALUInt64(ALU_OP_ADDINT64, parmA, parmB)
	if outA != 50 {
		t.Errorf("OP_ADDINT64 Expected 50.0, got %d %d %b", outA, outB, flags)
	}
	outA, outB, flags = ALUInt64(ALU_OP_MULTINT64, parmA, parmB)
	if outA != 600 {
		t.Errorf("OP_MULTINT64 Expected 50.0, got %d %d %b", outA, outB, flags)
	}
	outA, outB, flags = ALUInt64(ALU_OP_DIVINT64, parmA, parmB)
	if outA != 1 || outB != 10 {
		t.Errorf("OP_DIVINT64 Expected 1, 10, got %d %d %b", outA, outB, flags)
	}
	outA, outB, flags = ALUInt64(ALU_OP_DIVINT64, parmA, 0)
	if outA != 0 || outB != 0 {
		t.Errorf("OP_DIVINT64 Expected 1, 10, got %d %d %b", outA, outB, flags)
	}
	outA, outB, flags = ALUInt64(ALU_OP_ANDINT64, parmA, parmB)
	if outA != parmA&parmB {
		t.Errorf("OP_ANDINT64 Expected 10, got %d %d %b", outA, outB, flags)
	}
	outA, outB, flags = ALUInt64(ALU_OP_ORINT64, parmA, parmB)
	if outA != parmA|parmB {
		t.Errorf("OP_ORINT64 Expected 10, got %d %d %b", outA, outB, flags)
	}
	outA, outB, flags = ALUInt64(ALU_OP_XORINT64, parmA, parmB)
	if outA != parmA^parmB {
		t.Errorf("OP_XORINT64 Expected 10, got %d %d %b", outA, outB, flags)
	}
	outA, outB, flags = ALUInt64(ALU_OP_SHLINT64, parmA, parmB)
	if outA != parmA<<parmB {
		t.Errorf("OP_SHLINT64 Expected 10, got %d %d %b", outA, outB, flags)
	}
	outA, outB, flags = ALUInt64(ALU_OP_SHRINT64, parmA, parmB)
	if outA != parmA>>parmB {
		t.Errorf("OP_SHRINT64 Expected 10, got %d %d %b", outA, outB, flags)
	}
}
