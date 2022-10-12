package main

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/srikrsna/nidhi/gen/nidhi"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"google.golang.org/protobuf/types/pluginpb"
)

const (
	generatedFilenameExtension = ".nidhi.go"
	generatedPackageSuffix     = "nidhi"
	nidhiPkg                   = protogen.GoImportPath("github.com/srikrsna/nidhi")
	contextPkg                 = protogen.GoImportPath("context")
	sqlPkg                     = protogen.GoImportPath("database/sql")
	jsoniterPkg                = protogen.GoImportPath("github.com/json-iterator/go")
	protojsonPkg               = protogen.GoImportPath("google.golang.org/protobuf/encoding/protojson")
	base64Pkg                  = protogen.GoImportPath("encoding/base64")
	jsonPkg                    = protogen.GoImportPath("encoding/json")
)

var (
	wktSet = map[protoreflect.FullName]bool{
		(new(structpb.NullValue)).Descriptor().FullName():                  true,
		(&structpb.Struct{}).ProtoReflect().Descriptor().FullName():        true,
		(&structpb.ListValue{}).ProtoReflect().Descriptor().FullName():     true,
		(&structpb.Value{}).ProtoReflect().Descriptor().FullName():         true,
		(&fieldmaskpb.FieldMask{}).ProtoReflect().Descriptor().FullName():  true,
		(&timestamppb.Timestamp{}).ProtoReflect().Descriptor().FullName():  true,
		(&durationpb.Duration{}).ProtoReflect().Descriptor().FullName():    true,
		(&anypb.Any{}).ProtoReflect().Descriptor().FullName():              true,
		(&emptypb.Empty{}).ProtoReflect().Descriptor().FullName():          true,
		(&wrapperspb.BoolValue{}).ProtoReflect().Descriptor().FullName():   true,
		(&wrapperspb.StringValue{}).ProtoReflect().Descriptor().FullName(): true,
		(&wrapperspb.BytesValue{}).ProtoReflect().Descriptor().FullName():  true,
		(&wrapperspb.Int32Value{}).ProtoReflect().Descriptor().FullName():  true,
		(&wrapperspb.Int64Value{}).ProtoReflect().Descriptor().FullName():  true,
		(&wrapperspb.UInt32Value{}).ProtoReflect().Descriptor().FullName(): true,
		(&wrapperspb.UInt64Value{}).ProtoReflect().Descriptor().FullName(): true,
		(&wrapperspb.FloatValue{}).ProtoReflect().Descriptor().FullName():  true,
		(&wrapperspb.DoubleValue{}).ProtoReflect().Descriptor().FullName(): true,
	}
)

func main() {
	protogen.Options{}.Run(func(plugin *protogen.Plugin) error {
		plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, file := range plugin.Files {
			if file.Generate {
				gen(plugin, file)
			}
		}
		return nil
	})
}

func gen(plugin *protogen.Plugin, file *protogen.File) {
	type messageToGen struct {
		*protogen.Message
		*protogen.Field
	}
	var messagesToGen []messageToGen
	for _, message := range file.Messages {
		for _, field := range message.Fields {
			if proto.GetExtension(field.Desc.Options(), nidhi.E_DocumentId).(bool) {
				messagesToGen = append(messagesToGen, messageToGen{message, field})
			}
		}
	}
	if len(messagesToGen) == 0 {
		return
	}
	file.GoPackageName += generatedPackageSuffix
	dir := filepath.Dir(file.GeneratedFilenamePrefix)
	base := filepath.Base(file.GeneratedFilenamePrefix)
	file.GeneratedFilenamePrefix = filepath.Join(
		dir,
		string(file.GoPackageName),
		base,
	)
	genFile := plugin.NewGeneratedFile(
		file.GeneratedFilenamePrefix+generatedFilenameExtension,
		protogen.GoImportPath(path.Join(
			string(file.GoImportPath),
			string(file.GoPackageName),
		)),
	)
	genHeader(genFile, file)
	for _, message := range messagesToGen {
		genMessage(genFile, message.Message, message.Field, plugin)
	}
}

