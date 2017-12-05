package library

import (
	. "connection"
	"fmt"
	"os"
	"os/exec"
	"log"
	"io/ioutil"
)

func ProcessLiquibase() {
	//先导出
	exportLiquibaseTable()
	//处理
	process()
	//重新初始化数据库链接
	Init()
}

func process() {
	
	//先遍历租户删除liquibase的记录表
	for _, value := range getTenants() {
		//切库执行
		DBConfig.Database = value
		//如果这个库不存在或者链接出错，直接跳过
		err := ChangeDB()
		if err != nil {
			Logger.Println("error by connect database:", DBConfig.Database)
			continue
		}
		//删除changelog表
		deleteLiquibaseTable()
		fmt.Println("success delete liquibase table in database:", DBConfig.Database)
		createLiquibaseTable()
		fmt.Println("success create liquibase table in database:", DBConfig.Database)
	}
	
}

func getTenants() []string {
	rows, err := DB.Query("select id from tenements")
	CheckErr(err)
	
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	
	var result []string
	result = append(result, DBConfig.DefaultTenant)
	
	for rows.Next() {
		//将行数据保存到record字典
		err = rows.Scan(scanArgs...)
		record := make(map[string]string)
		
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
		result = append(result, DBConfig.TenantPrefix+record["id"])
	}
	
	return result
}

func deleteLiquibaseTable() {
	
	_, err := DB.Query("DROP TABLE IF EXISTS `DATABASECHANGELOG`,`DATABASECHANGELOGLOCK`,`databasechangelog`,`databasechangeloglock`")
	
	if err != nil {
		fmt.Printf("delete changelog table in "+DBConfig.Database+" error: %v\n", err)
		os.Exit(-1)
	}
}

func createLiquibaseTable() {
	
	//直接用mysql导入
	var options []string

	options = append(options, fmt.Sprintf(`-h%v`, DBConfig.Host))
	options = append(options, fmt.Sprintf(`-u%v`, DBConfig.Username))
	options = append(options, fmt.Sprintf(`-P%v`, DBConfig.Port))
	options = append(options, fmt.Sprintf(`-p%v`, DBConfig.Password))
	options = append(options, fmt.Sprintf(`-D%v`, DBConfig.Database))
	options = append(options, "-e source /www/go/liquibase.sql;")

	cmd := exec.Command("mysql", options...)
	
	stdout, err := cmd.StdoutPipe()
	
	if err != nil {
		Logger.Println(err)
	}
	
	if err := cmd.Start(); err != nil {
		Logger.Println(err)
	}
	Logger.Println(stdout)
	/*cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		Logger.Println("cmd.Output: ", err)
	}*/
	
}

func exportLiquibaseTable() {
	
	DBConfig.Database = DBConfig.BasicTenant
	//如果这个库不存在或者链接出错，直接跳过
	err := ChangeDB()
	if err != nil {
		Logger.Println("error by connect database:", DBConfig.Database)
	}
	//先改成小写再导出
	rows, err := DB.Query("show tables like 'databasechangelog'")
	defer rows.Close()
	var name string
	for rows.Next() {
		if err := rows.Scan(&name); err != nil {
			Logger.Println(err)
		}
	}
	/*if name != "DATABASECHANGELOG" {
		DB.Query("rename table `databasechangeloglock` TO`DATABASECHANGELOGLOCK`")
		DB.Query("rename table `databasechangelog` TO`DATABASECHANGELOG`")
	}*/
	
	if name != "databasechangelog" {
		DB.Query("rename table `DATABASECHANGELOGLOCK` TO`databasechangeloglock`")
		DB.Query("rename table `DATABASECHANGELOG` TO`databasechangelog`")
	}
	
	//直接用mysqldump导出
	var options []string
	
	options = append(options, "--skip-comments")
	options = append(options, fmt.Sprintf(`-h%v`, DBConfig.Host))
	options = append(options, fmt.Sprintf(`-u%v`, DBConfig.Username))
	options = append(options, fmt.Sprintf(`-P%v`, DBConfig.Port))
	options = append(options, fmt.Sprintf(`-p%v`, DBConfig.Password))
	options = append(options, DBConfig.Database)
	options = append(options, "databasechangelog")
	options = append(options, "databasechangeloglock")
	
	cmd := exec.Command("mysqldump", options...)

	stdout, err := cmd.StdoutPipe()
	
	if err != nil {
		log.Fatal(err)
	}
	
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	
	bytes, err := ioutil.ReadAll(stdout)
	
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("./liquibase.sql", bytes, 0644)
	if err != nil {
		panic(err)
	}
	
	Init()
}
