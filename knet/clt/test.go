package clt

import (
	"encoding/json"
	"fmt"
	"game.pb"
	"gamelord.pb"
	"log"
	"match.pb"
	"watch.pb"
)

func (self *Client) WatchTest() {

	self.readyNtf()
	self.startNtf()
	self.boutEndNtf()
	self.confirmBoutResultNtf()
	self.playerSelectSeatAck()
	self.teammateSelectSeatNtf()
	self.initCardNtf()
	self.callScoreAck()
	self.callScoreNtf()
	self.moveCardNtf()
	self.callScoreResultNtf()
	self.roundTakeoutStartNtf()
	self.takeoutCardAck()
	self.takeoutCardNtf()
	self.lordResultNtf()
}

//测试通用消息

func (self *Client) readyNtf() {
	lord := &struct {
		ID      string          `json:"id"`
		Content *watch.WatchNtf `json:"content"`
	}{
		Content: &watch.WatchNtf{
			MatchId: 1234,
			TableId: 1,
			Seat:    2,
			GameId:  1007,
		},
	}

	lord.ID = "GameReadyNtf"
	lord.Content.GameReadyNtf = &game.GameReadyNtf{
		Type:   match.MatchType_Team,
		Player: make([]*game.SeatPlayer, 0),
	}

	for i := 0; i < 4; i++ {
		p := &game.SeatPlayer{
			Pos:        (match.SeatDirection)(i),
			Teamname:   "AKB48",
			Playername: fmt.Sprintf("player_auto_%d", i),
			Playerid:   (int32)(i),
		}
		lord.Content.GameReadyNtf.Player = append(lord.Content.GameReadyNtf.Player, p)
	}

	b, err := json.Marshal(lord)
	if err != nil {
		log.Println(err)
		return
	}
	self.Send(b)
}

func (self *Client) startNtf() {
	lord := &struct {
		ID      string          `json:"id"`
		Content *watch.WatchNtf `json:"content"`
	}{
		Content: &watch.WatchNtf{
			MatchId: 1234,
			TableId: 1,
			Seat:    2,
			GameId:  1007,
		},
	}

	lord.ID = "GameStartNtf"
	lord.Content.GameStartNtf = &game.GameStartNtf{
		Hand: 1234,
	}
	b, err := json.Marshal(lord)
	if err != nil {
		log.Println(err)
		return
	}
	self.Send(b)
}

func (self *Client) boutEndNtf() {
	lord := &struct {
		ID      string          `json:"id"`
		Content *watch.WatchNtf `json:"content"`
	}{
		Content: &watch.WatchNtf{
			MatchId: 1234,
			TableId: 1,
			Seat:    2,
			GameId:  1007,
		},
	}

	lord.ID = "BoutEndNtf"
	lord.Content.BoutEndNtf = &game.BoutEndNtf{
		Player: make([]*game.PlayerScore, 0),
	}

	for i := 0; i < 4; i++ {
		p := &game.PlayerScore{
			Pos:   (match.SeatDirection)(i),
			Score: (int32)(100 * i),
		}
		lord.Content.BoutEndNtf.Player = append(lord.Content.BoutEndNtf.Player, p)
	}

	b, err := json.Marshal(lord)
	if err != nil {
		log.Println(err)
		return
	}
	self.Send(b)
}

func (self *Client) confirmBoutResultNtf() {
	lord := &struct {
		ID      string          `json:"id"`
		Content *watch.WatchNtf `json:"content"`
	}{
		Content: &watch.WatchNtf{
			MatchId: 1234,
			TableId: 1,
			Seat:    2,
			GameId:  1007,
		},
	}

	lord.ID = "ConfirmBoutResultNtf"
	lord.Content.ConfirmBoutResultNtf = &game.ConfirmBoutResultNtf{
		Pos: match.SeatDirection_NorthSeat,
	}

	b, err := json.Marshal(lord)
	if err != nil {
		log.Println(err)
		return
	}
	self.Send(b)
}

