#!/bin/sh
set -o errexit
CLUSTER_NAME="${CLUSTER_NAME:-local}"

# create a cluster with the local registry enabled 
echo "🐳 Creating k3d with local Registry"
k3d cluster create -i rancher/k3s:v1.18.8-k3s1 \
  -a 2 --volume $(pwd)/registries.yaml:/etc/rancher/k3s/registries.yaml "${CLUSTER_NAME}"
# Annotate nodes with registry info for Tilt to auto-detect
echo "⏳ Waiting for node(s) + annotating with registry info..."
DONE=""
timeout=$(($(date +%s) + 30))
until [[ $(date +%s) -gt $timeout ]]; do
  nodes=$(kubectl get nodes -o go-template --template='{{range .items}}{{printf "%s\n" .metadata.name}}{{end}}')
  if [ ! -z "${nodes}" ]; then
    for node in $nodes; do
      kubectl annotate node "${node}" \
              tilt.dev/registry=localhost:${reg_port} \
              tilt.dev/registry-from-cluster=localhost:${reg_port}
    done
    DONE=true
    break
  fi
  sleep 0.2
done
if [ -z "$DONE" ]; then
  echo "🚨 Timed out waiting for node(s) to be up"
  exit 1
fi

# create registry container unless it already exists
echo "👀 Checking Docker Registry"
reg_name='registry.localhost'
reg_port='5000'
running="$(docker inspect -f '{{.State.Running}}' "${reg_name}" 2>/dev/null || true)"
if [ "${running}" != 'true' ]; then
  echo "🌵 Creating Docker Registry"
  docker container run \
    -d --name "${reg_name}" \
    -v local_registry:/var/lib/registry \
    --restart=always -e REGISTRY_HTTP_ADDR=0.0.0.0:5000 \
    -p "${reg_port}:5000" \
    registry:2
  docker network connect "k3d-${CLUSTER_NAME}" "${reg_name}"
else
   echo "👍 Docker Registry is existing."
fi

echo "🎉 You are all set!! Run kubectl get nodes to check status."