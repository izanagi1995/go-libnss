package nss

//#include <shadow.h>
//#include <errno.h>
import "C"

import (
	"bytes"
	"syscall"
	"unsafe"

	. "github.com/izanagi1995/go-libnss/structs"
)

var entries_shadow = make([]Shadow, 0)
var entry_index_shadow int

//export go_setspent
func go_setspent() Status {
	var status Status
	status, entries_shadow = implemented.ShadowAll()
	entry_index_shadow = 0
	return status
}

//export go_endspent
func go_endspent() Status {
	entries_shadow = make([]Shadow, 0)
	entry_index_shadow = 0
	return StatusSuccess
}

//export go_getspent_r
func go_getspent_r(spwd *C.struct_spwd, buf *C.char, buflen C.size_t, errnop *C.int) Status {
	if entry_index_shadow == len(entries_shadow) {
		return StatusNotfound
	}
	setCShadow(&entries_shadow[entry_index_shadow], spwd, buf, buflen, errnop)
	entry_index_shadow++
	return StatusSuccess
}

//export go_getspnam_r
func go_getspnam_r(name string, spwd *C.struct_spwd, buf *C.char, buflen C.size_t, errnop *C.int) Status {
	status, shadow := implemented.ShadowByName(name)
	if status != StatusSuccess {
		return status
	}
	setCShadow(&shadow, spwd, buf, buflen, errnop)
	return StatusSuccess
}

func setCShadow(p *Shadow, spwd *C.struct_spwd, buf *C.char, buflen C.size_t, errnop *C.int) Status {
	if len(p.Username)+len(p.Password)+7 > int(buflen) {
		*errnop = C.int(syscall.EAGAIN)
		return StatusTryagain
	}

	gobuf := C.GoBytes(unsafe.Pointer(buf), C.int(buflen))
	b := bytes.NewBuffer(gobuf)
	b.Reset()

	spwd.sp_namp = (*C.char)(unsafe.Pointer(&gobuf[b.Len()]))
	b.WriteString(p.Username)
	b.WriteByte(0)

	spwd.sp_lstchg = C.long(p.LastChange)
	spwd.sp_min = C.long(p.MinChange)
	spwd.sp_max = C.long(p.MaxChange)
	spwd.sp_warn = C.long(p.PasswordWarn)
	spwd.sp_inact = C.long(p.InactiveLockout)
	spwd.sp_expire = C.long(p.ExpirationDate)
	spwd.sp_flag = C.ulong(p.Reserved)

	spwd.sp_pwdp = (*C.char)(unsafe.Pointer(&gobuf[b.Len()]))
	b.WriteString(p.Password)
	b.WriteByte(0)

	return StatusSuccess
}
