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
      image: "hello-node:ea1e422cf8e5e1652db6afbca542c5bec1e1fc8ed5c75825300fbe5b1a7533ba",
      imagePullPolicy: "Never",
      name: "hello-node",
      replicas: 1,
      servicePort: 80,
      type: "NodePort",
    },
  },
}
