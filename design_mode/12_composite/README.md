# 组合模式
> 组合模式统一对象和对象集，使得使用相同接口使用对象和对象集。

组合模式常用于树状结构，用于统一叶子节点和树节点的访问，并且可以用于应用某一操作到所有子节点。

示例：以飞书文档接口为例，一个目录下面可以包含多个子文件
> 我们最先想到的做法就是：将文件和目录放在一个类中，新增一个字段，用于判断是文件还是目录。但这样并不优雅。因为文件和目录是不同的，各自有各自的特性，将特有的内容放到一个类里，不满足单一职责原则。
> 下面将展示通过组合模式来实现文档管理结构
1. 定义FileSystem interface，抽取文档和目录的公共部分
    - Display(separator string)
2. 定义FileCommon struct，抽取文件和目录的公共部分
    - fileName string
    - SetFileName(fileName string)
3. 定义FileNode struct，用于表示文件类。并实现Display方法
    - FileCommon
4. 定义DirectoryNode struct，用于表示目录类。并实现Display方法。因为目录下可以存放多个文件，因此需要提供addFile方法
    - FileCommon
    - nodes []FileSystemNode
    - Add(f FileSystemNode)