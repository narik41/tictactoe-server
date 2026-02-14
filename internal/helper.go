package internal

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

func UUID(prefix string) string {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		log.Printf("UUID generation error: %v", err)
		return ""
	}
	return fmt.Sprintf("%s-%s", prefix, newUUID.String())
}

func GetNPTToUtcInMillisecond() int64 {
	nptLocation, err := time.LoadLocation("Asia/Kathmandu")
	if err != nil {
		return time.Now().UTC().UnixMilli()
	}
	nptTime := time.Now().In(nptLocation)
	return nptTime.UTC().UnixMilli()
}
