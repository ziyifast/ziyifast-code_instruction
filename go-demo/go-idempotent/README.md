Go语言幂等性实现指南

1. 什么是幂等性？
   幂等性（Idempotence） 是数学和计算机科学中的重要概念，指的是一个操作执行多次与执行一次的效果相同。
   在分布式系统和API设计中，幂等性意味着：
- 客户端对同一操作发起一次或多次请求的结果完全一致
- 不会因为多次调用而产生副作用（如重复扣款、重复创建资源等）
- 是构建可靠分布式系统的基石

接口幂等性就是用户对于同一操作发起的一次请求或者多次请求的结果是一致的，不会因为多次点击而产生了副作用。举个最简单的例子，支付过程中，用户购买商品后支付，支付扣款成功，但是返回结果的时候网络异常，此时钱已经扣了，用户再次点击按钮，此时会进行第二次扣款，返回结果成功，用户查询余额返发现多扣钱了，流水记录也变成了两条，这就没有保证接口的幂等性。
简单的说就是 一个用户对于同一个操作发起一次或多次的请求，请求的结果一致。不会因为多次点击而产生多条数据。

2. 为什么需要幂等性？
   2.1 常见场景
- 网络不稳定：客户端超时重试
- 用户重复操作：多次点击提交按钮
- 微服务重试机制：服务间调用自动重试
- 消息队列重复消费：MQ至少一次投递语义

2.2 缺乏幂等性的风险
- ⚠️重复扣款或支付
- ⚠️重复创建订单或资源
- ⚠️数据不一致状态
- ⚠️用户体验差（操作结果不确定）

2.3 需要幂等性的接口
1. Get 方法：幂等。用于获取资源，其一般不会也不应当对系统资源进行改变，所以是幂等的。
2. Post 方法：非幂等。一般用于创建新的资源。其每次执行都会新增数据，所以不是幂等的。
3. Put 方法：视情况而定。一般用于修改资源。该操作则分情况来判断是不是满足幂等，更新操作中直接根据某个值进行更新，也能保持幂等。不过执行累加操作的更新是非幂等。
4. Delete 方法：视情况定。一般用于删除资源。该操作则分情况来判断是不是满足幂等，当根据唯一值进行删除时，删除同一个数据多次执行效果一样。不过需要注意，带查询条件的删除则就不一定满足幂等了。例如在根据条件删除一批数据后，这时候新增加了一条数据也满足条件，然后又执行了一次删除，那么将会导致新增加的这条满足条件数据也被删除。

接口幂等性的处理方式有很多，数据库唯一主键、数据库乐观锁、令牌表+唯一约束、下游传递唯一序列号、同步锁（单体项目）、分布式锁如redis 等等，这里只详细阐述使用token令牌的方式：
[图片]
整体流程：
1. 客户端请求服务器接口获取token。
2. 服务器将token返给客户端的同时将信息（这里包括用户信息、和token）存储到redis中。
3. 请求业务接口时，将token放入header中进行接口请求。
4. 服务器通过用户信息和token检查token是否还存在，如果存在就删除，如果不存在直接返回结果。
5. 响应服务器请求结果。



3. 常见的错误实现方式

3.1 完全没有幂等性保证
func (ps *PaymentService) ProcessPaymentWithoutIdempotent(request PaymentRequest) PaymentResult {
// 直接处理支付，没有幂等性检查
// 每次调用都会执行完整的支付流程
// 用户重复点击会多次扣款
return ps.executePayment(request)
}
问题：用户重复点击会多次扣款。

3.2 简单去重方案
func (ps *PaymentService) ProcessPaymentWithSimpleDedup(request PaymentRequest) PaymentResult {
// 简单用用户ID+金额做key
key := fmt.Sprintf("payment:%d:%.2f", request.UserID, request.Amount)

    // 如果key存在，认为是重复请求
    if _, exists := ps.records[key]; exists {
        return PaymentResult{Status: "DUPLICATE"}
    }
    
    // 处理支付并记录
    result := ps.executePayment(request)
    ps.records[key] = true
    return result
}
问题：
- 用户无法连续支付相同金额
- 处理失败后无法重试
- 无法区分不同业务请求

3.3 考虑不周的幂等token方案
func (ps *PaymentService) ProcessPaymentWithFlawedToken(request PaymentRequest) PaymentResult {
// 检查token是否已使用
if result, exists := ps.records[request.IdempotentToken]; exists {
return *result
}

    // 处理支付逻辑...
    result := ps.executePayment(request)
    
    // 保存结果
    ps.records[request.IdempotentToken] = &result
    return result
}

问题：
- ❌ 存在并发问题，多个请求可能同时通过检查
- ❌ 失败后无法重试（一旦记录，永远返回相同结果）
- ❌ 缺乏状态管理

