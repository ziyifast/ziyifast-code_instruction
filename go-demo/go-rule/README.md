# [代码规范篇]Go代码规范
![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/95653709744340e2a31e1cc13ec6760e.png)



作为编程后端开发新人，需要从一开始就养成一个比较不错的习惯，编写代码时，需要注意安全、并发、命名等问题，本文将会给大家讲解一些日常Go开发工作中需要注意的点。
> Google官方文档：https://google.github.io/styleguide/go/guide

## 1. 安全性
原则：
- 最小权限原则：只赋予程序运行所需的最小权限。
- 输入校验：对所有外部输入（用户输入、网络请求、文件内容）进行严格的长度、类型、格式校验。
- 敏感信息保护：避免硬编码密钥、密码等敏感信息。使用环境变量或安全的配置管理服务。

```go
// 不推荐：硬编码敏感配置
var dbPassword = "root123"

// 推荐：通过环境变量读取
func getDBPassword() string {
    return os.Getenv("DB_PASSWORD")
}

// 注意：使用HTTPS、避免SQL注入等也是安全的重要部分。
```
## 2. 并发安全
原则
- 明确数据所有权：清晰界定哪些协程可以访问和修改哪些数据。
- 锁范围最小化：减少锁的持有时间，只锁住临界区，防止死锁和不必要的性能损耗。
- 优先使用通道：在适合的场景下，使用 channel 进行协程间通信，而非共享内存。
- 使用 sync.RWMutex：在读多写少的场景下，使用读写锁提升性能。

```go
// 不推荐：在耗时操作中持有锁
func badExample(mu *sync.Mutex, data map[string]int) {
    mu.Lock()
    defer mu.Unlock() // 锁会持续到函数结束
    // 模拟耗时操作（如IO、复杂计算）
    time.Sleep(1 * time.Second)
    data["key"] = 1
}

// 推荐：尽快释放锁，耗时操作不持有锁
func goodExample(mu *sync.Mutex, data *map[string]int, result *int) {
    mu.Lock()
    (*data)["key"] = 1
    mu.Unlock() // 操作完共享数据立即释放锁

    // 耗时操作与共享数据无关，在此执行
    *result = calculateExpensiveResult()
}

// 推荐：使用 RWMutex
func goodExampleWithRWMutex(rwMu *sync.RWMutex, data map[string]int) int {
    rwMu.RLock()         // 读锁
    defer rwMu.RUnlock() // 释放读锁
    return data["key"]   // 多个读协程可以同时进入
}
```
## 3. 资源管理
原则：
- 空指针检查：访问用户传入的指针、接口、映射或切片前必须进行有效性判断。
- 资源释放：使用 defer 及时关闭文件、数据库连接、网络连接等资源，确保即使发生错误也能正确释放。

```go
// 不推荐：未检查指针是否为空
func processUser(u *User) {
    fmt.Println(u.Name) // 如果 u 为 nil，会 panic
}

// 推荐：先检查指针是否为空
func processUser(u *User) {
    if u == nil {
        return // 或返回错误
    }
    fmt.Println(u.Name)
}

// 推荐：使用 defer 管理资源
func readFile(filename string) (string, error) {
    f, err := os.Open(filename)
    if err != nil {
        return "", err
    }
    defer f.Close() // 确保文件一定会被关闭

    data, err := io.ReadAll(f)
    if err != nil {
        return "", err
    }
    return string(data), nil
}
```
## 4. 复杂性控制
原则：
- 简洁优先：用最少的、清晰的代码实现功能。一个函数只做一件事。
- 避免过度设计：不引入不必要的抽象、依赖或复杂设计模式。优先使用标准库。
- 函数长度：如果函数过长（例如超过 50 行），考虑将其拆分为多个更小的函数。

```go
// 不推荐：函数冗长，职责不清
func processDataAndSend(data []byte) error {
    // ... 很长很复杂的解析逻辑
    // ... 很长很复杂的处理逻辑
    // ... 很长很复杂的发送逻辑
    return nil
}

// 推荐：拆分为多个单一职责的函数
func processData(data []byte) (*ProcessedResult, error) {
    // 只负责解析
    return &ProcessedResult{}, nil
}

func sendResult(result *ProcessedResult) error {
    // 只负责发送
    return nil
}

// 主函数变得非常简洁
func mainLogic(data []byte) error {
    result, err := processData(data)
    if err != nil {
        return err
    }
    return sendResult(result)
}
```
## 5. 可测性
原则：
- 单元测试覆盖：关键逻辑和公开函数必须有单元测试。
- 依赖注入与接口：通过定义接口而非依赖具体实现，以便使用 Mock 进行测试。
- “接受接口，返回结构体”：函数参数尽可能使用接口类型，增加灵活性。

