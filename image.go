package main

import (
	"fmt"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

func invertImageBytes(data []byte) []byte {
	out := make([]byte, len(data))
	for i, b := range data {
		out[i] = 255 - b
	}
	return out
}

func invertPageImages(xrt *model.XRefTable, pageDict types.Dict) error {
	resourcesObj, ok := pageDict["Resources"]
	if !ok {
		return nil
	}
	resourcesDict, err := xrt.DereferenceDict(resourcesObj)
	if err != nil || resourcesDict == nil {
		return nil
	}

	xObjEntry, ok := resourcesDict["XObject"]
	if !ok {
		return nil
	}
	xObjDict, err := xrt.DereferenceDict(xObjEntry)
	if err != nil || xObjDict == nil {
		return nil
	}

	for name, obj := range xObjDict {
		if err := invertImageXObject(xrt, xObjDict, name, obj); err != nil {
			return fmt.Errorf("XObject %q: %w", name, err)
		}
	}
	return nil
}

func invertImageXObject(xrt *model.XRefTable, xObjDict types.Dict, name string, obj types.Object) error {
	sd, _, err := xrt.DereferenceStreamDict(obj)
	if err != nil || sd == nil {
		return nil
	}

	if !sd.Image() {
		return nil
	}

	// Only handle FlateDecode; skip others (DCT/JPEG etc.) silently.
	if !sd.HasSoleFilterNamed("FlateDecode") {
		return nil
	}

	// Decode the stream to get raw pixel bytes.
	if err := sd.Decode(); err != nil {
		return fmt.Errorf("decode stream: %w", err)
	}
	if sd.Content == nil {
		return nil
	}

	inverted := invertImageBytes(sd.Content)

	newSD, err := xrt.NewStreamDictForBuf(inverted)
	if err != nil {
		return fmt.Errorf("new stream dict: %w", err)
	}
	// Carry over image metadata (dimensions, color space, BitsPerComponent, etc.)
	for k, v := range sd.Dict {
		if k != "Filter" && k != "Length" && k != "DecodeParms" {
			newSD.Dict[k] = v
		}
	}
	if err := newSD.Encode(); err != nil {
		return fmt.Errorf("encode stream: %w", err)
	}

	indRef, err := xrt.IndRefForNewObject(*newSD)
	if err != nil {
		return fmt.Errorf("insert stream: %w", err)
	}
	xObjDict[name] = *indRef
	return nil
}
