package parsing

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
	"sync"
)

type ParseSet struct {
	FileSet   *token.FileSet
	Files     map[uint]File
	Functions map[uint]Function
	Strings   sync.Map
}

type Param struct {
	Name string
	Type string
}

type File struct {
	Id   uint
	Path string
}

type Function struct {
	Id              uint
	FileId          uint
	StartLine       uint
	StartColumn     uint
	EndLine         uint
	EndColumn       uint
	LineLength      uint
	Name            string
	Params          []Param
	Results         []Param
	Receiver        *Param
	Body            string
	PrettyPrintBody []string
}

func FileList(dir string) ([]string, error) {
	retval := make([]string, 0, 100)
	err := filepath.Walk(dir,
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Ext(path) == ".go" {
				retval = append(retval, path)
			}
			return nil
		})
	return retval, err
}

func Parse(files []string) (*ParseSet, error) {
	parseSet := ParseSet{
		FileSet:   token.NewFileSet(),
		Files:     make(map[uint]File),
		Functions: make(map[uint]Function),
		Strings:   sync.Map{},
	}

	var wg sync.WaitGroup
	wg.Add(len(files))

	cfunc := make(chan Function, 100)

	go func() {
		wg.Wait()
		close(cfunc)
	}()

	for i, file := range files {
		file := File{
			Id:   uint(i + 1),
			Path: file,
		}

		parseSet.Files[file.Id] = file

		go func(file File, cfunc chan Function) {
			defer wg.Done()
			defer func() {
				if perr := recover(); perr != nil {
					log.Println("parsing failed for file", file.Path)
				}
			}()

			functions, errs := parseFile(file, parseSet.FileSet, &parseSet.Strings)
			if len(errs) > 0 {
				log.Println("prasing failed for file", file)
			} else {
				for _, function := range functions {
					cfunc <- function
				}
			}
		}(file, cfunc)
	}

	functionId := uint(1)
	for function := range cfunc {
		function.Id = functionId
		parseSet.Functions[functionId] = function
		functionId = functionId + 1
	}

	return &parseSet, nil
}

func parseFile(file File, fset *token.FileSet, intern *sync.Map) ([]Function, []error) {
	retval := make([]Function, 0, 10)

	errors := make([]error, 0, 1)

	f, err := parser.ParseFile(fset, file.Path, nil, 0)
	if err != nil {
		errors = append(errors, err)
		return nil, errors
	}

	for _, d := range f.Decls {
		switch d := d.(type) {
		case *ast.FuncDecl:
			startPosition := fset.PositionFor(d.Pos(), true)
			endPosition := fset.PositionFor(d.End(), true)

			params, err := getFunctionParams(d, fset)
			if err != nil {
				errors = append(errors, err)
				continue
			}

			results, err := getFunctionReturns(d, fset)
			if err != nil {
				errors = append(errors, err)
				continue
			}

			receiver, err := getFunctionReceiver(d, fset)
			if err != nil {
				errors = append(errors, err)
				continue
			}

			body, err := getBody(d, fset)
			if err != nil {
				errors = append(errors, err)
				continue
			}

			ppbody, err := getPrettyPrintBody(d, fset, intern)
			if err != nil {
				errors = append(errors, err)
				continue
			}

			if len(ppbody) == 0 {
				continue
			}

			retval = append(retval, Function{
				Id:              uint(0),
				FileId:          file.Id,
				StartLine:       uint(startPosition.Line),
				StartColumn:     uint(startPosition.Column),
				EndLine:         uint(endPosition.Line),
				EndColumn:       uint(endPosition.Column),
				LineLength:      uint(endPosition.Line - startPosition.Line + 1),
				Name:            d.Name.Name,
				Params:          params,
				Results:         results,
				Receiver:        receiver,
				Body:            body,
				PrettyPrintBody: ppbody,
			})
		}
	}
	return retval, errors
}

func getFunctionParams(f *ast.FuncDecl, fset *token.FileSet) ([]Param, error) {
	if f.Type.Params.List != nil {
		return extractParams(f.Type.Params, fset)
	} else {
		return nil, nil
	}
}

func getFunctionReceiver(f *ast.FuncDecl, fset *token.FileSet) (*Param, error) {
	if f.Recv != nil && len(f.Recv.List) > 0 {
		params, err := extractParams(f.Recv, fset)
		return &params[0], err
	} else {
		return nil, nil
	}
}

func getFunctionReturns(f *ast.FuncDecl, fset *token.FileSet) ([]Param, error) {
	if f.Type.Results != nil {
		return extractParams(f.Type.Results, fset)
	} else {
		return nil, nil
	}
}

func extractParams(f *ast.FieldList, fset *token.FileSet) ([]Param, error) {
	params := make([]Param, 0, len(f.List))
	for _, param := range f.List {
		var names bytes.Buffer
		var types bytes.Buffer

		if len(param.Names) > 0 {
			printer.Fprint(&names, fset, param.Names[0])
		}
		err := printer.Fprint(&types, fset, param.Type)
		if err != nil {
			return nil, err
		}

		params = append(params, Param{
			Name: names.String(),
			Type: types.String(),
		})
	}
	return params, nil
}

func getBody(f *ast.FuncDecl, fset *token.FileSet) (string, error) {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, fset, f.Body)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func getPrettyPrintBody(f *ast.FuncDecl, fset *token.FileSet, intern *sync.Map) ([]string, error) {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, fset, f.Body)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.ReplaceAll(buf.String(), "\r\n", "\n"), "\n")
	plines := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			line, _ := intern.LoadOrStore(line, line)
			plines = append(plines, line.(string))
		}
	}
	return plines, nil
}
