package ds

type RoomResponse struct {
	Host     string `json:"host"` // зачем здесь хост?
	RoomName string `json:"room_name"`
	Players  int    `json:"players"`
}
