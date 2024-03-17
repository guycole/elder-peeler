package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/aws-sdk-go/aws"
)

const banner = "elder-peeler 0.0"

type prometheusMetrics struct {
	messagesAvailable prometheus.Gauge
	messagesConsumed  prometheus.Counter
}

func newMetrics(register prometheus.Registerer) *prometheusMetrics {
	pco := prometheus.CounterOpts{
		Name: "peeler_messages_consumed",
		Help: "peeler messages consumed",
	}

	pgo := prometheus.GaugeOpts{
		Name: "peeler_messages_available",
		Help: "peeler messages available",
	}

	metrics := &prometheusMetrics{
		messagesAvailable: prometheus.NewGauge(pgo),
		messagesConsumed:  prometheus.NewCounter(pco),
	}

	register.MustRegister(metrics.messagesAvailable)
	register.MustRegister(metrics.messagesConsumed)

	return metrics
}

func getRoot(ww http.ResponseWriter, rr *http.Request) {
	log.Println("get root")
	ww.Header().Set("Content-Type", "text/html")
	ww.WriteHeader(http.StatusOK)
	fmt.Fprintf(ww, banner)
}

func sqsClient() *sqs.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := sqs.NewFromConfig(cfg)

	return client
}

func sqsDelete(messages []types.Message, queueUrl string) {
	log.Println("begin sqs delete")

	client := sqsClient()

	targets := make([]types.DeleteMessageBatchRequestEntry, len(messages))
	for ndx := range messages {
		targets[ndx].Id = aws.String(fmt.Sprintf("%v", ndx))
		targets[ndx].ReceiptHandle = messages[ndx].ReceiptHandle
	}

	dmbi := sqs.DeleteMessageBatchInput{
		Entries:  targets,
		QueueUrl: aws.String(queueUrl),
	}

	_, err := client.DeleteMessageBatch(context.TODO(), &dmbi)
	if err == nil {
		log.Printf("deleted:%d\n", len(messages))
	} else {
		log.Printf("delete failure: %v\n", err)
	}

	log.Println("end sqs delete")
}

// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/sqs
func sqsPoll(metrics *prometheusMetrics, queueUrl string) {
	log.Println("begin sqs poll")

	client := sqsClient()

	// discover waiting message population

	gcai := sqs.GetQueueAttributesInput{
		AttributeNames: []types.QueueAttributeName{"All"},
		QueueUrl:       aws.String(queueUrl),
	}

	attribute, err := client.GetQueueAttributes(context.TODO(), &gcai)
	if err != nil {
		log.Fatalf("attribute, %v", err)
	}

	//log.Println(attribute)
	//log.Println(reflect.TypeOf(attribute))

	messagePopulation, err := strconv.Atoi(attribute.Attributes["ApproximateNumberOfMessages"])
	log.Println(messagePopulation)
	log.Println(reflect.TypeOf(messagePopulation))

	metrics.messagesAvailable.Set(float64(messagePopulation))

	// attempt to read waiting messages

	rmi := sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueUrl),
		MaxNumberOfMessages: 5,
		WaitTimeSeconds:     5,
	}

	result, err := client.ReceiveMessage(context.TODO(), &rmi)
	if err != nil {
		log.Printf("receive failure, %v", err)
	}

	messages := result.Messages
	log.Println(len(messages))
	log.Println(reflect.TypeOf(messages))

	for _, message := range messages {
		log.Println(message)
		log.Println(reflect.TypeOf(message))
		//log.Println(*message.ReceiptHandle)
		log.Println(*message.MessageId)
		log.Println(*message.Body)

		metrics.messagesConsumed.Inc()
	}

	if len(messages) > 0 {
		sqsDelete(messages, queueUrl)
	}

	log.Println("end sqs poll")
}

func pollBoy(metrics *prometheusMetrics, queueUrl string) {
	for {
		sqsPoll(metrics, queueUrl)
		log.Println("...sleeping...")
		time.Sleep(23 * time.Second)
	}
}

func main() {
	log.Println(banner)

	registry := prometheus.NewRegistry()
	metrics := newMetrics(registry)
	metrics.messagesAvailable.Set(0.0)

	//updateMetrics()

	go pollBoy(metrics, os.Getenv("SQS_URL"))

	router := mux.NewRouter()
	router.HandleFunc("/", getRoot).Methods("GET")
	router.Path("/metrics").Handler(promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry}))
	log.Fatal(http.ListenAndServe(":80", router))
}
