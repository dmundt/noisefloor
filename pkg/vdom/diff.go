package vdom

func diffElementTreesRecursive(old *Element, new *Element) []Patch {
	patchList := []Patch{}

	// full patch from this node if any changes in children
	// TODO - this should check for any changes, not just total count
	if len(new.Children) != len(old.Children) {
		patch := Patch{Type: Replace, Path: old.Path, Element: *new}
		patchList = append(patchList, patch)
		return patchList
	}

	if new.Type == Text {
		// add patche for text element changed
		if new.Attrs["Text"] != old.Attrs["Text"] {
			patch := Patch{Type: TextSet, Path: old.Path[:len(old.Path)-1], Attr: Attr{Name: "Text", Value: new.Attrs["Text"]}}
			patchList = append(patchList, patch)
		}
		return patchList
	}

	// add patches for changed (or new) attributes
	oldAttrs := old.Attrs
	for key := range new.Attrs {
		value, present := old.Attrs[key]
		if !present || value != new.Attrs[key] {
			setType := AttrSet
			if key == "value" {
				setType = ValueSet
			}
			patch := Patch{Type: setType, Path: old.Path, Attr: Attr{Name: key, Value: new.Attrs[key]}}
			patchList = append(patchList, patch)
		}
		delete(oldAttrs, key)
	}
	// add patches for removed attributes
	for key := range oldAttrs {
		patch := Patch{Type: AttrRemove, Path: old.Path, Attr: Attr{Name: key}}
		patchList = append(patchList, patch)
	}

	// add patches for children
	for i := 0; i < len(new.Children); i++ {
		patchList = append(patchList, diffElementTreesRecursive(&old.Children[i], &new.Children[i])...)
	}

	return patchList
}

func diffElementTrees(old *Element, new *Element) PatchList {
	patch := diffElementTreesRecursive(old, new)
	patchList := PatchList{SVGNamespace: svgNamespace, Patch: patch}
	// patchList := PatchList{SVGNamespace: svgNamespace, Patch: []Patch{Patch{Type: Replace, Path: []int{}, Element: *new}}}

	return patchList
}
