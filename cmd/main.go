package main

import (
	"context"
	"database/sql"
	"fmt"
	log2 "log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	awsv1 "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sesv2"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/caarlos0/env/v6"
	"github.com/dgrijalva/jwt-go"
	jwt2 "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	http2 "github.com/go-kit/kit/transport/http"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go/v72/client"
	mysql2 "gorm.io/driver/mysql"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"

	auth2 "github.com/Equiphunter-com/equipment-hunter-api/pkg/auth"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/auth/jwt/claims"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/auth/jwt/hmac"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/checkout"
	cities2 "github.com/Equiphunter-com/equipment-hunter-api/pkg/cities"
	companies2 "github.com/Equiphunter-com/equipment-hunter-api/pkg/companies"
	repositories2 "github.com/Equiphunter-com/equipment-hunter-api/pkg/companies/repositories"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/coupons"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/email"
	equipmentcategories2 "github.com/Equiphunter-com/equipment-hunter-api/pkg/equipmentcategories"
	equipments2 "github.com/Equiphunter-com/equipment-hunter-api/pkg/equipments"
	equipmentsubcategories2 "github.com/Equiphunter-com/equipment-hunter-api/pkg/equipmentsubcategories"
	media2 "github.com/Equiphunter-com/equipment-hunter-api/pkg/media"
	proposals2 "github.com/Equiphunter-com/equipment-hunter-api/pkg/proposals"
	quotes2 "github.com/Equiphunter-com/equipment-hunter-api/pkg/quotes"
	repositories4 "github.com/Equiphunter-com/equipment-hunter-api/pkg/quotes/repositories"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/repositories/gormclient"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/repositories/mysql"
	ses2 "github.com/Equiphunter-com/equipment-hunter-api/pkg/ses"
	sublists2 "github.com/Equiphunter-com/equipment-hunter-api/pkg/sublists"
	repositories3 "github.com/Equiphunter-com/equipment-hunter-api/pkg/sublists/repositories"
	supplycategories2 "github.com/Equiphunter-com/equipment-hunter-api/pkg/supplycategories"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/auth"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/cities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/companies"
	coupons2 "github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/coupons"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/equipmentcategories"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/equipments"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/equipmentsubcategories"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/media"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/proposals"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/quotes"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/sublists"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/supplycategories"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/users"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/vendorpurchases"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/handlers"
	users2 "github.com/Equiphunter-com/equipment-hunter-api/pkg/users"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/users/repositories"
)

type AWSLamdaHandler func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

type DBConfig struct {
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PASSWORD,required"`
	Host     string `env:"DB_HOST,required"`
	Port     string `env:"DB_PORT,required"`
	Name     string `env:"DB_NAME,required"`
}

type AWSBucket struct {
	Name string `env:"AWS_BUCKET,required"`
}

type AuthConfig struct {
	HashCost             int    `env:"HASH_COST" envDefault:"10"`
	JWTKey               string `env:"JWT_KEY,required"`
	JWTExpiration        int64  `env:"JWT_EXPIRATION,required"`
	JWTRefreshExpiration int64  `env:"JWT_REFRESH,required"`
	FrontURL             string `env:"FRONT_URL,required"`
	StripeKey            string `env:"STRIPE_KEY,required"`
}

