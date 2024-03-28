# elder-peeler

Simple "Kubernetes Event Driven Autoscaling" [KEDA](https://github.com/kedacore/keda) demonstration for [AWS EKS](https://aws.amazon.com/eks/). 

This repository:
1. Deploys a simple SQS consumer (and shares SQS values w/prometheus).
1. Deploys KEDA to cluster, to scale the consumer depending upon SQS message population.

## SQS Consumer Application
1.  Assumes same AWS account, region, etc.
1.  Create a [AWS SQS queue](https://aws.amazon.com/sqs/)
1.  Create a IAM role
    1.  [example](https://github.com/guycole/elder-peeler/blob/main/iam_role.json)
    1.  Be sure to add SQS permissions
1.  Create a namespace for your test
    1.  ```kubectl create ns guytest```
1.  Create service account
    1.  Be sure to update the annotation to reflect the IAM role.
    1.  ```kubectl apply -f service_account.yaml```
    1.  Create an IRSA account that grants access to SQS
        1. iam_role.json is an example
        1. be sure to add SQS access to role
1.  Deploy the SQS consumer to your cluster.
    1.  "elder-peeler" is the name of the SQS consumer.
    1.  There is a Dockerfile to help you get started.
        1. SQS queue name is defined within Dockerfile
    1.  There is a deployment.yaml to help you get started
1.  Sanity test.
    1.  There should be an elder-peeler pod active in your namespace.
    1.  Write some messages to SQS (via the console or try [producer.py](https://github.com/guycole/aws-lab/tree/master/sqs_driver)
    1.  Tail the pod log and it should be writing happy messages about message consumption.  
1.  Optional prometheus step
    1.  Apply pod-monitor.yaml and elder-peeler and prometheus will scrape elder-peeler
        1. exposes peeler_messages_available and peeler_messages_consumed
## Apply KEDA
1.  Create a IAM role
    1.  [example](https://github.com/guycole/elder-peeler/blob/main/keda_role.json)
    1.  Be sure to add SQS permissions
1.  Install KEDA via helm
    1.  ```helm install keda kedacore/keda --namespace keda --create-namespace --version 2.13.2 -f values.yaml```
    1.  It takes a minute for gRPC to get connected
1.  Apply ScaledObject 
    1. ```kubectl apply -f scaled.yaml```
1.  Test KEDA
    1. Write messages to your SQS queue.
    1. Before KEDA, there would have been only a single replica for all traffic but now the scaler should create many pods handle the increase in SQS traffic.
    1. After the SQS queue is emptied, the pods will all be removed by KEDA 
1.  Cleanup
    1. Be sure to read cleanup instructions on KEDA website, it isn't enough to invoke ```helm delete```
