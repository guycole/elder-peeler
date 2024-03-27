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
1.  Apply KEDA
    1. xxx
1.  Test KEDA
    1. xxx
