package main

import (
	"back-go/services/account"
	"back-go/services/admin"
	"back-go/services/auth"
	"back-go/services/config"
	"back-go/services/email"
	"back-go/services/models"
	"back-go/services/order"
	_package "back-go/services/package"
	"back-go/services/pricing"
	"back-go/services/repository"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/rs/cors"
	"golang.org/x/time/rate"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	stdjwt "github.com/golang-jwt/jwt/v4"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stdout)

	var configuration config.Config
	{
		var flagconfig string
		flag.StringVar(&flagconfig, "c", "config.yaml", "configuration file location")
		flag.Parse()
		path, _ := filepath.Abs(flagconfig)
		file, err := os.Open(path)
		if errors.Is(err, os.ErrNotExist) {
			_ = logger.Log("err", "configuration file does not exist: "+path)
			return
		}
		if err != nil {
			panic(err)
			return
		}
		configuration, err = config.ReadConfig(file)
		if err != nil {
			panic(err)
			return
		}
	}
	db, err := gorm.Open(postgres.Open(configuration.Dsn.String()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database" + err.Error())
	}

	err = db.AutoMigrate(models.Order{}, models.Package{}, models.Status{}, models.Account{}, models.TokenRecord{}, models.Admin{}, models.Pricing{})
	if err != nil {
		return
	}

	hash, _ := argon2id.CreateHash("admin", argon2id.DefaultParams)
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&models.Admin{
		Username: "admin",
		Password: hash,
	})

	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&models.Pricing{
		Model: gorm.Model{
			ID: 1,
		},
		KmPrice:   15,
		BasePrice: 1000,
	})

	fieldKeys := []string{"method", "error"}
	promNamespace := "package_app"
	requestCount := func(subsystem string) *kitprometheus.Counter {
		return kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: promNamespace,
			Subsystem: subsystem,
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys)
	}
	requestLatency := func(subsystem string) *kitprometheus.Histogram {
		return kitprometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
			Namespace: promNamespace,
			Subsystem: subsystem,
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys)
	}

	authRepository := repository.JWTRepository{Db: db}
	accountRepository := repository.AccountRepository{Db: db}
	orderRepository := repository.OrderRepository{Db: db}
	packageRepository := repository.PackageRepository{Db: db}
	adminRepository := repository.AdminRepository{Db: db}
	pricingRepository := repository.PricingRepository{Db: db}

	frontendLocation := configuration.Frontend.Host + ":" + strconv.Itoa(configuration.Frontend.Port)

	var emailService email.Service
	{
		if configuration.EmailDev {
			emailService = email.CreateConsoleService(os.Stdout)
		} else {
			err = email.TestEmailConfig(configuration.Smtp)
			if err != nil {
				_ = logger.Log("error", err)
				panic("smtp test failed, check configuration")

			}
			emailService = email.CreateEmailService(configuration.Smtp)
		}

	}
	var pricingService pricing.Service
	{
		pricingService = pricing.CreatePricingService(configuration.GeolocationAPIKey, rate.NewLimiter(rate.Every(time.Second*2), 1), pricingRepository, configuration.PricingDev)
		pricingService = pricing.LoggingMiddleware{Logger: logger, Next: pricingService}
		pricingService = pricing.InstrumentingMiddleware{
			RequestCount:   requestCount("pricing_service"),
			RequestLatency: requestLatency("pricing_service"),
			Next:           pricingService,
		}
	}

	var accountService account.Service
	{
		accountService = account.CreateAccountService(accountRepository, emailService, frontendLocation)
		accountService = account.LoggingMiddleware{Logger: logger, Next: accountService}
		accountService = account.InstrumentingMiddleware{
			RequestCount:   requestCount("account_service"),
			RequestLatency: requestLatency("account_service"),
			Next:           accountService,
		}
	}

	var authService auth.Service
	{
		authService = auth.CreateAuthService([]byte(configuration.JWTSecretKey), accountRepository, authRepository)
		authService = auth.LoggingMiddleware{Logger: logger, Next: authService}
		authService = auth.InstrumentingMiddleware{
			RequestCount:   requestCount("auth_service"),
			RequestLatency: requestLatency("auth_service"),
			Next:           authService,
		}
	}

	var orderService order.Service
	{
		orderService = order.CreateOrderService(orderRepository, pricingService)
		orderService = order.LoggingMiddleware{Logger: logger, Next: orderService}
		orderService = order.InstrumentingMiddleware{
			RequestCount:   requestCount("order_service"),
			RequestLatency: requestLatency("order_service"),
			Next:           orderService,
		}
	}

	var packageService _package.Service
	{
		packageService = _package.CreatePackageService(packageRepository)
		packageService = _package.LoggingMiddleware{Logger: logger, Next: packageService}
		packageService = _package.InstrumentingMiddleware{
			RequestCount:   requestCount("package_service"),
			RequestLatency: requestLatency("package_service"),
			Next:           packageService,
		}
	}

	var adminService admin.Service
	{
		adminService = admin.CreateAdminService([]byte(configuration.JWTSecretKey), adminRepository, authRepository)
	}

	keyFunction := func(token *stdjwt.Token) (interface{}, error) { return []byte(configuration.JWTSecretKey), nil }
	Auth := func(e endpoint.Endpoint) endpoint.Endpoint {
		return JWTParser(keyFunction, stdjwt.SigningMethodHS256, authRepository, false)(e)
	}
	AdminAuth := func(e endpoint.Endpoint) endpoint.Endpoint {
		return JWTParser(keyFunction, stdjwt.SigningMethodHS256, authRepository, true)(e)
	}
	/** Admin **/

	adminGetAllOrdersHandler := httptransport.NewServer(
		AdminAuth(admin.MakeGetAllOrdersAdminEndpoint(orderService)),
		httptransport.NopRequestDecoder,
		EncodeResponse,
		httptransport.ServerBefore(jwt.HTTPToContext()),
		httptransport.ServerErrorEncoder(JSONErrorEncoder),
	)

	adminGetPricingHandler := httptransport.NewServer(
		AdminAuth(admin.MakeGetPricingAdminEndpoint(pricingService)),
		httptransport.NopRequestDecoder,
		EncodeResponse,
		httptransport.ServerBefore(jwt.HTTPToContext()),
		httptransport.ServerErrorEncoder(JSONErrorEncoder),
	)

	adminSetPricingHandler := httptransport.NewServer(
		AdminAuth(admin.MakeSetPricingAdminEndpoint(pricingService)),
		pricing.DecodePricingRequest,
		EncodeResponse,
		httptransport.ServerBefore(jwt.HTTPToContext()),
		httptransport.ServerErrorEncoder(JSONErrorEncoder),
	)

	adminGetHandler := httptransport.NewServer(
		AdminAuth(admin.MakeAdminGetEndpoint(adminService)),
		httptransport.NopRequestDecoder,
		EncodeResponse,
		httptransport.ServerBefore(jwt.HTTPToContext()),
		httptransport.ServerErrorEncoder(JSONErrorEncoder),
	)

	adminGetAllHandler := httptransport.NewServer(
		AdminAuth(admin.MakeGetAdminsEndpoint(adminService)),
		httptransport.NopRequestDecoder,
		EncodeResponse,
		httptransport.ServerBefore(jwt.HTTPToContext()),
		httptransport.ServerErrorEncoder(JSONErrorEncoder),
	)

	adminDeleteHandler := httptransport.NewServer(
		AdminAuth(admin.MakeDeleteAdminEndpoint(adminService)),
		admin.DecodeIDHeader,
		EncodeResponse,
		httptransport.ServerBefore(jwt.HTTPToContext()),
		httptransport.ServerErrorEncoder(JSONErrorEncoder),
	)

	adminCreateHandler := httptransport.NewServer(
		AdminAuth(admin.MakeCreateAdminEndpoint(adminService)),
		admin.DecodeRequest,
		EncodeResponse,
		httptransport.ServerBefore(jwt.HTTPToContext()),
		httptransport.ServerErrorEncoder(JSONErrorEncoder),
	)

	adminLoginHandler := httptransport.NewServer(
		admin.MakeLoginEndpoint(adminService),
		admin.DecodeRequest,
		EncodeResponse,
		httptransport.ServerErrorEncoder(JSONErrorEncoder),
	)

	adminLogoutHandler := httptransport.NewServer(
		AdminAuth(admin.MakeLogoutEndpoint(adminService)),
		httptransport.NopRequestDecoder,
		EncodeResponse,
		httptransport.ServerBefore(jwt.HTTPToContext()),
		httptransport.ServerErrorEncoder(JSONErrorEncoder),
	)

	adminRefreshHandler := httptransport.NewServer(
		AdminAuth(admin.MakeRefreshEndpoint(adminService)),
		httptransport.NopRequestDecoder,
		EncodeResponse,
		httptransport.ServerBefore(jwt.HTTPToContext()),
		httptransport.ServerErrorEncoder(JSONErrorEncoder),
	)

	adminGetPackageHandler := httptransport.NewServer(
		AdminAuth(admin.MakeAdminGetPackageEndpoint(packageService)),
		admin.DecodeIDHeader,
		EncodeResponse,
		httptransport.ServerBefore(jwt.HTTPToContext()),
		httptransport.ServerErrorEncoder(JSONErrorEncoder),
	)

	adminDeletePackageHandler := httptransport.NewServer(
		AdminAuth(admin.MakeAdminDeletePackageEndpoint(packageService)),
		admin.DecodeIDHeader,
		EncodeResponse,
		httptransport.ServerBefore(jwt.HTTPToContext()),
		httptransport.ServerErrorEncoder(JSONErrorEncoder),
	)

	adminUpdatePackageHandler := httptransport.NewServer(
		AdminAuth(admin.MakeAdminUpdatePackageEndpoint(packageService)),
		admin.DecodeUpdatePackageRequest,
		EncodeResponse,
		httptransport.ServerBefore(jwt.HTTPToContext()),
		httptransport.ServerErrorEncoder(JSONErrorEncoder),
	)

	adminGetOrderHandler := httptransport.NewServer(
		AdminAuth(admin.MakeAdminGetOrderEndpoint(orderService)),
		admin.DecodeIDHeader,
		EncodeResponse,
		httptransport.ServerBefore(jwt.HTTPToContext()),
		httptransport.ServerErrorEncoder(JSONErrorEncoder),
	)

	adminDeleteOrderHandler := httptransport.NewServer(
		AdminAuth(admin.MakeAdminDeleteOrderEndpoint(orderService)),
		admin.DecodeIDHeader,
		EncodeResponse,
		httptransport.ServerBefore(jwt.HTTPToContext()),
		httptransport.ServerErrorEncoder(JSONErrorEncoder),
	)

	adminUpdateOrderHandler := httptransport.NewServer(
		AdminAuth(admin.MakeAdminUpdateOrderEndpoint(orderService)),
		admin.DecodeUpdateOrderRequest,
		EncodeResponse,
		httptransport.ServerBefore(jwt.HTTPToContext()),
		httptransport.ServerErrorEncoder(JSONErrorEncoder),
	)

	adminAddStatusHandler := httptransport.NewServer(
		AdminAuth(admin.MakeAdminAddStatusEndpoint(packageService)),
		_package.DecodeAddPackageStatusRequest,
		EncodeResponse,
		httptransport.ServerBefore(jwt.HTTPToContext()),
		httptransport.ServerErrorEncoder(JSONErrorEncoder),
	)

	/** Admin **/

	createOrderHandler := httptransport.NewServer(
		Auth(order.MakeCreateOrderRequest(orderService)),
		order.DecodeCreateOrderRequest,
		EncodeResponse,
		httptransport.ServerBefore(jwt.HTTPToContext()),
		httptransport.ServerErrorEncoder(JSONErrorEncoder),
	)

	getAllOrderHandler := httptransport.NewServer(
		Auth(order.MakeGetAllOrdersRequest(orderService)),
		httptransport.NopRequestDecoder,
		EncodeResponse,
		httptransport.ServerBefore(jwt.HTTPToContext()),
		httptransport.ServerErrorEncoder(JSONErrorEncoder),
	)

	getPackageStatusHandler := httptransport.NewServer(
		_package.MakeGetPackageStatus(packageService),
		_package.DecodeGetPackageStatusRequest,
		EncodeResponse,
		httptransport.ServerErrorEncoder(JSONErrorEncoder),
	)
	calculatePriceHandler := httptransport.NewServer(
		pricing.MakeCalculatePriceRequest(pricingService),
		pricing.DecodeCalculatePriceRequest,
		EncodeResponse,
		httptransport.ServerErrorEncoder(JSONErrorEncoder),
	)
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
		Auth(auth.MakeLogoutEndpoint(authService)),
		httptransport.NopRequestDecoder,
		EncodeResponse,
		httptransport.ServerErrorEncoder(JSONErrorEncoder),
		httptransport.ServerBefore(jwt.HTTPToContext()),
	)
	refreshHandler := httptransport.NewServer(
		Auth(auth.MakeRefreshEndpoint(authService)),
		httptransport.NopRequestDecoder,
		EncodeResponse,
		httptransport.ServerErrorEncoder(JSONErrorEncoder),
		httptransport.ServerBefore(jwt.HTTPToContext()),
	)

	mux := http.NewServeMux()
	mux.Handle("POST /login", loginHandler)
	mux.Handle("POST /register", registerHandler)
	mux.Handle("POST /logout", logoutHandler)
	mux.Handle("GET /refresh", refreshHandler)
	mux.Handle("POST /price", calculatePriceHandler)
	mux.Handle("POST /order/create", createOrderHandler)
	mux.Handle("GET /metrics", promhttp.Handler())
	mux.Handle("GET /get", getHandler)
	mux.Handle("POST /verify", verifyHandler)
	mux.Handle("GET /track", getPackageStatusHandler)
	mux.Handle("GET /getall", getAllOrderHandler)

	mux.Handle("GET /admin/get", adminGetHandler)
	mux.Handle("POST /admin/login", adminLoginHandler)
	mux.Handle("POST /admin/logout", adminLogoutHandler)
	mux.Handle("GET /admin/refresh", adminRefreshHandler)
	mux.Handle("GET /admin/getall", adminGetAllHandler)
	mux.Handle("DELETE /admin/user", adminDeleteHandler)
	mux.Handle("PUT /admin/user", adminCreateHandler)

	mux.Handle("GET /admin/order", adminGetOrderHandler)
	mux.Handle("GET /admin/orders", adminGetAllOrdersHandler)
	mux.Handle("DELETE /admin/order", adminDeleteOrderHandler)
	mux.Handle("POST /admin/order", adminUpdateOrderHandler)
	mux.Handle("GET /admin/package", adminGetPackageHandler)
	mux.Handle("DELETE /admin/package", adminDeletePackageHandler)
	mux.Handle("POST /admin/package", adminUpdatePackageHandler)
	mux.Handle("POST /admin/status", adminAddStatusHandler)
	mux.Handle("GET /admin/pricing", adminGetPricingHandler)
	mux.Handle("POST /admin/pricing", adminSetPricingHandler)

	//GET /admin/orders

	opt := cors.Options{
		AllowedOrigins: []string{configuration.Frontend.Host + ":" + strconv.Itoa(configuration.Frontend.Port)},
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
	_ = logger.Log("msg", "HTTP", "addr", ":"+strconv.Itoa(configuration.Port))
	_ = logger.Log("err", http.ListenAndServe(":"+strconv.Itoa(configuration.Port), handler))
}

