Validators: {
	apps: v1: Deployment: #Deployment
	v1: Service: #Service
}

Generators: {
	apps: v1: Deployment: #DeploymentGenerator
}

#Deployment: #Resource & {
	apiVersion: "apps/v1"
	kind:       "Deployment"

	metadata: #Metadata

	spec: {
		selector: matchLabels: [string]: string
		replicas: int
		strategy: _
		template: _
	}
}

#DeploymentGenerator: X=#Deployment & {
	_$$resource: {
		[string]: _ & {
			metadata: X.metadata
		}
	}

	_$$resource: {
		if X.spec.replicas > 1 {
			HorizontalPodAutoscaler: #HorizontalPodAutoscaler & {
				spec: {
					minReplicas: X.spec.replicas
					maxReplicas: minReplicas * 5

					metrics: [...] | *[{
						type: "ContainerResource"
						containerResource: {
							container: X.spec.template.spec.containers[0].name
							name:      "cpu"
							target: {
								averageUtilization: 70
								type:               "Utilization"
							}
						}
					}, ...]

					scaleTargetRef: {
						apiVersion: X.apiVersion
						kind:       X.kind
						name:       X.metadata.name
					}
				}
			}
		}

		if X.spec.replicas > 1 {
			PodDisruptionBudget: #PodDisruptionBudget & {
				spec: {
					minAvailable: "50%"
					selector:     X.metadata.labels
				}
			}
		}
	}

	_$$resources: [for _, r in _$$resource {r}]
}

#HorizontalPodAutoscaler: #Resource & {
	apiVersion: "autoscaling/v2"
	kind:       "HorizontalPodAutoscaler"

	metadata: #Metadata

	spec: {
		maxReplicas: int
		metrics: [...]
		minReplicas:    int
		scaleTargetRef: _
	}
}

#Metadata: {
	name:      string
	namespace: string

	annotations?: [string]: string
	labels?: [string]:      string
}

#PodDisruptionBudget: #Resource & {
	apiVersion: "core/v1"
	kind:       "PodDisruptionBudget"

	metadata: #Metadata

	spec: {
		minAvailable: string
		selector: [string]: string
	}
}

#Resource: {
	apiVersion: string
	kind:       string

	metadata: _
	spec:     _

	...
}

#Service: #Resource & {
	apiVersion: "core/v1"
	kind:       "Service"

	metadata: _

	spec: {
		ports: [...{
			name:       string
			port:       int
			targetPort: int
			protocol:   string
		}]
		selector: [string]: string
		type: "ClusterIP"
	}
}