func genHeader(g *protogen.GeneratedFile, file *protogen.File) {
	g.P("// Code generated by ", filepath.Base(os.Args[0]), ". DO NOT EDIT.")
	g.P("//")
	if file.Proto.GetOptions().GetDeprecated() {
		g.P("//", file.Desc.Path(), " is a deprecated file.")
	} else {
		g.P("// Source: ", file.Desc.Path())
	}
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()
}

func genMessage(g *protogen.GeneratedFile, message *protogen.Message, idField *protogen.Field, plugin *protogen.Plugin) {
	genUpdates(g, message, idField, plugin)
	genSchema(g, message, idField)
	genQuery(g, message)
	genConj(g, message)
	genNewStore(g, message, idField)
	genFields(g, message, idField)
}

func genUpdates(g *protogen.GeneratedFile, message *protogen.Message, idField *protogen.Field, plugin *protogen.Plugin) {
	name := message.GoIdent.GoName
	g.P("// ", name, "Updates is an updates type that can be passed to")
	g.P("// [*nidhi.Store.Update] and [*nidhi.Store.UpdateMany]")
	g.P("type ", name, "Updates struct {")
	for _, field := range message.Fields {
		if field == idField {
			continue
		}
		g.P(append([]any{field.GoName, " "}, append(lkpFieldType(field, plugin), " `json:\"", field.Desc.JSONName(), ",omitempty\"`")...)...)
	}
	g.P("}")
	g.P("")
	g.P("func (u *", name, "Updates) WriteJSON(w *", jsoniterPkg.Ident("Stream"), ") {")
	g.P("if u == nil {")
	g.P("w.WriteEmptyObject()")
	g.P("return")
	g.P("}")
	g.P("first := true")
	g.P("w.WriteObjectStart()")
	for _, field := range message.Fields {
		if field == idField {
			continue
		}
		g.P("if u.", field.GoName, " != nil {")
		g.P("if !first {")
		g.P("w.WriteMore()")
		g.P("}")
		g.P(`w.WriteObjectField("`, field.Desc.JSONName(), `")`)
		genUpdateFieldMarshaler(g, field.Desc, "*u."+field.GoName, false)
		g.P("first = false")
		g.P("}")
	}
	g.P("w.WriteObjectEnd()")
	g.P("}")
}

func genUpdateFieldMarshaler(g *protogen.GeneratedFile, fd protoreflect.FieldDescriptor, name string, inList bool) {
	switch {
	case fd.IsList() && !inList:
		g.P("w.WriteArrayStart()")
		g.P("ap := false")
		g.P("for _, v := range ", name, " {")
		g.P("if !ap {")
		g.P("w.WriteMore()")
		g.P("}")
		genUpdateFieldMarshaler(g, fd, "v", true)
		g.P("ap = false")
		g.P("}")
		g.P("w.WriteArrayEnd()")
	case fd.IsMap():
		g.P("w.WriteObjectStart()")
		g.P("mp := false")
		g.P("for k, v := range ", name, " {")
		g.P("if !mp {")
		g.P("w.WriteMore()")
		g.P("}")
		switch fd.MapKey().Kind() {
		case protoreflect.BoolKind:
			g.P("w.WriteObjectField(strconv.FormatBool(k))")
		case protoreflect.StringKind:
			g.P("w.WriteObjectField(k)")
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
			g.P("w.WriteObjectField(strconv.FormatInt(int64(k), 10))")
		case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
			g.P("w.WriteObjectField(strconv.FormatInt(k, 10))")
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
			g.P("w.WriteObjectField(strconv.FormatUint(uint64(k), 10))")
		case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
			g.P("w.WriteObjectField(strconv.FormatUint(k, 10))")
		}
		genUpdateFieldMarshaler(g, fd.MapValue(), "v", false)
		g.P("mp = false")
		g.P("}")
		g.P("w.WriteObjectEnd()")
	case fd.Message() != nil:
		if wktSet[fd.Message().FullName()] {
			g.P("data, err := ", protojsonPkg.Ident("Marshal"), "(", name, ")")
			g.P("if err != nil {")
			g.P("w.Error = err")
			g.P("return")
			g.P("}")
			g.P("w.WriteVal(", jsonPkg.Ident("RawMessage"), "(data))")
		} else {
			g.P("if m, ok := any(", name, ").(interface{ WriteJSON(*", jsoniterPkg.Ident("Stream"), ") }); ok {")
			g.P("m.WriteJSON(w)")
			g.P("} else {")
			g.P("data, err := ", protojsonPkg.Ident("Marshal"), "(", name, ")")
			g.P("if err != nil {")
			g.P("w.Error = err")
			g.P("return")
			g.P("}")
			g.P("w.WriteVal(", jsonPkg.Ident("RawMessage"), "(data))")
			g.P("}")
		}
	case fd.Enum() != nil:
		g.P("w.WriteString((", name, ").String())")
	case fd.Kind() == protoreflect.BytesKind:
		g.P("w.WriteString(", base64Pkg.Ident("StdEncoding.EncodeToString"), "(", name, "))")
	default:
		g.P("w.Write", strings.Title(lkpScalarFieldTyp(fd.Kind())), "(", name, ")")
		// Only scalar fields
	}
}

