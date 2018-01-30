package app

type wsMessage struct {
	Status string      `json:"status"` //
	Body   interface{} `json:"body"`
	//type
	//code
}

type statement struct {
	Players map[int]playerStatement
	Wall    []string `json:"wall"`
	Wind    int      `json:"wind"` //wind of game
	Step    int      `json:"step"` //current player number
	Reserve []string `json:"reserve"`
}

type playerStatement struct {
	Hand        []string   `json:"hand"`    // hand
	Discard     []string   `json:"discard"` //
	Open        [][]string `json:"open"`    //open
	CurrentTail string     `json:"current_tail"`
	Wind        int        `json:"wind"`
}
