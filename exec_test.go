package boolean

import (
	"testing"
)

var r bool

func BenchmarkExecute(b *testing.B) {
	text := "p and q"
	data := map[string]bool{"p": true, "q": false}
	for i := 0; i < b.N; i++ {
		tree, _ := New().Parse(text)
		result, _ := tree.Execute(data)
		r = result
	}
}
