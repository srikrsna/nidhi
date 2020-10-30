package nidhi

import (
	"text/template"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"github.com/markbates/pkger"
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

	tpl := template.New("nidhi").Funcs(map[string]interface{}{})

	f, err := pkger.Open("/pkg/protoc-gen-nidhi/templates/header.tmpl")
	if err != nil {
		m.Fail(err)
	}
	defer f.Close()

	m.tpl = template.Must(tpl.Parse(``))
}

func (m *Module) Execute(files map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {
	for _, file := range files {
		type Root struct {
			pgs.Message
			Prefix string
		}

		roots := make([]Root, 0, len(file.AllMessages()))
		for _, msg := range file.Messages() {
			var prefix string
			found, err := msg.Extension(nidhipb.E_Prefix, &prefix)
			if err != nil {
				m.Fail(err)
			}

			if found {
				roots = append(roots, Root{msg, prefix})
			}
		}

		if len(roots) == 0 {
			continue
		}

		// Found some
		name := m.goContext.OutputPath(file).SetExt(".nidhi.go")
		m.AddGeneratorTemplateFile(name.String(), nil, file)
		m.goContext.Name(file.Package())
		for _, root := range roots {
			m.AddGeneratorTemplateAppend(name.String(), nil, root)
		}
	}

	return m.Artifacts()
}
