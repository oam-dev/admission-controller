package v1alpha1

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	SchemeBuilder.Register(&ComponentSchematic{}, &ComponentSchematicList{})
	SchemeBuilder.Register(&Scope{}, &ScopeList{}, &ApplicationScopeList{}, &ApplicationScope{})
	SchemeBuilder.Register(&Trait{}, &TraitList{})
	SchemeBuilder.Register(&ApplicationConfiguration{}, &ApplicationConfigurationList{})
	AddToSchemes = append(AddToSchemes, SchemeBuilder.AddToScheme)
}
