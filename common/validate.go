package common

import (
	"regexp"
)

// const mobileRex = "/^[1](([3][0-9])|([4][0,1,4-9])|([5][0-3,5-9])|([6][2,5,6,7])|([7][0-8])|([8][0-9])|([9][0-3,5-9]))[0-9]{8}$/"
const mobileRex = "^1[3-9]\\d{9}$"

func IsValidateMobile(mobile string) bool {
	if mobile == "" {
		return false
	}
	mobileRegExp := regexp.MustCompile(mobileRex)
	//fmt.Println(mobileRegExp.MatchString(mobile))
	return mobileRegExp.MatchString(mobile)
}
