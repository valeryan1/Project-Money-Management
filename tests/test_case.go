package tests

import (
	"github.com/goravel/framework/testing"

	"spendid/bootstrap"
)

func init() {
	bootstrap.Boot()
}

type TestCase struct {
	testing.TestCase
}
