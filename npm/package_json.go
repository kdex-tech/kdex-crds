package npm

import (
	"fmt"
	"strings"
)

func (p *PackageJSON) hasESModule() error {
	if p.Browser != "" {
		return nil
	}

	if p.Type == "module" {
		return nil
	}

	if p.Module != "" {
		return nil
	}

	if p.Exports != nil {
		if strings.HasSuffix(p.Exports.Single, ".mjs") {
			return nil
		}

		if p.Exports.Multiple != nil {
			_, ok := p.Exports.Multiple["browser"]

			if ok {
				return nil
			}

			_, ok = p.Exports.Multiple["import"]

			if ok {
				return nil
			}
		}
	}

	if strings.HasSuffix(p.Main, ".mjs") {
		return nil
	}

	return fmt.Errorf("package does not contain an ES module")
}
