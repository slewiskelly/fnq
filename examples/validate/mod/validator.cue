import (
	"list"
	"struct"
)

Validators: {
	apps: v1: Deployment:                     #Deployment
	autoscaling: v2: HorizontalPodAutoscaler: #HorizontalPodAutoscaler
	core: v1: {
		PodDisruptionBudget: #PodDisruptionBudget
		Service:             #Service
	}
}

#Container: {
	name:  string
	image: string

	resources: {
		requests: {
			cpu:    number & >0
			memory: int & >0
		}
		limits: {
			cpu:    number & >=requests.cpu
			memory: int & requests.memory
		}
	}
}

#Deployment: #Resource & {
	apiVersion: "apps/v1"
	kind:       "Deployment"

	metadata: #Metadata & {
		labels: app: string
	}

	spec: {
		selector: {
			[string]: string
		} & struct.MinFields(1)

		strategy: _

		template: {
			metadata: #Metadata & {
				labels: app: string
			}

			spec: {
				containers: [...#Container] & list.MinItems(1)
			}
		}
	}
}

#HorizontalPodAutoscaler: #Resource & {
	apiVersion: "autoscaling/v2"
	kind:       "HorizontalPodAutoscaler"

	metadata: #Metadata

	spec: {
		minReplicas: int & >0
		maxReplicas: int & >=minReplicas

		metrics: [...] & list.MinItems(1)

		scaleTargetRef: {
			apiVersion: string
			kind:       string
			name:       string
		}

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

		selector: {
			[string]: string
		} & struct.MinFields(1)
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
		}] & list.MinItems(1)

		selector: {
			[string]: string
		} & struct.MinFields(1)

		type: string
	}
}
