// +build windows

package winapi

import (
	"fmt"
	"time"

	"golang.org/x/sys/windows/registry"

	so "github.com/iamacarpet/go-win64api/shared"
)

func InstalledSoftwareList() ([]so.Software, error) {
	sw64, err := GetSoftwareList(`SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`, "X64")
	if err != nil {
		return nil, err
	}
	sw32, err := GetSoftwareList(`SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall`, "X32")
	if err != nil {
		return nil, err
	}

	return append(sw64, sw32...), nil
}

func GetSoftwareList(baseKey string, arch string) ([]so.Software, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, baseKey, registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		return nil, fmt.Errorf("Error reading from registry: %s", err.Error())
	}
	defer k.Close()

	swList := make([]so.Software, 0)

	subkeys, err := k.ReadSubKeyNames(-1)
	if err != nil {
		return nil, fmt.Errorf("Error reading subkey list from registry: %s", err.Error())
	}
	for _, sw := range subkeys {
		sk, err := registry.OpenKey(registry.LOCAL_MACHINE, baseKey+`\`+sw, registry.QUERY_VALUE)
		if err != nil {
			return nil, fmt.Errorf("Error reading from registry (subkey %s): %s", sw, err.Error())
		}

		dn, _, err := sk.GetStringValue("DisplayName")
		if err == nil {
			swv := so.Software{DisplayName: dn, Arch: arch}

			dv, _, err := sk.GetStringValue("DisplayVersion")
			if err == nil {
				swv.DisplayVersion = dv
			}

			pub, _, err := sk.GetStringValue("Publisher")
			if err == nil {
				swv.Publisher = pub
			}

			id, _, err := sk.GetStringValue("InstallDate")
			if err == nil {
				swv.InstallDate, _ = time.Parse("20060102", id)
			}

			es, _, err := sk.GetIntegerValue("EstimatedSize")
			if err == nil {
				swv.EstimatedSize = es
			}

			cont, _, err := sk.GetStringValue("Contact")
			if err == nil {
				swv.Contact = cont
			}

			hlp, _, err := sk.GetStringValue("HelpLink")
			if err == nil {
				swv.HelpLink = hlp
			}

			isource, _, err := sk.GetStringValue("InstallSource")
			if err == nil {
				swv.InstallSource = isource
			}

			ilocaction, _, err := sk.GetStringValue("InstallLocation")
			if err == nil {
				swv.InstallLocation = ilocaction
			}

			ustring, _, err := sk.GetStringValue("UninstallString")
			if err == nil {
				swv.UninstallString = ustring
			}

			mver, _, err := sk.GetIntegerValue("VersionMajor")
			if err == nil {
				swv.VersionMajor = mver
			}

			mnver, _, err := sk.GetIntegerValue("VersionMinor")
			if err == nil {
				swv.VersionMinor = mnver
			}

			swList = append(swList, swv)
		}
		sk.Close()
	}

	return swList, nil
}
