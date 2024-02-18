# 备忘录模式
> 备忘录模式用于保存程序内部状态到外部，又不希望暴露内部状态的情形。

允许在不破坏封装性的前提下保存和恢复对象的内部状态，程序内部状态使用窄接口传递给外部进行存储，从而不暴露程序实现细节。

备忘录模式同时可以离线保存内部状态，如保存到数据库，文件等。

该模式涉及三个主要角色：
- Originator（发起人）：Originator 是拥有内部状态的对象
- Memento（备忘录）：Memento 是 Originator 的快照
- Caretaker（负责人）：Caretaker 负责备份和恢复 Memento。


示例：
> state暂时用string来代替,实际可用struct
> 案例演示：实现一个文本编辑器，提供撤销和重做功能

1. 定义Originator interface， Originator 发起人：用于保存或者恢复当前状态（备忘录）
   - Save() Memento
   - Restore(m Memento)
2. 定义Memento interface ，提供方法用于获取当前状态【备忘录：记录当前状态】
    - GetState() string
3. 定义TextMemento struct，实现Memento接口，保存state
4. 定义textEditor struct，实现Originator接口
    - state string
    - func Save() Memento 
    - func Restore(m Memento)
    - func SetState(state string)
5. 定义Caretaker struct 
   - mementos     []Memento
   - currentIndex int
   - func AddMemento(m Memento)
   - func Undo(t *TextEditor)
   - func Redo(t *TextEditor)


