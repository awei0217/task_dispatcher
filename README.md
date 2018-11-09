#task_dispatcher

        任务调度系统 2018-08-05
        
#author
                
        孙朋伟,王瑞,张宁,李桐葳
        
#realization language
        
        本系统采用google的golang语言编程实现，如果不了解golang,请自行学习        
        
#目录结构说明
        common       公共类库
        config       配置信息
        domain       业务属性和方法
        enum         枚举
        execute      任务执行类
        net_client   集群中各个节点网络通信客户端
        net_server   集群中各个节点网络通信服务端
        scheduler    任务调度模块
        task_manager 任务管理模块
        web          前端页面模块
        
            controller    控制层
            public        js和css库
            views         视图层
            

#本系统功能
        1、支持任务的增删改查（方便对任务的统一管理）
        2、支持任务的分组,分区,并发调用
        3、采用http协议通过url调用任务
        4、支持任务类型的划分（定时任务，手动任务）等
        5、支持机器down掉后的任务自动转移
        6、支持新增机器后的任务自动分配
        7、支持任务执行时日志的持久化和查看（可动态配置）
        8、支持任务和系统负载的监控
        
# 系统启动步骤

        1、首先根据项目中的sql创建库和表结构 
        2、修改config包下的mysql.go 文件，把其中的数据库host和username，password配置成自己的
           如果是在外网：请修改common包下的redis.go,把其中的addr，username，password配置成自己的
        4、修改common包下emial.go 邮件中配置项
            下载过程中可能遇到上述依赖包中依赖的golang.org的一些包无法下载,百度一下，有很多解决方法，辛苦自行学习
        5、第一次启动时请在 /config/task_slice.yml 中配置机器IP的分片
        6、在web报下运行main.go中的main方法，启动成功后,浏览器中访问localhost:8080
        
# 启动过程中如有问题，辛苦联系
        erp : sunpengwei               
        mail: sunpengwei1992@aliyun.com
                              
                                    
    



