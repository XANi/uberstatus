package syncthing

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

//	{
//  "version": 29,
//  "folders": [
//    {
//      "id": "Blog",
//      "label": "",
//      "filesystemType": "basic",
//      "path": "/home/user/sync/thing",
//      "type": "sendonly",
//      "paused": false,
//      "devices": [..]
//
//}
//

type STFolderConfig struct {
	Id     string `json:"id"`
	Label  string `json:"label"`
	Paused bool
}

type STConfig struct {
	Folders []STFolderConfig `json:"folders"`
}
type FolderStatusId int

const (
	StatusUnknown      FolderStatusId = iota //
	StatusIdle                               // idle
	StatusScanning                           // scanning
	StatusScanWaiting                        // scan-waiting
	StatusSyncing                            // syncing
	StatusSynPreparing                       // sync-preparing

)

//   "folder": "pdomp-znvc9",
//     "summary": {
//       "errors": 0,
//       "globalBytes": 47745486,
//       "globalDeleted": 468,
//       "globalDirectories": 807,
//       "globalFiles": 1161,
//       "globalSymlinks": 0,
//       "globalTotalItems": 2436,
//       "ignorePatterns": true,
//       "inSyncBytes": 47745486,
//       "inSyncFiles": 1161,
//       "invalid": "",
//       "localBytes": 47745486,
//       "localDeleted": 468,
//       "localDirectories": 807,
//       "localFiles": 1161,
//       "localSymlinks": 0,
//       "localTotalItems": 2436,
//       "needBytes": 0,
//       "needDeletes": 0,
//       "needDirectories": 0,
//       "needFiles": 0,
//       "needSymlinks": 0,
//       "needTotalItems": 0,
//       "pullErrors": 0,
//       "sequence": 22004,
//       "state": "idle",
//       "stateChanged": "2020-02-08T22:43:18.256904354+01:00",
//       "version": 22004
//     }

type STFolderStats struct {
	State        string    `json:"state"`
	StateChanged time.Time `json:"stateChanged"`
	OOSBytes     int       `json:"needBytes"`
	OOSItems     int       `json:"needItems"`
	Sequence     int       `json:"sequence"`
}

var apiClient = http.Client{
	Timeout: time.Second * 5,
}

var eventClient = http.Client{
	Timeout: time.Second * 120,
}

func (p *plugin) req(method, url string, body io.Reader) *http.Request {
	api, err := http.NewRequest(method, p.cfg.ServerAddr+url, body)
	if err != nil {
		p.l.Panicf("error creating http request: %s", err)
	}
	api.Header.Add("X-API-Key", p.cfg.ApiKey)
	return api
}

func (p *plugin) updateSyncthingFolders() {

	dirListR := p.req(http.MethodGet, "/rest/system/config", nil)
	res, err := apiClient.Do(dirListR)
	if err != nil {
		p.l.Errorf("error getting folder list: %s", err)
		return
	}
	var list STConfig
	err = json.NewDecoder(res.Body).Decode(&list)
	if err != nil {
		p.l.Errorf("error decoding folder list json: %s", err)
		return
	}

	idToName := make(map[string]string, len(list.Folders))
	for _, f := range list.Folders {
		if len(f.Label) > 0 {
			idToName[f.Id] = f.Label
		} else {
			idToName[f.Id] = f.Id
		}
	}
	p.Lock()
	p.folderIdToFolder = idToName
	p.Unlock()
	var alerted bool
	for id := range idToName {
		folderR := p.req(http.MethodGet, "/rest/db/status", nil)
		q := folderR.URL.Query()
		q.Add("folder", id)
		folderR.URL.RawQuery = q.Encode()
		res, err := apiClient.Do(folderR)
		if err != nil {
			// syncthing timeouts on that endpoint happen when sync is in progress, just ignore after first to not spam
			if !alerted {
				p.l.Warnf("error getting folder [%s] info: %s", id, err)
				alerted = true
			} else {
				p.l.Debugf("error getting folder [%s] info: %s", id, err)
			}
			continue
		} else {
			var folderStats STFolderStats
			err = json.NewDecoder(res.Body).Decode(&folderStats)
			if err != nil {
				p.l.Warnf("error decoding folder [%s] info: %s", id, err)
			} else {
				p.Lock()
				statusId := statusToStatusId(folderStats.State)
				if statusId == StatusUnknown {
					p.l.Warnf("unknown state on folder %s: %s", id, folderStats.State)
				}
				p.folderStatus[id] = statusId
				p.Unlock()

			}

		}
	}
}

func statusToStatusId(s string) FolderStatusId {
	switch s {
	case "idle":
		return StatusIdle
	case "scanning":
		return StatusScanning
	case "scan-waiting":
		return StatusScanWaiting
	case "syncing":
		return StatusSyncing
	case "sync-preparing":
		return StatusSynPreparing
	default:
		return StatusUnknown
	}
}