4. 正确的幂等性实现
   核心原则：
1. 先查后执行：在获取锁之前先查询是否已有处理结果
2. 双重检查：获取锁后再次检查结果，避免并发问题
3. 状态管理：区分处理中、成功、失败等状态
4. 允许重试：失败的请求应允许用户重试（根据业务需求）
5. 超时机制：处理中的请求应有超时时间，防止卡死
   4.1 常见方案介绍
   [图片]
   📢请注意：下面各方案实际流程还需要贴合各自真实业务场景，可适当增加或减少其中的流程
   目前业内主流的幂等性实现方案主要包含以下几种：
1.  select+insert+主键/唯一索引冲突
2. 直接insert + 主键/唯一索引冲突
3. 状态机幂等
4. 抽取防重表
5. token令牌
6. 悲观锁(如select for update)
7. 乐观锁
8. 分布式锁
   ...
   方案
   适用场景
   性能
   复杂度
   一致性
   SELECT+INSERT
   简单业务
   中
   低
   高
   直接INSERT
   高并发
   高
   低
   高
   状态机
   复杂状态流转
   中
   中
   高
   防重表
   复杂业务
   中
   中
   高
   Token令牌
   客户端参与
   高
   中
   高
   悲观锁
   强一致性
   低
   低
   最高
   乐观锁
   读多写少
   高
   中
   高
   分布式锁
   分布式环境
   中
   高
   高

使用建议：
简单场景：使用直接INSERT + 唯一索引
复杂业务：使用状态机 + 防重表
高并发场景：使用Token令牌 + 分布式锁
分布式系统：使用分布式锁 + 消息队列幂等
前端交互：使用Token令牌 + 客户端防重

1. SELECT+INSERT
   概念：交易请求过来，会先根据请求的唯一流水号 bizSeq字段，先select一下数据库的流水表
- 如果数据已经存在，就拦截是重复请求，直接返回成功；
- 如果数据不存在，就执行insert插入，如果insert成功，则直接返回成功，如果insert产生主键冲突异常，则捕获异常，接着直接返回成功。
  流程图：
  [图片]
  伪代码：
  // 幂等处理
  func idempotent(req Request) Response {
  // 1. 查询是否已处理
  record := selectByBizSeq(req.BizSeq)
  if record != nil {
  // 重复请求，直接返回成功
  log.Info("重复请求，直接返回成功，流水号：{}", req.BizSeq)
  return successResponse
  }

  // 2. 插入请求记录
  err := insert(req)
  if err != nil {
  // 3. 主键冲突说明是重复请求
  if isDuplicateKeyError(err) {
  log.Info("主键冲突，是重复请求，直接返回成功，流水号：{}", req.BizSeq)
  return successResponse
  }
  return errorResponse
  }

  // 4. 正常处理请求
  dealRequest(req)

  return successResponse
  }
2. 直接insert + 主键/唯一索引冲突
   概念：在上面的方案中我们都会先查一下流水表的交易请求，判断是否存在，然后不存在再插入请求记录。如果重复请求的概率比较低的话，我们可以直接插入请求，利用主键/唯一索引冲突，去判断是重复请求。
   流程图：
   [图片]
   PS：防重和幂等设计其实是有区别的。防重主要为了避免产生重复数据，把重复请求拦截下来即可。而幂等设计除了拦截已经处理的请求，还要求每次相同的请求都返回一样的效果。不过呢，很多时候，它们的处理流程可以是类似的。

伪代码：
// 幂等处理
func idempotent(req Request) Response {
// 1. 直接插入请求记录
err := insert(req)
if err != nil {
// 2. 主键冲突说明是重复请求
if isDuplicateKeyError(err) {
log.Info("主键冲突，是重复请求，直接返回成功，流水号：{}", req.BizSeq)
return successResponse
}
return errorResponse
}

    // 3. 正常处理请求
    dealRequest(req)
    return successResponse
}

3. 状态机幂等
   概念：很多业务表，都是有状态的，比如转账流水表，就会有0-待处理，1-处理中、2-成功、3-失败状态。转账流水更新的时候，都会涉及流水状态更新，即涉及状态机 (即状态变更图)。我们可以利用状态机实现幂等，一起来看下它是怎么实现的。
   比如转账成功后，把处理中的转账流水更新为成功状态，SQL这么写：
   update transfr_flow set status=2 where biz_seq=‘666’ and status=1;
   流程图：
   [图片]
   伪代码：
   底层原理：状态机是怎么实现幂等的呢？
