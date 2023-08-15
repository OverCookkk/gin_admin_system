package app

import (
    "os"
    "path/filepath"
    "time"

    "gin_admin_system/internal/app/config"
    "gin_admin_system/pkg/logger"
    rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

func InitLogger() (func(), error) {
    c := config.C.Log
    logger.SetLevel(logger.Level(c.Level))
    logger.SetFormatter(c.Format)

    var file *rotatelogs.RotateLogs
    if c.Output != "" {
        switch c.Output {
        case "stdout":
            logger.SetOutput(os.Stdout)
        case "stderr":
            logger.SetOutput(os.Stderr)
        case "file":
            if name := c.OutputFile; name != "" {
                _ = os.MkdirAll(filepath.Dir(name), 0777)

                f, err := rotatelogs.New(name+".%Y-%m-%d",
                    rotatelogs.WithLinkName(name),
                    rotatelogs.WithRotationTime(time.Duration(c.RotationTime)*time.Hour),
                    rotatelogs.WithRotationCount(uint(c.RotationCount)))
                if err != nil {
                    return nil, err
                }

                logger.SetOutput(f)
                file = f
            }
        }
    }

    return func() {
        if file != nil {
            file.Close()
        }
    }, nil
}
