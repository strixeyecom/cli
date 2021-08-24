package repository

import (
	`context`
	`fmt`
	`io`
	"log"
	`net`
	"os"
	`strings`
	"time"
	
	`github.com/docker/docker/api/types`
	`github.com/docker/docker/api/types/container`
	`github.com/docker/docker/api/types/filters`
	`github.com/docker/docker/client`
	`github.com/docker/docker/pkg/stdcopy`
	`github.com/docker/go-connections/nat`
	`github.com/pkg/errors`
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	
	`github.com/strixeyecom/cli/domain/repository`
)

/*
	Created by aomerk at 5/21/21 for project cli
*/

/*
	repository for database controls
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

// ConnectToAgentDB establishes a live connection to your agents database.
// Make sure to have permissions and network configurations so that use can connect to database. Usually,
// database ports and hosts are not public in enterprise networks. So, that part is on you to check.
func ConnectToAgentDB(dbConfig repository.Database) (*gorm.DB, error) {
	// orm logger config
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)
	
	// establish connection
	dsn := dbConfig.DSN()
	db, err := gorm.Open(
		mysql.New(
			mysql.Config{
				DSN:                       dsn,   // data source name
				DefaultStringSize:         256,   // default size for string fields
				DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
				DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
				DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
				SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
			},
		), &gorm.Config{
			Logger: newLogger,
		},
	)
	
	if err != nil {
		return nil, err
	}
	
	// handle possible errors.
	
	// one error is that we are using stack network hostnames like "db" "database" instead of ips or
	// resolvable domains
	_, err = net.LookupHost(dbConfig.DBAddr)
	if err != nil {
		return nil, errors.Wrap(
			err, `
db_addr field in your in your database config is not a resolvable hostname.

You can override it from your config. Check out https://docs.strixeye.com/cli/configuration#override`,
		)
	}
	return db, err
}

// CreateDatabase creates a temporary MySQL Database on a Docker container
func CreateDatabase(cliConfig repository.Database) error {
	const (
		imageName = "mysql"
		mysqlPort = "3306"
	)
	
	var (
		mysqlUser     = fmt.Sprintf("MYSQL_USER=%s", cliConfig.DBUser)
		mysqlPass     = fmt.Sprintf("MYSQL_PASSWORD=%s", cliConfig.DBPass)
		mysqlRootPass = fmt.Sprintf("MYSQL_ROOT_PASSWORD=%s", cliConfig.DBPass)
		mysqlDatabase = fmt.Sprintf("MYSQL_DATABASE=%s", cliConfig.DBName)
	)
	
	ctx := context.Background()
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	
	reader, err := dockerClient.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	
	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		return err
	}
	
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			mysqlPort: []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: cliConfig.DBPort,
				},
			},
		},
	}
	
	resp, err := dockerClient.ContainerCreate(
		ctx, &container.Config{
			Healthcheck: &container.HealthConfig{
				Interval: time.Second,
				Retries:  10,
				Test: []string{
					"CMD", "mysqladmin", "ping", "-u", "root",
					fmt.Sprintf("--password=%s", cliConfig.DBPass),
				},
			},
			Image: imageName,
			Tty:   false,
			ExposedPorts: nat.PortSet{
				mysqlPort: struct{}{},
			}, Env: []string{
				mysqlUser, mysqlPass, mysqlDatabase, mysqlRootPass,
			},
		}, hostConfig, nil, nil, cliConfig.TestContainerName(),
	)
	if err != nil {
		return err
	}

	err = dockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	filterArgs := filters.NewArgs(
		filters.KeyValuePair{Key: "container", Value: resp.ID},
		filters.KeyValuePair{Key: "event", Value: "health_status"},
	)
	statusCh, errCh := dockerClient.Events(
		ctx, types.EventsOptions{
			Filters: filterArgs,
		},
	)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case status := <-statusCh:
		if strings.Contains(status.Action, "healthy") {
			break
		}
	}
	
	out, err := dockerClient.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return err
	}
	
	_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	if err != nil {
		return err
	}
	
	// although I use everything I can, like, health checks, event monitors and so on,
	// there are still some problems while connecting immediately to database
	time.Sleep(time.Second * 5)
	
	return nil
}

// SetupDatabase implements table creation/migration on given database
func SetupDatabase(dbConfig repository.Database) error {
	var err error
	db, err := ConnectToAgentDB(dbConfig)
	for err != nil {
		db, err = ConnectToAgentDB(dbConfig)
		time.Sleep(time.Second * 2)
	}
	
	err = db.AutoMigrate(repository.Suspicion{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(repository.Suspect{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(repository.Trip{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(repository.Request{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(repository.Header{})
	if err != nil {
		return err
	}
	return nil
}

// RemoveDatabase removes created docker database container
func RemoveDatabase(dbConfig repository.Database) error {
	var (
		err error
	)

	ctx := context.Background()
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	// Find StrixEye Database container
	containers, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return err
	}
	var strixeyeContainerID string

	// Iterate through containers and find our one.
	for _, cont := range containers {
		for _, name := range cont.Names {
			if strings.Contains(name, dbConfig.TestContainerName()) {
				strixeyeContainerID = cont.ID
				break
			}
		}
	}

	// can't find container
	if strixeyeContainerID == "" {
		return err
	}

	// container found, stopping container
	timeout := time.Second * 5
	err = dockerClient.ContainerStop(ctx, strixeyeContainerID, &timeout)
	if err != nil {
		return err
	}
	// remove container and data
	err = dockerClient.ContainerRemove(ctx, strixeyeContainerID, types.ContainerRemoveOptions{})
	if err != nil {
		return err
	}
	return nil
}

// CreateDatabaseIFNotExists tries to connect to given database and creates a docker container if
// connection fails.
func CreateDatabaseIFNotExists(database repository.Database) error {
	var err error
	
	// Try to connect to existing database
	_, err = ConnectToAgentDB(database)
	if err == nil {
		return nil
	}
	
	// 	Create if non exists
	err = CreateDatabase(database)
	if err != nil {
		return err
	}
	
	return nil
}
