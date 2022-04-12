package cmd

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// migrateCmd represents the serve command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: migrate,
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	migrateCmd.PersistentFlags().StringP("dir", "d", "manifests", "Directory containing all manifests")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type Table struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string
	Metadata   TableMetadata
	Spec       TableSpec
}

type TableMetadata struct {
	Name string
}

type TableSpec struct {
	Columns []TableColumn
}

type TableColumn struct {
	Name string
	Type string

	PrimaryKey bool `yaml:"primaryKey"`
	Required   bool
}

func migrate(cmd *cobra.Command, _ []string) {
	manDir, err := cmd.PersistentFlags().GetString("dir")
	if err != nil {
		panic(err)
	}

	files, err := ioutil.ReadDir(manDir)
	if err != nil {
		panic(err)
	}

	var resources []*Table

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if strings.HasSuffix(f.Name(), ".yaml") || strings.HasSuffix(f.Name(), ".yml") {
			contents, err := ioutil.ReadFile(path.Join(manDir, f.Name()))
			if err != nil {
				panic(err)
			}

			resource := Table{}
			err = yaml.Unmarshal(contents, &resource)
			if err != nil {
				panic(err)
			}
			resources = append(resources, &resource)
		}
	}

	fmt.Printf("%+v\n", resources[0])

	LoadConfig()

	dsn := viper.GetString("dsn")
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal(err)
	}

	for _, tb := range resources {
		_, err = db.Exec(tableToCreateQuery(tb))
		if err != nil {
			panic(err)
		}
	}
}

func columnToSQLDefinition(col TableColumn) string {
	format := fmt.Sprintf("%s %s ", col.Name, col.Type)
	if col.Required {
		format += "not null "
	}
	if col.PrimaryKey {
		format += "primary key "
	}
	return format
}

func tableToCreateQuery(tb *Table) string {
	var cols []string

	for _, col := range tb.Spec.Columns {
		cols = append(cols, columnToSQLDefinition(col))
	}

	colsDef := strings.Join(cols, ",")
	s := fmt.Sprintf("create table if not exists `%s` (%s)", tb.Metadata.Name, colsDef)
	return s
}
