# 批量优化所有插件的错误处理

## 需要优化的插件列表

基于代码扫描，以下插件需要进行错误处理优化：

### 1. **kithttp 插件** (plugs/kithttp/)
- **文件**: kit.go, request.go, response.go, observer.go, permission.go
- **问题**: 多处使用 panic，缺少结构化错误处理
- **优先级**: 高

### 2. **gormq 插件** (plugs/gormq/)
- **文件**: gen.go, plug.go
- **问题**: panic 调用，错误信息不详细
- **优先级**: 高

### 3. **crud 插件** (plugs/crud/)
- **文件**: gen.go, plug.go
- **问题**: 已部分优化，需要完善
- **优先级**: 中

### 4. **temporal 插件** (plugs/temporal/)
- **文件**: gen.go, plug.go
- **问题**: panic 调用，缺少上下文信息
- **优先级**: 中

### 5. **otel 插件** (plugs/otel/)
- **文件**: gen.go, plug.go
- **问题**: 错误处理粗糙
- **优先级**: 中

### 6. **enum 插件** (plugs/enum/)
- **文件**: gen.go, plug.go
- **问题**: 基础错误处理需要改进
- **优先级**: 低

### 7. **copy 插件** (plugs/copy/)
- **文件**: gen.go, plug.go
- **问题**: 简单的错误处理
- **优先级**: 低

### 8. **alert 插件** (plugs/alert/)
- **文件**: gen.go, plug.go
- **问题**: 错误处理需要标准化
- **优先级**: 低

### 9. **log 插件** (plugs/log/)
- **文件**: gen.go, plug.go
- **问题**: 基础错误处理
- **优先级**: 低

### 10. **trace 插件** (plugs/trace/)
- **文件**: gen.go, plug.go
- **问题**: 错误处理需要改进
- **优先级**: 低

### 11. **do 插件** (plugs/do/)
- **文件**: gen.go, plug.go
- **问题**: 全局插件，需要特殊处理
- **优先级**: 中

## 优化模式

### 通用优化模式

1. **替换 panic 调用**
```go
// 之前
if err != nil {
    panic(err)
}

// 优化后
if err != nil {
    return common.ParseError("operation failed").
        WithCause(err).
        WithPlugin("@plugin-name").
        WithDetails("detailed description").
        Build()
}
```

2. **更新方法签名**
```go
// 之前
func (p *Plugin) Parse() {
    // ...
}

// 优化后
func (p *Plugin) Parse() error {
    // ...
    return nil
}
```

3. **添加上下文信息**
```go
return common.ValidationError("invalid annotation").
    WithPlugin("@plugin-name").
    WithInterface(interfaceName).
    WithMethod(methodName).
    WithAnnotation("@annotation-name").
    WithDetails("specific error description").
    Build()
```

4. **使用安全执行包装器**
```go
err := common.WithRecovery(func() error {
    // 可能会 panic 的代码
    return riskyOperation()
})
```

### 插件特定优化

#### kithttp 插件
- 重点优化 HTTP 路由解析错误
- 添加请求/响应类型验证错误
- 改进权限检查错误处理

#### gormq 插件
- 优化 GORM 查询构建错误
- 添加字段映射验证
- 改进 SQL 生成错误处理

#### crud 插件
- 完善 CRUD 操作错误处理
- 添加模型验证错误
- 改进关联关系错误处理

#### temporal 插件
- 优化工作流定义错误
- 添加活动验证错误
- 改进时间处理错误

## 实施计划

### 阶段 1: 高优先级插件 (1-2天)
1. ✅ kithttpclient (已完成)
2. 🔄 kithttp
3. 🔄 gormq
4. 🔄 crud (完善)

### 阶段 2: 中优先级插件 (2-3天)
1. temporal
2. otel
3. do

### 阶段 3: 低优先级插件 (1-2天)
1. enum
2. copy
3. alert
4. log
5. trace

### 阶段 4: 测试和验证 (1天)
1. 集成测试
2. 错误场景测试
3. 文档更新

## 验证清单

对于每个插件，确保：

- [ ] 所有 panic 调用都被替换
- [ ] 方法签名正确返回错误
- [ ] 错误信息包含足够的上下文
- [ ] 使用正确的错误代码
- [ ] 添加了详细的错误描述
- [ ] 错误链正确传播
- [ ] 位置信息准确
- [ ] 与 TUI 系统集成

## 测试策略

### 单元测试
```go
func TestPluginErrorHandling(t *testing.T) {
    // 测试各种错误场景
    // 验证错误信息的准确性
    // 检查错误代码的正确性
}
```

### 集成测试
```go
func TestPluginIntegration(t *testing.T) {
    // 测试插件与核心系统的集成
    // 验证错误传播
    // 检查 TUI 显示
}
```

### 错误场景测试
- 缺少注解
- 错误的注解格式
- 无效的类型定义
- 文件读写错误
- 网络错误（如适用）

## 预期效果

优化完成后，整个 GenX 项目将具备：

1. **统一的错误处理** - 所有插件使用相同的错误处理模式
2. **详细的错误信息** - 每个错误都包含足够的上下文和位置信息
3. **优雅的错误恢复** - 不再有突然的 panic 崩溃
4. **用户友好的错误显示** - TUI 中美观的错误展示
5. **开发者友好的调试** - 清晰的错误堆栈和上下文
6. **可维护的代码** - 标准化的错误处理模式

## 后续维护

1. **代码审查标准** - 确保新代码遵循错误处理规范
2. **文档更新** - 保持错误处理文档的最新状态
3. **持续改进** - 根据用户反馈优化错误信息
4. **监控和分析** - 收集错误统计信息，识别常见问题

这个优化将显著提升 GenX 项目的稳定性、可用性和开发体验。