func makeHandler() AWSLamdaHandler {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		var logger log.Logger
		{
			logger = log.NewLogfmtLogger(os.Stderr)
			logger = log.With(logger, "ts", log.DefaultTimestampUTC)
			logger = log.With(logger, "caller", log.DefaultCaller)
		}
		homeHandler := func(w http.ResponseWriter, req *http.Request) {
			logger.Log("base", "test")
			w.Header().Add("unfortunately-required-header", "")
		}
		r := mux.NewRouter()
		r.Use(transport.CORSHeadersPolicies([]string{http.MethodOptions, http.MethodPost, http.MethodGet, http.MethodPut, http.MethodPatch}))
		db, err := buildDB()
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
		}
		defer func() {
			_ = db.Close()
		}()

		newLogger := logger2.New(
			log2.New(os.Stdout, "\r\n", log2.LstdFlags), // io writer
			logger2.Config{
				SlowThreshold: time.Second,    // Slow SQL threshold
				LogLevel:      logger2.Silent, // Log level
				Colorful:      false,          // Disable color
			},
		)
		gormDB, err := gorm.Open(mysql2.New(mysql2.Config{
			Conn: db,
		}), &gorm.Config{Logger: newLogger})
		authConfig, err := buildAuthConfig()
		if err != nil {
			logger.Log("err", fmt.Sprintf("authConfig:%v", err))
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
		}
		gormDB = gormDB.Debug()
		r.HandleFunc("/", homeHandler)
		r = addSublistHandler(r, logger, db)
		r = addSupplyCategoriesHandler(r, logger, db, gormDB)
		r = addEquipmentCategoryHandler(r, logger, db)
		r = addQuotesHandler(r, logger, db)
		r = addSingleQuotesHandler(r, logger, gormDB)
		r = addSubEquipmentCategoryHandler(r, logger, db)
		r = addMediaHandler(r, logger, db)
		r = addAuthHandler(r, logger, gormDB, authConfig)
		r = addUserHandler(r, logger, gormDB, authConfig)
		r = addCompaniesHandler(r, logger, gormDB)
		r = addVendorPurchasesHandler(r, logger, gormDB, authConfig, db)
		r = addCouponsHandler(r, logger, gormDB, authConfig)
		r = addCitiesHandler(r, logger, gormDB)
		r = addSublistUserHandler(r, logger, gormDB, authConfig)
		r = addEquipmentUserHandler(r, logger, gormDB, authConfig)
		r = addUserQuotesHandler(r, logger, gormDB, authConfig)
		r = addEquipmentHandler(r, logger, gormDB)
		r = addUserSublistHandler(r, logger, gormDB, authConfig)
		r = addProposalHandler(r, logger, gormDB, authConfig)
		r = addBuyerProposalListHandler(r, logger, gormDB, authConfig)
		r = addSellerProposalListHandler(r, logger, gormDB, authConfig)
		r = addProposalListHandler(r, logger, db, authConfig)
		adapter := gorillamux.New(r)

		return adapter.ProxyWithContext(ctx, req)
	}
}

func addUserSublistHandler(r *mux.Router, logger log.Logger, db *gorm.DB, config *AuthConfig) *mux.Router {
	authMiddleware := makeAuthMiddleware([]byte(config.JWTKey))
	svc := sublists2.NewUserSublist(logger, db)
	ep := sublists.MakeUserSubEndpoint(svc)
	return handlers.MakeUserSubHandler(r, ep, authMiddleware, []http2.ServerOption{http2.ServerBefore(jwt2.HTTPToContext())})
}

func addEquipmentHandler(r *mux.Router, logger log.Logger, db *gorm.DB) *mux.Router {
	svc := equipments2.NewGetEquipmentService(logger, db)
	ep := equipments.MakeEquipmentEndpoint(svc)
	return handlers.MakeEquipmentHandler(r, ep)
}

func addProposalHandler(r *mux.Router, logger log.Logger, db *gorm.DB, config *AuthConfig) *mux.Router {
	authMiddleware := makeAuthMiddleware([]byte(config.JWTKey))
	senderEmail := "support@equiphunter.com"
	cfg, _ := external.LoadDefaultAWSConfig()
	cfg.Region = "us-east-1"
	cfg.Credentials = aws.NewStaticCredentialsProvider("111222223333333", "111222223333333", "")
	sesv2Cfg := awsv1.Config{
		Region:      &cfg.Region,
		Credentials: credentials.NewStaticCredentials("111222223333333", "111222223333333", ""),
	}
	sess := session.Must(session.NewSession(&sesv2Cfg))
	awsSes := ses.New(cfg)
	awsSesv2 := sesv2.New(sess)
	sesClient := ses2.NewClient(senderEmail, awsSes, awsSesv2, logger)
	s := email.NewSender(logger, *sesClient, config.FrontURL)
	svce := proposals2.NewSellerService(logger, db, *s)
	svcb := proposals2.NewBuyerService(logger, db, *s)
	eps := proposals.NewEndpoints(svce, svcb)
	r = handlers.MakeBuyerProposalHandler(r, *eps, authMiddleware, []http2.ServerOption{http2.ServerBefore(jwt2.HTTPToContext())})
	return handlers.MakeUserProposalHandler(r, *eps, authMiddleware, []http2.ServerOption{http2.ServerBefore(jwt2.HTTPToContext())})
}

