package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
	"warranty-api/services"
)

// You can also use the simplified format "<limit>-<period>"", with the given
// periods:
//
// * "S": second
// * "M": minute
// * "H": hour
// * "D": day
//
// Examples:
//
// * 5 reqs/second: "5-S"
// * 10 reqs/minute: "10-M"
// * 1000 reqs/hour: "1000-H"
// * 2000 reqs/day: "2000-D"
//
func RateLimitMiddleware(format string) gin.HandlerFunc {
	rate, err := limiter.NewRateFromFormatted(format)
	if err != nil {
		panic(err)
	}

	rateLimitRedisStore, err := sredis.NewStoreWithOptions(services.RedisConnect(), limiter.StoreOptions{
		Prefix:   "limiter_warranty",
		MaxRetry: 3,
	})
	middleware := mgin.NewMiddleware(limiter.New(rateLimitRedisStore, rate))
	return middleware
}
