// Package credential provides OS-specific credential storage implementations.
// This file contains macOS-specific credential store functions.
//go:build darwin
// +build darwin

package credential

/*
#cgo darwin CFLAGS: -mmacosx-version-min=10.10
#cgo darwin LDFLAGS: -framework Security -framework CoreFoundation

#include <Security/Security.h>
#include <CoreFoundation/CoreFoundation.h>

static CFStringRef create_cfstring(const char *str) {
    return CFStringCreateWithCString(kCFAllocatorDefault, str, kCFStringEncodingUTF8);
}

static CFDataRef create_cfdata(const char *str, size_t len) {
    return CFDataCreate(kCFAllocatorDefault, (const UInt8 *)str, len);
}

static os_status add_generic_password(CFStringRef service, CFStringRef account, CFStringRef value) {
    CFTypeRef keys[] = { kSecClass, kSecAttrService, kSecAttrAccount, kSecValueRef };
    CFTypeRef values[] = { kSecClassGenericPassword, service, account, value };

    CFDictionaryRef dict = CFDictionaryCreate(
        kCFAllocatorDefault,
        (const void **)keys,
        (const void **)values,
        4,
        &kCFCopyStringDictionaryKeyCallBacks,
        &kCFTypeDictionaryValueCallBacks
    );

    OSStatus status = SecItemAdd(dict, NULL);
    CFRelease(dict);
    return status;
}

static os_status get_generic_password(CFStringRef service, CFStringRef account, CFDataRef *data) {
    CFTypeRef keys[] = { kSecClass, kSecAttrService, kSecAttrAccount, kSecReturnData };
    CFTypeRef values[] = { kSecClassGenericPassword, service, account, kCFBooleanTrue };

    CFDictionaryRef dict = CFDictionaryCreate(
        kCFAllocatorDefault,
        (const void **)keys,
        (const void **)values,
        4,
        &kCFCopyStringDictionaryKeyCallBacks,
        &kCFTypeDictionaryValueCallBacks
    );

    OSStatus status = SecItemCopyMatching(dict, (CFTypeRef *)data);
    CFRelease(dict);
    return status;
}

static os_status delete_generic_password(CFStringRef service, CFStringRef account) {
    CFTypeRef keys[] = { kSecClass, kSecAttrService, kSecAttrAccount };
    CFTypeRef values[] = { kSecClassGenericPassword, service, account };

    CFDictionaryRef dict = CFDictionaryCreate(
        kCFAllocatorDefault,
        (const void **)keys,
        (const void **)values,
        3,
        &kCFCopyStringDictionaryKeyCallBacks,
        &kCFTypeDictionaryValueCallBacks
    );

    OSStatus status = SecItemDelete(dict);
    CFRelease(dict);
    return status;
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func macosKeychainStore(provider, key, value string) (bool, error) {
	service := C.create_cfstring("ProxyBridge")
	defer C.CFRelease(service)

	account := C.create_cfstring(fmt.Sprintf("%s:%s", provider, key))
	defer C.CFRelease(account)

	valueData := C.create_cfdata(value, C.size_t(len(value)))
	defer C.CFRelease(valueData)

	status := C.add_generic_password(service, account, valueData)
	if status != 0 {
		return false, fmt.Errorf("failed to store keychain item: %d", status)
	}

	return true, nil
}

func macosKeychainGet(provider, key string) (string, error) {
	service := C.create_cfstring("ProxyBridge")
	defer C.CFRelease(service)

	account := C.create_cfstring(fmt.Sprintf("%s:%s", provider, key))
	defer C.CFRelease(account)

	var data C.CFDataRef
	status := C.get_generic_password(service, account, &data)
	if status != 0 {
		return "", fmt.Errorf("credential not found: %d", status)
	}
	defer C.CFRelease(data)

	dataLen := C.CFDataGetLength(data)
	dataPtr := C.CFDataGetBytePtr(data)

	value := C.GoStringN((*C.char)(unsafe.Pointer(dataPtr)), dataLen)
	return value, nil
}

func macosKeychainDelete(provider, key string) error {
	service := C.create_cfstring("ProxyBridge")
	defer C.CFRelease(service)

	account := C.create_cfstring(fmt.Sprintf("%s:%s", provider, key))
	defer C.CFRelease(account)

	status := C.delete_generic_password(service, account)
	if status != 0 {
		return fmt.Errorf("failed to delete keychain item: %d", status)
	}

	return nil
}
