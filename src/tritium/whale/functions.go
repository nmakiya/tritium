package whale

import (
	"strings"
	"os"
	"gokogiri/html"
	"gokogiri/xml"
	"gokogiri/xpath"
	"fmt"
	//log "log4go"
	tp "athena/src/athena/proto"
	"rubex/lib"
	"css2xpath"
	"goconv"
)

func this_(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	returnValue = scope.Value
	return
}

func yield_(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	myYieldBlock := ctx.yieldBlock()
	ctx.Yields = ctx.Yields[:(len(ctx.Yields) - 1)]
	if ctx.yieldBlock() != nil {
		returnValue = ctx.runChildren(scope, myYieldBlock.Ins)
		if returnValue == nil {
			returnValue = "false"
		}
	} else {
		ctx.Log.Error("yield() failure")
	}
	ctx.Yields = append(ctx.Yields, myYieldBlock)
	return
}

func var_Text(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	val := ctx.Env[args[0].(string)]
	returnValue = val
	if len(ins.Children) > 0 {
		ts := &Scope{Value: returnValue}
		ctx.runChildren(ts, ins)
		returnValue = ts.Value
		ctx.Env[args[0].(string)] = returnValue.(string)
	}
	return
}

func var_Text_Text(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	returnValue = args[1].(string)
	ctx.Env[args[0].(string)] = returnValue.(string)

	if len(ins.Children) > 0 {
		ts := &Scope{Value: returnValue}
		ctx.runChildren(ts, ins)
		returnValue = ts.Value
		ctx.Env[args[0].(string)] = returnValue.(string)
	}
	return
}

func match_Text(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	// Setup stacks
	against, ok := args[0].(string)
	if !ok {
		ctx.Log.Error("AH!")
	}
	ctx.MatchStack = append(ctx.MatchStack, against)
	ctx.MatchShouldContinue = append(ctx.MatchShouldContinue, true)

	// Run children
	ctx.runChildren(scope, ins)

	if ctx.matchShouldContinue() {
		returnValue = "false"
	} else {
		returnValue = "true"
	}

	// Clear
	ctx.MatchShouldContinue = ctx.MatchShouldContinue[:len(ctx.MatchShouldContinue)-1]
	ctx.MatchStack = ctx.MatchStack[:len(ctx.MatchStack)-1]
	return
}

func with_Text(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	returnValue = "false"
	if ctx.matchShouldContinue() {
		if args[0].(string) == ctx.matchTarget() {
			ctx.MatchShouldContinue[len(ctx.MatchShouldContinue)-1] = false
			ctx.runChildren(scope, ins)
			returnValue = "true"
		}
	}
	return
}

func with_Regexp(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	returnValue = "false"
	if ctx.matchShouldContinue() {
		//println(matcher.MatchAgainst, matchWith)
		if (args[0].(*rubex.Regexp)).Match([]uint8(ctx.matchTarget())) {
			ctx.MatchShouldContinue[len(ctx.MatchShouldContinue)-1] = false
			ctx.runChildren(scope, ins)
			returnValue = "true"
		}
	}
	return
}

func not_Text(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	returnValue = "false"
	if ctx.matchShouldContinue() {
		if args[0].(string) != ctx.matchTarget() {
			ctx.MatchShouldContinue[len(ctx.MatchShouldContinue)-1] = false
			ctx.runChildren(scope, ins)
			returnValue = "true"
		}
	}
	return
}

func not_Regexp(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	returnValue = "false"
	if ctx.matchShouldContinue() {
		//println(matcher.MatchAgainst, matchWith)
		if !(args[0].(*rubex.Regexp)).Match([]uint8(ctx.matchTarget())) {
			ctx.MatchShouldContinue[len(ctx.MatchShouldContinue)-1] = false
			ctx.runChildren(scope, ins)
			returnValue = "true"
		}
	}
	return
}

func else_(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	returnValue = "false"
	if ctx.matchShouldContinue() {
		ctx.MatchShouldContinue[len(ctx.MatchShouldContinue)-1] = false
		ctx.runChildren(scope, ins)
		returnValue = "true"
	}
	return
}

func regexp_Text_Text(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	pattern := args[0].(string) + "/" + args[1].(string)
	if r := ctx.RegexpCache[pattern]; r != nil {
		returnValue = r
	} else {
		mode := rubex.ONIG_OPTION_DEFAULT
		if strings.Index(args[1].(string), "i") >= 0 {
			mode = rubex.ONIG_OPTION_IGNORECASE
		}
		if strings.Index(args[1].(string), "m") >= 0 {
			mode = rubex.ONIG_OPTION_MULTILINE
		}
		var err os.Error
		r, err = rubex.NewRegexp(args[0].(string), mode)
		if err != nil {
			ctx.Log.Error("Invalid regexp")
		}
		ctx.RegexpCache[pattern] = r
		returnValue = r
	}
	return
}