func genQuery(g *protogen.GeneratedFile, message *protogen.Message) {
	name := message.GoIdent.GoName
	g.P("// ", name, "Query is an alias for nidhi.Query[", name, "Field]")
	g.P("type ", name, "Query = nidhi.Query[", name, "Field]")
}

func genConj(g *protogen.GeneratedFile, message *protogen.Message) {
	name := message.GoIdent.GoName
	g.P("// ", name, "Conj is an alias for nidhi.Conj[", name, "Field]")
	g.P("type ", name, "Conj = nidhi.Conj[", name, "Field]")
}

func genSchema(g *protogen.GeneratedFile, message *protogen.Message, idField *protogen.Field) {
	g.P("// ", message.GoIdent.GoName, "Schema is the set of path selectors for the ", message.GoIdent.GoName, " type.")
	g.P("var ", message.GoIdent.GoName, "Schema = struct {")
	lowerMessageName := lowerMessageName(message)
	for _, field := range message.Fields {
		g.P(field.GoName, " ", lowerMessageName, field.GoName)
	}
	g.P("}{}")
}

func genNewStore(g *protogen.GeneratedFile, m *protogen.Message, id *protogen.Field) {
	g.P("// New", m.GoIdent.GoName, "Store is a document store for ", m.GoIdent.GoName)
	g.P("func New", m.GoIdent.GoName, "Store(")
	g.P("ctx ", contextPkg.Ident("Context"), ",")
	g.P("db *", sqlPkg.Ident("DB"), ",")
	g.P("opt ", nidhiPkg.Ident("StoreOptions"), ",")
	g.P(") (*", nidhiPkg.Ident("Store"), "[", m.GoIdent, "], error) {")
	g.P("return ", nidhiPkg.Ident("NewStore"), "(")
	g.P("ctx,")
	g.P("db,")
	g.P(`"`, strings.ReplaceAll(string(m.Desc.ParentFile().Package()), ".", "_"), `",`)
	g.P(`"`, strings.ToLower(m.GoIdent.GoName), `",`)
	g.P("[]string{")
	for _, field := range m.Fields {
		g.P(`"`, field.Desc.JSONName(), `",`)
	}
	g.P("},")
	g.P("func(x *", m.GoIdent, ") string { return x.", id.GoName, " },")
	g.P("func(x *", m.GoIdent, ", id string) {x.", id.GoName, " = id },")
	g.P("opt,")
	g.P(")")
	g.P("}")
}

func genFields(g *protogen.GeneratedFile, m *protogen.Message, id *protogen.Field) {
	fieldEmbed := genFieldInterface(g, m)
	genSelectors(g, m, fieldEmbed, lowerMessageName(m), "$", 0)
}

func genSelectors(g *protogen.GeneratedFile, m *protogen.Message, fieldEmbed, typePrefix, dbPrefix string, inSlice uint) {
	for _, field := range m.Fields {
		typeName := typePrefix + field.GoName
		g.P("type ", typeName, " struct{ ")
		g.P(fieldEmbed)
		if field.Message != nil {
			for _, subField := range field.Message.Fields {
				g.P(subField.GoName, " ", typePrefix, field.GoName, subField.GoName)
			}
		}
		g.P("}")
		genSelectorFunc(g, field, typeName, dbPrefix, inSlice)
		if field.Message != nil {
			var (
				sliceExt   string
				sliceCount uint
			)
			if field.Desc.IsList() {
				sliceExt = "[*]"
				sliceCount = 1
			}
			genSelectors(g, field.Message, fieldEmbed, typePrefix+field.GoName, dbPrefix+"."+field.Desc.JSONName()+sliceExt, inSlice+sliceCount)
		}
	}
}