- 第1次请求来时，bizSeq流水号是 666，该流水的状态是处理中，值是 1，要更新为2-成功的状态，所以该update语句可以正常更新数据，sql执行结果的影响行数是1，流水状态最后变成了2。
- 第2请求也过来了，如果它的流水号还是 666，因为该流水状态已经2-成功的状态了，所以更新结果是0，不会再处理业务逻辑，接口直接返回。
  // 幂等转账处理
  func idempotentTransfer(req Request) Response {
  bizSeq := req.BizSeq

  // 1. 更新状态，只有状态为1的记录才能更新为状态2
  rows := update("transfr_flow set status=2 where biz_seq=? and status=1", bizSeq)

  if rows == 1 {
  log.Info("更新成功,可以处理该请求")
  // 其他业务逻辑处理
  return successResponse
  } else if rows == 0 {
  log.Info("更新不成功，不处理该请求")
  // 不处理，直接返回
  return successResponse
  }

  log.Warn("数据异常")
  return successResponse
  }

4. 抽取防重表
   1和2的方案都是建立在业务流水表上bizSeq的唯一性上。很多时候，我们业务表唯一流水号希望后端系统生成，又或者我们希望防重功能与业务表分隔开来，这时候我们可以单独搞个防重表。当然防重表也是利用主键/索引的唯一性，如果插入防重表冲突即直接返回成功，如果插入成功，即去处理请求。
   概念：防重表是一种独立于业务表的幂等性控制机制，通过单独建立防重表，利用主键或唯一索引的约束特性，请求处理前先插入防重记录，利用数据库约束检测重复，插入失败表示重复请求，直接返回成功；插入成功则继续处理，实现防重逻辑与业务逻辑解耦，适用于后端生成唯一流水号的场景


流程图：
[图片]

伪代码：
// 幂等性处理 - 防重表方式
func handleWithDedupTable(ctx context.Context, requestID string) error {
// 1. 开始事务
tx := db.Begin()

    // 2. 插入防重表记录
    err := tx.Exec("INSERT INTO dedup_table (request_id, create_time) VALUES (?, ?)", 
                   requestID, time.Now()).Error
    
    if err != nil {
        // 3. 插入失败表示重复请求，回滚事务并返回成功
        tx.Rollback()
        return nil  // 幂等性处理：重复请求直接返回成功
    }
    
    // 4. 插入成功，执行业务逻辑
    bizErr := processBusinessLogic(tx, requestID)
    
    // 5. 根据业务处理结果提交或回滚事务
    if bizErr != nil {
        tx.Rollback()
        return bizErr
    }
    
    tx.Commit()
    return nil
}

// 防重表结构示例
/*
CREATE TABLE dedup_table (
request_id VARCHAR(64) PRIMARY KEY,
create_time DATETIME NOT NULL,
INDEX idx_create_time (create_time)
);
*/



5. token令牌
   概念：token 令牌方案一般包括两个请求阶段。阶段一是客户端请求申请获取token，服务端生成token返回。阶段二是客户端带着token请求，服务端校验token。详细流程：
1. 客户端发起请求，申请获取token。
2. 服务端生成全局唯一的token，保存到redis中（一般会设置一个过期时间），然后返回给客户端。
3. 客户端带着token，发起请求。
4. 服务端去redis确认token是否存在，一般用 redis.del(token)的方式，如果存在会删除成功，即处理业务逻辑，如果删除失败不处理业务逻辑，直接返回结果。

流程图：
[图片]
伪代码：
// Token令牌方式幂等处理

// 1. 申请token
func applyToken() string {
// 生成全局唯一token
token := generateUniqueToken()

    // 保存到redis，设置过期时间
    redis.Set(token, "unused", 24*time.Hour)
    
    return token
}

// 2. 幂等处理
func idempotentProcess(req Request) Response {
token := req.Token

    // 校验token是否存在
    deleted := redis.Del(token)
    if deleted > 0 {
        // token存在，处理业务逻辑
        processBusinessLogic(req)
        return successResponse
    } else {
        // token不存在，直接返回成功
        log.Info("token已使用或不存在，幂等返回")
        return successResponse
    }
}

6. 悲观锁(如select for update)
   概念：通俗点讲就是很悲观，每次去操作数据时，都觉得别人中途会修改，所以每次在拿数据的时候都会上锁。官方点讲就是，共享资源每次只给一个线程使用，其它线程阻塞，用完后再把资源转让给其它线程。悲观锁如何控制幂等的呢？就是加锁，一般配合事务来实现。举个更新订单的业务场景
   假设先查出订单，如果查到的是处理中状态，就处理完业务，再然后更新订单状态为完成。如果查到订单，并且是不是处理中的状态，则直接返回

流程图：
[图片]


