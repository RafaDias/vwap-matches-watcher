package domain

import "errors"

var ErrCannotListenExchange = errors.New("error trying to listen exchange service")
var ErrSubscriptionNotFound = errors.New("could not find consumer for subscription")
var ErrCannotConvertPrice = errors.New("could not convert price")
var ErrCannotConvertSize = errors.New("could not convert size")
