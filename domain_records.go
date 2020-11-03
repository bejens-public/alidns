package alidns

import (
	"github.com/libdns/libdns"
	"sync"
)

type domainRecords struct {
	records []libdns.Record
	mutex   sync.RWMutex
}

func (d *domainRecords) load(domain string) (record libdns.Record, ok bool) {

	d.mutex.RLock()
	defer d.mutex.RUnlock()

	for _, r := range d.records {
		if r.Name == domain {
			return r, true
		}
	}

	return record, false
}

func (d *domainRecords) store(record libdns.Record) {

	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.records = append(d.records, record)
}
