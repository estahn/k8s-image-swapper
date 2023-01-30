#!/usr/bin/env bash

# set -e

tag=$1

rm -f k8s-image-swapper
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

docker build -t k8s-image-swapper:"$tag" .


push_thing () {
  local tag=$1
  local account=$2
  local region=$3
  docker tag k8s-image-swapper:"$tag" "$account".dkr.ecr."$region".amazonaws.com/ghcr.io/estahn/k8s-image-swapper:"$tag"
  AWS_PROFILE=stageeng docker push "${account}".dkr.ecr."$region".amazonaws.com/ghcr.io/estahn/k8s-image-swapper:"$tag"
}

# for r in us-west-2 ap-southeast-2 us-east-1 eu-west-1; do
# for a in "520455238173" "035088524874"; do
#   for r in us-west-2 ap-southeast-2 us-east-1 eu-west-1; do
#     push_thing "$1" "${a}" "${r}" &
#   done
# done

push_thing "$1" "520455238173" "us-west-2"
# push_thing "$1" "035088524874" "us-wnest-2"

kubectl --context dev-us-west-2 -n kube-system set image deploy/k8s-image-swapper k8s-image-swapper=520455238173.dkr.ecr.us-west-2.amazonaws.com/ghcr.io/estahn/k8s-image-swapper:"$tag"
# kubectl --context stage-us-west-2 -n kube-system set image deploy/k8s-image-swapper k8s-image-swapper=035088524874.dkr.ecr.us-west-2.amazonaws.com/ghcr.io/estahn/k8s-image-swapper:"$tag"
