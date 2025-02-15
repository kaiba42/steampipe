package modconfig

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/turbot/go-kit/helpers"
)

// ModResources is a struct containing maps of all mod resource types
// This is provided to avoid db needing to reference workspace package
type ModResources struct {
	// the parent mod
	Mod *Mod
	// all mods (including deps)
	Mods                  map[string]*Mod
	Queries               map[string]*Query
	Controls              map[string]*Control
	Benchmarks            map[string]*Benchmark
	Variables             map[string]*Variable
	Dashboards            map[string]*Dashboard
	DashboardContainers   map[string]*DashboardContainer
	DashboardCards        map[string]*DashboardCard
	DashboardCharts       map[string]*DashboardChart
	DashboardFlows        map[string]*DashboardFlow
	DashboardHierarchies  map[string]*DashboardHierarchy
	DashboardImages       map[string]*DashboardImage
	DashboardInputs       map[string]map[string]*DashboardInput
	GlobalDashboardInputs map[string]*DashboardInput
	DashboardTables       map[string]*DashboardTable
	DashboardTexts        map[string]*DashboardText
	References            map[string]*ResourceReference

	Locals          map[string]*Local
	LocalQueries    map[string]*Query
	LocalControls   map[string]*Control
	LocalBenchmarks map[string]*Benchmark
}

func NewWorkspaceResourceMaps(mod *Mod) *ModResources {
	return &ModResources{
		Mod:                   mod,
		Mods:                  map[string]*Mod{mod.Name(): mod},
		Queries:               make(map[string]*Query),
		Controls:              make(map[string]*Control),
		Benchmarks:            make(map[string]*Benchmark),
		Locals:                make(map[string]*Local),
		Variables:             make(map[string]*Variable),
		Dashboards:            make(map[string]*Dashboard),
		DashboardContainers:   make(map[string]*DashboardContainer),
		DashboardCards:        make(map[string]*DashboardCard),
		DashboardCharts:       make(map[string]*DashboardChart),
		DashboardFlows:        make(map[string]*DashboardFlow),
		DashboardHierarchies:  make(map[string]*DashboardHierarchy),
		DashboardImages:       make(map[string]*DashboardImage),
		DashboardInputs:       make(map[string]map[string]*DashboardInput),
		GlobalDashboardInputs: make(map[string]*DashboardInput),
		DashboardTables:       make(map[string]*DashboardTable),
		DashboardTexts:        make(map[string]*DashboardText),
		References:            make(map[string]*ResourceReference),
		LocalQueries:          make(map[string]*Query),
		LocalControls:         make(map[string]*Control),
		LocalBenchmarks:       make(map[string]*Benchmark),
	}
}

func CreateWorkspaceResourceMapForQueries(queryProviders []QueryProvider, mod *Mod) *ModResources {
	res := NewWorkspaceResourceMaps(mod)
	for _, p := range queryProviders {
		res.addControlOrQuery(p)
	}
	return res
}

func (m *ModResources) QueryProviders() []QueryProvider {
	numDashboardInputs := 0
	for _, inputs := range m.DashboardInputs {
		numDashboardInputs += len(inputs)
	}
	res := make([]QueryProvider,
		len(m.Queries)+
			len(m.Controls)+
			len(m.DashboardCards)+
			len(m.DashboardCharts)+
			len(m.DashboardFlows)+
			len(m.DashboardHierarchies)+
			numDashboardInputs+
			len(m.GlobalDashboardInputs)+
			len(m.DashboardTables))

	idx := 0
	for _, p := range m.Queries {
		res[idx] = p
		idx++
	}
	for _, p := range m.Controls {
		res[idx] = p
		idx++
	}
	for _, p := range m.DashboardCards {
		res[idx] = p
		idx++
	}
	for _, p := range m.DashboardCharts {
		res[idx] = p
		idx++
	}
	for _, p := range m.DashboardFlows {
		res[idx] = p
		idx++
	}
	for _, p := range m.DashboardHierarchies {
		res[idx] = p
		idx++
	}
	for _, inputsForDashboard := range m.DashboardInputs {
		for _, p := range inputsForDashboard {
			res[idx] = p
			idx++
		}
	}
	for _, p := range m.GlobalDashboardInputs {
		res[idx] = p
		idx++
	}
	for _, p := range m.DashboardTables {
		res[idx] = p
		idx++
	}
	return res
}