func addBuyerProposalListHandler(r *mux.Router, logger log.Logger, db *gorm.DB, config *AuthConfig) *mux.Router {
	authMiddleware := makeAuthMiddleware([]byte(config.JWTKey))
	senderEmail := "support@equiphunter.com"
	cfg, _ := external.LoadDefaultAWSConfig()
	cfg.Region = "us-east-1"
	cfg.Credentials = aws.NewStaticCredentialsProvider("111222223333333", "111222223333333", "")
	sesv2Cfg := awsv1.Config{
		Region:      &cfg.Region,
		Credentials: credentials.NewStaticCredentials("111222223333333", "111222223333333", ""),
	}
	sess := session.Must(session.NewSession(&sesv2Cfg))
	awsSes := ses.New(cfg)
	awsSesv2 := sesv2.New(sess)
	sesClient := ses2.NewClient(senderEmail, awsSes, awsSesv2, logger)
	s := email.NewSender(logger, *sesClient, config.FrontURL)
	svc := proposals2.NewBuyerService(logger, db, *s)
	eps := proposals.MakeBuyerProposalEndpoint(svc)
	return handlers.MakeBuyProposalListHandler(r, eps, authMiddleware, []http2.ServerOption{http2.ServerBefore(jwt2.HTTPToContext())})
}
func addSellerProposalListHandler(r *mux.Router, logger log.Logger, db *gorm.DB, config *AuthConfig) *mux.Router {
	authMiddleware := makeAuthMiddleware([]byte(config.JWTKey))
	senderEmail := "support@equiphunter.com"
	cfg, _ := external.LoadDefaultAWSConfig()
	cfg.Region = "us-east-1"
	cfg.Credentials = aws.NewStaticCredentialsProvider("111222223333333", "111222223333333", "")
	sesv2Cfg := awsv1.Config{
		Region:      &cfg.Region,
		Credentials: credentials.NewStaticCredentials("111222223333333", "111222223333333", ""),
	}
	sess := session.Must(session.NewSession(&sesv2Cfg))
	awsSes := ses.New(cfg)
	awsSesv2 := sesv2.New(sess)
	sesClient := ses2.NewClient(senderEmail, awsSes, awsSesv2, logger)
	s := email.NewSender(logger, *sesClient, config.FrontURL)
	svc := proposals2.NewSellerService(logger, db, *s)
	eps := proposals.MakeSellerProposalEndpoint(svc)
	return handlers.MakeSellerProposalListHandler(r, eps, authMiddleware, []http2.ServerOption{http2.ServerBefore(jwt2.HTTPToContext())})
}

func addUserQuotesHandler(r *mux.Router, logger log.Logger, db *gorm.DB, config *AuthConfig) *mux.Router {
	authMiddleware := makeAuthMiddleware([]byte(config.JWTKey))
	svc := quotes2.NewUserQuotes(logger, db)
	eps := quotes.MakeUserQuoteEndpoint(svc)
	return handlers.MakeUserQuotesHandler(r, eps, authMiddleware, []http2.ServerOption{http2.ServerBefore(jwt2.HTTPToContext())})
}

func addEquipmentUserHandler(r *mux.Router, logger log.Logger, db *gorm.DB, config *AuthConfig) *mux.Router {
	authMiddleware := makeAuthMiddleware([]byte(config.JWTKey))
	senderEmail := "support@equiphunter.com"
	cfg, _ := external.LoadDefaultAWSConfig()
	cfg.Region = "us-east-1"
	cfg.Credentials = aws.NewStaticCredentialsProvider("111222223333333", "111222223333333", "")
	sesv2Cfg := awsv1.Config{
		Region:      &cfg.Region,
		Credentials: credentials.NewStaticCredentials("111222223333333", "111222223333333", ""),
	}
	sess := session.Must(session.NewSession(&sesv2Cfg))
	awsSes := ses.New(cfg)
	awsSesv2 := sesv2.New(sess)
	sesClient := ses2.NewClient(senderEmail, awsSes, awsSesv2, logger)
	s := email.NewSender(logger, *sesClient, config.FrontURL)
	svcs := quotes2.NewSupplyService(logger, db, *s)
	svce := quotes2.NewEquipmentService(logger, db, *s)
	eps := quotes.NewEndpoints(svce, svcs)
	r = handlers.MakeUserEquipmentHandler(r, *eps, authMiddleware, []http2.ServerOption{http2.ServerBefore(jwt2.HTTPToContext())})
	return handlers.MakeUserSupplyHandler(r, *eps, authMiddleware, []http2.ServerOption{http2.ServerBefore(jwt2.HTTPToContext())})
}