func export_Text(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	val := make([]string, 2)
	val[0] = args[0].(string)
	ts := &Scope{Value: nil}
	ctx.runChildren(ts, ins)
	val[1] = string(ts.Value.([]byte))
	ctx.Exports = append(ctx.Exports, val)
	return
}

func log_Text(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	ctx.Logs = append(ctx.Logs, args[0].(string))
	returnValue = args[0].(string)
	return
}

func concat_Text_Text(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	returnValue = args[0].(string) + args[1].(string)
	return
}

func concat_Text_Text_Text(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	returnValue = args[0].(string) + args[1].(string) + args[2].(string)
	return
}

func downcase_Text(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	returnValue = strings.ToLower(args[0].(string))
	return
}

func upcase_Text(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	returnValue = strings.ToUpper(args[0].(string))
	return
}

func set_Text(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	scope.Value = []byte(args[0].(string))
	return
}

func append_Text(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	scope.Value = scope.Value.(string) + args[0].(string)
	return
}

func prepend_Text(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	scope.Value = args[0].(string) + scope.Value.(string)
	return
}

func index_XMLNode(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	returnValue = fmt.Sprintf("%d", scope.Index+1)
	return
}

func index_Node(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	returnValue = fmt.Sprintf("%d", scope.Index+1)
	return
}

func replace_Text(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	ts := &Scope{Value: ""}
	ctx.runChildren(ts, ins)
	scope.Value = strings.Replace(scope.Value.(string), args[0].(string), ts.Value.(string), -1)
	return
}

func replace_Regexp(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	regexp := args[0].(*rubex.Regexp)
	scope.Value = regexp.GsubFunc(scope.Value.(string), func(match string, captures map[string]string) string {
		usesGlobal := (ctx.Env["use_global_replace_vars"] == "true")

		for name, capture := range captures {
			if usesGlobal {
				//println("setting $", name, "to", capture)
				ctx.Env[name] = capture
			}
			ctx.vars()[name] = capture
		}

		replacementScope := &Scope{Value: match}
		ctx.runChildren(replacementScope, ins)
		//println(ins.String())

		//println("Replacement:", replacementScope.Value.(string))
		return ctx.InnerReplacer.GsubFunc(replacementScope.Value.(string), func(_ string, numeric_captures map[string]string) string {
			capture := numeric_captures["1"]
			var val string
			if usesGlobal {
				val = ctx.Env[capture]
			} else {
				val = ctx.vars()[capture].(string)
			}
			return val
		})
	})
	returnValue = scope.Value
	return
}

func convert_encoding_Text_Text(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	input := scope.Value.(string)
	fromCode := args[0].(string)
	toCode := args[1].(string)
	ic, err := goconv.OpenWithFallback(fromCode, toCode, goconv.KEEP_UNRECOGNIZED)
	if err == nil {
		outputBytes, _ := ic.Conv([]byte(input))
		scope.Value = string(outputBytes)
		ic.Close()
	} else {
		scope.Value = input
	}
	returnValue = scope.Value
	return
}

func xml_Text_Text(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	doc, err := xml.ParseWithBuffer(scope.Value.([]byte), nil, nil, xml.DefaultParseOption, nil, ctx.OutputBuffer)
	if err != nil {
		ctx.Log.Error("xml err: %s", err.String())
		returnValue = "false"
		return
	}
	ns := &Scope{Value: doc}
	ctx.runChildren(ns, ins)
	output := doc.ToXml(nil)
	scope.Value = output
	returnValue = string(output)
	doc.Free()
	return
}

func html_doc_Text_Text(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	inputEncoding := args[0].(string)
	inputEncodingBytes := []byte(inputEncoding)
	outputEncoding := args[1].(string)
	outputEncodingBytes := []byte(outputEncoding)
	input := scope.Value.([]byte)
	doc, err := html.ParseWithBuffer(input, inputEncodingBytes, nil, html.DefaultParseOption, outputEncodingBytes, ctx.OutputBuffer)
	if err != nil {
		ctx.Log.Error("html_doc err: %s", err.String())
		returnValue = "false"
		return
	}
	ns := &Scope{Value: doc}
	ctx.runChildren(ns, ins)
	if err := doc.SetMetaEncoding(outputEncoding); err != nil {
		//ctx.Log.Warn("executing html:" + err.String())
	}
	output := doc.ToHtml(nil)
	scope.Value = output
	returnValue = string(output)
	doc.Free()
	return
}

