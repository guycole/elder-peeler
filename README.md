# elder-peeler
demonstration to consume from SQS and update prometheus

kubectl create sa elder-peeler -n guytest
(role) irsa-elder-peeler-us-west-2

prometheus: exposes peeler_messages_available and peeler_messages_consumed
