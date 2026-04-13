package routes

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vod/graph"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewGraphQLRouter(db *mongo.Database) *fiber.App {
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{DB: db}}))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	app := fiber.New()
	app.Get("/", adaptor.HTTPHandler(playground.Handler("GraphQL playground", "/query")))
	app.All("/query", adaptor.HTTPHandler(srv))

	return app
}

