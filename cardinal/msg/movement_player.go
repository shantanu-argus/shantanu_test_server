package msg

type MovementPlayerMsg struct {
	TargetNickname string  `json:"target"`
	Velocity       int     `json:"velocity"`
	Direction      int     `json:"direction"`
	LocationX      float64 `json:"locationX"`
	LocationY      float64 `json:"locationY"`
}

type MovementPlayerMsgReply struct {
	LocationX float64 `json:"locationX"`
	LocationY float64 `json:"locationY"`
	IsValid   bool    `json:"isValid"`
}