伪代码：
- 这里面order_id需要是索引或主键，要锁住这条记录就好，如果不是索引或者主键，会锁表的！
- 悲观锁在同一事务操作过程中，锁住了一行数据。别的请求过来只能等待，如果当前事务耗时比较长，就很影响接口性能。所以一般不建议用悲观锁做这个事情。
  begin;  // 1.开始事务
  select * from order where order_id='666' for update // 查询订单，判断状态,锁住这条记录
  if（status !=处理中）{
  //非处理中状态，直接返回；
  return ;
  }
  // 处理业务逻辑
  update order set status='完成' where order_id='666' // 更新完成
  commit; // 5.提交事务



7. 乐观锁
   概念：因为悲观锁有性能问题，这时可以尝试乐观锁。乐观锁在操作数据时,则非常乐观，认为别人不会同时在修改数据，因此乐观锁不会上锁。只是在执行更新的时候判断一下，在此期间别人是否修改了数据。其实就是给表的加多一列version版本号，每次更新记录version都升级一下（version=version+1）。具体流程就是先查出当前的版本号version，然后去更新修改数据时，确认下是不是刚刚查出的版本号，如果是才执行更新
   流程图：
   [图片]
   伪代码：
   //更新前，先查数据，查出的版本号是version =1
   select order_id，version from order where order_id='order_132123'；
   //使用version =1和订单Id一起作为条件，再去更新
   update order set version = version +1，status='P' where  order_id='order_132123' and version =1
   //最后更新成功，才可以处理业务逻辑，如果更新失败，默认为重复请求，直接返回。

8. 分布式锁
   概念：分布式锁实现幂等性的逻辑就是，请求过来时，先去尝试获得分布式锁，如果获得成功，就执行业务逻辑，反之获取失败的话，就舍弃请求直接返回成功/处理中状态。
- 分布式锁可以使用Redis，也可以使用ZooKeeper，不过还是Redis相对好点，因为较轻量级。
- Redis分布式锁，可以使用命令SET EX PX NX + 唯一流水号实现，分布式锁的key必须为业务的唯一标识哈
- Redis执行设置key的动作时，要设置过期时间哈，这个过期时间不能太短，太短拦截不了重复请求，也不能设置太长，会占存储空间。
  流程图：
  [图片]
  伪代码：
  // 幂等性处理 - 分布式锁方式
  func handleIdempotent(businessKey string) error {
  // 1. 尝试获取分布式锁
  lockKey := "idempotent_lock:" + businessKey
  lockValue := generateUniqueValue()

  acquired := redis.SetNX(lockKey, lockValue, 30*time.Second) // 30秒过期
  if !acquired {
  // 2. 获取锁失败，直接返回。提示用户稍后再试或订单正在处理中（实现幂等性）
  return nil
  }

  // 3. 获取锁成功，执行业务逻辑
  defer redis.Del(lockKey) // 释放锁

  // 执行具体业务处理
  return processBusinessLogic()
  }



4.2 整体案例(基于token)
注意：下面代码主要演示整体流程，里面代码并非遵守生产环境的所有规范，实际使用需要结合各自业务场景。
整体流程图：
[图片]
环境准备
因为需要用到分布式锁，因此需要本地有redis环境，可以本地直接安装也可以使用云上环境。我这里通过docker搭建一个redis。
docker run -d --name redis -v /Users/ziyi/docker-home/redis:/data -p 6379:6379 redis
[图片]

项目代码
1. main.go:
   package main

