package models

//XMMData contains the data that has to be delivered to the frontend for each XMM register
type XMMData struct {
	XmmID     string
	XmmValues []string
}