func (m *ModResources) Equals(other *ModResources) bool {
	if other == nil {
		return false
	}

	for name, query := range m.Queries {
		if otherQuery, ok := other.Queries[name]; !ok {
			return false
		} else if !query.Equals(otherQuery) {
			return false
		}
	}
	for name := range other.Queries {
		if _, ok := m.Queries[name]; !ok {
			return false
		}
	}

	for name, control := range m.Controls {
		if otherControl, ok := other.Controls[name]; !ok {
			return false
		} else if !control.Equals(otherControl) {
			return false
		}
	}
	for name := range other.Controls {
		if _, ok := m.Controls[name]; !ok {
			return false
		}
	}

	for name, benchmark := range m.Benchmarks {
		if otherBenchmark, ok := other.Benchmarks[name]; !ok {
			return false
		} else if !benchmark.Equals(otherBenchmark) {
			return false
		}
	}
	for name := range other.Benchmarks {
		if _, ok := m.Benchmarks[name]; !ok {
			return false
		}
	}

	for name, variable := range m.Variables {
		if otherVariable, ok := other.Variables[name]; !ok {
			return false
		} else if !variable.Equals(otherVariable) {
			return false
		}
	}
	for name := range other.Variables {
		if _, ok := m.Variables[name]; !ok {
			return false
		}
	}

	for name, dashboard := range m.Dashboards {
		if otherDashboard, ok := other.Dashboards[name]; !ok {
			return false
		} else if !dashboard.Equals(otherDashboard) {
			return false
		}
	}
	for name := range other.Dashboards {
		if _, ok := m.Dashboards[name]; !ok {
			return false
		}
	}

	for name, container := range m.DashboardContainers {
		if otherContainer, ok := other.DashboardContainers[name]; !ok {
			return false
		} else if !container.Equals(otherContainer) {
			return false
		}
	}
	for name := range other.DashboardContainers {
		if _, ok := m.DashboardContainers[name]; !ok {
			return false
		}
	}

	for name, cards := range m.DashboardCards {
		if otherCard, ok := other.DashboardCards[name]; !ok {
			return false
		} else if !cards.Equals(otherCard) {
			return false
		}
	}
	for name := range other.DashboardCards {
		if _, ok := m.DashboardCards[name]; !ok {
			return false
		}
	}

	for name, charts := range m.DashboardCharts {
		if otherChart, ok := other.DashboardCharts[name]; !ok {
			return false
		} else if !charts.Equals(otherChart) {
			return false
		}
	}
	for name := range other.DashboardCharts {
		if _, ok := m.DashboardCharts[name]; !ok {
			return false
		}
	}

	for name, flows := range m.DashboardFlows {
		if otherFlow, ok := other.DashboardFlows[name]; !ok {
			return false
		} else if !flows.Equals(otherFlow) {
			return false
		}
	}
	for name := range other.DashboardFlows {
		if _, ok := m.DashboardFlows[name]; !ok {
			return false
		}
	}

	for name, hierarchies := range m.DashboardHierarchies {
		if otherHierarchy, ok := other.DashboardHierarchies[name]; !ok {
			return false
		} else if !hierarchies.Equals(otherHierarchy) {
			return false
		}
	}
	for name := range other.DashboardHierarchies {
		if _, ok := m.DashboardHierarchies[name]; !ok {
			return false
		}
	}

	for name, images := range m.DashboardImages {
		if otherImage, ok := other.DashboardImages[name]; !ok {
			return false
		} else if !images.Equals(otherImage) {
			return false
		}
	}
	for name := range other.DashboardImages {
		if _, ok := m.DashboardImages[name]; !ok {
			return false
		}
	}

	for name, input := range m.GlobalDashboardInputs {
		if otherInput, ok := other.GlobalDashboardInputs[name]; !ok {
			return false
		} else if !input.Equals(otherInput) {
			return false
		}
	}
	for name := range other.DashboardInputs {
		if _, ok := m.DashboardInputs[name]; !ok {
			return false
		}
	}

	for dashboardName, inputsForDashboard := range m.DashboardInputs {
		if otherInputsForDashboard, ok := other.DashboardInputs[dashboardName]; !ok {
			return false
		} else {

			for name, input := range inputsForDashboard {
				if otherInput, ok := otherInputsForDashboard[name]; !ok {
					return false
				} else if !input.Equals(otherInput) {
					return false
				}
			}
			for name := range otherInputsForDashboard {
				if _, ok := inputsForDashboard[name]; !ok {
					return false
				}
			}

		}
	}
	for name := range other.DashboardInputs {
		if _, ok := m.DashboardInputs[name]; !ok {
			return false
		}
	}

	for name, tables := range m.DashboardTables {
		if otherTable, ok := other.DashboardTables[name]; !ok {
			return false
		} else if !tables.Equals(otherTable) {
			return false
		}
	}
	for name := range other.DashboardTables {
		if _, ok := m.DashboardTables[name]; !ok {
			return false
		}
	}

	for name, texts := range m.DashboardTexts {
		if otherText, ok := other.DashboardTexts[name]; !ok {
			return false
		} else if !texts.Equals(otherText) {
			return false
		}
	}
	for name := range other.DashboardTexts {
		if _, ok := m.DashboardTexts[name]; !ok {
			return false
		}
	}

	for name, reference := range m.References {
		if otherReference, ok := other.References[name]; !ok {
			return false
		} else if !reference.Equals(otherReference) {
			return false
		}
	}

	for name := range other.References {
		if _, ok := m.References[name]; !ok {
			return false
		}
	}

	for name := range other.Locals {
		if _, ok := m.Locals[name]; !ok {
			return false
		}
	}
	return true
}

