package protamoattr

import (
	"reflect"
	"strings"
)

const (
	protoTagKey      = "protobuf"
	protoOneofTagKey = "protobuf_oneof"
)

type tag struct {
	Name                         string
	Ignore                       bool
	OmitEmpty                    bool
	OmitEmptyElem                bool
	NullEmpty                    bool
	NullEmptyElem                bool
	AsString                     bool
	AsBinSet, AsNumSet, AsStrSet bool
	AsUnixTime                   bool
	ProtoOneOf                   bool
}

func (t *tag) parse(structTag reflect.StructTag, opts structFieldOptions) {
	t.parseAVTag(structTag)
	// Because MarshalOptions.TagKey must be explicitly set.
	if opts.TagKey != "" && opts.TagKey != defaultTagKey {
		t.parseStructTag(opts.TagKey, structTag)
	}

	// specific to protamo is that the protobuf field nr takes
	// presidence if it is set. To provide stable item structure even when
	// refactoring field names and order.
	t.parseProtoTag(structTag)
}

func (t *tag) parseProtoTag(structTag reflect.StructTag) {
	tagStr := structTag.Get(protoOneofTagKey)
	if len(tagStr) > 0 {
		t.ProtoOneOf = true
	}

	tagStr = structTag.Get(protoTagKey)
	if len(tagStr) == 0 {
		return
	}

	parts := strings.Split(tagStr, ",")
	if len(parts) < 2 {
		return
	}

	if name := parts[1]; name == "-" {
		t.Name = ""
		t.Ignore = true
	} else {
		t.Name = name
		t.Ignore = false
	}
}

func (t *tag) parseAVTag(structTag reflect.StructTag) {
	tagStr := structTag.Get(defaultTagKey)
	if len(tagStr) == 0 {
		return
	}

	t.parseTagStr(tagStr)
}

func (t *tag) parseStructTag(tag string, structTag reflect.StructTag) {
	tagStr := structTag.Get(tag)
	if len(tagStr) == 0 {
		return
	}

	t.parseTagStr(tagStr)
}

func (t *tag) parseTagStr(tagStr string) {
	parts := strings.Split(tagStr, ",")
	if len(parts) == 0 {
		return
	}

	if name := parts[0]; name == "-" {
		t.Name = ""
		t.Ignore = true
	} else {
		t.Name = name
		t.Ignore = false
	}

	for _, opt := range parts[1:] {
		switch opt {
		case "omitempty":
			t.OmitEmpty = true
		case "omitemptyelem":
			t.OmitEmptyElem = true
		case "nullempty":
			t.NullEmpty = true
		case "nullemptyelem":
			t.NullEmptyElem = true
		case "string":
			t.AsString = true
		case "binaryset":
			t.AsBinSet = true
		case "numberset":
			t.AsNumSet = true
		case "stringset":
			t.AsStrSet = true
		case "unixtime":
			t.AsUnixTime = true
		}
	}
}