func addSublistUserHandler(r *mux.Router, logger log.Logger, db *gorm.DB, config *AuthConfig) *mux.Router {
	authMiddleware := makeAuthMiddleware([]byte(config.JWTKey))
	senderEmail := "support@equiphunter.com"
	cfg, _ := external.LoadDefaultAWSConfig()
	cfg.Region = "us-east-1"
	cfg.Credentials = aws.NewStaticCredentialsProvider("111222223333333", "111222223333333", "")
	sesv2Cfg := awsv1.Config{
		Region:      &cfg.Region,
		Credentials: credentials.NewStaticCredentials("111222223333333", "111222223333333", ""),
	}
	sess := session.Must(session.NewSession(&sesv2Cfg))
	awsSes := ses.New(cfg)
	awsSesv2 := sesv2.New(sess)
	sesClient := ses2.NewClient(senderEmail, awsSes, awsSesv2, logger)
	s := email.NewSender(logger, *sesClient, config.FrontURL)
	svc := sublists2.NewServices(logger, db, *s)
	eps := sublists.NewEndpoints(svc)
	return handlers.MakeUserSublistHandler(r, eps, authMiddleware, []http2.ServerOption{http2.ServerBefore(jwt2.HTTPToContext())})
}

func addCitiesHandler(r *mux.Router, logger log.Logger, db *gorm.DB) *mux.Router {
	svc := cities2.NewGetCitiesService(logger, db)
	ce := cities.MakeGetCitiesEndpoint(svc)
	coe := cities.MakeGetCountiesEndpoint(svc)
	se := cities.MakeGetStatesEndpoint(svc)
	eps := cities.NewEndpoints(ce, coe, se)
	return handlers.MakeCitiesHandler(r, eps)
}

func addVendorPurchasesHandler(r *mux.Router, logger log.Logger, gormDB *gorm.DB, config *AuthConfig, db *sql.DB) *mux.Router {
	authMiddleware := makeAuthMiddleware([]byte(config.JWTKey))
	qr := gormclient.NewQuotesRepository(logger, gormDB)
	srr := repositories4.NewSupplyRequestRepository(logger, gormDB)
	eqrr := repositories4.NewEquipmentRequestRepository(logger, gormDB)
	mqr := mysql.NewQuotesRepositories(logger, db)
	srsvc := users2.NewVendorSupplyRequestsService(logger, qr, srr, mqr)
	ersvc := users2.NewEquipmentRequestsService(qr, logger, eqrr, mqr)
	sr := mysql.NewSublistsRepository(logger, db)
	gsr := repositories3.NewSublistRepository(logger, gormDB)
	ssvc := users2.NewGetVendorSublists(sr, logger, gsr)
	sc := buildStripeClient(config.StripeKey)
	senderEmail := "support@equiphunter.com"
	cfg, _ := external.LoadDefaultAWSConfig()
	cfg.Region = "us-east-1"
	cfg.Credentials = aws.NewStaticCredentialsProvider("111222223333333", "111222223333333", "")
	sesv2Cfg := awsv1.Config{
		Region:      &cfg.Region,
		Credentials: credentials.NewStaticCredentials("111222223333333", "111222223333333", ""),
	}
	sess := session.Must(session.NewSession(&sesv2Cfg))
	awsSes := ses.New(cfg)
	awsSesv2 := sesv2.New(sess)
	sesClient2 := ses2.NewClient(senderEmail, awsSes, awsSesv2, logger)
	s := email.NewSender(logger, *sesClient2, config.FrontURL)

	ce := checkout.NewService(logger, db, gormDB, sc, *s)
	eps := vendorpurchases.NewEndpoints(ssvc, ersvc, srsvc, ce)
	return handlers.MakeVendorPurchasesHandler(r, eps, authMiddleware, []http2.ServerOption{http2.ServerBefore(jwt2.HTTPToContext())})
}

