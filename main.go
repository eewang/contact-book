package main

import (
  "bytes"
  "fmt"
  "net/http"
  "database/sql"
  "encoding/json"

  "github.com/go-martini/martini"
  _ "github.com/go-sql-driver/mysql" // _ is a blank identifier that is used if you need a Go package to run some code upon init but don't need to use its functionality
  "github.com/joho/godotenv"
  "github.com/martini-contrib/binding"
  "upper.io/db"
  "upper.io/db/mysql"
  "upper.io/db/util/sqlutil"
)

var (
  qh = NewQueryHelper()
)

type EnvMap struct {
  Vars map[string]string
}

type QueryHelp struct {
  Rows            *sql.Rows
  Err             error
  DefaultRowCount string
  drv             *sql.DB
  Sess            db.Database
  dbSettings      db.Settings
  SQLString       bytes.Buffer
}

func NewEnvMap() *EnvMap {
  envMap := &EnvMap{}
  eMap, err := godotenv.Read()

  if err != nil {
    fmt.Printf("No configuration defined")
    return nil
  }

  envMap.Vars = eMap
  return envMap
}

func NewQueryHelper() *QueryHelp {
  envMap := NewEnvMap()

  settings := db.Settings{
    Host: envMap.Vars["MYSQL_HOST"],
    Database: envMap.Vars["MYSQL_DATABASE"],
    User: envMap.Vars["MYSQL_USERNAME"],
    Password: envMap.Vars["MYSQL_PASSWORD"],
  }

  return &QueryHelp{
    dbSettings: settings,
    DefaultRowCount: "10",
  }
}

func (q *QueryHelp) Query() *sql.Rows {
  q.Sess, q.Err = db.Open(mysql.Adapter, q.dbSettings)

  if q.Err != nil {
    return q.BadQuery()
  }
  defer q.Sess.Close()

  q.drv = q.Sess.Driver().(*sql.DB)
  q.Rows, q.Err = q.drv.Query(q.SQLString.String())

  if q.Err != nil {
    return q.BadQuery()
  }

  return q.Rows
}

func (q *QueryHelp) BadQuery() *sql.Rows {
  fmt.Printf("\n\n Error: %#v \n Query: %#v \n\n", q.Err, q.SQLString.String())
  q.Rows = nil
  return nil
}

func IndexGroups(params martini.Params, res http.ResponseWriter, req *http.Request) []byte {
  var groups []Group
  ret := GroupApiResponse{}

  qh.SQLString.Reset();
  qh.SQLString.WriteString("SELECT * FROM groups;")
  
  qh.Err = sqlutil.FetchRows(qh.Query(), &groups)

  if qh.Err != nil {
    fmt.Printf("Something went wrong with the query")
  }

  fmt.Printf("\n %v \n", groups)

  for _, group := range groups {
    ret.Response = append(ret.Response, group)
  }

  ret.Success = true
  ret.Message = "Here are the groups"

  responseJson, err := json.MarshalIndent(ret, "", "  ")
  if err != nil {
    fmt.Printf("Something went wrong with json.MarshalIndent")
    return []byte("")
  }
  return responseJson
}

type Group struct {
  Id   int    `db:"id"`
  Name string `db:"name"`
}

type GroupApiResponse struct {
  Success bool
  Message string
  Response []Group
}

// NOTE: struct attributes MUST be capitalized to appear in the JSON
// 
type Person struct {
  Id    int      `db:"id"`
  Name  string   `db:"name" form:"name" binding:"required"`
  Notes string   `db:"notes"`
  GroupId string `db:"group_id" form:"group_id" binding:"required"`
}

type PersonApiResponse struct {
  Success bool
  Message string
  Response []Person
}

func puts(text string) {
  fmt.Printf("\n %v \n", text)
}

func IndexPersons(params martini.Params, res http.ResponseWriter, req *http.Request) []byte {
  var people []Person
  ret := PersonApiResponse{}

  qh.SQLString.Reset();
  qh.SQLString.WriteString("SELECT * FROM persons;")
  
  qh.Err = sqlutil.FetchRows(qh.Query(), &people)

  if qh.Err != nil {
    fmt.Printf("Something went wrong with the query")
  }

  fmt.Printf("\n %v \n", people)

  for _, person := range people {
    ret.Response = append(ret.Response, person)
  }

  ret.Success = true
  ret.Message = "Here are the people"

  responseJson, err := json.MarshalIndent(ret, "", "  ")
  if err != nil {
    fmt.Printf("Something went wrong with json.MarshalIndent")
    return []byte("")
  }
  return responseJson
}

func HandleStatic(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("\n %v \n", req.URL.Path[1:])
  http.ServeFile(res, req, req.URL.Path[1:])
}

func CreatePerson(params martini.Params, person Person, res http.ResponseWriter, req *http.Request) []byte {
  fmt.Printf("The request was made with %s and %s \n", person.Name, person.GroupId)
  name := person.Name
  group := person.GroupId
  var people []Person

  qh.SQLString.Reset()
  qh.SQLString.WriteString(fmt.Sprintf("INSERT INTO persons(name,group_id) VALUES('%s','%s');", name, group))
  qh.Err = sqlutil.FetchRows(qh.Query(), &people)

  if qh.Err != nil {
    fmt.Printf("Something went wrong trying to create a person")
  }

  return []byte("")
}

func DeletePerson(params martini.Params, res http.ResponseWriter, req *http.Request) []byte {
  qh.SQLString.Reset()
  qh.SQLString.WriteString(fmt.Sprintf("DELETE FROM persons WHERE id=%s;", params["id"]))
  qh.Query()

  return []byte("")
}

func main() {
  m := martini.Classic()

  m.Get("/static", HandleStatic)
  m.Get("/javascript/react/react-with-addons.js", func(res http.ResponseWriter, req * http.Request){
    http.ServeFile(res, req, "javascript/react/react-with-addons.js")  
  })
  m.Get("/javascript/react/JSXTransformer.js", func(res http.ResponseWriter, req * http.Request){
    http.ServeFile(res, req, "javascript/react/JSXTransformer.js")  
  })

  m.Get("/persons", IndexPersons)
  m.Get("/groups", IndexGroups)
  m.Delete("/persons/:id", DeletePerson)

  m.Post("/persons", binding.Bind(Person{}), CreatePerson)

  m.RunOnAddr(":8000")
}
