package database

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type CustomLogger struct {
    logger.Interface
}

func (c *CustomLogger) LogMode(level logger.LogLevel) logger.Interface {
    return &CustomLogger{c.Interface.LogMode(level)}
}

func (c *CustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
    if err != nil && err == gorm.ErrRecordNotFound {
        // Suppress record not found message
        return
    }
    // Call the default logger for other messages
    // c.Interface.Trace(ctx, begin, fc, err)
}
