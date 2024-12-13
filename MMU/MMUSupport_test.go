package MMU

import (
	"reflect"
	"testing"
)

func TestMMU_Read(t *testing.T) {
	type fields struct {
		PageTable      []PageTableEntry
		TLB            []TLBEntry
		PhysicalMem    []byte
		Disk           [][]byte
		FreePages      []int
		FreeDiskSlots  []int
		TLBHitCount    int
		TLBMissCount   int
		PageFaultCount int
		SwapCount      int
		CurrentMode    int
	}
	type args struct {
		virtualAddr int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu := &MMU{
				PageTable:      tt.fields.PageTable,
				TLB:            tt.fields.TLB,
				PhysicalMem:    tt.fields.PhysicalMem,
				Disk:           tt.fields.Disk,
				FreePages:      tt.fields.FreePages,
				FreeDiskSlots:  tt.fields.FreeDiskSlots,
				TLBHitCount:    tt.fields.TLBHitCount,
				TLBMissCount:   tt.fields.TLBMissCount,
				PageFaultCount: tt.fields.PageFaultCount,
				SwapCount:      tt.fields.SwapCount,
				CurrentMode:    tt.fields.CurrentMode,
			}
			got, err := mmu.Read(tt.args.virtualAddr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Read() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMMU_Translate(t *testing.T) {
	type fields struct {
		PageTable      []PageTableEntry
		TLB            []TLBEntry
		PhysicalMem    []byte
		Disk           [][]byte
		FreePages      []int
		FreeDiskSlots  []int
		TLBHitCount    int
		TLBMissCount   int
		PageFaultCount int
		SwapCount      int
		CurrentMode    int
	}
	type args struct {
		virtualAddr int
		accessType  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu := &MMU{
				PageTable:      tt.fields.PageTable,
				TLB:            tt.fields.TLB,
				PhysicalMem:    tt.fields.PhysicalMem,
				Disk:           tt.fields.Disk,
				FreePages:      tt.fields.FreePages,
				FreeDiskSlots:  tt.fields.FreeDiskSlots,
				TLBHitCount:    tt.fields.TLBHitCount,
				TLBMissCount:   tt.fields.TLBMissCount,
				PageFaultCount: tt.fields.PageFaultCount,
				SwapCount:      tt.fields.SwapCount,
				CurrentMode:    tt.fields.CurrentMode,
			}
			got, err := mmu.Translate(tt.args.virtualAddr, tt.args.accessType)
			if (err != nil) != tt.wantErr {
				t.Errorf("Translate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Translate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMMU_Write(t *testing.T) {
	type fields struct {
		PageTable      []PageTableEntry
		TLB            []TLBEntry
		PhysicalMem    []byte
		Disk           [][]byte
		FreePages      []int
		FreeDiskSlots  []int
		TLBHitCount    int
		TLBMissCount   int
		PageFaultCount int
		SwapCount      int
		CurrentMode    int
	}
	type args struct {
		virtualAddr int
		value       byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu := &MMU{
				PageTable:      tt.fields.PageTable,
				TLB:            tt.fields.TLB,
				PhysicalMem:    tt.fields.PhysicalMem,
				Disk:           tt.fields.Disk,
				FreePages:      tt.fields.FreePages,
				FreeDiskSlots:  tt.fields.FreeDiskSlots,
				TLBHitCount:    tt.fields.TLBHitCount,
				TLBMissCount:   tt.fields.TLBMissCount,
				PageFaultCount: tt.fields.PageFaultCount,
				SwapCount:      tt.fields.SwapCount,
				CurrentMode:    tt.fields.CurrentMode,
			}
			if err := mmu.Write(tt.args.virtualAddr, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMMU_evictPage(t *testing.T) {
	type fields struct {
		PageTable      []PageTableEntry
		TLB            []TLBEntry
		PhysicalMem    []byte
		Disk           [][]byte
		FreePages      []int
		FreeDiskSlots  []int
		TLBHitCount    int
		TLBMissCount   int
		PageFaultCount int
		SwapCount      int
		CurrentMode    int
	}
	tests := []struct {
		name    string
		fields  fields
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu := &MMU{
				PageTable:      tt.fields.PageTable,
				TLB:            tt.fields.TLB,
				PhysicalMem:    tt.fields.PhysicalMem,
				Disk:           tt.fields.Disk,
				FreePages:      tt.fields.FreePages,
				FreeDiskSlots:  tt.fields.FreeDiskSlots,
				TLBHitCount:    tt.fields.TLBHitCount,
				TLBMissCount:   tt.fields.TLBMissCount,
				PageFaultCount: tt.fields.PageFaultCount,
				SwapCount:      tt.fields.SwapCount,
				CurrentMode:    tt.fields.CurrentMode,
			}
			got, err := mmu.evictPage()
			if (err != nil) != tt.wantErr {
				t.Errorf("evictPage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("evictPage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMMU_findInTLB(t *testing.T) {
	type fields struct {
		PageTable      []PageTableEntry
		TLB            []TLBEntry
		PhysicalMem    []byte
		Disk           [][]byte
		FreePages      []int
		FreeDiskSlots  []int
		TLBHitCount    int
		TLBMissCount   int
		PageFaultCount int
		SwapCount      int
		CurrentMode    int
	}
	type args struct {
		virtualPage int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu := &MMU{
				PageTable:      tt.fields.PageTable,
				TLB:            tt.fields.TLB,
				PhysicalMem:    tt.fields.PhysicalMem,
				Disk:           tt.fields.Disk,
				FreePages:      tt.fields.FreePages,
				FreeDiskSlots:  tt.fields.FreeDiskSlots,
				TLBHitCount:    tt.fields.TLBHitCount,
				TLBMissCount:   tt.fields.TLBMissCount,
				PageFaultCount: tt.fields.PageFaultCount,
				SwapCount:      tt.fields.SwapCount,
				CurrentMode:    tt.fields.CurrentMode,
			}
			got, got1 := mmu.findInTLB(tt.args.virtualPage)
			if got != tt.want {
				t.Errorf("findInTLB() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("findInTLB() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMMU_handlePageFault(t *testing.T) {
	type fields struct {
		PageTable      []PageTableEntry
		TLB            []TLBEntry
		PhysicalMem    []byte
		Disk           [][]byte
		FreePages      []int
		FreeDiskSlots  []int
		TLBHitCount    int
		TLBMissCount   int
		PageFaultCount int
		SwapCount      int
		CurrentMode    int
	}
	type args struct {
		virtualPage int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu := &MMU{
				PageTable:      tt.fields.PageTable,
				TLB:            tt.fields.TLB,
				PhysicalMem:    tt.fields.PhysicalMem,
				Disk:           tt.fields.Disk,
				FreePages:      tt.fields.FreePages,
				FreeDiskSlots:  tt.fields.FreeDiskSlots,
				TLBHitCount:    tt.fields.TLBHitCount,
				TLBMissCount:   tt.fields.TLBMissCount,
				PageFaultCount: tt.fields.PageFaultCount,
				SwapCount:      tt.fields.SwapCount,
				CurrentMode:    tt.fields.CurrentMode,
			}
			if err := mmu.handlePageFault(tt.args.virtualPage); (err != nil) != tt.wantErr {
				t.Errorf("handlePageFault() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMMU_updateTLB(t *testing.T) {
	type fields struct {
		PageTable      []PageTableEntry
		TLB            []TLBEntry
		PhysicalMem    []byte
		Disk           [][]byte
		FreePages      []int
		FreeDiskSlots  []int
		TLBHitCount    int
		TLBMissCount   int
		PageFaultCount int
		SwapCount      int
		CurrentMode    int
	}
	type args struct {
		virtualPage  int
		physicalPage int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu := &MMU{
				PageTable:      tt.fields.PageTable,
				TLB:            tt.fields.TLB,
				PhysicalMem:    tt.fields.PhysicalMem,
				Disk:           tt.fields.Disk,
				FreePages:      tt.fields.FreePages,
				FreeDiskSlots:  tt.fields.FreeDiskSlots,
				TLBHitCount:    tt.fields.TLBHitCount,
				TLBMissCount:   tt.fields.TLBMissCount,
				PageFaultCount: tt.fields.PageFaultCount,
				SwapCount:      tt.fields.SwapCount,
				CurrentMode:    tt.fields.CurrentMode,
			}
			mmu.updateTLB(tt.args.virtualPage, tt.args.physicalPage)
		})
	}
}

func TestNewMMU(t *testing.T) {
	type args struct {
		cnf MMUConfig
	}
	tests := []struct {
		name string
		args args
		want *MMU
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMMU(tt.args.cnf); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMMU() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMMUConfig(t *testing.T) {
	tests := []struct {
		name string
		want *MMUConfig
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMMUConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMMUConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