func genSelectorFunc(g *protogen.GeneratedFile, field *protogen.Field, fieldType, prefix string, inSlice uint) {
	dataType, cond, defaultValue := getDbTypeCondAndDefault(field, inSlice)
	g.P("func (", fieldType, ") Selector() string { return `JSON_VALUE(`+", nidhiPkg.Ident("ColDoc"), "+`::jsonb, '", prefix+"."+field.Desc.JSONName(), "' RETURNING ", dataType, " DEFAULT ", defaultValue, " ON EMPTY)` }")
	g.P("func (f ", fieldType, ") Is(c *", nidhiPkg.Ident(cond), ") (", fieldType, ", ", nidhiPkg.Ident("Cond"), ") { return f, c }")
}

func genFieldInterface(g *protogen.GeneratedFile, m *protogen.Message) string {
	g.P("type ", m.GoIdent.GoName, "Field interface {")
	g.P(nidhiPkg.Ident("Field"))
	g.P(lowerMessageName(m), "Field()")
	g.P("}")
	g.P("type base", m.GoIdent.GoName, "Field struct{}")
	fieldEmbed := "base" + m.GoIdent.GoName + "Field"
	g.P("func (", fieldEmbed, ") ", lowerMessageName(m), "Field() {}")
	return fieldEmbed
}

func lowerMessageName(m *protogen.Message) string {
	return strings.ToLower(m.GoIdent.GoName[:1]) + m.GoIdent.GoName[1:]
}

func getDbTypeCondAndDefault(field *protogen.Field, inSlice uint) (typ string, cond string, def string) {
	defer func() {
		if field.Desc.HasOptionalKeyword() {
			def = "NULL"
		}
		sliceCount := inSlice
		if field.Desc.IsList() {
			sliceCount++
		}
		if sliceCount > 0 && typ != "JSONB" {
			cond = strings.ReplaceAll(cond, "Cond", "SliceCond")
			typ += strings.Repeat("[]", int(sliceCount))
			def = "'{}'"
		}
	}()
	switch {
	case field.Desc.IsMap():
		return "JSONB", "JsonCond", "'{}'"
	case field.Message != nil:
		if wktSet[field.Message.Desc.FullName()] {
			return lkpWktTyp(field.Message.Desc.FullName())
		} else {
			return "JSONB", "JsonCond", "'{}'"
		}
	case field.Enum != nil:
		vd := field.Enum.Desc.Values().ByNumber(0)
		if vd == nil {
			log.Println("Enum doesn't have a ZERO based default value")
		}
		return "TEXT", "StringCond", "'" + string(vd.Name()) + "'"
	default:
		// Only scalar fields
		return lkpScalarTyp(field.Desc.Kind())
	}
}

func lkpScalarTyp(kind protoreflect.Kind) (string, string, string) {
	switch kind {
	case protoreflect.StringKind, protoreflect.BytesKind:
		return "TEXT", "StringCond", "''"
	case protoreflect.BoolKind:
		return "BOOLEAN", "BoolCond", "FALSE"
	case protoreflect.DoubleKind:
		return "DOUBLE PRECISION", "FloatCond", "0"
	case protoreflect.FloatKind:
		return "REAL", "FloatCond", "0"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind, protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "INTEGER", "IntCond", "0"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind, protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "BIGINT", "IntCond", "0"
	default:
		panic("this is a bug")

	}
}

