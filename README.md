# logger
## How to use
```go
import "github.com/EricChiou/logger"

logger.Init("log/")

logger.Trace.Println("some thing failed")
logger.Info.Println("some thing failed")
logger.Warn.Println("some thing failed")
logger.Error.Println("some thing failed")
```