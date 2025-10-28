package main

import (
	"back-go/services/account"
	"back-go/services/auth"
	"back-go/services/email"
	"back-go/services/models"
	"back-go/services/order"
	"back-go/services/pricing"
	"back-go/services/repository"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/rs/cors"
	"golang.org/x/time/rate"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	stdjwt "github.com/golang-jwt/jwt/v4"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Config struct {
	Dsn               string       `json:"dsn"`
	Smtp              email.Config `json:"smtp"`
	GeolocationAPIKey string       `json:"geolocation_api_key"`
	JWTSecretKey      string       `json:"jwt_secret_key"`
	EmailDev          bool         `json:"email_dev"`
}

func main() {

	var config Config
	{
		path, _ := filepath.Abs("./main/config.json")
		file, err := os.Open(path)
		if err != nil {
			panic(err)
			return
		}
		err = json.NewDecoder(file).Decode(&config)
		if err != nil {
			panic(err)
			return
		}
	}

	db, err := gorm.Open(postgres.Open(config.Dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database" + err.Error())
	}

	err = db.AutoMigrate(models.Order{}, models.Package{}, models.Status{}, models.Account{}, models.TokenRecord{})
	if err != nil {
		return
	}

	logger := log.NewLogfmtLogger(os.Stdout)

	fieldKeys := []string{"method", "error"}
	promNamespace := "package_app"
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: promNamespace,
		Subsystem: "auth_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: promNamespace,
		Subsystem: "auth_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: promNamespace,
		Subsystem: "auth_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{})

	var emailService email.Service
	{
		if config.EmailDev {
			emailService = email.CreateConsoleService(os.Stdout)
		} else {
			err = email.TestEmailConfig(config.Smtp)
			if err != nil {
				logger.Log("error", err)
			}
			emailService = email.CreateEmailService(config.Smtp)
		}

	}
	var pricingService pricing.Service
	{
		pricingService = pricing.CreatePricingService(config.GeolocationAPIKey, rate.NewLimiter(rate.Every(time.Second*2), 1), pricing.Pricing{BasePrice: 1000, KmPrice: 15})
		pricingService = pricing.LoggingMiddleware{Logger: logger, Next: pricingService}
		pricingService = pricing.InstrumentingMiddleware{
			RequestCount:   requestCount,
			RequestLatency: requestLatency,
			CountResult:    countResult,
			Next:           pricingService,
		}
	}

	var accountService account.Service
	{
		accountService = account.CreateAccountService(repository.AccountRepository{Db: db}, emailService)
		accountService = account.LoggingMiddleware{Logger: logger, Next: accountService}
		accountService = account.InstrumentingMiddleware{
			RequestCount:   requestCount,
			RequestLatency: requestLatency,
			CountResult:    countResult,
			Next:           accountService,
		}
	}

	jwtRepository := repository.JWTRepository{Db: db}

	var authService auth.Service
	{
		authService = auth.CreateAuthService([]byte(config.JWTSecretKey), repository.AccountRepository{Db: db}, jwtRepository)
		authService = auth.LoggingMiddleware{Logger: logger, Next: authService}
		authService = auth.InstrumentingMiddleware{
			RequestCount:   requestCount,
			RequestLatency: requestLatency,
			CountResult:    countResult,
			Next:           authService,
		}
	}

	var orderService order.Service
	{
		orderService = order.CreateOrderService(repository.OrderRepository{Db: db}, pricingService)
	}

	keyFunction := func(token *stdjwt.Token) (interface{}, error) { return []byte(config.JWTSecretKey), nil }
	Auth := func(e endpoint.Endpoint) endpoint.Endpoint {
		return JWTParser(keyFunction, stdjwt.SigningMethodHS256, jwtRepository)(e)
	}

	createOrderHandler := httptransport.NewServer(
		Auth(order.MakeCreateOrderRequest(orderService)),
		order.DecodeCreateOrderRequest,
		EncodeResponse,
		httptransport.ServerBefore(jwt.HTTPToContext()),
		httptransport.ServerErrorEncoder(JSONErrorEncoder),
	)

	//Pricing Service Endpoint
	calculatePriceHandler := httptransport.NewServer(
		pricing.MakeCalculatePriceRequest(pricingService),
		pricing.DecodeCalculatePriceRequest,
		EncodeResponse,
		httptransport.ServerErrorEncoder(JSONErrorEncoder),
	)

	//Account Service Endpoints
	registerHandler := httptransport.NewServer(
		account.MakeRegisterEndpoint(accountService),
		account.DecodeRegisterRequest,
		EncodeResponse,
		httptransport.ServerErrorEncoder(JSONErrorEncoder),
	)

	getHandler := httptransport.NewServer(
		Auth(account.MakeGetEndpoint(accountService)),
		httptransport.NopRequestDecoder,
		EncodeResponse,
		httptransport.ServerBefore(jwt.HTTPToContext()),
		httptransport.ServerErrorEncoder(JSONErrorEncoder),
	)

	//Auth Service Endpoints

	loginHandler := httptransport.NewServer(
		auth.MakeLoginEndpoint(authService),
		auth.DecodeLoginRequest,
		EncodeResponse,
		httptransport.ServerErrorEncoder(JSONErrorEncoder),
	)

	verifyHandler := httptransport.NewServer(
		account.MakeVerifyEndpoint(accountService),
		account.DecodeVerifyRequest,
		EncodeResponse,
		httptransport.ServerErrorEncoder(JSONErrorEncoder),
	)

	logoutHandler := httptransport.NewServer(
		JWTParser(keyFunction, stdjwt.SigningMethodHS256, jwtRepository)(auth.MakeLogoutEndpoint(authService)),
		httptransport.NopRequestDecoder,
		EncodeResponse,
		httptransport.ServerErrorEncoder(JSONErrorEncoder),
		httptransport.ServerBefore(jwt.HTTPToContext()),
	)
	refreshHandler := httptransport.NewServer(
		JWTParser(keyFunction, stdjwt.SigningMethodHS256, jwtRepository)(auth.MakeRefreshEndpoint(authService)),
		auth.DecodeRefreshRequest,
		EncodeResponse,
		httptransport.ServerErrorEncoder(JSONErrorEncoder),
		httptransport.ServerBefore(jwt.HTTPToContext()),
	)

	mux := http.NewServeMux()
	mux.Handle("POST /login", loginHandler)
	mux.Handle("POST /register", registerHandler)
	mux.Handle("POST /logout", logoutHandler)
	mux.Handle("POST /refresh", refreshHandler)
	mux.Handle("GET /price", calculatePriceHandler)
	mux.Handle("POST /order/create", createOrderHandler)
	mux.Handle("GET /metrics", promhttp.Handler())
	mux.Handle("GET /get", getHandler)
	mux.Handle("POST /verify", verifyHandler)

	//handler := cors.AllowAll().Handler(mux)
	opt := cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}
	handler := cors.New(opt).Handler(mux)
	logger.Log("msg", "HTTP", "addr", ":8080")
	logger.Log("err", http.ListenAndServe(":8080", handler))
}

// gokit jwt implementáció kibővítve az adatbázisból törölt tokenek szűrésével
func JWTParser(keyFunc stdjwt.Keyfunc, method stdjwt.SigningMethod, jwtRepository repository.JWTRepository) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			tokenString, ok := ctx.Value(jwt.JWTContextKey).(string)
			if !ok {
				return nil, jwt.ErrTokenContextMissing
			}

			newClaims := func() stdjwt.Claims {
				return &stdjwt.RegisteredClaims{}
			}

			token, err := stdjwt.ParseWithClaims(tokenString, newClaims(), func(token *stdjwt.Token) (interface{}, error) {
				if token.Method != method {
					return nil, jwt.ErrUnexpectedSigningMethod
				}

				return keyFunc(token)
			})

			if !token.Valid {
				id := token.Claims.(*stdjwt.RegisteredClaims).ID
				repoErr := jwtRepository.Delete(id)
				if repoErr != nil {
					err = errors.Join(err, repoErr)
				}
				return nil, err
			}
			id := token.Claims.(*stdjwt.RegisteredClaims).ID
			_, err = jwtRepository.Find(id)
			if err != nil {
				return nil, err
			}

			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, jwt.ErrTokenExpired
			}

			ctx = context.WithValue(ctx, jwt.JWTClaimsContextKey, token.Claims)

			return next(ctx, request)
		}
	}
}

func JSONErrorEncoder(ctx context.Context, err error, w http.ResponseWriter) {

	status := http.StatusInternalServerError
	if errors.Is(err, io.EOF) {
		err = errors.New("empty request")
		status = http.StatusBadRequest
	}

	if errors.Is(err, account.ErrInvalidEmail) ||
		errors.Is(err, account.ErrDuplicateEmail) ||
		errors.Is(err, account.ErrInvalidPhoneNumber) ||
		errors.Is(err, account.ErrInvalidVerificationCode) ||
		errors.Is(err, account.ErrNotVerified) ||
		errors.Is(err, account.ErrAccountNotExist) ||
		errors.Is(err, account.ErrAlreadyVerified) {
		status = http.StatusBadRequest
	}
	if errors.Is(err, auth.ErrInvalidCredentials) ||
		errors.Is(err, stdjwt.ErrTokenExpired) {
		status = http.StatusUnauthorized

	}

	contentType, body := "application/json; charset=utf-8", "{\"error\": \""+err.Error()+"\" }"
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)
	w.Write([]byte(body))
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
