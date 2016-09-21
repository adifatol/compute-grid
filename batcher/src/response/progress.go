package response

import (
	"batch"
	"time"
)

type ProgressHandler struct {
	inertia      int64 /* batches/second*/
	finished     int64
	totalBatches int64
	timeStarted  int64
}

func NewProgress(batchHandler *batch.ReadHandler) *ProgressHandler {
	p := new(ProgressHandler)

	/* Aproximage nr of Batches */
	totalBatches := batchHandler.FileSize / batch.BATCH_SIZE

	p.totalBatches = totalBatches
	p.inertia = 0
	p.finished = 0
	p.timeStarted = time.Now().Unix()

	return p
}

func Progress(progressHandler *ProgressHandler) *ProgressHandler{
	sec := time.Now().Unix() - progressHandler.timeStarted
	progressHandler.finished++
	progressHandler.inertia = progressHandler.finished / sec
	return progressHandler
}
