package main

import (
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	qt6 "github.com/mappu/miqt/qt"
	"io/ioutil"
	"os"
	"strings"
)

var tableCounter int

func initializeData(model *qt6.QStandardItemModel) {
	rootItem := model.InvisibleRootItem()
	const topLevelRowsToMake = 100

	for i := 0; i < topLevelRowsToMake; i++ {
		// 生成唯一ID
		//uniqueID := uuid.New().String()

		// 创建顶级节点
		item := qt6.NewQStandardItem2(fmt.Sprintf("Row %d", i+1))
		description := qt6.NewQStandardItem2(fmt.Sprintf("Some longer content for Row %d", i+1))
		rootItem.AppendRow([]*qt6.QStandardItem{item, description})

		if i%10 == 3 {
			// 创建一个容器子节点
			subItem := qt6.NewQStandardItem2(fmt.Sprintf("Sub Row Container %d", i+1))
			subDescription := qt6.NewQStandardItem2(fmt.Sprintf("Container for Row %d", i+1))
			item.AppendRow([]*qt6.QStandardItem{subItem, subDescription})

			// 创建子节点的子节点
			for j := 0; j < 5; j++ {
				//uniqueChildID := uuid.New().String()
				child := qt6.NewQStandardItem2(fmt.Sprintf("Sub Row %d", j+1))
				childDescription := qt6.NewQStandardItem2(fmt.Sprintf("Some longer content for Sub Row %d", j+1))
				subItem.AppendRow([]*qt6.QStandardItem{child, childDescription})

				if j < 2 {
					// 创建更深层次的子节点
					//uniqueSubChildID1 := uuid.New().String()
					subChild1 := qt6.NewQStandardItem2(fmt.Sprintf("Sub Sub Row %d", 1))
					subChild1Description := qt6.NewQStandardItem2(fmt.Sprintf("Some longer content for Sub Sub Row %d", 1))
					child.AppendRow([]*qt6.QStandardItem{subChild1, subChild1Description})

					//uniqueSubChildID2 := uuid.New().String()
					subChild2 := qt6.NewQStandardItem2(fmt.Sprintf("Sub Sub Row %d", 2))
					subChild2Description := qt6.NewQStandardItem2(fmt.Sprintf("Some longer content for Sub Sub Row %d", 2))
					child.AppendRow([]*qt6.QStandardItem{subChild2, subChild2Description})
				}
			}
		}
	}
}

// 其余代码保持不变
type CustomDelegate struct {
	*qt6.QStyledItemDelegate
	m_Grid  bool
	m_Color *qt6.QColor
	m_Tree  bool
}

func NewCustomDelegate(parent *qt6.QObject) *CustomDelegate {
	return &CustomDelegate{
		QStyledItemDelegate: qt6.NewQStyledItemDelegate2(parent),
		m_Grid:              true,
		m_Color:             qt6.NewQColor3(0, 0, 255), // 设置线条颜色为黑色
		m_Tree:              true,
	}
}

func (d *CustomDelegate) Paint(painter *qt6.QPainter, option *qt6.QStyleOptionViewItem, index *qt6.QModelIndex) {
	d.QStyledItemDelegate.Paint(painter, option, index)

	if d.m_Grid {
		painter.Save()
		pen := qt6.NewQPen3(d.m_Color)
		painter.SetPen(pen.Color())

		// 绘制垂直虚线
		//painter.DrawLine(option.Rect.Right(), option.Rect.Top(), option.Rect.Right(), option.Rect.Bottom())

		// 绘制底部水平虚线
		//left := option.Rect.Left()
		if d.m_Tree && index.Column() == 0 {
			//left += 24 // 根据需要调整左侧偏移
		}
		//painter.DrawLine(left, option.Rect.Bottom(), option.Rect.Right(), option.Rect.Bottom())

		painter.Restore()
	}
}

//vline.svg branch-more.svg branch-end.svg

//go:embed  CircledChevronDown.svg CircledChevronRight.svg
var resources embed.FS

