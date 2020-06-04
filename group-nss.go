package nss

//#include <grp.h>
//#include <errno.h>
//#include <stdlib.h>
/*
static char**makeCharArray(int size) {
	return calloc(sizeof(char*), size);
}

static void setArrayString(char **a, char *s, int n) {
	a[n] = s;
}

static void freeCharArray(char **a, int size) {
	int i;
	for (i = 0; i < size; i++)
		free(a[i]);
	free(a);
}
*/
import "C"

import (
	"bytes"
	"syscall"
	"unsafe"

	. "github.com/izanagi1995/go-libnss/structs"
)

var entries_group = make([]Group, 0)
var entry_index_group int

//export go_setgrent
func go_setgrent(stayopen C.int) Status {
	var status Status
	status, entries_group = implemented.GroupAll()
	entry_index_group = 0
	return status
}

//export go_endgrent
func go_endgrent() Status {
	entries_group = make([]Group, 0)
	entry_index_group = 0
	return StatusSuccess
}

//export go_getgrent_r
func go_getgrent_r(grp *C.struct_group, buf *C.char, buflen C.size_t, errnop *C.int) Status {
	if entry_index_group == len(entries_group) {
		return StatusNotfound
	}
	setCGroup(&entries_group[entry_index_group], grp, buf, buflen, errnop)
	entry_index_group++
	return StatusSuccess
}

//export go_getgrnam_r
func go_getgrnam_r(name string, grp *C.struct_group, buf *C.char, buflen C.size_t, errnop *C.int) Status {
	status, group := implemented.GroupByName(name)
	if status != StatusSuccess {
		return status
	}
	setCGroup(&group, grp, buf, buflen, errnop)
	return StatusSuccess
}

//export go_getgrgid_r
func go_getgrgid_r(gid uint, grp *C.struct_group, buf *C.char, buflen C.size_t, errnop *C.int) Status {
	status, group := implemented.GroupByGid(gid)
	if status != StatusSuccess {
		return status
	}
	setCGroup(&group, grp, buf, buflen, errnop)
	return StatusSuccess
}

// Sets the C values for libnss
func setCGroup(p *Group, grp *C.struct_group, buf *C.char, buflen C.size_t, errnop *C.int) Status {
	// TODO: Need to add length for Members....
	if len(p.Groupname)+len(p.Password)+5 > int(buflen) {
		*errnop = C.int(syscall.EAGAIN)
		return StatusTryagain
	}

	gobuf := C.GoBytes(unsafe.Pointer(buf), C.int(buflen))
	b := bytes.NewBuffer(gobuf)
	b.Reset()

	grp.gr_name = (*C.char)(unsafe.Pointer(&gobuf[b.Len()]))
	b.WriteString(p.Groupname)
	b.WriteByte(0)

	grp.gr_passwd = (*C.char)(unsafe.Pointer(&gobuf[b.Len()]))
	b.WriteString("x")
	b.WriteByte(0)

	grp.gr_gid = C.uint(p.GID)

	// ################ MAKING **C.char in GO!
	// Making a list of the members...
	// NOTE: There has to be a better way to do this.
	// I'm also making an assumption the process running this dies, freeing up the memory.

	cArray := C.malloc(C.size_t(len(p.Members)+1) * C.size_t(unsafe.Sizeof(uintptr(0))))

	// convert the C array to a Go Array so we can index it
	a := (*[1<<30 - 1]*C.char)(cArray)

	for idx, u := range p.Members {
		a[idx] = C.CString(u)
	}

	grp.gr_mem = (**C.char)(cArray)

	// ################ DONE MAKING **C.char in GO!

	return StatusSuccess
}
