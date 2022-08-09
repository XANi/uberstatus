package syncthing

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

//"id": 2038,
//  "globalID": 4589,
//  "time": "2020-02-09T00:07:05.243068465+01:00",
//  "type": "FolderCompletion",
//  "data": {
//    "completion": 100,
//    "device": "zzzzzzz-yyyyyyy-zzzzzzz-yyyyyyy-zzzzzzz-yyyyyyy-zzzzzzz-yyyyyyy",
//    "folder": "t9jyl-ecmyd",
//    "globalBytes": 639373,
//    "needBytes": 0,
//    "needDeletes": 0,
//    "needItems": 0
//  }
//}

type STEvent struct {
	Id   int       `json:"id"`
	Time time.Time `json:"time"`
	Type string    `json:"type"`
	Data json.RawMessage
}

type STEvents []STEvent

//    "duration": 0.000005238,
//    "folder": "Blog",
//    "from": "scan-waiting",
//    "to": "scanning"

type STEventStateChanged struct {
	Duration float32 `json:"duration"`
	FolderId string  `json:"folder"`
	From     string  `json:"from"`
	To       string  `json:"to"`
}

//   "completion": 31.28736597109183,
//   "device": "zzzzzzz-yyyyyyy-zzzzzzz-yyyyyyy-zzzzzzz-yyyyyyy-zzzzzzz-yyyyyyy",
//   "folder": "pdomp-znvc9",
//   "globalBytes": 152603086,
//   "needBytes": 104857600,
//   "needDeletes": 1,
//   "needItems": 1

type STEventFolderCompletion struct {
	Completion float32 `json:"completion"`
	FolderId   string  `json:"folder"`
}

type STEventFolderSummary struct {
	Folder  string        `json:"folder"`
	Summary STFolderStats `json:"summary"`
}

func (p *plugin) updateSynctingEvents() {
	eventId := 1
	for {
		time.Sleep(20 * time.Second) // must be at start or else continue blocks will skip it
		eventR := p.req(http.MethodGet, "/rest/events", nil)
		q := eventR.URL.Query()
		q.Add("since", strconv.Itoa(eventId))
		eventR.URL.RawQuery = q.Encode()
		res, err := eventClient.Do(eventR)
		if err != nil {
			continue
			// TODO ignore timeout, log every other error
		}
		var events STEvents
		err = json.NewDecoder(res.Body).Decode(&events)
		if err != nil {
			log.Warningf("Error decoding event: %s", err)
			continue
		}
		for _, ev := range events {
			if ev.Id < 1 {
				continue
			}
			eventId = ev.Id
			switch ev.Type {
			case "StateChanged":
				p.handleEventStateChanged(ev.Data)
			case "FolderCompletion":
				p.handleEventFolderCompletion(ev.Data)
			case "FolderSummary":
				p.handleEventFolderSummary(ev.Data)
			default:
				//log.Debugf("event: %s\n %s",ev.Type,string(ev.Data))
			}
		}
	}
}

func (p *plugin) handleEventStateChanged(m json.RawMessage) {
	var event STEventStateChanged
	err := json.Unmarshal(m, &event)
	if err != nil {
		log.Warningf("could not unmarshal event: %s", err)
		return
	}
	state := statusToStatusId(event.To)
	p.Lock()
	p.folderStatus[event.FolderId] = state
	p.Unlock()
}
func (p *plugin) handleEventFolderCompletion(m json.RawMessage) {
	var event STEventFolderCompletion
	err := json.Unmarshal(m, &event)
	if err != nil {
		log.Warningf("could not unmarshal event: %s", err)
		return
	}

	p.Lock()
	p.folderCompletion[event.FolderId] = event.Completion
	p.folderStatus[event.FolderId] = StatusSyncing
	p.Unlock()
}

func (p *plugin) handleEventFolderSummary(m json.RawMessage) {
	var event STEventFolderSummary
	err := json.Unmarshal(m, &event)
	if err != nil {
		log.Warningf("could not unmarshal event: %s", err)
		return
	}
	stateId := statusToStatusId(event.Summary.State)
	p.Lock()
	p.folderStatus[event.Folder] = stateId
	p.Unlock()
}
