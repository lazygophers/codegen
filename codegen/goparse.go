package codegen

import (
	"fmt"
	"github.com/lazygophers/utils/candy"
	"github.com/lazygophers/utils/osx"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/gookit/color"
	"github.com/lazygophers/log"
)

type GoContext struct {
	GoPathList          []string
	ModPath2RealPathMap map[string]string
	ModPath2ModNameMap  map[string]string
}
type GoFuncNode struct {
	goFunc *GoFunc
	goFile *GoFile
}

func (g *GoFuncNode) GetFunc() *GoFunc {
	return g.goFunc
}
func (g *GoFuncNode) GetFile() *GoFile {
	return g.goFile
}

type GoPackage struct {
	Name    string
	FileMap map[string]*GoFile
	FuncMap map[string]*GoFuncNode
}

type GoImport struct {
	Name string
	Path string
}

type GoReturn struct {
	Name string
}

func (p *GoReturn) String() string {
	return p.Name
}

func (p *GoReturn) StringWithColor() string {
	return color.Blue.Text(p.Name)
}

type GoParams struct {
	TypeName  string
	FieldName string

	// typename == func
	GoFunc *GoFunc
}

func (p *GoParams) String() string {
	switch p.TypeName {
	case "func":
		return fmt.Sprintf("%s", p.GoFunc.String())
	default:
		return fmt.Sprintf("%s %s", p.FieldName, p.TypeName)
	}
}

func (p *GoParams) StringWithColor() string {
	switch p.TypeName {
	case "func":
		return fmt.Sprintf("%s", p.GoFunc.StringWithColor())
	default:
		return fmt.Sprintf("%s %s", color.Cyan.Text(p.FieldName), color.Blue.Text(p.TypeName))
	}
}

type GoFuncCallArg struct {
	Names []string
}

type GoFuncCall struct {
	Names []string
	Args  []*GoFuncCallArg
}

type GoFunc struct {
	Name      string
	FuncCalls []*GoFuncCall

	Returns []*GoReturn
	Params  []*GoParams

	// 接收类型,struct的类型
	RecvType string
}

func (p *GoFunc) Clone() *GoFunc {
	return &GoFunc{
		Name:      p.Name,
		FuncCalls: p.FuncCalls,
		Returns:   p.Returns,
		Params:    p.Params,
		RecvType:  p.RecvType,
	}
}

func (p *GoFunc) String() string {
	return fmt.Sprintf("%s%s (%s) %s",
		func() string {
			if p.RecvType == "" {
				return ""
			}
			return fmt.Sprintf("(%s) ", p.RecvType)
		}(),
		p.Name,
		func() string {
			var typs []string

			for _, val := range p.Params {
				typs = append(typs, val.String())
			}

			return strings.Join(typs, ", ")
		}(),
		func() string {
			var typs []string

			for _, val := range p.Returns {
				typs = append(typs, val.String())
			}

			return strings.Join(typs, ", ")
		}(),
	)
}

func (p *GoFunc) StringWithColor() string {
	return fmt.Sprintf("%s%s (%s) [%s]",
		func() string {
			if p.RecvType == "" {
				return ""
			}
			return color.Blue.Sprintf("(%s) ", p.RecvType)
		}(),
		color.Yellow.Text(p.Name),
		func() string {
			var typs []string

			for _, val := range p.Params {
				typs = append(typs, val.StringWithColor())
			}

			return strings.Join(typs, ", ")
		}(),
		func() string {
			var typs []string

			for _, val := range p.Returns {
				typs = append(typs, val.StringWithColor())
			}

			return strings.Join(typs, ", ")
		}(),
	)
}

type GoStruct struct {
	Name string
}

type GoFile struct {
	Imports []*GoImport
	Funcs   []*GoFunc
	Structs []*GoStruct
}

