Generators: {
	"acme.com": v1: Application: #ApplicationGenerator
}

Validators: {
	"acme.com": v1: Application: #Application
}

#Application: X={
	apiVersion: "acme.com/v1"
	kind:       "Application"

	metadata: #Metadata & {
		name:      string
		namespace: string

		labels: app: X.metadata.name
	}

	spec: {
		expose: [Name=string]: {
			name: Name
			port: int
		}

		image: {
			registry: string
			name:     string
			tag:      string
		}

		resources: {
			requests: {
				cpu:    number | *1.0
				memory: int | *256M
			}
			limits: {
				cpu:    (number & >=X.spec.resources.requests.cpu) | *(X.spec.resources.requests.cpu * 3)
				memory: X.spec.resources.requests.memory
			}
		}

		scaling: {
			horizontal: {
				replicas: {
					min: int | *3
					max: (int & >=X.spec.replicas.min) | *5
				}
			}
		}
	}

	$patch: {
		deployment: {...}
		horizontalPodAutoscaler: {...}
		podDisruptionBudget: {...}
		service: {...}
	}
}

#ApplicationGenerator: X=#Application & {
	_$$resource: {
		[string]: _ & {
			metadata: X.metadata
		}
	}

	_$$resource: {
		Deployment: #Deployment & {
			spec: {
				selector: X.metadata.labels
				strategy: {
					rollingUpdate: {
						maxSurge:       "50%"
						maxUnavailable: "0%"
					}
					type: "RollingUpdate"
				}
				template: spec: {
					containers: [{
						name:  X.metadata.name
						image: "\(X.spec.image.registry)/\(X.spec.image.name):\(X.spec.image.tag)"
						ports: [for _, p in X.spec.expose {
							name:          p.name
							containerPort: p.port
							protocol:      "TCP"
						}]
						resources: X.spec.resources
					}, ...]
				}
			}
		} & X.$patch.deployment

		HorizontalPodAutoscaler: #HorizontalPodAutoscaler & {
			spec: {
				minReplicas: X.spec.scaling.horizontal.replicas.min
				maxReplicas: X.spec.scaling.horizontal.replicas.max

				metrics: [...] | *[{
					type: "ContainerResource"
					containerResource: {
						container: X._$$resource.Deployment.spec.template.spec.containers[0].name
						name:      "cpu"
						target: {
							averageUtilization: 70
							type:               "Utilization"
						}
					}
				}, ...]

				scaleTargetRef: {
					apiVersion: X._$$resource.Deployment.apiVersion
					kind:       X._$$resource.Deployment.kind
					name:       X._$$resource.Deployment.metadata.name
				}
			}
		} & X.$patch.horizontalPodAutoscaler

		if X.spec.scaling.horizontal.replicas.max > 1 {
			PodDisruptionBudget: #PodDisruptionBudget & {
				spec: {
					minAvailable: "50%"
					selector:     X._$$resource.Deployment.metadata.labels
				}
			} & X.$patch.podDisruptionBudget
		}

		if len(X.spec.expose) > 0 {
			Service: #Service & {
				spec: {
					ports: [for _, p in X.spec.expose {
						name:       p.name
						port:       p.port
						targetPort: p.port
						protocol:   "TCP"
					}]
					selector: X._$$resource.Deployment.metadata.labels
				}
			} & X.$patch.service
		}
	}

	_$$resources: [for _, r in _$$resource {r}]
}

#Deployment: #Resource & {
	apiVersion: "core/v1"
	kind:       "Deployment"

	metadata: #Metadata

	spec: {
		selector: [string]: string
		strategy: _
		template: _
	}
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