func (m *ModResources) PopulateReferences() {
	m.References = make(map[string]*ResourceReference)

	resourceFunc := func(resource HclResource) (bool, error) {
		parsedName, _ := ParseResourceName(resource.Name())
		if helpers.StringSliceContains(ReferenceBlocks, parsedName.ItemType) {
			for _, ref := range resource.GetReferences() {
				m.References[ref.String()] = ref
			}
		}

		// continue walking
		return true, nil
	}
	m.WalkResources(resourceFunc)
}

func (m *ModResources) Empty() bool {
	return len(m.Mods)+
		len(m.Queries)+
		len(m.Controls)+
		len(m.Benchmarks)+
		len(m.Variables)+
		len(m.Dashboards)+
		len(m.DashboardContainers)+
		len(m.DashboardCards)+
		len(m.DashboardCharts)+
		len(m.DashboardFlows)+
		len(m.DashboardHierarchies)+
		len(m.DashboardImages)+
		len(m.DashboardInputs)+
		len(m.DashboardTables)+
		len(m.DashboardTexts)+
		len(m.References) == 0
}

// this is used to create an optimized ModResources containing only the queries which will be run
func (m *ModResources) addControlOrQuery(provider QueryProvider) {
	switch p := provider.(type) {
	case *Query:
		if p != nil {
			m.Queries[p.FullName] = p
		}
	case *Control:
		if p != nil {
			m.Controls[p.FullName] = p
		}
	}
}

