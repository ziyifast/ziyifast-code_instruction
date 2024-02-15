# 模版方法模式
> 定义一个操作中的算法的骨架，而将一些步骤延迟到子类中。模板方法使得子类可以不改变一个算法的结构即可重定义该算法的某些特定步骤。
> - 通用步骤在抽象类中实现，变化的步骤在具体的子类中实现
> - 做饭，打开煤气，开火，（做饭）， 关火，关闭煤气。除了做饭其他步骤都是相同的，抽到抽象类中实现

示例：
1. 定义type Cooker interface{}，包含做饭的全部步骤，open()、openfire()、cook()、close()、closefire()等
2. 定义type CookMenu struct抽象类，实现做饭的通用步骤，如：打开煤气、打开开关、关闭煤气、关闭开关。具体的做饭内容cook()在子类中实现
3. 定义type ChaojiDan struct，继承CookMenu（属性里包含CookMenu），实现具体做饭的步骤：炒鸡蛋
    - CookMenu
    - func (ChaoJiDan) cook() {
         fmt.Println("做炒鸡蛋")
      }

    