func (p *GoPackage) buildFuncMap() {
	p.FuncMap = map[string]*GoFuncNode{}
	for _, v := range p.FileMap {
		for _, x := range v.Funcs {
			p.FuncMap[x.Name] = &GoFuncNode{
				goFunc: x,
				goFile: v,
			}
		}
	}
}
func (p *GoPackage) GetFunc(name string) *GoFuncNode {
	return p.FuncMap[name]
}
func (ctx *GoContext) initGoPath() {
	if len(ctx.GoPathList) > 0 {
		return
	}
	v := os.Getenv("GOPATH")
	if v != "" {
		var vv []string
		if runtime.GOOS == "windows" {
			vv = strings.Split(v, ";")
		} else {
			vv = strings.Split(v, ":")
		}
		for _, x := range vv {
			if x != "" {
				ctx.GoPathList = append(ctx.GoPathList, x)
			}
		}
	}
}
func (ctx *GoContext) FindPackagePath(modPath string) string {
	ctx.initGoPath()
	if ctx.ModPath2RealPathMap == nil {
		ctx.ModPath2RealPathMap = map[string]string{}
	}
	if v, ok := ctx.ModPath2RealPathMap[modPath]; ok {
		return v
	}
	for _, v := range ctx.GoPathList {
		r := fmt.Sprintf("%s/src/%s", v, modPath)
		if osx.IsFile(r) {
			ctx.ModPath2RealPathMap[modPath] = r
			return r
		}
		r = fmt.Sprintf("%s/%s", v, modPath)
		if osx.IsFile(r) {
			ctx.ModPath2RealPathMap[modPath] = r
			return r
		}
	}
	ctx.ModPath2RealPathMap[modPath] = ""
	return ""
}
func (ctx *GoContext) FindPackageName(modPath string) string {
	if ctx.ModPath2ModNameMap == nil {
		ctx.ModPath2ModNameMap = map[string]string{}
	}
	if v, ok := ctx.ModPath2ModNameMap[modPath]; ok {
		return v
	}
	realPath := ctx.FindPackagePath(modPath)
	if realPath == "" {
		ctx.ModPath2ModNameMap[modPath] = ""
		return ""
	}
	files, err := ioutil.ReadDir(realPath)
	if err != nil {
		log.Fatalf("read dir %s err %v", realPath, err)
	}
	var name string
	for _, f := range files {
		if !strings.HasSuffix(strings.ToLower(f.Name()), ".go") {
			continue
		}
		filePath := fmt.Sprintf("%s/%s", realPath, f.Name())
		buf, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatalf("read file %s err %v", filePath, err)
		}
		const (
			pkg = "package"
		)
		var p int
		for {
			tmp := strings.Index(string(buf[p:]), pkg)
			if tmp < 0 {
				p = -1
				break
			}
			p += tmp
			// 看看是不是 package 是一个新行，排除一下 package 在注释里面了
			ok := true
			if p > 0 {
				c := buf[p-1]
				if !(c == '\r' || c == '\n') {
					ok = false
				}
			}
			p += len(pkg)
			if ok {
				break
			}
		}
		if p < 0 {
			continue
		}
		for p < len(buf) {
			if buf[p] == ' ' || buf[p] == '\t' {
				p++
			} else {
				break
			}
		}
		b := p
		for p < len(buf) {
			if buf[p] == '\r' || buf[p] == '\n' {
				break
			} else {
				p++
			}
		}
		if b < p {
			tmp := string(buf[b:p])
			if strings.HasSuffix(tmp, "_test") {
				continue
			}
			name = tmp
			break
		}
	}
	ctx.ModPath2ModNameMap[modPath] = name
	return name
}
func removeStrQuote(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}
func fetchTypeName(ctx *GoContext, expr ast.Expr) string {
	var names []string
	isPtr := false
C:
	switch x := expr.(type) {
	case *ast.Ident:
		names = append(names, x.Name)
	case *ast.StarExpr:
		isPtr = true
		expr = x.X
		goto C
	case *ast.SelectorExpr:
		if x.Sel != nil {
			names = append(names, x.Sel.Name)
		}
		expr = x.X
		goto C
	}
	candy.Reverse(names)
	n := strings.Join(names, ".")
	if isPtr {
		n = "*" + n
	}
	return n
}
func fetchImports(ctx *GoContext, f *ast.File, goFile *GoFile) {
	for _, i := range f.Imports {
		goImport := &GoImport{}
		if i.Name != nil {
			goImport.Name = i.Name.Name
		}
		if i.Path != nil {
			goImport.Path = removeStrQuote(i.Path.Value)
		}
		if goImport.Path == "" {
			continue
		}
		if goImport.Name == "" {
			name := ctx.FindPackageName(goImport.Path)
			if name != "" {
				goImport.Name = name
			} else {
				goImport.Name = path.Base(goImport.Path)
			}
		}
		goFile.Imports = append(goFile.Imports, goImport)
	}
}
func getSelectorExprIndents(ctx *GoContext, c ast.Expr, goFunc *GoFunc) []string {
	var names []string
	for {
		cont := false
		switch x := c.(type) {
		case *ast.Ident:
			names = append(names, x.Name)
		case *ast.SelectorExpr:
			if x.Sel != nil {
				names = append(names, x.Sel.Name)
			}
			c = x.X
			cont = true
		case *ast.CallExpr:
			fetchCallExpr(ctx, x, goFunc)
		case *ast.BasicLit:
		// log.Infof("%+v", x)// 调用函数的传入参数
		case *ast.FuncLit:
			if x.Body != nil {
				fetchStmtList(ctx, x.Body.List, goFunc)
			}
		case *ast.UnaryExpr:
		// log.Infof("%+v", x)// 各种变量
		case *ast.CompositeLit:
		// for _, elt := range x.Elts {
		//	log.Infof("%+v", elt)
		// }
		case *ast.IndexExpr:
		// log.Infof("%+v", x.X)
		case *ast.BinaryExpr:
		// log.Infof("%+v", x.X)
		default:
			// log.Infof("%+v", reflect.TypeOf(c))
		}
		if !cont {
			break
		}
	}
	candy.Reverse(names)
	return names
}
func fetchCallExprFunc(ctx *GoContext, c ast.Expr, goFunc *GoFunc) *GoFuncCall {
	var names []string
	goCall := &GoFuncCall{}
	switch x := c.(type) {
	case *ast.Ident:
		names = append(names, x.Name)
	case *ast.SelectorExpr:
		names = getSelectorExprIndents(ctx, c, goFunc)
	case *ast.FuncLit:
		if x.Body != nil {
			fetchStmtList(ctx, x.Body.List, goFunc)
		}
	}
	if len(names) == 0 {
		return nil
	}
	goCall.Names = names
	return goCall
}
func fetchExpr(ctx *GoContext, c ast.Expr, goFunc *GoFunc) {
	// 表达式这里，用到的不多，后续需要再补充
	switch x := c.(type) {
	case *ast.CallExpr:
		fetchCallExpr(ctx, x, goFunc)
	case *ast.FuncLit:
		if x.Body != nil {
			fetchStmtList(ctx, x.Body.List, goFunc)
		}
	case *ast.BinaryExpr:
		fetchExpr(ctx, x.X, goFunc)
		fetchExpr(ctx, x.Y, goFunc)
	case *ast.UnaryExpr:
		fetchExpr(ctx, x.X, goFunc)
	case *ast.ParenExpr:
		fetchExpr(ctx, x.X, goFunc)
	case *ast.StarExpr:
		fetchExpr(ctx, x.X, goFunc)
	}
}
func fetchCallExpr(ctx *GoContext, c *ast.CallExpr, goFunc *GoFunc) {
	// 函数
	goCall := fetchCallExprFunc(ctx, c.Fun, goFunc)
	// 参数
	for _, v := range c.Args {
		names := getSelectorExprIndents(ctx, v, goFunc)
		if goCall != nil {
			goCall.Args = append(goCall.Args, &GoFuncCallArg{Names: names})
		}
		if len(names) == 0 {
			fetchExpr(ctx, v, goFunc)
		}
	}
	// TODO 返回值
	if goCall != nil {
		goFunc.FuncCalls = append(goFunc.FuncCalls, goCall)
	}
}
func fetchReturnStmt(ctx *GoContext, c *ast.ReturnStmt, goFunc *GoFunc) {
	for _, v := range c.Results {
		fetchExpr(ctx, v, goFunc)
	}
}
func fetchStmtList(ctx *GoContext, stmtList []ast.Stmt, goFunc *GoFunc) {
	for _, st := range stmtList {
		switch t := st.(type) {
		case *ast.ExprStmt:
			switch x := t.X.(type) {
			case *ast.CallExpr:
				fetchCallExpr(ctx, x, goFunc)
			}
		case *ast.IfStmt:
			if t.Body != nil {
				fetchStmtList(ctx, t.Body.List, goFunc)
			}
			switch y := t.Else.(type) {
			case *ast.BlockStmt:
				fetchStmtList(ctx, y.List, goFunc)
			}
		case *ast.ForStmt:
			if t.Body != nil {
				fetchStmtList(ctx, t.Body.List, goFunc)
			}
		case *ast.BlockStmt:
			fetchStmtList(ctx, t.List, goFunc)
		case *ast.GoStmt:
			fetchCallExpr(ctx, t.Call, goFunc)
		case *ast.ReturnStmt:
			fetchReturnStmt(ctx, t, goFunc)
		case *ast.CaseClause:
			fetchStmtList(ctx, t.Body, goFunc)
		case *ast.SwitchStmt:
			if t.Body != nil {
				fetchStmtList(ctx, t.Body.List, goFunc)
			}
		case *ast.SelectStmt:
			if t.Body != nil {
				fetchStmtList(ctx, t.Body.List, goFunc)
			}
		case *ast.CommClause:
			fetchStmtList(ctx, t.Body, goFunc)
		case *ast.RangeStmt:
			if t.Body != nil {
				fetchStmtList(ctx, t.Body.List, goFunc)
			}
		case *ast.AssignStmt:
			for _, e := range t.Lhs {
				fetchExpr(ctx, e, goFunc)
			}
			for _, e := range t.Rhs {
				fetchExpr(ctx, e, goFunc)
			}
		default:
			// log.Infof("%+v", reflect.TypeOf(st))
		}
	}
}