```go
// 定义接口以便 mock
type UserDB interface {
    GetUser(id int) (*User, error)
    SaveUser(u *User) error
}

// 业务逻辑层依赖接口
type UserService struct {
    db UserDB
}

func (s *UserService) GetUserProfile(id int) (*Profile, error) {
    user, err := s.db.GetUser(id) // 依赖抽象，而非具体数据库
    // ... 业务逻辑
    return &Profile{}, nil
}

// 测试中使用 mock
type MockDB struct {
    // 实现 UserDB 接口所需的方法
}

func (m *MockDB) GetUser(id int) (*User, error) {
    // 返回预设的测试数据
    return &User{Name: "Test User"}, nil
}
func (m *MockDB) SaveUser(u *User) error { return nil }

func TestGetUserProfile(t *testing.T) {
    mockDB := &MockDB{}
    service := &UserService{db: mockDB}
    profile, err := service.GetUserProfile(1)
    assert.NoError(t, err)
    assert.Equal(t, "Test User", profile.Name) // 基于 mock 数据进行断言
}
```
## 6. 作用域控制
原则：
- 最小作用域：变量应在最近使用的地方声明，并限制其作用域。
- 合理导出：区分可导出（大写开头）与不可导出（小写开头）的标识符。只将需要被外部包使用的部分导出。

```go
// 不推荐：不必要的全局变量
var globalCounter int
func increment() {
    globalCounter++ // 难以测试和维护
}

// 推荐：使用局部变量或依赖注入
func increment(counter *int) {
    *counter++
}

// 不推荐：导出内部工具函数
func FormatTime(t time.Time) string { // 外部包能调用它
    return t.Format("2006-01-02")
}

// 推荐：内部函数不导出
func formatTime(t time.Time) string { // 仅在包内可用
    return t.Format("2006-01-02")
}
```
## 7. 可读性
原则:
- 文档注释：公共的包、类型、函数、常量必须添加 // Comment 注释，说明其目的、行为和参数返回值。
- 清晰命名：变量、函数、包名应具有明确的含义，避免缩写（除非是广为人知的）。
- 代码表达意图：代码本身应该清晰，让注释成为“为什么这么做”的补充，而非“做了什么”的重复。

```go
// 推荐：为公开元素添加详细注释
// ConfigReader is responsible for loading configuration from a file.
// It supports both JSON and YAML formats.
type ConfigReader struct {
    Path string // Path to the config file.
}

// Load reads the configuration from the file specified in Path.
// Returns an error if the file cannot be read or parsed.
func (c *ConfigReader) Load() error {
    // 实现细节...
    return nil
}

// 不推荐：命名模糊
var d int    // 什么的时间？天？ duration?
var f func() // 什么函数？

// 推荐：命名明确，见名知意
var timeoutDuration time.Duration
var onSuccessCallback func()
```

## 8. 错误处理
原则：
- 正确处理错误：不能忽略错误，应返回或记录。
- 具体错误信息：提供具体的错误描述。

```go
// 不推荐：忽略错误
func readFile(filename string) string {
    data, _ := ioutil.ReadFile(filename)
    return string(data)
}

// 推荐：处理错误
func readFile(filename string) (string, error) {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return "", fmt.Errorf("failed to read file %s: %w", filename, err)
    }
    return string(data), nil
}
```
## 9. 命名规范
原则：
- 表意明确：命名应清晰表达其用途。
- 避免冲突：变量名不应与包名重复。

```go
// 不推荐：命名模糊
var tmp string
var data []byte

// 推荐：命名明确
var userName string
var activeConnections int

// 不推荐：变量名与包名相同
package user
var user User // 混淆：user.User 指的是变量还是类型？

// 推荐：区分变量名与包名
package user
var currentUser User // 清晰：user.currentUser

// 推荐：接口命名
type Fetcher interface {
    Fetch(url string) ([]byte, error)
}
```
## 10. 一致性维护
原则：
- 注释与代码一致：确保注释描述与实际代码行为相符。修改代码后务必更新注释。
- 统一命名：同一概念在整个项目或包中应使用相同的名称和术语。
- 格式化：使用 gofmt 或 goimports 工具统一代码格式。

```go
// 不推荐：注释与实现不符
// GetActiveUsers returns all users
func GetActiveUsers() []User {
    // 实际返回所有用户，包括非活跃用户
    return allUsers
}

// 推荐：注释与实现一致
// GetAllUsers returns all users
func GetAllUsers() []User {
    return allUsers
}
```
## 11. 惯用法遵循
原则：
- context 使用：关键函数使用 context 控制超时。
- Go 风格编码：遵循 Go 语言惯用法。

```go
// 不推荐：未使用 context
func fetchData() (*Data, error) {
    // 数据获取逻辑
    return &Data{}, nil
}

// 推荐：使用 context
func fetchData(ctx context.Context) (*Data, error) {
    // 数据获取逻辑
    return &Data{}, nil
}

// 不推荐：使用 panic 处理错误
func divide(a, b int) int {
    if b == 0 {
        panic("division by zero")
    }
    return a / b
}

// 推荐：返回错误
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}
```
## 12. 格式规范
原则：
- import 分类排序：按标准库、第三方、本地库分组并排序。
- 合理空行：避免不必要的空行。

```go
// 不推荐：import 无分类
import (
    "github.com/gin-gonic/gin"
    "fmt"
    "myproject/utils"
    "net/http"
)

// 推荐：import 分类排序
import (
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"

    "myproject/utils"
)

// 不推荐：无意义空行
func example() {
    
    fmt.Println("Hello")
    
    
    fmt.Println("World")
    
}
```
参考资料：
- https://go-lang.org.cn/ref/spec
- https://google.github.io/styleguide/go/guide
- https://github.com/xxjwxc/uber_go_guide_cn