func (self *Client) playerSelectSeatAck() {
	lord := &struct {
		ID      string          `json:"id"`
		Content *watch.WatchNtf `json:"content"`
	}{
		Content: &watch.WatchNtf{
			MatchId: 1234,
			TableId: 1,
			Seat:    2,
			GameId:  1007,
		},
	}

	lord.ID = "PlayerSelectSeatAck"
	lord.Content.PlayerSelectSeatAck = &game.PlayerSelectSeatAck{
		PlayerNum: 1,
		Result:    match.SelectSeatResult_INVALID_SEAT,
	}

	b, err := json.Marshal(lord)
	if err != nil {
		log.Println(err)
		return
	}
	self.Send(b)
}

func (self *Client) teammateSelectSeatNtf() {
	lord := &struct {
		ID      string          `json:"id"`
		Content *watch.WatchNtf `json:"content"`
	}{
		Content: &watch.WatchNtf{
			MatchId: 1234,
			TableId: 1,
			Seat:    2,
			GameId:  1007,
		},
	}

	lord.ID = "TeammateSelectSeatNtf"
	lord.Content.TeammateSelectSeatNtf = &game.TeammateSelectSeatNtf{
		Matchid:   1234,
		Tableid:   1111,
		Roomtype:  match.RoomType_ClosedRoom,
		PlayerNum: 4,
		Seat:      match.SeatDirection_NorthSeat,
	}

	b, err := json.Marshal(lord)
	if err != nil {
		log.Println(err)
		return
	}
	self.Send(b)
}

//测试游戏消息

func (self *Client) initCardNtf() {
	ntf := &struct {
		ID      string          `json:"id"`
		Content *watch.WatchNtf `json:"content"`
	}{
		Content: &watch.WatchNtf{
			MatchId: 1234,
			TableId: 1,
			Seat:    2,
			GameId:  1007,
		},
	}

	ntf.ID = "InitCardNtf"
	ntf.Content.WatchLordAck = &watch.WatchLordAck{
		InitCardNtf: &lord.InitCardNtf{
			PosFirts:         match.SeatDirection_EastSeat,
			CountBottomcards: 3,
			Next: &lord.NextAction{
				Pos:     match.SeatDirection_NorthSeat,
				Action:  lord.Action_ACTION_CALL,
				Timeout: 10,
				Step:    1,
			},
		},
	}

	//init cards

	ntf.Content.WatchLordAck.InitCardNtf.Players = make([]*lord.PlayerCards, 0)

	for i := 0; i < 3; i++ {
		p := &lord.PlayerCards{
			Pos:   (match.SeatDirection)(i),
			Count: 16,
			Cards: make([]*lord.Card, 0),
		}

		for j := 0; j < (int)(p.Count); j++ {
			p.Cards = append(p.Cards, &lord.Card{
				Id: 53,
			})
		}

		ntf.Content.WatchLordAck.InitCardNtf.Players = append(ntf.Content.WatchLordAck.InitCardNtf.Players, p)
	}

	b, err := json.Marshal(ntf)
	if err != nil {
		log.Println(err)
		return
	}

	self.Send(b)
}

func (self *Client) callScoreAck() {

	ntf := &struct {
		ID      string          `json:"id"`
		Content *watch.WatchNtf `json:"content"`
	}{
		Content: &watch.WatchNtf{
			MatchId: 1234,
			TableId: 1,
			Seat:    2,
			GameId:  1007,
		},
	}

	ntf.ID = "CallScoreAck"
	ntf.Content.WatchLordAck = &watch.WatchLordAck{
		CallScoreAck: &lord.CallScoreAck{
			Result: lord.CallScoreAck_RESULT_INVALID_USER,
		},
	}

	b, err := json.Marshal(ntf)
	if err != nil {
		log.Println(err)
		return
	}

	self.Send(b)
}

