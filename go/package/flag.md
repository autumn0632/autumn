# flag包&log包

flag包实现了命令行参数的解析

log包实现了简单的日志服务。

> golang标准库的日志框架非常简单，仅仅提供了print，panic和fatal三个函数对于更精细的日志级别、日志文件分割以及日志分发等方面并没有提供支持。

**示例1**

```go
import flag
/*
1. 注册flag
2. 解析flag
3. 执行函数
*/
/*主要函数：
1. usage ：命令行help
2. vaildateArgs：校验参数个数
3. runParses：解析参数
*/

type CLI struct{}
 

var (
	h    bool
	v, V bool
	t, T bool
	q    *bool

	s, p, c, g string
)

func (cli *CLI) usage() {
    	fmt.Fprintf(os.Stderr, `nginx version: nginx/1.10.0
Usage: nginx [-hvVtTq] [-s signal] [-c filename] [-p prefix] [-g directives]

Options:
`)

	//打印所有已定义参数的默认值（调用 VisitAll 实现）
	flag.PrintDefaults()
}

func (cli *CLI) vaildateArgs() {
    if len(os.Args < 2) {
        cli.usage()
        os.Exit(1)
    }
}
func (cli *CLI) runParses1() {
    //fir：注册flag
    //sec：解析flag
    //thi: 解析成功，进行对应执行函数，解析失败，打印usage
    
 	//注册flag
	flag.BoolVar(&h, "h", false, "this help")

	flag.BoolVar(&v, "v", false, "show version and exit")
	flag.BoolVar(&V, "V", false, "show version and configure options then exit")

	flag.BoolVar(&t, "t", false, "test configuration and exit")
	flag.BoolVar(&T, "T", false, "test configuration, dump it and exit")

	flag.StringVar(&s, "s", "", "send `signal` to a master process: stop, quit, reopen, reload")
	flag.StringVar(&p, "p", "/usr/local/nginx/", "set `prefix` path")
	flag.StringVar(&c, "c", "conf/nginx.conf", "set configuration `file`")
	flag.StringVar(&g, "g", "conf/nginx.conf", "set global `directives` out of configuration file")

	cli.vaildateArgs()
	//change the default usage
	flag.Usage = cli.usage
	
    //解析flag
	flag.Parse()

	if v {
		fmt.Println("version:5.6")
	}
    else if h {
        cli.usage()
    }
    
}

```



**示例2**

```go
import flag

type CLI struct{}

func (cli *CLI1) usage() {
	fmt.Println("Usage:")
	fmt.Println("  createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  createwallet - Generates a new key-pair and saves it into the wallet file")
	fmt.Println("  getbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("  listaddresses - Lists all addresses from the wallet file")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
	fmt.Println("  reindexutxo - Rebuilds the UTXO set")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT -mine - Send AMOUNT of coins from FROM address to TO. Mine on the same node, when -mine is set.")
	fmt.Println("  startnode -miner ADDRESS - Start a node with ID specified in NODE_ID env. var. -miner enables mining")
}

func (cli *CLI) vaildateArgs() {
        if len(os.Args < 2) {
        cli.usage()
        os.Exit(1)
    }
}

func (cli *CLI1) runParses() {
	cli.vaildateArgs()

	//NewFlagSet创建一个新的、名为name，采用errorHandling为错误处理策略的FlagSet。
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	//String用指定的名称、默认值、使用信息注册一个string类型flag。返回一个保存了该flag的值的指针。
	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")

	switch os.Args[1] {
	case "getbalance":
		//从arguments中解析注册的flag。
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.usage()
		os.Exit(1)
	}

	//返回是否f.Parse已经被调用过。
	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		//cli.getBalance(*getBalanceAddress)
		fmt.Println("run getbalance.")
	}
}
```

