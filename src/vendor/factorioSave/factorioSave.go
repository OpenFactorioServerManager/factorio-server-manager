package factorioSave

import (
	"log"
	"encoding/binary"
	"errors"
	"io"
	"strconv"
	"github.com/Masterminds/semver"
)

type version16 struct {
	versionShort16
	Revision    uint16  `json:"revision"`
}
type versionShort16 struct {
	Major       uint16  `json:"major"`
	Minor       uint16  `json:"minor"`
	Build       uint16  `json:"build"`
}
type versionShort8 struct {
	Major       uint8   `json:"major"`
	Minor       uint8   `json:"minor"`
	Build       uint8   `json:"build"`
}
type Header struct {
	FactorioVersion           version16     `json:"factorio_version"`
	Campaign                  string        `json:"campaign"`
	Name                      string        `json:"name"`
	BaseMod                   string        `json:"base_mod"`
	Difficulty                uint8         `json:"difficulty"`
	Finished                  bool          `json:"finished"`
	PlayerWon                 bool          `json:"player_won"`
	NextLevel                 string        `json:"next_level"`
	CanContinue               bool          `json:"can_continue"`
	FinishedButContinuing     bool          `json:"finished_but_continuing"`
	SavingReplay              bool          `json:"saving_replay"`
	AllowNonAdminDebugOptions bool          `json:"allow_non_admin_debug_options"`
	LoadedFrom                versionShort8 `json:"loaded_from"`
	LoadedFromBuild           uint16        `json:"loaded_from_build"`
	AllowedCommads            uint8         `json:"allowed_commads"`
	NumMods                   uint8         `json:"num_mods"`
	Mods                      []singleMod   `json:"mods"`
}
type singleMod struct {
	Name    string          `json:"name"`
	Version versionShort8   `json:"version"`
	CRC     uint32          `json:"crc"`
}

var ErrorIncompatible = errors.New("incompatible save")
var data Header

func ReadHeader(filePath string) (Header, error) {
	var err error

	datFile, err := openSave(filePath)
	if err != nil {
		log.Printf("error opening file: %s", err)
		return data, err
	}
	defer datFile.Close()

	data.FactorioVersion, err = readVersion16(datFile)
	if err != nil {
		log.Printf("Cant read FactorioVersion: %s", err)
		return data, err
	}

	Constraint, _ := semver.NewConstraint("0.16.0 - 0.17.0")
	Compatible, err := data.FactorioVersion.CheckCompatibility(Constraint)
	if err != nil {
		log.Printf("Error checking compatibility: %s", err)
		return data, err
	}
	if !Compatible {
		log.Printf("NOT COMPATIBLE Save-File")
		log.Println(data)
		return data, ErrorIncompatible
	}

	data.Campaign, err = readUTF8String(datFile)
	if err != nil {
		log.Printf("Cant read Campaign: %s", err)
		return data, err
	}

	data.Name, err = readUTF8String(datFile)
	if err != nil {
		log.Printf("Cant read Name: %s", err)
		return data, err
	}

	data.BaseMod, err = readUTF8String(datFile)
	if err != nil {
		log.Printf("Cant read BaseMod: %s", err)
		return data, err
	}

	data.Difficulty, err = readUint8(datFile)
	if err != nil {
		log.Printf("Cant read Difficulty: %s", err)
		return data, err
	}

	data.Finished, err = readBool(datFile)
	if err != nil {
		log.Printf("Couln't read Finished bool: %s", err)
		return data, err
	}

	data.PlayerWon, err = readBool(datFile)
	if err != nil {
		log.Printf("Couldn't read PlayerWon: %s", err)
		return data, err
	}

	data.NextLevel, err = readUTF8String(datFile)
	if err != nil {
		log.Printf("Couldn't read NextLevel: %s", err)
		return data, err
	}

	data.CanContinue, err = readBool(datFile)
	if err != nil {
		log.Printf("Couldn't read CanContinue: %s", err)
		return data, err
	}

	data.FinishedButContinuing, err = readBool(datFile)
	if err != nil {
		log.Printf("Couldn't read FinishedButContinuing: %s", err)
		return data, err
	}

	data.SavingReplay, err = readBool(datFile)
	if err != nil {
		log.Printf("Couldn't read SavingReplay: %s", err)
		return data, err
	}

	Constraint, _ = semver.NewConstraint(">= 0.16.0")
	Used, err := data.FactorioVersion.CheckCompatibility(Constraint)
	if err != nil {
		log.Printf("Error checking if used: %s", err)
		return data, err
	}
	if Used {
		data.AllowNonAdminDebugOptions, err = readBool(datFile)
		if err != nil {
			log.Printf("Couldn't read allow_non_admin_debug_options: %s", err)
			return data, err
		}
	}

	data.LoadedFrom, err = readVersionShort8(datFile)
	if err != nil {
		log.Printf("Couldn't read LoadedFrom: %s", err)
		return data, err
	}

	data.LoadedFromBuild, err = readUint16(datFile)
	if err != nil {
		log.Printf("Couldn't read LoadedFromBuild: %s", err)
		return data, err
	}

	data.AllowedCommads, err = readUint8(datFile)
	if err != nil {
		log.Printf("Couldn't read AllowedCommands: %s", err)
		return data, err
	}

	data.NumMods, err = readUint8(datFile)
	if err != nil {
		log.Printf("Couldn't read NumMods: %s", err)
		return data, err
	}

	for i := uint8(0); i < data.NumMods; i++ {
		SingleMod, err := readSingleMod(datFile)
		if err != nil {
			log.Printf("Couldn't read SingleMod: %s", err)
			return data, err
		}

		data.Mods = append(data.Mods, SingleMod)
	}

	log.Println(data)
	return data, nil
}