func html_fragment_Text(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	inputEncoding := args[0].(string)
	inputEncodingBytes := []byte(inputEncoding)
	input := scope.Value.([]byte)
	fragment, err := html.ParseFragment(input, inputEncodingBytes, nil, html.DefaultParseOption, html.DefaultEncodingBytes, ctx.OutputBuffer)
	if err != nil {
		ctx.Log.Error("html_fragment err: %s", err.String())
		returnValue = "false"
		return
	}
	ns := &Scope{Value: fragment}
	ctx.runChildren(ns, ins)
	//output is always utf-8 because the content is internal to Doc.
	scope.Value = ns.Value.(*xml.DocumentFragment).Content()
	returnValue = scope.Value
	fragment.Node.MyDocument().Free()
	return
}

func select_Text(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	node := scope.Value.(xml.Node)

	xpathStr := args[0].(string)
	expr := ctx.XPathCache[xpathStr]
	if expr == nil {
		expr = xpath.Compile(xpathStr)
		if expr == nil {
			ctx.Logs = append(ctx.Logs, "Invalid XPath used: "+xpathStr)
			returnValue = "false"
			return
		}
		ctx.XPathCache[xpathStr] = expr
	}
	nodes, err := node.Search(expr)
	if err != nil {
		ctx.Log.Error("select err: %s", err.String())
		returnValue = "false"
		return
	}

	if len(nodes) == 0 {
		returnValue = "0"
	} else {
		returnValue = fmt.Sprintf("%d", len(nodes))
	}

	for index, node := range nodes {
		if node != nil && node.IsValid() {
			t := node.NodeType()
			if t == xml.XML_ELEMENT_NODE {
				ns := &Scope{Value: node, Index: index}
				ctx.runChildren(ns, ins)
			} else if t == xml.XML_TEXT_NODE {
				ctx.Logs = append(ctx.Logs, "You have just selected a text() node... THIS IS A TERRIBLE IDEA. Please run 'moov check' and sort it out!")
			}
		}
	}
	return
}

func remove_(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	node := scope.Value.(xml.Node)
	node.Remove()
	return
}

func remove_Text(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	node := scope.Value.(xml.Node)

	xpathStr := args[0].(string)
	expr := ctx.XPathCache[xpathStr]
	if expr == nil {
		expr = xpath.Compile(xpathStr)
		if expr == nil {
			ctx.Logs = append(ctx.Logs, "Invalid XPath used: "+xpathStr)
			returnValue = "0"
			return
		}
		ctx.XPathCache[xpathStr] = expr
	}
	nodes, err := node.Search(expr)
	if err != nil {
		ctx.Log.Error("select err: %s", err.String())
		returnValue = "false"
		return
	}

	if len(nodes) == 0 {
		returnValue = "0"
	} else {
		returnValue = fmt.Sprintf("%d", len(nodes))
	}

	for _, node := range nodes {
		if node != nil {
			node.Remove()
		}
	}

	return
}

func position_Text(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	returnValue = Positions[args[0].(string)]
	return
}

func insert_at_Position_Text(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	node := scope.Value.(xml.Node)
	position := args[0].(Position)
	tagName := args[1].(string)
	element := node.MyDocument().CreateElementNode(tagName)
	MoveFunc(element, node, position)
	ns := &Scope{Value: element}
	ctx.runChildren(ns, ins)
	returnValue = "true"
	return
}

func attribute_Text(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	node := scope.Value.(xml.Node)
	name := args[0].(string)
	attr := node.Attribute(name)
	if attr == nil {
		node.SetAttr(name, "")
		attr = node.Attribute(name)
	}
	if attr != nil {
		println("attr", attr.String())
		as := &Scope{Value: attr}
		ctx.runChildren(as, ins)
		if attr.Value() == "" {
			attr.Remove()
		}
		returnValue = "true"
	}
	return
}

func value(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	node := scope.Value.(xml.Node)
	ts := &Scope{Value: node.Content()}
	ctx.runChildren(ts, ins)

	val := ts.Value.([]byte)
	if attr, ok := node.(*xml.AttributeNode); ok {
		attr.SetValue(val)
	}
	returnValue = val
	return
}

func move_XMLNode_XMLNode_Position(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	MoveFunc(args[0].(xml.Node), args[1].(xml.Node), args[2].(Position))
	return
}

func inner(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	node := scope.Value.(xml.Node)
	ts := &Scope{Value: node.Content()}
	ctx.runChildren(ts, ins)
	val := ts.Value.([]byte)
	node.SetInnerHtml(val)
	returnValue = val
	return
}

func equal_XMLNode_XMLNode(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	returnValue = "false"
	node1 := args[0].(xml.Node)
	node2 := args[1].(xml.Node)
	if node1.NodePtr() == node2.NodePtr() {
		returnValue = "true"
	}
	return
}