func (self *Client) callScoreNtf() {

	ntf := &struct {
		ID      string          `json:"id"`
		Content *watch.WatchNtf `json:"content"`
	}{
		Content: &watch.WatchNtf{
			MatchId: 1234,
			TableId: 1,
			Seat:    2,
			GameId:  1007,
		},
	}

	ntf.ID = "CallScoreNtf"
	ntf.Content.WatchLordAck = &watch.WatchLordAck{
		CallScoreNtf: &lord.CallScoreNtf{
			Pos:   match.SeatDirection_SouthSeat,
			Score: 3,
			Next: &lord.NextAction{
				Pos:     match.SeatDirection_NorthSeat,
				Action:  lord.Action_ACTION_TAKEOUT,
				Timeout: 10,
				Step:    1,
			},
		},
	}

	b, err := json.Marshal(ntf)
	if err != nil {
		log.Println(err)
		return
	}

	self.Send(b)
}

func (self *Client) moveCardNtf() {
	ntf := &struct {
		ID      string          `json:"id"`
		Content *watch.WatchNtf `json:"content"`
	}{
		Content: &watch.WatchNtf{
			MatchId: 1234,
			TableId: 1,
			Seat:    2,
			GameId:  1007,
		},
	}

	ntf.ID = "MoveCardNtf"
	ntf.Content.WatchLordAck = &watch.WatchLordAck{
		MoveCardNtf: &lord.MoveCardNtf{
			From:  match.SeatDirection_SouthSeat,
			To:    match.SeatDirection_NorthSeat,
			Cards: make([]*lord.Card, 0),
		},
	}

	for i := 0; i < 16; i++ {
		ntf.Content.WatchLordAck.MoveCardNtf.Cards =
			append(ntf.Content.WatchLordAck.MoveCardNtf.Cards, &lord.Card{
				Id: 51,
			})
	}

	b, err := json.Marshal(ntf)
	if err != nil {
		log.Println(err)
		return
	}

	self.Send(b)
}

func (self *Client) callScoreResultNtf() {
	ntf := &struct {
		ID      string          `json:"id"`
		Content *watch.WatchNtf `json:"content"`
	}{
		Content: &watch.WatchNtf{
			MatchId: 1234,
			TableId: 1,
			Seat:    2,
			GameId:  1007,
		},
	}

	moveCardNtf := &lord.MoveCardNtf{
		From:  match.SeatDirection_SouthSeat,
		To:    match.SeatDirection_NorthSeat,
		Cards: make([]*lord.Card, 0),
	}

	for i := 0; i < 16; i++ {
		moveCardNtf.Cards =
			append(moveCardNtf.Cards, &lord.Card{
				Id: 51,
			})
	}

	ntf.ID = "CallScoreResultNtf"
	ntf.Content.WatchLordAck = &watch.WatchLordAck{
		CallScoreResultNtf: &lord.CallScoreResultNtf{
			Lord:        match.SeatDirection_SouthSeat,
			CallScore:   3,
			BottomCards: make([]*lord.Card, 0),
			Movecard:    moveCardNtf,
		},
	}

	b, err := json.Marshal(ntf)
	if err != nil {
		log.Println(err)
		return
	}

	self.Send(b)
}

func (self *Client) roundTakeoutStartNtf() {
	ntf := &struct {
		ID      string          `json:"id"`
		Content *watch.WatchNtf `json:"content"`
	}{
		Content: &watch.WatchNtf{
			MatchId: 1234,
			TableId: 1,
			Seat:    2,
			GameId:  1007,
		},
	}

	moveCardNtf := &lord.MoveCardNtf{
		From:  match.SeatDirection_SouthSeat,
		To:    match.SeatDirection_NorthSeat,
		Cards: make([]*lord.Card, 0),
	}

	for i := 0; i < 16; i++ {
		moveCardNtf.Cards =
			append(moveCardNtf.Cards, &lord.Card{
				Id: 51,
			})
	}

	ntf.ID = "RoundTakeoutStartNtf"
	ntf.Content.WatchLordAck = &watch.WatchLordAck{
		RoundTakeoutStartNtf: &lord.RoundTakeoutStartNtf{
			Next: &lord.NextAction{
				Pos:     match.SeatDirection_NorthSeat,
				Action:  lord.Action_ACTION_TAKEOUT,
				Timeout: 10,
				Step:    1,
			},
		},
	}

	b, err := json.Marshal(ntf)
	if err != nil {
		log.Println(err)
		return
	}

	self.Send(b)
}

