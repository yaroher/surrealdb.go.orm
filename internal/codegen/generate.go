package codegen

import "fmt"

func Generate(opts Options) error {
	if len(opts.Dirs) == 0 {
		return fmt.Errorf("no directories provided")
	}
	for _, dir := range opts.Dirs {
		pkg, err := ParseDir(dir)
		if err != nil {
			return err
		}
		if len(pkg.Models) == 0 {
			continue
		}
		if err := Render(dir, pkg); err != nil {
			return err
		}
	}
	return nil
}