func move_children_to_XMLNode_Position(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	node := scope.Value.(xml.Node)
	destNode := args[0].(xml.Node)
	if destNode.NodeType() == xml.XML_ELEMENT_NODE {
		child := node.FirstChild()
		for child != nil {
			nextChild := child.NextSibling()
			if child.NodePtr() != destNode.NodePtr() {
				returnValue = "true"
				MoveFunc(child, destNode, args[1].(Position))
			}
			child = nextChild
		}
	}
	return
}

func name(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	node := scope.Value.(xml.Node)
	ts := &Scope{Value: node.Name()}
	ctx.runChildren(ts, ins)
	node.SetName(ts.Value.(string))
	returnValue = ts.Value.(string)
	return
}

func text(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	node := scope.Value.(xml.Node)
	ts := &Scope{Value: node.Content()}
	ctx.runChildren(ts, ins)
	val := ts.Value.(string)
	node.SetContent(val)
	returnValue = val
	return
}

func dup(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	node := scope.Value.(xml.Node)
	newNode := node.Duplicate(1)
	if newNode.NodeType() == xml.XML_ELEMENT_NODE {
		MoveFunc(newNode, node, AFTER)
	}
	ns := &Scope{Value: newNode}
	ctx.runChildren(ns, ins)
	return
}

func fetch_Text(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	node := scope.Value.(xml.Node)
	xpathStr := args[0].(string)
	expr := ctx.XPathCache[xpathStr]
	if expr == nil {
		expr = xpath.Compile(xpathStr)
		if expr == nil {
			ctx.Logs = append(ctx.Logs, "Invalid XPath used: "+xpathStr)
			returnValue = "false"
			return
		}
		ctx.XPathCache[xpathStr] = expr
	}
	nodes, err := node.Search(expr)

	if err == nil && len(nodes) > 0 {
		node := nodes[0]
		if node.NodeType() == xml.XML_ATTRIBUTE_NODE {
			returnValue = node.Content()
		} else {
			returnValue = node.String()
		}
	}
	if len(ins.Children) > 0 {
		ts := &Scope{Value: returnValue}
		ctx.runChildren(ts, ins)
		returnValue = ts.Value
	}
	return
}

func deprecated_Text(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	ctx.Log.Info(args[0].(string))
	return
}

func cdata_Text(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	node := scope.Value.(xml.Node)
	if node.NodeType() == xml.XML_ELEMENT_NODE {
		content := args[0].(string)
		cdata := node.MyDocument().CreateCData(content)
		first := node.FirstChild()
		if first != nil {
			node.ResetChildren()
		}
		node.AddChild(cdata)
	}
	return
}

func inject_at_Position_Text(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	node := scope.Value.(xml.Node)
	position := args[0].(Position)
	input := args[1].([]byte)

	nodes, err := node.Coerce(input)
	if err == nil {
		for _, n := range nodes {
			if position == BEFORE {
				node.InsertBefore(n)
			} else if position == AFTER {
				node.InsertAfter(n)
			} else if position == TOP {
				node.InsertBegin(n)
			} else if position == BOTTOM {
				node.InsertEnd(n)
			}
		}
	}
	if len(nodes) > 0 {
		first := nodes[0]
		if first.NodeType() == xml.XML_ELEMENT_NODE {
			// successfully ran scope
			returnValue = "true"
			ns := &Scope{Value: first}
			ctx.runChildren(ns, ins)
		}
	} else {
		returnValue = "false"
	}
	return
}

func path(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	returnValue = scope.Value.(xml.Node).Path()
	return
}

func css_Text(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	returnValue = css2xpath.Convert(args[0].(string), css2xpath.LOCAL)
	return
}

func wrap_text_children_Text(ctx *Ctx, scope *Scope, ins *tp.Instruction, args []interface{}) (returnValue interface{}) {
	returnValue = "false"
	node := scope.Value.(xml.Node)
	if textNodes, err := node.Search("./text()"); err == nil {
		tagName := args[0].(string)
		tag := fmt.Sprintf("<%s />", tagName)
		for index, textNode := range textNodes {
			textNode.Wrap(tag)
			ns := &Scope{textNode, index}
			ctx.runChildren(ns, ins)
		}
	}
	return
}

/*

	// ATTRIBUTE FUNCTIONS
	case "attribute.Text":
		node := scope.Value.(xml.Node)
		name := args[0].(string)
		if _, ok := node.(*xml.Element); ok {
			attr, _ := node.Attribute(name)
			as := &Scope{Value: attr}
			ctx.runChildren(as, ins)
			if attr.IsLinked() && (attr.Content() == "") {
				attr.Remove()
			}
			if !attr.IsLinked() {
				attr.Free()
			}
			returnValue = "true"
		}
	case "dump":
		returnValue = scope.Value.(xml.Node).String()

	default:
		ctx.Log.Error("Must implement " + fun.Name)
	}
*/
