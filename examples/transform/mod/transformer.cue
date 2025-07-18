Transformers: {
	apps: v1: Deployment: #Deployment
}

#Deployment: #Resource & {
	apiVersion: "apps/v1"
	kind:       "Deployment"

	metadata: #Metadata

	spec: {
		selector: matchLabels: [string]: string
		replicas: int
		template: {
			metadata: _

			spec: {
				containers: [..._ & {
					securityContext: {
						privileged:             bool | *false
						readOnlyRootFilesystem: bool | *true
					}
					...
				}]

				securityContext: {
					securityContext: {
						runAsNonRoot: bool | *true
						runAsUser:    int | *10001
						runAsGroup:   int | *10001
					}
				}
				...
			}
		}
	}
}

#Metadata: {
	name:      string
	namespace: string

	annotations?: [string]: string
	labels?: [string]:      string
}

#Resource: {
	apiVersion: string
	kind:       string

	metadata: _
	spec:     _

	...
}
