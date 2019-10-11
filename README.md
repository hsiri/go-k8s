# go-k8s
All kubernetes example objects used in this techtalk are found in `k8s` directory

## Required Installation
- kubectl
	- Please follow steps here :  ([Kubectl]( https://kubernetes.io/docs/tasks/tools/install-kubectl/))

- kops for AWS
	- Please follow steps here :  ([KOPS]( https://kubernetes.io/docs/setup/production-environment/tools/kops/))
	- Since we are using AWS in the talk, please make sure you have an AWS account with **EC2, S3, IAM FULL ACCESS permission**

## Setting up k8s local cluster in AWS (to run in your machine)

#### Configure AWS
```
$ aws configure
#Provide your AWS Access Key ID, AWS Secret Access Key, AWS Default Region(ap-south-1) and Default output format (json)
# Generate access and secret from here - (Link)
# Because "aws configure" doesn't export these vars for kops to use, we export them now

export AWS_ACCESS_KEY_ID=$(aws configure get aws_access_key_id)
export AWS_SECRET_ACCESS_KEY=$(aws configure get aws_secret_access_key)
```

#### Create S3 Bucket (change on your desired bucket name)
```
aws s3api create-bucket -- bucket demo-k8s-local -- region ap-southeast-1 -- create-bucket-configuration LocationConstraint=ap-southeast-1
aws s3api put-bucket-versioning --bucket demo-k8s-local --versioning-configuration Status=Enabled
aws s3api put-bucket-encryption --bucket demo-k8s-local --server-side-encryption-configuration '{"Rules":[{"ApplyServerSideEncryptionByDefault":{"SSEAlgorithm":"AES256"}}]}'
```

#### Creating your first cluster
For testing purposes, we will be creating a gossip-based cluster only. To do it, name your cluster _<desired_name>.k8s.local_
```
export NAME=demo.k8s.local
export KOPS_STATE_STORE=s3://demo-k8s-local

kops create cluster --master-size=t3a.medium --node-size=t3a.medium --zones=ap-southeast-1a ${NAME}
kops update cluster --name demo.k8s.local --yes
```

This will take time to create the cluster, so just wait until it finishes. Also, you can check on AWS EC2 dashboard and there should be 3 running instances(1 master node and 2 worker nodes) of the created cluster.

#### Validating your cluster
```
kops validate cluster

Using cluster from kubectl context: demo.k8s.local

Validating cluster demo.k8s.local

INSTANCE GROUPS
NAME			ROLE	MACHINETYPE	MIN	MAX	SUBNETS
master-ap-southeast-1a	Master	t3a.medium	1	1	ap-southeast-1a
nodes			Node	t3a.medium	2	2	ap-southeast-1a

NODE STATUS
NAME							ROLE	READY
ip-xxxx.ap-southeast-1.compute.internal	master	True
ip-xxxx.ap-southeast-1.compute.internal		node	True
ip-xxxx.ap-southeast-1.compute.internal	node	True

Your cluster demo.k8s.local is ready
```

## Kubernetes Dashboard
For easy monitoring, we can install the GUI of Kubernetes

```
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v1.10.1/src/deploy/recommended/kubernetes-dashboard.yaml

kubectl create serviceaccount dashboard -n default

kubectl create clusterrolebinding dashboard-admin -n default \
  --clusterrole=cluster-admin \
  --serviceaccount=default:dashboard

kubectl get secret $(kubectl get serviceaccount dashboard -o jsonpath="{.secrets[0].name}") -o jsonpath="{.data.token}" | base64 --decode
# save your token somewhere else for dashboard login

kubectl proxy

# open browser
# dashboard link : http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/#!/login
# paste your generated token here
```

## Pod
```
kubectl create -f k8s/pod.yml
kubectl delete pod go-k8s
```

## Deployment
```
kubectl create -f k8s/deployment.yml
```

## Services
```
kubectl create -f k8s/service.yml
```

## Secrets
```
kubectl create secret generic go-k8s-credentials --save-config --from-env-file=env.local --dry-run -o yaml | kubectl apply -f -
```

## Updating deployment if docker image is also updated
```
kubectl set image deployments/go-k8s go-k8s=hsiri/go-k8s:1.2.0

kubectl replace -f k8s-deployment.yml
```