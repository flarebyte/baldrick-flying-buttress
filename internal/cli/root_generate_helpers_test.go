package cli

import "sync/atomic"

func setMarkdownReportWorkersForTest(workers int) func() {
	previous := int(atomic.LoadInt32(&markdownReportWorkers))
	atomic.StoreInt32(&markdownReportWorkers, int32(workers))
	return func() {
		atomic.StoreInt32(&markdownReportWorkers, int32(previous))
	}
}
