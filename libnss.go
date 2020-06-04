package nss

// #include <nss.h>
import "C"
import (
	"errors"

	. "github.com/izanagi1995/go-libnss/structs"
)

const (
	StatusTryagain = C.NSS_STATUS_TRYAGAIN
	StatusUnavail  = C.NSS_STATUS_UNAVAIL
	StatusNotfound = C.NSS_STATUS_NOTFOUND
	StatusSuccess  = C.NSS_STATUS_SUCCESS
)

type Status int32

var ErrNotFound error = errors.New("not found")

type Implementation interface {
	//PasswdOpen() (Status)
	//PasswdClose() (Status)
	PasswdAll() (Status, []Passwd)
	PasswdByName(name string) (Status, Passwd)
	PasswdByUid(uid uint) (Status, Passwd)

	//GroupOpen() (Status)
	//GroupClose() (Status)
	GroupAll() (Status, []Group)
	GroupByName(name string) (Status, Group)
	GroupByGid(gid uint) (Status, Group)

	//ShadowOpen() (Status)
	//ShadowClose() (Status)
	ShadowAll() (Status, []Shadow)
	ShadowByName(name string) (Status, Shadow)
}

// Prototype structure for people to embed
type LIBNSS struct{}

// setpwent
/*func (self LIBNSS) PasswdOpen() (Status) {
	return StatusUnavail
}

// endpwent
func (self LIBNSS) PasswdClose() (Status) {
	return StatusUnavail
}*/

// getpwent
func (self LIBNSS) PasswdAll() (Status, []Passwd) {
	return StatusUnavail, []Passwd{}
}

// getpwnam
func (self LIBNSS) PasswdByName(name string) (Status, Passwd) {
	return StatusUnavail, Passwd{}
}

// getpwuid
func (self LIBNSS) PasswdByUid(uid uint) (Status, Passwd) {
	return StatusUnavail, Passwd{}
}

// setgrent
/*func (self LIBNSS) GroupOpen() (Status) {
	return StatusUnavail
}

// endgrent
func (self LIBNSS) GroupClose() (Status) {
	return StatusUnavail
}*/

// endgrent
func (self LIBNSS) GroupAll() (Status, []Group) {
	return StatusUnavail, []Group{}
}

// getgrent
func (self LIBNSS) GroupByName(name string) (Status, Group) {
	return StatusUnavail, Group{}
}

// getgrnam
func (self LIBNSS) GroupByGid(gid uint) (Status, Group) {
	return StatusUnavail, Group{}
}

// getgrgid
/*func (self LIBNSS) ShadowOpen() (Status) {
	return StatusUnavail
}

// setspent
func (self LIBNSS) ShadowClose() (Status) {
	return StatusUnavail
}*/

// endspent
func (self LIBNSS) ShadowAll() (Status, []Shadow) {
	return StatusUnavail, []Shadow{}
}

// getspent
func (self LIBNSS) ShadowByName(name string) (Status, Shadow) {
	return StatusUnavail, Shadow{}
}

var implemented Implementation

// setspnam
func SetImpl(i Implementation) {
	implemented = i
}
