package model

type Request struct {
    Action  string `json:"action"`
    Filter  string `json:"filter"`
    StartAt string `json:"startAt"`
    EndAt   string `json:"endAt"`
}