// WalkResources calls resourceFunc for every resource in the mod
// if any resourceFunc returns false or an error, return immediately
func (m *ModResources) WalkResources(resourceFunc func(item HclResource) (bool, error)) error {
	for _, r := range m.Queries {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range m.Controls {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range m.Benchmarks {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range m.Dashboards {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range m.DashboardContainers {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range m.DashboardCards {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range m.DashboardCharts {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range m.DashboardFlows {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range m.DashboardHierarchies {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range m.DashboardImages {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range m.GlobalDashboardInputs {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, inputsForDashboard := range m.DashboardInputs {
		for _, r := range inputsForDashboard {
			if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
				return err
			}
		}
	}
	for _, r := range m.DashboardTables {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range m.DashboardTexts {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range m.Variables {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	for _, r := range m.Locals {
		if continueWalking, err := resourceFunc(r); err != nil || !continueWalking {
			return err
		}
	}
	return nil
}

func (m *ModResources) AddResource(item HclResource) hcl.Diagnostics {
	var diags hcl.Diagnostics
	switch r := item.(type) {
	case *Query:
		name := r.Name()
		if existing, ok := m.Queries[name]; ok {
			diags = append(diags, checkForDuplicate(existing, item)...)
			break
		}
		m.Queries[name] = r
		// also add to LocalQueries
		m.LocalQueries[r.GetUnqualifiedName()] = r

	case *Control:
		name := r.Name()
		if existing, ok := m.Controls[name]; ok {
			diags = append(diags, checkForDuplicate(existing, item)...)
			break
		}
		m.Controls[name] = r
		// also add to LocalControls
		m.LocalControls[r.GetUnqualifiedName()] = r

	case *Benchmark:
		name := r.Name()
		if existing, ok := m.Benchmarks[name]; ok {
			diags = append(diags, checkForDuplicate(existing, item)...)
			break
		}
		m.Benchmarks[name] = r
		// also add to LocalBenchmarks
		m.LocalBenchmarks[r.GetUnqualifiedName()] = r

	case *Dashboard:
		name := r.Name()
		if existing, ok := m.Dashboards[name]; ok {
			diags = append(diags, checkForDuplicate(existing, item)...)
			break
		}
		m.Dashboards[name] = r

	case *DashboardContainer:
		name := r.Name()
		if existing, ok := m.DashboardContainers[name]; ok {
			diags = append(diags, checkForDuplicate(existing, item)...)
			break
		}
		m.DashboardContainers[name] = r

	case *DashboardCard:
		name := r.Name()
		if existing, ok := m.DashboardCards[name]; ok {
			diags = append(diags, checkForDuplicate(existing, item)...)
			break
		} else {
			m.DashboardCards[name] = r
		}

	case *DashboardChart:
		name := r.Name()
		if existing, ok := m.DashboardCharts[name]; ok {
			diags = append(diags, checkForDuplicate(existing, item)...)
			break
		}
		m.DashboardCharts[name] = r

	case *DashboardFlow:
		name := r.Name()
		if existing, ok := m.DashboardFlows[name]; ok {
			diags = append(diags, checkForDuplicate(existing, item)...)
			break
		}
		m.DashboardFlows[name] = r

	case *DashboardHierarchy:
		name := r.Name()
		if existing, ok := m.DashboardHierarchies[name]; ok {
			diags = append(diags, checkForDuplicate(existing, item)...)
			break
		}
		m.DashboardHierarchies[name] = r

	case *DashboardImage:
		name := r.Name()
		if existing, ok := m.DashboardImages[name]; ok {
			diags = append(diags, checkForDuplicate(existing, item)...)
			break
		}
		m.DashboardImages[name] = r

	case *DashboardInput:
		// if input has a dashboard asssigned, add to DashboardInputs
		name := r.Name()
		if dashboardName := r.DashboardName; dashboardName != "" {
			inputsForDashboard := m.DashboardInputs[dashboardName]
			if inputsForDashboard == nil {
				inputsForDashboard = make(map[string]*DashboardInput)
				m.DashboardInputs[dashboardName] = inputsForDashboard
			}
			// no need to check for dupes as we have already checked before adding the input to th m od
			inputsForDashboard[name] = r
			break
		}

		// so Dashboard Input must be global
		if existing, ok := m.GlobalDashboardInputs[name]; ok {
			diags = append(diags, checkForDuplicate(existing, item)...)
			break
		}
		m.GlobalDashboardInputs[name] = r

	case *DashboardTable:
		name := r.Name()
		if existing, ok := m.DashboardTables[name]; ok {
			diags = append(diags, checkForDuplicate(existing, item)...)
			break
		}
		m.DashboardTables[name] = r

	case *DashboardText:
		name := r.Name()
		if existing, ok := m.DashboardTexts[name]; ok {
			diags = append(diags, checkForDuplicate(existing, item)...)
			break
		}
		m.DashboardTexts[name] = r

	case *Variable:
		// NOTE: add variable by unqualified name
		name := r.UnqualifiedName
		if existing, ok := m.Variables[name]; ok {
			diags = append(diags, checkForDuplicate(existing, item)...)
			break
		}
		m.Variables[name] = r

	case *Local:
		name := r.Name()
		if existing, ok := m.Locals[name]; ok {
			diags = append(diags, checkForDuplicate(existing, item)...)
			break
		}
		m.Locals[name] = r

	}
	return diags
}

func (m *ModResources) Merge(others []*ModResources) *ModResources {
	res := NewWorkspaceResourceMaps(m.Mod)
	sourceMaps := append([]*ModResources{m}, others...)

	// take local resources from ourselves
	for k, v := range m.LocalQueries {
		res.LocalQueries[k] = v
	}
	for k, v := range m.LocalControls {
		res.LocalControls[k] = v
	}
	for k, v := range m.LocalBenchmarks {
		res.LocalBenchmarks[k] = v
	}

	for _, source := range sourceMaps {
		for k, v := range source.Mods {
			res.Mods[k] = v
		}
		for k, v := range source.Queries {
			res.Queries[k] = v
		}
		for k, v := range source.Controls {
			res.Controls[k] = v
		}
		for k, v := range source.Benchmarks {
			res.Benchmarks[k] = v
		}
		for k, v := range source.Locals {
			res.Locals[k] = v
		}
		for k, v := range source.Variables {
			// NOTE: only include variables from root mod  - we add in the others separately
			if v.Mod.FullName == m.Mod.FullName {
				res.Variables[k] = v
			}
		}
		for k, v := range source.Dashboards {
			res.Dashboards[k] = v
		}
		for k, v := range source.DashboardContainers {
			res.DashboardContainers[k] = v
		}
		for k, v := range source.DashboardCards {
			res.DashboardCards[k] = v
		}
		for k, v := range source.DashboardCharts {
			res.DashboardCharts[k] = v
		}
		for k, v := range source.DashboardFlows {
			res.DashboardFlows[k] = v
		}
		for k, v := range source.DashboardHierarchies {
			res.DashboardHierarchies[k] = v
		}
		for k, v := range source.DashboardImages {
			res.DashboardImages[k] = v
		}
		for k, v := range source.DashboardInputs {
			res.DashboardInputs[k] = v
		}
		for k, v := range source.GlobalDashboardInputs {
			res.GlobalDashboardInputs[k] = v
		}
		for k, v := range source.DashboardTables {
			res.DashboardTables[k] = v
		}
		for k, v := range source.DashboardTexts {
			res.DashboardTexts[k] = v
		}
	}

	return res
}
