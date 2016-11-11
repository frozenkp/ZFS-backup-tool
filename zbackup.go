package main

import(
  "strings"
  "fmt"
  "os/exec"
  "flag"
  "strconv"
  "time"
)

type cmd struct{
  mode            string          //create(default) , list , delete
  target_dataset  string
  id              int             //0(default)
  rotation        int             //20(default)
  time            string
}

func main(){
  //set flag
  list:=flag.Bool("list",false,"flag of list")
  del:=flag.Bool("delete",false,"flag of del")
  help:=flag.Bool("help",false,"flag of help")

  flag.Parse()

  //help page
  if(*help==true){
    fmt.Println("This is a tool to manage zfs snapshots")
    fmt.Println("Command: zbackup ( target_dataset [rotation_count] | --list [target_dataset [ID]] | --delete [target_dataset [ID]] )")
    fmt.Println("")
    fmt.Println("Create\tzbackup target_dataset [rotation_count]")
    fmt.Println("\tcreate a snapshot for target_dataset with rotation_count")
    fmt.Println("\trotation count: The sum of snapshots of this snapshot,it will remove old snapshot automatically (default=20)")
    fmt.Println("")
    fmt.Println("List\tzbackup --list [target_dataset [ID]]")
    fmt.Println("\tlist all snapshots created by this tool")
    fmt.Println("\tyou can use target_dataset and id to modify your list")
    fmt.Println("")
    fmt.Println("Delete\tzbackup --delete [target_dataset [ID]]")
    fmt.Println("\tdelete snapshots")
    fmt.Println("")
    fmt.Println("Help\tzbackup --help")
    return
  }

  cmdin:=cmd{"create","",0,20,timeModify(time.Now())}
  args:=flag.Args()
  var args1 int = 0
  fmt.Println(cmdin)
  //set cmd
  if len(args)>0{
    cmdin.target_dataset=args[0]
    if len(args)>1 {
      args1,_=strconv.Atoi(args[1])
    }
    if *list || *del {
      if *list {
        cmdin.mode="list"
      }else{
        cmdin.mode="delete"
      }
      cmdin.id=args1
    }else{
      cmdin.mode="create"
      if args1!=0 {
        cmdin.rotation=args1
      }
    }
  }else{
    if *list {
      cmdin.mode="list"
    }else{
      cmdin.mode="delete"
    }
  }

  //process
  fmt.Println(cmdin)
  cmdin.process()
}

func timeModify(t time.Time)string{
  TimeStr:=fmt.Sprintln(t)
  StrSpl:=strings.Split(TimeStr," ")
  StrSpltime:=strings.Split(StrSpl[1],".")
  return StrSpl[0]+"_"+StrSpltime[0]
}

func (cmdin cmd) process(){
  switch cmdin.mode {
    case "create":
      cmdin.create()
    case "list":
      cmdin.list()
    case "delete":
      cmdin.del()
  }
}

func (cmdin cmd) create(){
  //create
  bkname:=cmdin.target_dataset+"@zbk_"+cmdin.time
  c:=exec.Command("zfs","snapshot",bkname)
  c.Run()
  ListOut,_:=exec.Command("zfs","list","-r","-t","snapshot","-o","name",cmdin.target_dataset).Output()
  strList:=strings.Split(string(ListOut),"\n")
  //ratation count delete
  //0 -> "name" , len(strList)-1 -> ""
  //target_dataset && zbk_
  chosenList:=make([]string,0)
  for a:=1;a<len(strList)-1;a++{
    s:=strings.Split(strList[a],"@")
    if s[0]==cmdin.target_dataset && strings.Contains(s[1],"zbk_") {
      chosenList=append(chosenList,strList[a])
    }
  }
  //delete
  for a:=0;a<len(chosenList)-cmdin.rotation;a++{
    c:=exec.Command("zfs","destroy",chosenList[a])
    c.Run()
  }
}