import (
"context"
"crypto/rand"
"encoding/hex"
"fmt"
"github.com/go-redis/redis"
"log"
"math/big"
"net/http"
"strings"
"sync"
"time"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

// PaymentStatus 定义支付状态
type PaymentStatus string

const (
StatusProcessing PaymentStatus = "PROCESSING"
StatusSuccess    PaymentStatus = "SUCCESS"
StatusFailed     PaymentStatus = "FAILED"
)

// PaymentRequest 支付请求结构
type PaymentRequest struct {
IdempotentToken string  `json:"idempotent_token"`
Amount          float64 `json:"amount"`
UserID          string  `json:"user_id"`
Description     string  `json:"description"`
}

// PaymentResult 支付结果结构
type PaymentResult struct {
Status    string `json:"status"`
Message   string `json:"message"`
OrderID   string `json:"order_id,omitempty"`
RequestID string `json:"request_id,omitempty"`
}

// PaymentRecord 支付记录结构
type PaymentRecord struct {
IdempotentKey string         `json:"idempotent_key"`
Status        PaymentStatus  `json:"status"`
Result        *PaymentResult `json:"result,omitempty"`
CreatedAt     time.Time      `json:"created_at"`
UpdatedAt     time.Time      `json:"updated_at"`
ErrorMessage  string         `json:"error_message,omitempty"`
}

// RedisLockService Redis分布式锁服务
type RedisLockService struct {
client redis.Cmdable
ctx    context.Context
}

func NewRedisLockService(client redis.Cmdable) *RedisLockService {
return &RedisLockService{
client: client,
ctx:    context.Background(),
}
}

// generateLockValue 生成锁的唯一标识值
func (r *RedisLockService) generateLockValue() (string, error) {
bytes := make([]byte, 16)
if _, err := rand.Read(bytes); err != nil {
return "", err
}
return hex.EncodeToString(bytes), nil
}

// AcquireLock 获取分布式锁
func (r *RedisLockService) AcquireLock(key string, expiration time.Duration) (string, bool) {
value, err := r.generateLockValue()
if err != nil {
return "", false
}

    // 使用SET命令的NX和EX选项原子性地获取锁
    success, err := r.client.SetNX(key, value, expiration).Result()
    if err != nil {
       return "", false
    }

    if !success {
       return "", false
    }

    return value, true
}

// TryAcquireLock 尝试获取分布式锁，支持等待时间
func (r *RedisLockService) TryAcquireLock(key string, waitTime, expiration time.Duration) (string, bool) {
value, acquired := r.AcquireLock(key, expiration)
if acquired {
return value, true
}

    // 如果获取锁失败，则等待并重试
    deadline := time.Now().Add(waitTime)
    for time.Now().Before(deadline) {
       time.Sleep(50 * time.Millisecond) // 短暂休眠避免过度竞争
       value, acquired := r.AcquireLock(key, expiration)
       if acquired {
          return value, true
       }
    }

    return "", false
}

// ReleaseLock 释放分布式锁（使用Lua脚本确保原子性）
func (r *RedisLockService) ReleaseLock(key, value string) bool {
// Lua脚本原子性地检查并删除锁
luaScript := `
if redis.call("GET", KEYS[1]) == ARGV[1] then
return redis.call("DEL", KEYS[1])
else
return 0
end
`

    result, err := r.client.Eval(luaScript, []string{key}, value).Result()
    if err != nil {
       return false
    }

    return result.(int64) == 1
}

// IsHeldByCurrent 检查锁是否由当前实例持有
func (r *RedisLockService) IsHeldByCurrent(key, value string) bool {
val, err := r.client.Get(key).Result()
if err != nil {
return false
}
return val == value
}

// PaymentRepository 支付记录存储接口
type PaymentRepository interface {
GetPaymentRecord(idempotentKey string) *PaymentRecord
CreatePaymentRecord(idempotentKey string, record *PaymentRecord) error
UpdatePaymentRecord(idempotentKey string, status PaymentStatus, result *PaymentResult, errorMsg string) error
DeletePaymentRecord(idempotentKey string) error
}

// MemoryPaymentRepository 内存存储实现（模拟数据库）
type MemoryPaymentRepository struct {
records map[string]*PaymentRecord
mutex   sync.RWMutex
}

func NewMemoryPaymentRepository() *MemoryPaymentRepository {
return &MemoryPaymentRepository{
records: make(map[string]*PaymentRecord),
}
}

/*
*

     实际根据业务场景来：
     1. 如果只需要短期幂等，那么可以把数据存在redis中，并设置TTL过期时间
     2. 如果需要长期幂等，那么可以把数据存在数据库中
    *
*/
func (m *MemoryPaymentRepository) GetPaymentRecord(idempotentKey string) *PaymentRecord {
m.mutex.RLock()
defer m.mutex.RUnlock()

    if record, exists := m.records[idempotentKey]; exists {
       // 返回副本避免并发修改
       copy := *record
       if record.Result != nil {
          resultCopy := *record.Result
          copy.Result = &resultCopy
       }
       return &copy
    }
    return nil
}

func (m *MemoryPaymentRepository) CreatePaymentRecord(idempotentKey string, record *PaymentRecord) error {
m.mutex.Lock()
defer m.mutex.Unlock()

    if _, exists := m.records[idempotentKey]; exists {
       return fmt.Errorf("记录已存在")
    }

    m.records[idempotentKey] = record
    return nil
}

func (m *MemoryPaymentRepository) UpdatePaymentRecord(idempotentKey string, status PaymentStatus, result *PaymentResult, errorMsg string) error {
m.mutex.Lock()
defer m.mutex.Unlock()

    if record, exists := m.records[idempotentKey]; exists {
       record.Status = status
       record.Result = result
       record.ErrorMessage = errorMsg
       record.UpdatedAt = time.Now()
       return nil
    }
    return fmt.Errorf("记录不存在")
}

func (m *MemoryPaymentRepository) DeletePaymentRecord(idempotentKey string) error {
m.mutex.Lock()
defer m.mutex.Unlock()

    delete(m.records, idempotentKey)
    return nil
}

// PaymentService 支付服务
type PaymentService struct {
lockService LockService
repository  PaymentRepository
requestID   string
}

// LockService 锁服务接口
type LockService interface {
AcquireLock(key string, expiration time.Duration) (string, bool)
TryAcquireLock(key string, waitTime, expiration time.Duration) (string, bool)
ReleaseLock(key, value string) bool
IsHeldByCurrent(key, value string) bool
}

func NewPaymentService(lockService LockService, repository PaymentRepository) *PaymentService {
return &PaymentService{
lockService: lockService,
repository:  repository,
}
}

// ProcessPaymentWithIdempotent 处理幂等性支付请求
func (ps *PaymentService) ProcessPaymentWithIdempotent(request PaymentRequest) PaymentResult {
idempotentKey := request.IdempotentToken

    // 参数校验。通常除了验空之外，还需要验证幂等key的有效性，防篡改等。
    // 这里为了代码简单就跳过对key的有效性校验
    if idempotentKey == "" {
       return PaymentResult{Status: "ERROR", Message: "缺少幂等token"}
    }

    // 1. 先查询是否已有处理结果。
    // 这里可以根据业务场景使用多级缓存，比如 先查本地缓存(L1缓存) -> 查询Redis(L2缓存) -> 再查db等
    if existingRecord := ps.getPaymentRecord(idempotentKey); existingRecord != nil {
       switch existingRecord.Status {
       case StatusSuccess:
          result := *existingRecord.Result
          result.RequestID = ps.requestID
          return result
       case StatusProcessing:
          // 检查是否超时 (通常设置为业务处理超时时间)
          if time.Since(existingRecord.UpdatedAt) > 60*time.Second {
             // 超时，允许重试，更新状态为失败
             ps.updateRecordStatus(idempotentKey, StatusFailed, nil, "处理超时")
          } else {
             return PaymentResult{
                Status:    "PROCESSING",
                Message:   "请求正在处理中",
                RequestID: ps.requestID,
             }
          }
       case StatusFailed:
          // 根据业务决定是否允许重试
          if ps.allowRetry(request) {
             // 清除失败记录，允许重试
             ps.repository.DeletePaymentRecord(idempotentKey)
          } else {
             return PaymentResult{
                Status:    "FAILED",
                Message:   "请求已失败且不可重试",
                RequestID: ps.requestID,
             }
          }
       }
    }

    // 2. 获取分布式锁
    lockKey := "payment_lock:" + idempotentKey
    lockValue, acquired := ps.lockService.TryAcquireLock(lockKey, 3*time.Second, 10*time.Second)
    if !acquired {
       return PaymentResult{
          Status:    "RETRY",
          Message:   "系统繁忙，请稍后重试", // 获取锁失败，可提示用户稍后再试 或 订单正在处理中
          RequestID: ps.requestID,
       }
    }

    // 确保锁被释放
    defer func() {
       if ps.lockService.IsHeldByCurrent(lockKey, lockValue) {
          ps.lockService.ReleaseLock(lockKey, lockValue)
       }
    }()

    // 3. 双重检查：获取锁后再次检查
    if existingRecord := ps.getPaymentRecord(idempotentKey); existingRecord != nil {
       if existingRecord.Status == StatusSuccess {
          result := *existingRecord.Result
          result.RequestID = ps.requestID
          return result
       }
    }

    // 4. 创建处理中记录
    if err := ps.createProcessingRecord(idempotentKey); err != nil {
       return PaymentResult{
          Status:    "ERROR",
          Message:   "创建处理记录失败: " + err.Error(),
          RequestID: ps.requestID,
       }
    }

    // 5. 执行业务逻辑
    var result PaymentResult
    func() {
       defer func() {
          if r := recover(); r != nil {
             result = PaymentResult{
                Status:    "FAILED",
                Message:   "处理过程发生异常",
                RequestID: ps.requestID,
             }
             ps.repository.UpdatePaymentRecord(idempotentKey, StatusFailed, &result, "处理过程发生异常")
          }
       }()

       result = ps.executePayment(request)
       result.RequestID = ps.requestID
       ps.repository.UpdatePaymentRecord(idempotentKey, StatusSuccess, &result, "")
    }()

    return result
}

// getPaymentRecord 获取支付记录
func (ps *PaymentService) getPaymentRecord(idempotentKey string) *PaymentRecord {
return ps.repository.GetPaymentRecord(idempotentKey)
}

// updateRecordStatus 更新记录状态
func (ps *PaymentService) updateRecordStatus(idempotentKey string, status PaymentStatus, result *PaymentResult, errorMsg string) {
ps.repository.UpdatePaymentRecord(idempotentKey, status, result, errorMsg)
}

// createProcessingRecord 创建处理中记录
func (ps *PaymentService) createProcessingRecord(idempotentKey string) error {
record := &PaymentRecord{
IdempotentKey: idempotentKey,
Status:        StatusProcessing,
CreatedAt:     time.Now(),
UpdatedAt:     time.Now(),
}
return ps.repository.CreatePaymentRecord(idempotentKey, record)
}

// allowRetry 是否允许重试（根据业务需求实现）
func (ps *PaymentService) allowRetry(request PaymentRequest) bool {
// 根据业务规则判断是否允许重试，例如检查失败次数等
return true
}

// executePayment 执行支付逻辑（模拟）
func (ps *PaymentService) executePayment(request PaymentRequest) PaymentResult {
log.Printf("开始处理支付请求: 用户=%s, 金额=%.2f, 描述=%s",
request.UserID, request.Amount, request.Description)

    // 模拟支付处理时间（100-500ms）
    n, _ := rand.Int(rand.Reader, big.NewInt(400))
    time.Sleep(time.Duration(100+n.Int64()) * time.Millisecond)

    // 模拟随机失败(5%概率失败)
    failureRate, _ := rand.Int(rand.Reader, big.NewInt(100))
    if failureRate.Int64() < 5 {
       log.Printf("支付处理失败: 用户=%s, 金额=%.2f", request.UserID, request.Amount)
       return PaymentResult{
          Status:  "FAILED",
          Message: "支付网关暂时不可用，请稍后重试",
       }
    }

    orderID := "PAY_" + strings.ToUpper(uuid.New().String()[:8])
    log.Printf("支付处理成功: 用户=%s, 金额=%.2f, 订单=%s",
       request.UserID, request.Amount, orderID)

    return PaymentResult{
       Status:  "SUCCESS",
       Message: fmt.Sprintf("支付成功，金额: %.2f元", request.Amount),
       OrderID: orderID,
    }
}

// HTTP handlers
type PaymentHandler struct {
paymentService *PaymentService
repository     PaymentRepository
}

func NewPaymentHandler(paymentService *PaymentService, repository PaymentRepository) *PaymentHandler {
return &PaymentHandler{
paymentService: paymentService,
repository:     repository,
}
}

func (h *PaymentHandler) ProcessPayment(c *gin.Context) {
var request PaymentRequest
if err := c.ShouldBindJSON(&request); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数: " + err.Error()})
return
}

    // 为每次请求生成唯一ID用于追踪
    h.paymentService.requestID = "REQ_" + strings.ToUpper(uuid.New().String()[:8])

    result := h.paymentService.ProcessPaymentWithIdempotent(request)

    // 根据状态码返回不同的HTTP状态
    var statusCode int
    switch result.Status {
    case "SUCCESS":
       statusCode = http.StatusOK
    case "PROCESSING":
       statusCode = http.StatusAccepted
    case "FAILED", "ERROR":
       statusCode = http.StatusInternalServerError
    case "RETRY":
       statusCode = http.StatusTooManyRequests
    default:
       statusCode = http.StatusOK
    }

    c.JSON(statusCode, result)
}

