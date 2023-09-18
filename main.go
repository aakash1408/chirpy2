package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type apiconfig struct {
	fileServerHits int
}

func main() {
	const port = "8080"
	const filePathRoot = "."

	apiCfg := apiconfig{
		fileServerHits: 0,
	}

	//mux := http.NewServeMux()
	r := chi.NewRouter()
	fsHandler := apiCfg.middlewareMetric(http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot))))
	r.Handle("/app", fsHandler)
	r.Handle("/app/*", fsHandler)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Get("/reset", apiCfg.handlerReset)
	apiRouter.Post("/validate_chirp", handlerChirpsValidate)
	r.Mount("/api", apiRouter)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiCfg.handlerMetrics)
	r.Mount("/admin", adminRouter)

	corsMux := middlewareCors(r)

	srv := http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}
	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
