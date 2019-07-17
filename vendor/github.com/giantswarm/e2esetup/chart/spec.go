package chart

type LegacyFramework interface {
	CreateNamespace(string) error
	Teardown()
}
