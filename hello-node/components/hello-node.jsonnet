local params = std.extVar("__ksonnet/params").components["hello-node"];
local k = import "k.libsonnet";
local deployment = k.apps.v1beta1.deployment;
local container = k.apps.v1beta1.deployment.mixin.spec.template.spec.containersType;
local containerPort = container.portsType;
local service = k.core.v1.service;
local servicePort = k.core.v1.service.mixin.spec.portsType;

local image = std.extVar("image");

local targetPort = params.containerPort;
local nodePort = params.nodePort;
local imagePullPolicy = params.imagePullPolicy;
local labels = {app: params.name};

local appService = service
  .new(
    params.name,
    labels,
    servicePort
      .new(params.servicePort, targetPort)
      .withNodePort(nodePort))
  .withType(params.type);

local appDeployment = deployment
  .new(
    params.name,
    params.replicas,
    container
      .new(params.name, image)
      .withPorts(containerPort.new(targetPort))
      .withImagePullPolicy(imagePullPolicy),
    labels)
    .withTerminationGracePeriodSeconds(0);

k.core.v1.list.new([appService, appDeployment])
