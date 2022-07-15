package v1alpha1

#apiVersion: string
#kind:       string
#metadata: {
	// TODO: fix regex - requires at least two chars
	name: =~"^[a-z0-9][a-z0-9_]{1,61}[a-z0-9]$" & string
	// optional
	displayName?: string
	// map[string]string, optional
	labels?: [string]: string
}