func parseFuncReturns(ctx *GoContext, list []*ast.Field) []*GoReturn {
	var returns []*GoReturn
	for _, field := range list {
		returns = append(returns, &GoReturn{
			Name: fetchTypeName(ctx, field.Type),
		})
	}
	return returns
}

func parseFuncParams(ctx *GoContext, list []*ast.Field) []*GoParams {
	var params []*GoParams

	for _, field := range list {
		switch x := field.Type.(type) {
		case *ast.StarExpr, *ast.Ident:
			for _, name := range field.Names {
				params = append(params, &GoParams{
					TypeName:  fetchTypeName(ctx, field.Type),
					FieldName: name.Name,
				})
			}
		case *ast.FuncType:
			goFunc := &GoFunc{
				Returns: nil,
				Params:  nil,
			}

			parseFuncType(ctx, goFunc, x)

			for _, name := range field.Names {
				_gf := goFunc.Clone()
				_gf.Name = name.Name

				param := &GoParams{
					TypeName:  "func",
					FieldName: name.Name,
					GoFunc:    _gf,
				}
				params = append(params, param)
			}
			// default:
			//	log.Info(fetchTypeName(ctx, field.Type))
			//	log.Info(field.Type)
			//	log.Info(field.Tag)
			//	log.Info(field.Names)
			//	log.Info(field.Comment)
			//	log.Info(field.Doc)
			//	log.Info(reflect.TypeOf(x))
		}
	}

	return params
}

