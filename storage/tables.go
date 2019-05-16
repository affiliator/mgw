package storage

import "github.com/jinzhu/gorm"

func Tables() []Entity {
	return []Entity{
		Message{},
	}
}

type Entity interface{}
type Message struct {
	gorm.Model
	From      string // Store as format: Pascal Krason <p.krason@padr.io>
	To        string // make multiple possible
	Subject   string
	Body      string
	TLS       bool
	LastState string
}

//type Recipient struct {
//	gorm.Model
//	Name string
//	Host string
//	Type string // Make enum
//}
//
//type StateHistory struct {
//	gorm.Model
//
//}

/*
"`subject`, " +
"`body`, " +
"`charset`, " +
"`mail`, " +
"`spam_score`, " +
"`hash`, " +
"`content_type`, " +
"`recipient`, " +
"`has_attach`, " +
"`ip_addr`, " +
"`return_path`, */