func (self *Client) takeoutCardAck() {
	ntf := &struct {
		ID      string          `json:"id"`
		Content *watch.WatchNtf `json:"content"`
	}{
		Content: &watch.WatchNtf{
			MatchId: 1234,
			TableId: 1,
			Seat:    2,
			GameId:  1007,
		},
	}

	ntf.ID = "TakeoutCardAck"
	ntf.Content.WatchLordAck = &watch.WatchLordAck{
		TakeoutCardAck: &lord.TakeoutCardAck{
			Result: lord.TakeoutCardAck_RESLUT_INVALID_CARD,
		},
	}

	b, err := json.Marshal(ntf)
	if err != nil {
		log.Println(err)
		return
	}

	self.Send(b)
}

func (self *Client) takeoutCardNtf() {
	ntf := &struct {
		ID      string          `json:"id"`
		Content *watch.WatchNtf `json:"content"`
	}{
		Content: &watch.WatchNtf{
			MatchId: 1234,
			TableId: 1,
			Seat:    2,
			GameId:  1007,
		},
	}

	ntf.ID = "TakeoutCardNtf"
	ntf.Content.WatchLordAck = &watch.WatchLordAck{
		TakeoutCardNtf: &lord.TakeoutCardNtf{
			Pos:            match.SeatDirection_EastSeat,
			CountLeftcards: 13,
			Desc: &lord.PlayerTakeoutCardsDesc{
				Cards:     []*lord.Card{{Id: 51}, {Id: 51}, {Id: 51}},
				Type:      lord.CardsType_CARDS_TYPE_AA,
				Mainpoint: lord.CardPoint_CARD_POINT_3,
				Subcards:  []*lord.Card{{Id: 51}, {Id: 51}, {Id: 51}},
				Mainlen:   2,
			},
			Next: &lord.NextAction{
				Action:  lord.Action_ACTION_TAKEOUT,
				Pos:     match.SeatDirection_NorthSeat,
				Step:    1,
				Timeout: 10,
			},
		},
	}

	b, err := json.Marshal(ntf)
	if err != nil {
		log.Println(err)
		return
	}

	self.Send(b)
}

func (self *Client) lordResultNtf() {

	ntf := &struct {
		ID      string          `json:"id"`
		Content *watch.WatchNtf `json:"content"`
	}{
		Content: &watch.WatchNtf{
			MatchId: 1234,
			TableId: 1,
			Seat:    2,
			GameId:  1007,
		},
	}

	ntf.ID = "LordResultNtf"
	ntf.Content.WatchLordAck = &watch.WatchLordAck{
		LordResultNtf: &lord.LordResultNtf{
			Hand:   10,
			Result: lord.LordResultNtf_RESULT_LORD_WIN,
			PlayerCards: []*lord.PlayerCards{
				{Pos: match.SeatDirection_EastSeat, Count: 10, Cards: []*lord.Card{{Id: 51}, {Id: 51}, {Id: 51}}},
				{Pos: match.SeatDirection_SouthSeat, Count: 10, Cards: []*lord.Card{{Id: 51}, {Id: 51}, {Id: 51}}},
				{Pos: match.SeatDirection_NorthSeat, Count: 10, Cards: []*lord.Card{{Id: 51}, {Id: 51}, {Id: 51}}},
			},
			Scores: []*game.PlayerScore{
				{match.SeatDirection_EastSeat, 10},
				{match.SeatDirection_SouthSeat, 10},
				{match.SeatDirection_NorthSeat, 10},
			},
			State: &lord.LordResultNtf_GameState{
				Spring: true,
				Bomb:   2,
			},
			TotalScores: []*game.PlayerScore{
				{match.SeatDirection_EastSeat, 100},
				{match.SeatDirection_SouthSeat, 100},
				{match.SeatDirection_NorthSeat, 100},
			},
			Timeout: 10,
		},
	}

	b, err := json.Marshal(ntf)
	if err != nil {
		log.Println(err)
		return
	}

	self.Send(b)
}
