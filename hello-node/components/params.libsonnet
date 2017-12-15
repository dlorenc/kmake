{
  global: {
    // User-defined global parameters; accessible to all component and environments, Ex:
    // replicas: 4,
  },
  components: {
    // Component-level parameters, defined initially from 'ks prototype use ...'
    // Each object below should correspond to a component in the components/ directory
    "hello-node": {
      containerPort: 8080,
      imagePullPolicy: "IfNotPresent",
      name: "hello-node",
      nodePort: 30005,
      replicas: 1,
      servicePort: 80,
      type: "NodePort",
    },
  },
}
