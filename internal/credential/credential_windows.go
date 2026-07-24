// Package credential provides OS-specific credential storage implementations.
// This file contains Windows-specific credential store functions.
//go:build windows
// +build windows

package credential

// #include <windows.h>
// #include <wincred.h>
// #cgo LDFLAGS: -lcrypt32 -ladvapi32
import "C"

import (
	"fmt"
	"unsafe"
)

// windowsCredentialStore stores a credential in Windows Credential Manager
func windowsCredentialStore(provider, key, value string) (bool, error) {
	targetName := fmt.Sprintf("ProxyBridge:%s:%s", provider, key)

	credential := C.CREDENTIAL{
		TargetName:     C.LPWSTR(unsafe.Pointer(C.CWideString(targetName))),
		ContentType:    C.LPWSTR(unsafe.Pointer(C.CWideString("ProxyBridge Credential"))),
		Value:          C.LPBYTE(unsafe.Pointer(C.CWideString(value))),
		ValueSize:      C.DWORD(len(value)),
		Persist:        C.DWORD(C.CRED_PERSIST_LOCAL_MACHINE),
		AttributeCount: 0,
		Attributes:     nil,
	}

	if !C.CredWriteW((*C.CREDENTIAL)(&credential), 0) {
		return false, fmt.Errorf("failed to write credential")
	}

	return true, nil
}

// windowsCredentialGet retrieves a credential from Windows Credential Manager
func windowsCredentialGet(provider, key string) (string, error) {
	targetName := fmt.Sprintf("ProxyBridge:%s:%s", provider, key)
	var credential *C.CREDENTIAL

	if !C.CredReadW(C.LPWSTR(unsafe.Pointer(C.CWideString(targetName))), C.CRED_TYPE_GENERIC, 0, &credential) {
		return "", fmt.Errorf("credential not found")
	}

	defer C.CredFree(C.LPVOID(credential))

	value := C.GoStringN((*C.char)(unsafe.Pointer(credential.Value)), C.int(credential.ValueSize))
	return value, nil
}

// windowsCredentialDelete deletes a credential from Windows Credential Manager
func windowsCredentialDelete(provider, key string) error {
	targetName := fmt.Sprintf("ProxyBridge:%s:%s", provider, key)

	if !C.CredDeleteW(C.LPWSTR(unsafe.Pointer(C.CWideString(targetName))), C.CRED_TYPE_GENERIC, 0) {
		return fmt.Errorf("failed to delete credential")
	}

	return nil
}

// CWideString converts a Go string to a Windows wide string
func CWideString(s string) *C.wchar_t {
	wide := syscall.StringToUTF16(s)
	return &wide[0]
}