// encodeResource 将资源文件编码为Base64字符串
func encodeResource(filename string) string {
	file, err := resources.Open(filename)
	if err != nil {
		fmt.Println("无法打开资源文件:", err)
		return ""
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("无法读取资源文件:", err)
		return ""
	}

	return base64.StdEncoding.EncodeToString(data)
}

func myRealMainFunc() {
	qt6.NewQApplication(os.Args)

	// 创建主窗口
	window := qt6.NewQMainWindow2()
	window.SetWindowTitle("Tree View Demo")
	window.SetMinimumSize2(800, 600)

	// 创建布局
	layout := qt6.NewQVBoxLayout(nil)

	// 创建模型
	model := qt6.NewQStandardItemModel()
	model.SetHorizontalHeaderLabels([]string{"名称", "描述"})

	// 创建树形视图
	treeView := qt6.NewQTreeView(nil)
	treeView.SetSortingEnabled(true)
	treeView.SetAlternatingRowColors(true)
	treeView.SetAllColumnsShowFocus(true)
	treeView.SetModel(model.QAbstractItemModel)
	treeView.SetAcceptDrops(true)
	treeView.SetDragEnabled(true)
	treeView.SetDropIndicatorShown(true)
	treeView.SetUniformRowHeights(true)
	treeView.SetSelectionMode(qt6.QAbstractItemView__ExtendedSelection)
	// 设置样式表
	// 设置样式表以绘制虚线

	//treeView.SetStyleSheet(`
	//	QTreeView::branch:has-siblings:!adjoins-item {
	//		border-image: url(data:image/svg+xml;base64,` + encodeResource("vline.svg") + `); /* 绘制垂直虚线 */
	//	}
	//	QTreeView::branch:has-siblings:adjoins-item {
	//		border-image: url(data:image/svg+xml;base64,` + encodeResource("branch-more.svg") + `); /* 连接下级节点的线条 */
	//	}
	//	QTreeView::branch:!has-children:!has-siblings:adjoins-item {
	//		border-image: url(data:image/svg+xml;base64,` + encodeResource("branch-end.svg") + `); /* 最后一个节点 */
	//	}
	//	QTreeView::branch:has-children:!has-siblings:closed,
	//	QTreeView::branch:closed:has-children:has-siblings {
	//		image: url(data:image/svg+xml;base64,` + encodeResource("CircledChevronDown.svg") + `); /* 折叠图标 */
	//	}
	//	QTreeView::branch:open:has-children:!has-siblings,
	//	QTreeView::branch:open:has-children:has-siblings {
	//		image: url(data:image/svg+xml;base64,` + encodeResource("CircledChevronRight.svg") + `); /* 展开图标 */
	//	}
	//`)
	// 创建并设置自定义委托
	//delegate := NewCustomDelegate(treeView.QObject)
	//treeView.SetItemDelegate(delegate.QAbstractItemDelegate)

	// 处理拖放事件
	treeView.OnDropEvent(func(super func(event *qt6.QDropEvent), event *qt6.QDropEvent) {
		if event.Source() == treeView.QObject {
			super(event)
			return
		}
		if event.IsAccepted() {
			urls := event.MimeData().Urls()
			if len(urls) > 0 {
				filePath := urls[0].Path()
				if strings.HasSuffix(filePath, ".json") {
					loadDataFromJSON(model, filePath)
					treeView.ExpandAll()
				}
			}
		}
	})

	treeView.SetSelectionMode(qt6.QAbstractItemView__ExtendedSelection)

	layout.AddWidget3(treeView.QWidget, 1, qt6.AlignLeft)

	// 添加输入框和标签
	//nameLabel := qt6.NewQLabel6("名称:", nil, 0)
	nameLineEdit := qt6.NewQLineEdit(nil)
	//descriptionLabel := qt6.NewQLabel6("描述:", nil, 0)
	descriptionLineEdit := qt6.NewQLineEdit(nil)

	//layout.AddWidget3(nameLabel.QWidget, 0, qt6.AlignLeft)
	//layout.AddWidget3(nameLineEdit.QWidget, 0, qt6.AlignLeft)
	//layout.AddWidget3(descriptionLabel.QWidget, 0, qt6.AlignLeft)
	//layout.AddWidget3(descriptionLineEdit.QWidget, 0, qt6.AlignLeft)
	//
	//// 添加按钮
	//addButton := qt6.NewQPushButton3("添加")
	//editButton := qt6.NewQPushButton3("编辑")
	//deleteButton := qt6.NewQPushButton3("删除")
	//
	//layout.AddWidget3(addButton.QWidget, 0, qt6.AlignLeft)
	//layout.AddWidget3(editButton.QWidget, 0, qt6.AlignLeft)
	//layout.AddWidget3(deleteButton.QWidget, 0, qt6.AlignLeft)

	// 初始化一些预填充的数据
	initializeData(model)
	treeView.ExpandAll()
	// 调整列宽以适应内容 SizeColumnsToFit
	treeView.ResizeColumnToContents(0)
	treeView.ResizeColumnToContents(1)

	// 设置主窗口的中心部件
	//centralWidget := qt6.NewQWidget(nil)
	//centralWidget.SetLayout(layout.Layout())
	//window.SetCentralWidget(centralWidget)
	window.SetCentralWidget(treeView.QWidget)

	// 连接信号和槽
	//addButton.OnClicked1(func(checked bool) {
	//	addNode(model, treeView, nameLineEdit, descriptionLineEdit)
	//})
	//
	//editButton.OnClicked1(func(checked bool) {
	//	editNode(model, treeView, nameLineEdit, descriptionLineEdit)
	//})
	//
	//deleteButton.OnClicked1(func(checked bool) {
	//	deleteNode(model, treeView)
	//})

	// 为树形视图设置右键菜单
	treeView.OnCustomContextMenuRequested(func(pos *qt6.QPoint) {
		menu := qt6.NewQMenu2()
		menu.AddAction("添加") //todo add icon
		menu.AddAction("编辑")
		menu.AddAction("删除")
		menu.AddAction("保存到JSON")

		menu.Actions()[0].OnTriggered1(func(checked bool) {
			addNode(model, treeView, nameLineEdit, descriptionLineEdit)
		})
		menu.Actions()[1].OnTriggered1(func(checked bool) {
			editNode(model, treeView, nameLineEdit, descriptionLineEdit)
		})
		menu.Actions()[2].OnTriggered1(func(checked bool) {
			deleteNode(model, treeView)
		})
		menu.Actions()[3].OnTriggered1(func(checked bool) {
			saveDataToJSON(model, "tree_data.json")
		})
		menu.ExecWithPos(qt6.QCursor_Pos())
	})

	// 设置右键菜单策略
	treeView.SetContextMenuPolicy(qt6.CustomContextMenu)

	window.Show()
	qt6.QApplication_Exec()
}

