package Configuration

import (
	"os"
	"testing"
)

func TestConfigObject_Export(t *testing.T) {
	s, err := MockConfig()
	if err != nil {
		t.Error(err)
	}
	c, err := LoadConfiguration(s)
	if err != nil {
		t.Error(err)
	}
	b := c.Save()
	f, err := os.Create("/tmp/sysgen.json")
	if err != nil {
		t.Error(err)
	}
	n, err := f.WriteString(b)
	if err != nil {
		t.Error(err)
	}
	if n != len(b) {
		t.Error("Failed to write all bytes")
	}
	f.Close()
}
