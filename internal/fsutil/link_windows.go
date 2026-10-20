package fsutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

// junctions is whether the platform has directory junctions.
const junctions = true

// symlink is os.Symlink; tests replace it to see the fallback.
var symlink = os.Symlink

// Link makes path a link to target. Windows creates symbolic links only in
// Developer Mode or for an administrator. Without that right a folder
// becomes a directory junction, which needs none; a file cannot be linked.
func Link(target, path string, dir bool) error {
	err := symlink(target, path)
	if err == nil || !errors.Is(err, windows.ERROR_PRIVILEGE_NOT_HELD) {
		return err
	}
	if dir {
		return junction(target, path)
	}
	return fmt.Errorf("cannot link %s: Windows links files only in Developer Mode (Settings > System > For developers). Turn it on, or keep this item local", path)
}

// junction creates path as a directory junction (a mount point reparse
// point) to the absolute folder target.
func junction(target, path string) error {
	target, err := filepath.Abs(target)
	if err != nil {
		return err
	}
	sub, err := windows.UTF16FromString(`\??\` + target)
	if err != nil {
		return err
	}
	prt, err := windows.UTF16FromString(target)
	if err != nil {
		return err
	}
	// REPARSE_DATA_BUFFER with a MountPointReparseBuffer: tag, data length,
	// reserved, then offset and length (in bytes, without the NUL) of the
	// substitute and print names, which follow with their NULs.
	names := append(sub, prt...)
	data := make([]byte, 16+2*len(names))
	le := func(off int, v uint32, size int) {
		for i := range size {
			data[off+i] = byte(v >> (8 * i))
		}
	}
	le(0, windows.IO_REPARSE_TAG_MOUNT_POINT, 4)
	le(4, uint32(len(data)-8), 2)
	le(8, 0, 2)
	le(10, uint32(2*(len(sub)-1)), 2)
	le(12, uint32(2*len(sub)), 2)
	le(14, uint32(2*(len(prt)-1)), 2)
	for i, c := range names {
		le(16+2*i, uint32(c), 2)
	}
	if err := os.Mkdir(path, 0o700); err != nil {
		return err
	}
	p, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	h, err := windows.CreateFile(p, windows.GENERIC_WRITE, 0, nil, windows.OPEN_EXISTING,
		windows.FILE_FLAG_OPEN_REPARSE_POINT|windows.FILE_FLAG_BACKUP_SEMANTICS, 0)
	if err != nil {
		os.Remove(path)
		return err
	}
	var n uint32
	err = windows.DeviceIoControl(h, windows.FSCTL_SET_REPARSE_POINT, &data[0], uint32(len(data)), nil, 0, &n, nil)
	_ = windows.CloseHandle(h)
	if err != nil {
		os.Remove(path)
		return fmt.Errorf("could not create a junction at %s: %w", path, err)
	}
	return nil
}
