# 状态模式
> 状态模式用于分离状态和行为。

示例：
1. 定义ActionState interface接口，包含：状态所关联的 查看、评论、发布行为
    - View()
    - Comment()
    - Post()
2. 定义ActionState的实现类
   - NormalState：账号正常状态
   - RestrictState：账号受限制状态
   - CloseState：账号被封状态
3. 定义type Account struct
   - state *ActionState
   - healthValue：健康值对应不同的状态
   - 分别实现View、Comment、Post方法，调用state的方法：a.state.View()
   - 实现changeState方法，根据healthValue设置state
4. 定义NewAccount(healthValue int) *Account：创建一个账号并设置健康值
