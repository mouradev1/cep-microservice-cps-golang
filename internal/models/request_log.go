package models

import "time"

type RequestLog struct {
    Ip          string    `bson:"ip" json:"ip"`
    Count       int       `bson:"count" json:"count"`
    LastRequest time.Time `bson:"lastRequest" json:"lastRequest"`
    CreatedAt   time.Time `bson:"createdAt" json:"createdAt"`
}