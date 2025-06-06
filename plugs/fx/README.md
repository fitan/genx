# @fx Plugin - Uber FX 依赖注入自动化

基于 `@do` 插件的设计思路，为 Uber FX 框架提供自动化的依赖注入配置。

## 🎯 功能特性

- ✅ 自动扫描 `@fx` 注解的函数
- ✅ 生成 fx.Option 配置
- ✅ 支持 Provide、Invoke、Decorate
- ✅ 支持接口绑定和分组
- ✅ 智能包名冲突处理
- ✅ 确定性排序保证一致性

## 📝 使用方式

### 1. 服务提供者 (Provide)

```go
// services/user_service.go

// @fx provide
func NewUserService(db *gorm.DB, logger *slog.Logger) *UserService {
    return &UserService{db: db, logger: logger}
}

// @fx provide UserServiceInterface  // 绑定到接口
func NewUserServiceImpl(db *gorm.DB) UserServiceInterface {
    return &UserService{db: db}
}

// @fx provide "" handlers  // 加入到 handlers 组
func NewUserHandler(svc *UserService) *UserHandler {
    return &UserHandler{svc: svc}
}
```

### 2. 启动时调用 (Invoke)

```go
// app/startup.go

// @fx invoke
func StartHTTPServer(handler *gin.Engine, lifecycle fx.Lifecycle) {
    srv := &http.Server{
        Addr:    ":8080",
        Handler: handler,
    }
    
    lifecycle.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            go srv.ListenAndServe()
            return nil
        },
        OnStop: func(ctx context.Context) error {
            return srv.Shutdown(ctx)
        },
    })
}

// @fx invoke
func InitDatabase(db *gorm.DB) error {
    return db.AutoMigrate(&User{}, &Order{})
}
```

### 3. 装饰器 (Decorate)

```go
// middleware/logging.go

// @fx decorate UserServiceInterface  // 装饰接口
func NewLoggingUserService(svc UserServiceInterface, logger *slog.Logger) UserServiceInterface {
    return &LoggingUserService{
        next:   svc,
        logger: logger,
    }
}

// @fx decorate
func NewCachedUserService(svc *UserService, cache *redis.Client) *UserService {
    return &CachedUserService{
        UserService: svc,
        cache:       cache,
    }
}
```

### 4. 应用入口标记

```go
// cmd/app/main.go

// @fx app
func main() {
    // 这个函数体不重要，只是标记在这个包生成 fx_app.go
    app := NewApp()
    app.Run()
}
```

## 🏗️ 生成的代码

### 自动生成的 `cmd/app/fx_app.go`:

```go
package app

import (
    "go.uber.org/fx"
    services "myproject/services"
    middleware "myproject/middleware"
    app1 "myproject/app"
)

// FxOptions 返回所有自动收集的 fx.Option
func FxOptions() []fx.Option {
    return []fx.Option{
        // Provide 选项
        fx.Provide(services.NewUserService),
        fx.Provide(fx.Annotate(
            services.NewUserServiceImpl,
            fx.As(new(UserServiceInterface)),
        )),
        fx.Provide(fx.Annotate(
            services.NewUserHandler,
            fx.ResultTags(`group:"handlers"`),
        )),
        
        // Decorate 选项
        fx.Decorate(fx.Annotate(
            middleware.NewLoggingUserService,
            fx.As(new(UserServiceInterface)),
        )),
        fx.Decorate(middleware.NewCachedUserService),
        
        // Invoke 选项
        fx.Invoke(app1.StartHTTPServer),
        fx.Invoke(app1.InitDatabase),
    }
}

// NewApp 创建配置好的 fx.App
func NewApp(opts ...fx.Option) *fx.App {
    allOpts := append(FxOptions(), opts...)
    return fx.New(allOpts...)
}
```

## 🚀 使用生成的代码

```go
// cmd/app/main.go
func main() {
    // 方式1: 使用生成的 NewApp
    app := NewApp(
        fx.Provide(NewConfig),  // 额外的配置
        fx.WithLogger(func() fxevent.Logger {
            return &fxevent.ConsoleLogger{W: os.Stdout}
        }),
    )
    
    // 方式2: 手动组合
    app := fx.New(
        append(FxOptions(), 
            fx.Provide(NewConfig),
            fx.WithLogger(...),
        )...,
    )
    
    app.Run()
}
```

## 🔗 与其他插件协作

其他 genx 插件生成的代码可以自动集成：

```go
// @log 插件生成的代码会包含 @fx provide 注解
// @fx provide
func NewLogging(logger *slog.Logger) func(next Service) Service {
    return func(next Service) Service {
        return &logging{logger: logger, next: next}
    }
}

// @temporal 插件生成的代码
// @fx provide
func InitTemporal(client client.Client, worker worker.Worker) (*Temporal, error) {
    // ...
}
```

## 📊 注解语法

### @fx provide [interface] [group]
- `interface`: 可选，绑定到指定接口类型
- `group`: 可选，加入到指定组

### @fx invoke
- 标记需要在应用启动时调用的函数

### @fx decorate [interface]
- `interface`: 可选，装饰指定接口类型

### @fx app
- 标记应用入口，在此包生成 fx_app.go

## 🎯 优势

1. **零样板代码**: 自动生成所有 fx.Option 配置
2. **类型安全**: 编译时检查依赖关系
3. **模块化**: 每个包独立定义依赖
4. **可扩展**: 支持额外的 fx.Option
5. **一致性**: 确定性排序保证构建一致性

## 🔄 迁移指南

### 从 @do 迁移到 @fx

```go
// 原来的 @do 注解
// @do provide
func NewUserService() *UserService { ... }

// 改为 @fx 注解  
// @fx provide
func NewUserService() *UserService { ... }

// 原来的手动初始化
func main() {
    injector := do.New()
    doInit(injector)
    svc := do.MustInvoke[*UserService](injector)
}

// 改为 fx 方式
func main() {
    app := NewApp()
    app.Run()
}
```

这样你就可以享受 fx 的强大功能，同时保持 genx 的自动化便利！