func addCouponsHandler(r *mux.Router, logger log.Logger, db *gorm.DB, config *AuthConfig) *mux.Router {
	authMiddleware := makeAuthMiddleware([]byte(config.JWTKey))
	svc := coupons.NewService(logger, db)
	eps := coupons2.NewEndpoints(svc)
	return handlers.MakeCouponsHandler(r, eps, authMiddleware, []http2.ServerOption{http2.ServerBefore(jwt2.HTTPToContext())})
}

func addCompaniesHandler(r *mux.Router, logger log.Logger, db *gorm.DB) *mux.Router {
	cr := repositories2.NewListCompaniesRepository(logger, db)
	lsvc := companies2.NewListService(logger, cr)
	svc := companies.NewServices(lsvc)
	eps := companies.NewEndpoints(svc)
	return handlers.MakeCompaniesHandler(r, eps)
}

func addUserHandler(r *mux.Router, logger log.Logger, db *gorm.DB, config *AuthConfig) *mux.Router {
	senderEmail := "support@equiphunter.com"
	cfg, _ := external.LoadDefaultAWSConfig()
	cfg.Region = "us-east-1"
	cfg.Credentials = aws.NewStaticCredentialsProvider("111222223333333", "111222223333333", "")
	sesv2Cfg := awsv1.Config{
		Region:      &cfg.Region,
		Credentials: credentials.NewStaticCredentials("111222223333333", "111222223333333", ""),
	}
	sess := session.Must(session.NewSession(&sesv2Cfg))
	awsSes := ses.New(cfg)
	awsSesv2 := sesv2.New(sess)
	snsClient := sns.New(cfg)
	sesClient2 := ses2.NewClient(senderEmail, awsSes, awsSesv2, logger)
	s := email.NewSender(logger, *sesClient2, config.FrontURL)
	hs := auth2.NewHasher(config.HashCost)
	ur := mysql.NewUser(logger, db)
	ccr := repositories.NewCouponRepository(logger, db)
	cvcr := repositories.NewVerificationCodeRepository(logger, db)
	n := users2.NewUserNotifier(logger, db, cvcr, ccr, snsClient, *s)
	iss := hmac.NewIssuer([]byte(config.JWTKey), config.JWTExpiration, config.JWTRefreshExpiration)
	ls := auth2.NewLoginService(ur, hs, iss)
	ucr := repositories.NewCreate(logger, db)
	csvc := users2.NewCreateService(logger, ucr, ls, hs, cvcr, n, *s)
	vs := users2.NewVerifyService(logger, cvcr, db, s)
	usvc := users2.NewResourceService(logger, db, hs)
	svc := users.NewServices(csvc, vs, usvc, usvc)
	eps := users.NewEndpoints(svc)
	authMiddleware := makeAuthMiddleware([]byte(config.JWTKey))
	return handlers.MakeUserHandler(r, eps, authMiddleware, []http2.ServerOption{http2.ServerBefore(jwt2.HTTPToContext())})
}

func makeAuthMiddleware(key []byte) endpoint.Middleware {
	return jwt2.NewParser(func(token *jwt.Token) (interface{}, error) {
		return key, nil
	}, jwt.SigningMethodHS256, func() jwt.Claims {
		return &claims.UserClaims{}
	})
}

