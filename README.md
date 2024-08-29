# 某节课视频下载器

> [!WARNING]
>
> 注意：必须开通会员才能下载相应的课程视频，当前只能下载视频，并只会保存为[.ts]后缀，需要**F12**打开浏览器的调试模式，获取该网站的  **Authorization**, **Cookie**, **Sjk-Apikey** ，下载视频是，需要进入相应的课程地址，例如 https://study.sanjieke.cn/course/0/34003473/34734298 ，34003473 为课程id

# 使用方法

## 1 编译应用程序

```shell
make sdl
```

## 2 执行应用程序（以Windows为例）

```shell
./sdl.exe # 会进行相应的输入，需要输入相应的值
```

- **Authorization**：**需要F12** 在相应的课程下面获取相应的**Header值**
- **Cookie**：**需要F12** 在相应的课程下面获取相应的**Header值**
- **Sjk-Apikey**: **需要F12** 在相应的课程下面获取相应的**Header值**
- **课程Id**: 例如 https://study.sanjieke.cn/course/0/34003473/34734298 ，34003473 为课程id
- **输出目录**: 需要下载的目录地址

## 3.以配置文件形式执行（以Windows为例）

```
./sdl.exe -c config.yaml # -c 可选，如果使用，则会用配置文件读取配置
```