func lkpWktTyp(name protoreflect.FullName) (string, string, string) {
	switch name {
	case (&structpb.Struct{}).ProtoReflect().Descriptor().FullName(),
		(&structpb.ListValue{}).ProtoReflect().Descriptor().FullName(),
		(&structpb.Value{}).ProtoReflect().Descriptor().FullName(),
		(&anypb.Any{}).ProtoReflect().Descriptor().FullName(),
		(new(structpb.NullValue)).Descriptor().FullName():
		return "JSONB", "JsonCond", "NULL"
	case (&fieldmaskpb.FieldMask{}).ProtoReflect().Descriptor().FullName():
		return "TEXT", "StringCond", "''"
	case (&timestamppb.Timestamp{}).ProtoReflect().Descriptor().FullName():
		return "TIMESTAMP", "TimeCond", "'1970-01-01 00:00:00'"
	case (&durationpb.Duration{}).ProtoReflect().Descriptor().FullName():
		return "SECOND P 6 ", "UNKNOWN", "'0s'"
	case (&emptypb.Empty{}).ProtoReflect().Descriptor().FullName():
		return "JSONB", "JsonCond", "'{}'"
	case (&wrapperspb.BoolValue{}).ProtoReflect().Descriptor().FullName():
		return "BOOLEAN", "BoolCond", "FALSE"
	case (&wrapperspb.StringValue{}).ProtoReflect().Descriptor().FullName(), (&wrapperspb.BytesValue{}).ProtoReflect().Descriptor().FullName():
		return "TEXT", "StringCond", "''"
	case (&wrapperspb.Int32Value{}).ProtoReflect().Descriptor().FullName():
		return "INTEGER", "IntCond", "0"
	case (&wrapperspb.Int64Value{}).ProtoReflect().Descriptor().FullName():
		return "BIGINT", "IntCond", "0"
	case (&wrapperspb.UInt32Value{}).ProtoReflect().Descriptor().FullName():
		return "INTEGER", "IntCond", "0"
	case (&wrapperspb.UInt64Value{}).ProtoReflect().Descriptor().FullName():
		return "BIGINT", "IntCond", "0"
	case (&wrapperspb.FloatValue{}).ProtoReflect().Descriptor().FullName():
		return "REAL", "FloatCond", "0"
	case (&wrapperspb.DoubleValue{}).ProtoReflect().Descriptor().FullName():
		return "DOUBLE PRECISION", "FloatCond", "0"
	default:
		panic("unknown wkt")
	}
}

func lkpFieldType(field *protogen.Field, plugin *protogen.Plugin) (fn []any) {
	defer func() {
		if field.Desc.HasOptionalKeyword() {
			fn = append([]any{"*"}, fn...)
		}
		if field.Desc.IsList() {
			fn[0] = "*[]"
		}
	}()
	switch {
	case field.Desc.IsMap():
		return append([]any{"*", "map[", lkpScalarFieldTyp(field.Desc.MapKey().Kind()), "]"}, lkpMapValueImport(field.Desc.MapValue(), plugin)...)
	case field.Message != nil:
		return []any{"*", "*", field.Message.GoIdent}
	case field.Enum != nil:
		return []any{"*", field.Enum.GoIdent}
	default:
		// Only scalar fields
		return []any{"*", lkpScalarFieldTyp(field.Desc.Kind())}
	}
}

func lkpScalarFieldTyp(kind protoreflect.Kind) string {
	switch kind {
	case protoreflect.StringKind:
		return "string"
	case protoreflect.BytesKind:
		return "[]byte"
	case protoreflect.BoolKind:
		return "bool"
	case protoreflect.DoubleKind:
		return "float64"
	case protoreflect.FloatKind:
		return "float32"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return "int32"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return "int64"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "uint32"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "uint64"
	default:
		// This should never happen
		return "Unknown"
	}
}

func lkpMapValueImport(fd protoreflect.FieldDescriptor, plugin *protogen.Plugin) []any {
	file := plugin.FilesByPath[fd.ParentFile().Path()]
	switch {
	case fd.Message() != nil:
		for _, msg := range file.Messages {
			if msg.Desc.FullName() == fd.Message().FullName() {
				return []any{"*", msg.GoIdent}
			}
		}
	case fd.Enum() != nil:
		for _, enum := range file.Enums {
			if enum.Desc.FullName() == fd.Enum().FullName() {
				return []any{enum.GoIdent}
			}
		}
	default:
		return []any{lkpScalarFieldTyp(fd.Kind())}
	}
	panic("should not happen")
}
