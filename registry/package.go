package registry

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"gopkg.in/go-playground/validator.v9"
)

type Version struct {
	Name    string `json:"name"    validate:"required"`
	Version string `json:"version" validate:"required"`
}

type Attachment struct {
	Data string `json:"data" validate:"required,base64"`
}

type Package struct {
	Name   string `json:"name"   validate:"required"`
	Author string `json:"author" validate:""`

	DistTags    map[string]string     `json:"dist-tags"`
	Versions    map[string]Version    `json:"versions"     validate:"required,gt=0"`
	Attachments map[string]Attachment `json:"_attachments" validate:"required,gt=0"`
}

var validate = validator.New()

func getVersion(m map[string]Version) Version {
	for _, v := range m {
		return v
	}

	return Version{}
}

func getAttachment(a map[string]Attachment) Attachment {
	for _, a := range a {
		return a
	}

	return Attachment{}
}

// PublishPackage is used to create a new version of package
// if package exists or create a new package if it doesn't
func (core Core) PublishPackage(newData string, user User, registryURL string) error {
	// 1. unmarshal new data
	// 2. validate new data
	// 3. get new version
	// 4. get existing data
	// 5. if package doesn't exist, create new
	// 6. unmarshal existing data
	// 7. if version already exists, return error
	// 8. update dist-tags
	// 9. upload tarball
	// 10. update version tarball url
	// 11. add new version
	// 12. save to database

	var err error

	// 1. unmarshal new data
	var newPkg Package
	err = json.Unmarshal([]byte(newData), &newPkg)
	if err != nil {
		return err
	}

	// 2. validate new data
	err = validate.Struct(newPkg)
	if err != nil {
		return err
	}

	// 3. get new version
	var newVersion = getVersion(newPkg.Versions)

	// 4. get existing data
	existingData, err := core.database.GetPackage(newPkg.Name)
	if err != nil {
		return err
	}

	// 5. if package doesn't exist, create new
	if existingData == "" {
		existingData = fmt.Sprintf(`{"name":"%s","author":"%s","versions":{}}`, newPkg.Name, user.Username)
	}

	// 6. unmarshal existing data
	var existingPkg Package
	err = json.Unmarshal([]byte(existingData), &existingPkg)
	if err != nil {
		return err
	}

	// 7. if version already exists, return error
	if existingPkg.Versions[newVersion.Version].Version != "" {
		return fmt.Errorf("version already exists")
	}

	// 8. update dist-tags
	for k, v := range newPkg.DistTags {
		existingData, err = sjson.Set(existingData, "dist-tags."+k, v)
		if err != nil {
			return err
		}
	}

	// 9. upload tarball
	var attachment = getAttachment(newPkg.Attachments)
	buffer, err := base64.StdEncoding.DecodeString(attachment.Data)
	if err != nil {
		return err
	}
	err = core.storage.WriteTarball(newPkg.Name, newVersion.Version, buffer)
	if err != nil {
		return err
	}

	// 10. update version tarball url
	var escapedVersionNumber = strings.ReplaceAll(newVersion.Version, `.`, `\.`)
	newVersionData, err := sjson.Set(
		gjson.Get(newData, "versions."+escapedVersionNumber).String(),
		"dist.tarball",
		fmt.Sprintf("http://%s/%s/-/%s.tgz", registryURL, newPkg.Name, newVersion.Version),
	)
	if err != nil {
		return err
	}

	// 11. add new version
	existingData, err = sjson.SetRaw(existingData, "versions."+escapedVersionNumber, newVersionData)
	if err != nil {
		return err
	}

	// 12. save to database
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