func (cmdin cmd) list(){
  //--list dataset [ID]
  if cmdin.target_dataset != "" {
    ListOut,_:=exec.Command("zfs","list","-r","-t","snapshot","-o","name",cmdin.target_dataset).Output()
    strList:=strings.Split(string(ListOut),"\n")
    //0 -> "name" , len(strList)-1 ->""
    //target_dataset && zbk_
    chosenList:=make([]string,0)
    for a:=1;a<len(strList)-1;a++{
      s:=strings.Split(strList[a],"@")
      if s[0]==cmdin.target_dataset && strings.Contains(s[1],"zbk_"){
        chosenList=append(chosenList,s[1])
      }
    }
    //output list
    //ID Dataset Time
    fmt.Printf("ID\t\tDataset\t\t\tTime\n")
    if cmdin.id!=0 {
      ss:=strings.Split(chosenList[cmdin.id-1],"_")
      fmt.Printf("%d\t\t%s\t\t%s %s\n",cmdin.id,cmdin.target_dataset,ss[1],ss[2])
    }else{
      for a:=0;a<len(chosenList);a++{
        ss:=strings.Split(chosenList[a],"_")
        fmt.Printf("%d\t\t%s\t\t%s %s\n",a+1,cmdin.target_dataset,ss[1],ss[2])
      }
    }
  }else{
    //--list
    ListOut,_:=exec.Command("zfs","list","-r","-t","snapshot","-o","name").Output()
    strList:=strings.Split(string(ListOut),"\n")
    //0 -> "name" , len(strList)-1 ->""
    //zbk_
    chosenList:=make([]string,0)
    for a:=1;a<len(strList)-1;a++{
      s:=strings.Split(strList[a],"@")
      if strings.Contains(s[1],"zbk_"){
        chosenList=append(chosenList,strList[a])
      }
    }
    //output list
    //ID Dataset Time
    fmt.Printf("ID\t\tDataset\t\t\tTime\n")
    for a:=0;a<len(chosenList);a++{
      ss:=strings.Split(chosenList[a],"@")
      sss:=strings.Split(ss[1],"_")
      fmt.Printf("%d\t\t%s\t\t%s %s\n",a+1,ss[0],sss[1],sss[2])
    }
  }
}

func (cmdin cmd) del(){
  //--delete target_dataset [ID]
  if cmdin.target_dataset!=""{
    ListOut,_:=exec.Command("zfs","list","-r","-t","snapshot","-o","name",cmdin.target_dataset).Output()
    strList:=strings.Split(string(ListOut),"\n")
    //0 -> "name" , len(strList)-1 -> ""
    //target_dataset && zbk_
    chosenList:=make([]string,0)
    for a:=1;a<len(strList)-1;a++{
      s:=strings.Split(strList[a],"@")
      if s[0]==cmdin.target_dataset && strings.Contains(s[1],"zbk_") {
        chosenList=append(chosenList,strList[a])
      }
    }
    //delete
    if cmdin.id!=0 {
      c:=exec.Command("zfs","destroy",chosenList[cmdin.id-1])
      c.Run()
    }else{
      for a:=0;a<len(chosenList);a++{
        c:=exec.Command("zfs","destroy",chosenList[a])
        c.Run()
      }
    }
  }else{
    //--delete
    ListOut,_:=exec.Command("zfs","list","-r","-t","snapshot","-o","name").Output()
    strList:=strings.Split(string(ListOut),"\n")
    //0 -> "name" , len(strList)-1 -> ""
    //zbk_
    chosenList:=make([]string,0)
    for a:=1;a<len(strList)-1;a++{
      s:=strings.Split(strList[a],"@")
      if strings.Contains(s[1],"zbk_") {
        chosenList=append(chosenList,strList[a])
      }
    }
    for a:=0;a<len(chosenList);a++{
      c:=exec.Command("zfs","destroy",chosenList[a])
      c.Run()
    }
  }
}

