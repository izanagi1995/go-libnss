package nss

//#include <pwd.h>
//#include <errno.h>
import "C"

import (
	"bytes"
	"syscall"
	"unsafe"

	. "github.com/izanagi1995/go-libnss/structs"
)

var entries_passwd = make([]Passwd, 0)
var entry_index_passwd int

//export go_setpwent
func go_setpwent() Status {
	var status Status
	status, entries_passwd = implemented.PasswdAll()
	entry_index_passwd = 0
	return status
}

//export go_endpwent
func go_endpwent() Status {
	entries_passwd = make([]Passwd, 0)
	entry_index_passwd = 0
	return StatusSuccess
}

//export go_getpwent
func go_getpwent(passwd *C.struct_passwd, buf *C.char, buflen C.size_t, errnop *C.int) Status {
	if entry_index_passwd == len(entries_passwd) {
		return StatusNotfound
	}
	setCPasswd(&entries_passwd[entry_index_passwd], passwd, buf, buflen, errnop)
	entry_index_passwd++
	return StatusSuccess
}

//export go_getpwnam
func go_getpwnam(name string, passwd *C.struct_passwd, buf *C.char, buflen C.size_t, errnop *C.int) Status {
	status, pwd := implemented.PasswdByName(name)
	if status != StatusSuccess {
		return status
	}
	setCPasswd(&pwd, passwd, buf, buflen, errnop)
	return StatusSuccess
}

//export go_getpwuid
func go_getpwuid(uid uint, passwd *C.struct_passwd, buf *C.char, buflen C.size_t, errnop *C.int) Status {
	status, pwd := implemented.PasswdByUid(uid)
	if status != StatusSuccess {
		return status
	}
	setCPasswd(&pwd, passwd, buf, buflen, errnop)
	return StatusSuccess
}

// Sets the C values for libnss
func setCPasswd(p *Passwd, passwd *C.struct_passwd, buf *C.char, buflen C.size_t, errnop *C.int) Status {
	if len(p.Username)+len(p.Password)+len(p.Gecos)+len(p.Dir)+len(p.Shell)+5 > int(buflen) {
		*errnop = C.int(syscall.EAGAIN)
		return StatusTryagain
	}

	gobuf := C.GoBytes(unsafe.Pointer(buf), C.int(buflen))
	b := bytes.NewBuffer(gobuf)
	b.Reset()

	passwd.pw_name = (*C.char)(unsafe.Pointer(&gobuf[b.Len()]))
	b.WriteString(p.Username)
	b.WriteByte(0)

	passwd.pw_passwd = (*C.char)(unsafe.Pointer(&gobuf[b.Len()]))
	b.WriteString(p.Password)
	b.WriteByte(0)

	passwd.pw_gecos = (*C.char)(unsafe.Pointer(&gobuf[b.Len()]))
	b.WriteString(p.Gecos)
	b.WriteByte(0)

	passwd.pw_dir = (*C.char)(unsafe.Pointer(&gobuf[b.Len()]))
	b.WriteString(p.Dir)
	b.WriteByte(0)

	passwd.pw_shell = (*C.char)(unsafe.Pointer(&gobuf[b.Len()]))
	b.WriteString(p.Shell)
	b.WriteByte(0)

	passwd.pw_uid = C.uint(p.UID)
	passwd.pw_gid = C.uint(p.GID)

	return StatusSuccess
}