func (h *PaymentHandler) GetPaymentStatus(c *gin.Context) {
idempotentToken := c.Query("token")
if idempotentToken == "" {
c.JSON(http.StatusBadRequest, gin.H{"error": "缺少token参数"})
return
}

    record := h.repository.GetPaymentRecord(idempotentToken)
    if record == nil {
       c.JSON(http.StatusNotFound, gin.H{"error": "未找到对应的支付记录"})
       return
    }

    c.JSON(http.StatusOK, record)
}

func (h *PaymentHandler) ListAllRecords(c *gin.Context) {
// 这里简单返回所有记录（实际生产中应该分页）
c.JSON(http.StatusOK, h.repository.(*MemoryPaymentRepository).records)
}

func main() {
// 初始化Redis客户端
redisClient := redis.NewClient(&redis.Options{
Addr: "localhost:6379",
DB:   0,
})

    // 测试Redis连接
    if err := redisClient.Ping().Err(); err != nil {
       log.Fatal("无法连接到Redis: ", err)
    }
    log.Println("成功连接到Redis")

    // 初始化服务
    lockService := NewRedisLockService(redisClient)
    repository := NewMemoryPaymentRepository()
    paymentService := NewPaymentService(lockService, repository)
    paymentHandler := NewPaymentHandler(paymentService, repository)

    // 初始化Gin路由器
    r := gin.Default()

    // API路由
    r.POST("/payment", paymentHandler.ProcessPayment)
    r.GET("/payment/status", paymentHandler.GetPaymentStatus)
    r.GET("/payment/records", paymentHandler.ListAllRecords)

    // 健康检查
    r.GET("/health", func(c *gin.Context) {
       c.JSON(http.StatusOK, gin.H{"status": "healthy"})
    })

    log.Println("服务器启动在端口 8080")
    log.Println("API endpoints:")
    log.Println("  POST /payment - 处理支付请求")
    log.Println("  GET  /payment/status?token={token} - 查询支付状态")
    log.Println("  GET  /payment/records - 查看所有支付记录")
    log.Println("  GET  /health - 健康检查")

    if err := r.Run(":8080"); err != nil {
       log.Fatal("服务器启动失败: ", err)
    }
}
2. 测试脚本
   test.sh:
