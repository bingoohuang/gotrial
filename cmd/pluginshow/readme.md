# [Go 插件系统](https://blog.yumaojun.net/2018/04/18/go-plugin/)

Go 插件系统
 
发表于 2018-04-18 |  更新于: 2018-04-24 |  分类于 开发语言 ， Golang |  阅读次数: 417

通过使用插件在运行时扩展程序的功能, 而无需重新编译程序, 这是一个很常见的功能需求, 特别是在模块化设计的程序里面, 比如Nginx的模块系统。 在C/C++中通过使用动态库的方式可以实现动态加载, 但是Go直到1.8官方才开始支持, 下面将介绍Go如何基于动态链接库来实现动态加载。

动态加载的优劣

优点:

> 动态加载, 也称热加载, 每次升级时不用重新编译整个工程，重新部署服务, 而是添加插件时进行动态更新。这对于很多比较重型的服务来说非常重要。

缺点:

> 带来一定的安全风险, 如果一些非法模块被注入如何防范
> 给系统带来一定的不稳定的因素, 如果模块有问题, 没有经过良好的测试, 容易导致服务崩溃
> 为版本管理带来了难题, 特别是在微服务的今天, 同一个服务, 加载了不同的插件, 应该怎么管理版本, 插件版本应该如何管理
> 因此请慎重考虑, 是使用动态插件还是在源码里面进行插件化。

Go的插件系统:Plugin

从1.8版开始, 官方提供了这种插件化的手段: plugin. 此功能使程序员可以使用动态链接库构建松散耦合的模块化程序，可以在运行时动态加载和绑定。

Go插件是使用-buildmode = plugin标记编译的一个包, 用于生成一个共享对象(.so)库文件。 Go包中的导出的函数和变量被公开为ELF符号，可以使用plugin包在运行时查找并绑定ELF符号。Go编译器能够使用build flag -buildmode = c-shared创建C风格的动态共享库。

1.8版本插件功能只能在Linux上使用。 1.10也可以在Mac上运行。

下面将介绍使用Go插件系统创建模块化软件的一些开发原则, 并提供一个功能齐全的示例。

1. 插件开发原则

    使用Go插件创建模块化程序需要遵循与常规Go软件包一样严格的软件实践。然而，插件引入了新的设计问题，因为它们的解耦性质被放大了。因此我们在设计可插拔系统时, 有一些原则需要关注:

1. 插件独立
    
    应该将插件视为与其他组件分离的独立组件。这允许插件独立于他们的消费者，并拥有自己的开发和部署生命周期。注意插件的可用性很重要, 因为它有肯能为整个系统带来不稳定的因素, 因此系统必须为插件集成提供一个简单的封装层, 插件开发人员将系统视为黑盒，不作为所提供的合约以外的假设, 从而保证插件自身的可用性。

1. 使用接口类型作为边界
    
    Go插件可以导出任何类型的包函数和变量。您可以设计插件来将其功能解耦为一组松散的函数。缺点是您必须单独查找和绑定每个函数符号。
然而，更为简单的方法是使用接口类型。创建导出功能的接口提供了统一简洁的交互，并具有清晰的功能划分。解析到接口的符号将提供对该功能的整个方法集的访问，而不仅仅是一个方法。

1. Unix模块化原则
   
   插件代码应该设计成只关注一个功能点。

1. 版本控制

   插件是不透明而独立的实体，应该进行版本控制，以向用户提示其支持的功能。这里的一个建议是在命名共享对象文件时使用语义版本控制。例如，上面的文件编译插件可以命名为eng.so.1.0.0。

插件开发示例:

我以我遇到的一个实际需求为例, 在开发物联网接入组件的时候, 需要动态支持物解析, 下面就开发一个物解析的插件系统。

下面是项目结构, parser.go是接口规约, main.go是主程序, plugins存放多个插件包

```
.
├── main.go
├── parser.go
└── plugins
    ├── car
    │   └── car.go
    └── phone
        └── phone.go
```

编写插件
编写主程序接口规约: main.go

```go
package main
// Parser use to parse things
type Parser interface {
	Parse([]byte) (meta map[string]string, data map[string]float64, err error)
}
```

根据接口规约编写插件: car.go

```go
package main
type car string
func (c *car) Parse([]byte) (meta map[string]string, data map[string]float64, err error) {
	meta = map[string]string{"key1": "carcar"}
	data = map[string]float64{"key1": 1}
	return meta, data, nil
}
var Car car
```

根据接口规约编写插件: phone.go

```go
package main
type phone string
func (p *phone) Parse([]byte) (meta map[string]string, data map[string]float64, err error) {
	meta = map[string]string{"key1": "phonephone"}
	data = map[string]float64{"key1": 2}
	return meta, data, nil
}
var Phone phone
```

编译插件
插件写完后将在plugins目录下编译插件:

```bash
$ cd plugins
$ go build -buildmode=plugin -o car.so car/car.go
$ go build -buildmode=plugin -o phone.so phone/phone.go
```

最终在plugins目录下会生成好我们编译好的插件:

```bash
$ ls *.so
car.so   phone.so
```

使用插件
插件的使用很简单, 大概步骤如下:

用plugin.Open()打开插件文件
用plguin.Lookup(“Export-Variable-Name”)查找导出的符号”Car”或者”Phone”。 请注意，符号名称与插件模块中定义的变量名称相匹配
使用该变量
主程序使用插件: main.go

```go
package main
import (
	"fmt"
	"plugin"
)
// Parser use to parse things
type Parser interface {
	Parse([]byte) (meta map[string]string, data map[string]float64, err error)
}
func pa() {
	plug, err := plugin.Open("./plugins/car.so")
	if err != nil {
		panic(err)
	}
	car, err := plug.Lookup("Car")
	if err != nil {
		panic(err)
	}
	p, ok := car.(Parser)
	if ok {
		meta, data, err := p.Parse([]byte("a"))
		if err != nil {
			panic(err)
		}
		fmt.Printf("meta: %v, data: %v \n", meta, data)
	}
}
func pb() {
	plug, err := plugin.Open("./plugins/phone.so")
	if err != nil {
		panic(err)
	}
	phone, err := plug.Lookup("Phone")
	if err != nil {
		panic(err)
	}
	p, ok := phone.(Parser)
	if ok {
		meta, data, _ := p.Parse([]byte("a"))
		fmt.Printf("meta: %v, data: %v \n", meta, data)
	}
}
func main() {
	pa()
	pb()
}
```

测试是否正常运行:

```bash
$ go run main.go
meta: map[key1:carcard], data: map[key1:1]
meta: map[key1:phonephone], data: map[key1:2]
```
