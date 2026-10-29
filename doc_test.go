package backlog_test

import (
	"go/ast"
	"go/doc"
	"go/doc/comment"
	"go/parser"
	"go/token"
	"io/fs"
	"strings"
	"testing"
)

// docComment is a doc comment of the root package, with a label that names
// where it belongs.
type docComment struct {
	where string
	text  string
}

// loadPackageDoc parses the non-test source files of the root package.
func loadPackageDoc(t *testing.T) *doc.Package {
	t.Helper()

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, ".", func(fi fs.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	pkg, ok := pkgs["backlog"]
	if !ok {
		t.Fatal("package backlog not found")
	}
	files := make([]*ast.File, 0, len(pkg.Files))
	for _, f := range pkg.Files {
		files = append(files, f)
	}

	p, err := doc.NewFromFiles(fset, files, "github.com/nattokin/go-backlog", doc.PreserveAST)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

// docComments returns every doc comment of the package: the package comment,
// constants, variables, functions, types, methods, and the fields of struct
// and interface types.
func docComments(p *doc.Package) []docComment {
	var out []docComment
	add := func(where, text string) {
		if text != "" {
			out = append(out, docComment{where: where, text: text})
		}
	}
	addValues := func(values []*doc.Value) {
		for _, v := range values {
			add(strings.Join(v.Names, ", "), v.Doc)
		}
	}

	add("package", p.Doc)
	addValues(p.Consts)
	addValues(p.Vars)
	for _, f := range p.Funcs {
		add(f.Name, f.Doc)
	}
	for _, typ := range p.Types {
		add(typ.Name, typ.Doc)
		addValues(typ.Consts)
		addValues(typ.Vars)
		for _, f := range typ.Funcs {
			add(f.Name, f.Doc)
		}
		for _, m := range typ.Methods {
			add(typ.Name+"."+m.Name, m.Doc)
		}
		for _, spec := range typ.Decl.Specs {
			var fields *ast.FieldList
			switch tt := spec.(*ast.TypeSpec).Type.(type) {
			case *ast.StructType:
				fields = tt.Fields
			case *ast.InterfaceType:
				fields = tt.Methods
			}
			if fields == nil {
				continue
			}
			for _, f := range fields.List {
				if f.Doc == nil {
					continue
				}
				for _, name := range f.Names {
					add(typ.Name+"."+name.Name, f.Doc.Text())
				}
			}
		}
	}
	return out
}

// docLinks returns the doc links in d.
func docLinks(d *comment.Doc) []*comment.DocLink {
	var out []*comment.DocLink
	var walkText func([]comment.Text)
	walkText = func(texts []comment.Text) {
		for _, t := range texts {
			switch t := t.(type) {
			case *comment.DocLink:
				out = append(out, t)
			case *comment.Link:
				walkText(t.Text)
			}
		}
	}
	var walkBlock func(comment.Block)
	walkBlock = func(b comment.Block) {
		switch b := b.(type) {
		case *comment.Paragraph:
			walkText(b.Text)
		case *comment.Heading:
			walkText(b.Text)
		case *comment.List:
			for _, item := range b.Items {
				for _, c := range item.Content {
					walkBlock(c)
				}
			}
		}
	}
	for _, b := range d.Content {
		walkBlock(b)
	}
	return out
}

// linkName formats l as it is written in a comment, such as "IssueService.List"
// or "errors.As".
func linkName(l *comment.DocLink) string {
	name := l.Name
	if l.Recv != "" {
		name = l.Recv + "." + name
	}
	if l.ImportPath != "" {
		name = l.ImportPath + "." + name
	}
	return name
}

// TestDocLinksResolve checks that every doc link in the package resolves the
// way go doc and pkg.go.dev resolve it. A link that does not resolve is shown
// as plain text with its brackets, and no other tool reports it.
func TestDocLinksResolve(t *testing.T) {
	p := loadPackageDoc(t)

	// strict resolves links the same way as go doc and pkg.go.dev.
	strict := p.Parser()
	// loose accepts every [X] and [X.Y] as a link.
	loose := &comment.Parser{
		LookupPackage: func(name string) (string, bool) { return name, true },
		LookupSym:     func(recv, name string) bool { return true },
	}

	for _, dc := range docComments(p) {
		resolved := map[string]bool{}
		for _, l := range docLinks(strict.Parse(dc.text)) {
			resolved[linkName(l)] = true
		}
		for _, l := range docLinks(loose.Parse(dc.text)) {
			if name := linkName(l); !resolved[name] {
				t.Errorf("%s: doc link [%s] does not resolve", dc.where, name)
			}
		}
	}
}

// TestDocOptionListsMatchService checks that the options a method lists in its
// doc comment belong to the option service of its receiver, that is, the type
// of the receiver's Option field.
func TestDocOptionListsMatchService(t *testing.T) {
	p := loadPackageDoc(t)
	strict := p.Parser()

	for _, typ := range p.Types {
		optionType := optionFieldType(typ)
		if optionType == "" {
			continue
		}
		for _, m := range typ.Methods {
			for _, l := range docLinks(strict.Parse(m.Doc)) {
				// A link to an option service itself ([XOptionService]) or to
				// one of its methods ([XOptionService.WithY]).
				service := l.Recv
				if service == "" {
					service = l.Name
				}
				if l.ImportPath != "" || !strings.HasSuffix(service, "OptionService") {
					continue
				}
				if service != optionType {
					t.Errorf("%s.%s: doc comment refers to [%s], want options of %s",
						typ.Name, m.Name, linkName(l), optionType)
				}
			}
		}
	}
}

// optionFieldType returns the type name of the Option field of typ, or "" if
// typ has no such field.
func optionFieldType(typ *doc.Type) string {
	for _, spec := range typ.Decl.Specs {
		st, ok := spec.(*ast.TypeSpec).Type.(*ast.StructType)
		if !ok {
			continue
		}
		for _, f := range st.Fields.List {
			for _, name := range f.Names {
				if name.Name != "Option" {
					continue
				}
				if star, ok := f.Type.(*ast.StarExpr); ok {
					if ident, ok := star.X.(*ast.Ident); ok {
						return ident.Name
					}
				}
			}
		}
	}
	return ""
}
