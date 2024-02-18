# 命令模式

> 命令模式本质是把某个对象的方法调用封装到对象中，方便传递、存储、调用。

示例中把主板单中的启动(start)方法和重启(reboot)方法封装为命令对象，再传递到主机(box)对象中。于两个按钮进行绑定：

第一个机箱(box1)设置按钮1(button1) 为开机按钮2(button2)为重启。
第二个机箱(box1)设置按钮2(button2) 为开机按钮1(button1)为重启。
从而得到配置灵活性。

除了配置灵活外，使用命令模式还可以用作：

批处理
任务队列
undo, redo
等把具体命令封装到对象中使用的场合

示例：
1. 定义ICommand interface接口，包含execute()方法
2. 定义Invoker struct，包含commands []ICommand
    - 包含Call()：用于执行commands中的所有方法
3. 定义TV struct，包含shutdown和turnOn方法
4. 定义ShutdownCommand struct，包含tv，实现execute()方法
    - t.tv.Shutdown()
5. 定义TurnOnCommand struct，包含tv，实现execute()方法
    - t.tv.turnOn()
