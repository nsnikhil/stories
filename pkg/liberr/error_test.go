package liberr_test

import (
	"errors"
	"fmt"
	"github.com/nsnikhil/stories/pkg/liberr"
	"testing"
)

//TODO: COMPLETE THE TESTS
func TestCreateError(t *testing.T) {
	err := errors.New("record not found")

	ae := liberr.WithArgs(liberr.Severity("error"), liberr.Operation("db.query"), err)

	fmt.Println(ae)

	ke := liberr.WithKind("resource not found", ae)

	fmt.Println(ke)
}