func JWTParser(keyFunc stdjwt.Keyfunc, method stdjwt.SigningMethod, jwtRepository repository.JWTRepository, admin bool) endpoint.Middleware {
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
			if err != nil {
				return nil, err
			}
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
			if (admin && !slices.Contains(token.Claims.(*stdjwt.RegisteredClaims).Audience, "admin")) || !admin && !slices.Contains(token.Claims.(*stdjwt.RegisteredClaims).Audience, "user") {
				return nil, stdjwt.ErrTokenInvalidAudience
			}
			ctx = context.WithValue(ctx, jwt.JWTClaimsContextKey, token.Claims)

			return next(ctx, request)
		}
	}
}

func JSONErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {

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
		errors.Is(err, account.ErrAlreadyVerified) ||
		errors.Is(err, pricing.ErrInvalidValue) ||
		errors.Is(err, jwt.ErrTokenContextMissing) {
		status = http.StatusBadRequest
	}
	if errors.Is(err, auth.ErrInvalidCredentials) ||
		errors.Is(err, stdjwt.ErrTokenExpired) ||
		errors.Is(err, jwt.ErrTokenExpired) ||
		errors.Is(err, jwt.ErrTokenInvalid) ||
		errors.Is(err, stdjwt.ErrTokenInvalidAudience) {
		status = http.StatusUnauthorized

	}

	contentType, body := "application/json; charset=utf-8", "{\"error\": \""+err.Error()+"\" }"
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)
	_, _ = w.Write([]byte(body))
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
