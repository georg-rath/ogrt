package inspect

import (
	"bytes"
	"debug/elf"
	"encoding/binary"
	"fmt"
	"io"
)

var ogrtNoteType = []byte{0x54, 0x52, 0x47, 0x4F}

type NoteInfo struct {
	SectionName string
	NoteName    string
	Version     byte
	UUID        string
	Allocatable bool
}

func (n NoteInfo) String() string {
	return fmt.Sprintf("SectionName: %s\nNoteName: %s\nVersion: %d\nUUID: %s\nAlloc: %t",
		n.SectionName, n.NoteName, n.Version, n.UUID, n.Allocatable)
}

// Find OGRT signatures in ELF
func FindSignatures(f io.ReaderAt) (notes []NoteInfo, err error) {
	e, err := elf.NewFile(f)
	if err != nil {
		return nil, err
	}

	for _, s := range e.Sections {
		if s.Type == elf.SHT_NOTE {
			sectionData, err := s.Data()
			if err != nil {
				return nil, err
			}
			if len(sectionData) >= 12 {
				if bytes.Compare(ogrtNoteType, sectionData[8:12]) == 0 {
					nameLength := binary.LittleEndian.Uint32(sectionData[0:4])
					ni := NoteInfo{
						SectionName: s.Name,
						NoteName:    string(sectionData[12 : 12+nameLength]),
						Version:     sectionData[20:21][0],
						Allocatable: s.Flags&elf.SHF_ALLOC != 0,
					}
					switch ni.Version {
					case 0x01:
						ni.UUID = string(sectionData[21:58])
					}
					notes = append(notes, ni)
				}
			}
		}
	}
	return notes, nil
}
