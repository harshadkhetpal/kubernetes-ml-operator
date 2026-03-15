# ☸️ Kubernetes ML Operator — Custom Resource for ML Workloads

[![Go](https://img.shields.io/badge/Go-1.22-00ADD8?style=flat-square&logo=go&logoColor=white)](https://golang.org)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-1.29-326CE5?style=flat-square&logo=kubernetes&logoColor=white)](https://kubernetes.io)
[![Operator SDK](https://img.shields.io/badge/Operator--SDK-1.34-EE0000?style=flat-square)](https://sdk.operatorframework.io)
[![License: MIT](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)

> A Kubernetes Operator built with Operator SDK that manages `MLTrainingJob` and `MLServingEndpoint` Custom Resources. Automates GPU provisioning, distributed training, and zero-downtime model rollouts.

## 🎯 Custom Resources

### MLTrainingJob
```yaml
apiVersion: mlops.harshadkhetpal.dev/v1alpha1
kind: MLTrainingJob
metadata:
  name: fraud-detection-v2
spec:
  image: harshadkhetpal/fraud-model:latest
  framework: pytorch       # pytorch | tensorflow | xgboost
  replicas: 4
  gpuPerReplica: 1
  dataset: s3://data/fraud/2024
  hyperparameters:
    epochs: "50"
    lr: "0.001"
  tracking:
    mlflow_uri: http://mlflow:5000
    experiment: fraud-detection
```

### MLServingEndpoint
```yaml
apiVersion: mlops.harshadkhetpal.dev/v1alpha1
kind: MLServingEndpoint
metadata:
  name: fraud-endpoint
spec:
  modelUri: s3://models/fraud-detection/v2
  replicas: 3
  autoscaling:
    minReplicas: 2
    maxReplicas: 20
    targetRPS: 1000
```

## 🚀 Install
```bash
kubectl apply -f config/crd/bases/
kubectl apply -f config/rbac/
kubectl apply -f config/manager/
```
