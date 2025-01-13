package cmd

import (
	"errors"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Module is a command to create a new module
func Module() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "module:create [name...]",
		Short: "Create a new module",
		Run: func(cmd *cobra.Command, args []string) {
			tmplModuleStr := `package {{.Name}}

type {{.PName}}Module struct {
	Ctl *{{.PName}}Controller
	Svc *{{.PName}}Service
}

func New() *{{.PName}}Module {
	svc := newService()
	return &{{.PName}}Module{
		Ctl: newController(svc),
		Svc: svc,
	}
}
`
			tmplServiceStr := `package {{.Name}}

type {{.PName}}Service struct {}

func newService() *{{.PName}}Service {
	return &{{.PName}}Service{}
}
`
			tmplControllerStr := `package {{.Name}}

import (
	"app/app/modules/base"

	"github.com/gin-gonic/gin"
)

type {{.PName}}Controller struct {
	{{.Name}}Svc *{{.PName}}Svc
}

func newController({{.Name}}Service *{{.PName}}Service) *{{.PName}}Controller {
	return &{{.PName}}Controller{
		{{.Name}}Svc: {{.Name}}Service,
	}
}

func (*{{.PName}}Controller) Get(ctx *gin.Context) {
	base.Success(ctx, "ok")
}
`
			tmplDTOStr := `package {{.Name}}dto`
			tmplENTStr := `package {{.Name}}ent`

			tmplModule := must(template.New("module").Parse(tmplModuleStr))

			tmplSVC := must(template.New("service").Parse(tmplServiceStr))

			tmplCTL := must(template.New("controller").Parse(tmplControllerStr))

			tmplDTO := must(template.New("dto").Parse(tmplDTOStr))
			tmplENT := must(template.New("ent").Parse(tmplENTStr))

			for _, arg := range args {

				name := cases.Lower(language.English).String(arg)
				pName := cases.Title(language.English).String(arg)

				data := map[string]string{
					"Name":  name,
					"PName": pName,
				}

				mod := strings.Builder{}
				svc := strings.Builder{}
				ctl := strings.Builder{}
				dto := strings.Builder{}
				ent := strings.Builder{}

				panicErr(errors.Join(
					tmplModule.Execute(&mod, data),
					tmplSVC.Execute(&svc, data),
					tmplCTL.Execute(&ctl, data),
					tmplDTO.Execute(&dto, data),
					tmplENT.Execute(&ent, data),

					os.MkdirAll("./app/modules/"+name+"/dto", 0755),
					os.MkdirAll("./app/modules/"+name+"/ent", 0755),
				))

				modFile := must(os.Create("./app/modules/" + name + "/" + name + ".mod.go"))

				svcFile := must(os.Create("./app/modules/" + name + "/" + name + ".svc.go"))

				ctlFile := must(os.Create("./app/modules/" + name + "/" + name + ".ctl.go"))

				dtoFile := must(os.Create("./app/modules/" + name + "/dto/" + name + ".dto.go"))
				entFile := must(os.Create("./app/modules/" + name + "/ent/" + name + ".ent.go"))

				panicErr(errors.Join(
					getError(modFile.WriteString(mod.String())),
					getError(svcFile.WriteString(svc.String())),
					getError(ctlFile.WriteString(ctl.String())),
					getError(dtoFile.WriteString(dto.String())),
					getError(entFile.WriteString(ent.String())),

					modFile.Close(),
					svcFile.Close(),
					ctlFile.Close(),
					dtoFile.Close(),
					entFile.Close(),
				))

			}
		},
	}
	return cmd
}

func must[T any](t T, err error) T {
	panicErr(err)
	return t
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getError(_ any, err error) error {
	return err
}
