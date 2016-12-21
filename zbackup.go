package main

import(
  "strings"
  "fmt"
  "os/exec"
  "bufio"
  "flag"
  "strconv"
  "time"
  "github.com/jasonlvhit/gocron"
  "io"
  "log"
  "os"
  "os/signal"
  "syscall"
)

type cmd struct{
  mode            string          //create(default) , list , delete , daemonHandle , daemon
  target_dataset  string
  id              int             //0(default)
  rotation        int             //20(default)
  time            string
  config          string          // /usr/local/etc/zbackup.conf (default)
}

type schedule struct{
  dataset   string
  rotation  int
  mode      string        // m: minutes, h: hours, d: days, w: weeks
  interval  int
  enable    bool
}

func main(){
  //set flag
  list:=flag.Bool("list",false,"flag of list")
  del:=flag.Bool("delete",false,"flag of del")
  help:=flag.Bool("help",false,"flag of help")
  daemon:=flag.Bool("daemon",false,"flag of daemonHandle")
  d:=flag.Bool("d",false,"flag of daemon(d)")
  config:=flag.String("config","","flag of config")
  c:=flag.String("c","","flag of config(c)")
  t:=flag.Bool("t",false,"flag of daemon")

  flag.Parse()

  //help page
  if(*help==true){
    fmt.Println("This is a tool to manage zfs snapshots")
    fmt.Println("Command: zbackup ( target_dataset [rotation_count] | --list [target_dataset [ID]] | --delete [target_dataset [ID]] | --daemon [--config /path/to/your/conf])")
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
    fmt.Println("Daemon\tzbackup --daemon [--config /path/to/your/conf]")
    fmt.Println("\tBackup automatically background.")
    fmt.Println("\tYou can use custom conf path with \"--config path/to/your/conf\",or zbackup will use /usr/local/etc/zbackup.conf in default.")
    fmt.Println("")
    fmt.Println("Help\tzbackup --help")
    return
  }

  cmdin:=cmd{"create","",0,20,timeModify(time.Now()),"/usr/local/etc/zbackup.conf"}
  args:=flag.Args()
  var args1 int = 0

  //set cmd
  if *t {
    cmdin.mode="daemon"
    if *c!=""{
      cmdin.config=*c
    }
  }else if *d||*daemon {
    cmdin.mode="daemonHandle"
    if *c!=""{
      cmdin.config=*c
    }else if *config!=""{
      cmdin.config=*config
    }
  }else if len(args)>0{
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
  cmdin.process()
}

func timeModify(t time.Time)string{
  TimeStr:=fmt.Sprintln(t)
  StrSpl:=strings.Split(TimeStr," ")
  StrSpltime:=strings.Split(StrSpl[1],".")
  return StrSpl[0]+"_"+StrSpltime[0]+"_"+StrSpltime[1]
}

func (cmdin cmd) process(){
  switch cmdin.mode {
  case "create":
    cmdin.create()
  case "list":
    cmdin.list()
  case "delete":
    cmdin.del()
  case "daemonHandle":
    cmdin.daemonHandle()
  case "daemon":
    cmdin.daemon()
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

//start dameon
func (cmdin cmd) daemonHandle(){
  if !exist(cmdin.config) {
    log.Printf("no such config file\n")
    return
  }
  cmd := exec.Command("zbackup","-t","-c",cmdin.config)
  cmd.Stdout = os.Stdout
  err := cmd.Start()
  if err != nil {
    log.Fatal(err)
  }
  log.Printf("pid: %d\n", cmd.Process.Pid)
}

//daemon
//dataset rotation mode interval enable
func (cmdin cmd) daemon(){
  //process
  for true{
    //check file valid
    if !exist(cmdin.config) {
      return
    }
    //schedule list
    f,_:=os.Open(cmdin.config)
    defer f.Close()
    fbuf:=bufio.NewReader(f)
    schS:=make([]schedule,0)
    var sch schedule
    var exe bool=true
    var newone bool
    for exe {
      sch.enable=true
      newone=true
      for true{
        p,err:=fbuf.Peek(1)
        if err==io.EOF{
          exe=false
          break
        }
        if p[0]=='[' {
          if newone{
            newone=false
            l,_,_:=fbuf.ReadLine()
            sch.dataset=strings.TrimRight(strings.TrimLeft(strings.Split(string(l),"]")[0], "[ "), " ")
          }else{
            break
          }
        }else if p[0]=='e'{
          l,_,_:=fbuf.ReadLine()
          en:=strings.Split(strings.TrimRight(strings.Split(string(l),"#")[0], " \r\n"), "=")[1]
          if en=="no"{
            sch.enable=false
          }
        }else if p[0]=='p'{
          l,_,_:=fbuf.ReadLine()
          stmp:=strings.Split(strings.TrimRight(strings.Split(string(l),"#")[0], " \r\n"), "=")[1]
          sch.mode=string(stmp[len(stmp)-1])
          sch.rotation,_=strconv.Atoi(strings.Split(stmp[:len(stmp)-1],"x")[0])
          sch.interval,_=strconv.Atoi(strings.Split(stmp[:len(stmp)-1],"x")[1])
        }else{
          fbuf.ReadLine()
        }
      }
      if sch.enable {
        schS=append(schS,sch)
      }
    }
    //add schedule
    cron := gocron.NewScheduler()
    for a:=0;a<len(schS);a++ {
      switch schS[a].mode{
      case "m":
        cron.Every(uint64(schS[a].interval)).Minutes().Do(schS[a].task)
      case "h":
        cron.Every(uint64(schS[a].interval)).Hours().Do(schS[a].task)
      case "d":
        cron.Every(uint64(schS[a].interval)).Days().Do(schS[a].task)
      case "w":
        cron.Every(uint64(schS[a].interval)).Weeks().Do(schS[a].task)
      }
    }
    crch:=cron.Start()
    //reload
    c:=make(chan os.Signal,1)
    signal.Notify(c,syscall.SIGHUP)
    go func(){
      for range c {
        cron.Clear()
        close(crch)
        return
      }
    }()
    <-crch
  }
}

func (sch schedule) task(){
  cmdin:=cmd{"create",sch.dataset,0,sch.rotation,timeModify(time.Now()),""}
  cmdin.create()
}

func exist(path string)bool{
  _,err:=os.Stat(path)
  if err == nil {
    return true
  }else{
    return false
  }
}