func addNode(model *qt6.QStandardItemModel, treeView *qt6.QTreeView, nameLineEdit *qt6.QLineEdit, descriptionLineEdit *qt6.QLineEdit) {
	name := nameLineEdit.Text()
	description := descriptionLineEdit.Text()
	if name != "" && description != "" {
		index := treeView.CurrentIndex()
		var parentItem *qt6.QStandardItem
		if index.IsValid() {
			parentItem = model.ItemFromIndex(index)
		} else {
			parentItem = model.InvisibleRootItem()
		}

		newItem := qt6.NewQStandardItem2(name)
		newDescriptionItem := qt6.NewQStandardItem2(description)
		parentItem.AppendRow([]*qt6.QStandardItem{newItem, newDescriptionItem})

		nameLineEdit.Clear()
		descriptionLineEdit.Clear()
	} else {
		qt6.QMessageBox_Warning2(nil, "警告", "名称和描述不能为空", int(qt6.QMessageBox__NoButton), int(qt6.QMessageBox__NoButton))
	}
}

func editNode(model *qt6.QStandardItemModel, treeView *qt6.QTreeView, nameLineEdit *qt6.QLineEdit, descriptionLineEdit *qt6.QLineEdit) {
	index := treeView.CurrentIndex()
	if index.IsValid() {
		name := nameLineEdit.Text()
		description := descriptionLineEdit.Text()
		if name != "" && description != "" {
			item := model.ItemFromIndex(index)
			item.SetText(name)

			descriptionIndex := model.Index(index.Row(), 1, index.Parent())
			descriptionItem := model.ItemFromIndex(descriptionIndex)
			descriptionItem.SetText(description)

			nameLineEdit.Clear()
			descriptionLineEdit.Clear()
		} else {
			qt6.QMessageBox_Warning2(nil, "警告", "名称和描述不能为空", int(qt6.QMessageBox__NoButton), int(qt6.QMessageBox__NoButton))
		}
	} else {
		qt6.QMessageBox_Warning2(nil, "警告", "请选择一个项进行编辑", int(qt6.QMessageBox__NoButton), int(qt6.QMessageBox__NoButton))
	}
}

