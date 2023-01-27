package config

import "time"

var ROLEDEFINE = map[string]int{
	"ADMIN": 1,
	"USER":  2,
}

// unit(second)
const (
	RATELIMITPERMINUTE            = 100000           // ทั้งหมด 100,000 ครั้งต่อนาที
	RATELIMITPERMINUTE_IP         = 120              // 1 คน ยิงมาได้ไม่เกิน 120 ครั้งต่อนาที
	TOKENEXPIRETIME               = 60 * 60 * 3      // Token has expire time 3 hours(60 * 60 * 3)
	DURATIONCHECKINGTOKENISEXPIRE = 60 * 60 * 3      // repeat checking token is expire every 3 hours (60 * 60 * 3)
	REDISCACHEEXPIRETIME          = time.Minute * 5  // Redis cache expire time
	REDISCACHEEDATA               = time.Minute * 10 // Redis data expire time
	REDISCACHEEDATALONG           = time.Hour * 24   // Redis data expire time (long)
	REDISCACHEEDATANULL           = time.Minute * 2  // Redis data expire time (null)
)
