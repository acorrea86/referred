package handler

import (
	"log"
	"net/http"

	"blumer-ms-refers/events"
	"blumer-ms-refers/graph"
	"blumer-ms-refers/repository"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

// Handler is the struct that contains the repository and the consumer
type Handler struct {
	Repository *repository.Repository
	Reducer    *events.KafkaReducer
	Producer   *events.KafkaProducer
	Port       string
}

// Handler is the function that handles the requests
func (h *Handler) Handler() {

	resolver := graph.NewResolver(h.Repository)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	go h.Reducer.StartConsumer()

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", h.Port)
	log.Fatal(http.ListenAndServe(":"+h.Port, nil))
}

// NewHandler creates a new handler
func NewHandler(
	repository *repository.Repository,
	reducer *events.KafkaReducer,
	producer *events.KafkaProducer,
	port string,
) *Handler {
	return &Handler{
		Repository: repository,
		Reducer:    reducer,
		Producer:   producer,
		Port:       port,
	}
}