func deleteNode(model *qt6.QStandardItemModel, treeView *qt6.QTreeView) {
	index := treeView.CurrentIndex()
	if index.IsValid() {
		item := model.ItemFromIndex(index)
		if item != nil {
			model.RemoveRow(index.Row())
		}
	} else {
		qt6.QMessageBox_Warning2(nil, "警告", "请选择一个项进行删除", int(qt6.QMessageBox__NoButton), int(qt6.QMessageBox__NoButton))
	}
}

func loadDataFromJSON(model *qt6.QStandardItemModel, filePath string) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		qt6.QMessageBox_Warning2(nil, "错误", "无法读取文件", int(qt6.QMessageBox__NoButton), int(qt6.QMessageBox__NoButton))
		return
	}

	var treeData TreeData
	err = json.Unmarshal(data, &treeData)
	if err != nil {
		qt6.QMessageBox_Warning2(nil, "错误", "无法解析JSON文件", int(qt6.QMessageBox__NoButton), int(qt6.QMessageBox__NoButton))
		return
	}

	model.Clear()
	model.SetHorizontalHeaderLabels([]string{"名称", "描述"})
	rootItem := model.InvisibleRootItem()

	for _, node := range treeData.Nodes {
		addNodeRecursively(rootItem, node)
	}
}

func addNodeRecursively(parentItem *qt6.QStandardItem, node TreeNode) {
	newItem := qt6.NewQStandardItem2(node.Name)
	newDescriptionItem := qt6.NewQStandardItem2(node.Description)
	parentItem.AppendRow([]*qt6.QStandardItem{newItem, newDescriptionItem})

	for _, childNode := range node.Children {
		addNodeRecursively(newItem, childNode)
	}
}

func saveDataToJSON(model *qt6.QStandardItemModel, filePath string) {
	rootItem := model.InvisibleRootItem()
	nodes := getNodesRecursively(rootItem)

	treeData := TreeData{Nodes: nodes}
	data, err := json.MarshalIndent(treeData, "", "    ")
	if err != nil {
		qt6.QMessageBox_Warning2(nil, "错误", "无法转换为JSON", int(qt6.QMessageBox__NoButton), int(qt6.QMessageBox__NoButton))
		return
	}

	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		qt6.QMessageBox_Warning2(nil, "错误", "无法保存文件", int(qt6.QMessageBox__NoButton), int(qt6.QMessageBox__NoButton))
		return
	}

	qt6.QMessageBox_Information5(nil, "信息", "文件保存成功", qt6.QMessageBox__NoButton, qt6.QMessageBox__NoButton)
}

func getNodesRecursively(item *qt6.QStandardItem) []TreeNode {
	var nodes []TreeNode
	rowCount := item.RowCount()

	for i := 0; i < rowCount; i++ {
		nameItem := item.Child(i)
		descriptionItem := item.Child(i + 1)
		node := TreeNode{
			Name:        nameItem.Text(),
			Description: descriptionItem.Text(),
		}

		if nameItem.HasChildren() {
			node.Children = getNodesRecursively(nameItem)
		}

		nodes = append(nodes, node)
	}

	return nodes
}

type TreeData struct {
	Nodes []TreeNode `json:"nodes"`
}

type TreeNode struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Children    []TreeNode `json:"children,omitempty"`
}