func parseFuncType(ctx *GoContext, goFunc *GoFunc, f *ast.FuncType) {
	if f.Results != nil {
		goFunc.Returns = append(goFunc.Returns, parseFuncReturns(ctx, f.Results.List)...)
	}

	if f.Params != nil {
		goFunc.Params = append(goFunc.Params, parseFuncParams(ctx, f.Params.List)...)
	}
}

func parseFuncDecl(ctx *GoContext, f *ast.FuncDecl) *GoFunc {
	goFunc := &GoFunc{}
	if f.Name != nil {
		goFunc.Name = f.Name.Name
		if f.Recv != nil {
			for _, r := range f.Recv.List {
				rt := fetchTypeName(ctx, r.Type)
				if rt != "" {
					goFunc.RecvType = rt
					break
				}
			}
		}

		if f.Type != nil {
			parseFuncType(ctx, goFunc, f.Type)
		}
	}
	if f.Body != nil {
		fetchStmtList(ctx, f.Body.List, goFunc)
	}
	return goFunc
}

func fetchFunc(ctx *GoContext, f *ast.FuncDecl, goFile *GoFile) {
	log.Info(f.Doc)
	goFunc := parseFuncDecl(ctx, f)
	if goFunc == nil {
		return
	}

	goFile.Funcs = append(goFile.Funcs, goFunc)
}
func fetchType(ctx *GoContext, ts *ast.TypeSpec, goFile *GoFile) {
	if ts.Type == nil {
		return
	}
	if ts.Name == nil {
		return
	}
	name := ts.Name.Name
	if name == "" {
		return
	}
	switch ts.Type.(type) {
	case *ast.StructType:
		goFile.Structs = append(goFile.Structs, &GoStruct{
			Name: name,
		})
	}
}
func ParseGoFile(filePath string) (*GoFile, error) {
	fSet := token.NewFileSet()
	f, err := parser.ParseFile(fSet, filePath, nil, 0)
	if err != nil {
		log.Errorf("parse go file %s err %v", filePath, err)
		return nil, err
	}
	var ctx GoContext
	goFile := &GoFile{}
	fetchImports(&ctx, f, goFile)
	for _, decl := range f.Decls {
		switch t := decl.(type) {
		case *ast.FuncDecl:
			fetchFunc(&ctx, t, goFile)
		case *ast.GenDecl:
			switch t.Tok {
			case token.TYPE:
				// 类型定义
				for _, spec := range t.Specs {
					switch x := spec.(type) {
					case *ast.TypeSpec:
						fetchType(&ctx, x, goFile)
					}
				}
			}
		}
	}
	return goFile, nil
}
func ParseGoDir(dir string) (map[string]*GoPackage, error) {
	fSet := token.NewFileSet()
	m, err := parser.ParseDir(fSet, dir, nil, 0)
	if err != nil {
		log.Errorf("parse dir %s err %v", dir, err)
		return nil, err
	}
	var ctx GoContext
	var pkgMap = map[string]*GoPackage{}
	for k, v := range m {
		pkg := &GoPackage{Name: k, FileMap: map[string]*GoFile{}}
		pkgMap[k] = pkg
		for fp, f := range v.Files {
			goFile := &GoFile{}
			fetchImports(&ctx, f, goFile)
			for _, decl := range f.Decls {
				switch t := decl.(type) {
				case *ast.FuncDecl: // func
					fetchFunc(&ctx, t, goFile)
					// case *ast.GenDecl:// type
					//	for _, spec := range t.Specs {
					//		log.Infof("%+v", spec)
					//	}
					// default:
					//	log.Info(reflect.TypeOf(decl))
				}
			}
			pkg.FileMap[fp] = goFile
		}
		pkg.buildFuncMap()
	}
	return pkgMap, nil
}