# 1. 生成幂等token
TOKEN=$(uuidgen)
echo "使用token: $TOKEN"

# 2. 发送支付请求
curl -X POST http://localhost:8080/payment \
-H "Content-Type: application/json" \
-d '{
"idempotent_token": "'$TOKEN'",
"amount": 100.50,
"user_id": "user_123",
"description": "测试支付"
}'

# 3. 使用相同token再次发送（测试幂等性）
curl -X POST http://localhost:8080/payment \
-H "Content-Type: application/json" \
-d '{
"idempotent_token": "'$TOKEN'",
"amount": 100.50,
"user_id": "user_123",
"description": "测试支付"
}'

# 4. 查询支付状态
curl "http://localhost:8080/payment/status?token=$TOKEN"

# 5. 查看所有记录
curl http://localhost:8080/payment/records

测试验证
1. 以debug方式运行程序
2. 发起请求
# 1. 生成幂等token
TOKEN=$(uuidgen)
echo "使用token: $TOKEN"

# 2. 发送支付请求
curl -X POST http://localhost:8080/payment \
-H "Content-Type: application/json" \
-d '{
"idempotent_token": "'$TOKEN'",
"amount": 100.50,
"user_id": "user_123",
"description": "测试支付"
}'
3. 首先会去检查当前订单是否有处理结果
   [图片]
   很明显，这是我们第一次请求，因此这笔订单不会有处理结果。然后我们会去请求redis获取锁，因为是该笔订单是第一次进来，所以可以获取锁成功。
   [图片]
   redis自然也会有对应的key：
   [图片]
   最后走订单处理流程，处理订单成功：
   [图片]
   [图片]