func buildAuthConfig() (*AuthConfig, error) {
	config := &AuthConfig{}
	err := env.Parse(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func addAuthHandler(r *mux.Router, logger log.Logger, db *gorm.DB, config *AuthConfig) *mux.Router {
	ur := mysql.NewUser(logger, db)
	hs := auth2.NewHasher(config.HashCost)
	iss := hmac.NewIssuer([]byte(config.JWTKey), config.JWTExpiration, config.JWTRefreshExpiration)
	lsvc := auth2.NewLoginService(ur, hs, iss)
	svc := auth.Services{Login: lsvc}
	eps := auth.NewAuth(svc)
	return handlers.MakeAuthHandler(r, eps)
}

func addSubEquipmentCategoryHandler(r *mux.Router, logger log.Logger, db *sql.DB) *mux.Router {
	escr := mysql.NewEquipmentSubcategoriesRepository(logger, db)
	svc := equipmentsubcategories2.NewGetService(logger, escr)
	esce := equipmentsubcategories.MakeGetEquipmentSubcategoriesEndpoint(svc)
	return handlers.MakeEquipmentSubcategoriesHandler(r, esce)
}

func addProposalListHandler(r *mux.Router, logger log.Logger, db *sql.DB, config *AuthConfig) *mux.Router {
	authMiddleware := makeAuthMiddleware([]byte(config.JWTKey))
	qr := mysql.NewProposalsRepositories(logger, db)
	qs := proposals2.NewListService(logger, qr)
	qe := proposals.MakeListProposalsEndpoint(qs)
	return handlers.MakeProposalListHandler(r, qe, authMiddleware, []http2.ServerOption{http2.ServerBefore(jwt2.HTTPToContext())})
}

func addQuotesHandler(r *mux.Router, logger log.Logger, db *sql.DB) *mux.Router {
	qr := mysql.NewQuotesRepositories(logger, db)
	qs := quotes2.NewListService(logger, qr)
	qe := quotes.MakeListQuotesEndpoint(qs)
	return handlers.MakeQuotesHandler(r, qe)
}

func addSingleQuotesHandler(r *mux.Router, logger log.Logger, db *gorm.DB) *mux.Router {
	sqr := repositories4.NewSingleQuoteRepository(logger, db)
	sqs := quotes2.NewSingleQuoteService(logger, sqr)
	svc := quotes.NewServices(sqs)
	eps := quotes.NewSingleEndpoints(svc)
	return handlers.MakeSingleQuoteHandler(r, eps)
}

func addSupplyCategoriesHandler(r *mux.Router, logger log.Logger, db *sql.DB, gdb *gorm.DB) *mux.Router {
	scr := mysql.NewSupplyCategoriesRepository(logger, db)
	svc := supplycategories2.NewListService(logger, scr, gdb)
	sce := supplycategories.MakeListSupplyCategoriesEndpoint(svc)
	se := supplycategories.MakeSuppliesEndpoint(svc)
	r = handlers.MakeSupplyCategoriesHandler(r, sce)
	return handlers.MakeSuppliesHandler(r, se)
}

func addSublistHandler(r *mux.Router, logger log.Logger, db *sql.DB) *mux.Router {
	sr := mysql.NewSublistsRepository(logger, db)
	lss := sublists2.NewListSublistsService(logger, sr)
	gse := sublists.MakeListSublistsEndpoint(lss)
	return handlers.MakeGetSublistHandler(r, gse)
}

func addEquipmentCategoryHandler(r *mux.Router, logger log.Logger, db *sql.DB) *mux.Router {
	ecr := mysql.NewEquipmentCategoriesRepository(logger, db)
	lecs := equipmentcategories2.NewListService(logger, ecr)
	ece := equipmentcategories.MakeListEndpoint(lecs)
	return handlers.MakeEquipmentCategoriesHandler(r, ece)
}

func addMediaHandler(r *mux.Router, logger log.Logger, db *sql.DB) *mux.Router {
	bucketConfig, _ := buildAWSBucket()
	ar := mysql.NewMediaRepository(db, logger)
	svc := media2.NewGetAboutUsService(logger, ar, bucketConfig.Name)
	ep := media.MakeAboutUsEndpoint(svc)
	return handlers.MakeMediaHandler(r, ep)
}

func (dbc DBConfig) BuildConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbc.User, dbc.Password, dbc.Host, dbc.Port, dbc.Name)
}

func buildDB() (*sql.DB, error) {
	dbConfig := DBConfig{}
	err := env.Parse(&dbConfig)
	if err != nil {
		return nil, err
	}
	return sql.Open("mysql", dbConfig.BuildConnectionString())
}

func buildAWSBucket() (*AWSBucket, error) {
	awsBucket := AWSBucket{}
	err := env.Parse(&awsBucket)
	if err != nil {
		return nil, err
	}
	return &awsBucket, nil
}

func main() {
	fmt.Println("==================================================================")
	lambda.Start(makeHandler())
}

func buildStripeClient(key string) *client.API {
	sc := &client.API{}
	sc.Init(key, nil)
	return sc
}
