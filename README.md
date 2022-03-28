# gomod
`gomod` 是一款分析 Goalng 项目三方依赖仓库 Star 、Fork 的工具。

支持用户按统计数据的升降序查看，默认按 Fork 降序排序，再次按 `s` 则以相反的顺序排序; 按 `f` 键可以切换到按 Fork 排序，再次按 `f` 则以相反顺序排序。

**为什么开发这个工具？**

在日常工作中，经常会关注一些优秀的项目如 [kubernetes](https://github.com/kubernetes/kubernetes)、[istio](https://github.com/istio/istio)、[etcd](https://github.com/etcd-io/etcd)、[nsq](https://github.com/nsqio/nsq)等。而这些项目的成功离不开它采用的那些三方库，于是通过这些项目发现一些不错的三方库慢慢变成了一种学习途径。对于这些三方库来讲它们的 Star 、Fork、Watch 数量成为了一个行业评价其是否优秀的标准。为了方便查看这些数据信息，减少人肉操作因此开发了此工具。

# 实现原理
由于使用官方API (https://api.github.com) 接口存在限率问题，所以采用了最直接的采集页面内容来分析的方法。

通过远程抓取 `https://github.com/user/repo` 仓库的页面的信息，来获取仓库的 `Star`、`Fork` 统计数据。

1. 读取 `go.mod` 文件，分析出 `Require` 指令中的所有依赖仓库
2. 远程抓取仓库页面，统计出 Star、Fork 信息
3. 输出统计结果，支持按 Star 和 Fork 的升序和降序排序



# 安装

```
$ go install github.com/cfanbo/gomod@latest
```

 此命令将自动编译并将生成的可执行文件存放在 `$GOPATH/bin` 目录里，推荐此方法。

您也可以下载到本机安装

```
$ git clone https://github.com/cfanbo/gomod.git
$ cd gomod
$ make install
```



# 效果展示

```
$ cd client-go
$ gomod
Press `ESC` to quit.  Press the [s|f|m|g] key to sort
Modules: 227 total, 227 success, 0 failed.  Display range 1 ~ 41

STAR   FORK   SHARE  MODULE                                  GITHUB
---------------------------------------------------------------------------------------------------------------
62505  18014         github.com/docker/docker                https://github.com/docker/docker
41680  6954          github.com/prometheus/prometheus        https://github.com/prometheus/prometheus
25800  2256          github.com/spf13/cobra                  https://github.com/spf13/cobra
21394  6117          helm.sh/helm/v3                         https://github.com/helm/helm
20173  2116          github.com/sirupsen/logrus              https://github.com/sirupsen/logrus
18686  1627          github.com/spf13/viper                  https://github.com/spf13/viper
16691  2809          github.com/gorilla/websocket            https://github.com/gorilla/websocket
16241  1498          github.com/gorilla/mux                  https://github.com/gorilla/mux
15812  1221          github.com/stretchr/testify             https://github.com/stretchr/testify
15586  3492          google.golang.org/grpc                  https://github.com/grpc/grpc-go
15278  1126          go.uber.org/zap                         https://github.com/uber-go/zap
11813  2626          github.com/antlr/antlr4/runtime/Go/antlrhttps://github.com/antlr/antlr4
11256  1268          github.com/golang/groupcache            https://github.com/golang/groupcache
10623  863           github.com/json-iterator/go             https://github.com/json-iterator/go
8322   1517          github.com/golang/protobuf              https://github.com/golang/protobuf
8200   1714   Y      sigs.k8s.io/kustomize/api               https://github.com/kubernetes-sigs/kustomize
8200   1714   Y      sigs.k8s.io/kustomize/kyaml             https://github.com/kubernetes-sigs/kustomize
7646   622           github.com/pkg/errors                   https://github.com/pkg/errors
6703   740           github.com/fsnotify/fsnotify            https://github.com/fsnotify/fsnotify
6629   2077          github.com/docker/distribution          https://github.com/docker/distribution
6429   858           github.com/lucas-clemente/quic-go       https://github.com/lucas-clemente/quic-go
```

输出字段：

`Star` 仓库的 Star 数量；如果统计数据超时则显示为 `?`

`Fork` 仓库的 Fork 数量；如果统计数据超时则显示为 `?`

`Share` 是否存在相同的仓库；如果引用了同一个仓库的多个包，则此列将显示 `Y`

`Module` 当前项目依赖的三方库名称

`GitHub` 三方库在 `github.com`  上的托管地址



可以通过多次按 `s`、 `f`、`m` 或`g` 键查看效果，如果要退出查看模式按 `q` 或 `Esc`键即可。  
如果要访问仓库在 github.com 上的托管地址，直接 `Enter` 键即可; 如果要访问在 pkg.go.dev 网站的地址，直接按`空格键` 即可。

注：`github.com` 网站显示的数据是对其进行四舍五入后的结果。



# 使用方法

根据用户习惯支持多种用法。

### 1. 本地项目
一、适用于查看用户本地项目的场景，在 Golang 项目里的任意目录即可执行此命令

```
$ gomod
```

此命令将查看当前项目文件 `go.mod` 中的所有依赖信息。



二、指定 `go.mod` 文件路径

```
$ gomod workspace/client-go/go.mod
```

此方法比较灵活, 可以查看任意位置的 `go.mod` 文件依赖信息。



### 2. 远程项目

如果想查看托管在 `github.com` 上的任何项目信息，可通过以下几种方法：

方法一

```
$ gomod github.com/kubernetes/client-go
```



方法二

```
$ gomod github.com/kubernetes/client-go/blob/master/CHANGELOG.md
```

指定 `github.com/user/repo` 仓库任意路径或文件地址，即使这个网站 `github.com/user/repo` 后面的路径不存在。如：

```
$ gomod github.com/kubernetes/client-go/noexists
```



方法三：

指定一个远程的 `go.mod` 网址

```
$ gomod https://raw.githubusercontent.com/kubernetes/client-go/master/go.mod
```

或 第三方网站

```
$ gomod https://example.com/go.mod
```




# 常见问题
1. 展示结果里统计数字会出现 ? 符号
由于每次执行统计时都需要实时访问 https://pkg.go.dev 和 https://github.com 两个网站，因此请确保这网站可以正常访问;
如果在本地网络正常的情况下出现此问题，一般是由于网络不稳定的原因，可多次执行此命令或等待其自动请求完成。
