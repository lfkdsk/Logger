package logger

import (
	"testing"
)

func TestExample(t *testing.T) {

	log, _ := New("", "fuck", true)
	log.SetMaxLevel(R_level)
	log.WTF("fuck: %d , %s", 112, "lfkds")
	log.D("do you love Me?")
	log.Close()


}


