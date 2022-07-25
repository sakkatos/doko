package utils

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds/rdsutils"
	"github.com/go-pg/pg/v10"
)

type (
	DBConnector struct {
		Session *session.Session
		Creds   *credentials.Credentials
		DB      *pg.DB
		Options *DBConnectorOptions
	}

	DBConnectorOptions struct {
		Region   string
		User     string
		Name     string
		Addr     string
		Password string
		Cert     string
		PoolSize int
	}
)

func InitDBConnector() *DBConnector {
	connector := &DBConnector{
		Options: &DBConnectorOptions{
			Region:   Getenv("AWS_REGION", "ap-southeast-1"),
			User:     Getenv("DB_USER", "postgres"),
			Name:     Getenv("DB_NAME", ""),
			Addr:     fmt.Sprintf("%s:5432", Getenv("DB_HOST", "")),
			Password: Getenv("DB_PASSWORD", ""),
			Cert:     Getenv("RDS_ROOT_CERT", ""),
			PoolSize: GetenvAsInt("DB_POOL_SIZE", 10),
		},
	}

	if connector.Options.Password == "" {
		connector.GetIAMSession()
		connector.Options.Password = GetRDSConnectToken(connector.Options, connector.Creds)
		connector.ConnectWithSSL()
	} else {
		connector.Connect()
	}

	return connector
}

func (c *DBConnector) GetIAMSession() {
	sess := session.Must(session.NewSession(
		&aws.Config{
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
	))

	c.Session = sess
	c.Creds = sess.Config.Credentials
}

func (c *DBConnector) ConnectWithSSL() {
	rootCerts := x509.NewCertPool()

	if ok := rootCerts.AppendCertsFromPEM([]byte(c.Options.Cert)); !ok {
		log.Fatal("Failed to append PEM.")
	}

	options := &pg.Options{
		User:     c.Options.User,
		Database: c.Options.Name,
		Addr:     c.Options.Addr,
		Password: c.Options.Password,
		PoolSize: c.Options.PoolSize,
		TLSConfig: &tls.Config{
			RootCAs:    rootCerts,
			ServerName: Getenv("DB_HOST", ""),
		},
	}

	c.DB = pg.Connect(options)
}

func (c *DBConnector) Connect() {
	options := &pg.Options{
		User:     c.Options.User,
		Database: c.Options.Name,
		Addr:     c.Options.Addr,
		Password: c.Options.Password,
	}

	c.DB = pg.Connect(options)
}

func MigrationDBConnect() *pg.DB {
	opts := &DBConnectorOptions{
		User:     Getenv("DB_USER", ""),
		Name:     Getenv("DB_NAME", ""),
		Addr:     fmt.Sprintf("%s:5432", Getenv("DB_HOST", "")),
		Password: Getenv("DB_PASSWORD", ""),
		Cert:     Getenv("RDS_ROOT_CERT", ""),
	}

	connector := &DBConnector{
		Options: opts,
	}

	if opts.Cert != "" {
		connector.ConnectWithSSL()
	} else {
		connector.Connect()
	}

	return connector.DB
}

func GetRDSConnectToken(opts *DBConnectorOptions, creds *credentials.Credentials) string {
	if creds.IsExpired() {
		creds.Get()
	}

	authToken, err := rdsutils.BuildAuthToken(opts.Addr, opts.Region, opts.User, creds)

	if err == nil {
		log.Println("DB authorization token requested successfully")
	} else {
		log.Println(err)
	}

	return authToken
}
