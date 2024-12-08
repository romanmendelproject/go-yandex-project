// Модуль управления сигналами завершения программы
package signal

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignal(t *testing.T) {
	termChan := make(chan os.Signal, 1)
	signal := Signal()

	assert.IsType(t, reflect.TypeOf(signal), reflect.TypeOf(termChan))
}
