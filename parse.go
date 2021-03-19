package structural

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"strings"
	"text/template"
)

var walkCount = 0

func newAstTool(f *token.FileSet) func(n ast.Node) astTool {
	return func(n ast.Node) astTool { return astTool{n, f, bytes.NewBuffer(nil)} }
}

type writerStringer interface {
	io.Writer
	fmt.Stringer
}

type astTool struct {
	ast.Node
	*token.FileSet
	writerStringer
}

func (a astTool) String() string {
	if err := printer.Fprint(a.writerStringer, a.FileSet, a.Node); err != nil {
		panic(err.Error())
	}
	return a.writerStringer.String()
}

func (a astTool) Matches(test string) bool {
	return a.String() == test
}

func Generate(w io.Writer, input string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, input, nil, parser.AllErrors)
	if err != nil {
		return err
	}
	apis, err := readDefinition(node, fset)
	if err != nil {
		return err
	}
	var api *APIData
	for _, v := range apis {
		api = v
	}
	return template.Must(template.ParseGlob("templates/**")).ExecuteTemplate(w, "main", api)
}

type (
	QueryParamData struct {
		Key                 string
		Value               string
		CoerceTo            string
		CoercedVariableName string
	}
	RouteData struct {
		Handler              string
		Type                 string
		Path                 string
		RequestBodyType      string
		RequestBodyEncoding  string
		ResponseBodyType     string
		ResponseBodyEncoding string
		QueryParamData       []QueryParamData
		QueryParamArguments  []string
	}

	APIData struct {
		Name      string
		Routes    []*RouteData
		Encodings map[string]struct{}
		Package   string
	}
)

func (a *APIData) findRouteData(routeType string) *RouteData {
	for idx, route := range a.Routes {
		if route.Type == routeType {
			return a.Routes[idx]
		}
	}
	return nil
}

type Pass uint

const (
	APIPass Pass = iota + 1
	HandlersPass
)

type CollectTypes struct {
	rootNode         ast.Node
	pass             Pass
	structuralImport string
	handlerTypes     []ast.Node
	apiTypes         []ast.Node
	fset             *token.FileSet
	apis             map[string]*APIData
	pkg              string
	astTool          func(ast.Node) astTool
}

func (c *CollectTypes) isType(node *ast.TypeSpec, test string) bool {
	if it, ok := node.Type.(*ast.StructType); ok {
		for _, f := range it.Fields.List {
			if c.astTool(f.Type).Matches(test) {
				return true
			}
		}
	}
	return false
}

func (c *CollectTypes) handleAPIType(node ast.Node, typeSpec *ast.TypeSpec, structType *ast.StructType) {
	name := typeSpec.Name.Name
	if _, exist := c.apis[name]; !exist {
		c.apis[name] = &APIData{Name: name, Encodings: map[string]struct{}{}, Package: c.pkg}
		for _, field := range structType.Fields.List {
			if field.Tag != nil {
				path := strings.Replace(field.Tag.Value, "`", "", -1)
				path = strings.Replace(path, "route:", "", -1)
				c.apis[name].Routes = append(c.apis[name].Routes, &RouteData{
					Path:    path,
					Handler: field.Names[0].Name,
					Type:    c.astTool(field.Type).String(),
				})
			}
		}
	}
}

func (c *CollectTypes) handleRouteType(node ast.Node, typeSpec *ast.TypeSpec, structType *ast.StructType) {
	var api *APIData
	for _, api = range c.apis {
		var route *RouteData
		if route = api.findRouteData(typeSpec.Name.Name); route == nil {
			// This ain't a route in the current API
			continue
		}
		for _, field := range structType.Fields.List {
			if c.astTool(field.Type).String() != "Route" {
				if len(field.Names) > 0 {
					name := field.Names[0].Name
					switch name {
					case "RequestBody":
						route.RequestBodyType = c.astTool(field.Type).String()
					case "Response":
						route.ResponseBodyType = c.astTool(field.Type).String()
					case "RequestBodyEncoding":
						route.RequestBodyEncoding = strings.Replace(
							strings.ToLower(c.astTool(field.Type).String()),
							c.structuralImport+".", "", -1)
						api.Encodings[route.RequestBodyEncoding] = struct{}{}
					case "ResponseBodyEncoding":
						route.ResponseBodyEncoding = strings.Replace(
							strings.ToLower(c.astTool(field.Type).String()),
							c.structuralImport+".", "", -1)
						api.Encodings[route.ResponseBodyEncoding] = struct{}{}
					}
					if field.Tag != nil {
						fieldType := c.astTool(field.Type).String()
						fieldName := field.Names[0].Name
						fieldTag := field.Tag.Value
						route.QueryParamData = append(route.QueryParamData, QueryParamData{
							Key:                 strings.Replace(strings.Replace(fieldTag, "`arg:\"", "", -1), "\"`", "", -1),
							Value:               fieldName,
							CoerceTo:            fieldType,
							CoercedVariableName: fmt.Sprintf("_arg%d", walkCount),
						})
						route.QueryParamArguments = append(route.QueryParamArguments, fmt.Sprintf("_arg%d", walkCount))
					}
				}
			}
		}
	}
}

func (c *CollectTypes) route() string {
	return fmt.Sprintf("%s.Route", c.structuralImport)
}
func (c *CollectTypes) api() string {
	return fmt.Sprintf("%s.API", c.structuralImport)
}

func (c *CollectTypes) apiPass(node ast.Node) {
	if it, ok := node.(*ast.TypeSpec); ok {
		if c.isType(it, c.api()) && c.pass == APIPass {
			c.handleAPIType(c.rootNode, it, it.Type.(*ast.StructType))
		}
	}
}

func (c *CollectTypes) handlerPass(node ast.Node) {
	if it, ok := node.(*ast.TypeSpec); ok {
		if c.isType(it, c.route()) {
			c.handleRouteType(c.rootNode, it, it.Type.(*ast.StructType))
		}
	}
}

func (c *CollectTypes) setStructuralImportPath(node ast.Node) {
	var (
		importSpec *ast.ImportSpec
		ok         bool
	)
	if importSpec, ok = node.(*ast.ImportSpec); !ok {
		return
	}
	if importSpec.Name != nil {
		c.structuralImport = importSpec.Name.Name
		return
	}
	importParts := strings.Split(importSpec.Path.Value, "/")
	structuralImport := importParts[len(importParts)-1]
	c.structuralImport = strings.Replace(structuralImport, `"`, "", -1)
}

func (c *CollectTypes) setPackage(node ast.Node) {
	if it, ok := node.(*ast.File); ok {
		c.pkg = it.Name.Name
	}
}

func (c *CollectTypes) Visit(node ast.Node) ast.Visitor {
	walkCount++
	c.setPackage(node)
	c.setStructuralImportPath(node)
	switch c.pass {
	case APIPass:
		c.apiPass(node)
	case HandlersPass:
		c.handlerPass(node)
	default:
		return nil
	}
	return c
}

func readDefinition(node ast.Node, fset *token.FileSet) (map[string]*APIData, error) {
	c := &CollectTypes{
		rootNode:     node,
		apiTypes:     []ast.Node{},
		handlerTypes: []ast.Node{},
		apis:         map[string]*APIData{},
		fset:         fset,
		astTool:      newAstTool(fset),
	}

	for _, pass := range []Pass{APIPass, HandlersPass} {
		c.pass = pass
		ast.Walk(c, node)
	}
	return c.apis, nil
}