func readUTF8String(file io.ReadCloser) (string, error) {
	var err error
	infoByte := make([]byte, 1)

	_, err = file.Read(infoByte)
	if err != nil {
		log.Printf("Error reading infoByte: %s", err)
		return "", nil
	}
	stringLengthInBytes := int8(infoByte[0])

	stringBytes := make([]byte, stringLengthInBytes)
	_, err = file.Read(stringBytes)
	if err != nil {
		log.Printf("error reading bytes: %s", err)
		return "", err
	}
	finalizedString := string(stringBytes[:])

	return finalizedString, nil
}

func readUint8(file io.ReadCloser) (uint8, error) {
	var err error
	var temp [1]byte
	_, err = file.Read(temp[:])
	if err != nil {
		log.Printf("error reading byte: %s", err)
		return 0, nil
	}

	return uint8(temp[0]), nil
}

func readUint16(file io.ReadCloser) (uint16, error) {
	var err error
	var temp [2]byte

	_, err = file.Read(temp[:])
	if err != nil {
		log.Printf("error reading bytes: %s", err)
		return 0, err
	}

	return binary.LittleEndian.Uint16(temp[:]), nil
}

func readUint32(file io.ReadCloser) (uint32, error) {
	var err error
	var temp [4]byte

	_, err = file.Read(temp[:])
	if err != nil {
		log.Printf("error reading bytes: %s", err)
		return 0, err
	}

	return binary.LittleEndian.Uint32(temp[:]), nil
}

func readBool(file io.ReadCloser) (bool, error) {
	byteAsInt, err := readUint8(file)
	if err != nil {
		log.Printf("error loading Uint8: %s", err)
		return false, err
	}

	return byteAsInt != 0, nil
}

func readVersion16(file io.ReadCloser) (version16, error) {
	var Version version16
	var VersionShort versionShort16
	var err error

	VersionShort, err = readVersionShort16(file)
	if err != nil {
		log.Printf("error reading VersionShort")
		return Version, err
	}

	Version.Major = VersionShort.Major
	Version.Minor = VersionShort.Minor
	Version.Build = VersionShort.Build

	Version.Revision, err = readUint16(file)
	if err != nil {
		log.Printf("error reading revision: %s", err)
		return Version, err
	}

	return Version, nil
}

func readVersionShort16(file io.ReadCloser) (versionShort16, error) {
	var Version versionShort16
	var err error

	Version.Major, err = readUint16(file)
	if err != nil {
		log.Printf("error reading major: %s", err)
		return Version, err
	}

	Version.Minor, err = readUint16(file)
	if err != nil {
		log.Printf("error reading minor: %s", err)
		return Version, err
	}

	Version.Build, err = readUint16(file)
	if err != nil {
		log.Printf("error reading build: %s", err)
		return Version, err
	}

	return Version, err
}

func readVersionShort8(file io.ReadCloser) (versionShort8, error) {
	var Version versionShort8
	var err error

	Version.Major, err = readUint8(file)
	if err != nil {
		log.Printf("error reading major: %s", err)
		return Version, err
	}

	Version.Minor, err = readUint8(file)
	if err != nil {
		log.Printf("error reading minor: %s", err)
		return Version, err
	}

	Version.Build, err = readUint8(file)
	if err != nil {
		log.Printf("error reading build: %s", err)
		return Version, err
	}

	return Version, nil
}

func readSingleMod(file io.ReadCloser) (singleMod, error) {
	var Mod singleMod
	var err error

	Mod.Name, err = readUTF8String(file)
	if err != nil {
		log.Printf("error loading modName: %s", err)
		return Mod, err
	}

	Mod.Version, err = readVersionShort8(file)
	if err != nil {
		log.Printf("error loading modVersion: %s", err)
		return Mod, err
	}

	Constraint, _ := semver.NewConstraint("> 0.15.0")
	Used, err := data.FactorioVersion.CheckCompatibility(Constraint)
	if err != nil {
		log.Printf("Error checking used of CRC: %s", err)
		return Mod, err
	}
	if Used {
		Mod.CRC, err = readUint32(file)
		if err != nil {
			log.Printf("error loading CRC: %s", err)
			return Mod, err
		}
	}

	return Mod, err
}

func (Version *versionShort16) CheckCompatibility(constraints *semver.Constraints) (bool, error) {
	Ver, err := semver.NewVersion(strconv.Itoa(int(Version.Major)) + "." + strconv.Itoa(int(Version.Minor)) + "." + strconv.Itoa(int(Version.Build)))
	if err != nil {
		log.Printf("Error creating semver-version: %s", err)
		return false, err
	}

	return constraints.Check(Ver), nil
}
