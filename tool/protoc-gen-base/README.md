# pb_generator proto代码生成器

> [原文件参考](https://raw.githubusercontent.com/golang/protobuf/master/protoc-gen-go/generator/generator.go)

## 修改
1. 每次生成只支持一个插件(plugin),原因是希望每个插件生成独立的pb.xx.go文件
2. 不负责生成protoc-gen-go生成的内容

## 插件Example
* [myplugin](protoc-gen-example/myplugin/myplugin.go)

## 命令说明
```shell
protoc \ 
--proto_path=. \
--proto_path=./pb \
--go_out=paths=source_relative:pb \
--example_out=plugins=gen_func+gen_interface,paths=source_relative:pb_example \
pb_example/example.proto
```

* --proto_path=. --proto_path=./pb  
	proto文件的搜索目录
* --go_out=paths=source_relative:./pb  
	调用protoc-gen-go 生成go代码,paths=source_relative表示生成相对目录
* --example_out=plugins=gen_func+gen_interface,paths=source_relative:.  
	调用protoc-gen-example的gen_func和gen_interface插件，生成目录是相对当前目录
* pb_example/example.proto  
	需要生成的proto文件