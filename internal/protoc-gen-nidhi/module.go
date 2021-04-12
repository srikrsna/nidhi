// Package pgn ...
package pgn

import (
	"embed"
	"io/fs"
	"io/ioutil"
	"strconv"
	"strings"
	"text/template"

	"github.com/gertd/go-pluralize"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"

	nidhipb "github.com/srikrsna/nidhi/nidhi"
)

type Module struct {
	*pgs.ModuleBase

	goContext pgsgo.Context
	tpl       *template.Template
}

func New() pgs.Module { return &Module{ModuleBase: &pgs.ModuleBase{}} }

func (m *Module) Name() string {
	return "nidhi"
}

func (m *Module) InitContext(c pgs.BuildContext) {
	m.ModuleBase.InitContext(c)
	m.goContext = pgsgo.InitContext(c.Parameters())

	pc := pluralize.NewClient()

	tpl, err := parseTemplates(map[string]interface{}{
		"Name":        m.goContext.Name,
		"ImportPath":  m.goContext.ImportPath,
		"PackageName": m.goContext.PackageName,
		"LowerCamel":  func(n pgs.Name) pgs.Name { return n.LowerCamelCase() },
		"Plural":      func(n pgs.Name) string { return pc.Plural(n.String()) },
		"IsString": func(f pgs.Field) bool {
			return !f.Type().IsRepeated() && f.Type().ProtoType() == pgs.StringT
		},
		"IsBool": func(f pgs.Field) bool {
			return !f.Type().IsRepeated() && f.Type().ProtoType() == pgs.BoolT
		},
		"IsBytes": func(f pgs.Field) bool {
			return !f.Type().IsRepeated() && f.Type().ProtoType() == pgs.BytesT
		},
		"IsWKT": func(f pgs.Field) bool {
			return f.Type().IsEmbed() && (f.Type().Embed().IsWellKnown() || f.Type().Embed().Name() == "FieldMask")
		},
		"IsWKTTime": func(f pgs.Field) bool {
			return f.Type().Embed().WellKnownType() == pgs.TimestampWKT
		},
		"IsWKTAny": func(f pgs.Field) bool {
			return f.Type().Embed().WellKnownType() == pgs.AnyWKT
		},
		"IsWKTFieldMask": func(f pgs.Field) bool {
			return f.Type().Embed().Name() == "FieldMask"
		},
		"IsWKTDuration": func(f pgs.Field) bool {
			return f.Type().Embed().WellKnownType() == pgs.DurationWKT
		},
		"IsEnum": func(f pgs.Field) bool {
			return f.Type().IsEnum()
		},
		"Fields": func(msg pgs.Message) string {
			var ff string
			for _, f := range msg.Fields() {
				ff += strconv.Quote(f.Name().LowerCamelCase().String()) + ", "
			}

			return ff
		},
		"OneOfOption": func(f pgs.Field) pgs.Name {
			return m.goContext.OneofOption(f)
		},
		"GoType": func(f pgs.Field) pgsgo.TypeName {
			return m.goContext.Type(f)
		},
		"Capitalise": func(name pgsgo.TypeName) string {
			return pgs.Name(name).UpperCamelCase().String()
		},
	})
	if err != nil {
		m.Fail(err)
	}

	m.tpl = tpl
}

func (m *Module) Execute(files map[string]pgs.File, _ map[string]pgs.Package) []pgs.Artifact {
	type Root struct {
		pgs.Message
	}

	for _, file := range files {
		if len(file.AllMessages()) == 0 {
			continue
		}

		roots := make([]Root, 0, len(file.AllMessages()))
		name := m.goContext.OutputPath(file).SetExt(".nidhi.go").String()

		imp := map[string]bool{}
		fileWrapper := struct {
			pgs.File
			Imps []string
		}{
			file,
			[]string{},
		}
		m.AddGeneratorTemplateFile(name, m.tpl.Lookup("header"), &fileWrapper)

		for _, msg := range file.Messages() {
			m.generateMarshaler(msg)

			for _, f := range msg.Fields() {
				if f.Type().IsEmbed() {
					if !f.Type().Embed().BuildTarget() && ((!f.Type().Embed().IsWellKnown() && f.Type().Embed().Name() != "FieldMask") || f.Type().Embed().WellKnownType() == pgs.AnyWKT) {
						imp[string(m.goContext.ImportPath(f.Type().Embed()))] = true
					}
				} else if (f.Type().IsRepeated() || f.Type().IsMap()) && f.Type().Element().IsEmbed() {
					if !f.Type().Element().Embed().BuildTarget() {
						imp[string(m.goContext.ImportPath(f.Type().Element().Embed()))] = true
					}
				}
			}

			var prefix bool
			found, err := msg.Extension(nidhipb.E_Root, &prefix)
			if err != nil {
				m.Fail(err)
			}

			if found {
				roots = append(roots, Root{msg})
			}
		}

		fileWrapper.Imps = make([]string, 0, len(imp)+3)
		for i := range imp {
			fileWrapper.Imps = append(fileWrapper.Imps, i)
		}

		if len(roots) == 0 {
			continue
		}

		fileWrapper.Imps = append(fileWrapper.Imps, "database/sql", "context", "github.com/srikrsna/nidhi")

		for _, root := range roots {
			m.AddGeneratorTemplateAppend(name, m.tpl.Lookup("fields-header"), root)
			m.AddGeneratorTemplateAppend(name, m.tpl.Lookup("store"), root)
			m.AddGeneratorTemplateAppend(name, m.tpl.Lookup("query"), root)
			m.generateSubQuery(name, root, root, root.Name())
		}
	}

	return m.Artifacts()
}

func (m *Module) shouldAddImport(msg pgs.Message) bool {
	return !msg.BuildTarget() && !msg.IsWellKnown() && msg.Name() != "FieldMask"
}

func (m *Module) generateSubQuery(name string, root, msg pgs.Message, parent pgs.Name) {
	type SubRoot struct {
		Root pgs.Message
		pgs.Field
		Parent pgs.Name
	}

	for _, f := range msg.Fields() {
		if !f.Type().IsEmbed() || f.Type().Embed().IsWellKnown() {
			continue
		}

		m.AddGeneratorTemplateAppend(name, m.tpl.Lookup("sub-query"), SubRoot{Root: root, Parent: parent, Field: f})

		m.generateSubQuery(name, root, f.Type().Embed(), parent.UpperCamelCase()+f.Name().UpperCamelCase())
	}
}

func (m *Module) generateMarshaler(msg pgs.Message) {
	if !msg.BuildTarget() {
		return
	}

	name := m.goContext.OutputPath(msg).SetExt(".nidhi.go").String()
	m.AddGeneratorTemplateAppend(name, m.tpl.Lookup("fields"), msg)
	m.AddGeneratorTemplateAppend(name, m.tpl.Lookup("json"), msg)
}

//go:embed templates/*.tmpl
var tfs embed.FS

func parseTemplates(fm template.FuncMap) (*template.Template, error) {
	tpl := template.New("nidhi").Funcs(fm)

	if err := fs.WalkDir(tfs, ".", func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := tfs.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		content, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}

		if _, err := tpl.New(strings.TrimSuffix(info.Name(), ".tmpl")).Parse(string(content)); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return tpl, nil
}
