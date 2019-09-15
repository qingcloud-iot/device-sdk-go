package mqtt

import "testing"

func TestGetInsIdWorker(t *testing.T) {

}

func BenchmarkGetInsIdWorker(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ins := GetInsIdWorker(1)
		id, err := ins.nextid()
		if err != nil {
			b.Fail()
		}
		b.Log(id)
	}
}