type ScanErrCodeContext struct {
	errMap    map[string]bool
	accessMap map[string]bool
	pkgMap    map[string]*GoPackage
	goCtx     GoContext
}

func NewScanErrCodeContext() *ScanErrCodeContext {
	return &ScanErrCodeContext{
		errMap:    map[string]bool{},
		accessMap: map[string]bool{},
		pkgMap:    map[string]*GoPackage{},
	}
}
func (p *ScanErrCodeContext) GetModList() []string {
	var list []string
	for k, v := range p.pkgMap {
		if k != "" && v != nil {
			list = append(list, k)
		}
	}
	return list
}
func (p *ScanErrCodeContext) GetErrListAndClear() []string {
	var errList []string
	for k := range p.errMap {
		errList = append(errList, k)
	}
	p.errMap = map[string]bool{}
	p.accessMap = map[string]bool{}
	return errList
}
func (p *ScanErrCodeContext) MergePkgMap(m map[string]*GoPackage) {
	for k, v := range m {
		p.pkgMap[k] = v
	}
}
func (p *ScanErrCodeContext) MergePkgMapWithPkgName(m map[string]*GoPackage, pkg string) {
	for _, v := range m {
		p.pkgMap[pkg] = v
	}
}
func (p *ScanErrCodeContext) Scan(mod, name string) {
	var pkg *GoPackage
	pkg = p.pkgMap[mod]
	if pkg == nil {
		log.Warnf("not found package %s", mod)
		return
	}
	// 看看是否处理过了
	full := fmt.Sprintf("%s.%s", mod, name)
	if p.accessMap[full] {
		return
	}
	p.accessMap[full] = true
	// find Pkg
	node := pkg.GetFunc(name)
	if node == nil {
		return
	}
	f := node.goFunc
	for _, v := range f.FuncCalls {
		if len(v.Names) == 2 &&
			v.Names[0] == "rpc" && strings.HasPrefix(v.Names[1], "CreateError") {
			if len(v.Args) > 0 {
				x := v.Args[0]
				if len(x.Names) == 2 && strings.HasPrefix(x.Names[1], "Err") {
					p.errMap[fmt.Sprintf("%s.%s", x.Names[0], x.Names[1])] = true
				}
			}
			continue
		}
		//
		// 支持这种形式:
		// XXXErr(xxx.ErrXXX)
		//
		// if len(v.Names) == 1 &&
		//	strings.HasSuffix(v.Names[0], "Err") {
		//	if len(v.Args) > 0 {
		//		x := v.Args[0]
		//		if len(x.Names) == 2 && strings.HasPrefix(x.Names[1], "Err") {
		//			p.errMap[fmt.Sprintf("%s.%s", x.Names[0], x.Names[1])] = true
		//		}
		//	}
		//	continue
		// }
		if len(v.Names) == 1 {
			p.Scan(mod, v.Names[0])
		} else if len(v.Names) == 2 {
			callMod := v.Names[0]
			callName := v.Names[1]
			// 1. 已经 load 过了
			if callPkg, ok := p.pkgMap[callMod]; ok {
				if callPkg != nil {
					p.Scan(callMod, callName)
				}
				continue
			}
			// 找下这个模块的 path
			goFile := node.goFile
			var impPath string
			for _, x := range goFile.Imports {
				if x.Name == callMod {
					impPath = x.Path
					break
				}
			}
			if impPath == "" {
				continue
			}
			fullPath := p.goCtx.FindPackagePath(impPath)
			if fullPath == "" {
				continue
			}
			files, err := ioutil.ReadDir(fullPath)
			if err != nil {
				log.Fatalf("read dir %s err:%v", fullPath, err)
				return
			}
			hasClient := false
			for _, f := range files {
				if strings.HasSuffix(f.Name(), "client.go") {
					hasClient = true
					break
				}
			}
			if !hasClient {
				p.pkgMap[callMod] = nil
				continue
			}
			// 看下下面有没有 impl
			fullPath += "/impl"
			if !osx.IsFile(fullPath) {
				p.pkgMap[callMod] = nil
				continue
			}
			m, err := ParseGoDir(fullPath)
			if err != nil {
				log.Errorf("parse go dir %s err %v", fullPath, err)
				p.pkgMap[callMod] = nil
				continue
			}
			p.MergePkgMapWithPkgName(m, callMod)
			// log.Infof("load dep module %s", callMod)
			p.Scan(callMod, callName)
		}
	}
}