4. 我们以相同订单号重新发送请求，验证幂等功能是否有效
# 使用token: B3C908D1-F5E7-4DC3-8BB8-9606138FB85B
TOKEN="B3C908D1-F5E7-4DC3-8BB8-9606138FB85B"

curl -X POST http://localhost:8080/payment \
-H "Content-Type: application/json" \
-d '{
"idempotent_token": "'$TOKEN'",
"amount": 100.50,
"user_id": "user_123",
"description": "测试支付"
}'
5. 可以看到我们已经查询之前的订单记录，且状态为成功，因此会直接返回订单成功，也不会去请求redis
   [图片]

[图片]






5. 实际应用建议
1. 首选数据库约束：只要可能，尽量使用唯一索引作为幂等的最终防线。它是所有方案中最可靠、最简单的。
2. 组合使用：一个完整的幂等方案通常是多种策略的组合。
- 示例：处理支付请求时，可以同时使用：
    - Token方案：防止前端重复提交。
    - 防重表+唯一索引：快速判断请求是否已处理。
    - 状态机+乐观锁：更新订单状态，防止并发更新。
    - 分布式锁：在跨多个服务的复杂业务流程中，锁定关键资源。
3. 根据业务场景选择：
- 创建操作：首选 唯一索引。
- 更新操作：首选 状态机 和 乐观锁。
- 分布式环境：必须引入 分布式锁 或基于Redis的Token方案。
- 前端交互：使用 Token令牌 防止用户重复点击。
  核心思想：幂等性是一个需要在多个层级（数据库、应用、业务）共同考虑的体系化设计，而不是一个单一的技术点。




6. 总结
   幂等性设计的核心要点：
1. 唯一标识：每个请求必须有唯一幂等key
2. 状态管理：明确区分处理中、成功、失败状态
3. 原子操作：检查-执行-保存结果需要原子性
4. 超时机制：处理中的请求需要有超时时间
5. 可重试性：根据业务需求设计合理的重试机制

最佳实践建议：
1. 对于敏感操作（支付、资金变动等）：必须由服务端生成带签名的幂等Token
2. Token应包含：用户信息、时间戳、随机数、签名
3. 对于低风险操作：可考虑接受前端生成的UUID，但需要结合用户会话验证
4. 始终验证：无论Token来自哪里，服务端都必须验证其有效性和合法性