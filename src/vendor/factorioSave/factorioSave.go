package factorioSave

import (
	"log"
	"os"
	"encoding/binary"
)

type version16 struct {
	versionShort16
	Revision    uint16
}
type versionShort16 struct {
	Major       uint16
	Minor       uint16
	Build       uint16
}
type versionShort8 struct {
	Major       uint8
	Minor       uint8
	Build       uint8
}
type Header struct {
	FactorioVersion           version16
	Campaign                  string
	Name                      string
	BaseMod                   string
	Difficulty                uint8
	Finished                  bool
	PlayerWon                 bool
	NextLevel                 string
	CanContinue               bool
	FinishedButContinuing     bool
	SavingReplay              bool
	AllowNonAdminDebugOptions bool
	LoadedFrom                versionShort8
	LoadedFromBuild           uint16
	AllowedCommads            uint8
	NumMods                   uint8
	Mods                      []singleMod
}
type singleMod struct {
	Name    string
	Version versionShort8
	CRC     uint32
}


func ReadHeader(filePath string) (Header, error) {
	var data Header

	fp, err := os.Open(filePath)
	if err != nil {
		log.Printf("error opening file: %s", err)
		return data, err
	}
	defer fp.Close()

	data.FactorioVersion, err = readVersion16(fp)
	if err != nil {
		log.Printf("Cant read FactorioVersion: %s", err)
		return data, err
	}

	if !data.FactorioVersion.CheckCompatibility(0, 16, 0) {
		log.Printf("NOT COMPATIBLE Save-File")
		return data, err
	}

	data.Campaign, err = readUTF8String(fp)
	if err != nil {
		log.Printf("Cant read Campaign: %s", err)
		return data, err
	}

	data.Name, err = readUTF8String(fp)
	if err != nil {
		log.Printf("Cant read Name: %s", err)
		return data, err
	}

	data.BaseMod, err = readUTF8String(fp)
	if err != nil {
		log.Printf("Cant read BaseMod: %s", err)
		return data, err
	}

	data.Difficulty, err = readUint8(fp)
	if err != nil {
		log.Printf("Cant read Difficulty: %s", err)
		return data, err
	}

	data.Finished, err = readBool(fp)
	if err != nil {
		log.Printf("Couln't read Finished bool: %s", err)
		return data, err
	}

	data.PlayerWon, err = readBool(fp)
	if err != nil {
		log.Printf("Couldn't read PlayerWon: %s", err)
		return data, err
	}

	data.NextLevel, err = readUTF8String(fp)
	if err != nil {
		log.Printf("Couldn't read NextLevel: %s", err)
		return data, err
	}

	data.CanContinue, err = readBool(fp)
	if err != nil {
		log.Printf("Couldn't read CanContinue: %s", err)
		return data, err
	}

	data.FinishedButContinuing, err = readBool(fp)
	if err != nil {
		log.Printf("Couldn't read FinishedButContinuing: %s", err)
		return data, err
	}

	data.SavingReplay, err = readBool(fp)
	if err != nil {
		log.Printf("Couldn't read SavingReplay: %s", err)
		return data, err
	}

	data.AllowNonAdminDebugOptions, err = readBool(fp)
	if err != nil {
		log.Printf("Couldn't read allow_non_admin_debug_options: %s", err)
		return data, err
	}

	data.LoadedFrom, err = readVersionShort8(fp)
	if err != nil {
		log.Printf("Couldn't read LoadedFrom: %s", err)
		return data, err
	}

	data.LoadedFromBuild, err = readUint16(fp)
	if err != nil {
		log.Printf("Couldn't read LoadedFromBuild: %s", err)
		return data, err
	}

	data.AllowedCommads, err = readUint8(fp)
	if err != nil {
		log.Printf("Couldn't read AllowedCommands: %s", err)
		return data, err
	}

	data.NumMods, err = readUint8(fp)
	if err != nil {
		log.Printf("Couldn't read NumMods: %s", err)
		return data, err
	}

	for i := uint8(0); i < data.NumMods; i++ {
		SingleMod, err := readSingleMod(fp)
		if err != nil {
			log.Printf("Couldn't read SingleMod: %s", err)
			return data, err
		}

		data.Mods = append(data.Mods, SingleMod)
	}

	log.Println(data)
	return data, nil
}

func readUTF8String(file *os.File) (string, error) {
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

func readUint8(file *os.File) (uint8, error) {
	var err error
	var temp [1]byte
	_, err = file.Read(temp[:])
	if err != nil {
		log.Printf("error reading byte: %s", err)
		return 0, nil
	}

	return uint8(temp[0]), nil
}

func readUint16(file *os.File) (uint16, error) {
	var err error
	var temp [2]byte

	_, err = file.Read(temp[:])
	if err != nil {
		log.Printf("error reading bytes: %s", err)
		return 0, err
	}

	return binary.LittleEndian.Uint16(temp[:]), nil
}

func readUint32(file *os.File) (uint32, error) {
	var err error
	var temp [4]byte

	_, err = file.Read(temp[:])
	if err != nil {
		log.Printf("error reading bytes: %s", err)
		return 0, err
	}

	return binary.LittleEndian.Uint32(temp[:]), nil
}

func readBool(file *os.File) (bool, error) {
	byteAsInt, err := readUint8(file)
	if err != nil {
		log.Printf("error loading Uint8: %s", err)
		return false, err
	}

	return byteAsInt != 0, nil
}

func readVersion16(file *os.File) (version16, error) {
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

func readVersionShort16(file *os.File) (versionShort16, error) {
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

func readVersionShort8(file *os.File) (versionShort8, error) {
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

func readSingleMod(file *os.File) (singleMod, error) {
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

	Mod.CRC, err = readUint32(file)
	if err != nil {
		log.Printf("error loading CRC: %s", err)
		return Mod, err
	}

	return Mod, err
}

func (Version *versionShort16) CheckCompatibility(Major uint16, Minor uint16, Build uint16) (bool) {
	return Major >= Version.Major && Minor >= Version.Minor && Build >= Version.Build
}
