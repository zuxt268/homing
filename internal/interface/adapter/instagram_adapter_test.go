package adapter

import (
	"fmt"
	"testing"

	"github.com/zuxt268/homing/internal/config"
)

func TestInstagramAdapter_GetAccount(t *testing.T) {
	add := config.Env.ADDRESS
	fmt.Println(add)
}
