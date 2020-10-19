package main

type fbPageEvent struct {
	Entry  []fbPageEventEntry `json:"entry"`
	Object string             `json:"object"`
}

type fbPageEventEntry struct {
	Changes []fbPageEventChange `json:"changes"`
	ID      int64               `json:"id,string"`
	Time    int64               `json:"time"`
}

type fbPageEventChange struct {
	Field string           `json:"field"`
	Value fbPageEventValue `json:"value"`
}

type fbPageEventValue struct {
	ID     int64  `json:"id,string"`
	Status string `json:"status"`
}
