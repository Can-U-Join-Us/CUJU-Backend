package test

import "testing"

func TestFoo(t *testing.T) {
    // todo test code
    result := 1
    actual := 0
    if result != actual {
        t.Fatal("fail")
        // t.Error()
        // t.Fail()
        // t.Log()
    }
    t.Fatal("Success")
}
