package hawkservice


import (
	"fmt"
	"log"
	"io/ioutil"
	"strings"
	"gopkg.in/yaml.v2"
	"../utils"
	"net/http"
	"encoding/json"
	"os/exec"
	"os"

)

//config for release.json


type version struct {
	Helm    string `json:"helm"`
	Service string `json:"service"`
} 

type releaseJson struct {
	Deploy  string `json:"deploy"`
	Service string `json:"service"`
	Version version `json:"version"`
}




//config for yml file
type cfg struct {
	
    Service struct{ 
	Arnrole string `yaml:"arnrole"`
    Image struct{
    Repository string `yaml:"repository"`
    } 
    
    }`yaml:"service"`
	HelmVersion 	 string `yaml:"helmVersion"`
	Error string
   }


   //response body config
   type Body struct {
	ResponseCode int
	Message      string
}




//func to parse the service configuration.yaml file 
func (c *cfg) getConfig(path string) *cfg {
	//var c config
	
	yamlFile, err := ioutil.ReadFile(path)
if err != nil {
	log.Printf("yamlFile.Get err   #%v ", err)

}

//fmt.Println(yamlFile)
err = yaml.Unmarshal(yamlFile, &c)
if err != nil {
	log.Fatalf("Unmarshal: %v", err)
}else{
	return c	
}

return c
}


//func to get index value for env

func indexOf(element string, data []string) (int) {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1    //not found.
 }


 //func to list all the folders present in the repo
func listDir( path string) []string{

	//fmt.Println(path)
dirList := make([]string, 100)

files, err := ioutil.ReadDir(path)
    if err != nil {
        log.Fatal(err)
    }

    for _, file := range files {
       // fmt.Println(file.Name())
		dirList=append(dirList,file.Name())
    }



	return dirList

}

func compareEnv(w http.ResponseWriter, r *http.Request){


conf := utils.ReadConfig()

repoPath:=conf.RepoPath
//fmt.Println(repoPath)

env1:=strings.TrimSpace(r.URL.Query().Get("env1"))
env2:=strings.TrimSpace(r.URL.Query().Get("env2"))

flag:=true

fmt.Println(utils.CheckEnv(env1,env2,conf.EnvRepo))
if utils.CheckEnv(env1,env2,conf.EnvRepo)!=true{
	log.Printf("not truee----")
	utils.RespondWithJSON("Please check the environment name", w, r)
	flag=false
}

if flag{

cmd := &exec.Cmd {
	Path: "./repo.sh",
	Args: []string{ "./repo.sh" },
	Stdout: os.Stdout,
	Stderr: os.Stdout,
}


err := cmd.Run()

if err != nil {
        log.Fatal(err)
    }

env1Path:=strings.TrimSpace(repoPath+"/"+env1)
//env1Path:="C:\\HAWK\\Repo\\ccgf-qastaging-hawk-config\\"
env2Path:=strings.TrimSpace(repoPath+"/"+env2)
//env2Path:="C:\\HAWK\\Repo\\ccgf-perf-hawk-config\\"


//get the actual env name from properties file
env1=conf.EnvList[indexOf(env1,conf.EnvRepo)]
env2=conf.EnvList[indexOf(env2,conf.EnvRepo)]

fmt.Printf("%s,%s",env1,env2)
fmt.Println(env1Path)

CC := r.URL.Query().Get("Email") 

var env1Config cfg
var env2Config cfg
var colorService string

jsondata := []releaseJson{}

dirListPerf:=listDir(env1Path)

//fmt.Println(dirListPerf)

htmlData:="<html>  <table style='backgound:#fff;border-collapse: collapse;' border = '1' cellpadding = '6'> <tr style='background:#000;color:#fff'><th>Service Name</th> <th>"+env1+" Service Version</th>  <th>"+env2+" Service Version </th> <th> "+env1+" Helm Version</th>  <th>"+env2+" Helm Version </th></tr>  "

//fmt.Println(dirListPerf)
for _,i := range dirListPerf{

	if i!="" && i!="README.md" && i!=".git" {
	
	env1Config.getConfig(env1Path+"/"+i+"/configuration.yaml")
	env2Config.getConfig(env2Path+"/"+i+"/configuration.yaml")
	
	//skip the service if repo empty
	//if env1Config.Service.Image.Repository=="" || env2Config.Service.Image.Repository=="" {
		if env1Config.Service.Image.Repository==""  {
			continue
		}
	
		if env2Config.Service.Image.Repository==""  {
			env2Config.Service.Image.Repository=":"
		}
	
	
		repoenv1:=strings.Split(env1Config.Service.Image.Repository,":")
		repoenv2:=strings.Split(env2Config.Service.Image.Repository,":")
		
	//	fmt.Println(i,repoenv2)
		if repoenv1[1]!=repoenv2[1]{
			colorService="#ff8080"
		}else{
			colorService="#66cc00"
		}
		
		helmVersionenv1:=""
		helmVersionenv2:=""
		
	
		if env1Config.HelmVersion !=""  {
			helmVersionenv1=env1Config.HelmVersion
			helmVersionenv2=env2Config.HelmVersion
	
			if env2Config.HelmVersion ==""  {
			helmVersionenv2=""}
			//fmt.Println(helmVersionenv1,helmVersionenv2)
		}

	htmlData=htmlData+"<tr> <td style='background:"+colorService+"'> <b>"+i+"</b></td><td>"+repoenv1[1]+"</td><td>"+repoenv2[1]+"</td><td>"+helmVersionenv1+"</td><td>"+helmVersionenv2+"</td></tr>"

	if colorService=="#ff8080"{
	newJsonDoc := &releaseJson{
        Deploy: i,
		Service:i,
		Version:version{
			Helm:helmVersionenv1,
			Service:repoenv1[1],
		},
    }

    jsondata = append(jsondata, *newJsonDoc)
}

}

}
//fmt.Println(htmlData)
outputJSON, _ := json.Marshal(jsondata)
//fmt.Println(outputJSON)
err = ioutil.WriteFile("release.json", outputJSON, 0644)
    if err != nil {
        log.Println(err)
    }
subject:="Service Version Comparison "+env1+" vs "+env2

utils.SendMail(htmlData, subject, CC)

utils.RespondWithJSON("Email Sent Successfully", w, r)

}

}



func test(w http.ResponseWriter, r *http.Request) {

	body := Body{ResponseCode: 200, Message: "OK"}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBody)

}