package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

type SaveHeader struct {
	FactorioVersion           version64 `json:"factorio_version"`
	Campaign                  string    `json:"campaign"`
	Name                      string    `json:"name"`
	BaseMod                   string    `json:"base_mod"`
	Difficulty                uint8     `json:"difficulty"`
	Finished                  bool      `json:"finished"`
	PlayerWon                 bool      `json:"player_won"`
	NextLevel                 string    `json:"next_level"`
	CanContinue               bool      `json:"can_continue"`
	FinishedButContinuing     bool      `json:"finished_but_continuing"`
	SavingReplay              bool      `json:"saving_replay"`
	AllowNonAdminDebugOptions bool      `json:"allow_non_admin_debug_options"`
	LoadedFrom                version24 `json:"loaded_from"`
	LoadedFromBuild           uint16    `json:"loaded_from_build"`
	AllowedCommands           uint8     `json:"allowed_commands"`
	Mods                      []Mod     `json:"mods"`
}

type Mod struct {
	Name    string    `json:"name"`
	Version version24 `json:"version"`
	CRC     uint32    `json:"crc"`
}

func (h *SaveHeader) ReadFrom(r io.Reader) (err error) {
	var scratch [8]byte

	_, err = r.Read(scratch[:8])
	if err != nil {
		return err
	}
	if err := h.FactorioVersion.UnmarshalBinary(scratch[:8]); err != nil {
		return fmt.Errorf("read FactorioVersion: %v", err)
	}

	atLeast016 := !h.FactorioVersion.Less(Version{0, 16, 0})

	h.Campaign, err = readString(r, atLeast016)
	if err != nil {
		return fmt.Errorf("read Campaign: %v", err)
	}

	h.Name, err = readString(r, atLeast016)
	if err != nil {
		return fmt.Errorf("read Name: %v", err)
	}

	h.BaseMod, err = readString(r, atLeast016)
	if err != nil {
		return fmt.Errorf("read BaseMod: %v", err)
	}

	_, err = r.Read(scratch[:1])
	if err != nil {
		return fmt.Errorf("read Difficulty: %v", err)
	}
	h.Difficulty = scratch[0]

	_, err = r.Read(scratch[:1])
	if err != nil {
		return fmt.Errorf("read Finished: %v", err)
	}
	h.Finished = scratch[0] != 0

	_, err = r.Read(scratch[:1])
	if err != nil {
		return fmt.Errorf("read PlayerWon: %v", err)
	}
	h.PlayerWon = scratch[0] != 0

	h.NextLevel, err = readString(r, atLeast016)
	if err != nil {
		return fmt.Errorf("read NextLevel: %v", err)
	}

	if !h.FactorioVersion.Less(Version{0, 12, 0}) {
		_, err = r.Read(scratch[:1])
		if err != nil {
			return fmt.Errorf("read CanContinue: %v", err)
		}
		h.CanContinue = scratch[0] != 0

		_, err = r.Read(scratch[:1])
		if err != nil {
			return fmt.Errorf("read FinishedButContinuing: %v", err)
		}
		h.FinishedButContinuing = scratch[0] != 0
	}

	_, err = r.Read(scratch[:1])
	if err != nil {
		return fmt.Errorf("read SavingReplay: %v", err)
	}
	h.SavingReplay = scratch[0] != 0

	if atLeast016 {
		_, err = r.Read(scratch[:1])
		if err != nil {
			return fmt.Errorf("read AllowNonAdminDebugOptions: %v", err)
		}
		h.AllowNonAdminDebugOptions = scratch[0] != 0
	}

	_, err = r.Read(scratch[:3])
	if err != nil {
		return err
	}
	if err := h.LoadedFrom.UnmarshalBinary(scratch[:3]); err != nil {
		return fmt.Errorf("read LoadedFrom: %v", err)
	}

	_, err = r.Read(scratch[:2])
	if err != nil {
		return fmt.Errorf("read LoadedFromBuild: %v", err)
	}
	h.LoadedFromBuild = binary.LittleEndian.Uint16(scratch[:2])

	_, err = r.Read(scratch[:1])
	if err != nil {
		return fmt.Errorf("read AllowedCommands: %v", err)
	}
	h.AllowedCommands = scratch[0]

	var n uint32
	if atLeast016 {
		n, err = readOptimUint32(r)
		if err != nil {
			return fmt.Errorf("read num mods: %v", err)
		}
	} else {
		_, err = r.Read(scratch[:4])
		if err != nil {
			return fmt.Errorf("read num mods: %v", err)
		}
		n = binary.LittleEndian.Uint32(scratch[:4])
	}

	for i := uint32(0); i < n; i++ {
		var m Mod
		if err = (&m).ReadFrom(r, h.FactorioVersion.Version); err != nil {
			return fmt.Errorf("read mod: %v", err)
		}
		h.Mods = append(h.Mods, m)
	}

	return nil
}

func readOptimUint32(r io.Reader) (uint32, error) {
	var b [4]byte
	_, err := r.Read(b[:1])
	if err != nil {
		return 0, err
	}
	if b[0] != 0xFF {
		return uint32(b[0]), nil
	}
	_, err = r.Read(b[:4])
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(b[:4]), nil
}

func readString(r io.Reader, optimized bool) (s string, err error) {
	var n uint32

	if optimized {
		n, err = readOptimUint32(r)
		if err != nil {
			return "", err
		}
	} else {
		var b [4]byte
		_, err := r.Read(b[:])
		if err != nil {
			return "", fmt.Errorf("failed to read string length: %v", err)
		}
		n = uint32(binary.LittleEndian.Uint32(b[:]))
	}

	d := make([]byte, n)
	_, err = r.Read(d)
	if err != nil {
		return "", fmt.Errorf("failed to read string: %v", err)
	}

	return string(d), nil
}

func (m *Mod) ReadFrom(r io.Reader, game Version) (err error) {
	m.Name, err = readString(r, true)
	if err != nil {
		return fmt.Errorf("read Name: %v", err)
	}

	var scratch [4]byte
	_, err = r.Read(scratch[:3])
	if err != nil {
		return err
	}
	if err := m.Version.UnmarshalBinary(scratch[:3]); err != nil {
		return fmt.Errorf("read Version: %v", err)
	}

	if game.Greater(Version{0, 15, 0}) {
		_, err = r.Read(scratch[:4])
		if err != nil {
			return err
		}
		m.CRC = binary.LittleEndian.Uint32(scratch[:4])
	}

	return nil
}
