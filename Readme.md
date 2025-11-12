# AWS EKS Commands


## Create Datavase 

```
docker exec -it postgres18 psql -U root
```

## DB Docs 
```
dbdocs build doc/db.dbml
dbdocs password --set secret --project simple_bank
```

## Update Kubeconfig

```bash
aws eks update-kubeconfig --name test-cluster --region eu-central-1
```

Adds the required configuration to `~/.kube/config` file to connect to the EKS cluster with kubectl.

## Switch Context

```bash
kubectl config use-context arn:aws:eks:eu-central-1:369270180377:cluster/test-cluster
```

Switches the current kubectl context to the specified EKS cluster.

## Apply AWS Auth ConfigMap

```bash
kubectl apply -f eks/aws-auth.yaml
```

Applies the aws-auth ConfigMap to the cluster. This configures which IAM users and roles can access the EKS cluster and their permissions.