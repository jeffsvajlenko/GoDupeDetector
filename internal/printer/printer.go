package printer

import (
	"GoDupeDetector/internal/detection"
	"GoDupeDetector/internal/parsing"
	"encoding/json"
	"io"
)

type cloneReport struct {
	Clones []clone `json:"clones"`
}

type clone struct {
	Function1 function `json:"function1"`
	Function2 function `json:"function2"`
}

type function struct {
	Path        string `json:"path"`
	StartLine   uint   `json:"start_line"`
	StartColumn uint   `json:"start_column"`
	EndLine     uint   `json:"end_line"`
	EndColumn   uint   `json:"end_column"`
	Body        string `json:"body"`
}

func PrintCloneReport(pset *parsing.ParseSet, cset *detection.CloneSet, output io.Writer) error {
	clones := make([]clone, 0, len(cset.Clones))
	for _, c := range cset.Clones {
		fu1 := pset.Functions[c.FunctionId1]
		f1 := pset.Files[fu1.FileId]
		fu2 := pset.Functions[c.FunctionId2]
		f2 := pset.Files[fu2.FileId]

		clones = append(clones, clone{
			Function1: function{
				Path:        f1.Path,
				StartLine:   fu1.StartLine,
				StartColumn: fu1.StartColumn,
				EndLine:     fu1.EndLine,
				EndColumn:   fu1.EndColumn,
				Body:        fu1.Body,
			},
			Function2: function{
				Path:        f2.Path,
				StartLine:   fu2.StartLine,
				StartColumn: fu2.StartColumn,
				EndLine:     fu2.EndLine,
				EndColumn:   fu2.EndColumn,
				Body:        fu2.Body,
			},
		})

	}

	report := cloneReport{
		Clones: clones,
	}

	data, err := json.MarshalIndent(report, "", "    ")
	if err != nil {
		return err
	}

	_, err = output.Write(data)
	if err != nil {
		return err
	}

	return nil
}
