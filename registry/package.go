package registry

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Version struct {
	Version string
}

type Attachment struct {
	Data string `json:"data"`
}

type Package struct {
	Name   string `json:"name"`
	Author string `json:"author"`

	DistTags    map[string]string     `json:"dist-tags"`
	Versions    map[string]Version    `json:"versions"`
	Attachments map[string]Attachment `json:"_attachments"`
}

func getVersion(m map[string]Version) Version {
	for _, v := range m {
		return v
		break
	}

	return Version{}
}

func getAttachment(a map[string]Attachment) Attachment {
	for _, a := range a {
		return a
		break
	}

	return Attachment{}
}

// PublishPackage is used to create a new version of package
// if package exists or create a new package if it doesn't
func (core Core) PublishPackage(newData string, user User, registryURL string) error {
	// 1. unmarshal new data
	// 2. if new data doesn't include version, return error
	// 3. get existing data
	// 4. if package doesn't exist, create new
	// 5. unmarshal existing data
	// 6. if version already exists, return error
	// 7. update dist-tags
	// 8. upload tarball
	// 9. update version tarball url
	// 10. add new version
	// 11. save to database

	var err error

	// 1. unmarshal new data
	var newPkg Package
	err = json.Unmarshal([]byte(newData), &newPkg)
	if err != nil {
		return err
	}

	// 2. if new data doesn't include version, return error
	var newVersion = getVersion(newPkg.Versions)
	if newVersion.Version == "" {
		return fmt.Errorf("missing version")
	}

	// 3. get existing data
	existingData, err := core.database.GetPackage(newPkg.Name)
	if err != nil {
		return err
	}

	// 4. if package doesn't exist, create new
	if existingData == "" {
		existingData = fmt.Sprintf(`{"name":"%s","author":"%s","versions":{}}`, newPkg.Name, user.Username)
	}

	// 5. unmarshal existing data
	var existingPkg Package
	err = json.Unmarshal([]byte(existingData), &existingPkg)
	if err != nil {
		return err
	}

	// 6. if version already exists, return error
	if existingPkg.Versions[newVersion.Version].Version != "" {
		return fmt.Errorf("version already exists")
	}

	// 7. update dist-tags
	for k, v := range newPkg.DistTags {
		existingData, err = sjson.Set(existingData, "dist-tags."+k, v)
		if err != nil {
			return err
		}
	}

	// 8. upload tarball
	var attachment = getAttachment(newPkg.Attachments)
	buffer, err := base64.StdEncoding.DecodeString(attachment.Data)
	if err != nil {
		return err
	}
	err = core.storage.WriteTarball(newPkg.Name, newVersion.Version, buffer)
	if err != nil {
		return err
	}

	// 9. update version tarball url
	var escapedVersionNumber = strings.ReplaceAll(newVersion.Version, `.`, `\.`)
	newVersionData, err := sjson.Set(
		gjson.Get(newData, "versions."+escapedVersionNumber).String(),
		"dist.tarball",
		fmt.Sprintf("http://%s/%s/-/%s.tgz", registryURL, newPkg.Name, newVersion.Version),
	)
	if err != nil {
		return err
	}

	// 10. add new version
	existingData, err = sjson.SetRaw(existingData, "versions."+escapedVersionNumber, newVersionData)
	if err != nil {
		return err
	}

	// 11. save to database
	return core.database.SetPackage(existingPkg.Name, existingData)
}

// GetPackage is used to get an existing package
func (core Core) GetPackage(name string) (string, error) {
	return core.database.GetPackage(name)
}

// GetPackageTarball is used to get a tarball for specific package version
func (core Core) GetPackageTarball(name string, version string) ([]byte, error) {
	return core.storage.ReadTarball(name, version)
}
