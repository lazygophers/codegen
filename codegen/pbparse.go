package codegen

import (
	"bytes"
	"github.com/emicklei/proto"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc/middleware/core"
	"github.com/lazygophers/utils/candy"
	"github.com/lazygophers/utils/json"
	"github.com/lazygophers/utils/stringx"
	"github.com/pterm/pterm"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type PbOption struct {
	option *proto.Option

	Name  string
	Value string
}

func (p *PbOption) walk() {
	p.Name = p.option.Name
	p.Value = p.option.Constant.Source
}

func NewPbOption(o *proto.Option) *PbOption {
	p := &PbOption{
		option: o,
	}
	p.walk()
	return p
}

type PbCommentTag struct {
	tag string

	lines []string
}

func (p *PbCommentTag) append(line string) {
	p.lines = append(p.lines, line)
}

func (p *PbCommentTag) Tag() string {
	return p.tag
}

func (p *PbCommentTag) Lines() []string {
	return p.lines
}

func (p *PbCommentTag) walk() {

}

func NewPbCommentTag(tag string) *PbCommentTag {
	return &PbCommentTag{
		tag: tag,
	}
}

type PbComment struct {
	comment *proto.Comment

	// 按照 tag 格式的注释
	tags map[string]*PbCommentTag
}

func (p *PbComment) Comment() *proto.Comment {
	return p.comment
}

func (p *PbComment) walk() {
	var lastTag string
	for _, line := range p.comment.Lines {
		for strings.HasPrefix(line, " ") {
			line = strings.TrimLeft(line, " ")
		}

		tag := lastTag
		if strings.HasPrefix(line, "@") {
			idx := strings.Index(line, ":")
			if idx > 0 {
				tag = line[1:idx]
				line = line[idx+1:]
			}
		}

		if tag != "" {
			lastTag = tag

			if _, ok := p.tags[tag]; !ok {
				p.tags[tag] = NewPbCommentTag(tag)
			}

			p.tags[tag].append(line)
		}
	}

	for _, tag := range p.tags {
		tag.walk()
	}
}

func NewPbComment(c *proto.Comment) *PbComment {
	p := &PbComment{
		comment: c,
		tags:    map[string]*PbCommentTag{},
	}

	p.walk()

	return p
}

type PbEnum struct {
	enum *proto.Enum

	Name     string
	FullName string

	fields []*PbEnumField
}

func (p *PbEnum) Enum() *proto.Enum {
	return p.enum
}

func (p *PbEnum) walk() {
	p.Name = p.enum.Name
	p.FullName = GetFullName(p.enum)

	for _, element := range p.enum.Elements {
		switch x := element.(type) {
		case *proto.EnumField:
			p.fields = append(p.fields, NewPbEnumField(x))

		case *proto.Comment:
			continue

		default:
			log.Panicf("unexpected enum field: %T", x)
		}
	}
}

func NewPbEnum(e *proto.Enum) *PbEnum {
	p := &PbEnum{
		enum: e,
	}

	p.walk()

	return p
}

type PbRpcGenOptions struct {
	GenTo        string
	Model        string
	Action       string
	Role         string
	SkipGenRoute bool
	Method       string
	Path         string
}

type PbRPC struct {
	rpc  *proto.RPC
	Name string

	options         map[string]map[string]string
	comment         *PbComment
	genOption       *PbRpcGenOptions
	requestType     string
	returnsType     string
	requestPackage  string
	responsePackage string
}

func (p *PbRPC) setDefaultPackage(pkg string) {
	if p.requestPackage == "" {
		p.requestPackage = pkg
	}

	if p.responsePackage == "" {
		p.responsePackage = pkg
	}
}

func (p *PbRPC) RPC() *proto.RPC {
	return p.rpc
}

func (p *PbRPC) RequestType() string {
	return p.requestType
}

func (p *PbRPC) ReturnsType() string {
	return p.returnsType
}

func (p *PbRPC) RequestPackage() string {
	return p.requestPackage
}

func (p *PbRPC) ResponsePackage() string {
	return p.responsePackage
}

func (p *PbRPC) walk() {
	for _, option := range p.rpc.Options {
		p.options[option.Name] = make(map[string]string, len(option.AggregatedConstants))
		for _, value := range option.AggregatedConstants {
			p.options[option.Name][value.Name] = value.Literal.Source
		}
	}

	if p.rpc.Comment != nil {
		p.comment = NewPbComment(p.rpc.Comment)
	}

	// 优先解析 options 的
	if v, ok := p.options["(lazygophers.lrpc.core.lazygen)"]; ok {
		var gen core.LazyGen
		buffer, err := json.Marshal(v)
		if err != nil {
			log.Panicf("err:%v", err)
			return
		}

		err = json.Unmarshal(buffer, &gen)
		if err != nil {
			log.Panicf("err:%v", err)
			return
		}

		log.Info(v)
		log.Info(gen)

		if gen.Role != "" {
			p.genOption.Role = gen.Role
		}

		if gen.SkipGenRoute {
			p.genOption.SkipGenRoute = gen.SkipGenRoute
		}
	} else {
		log.Info(p.options)
	}

	if v, ok := p.options["(lazygophers.lrpc.core.http)"]; ok {
		var gen core.Http
		buffer, err := json.Marshal(v)
		if err != nil {
			log.Panicf("err:%v", err)
			return
		}

		err = json.Unmarshal(buffer, &gen)
		if err != nil {
			log.Panicf("err:%v", err)
			return
		}

		if gen.Method != nil {
			p.genOption.Method = *gen.Method
		}

		if gen.Path != nil {
			p.genOption.Path = *gen.Path
		}
	}

	// 然后解析 注释的
	if p.comment != nil {
		if v, ok := p.comment.tags["gen"]; ok {
			p.genOption.GenTo = strings.ReplaceAll(strings.Join(v.Lines(), ""), " ", "")
		}

		if v, ok := p.comment.tags["model"]; ok {
			p.genOption.Model = strings.ReplaceAll(strings.Join(v.Lines(), ""), " ", "")
		}

		if v, ok := p.comment.tags["action"]; ok {
			p.genOption.Action = strings.ReplaceAll(strings.Join(v.Lines(), ""), " ", "")
		}

		if v, ok := p.comment.tags["role"]; ok {
			p.genOption.Role = strings.ReplaceAll(strings.Join(v.Lines(), ""), " ", "")
		}
	}

	if p.genOption.GenTo == "" {
		if p.genOption.Model != "" {
			p.genOption.GenTo = stringx.ToSnake(strings.TrimPrefix(p.genOption.Model, "Model"))
		} else {
			p.genOption.GenTo = "impl"
		}
	}

	if strings.Contains(p.rpc.RequestType, ".") {
		text := p.rpc.RequestType
		if idx := strings.LastIndex(text, "."); idx > 0 {
			p.requestType = text[idx+1:]
			text = text[:idx]
		}
		if idx := strings.LastIndex(text, "."); idx > 0 {
			p.requestPackage = text[idx+1:]
		} else {
			p.requestPackage = text
		}
	}

	if strings.Contains(p.rpc.ReturnsType, ".") {
		text := p.rpc.ReturnsType
		if idx := strings.LastIndex(text, "."); idx > 0 {
			p.requestType = text[idx+1:]
			text = text[:idx]
		}
		if idx := strings.LastIndex(text, "."); idx > 0 {
			p.responsePackage = text[idx+1:]
		} else {
			p.responsePackage = text
		}
	}

}

func NewPbRPC(rpc *proto.RPC) *PbRPC {
	p := &PbRPC{
		Name:    rpc.Name,
		rpc:     rpc,
		options: make(map[string]map[string]string, len(rpc.Options)),
		genOption: &PbRpcGenOptions{
			Method: "POST",
			Path:   "/" + rpc.Name,
		},
		requestType: rpc.RequestType,
		returnsType: rpc.ReturnsType,
	}

	p.walk()
	return p
}

type PbService struct {
	service *proto.Service
}

func (p *PbService) Service() *proto.Service {
	return p.service
}

func (p *PbService) walk() {
}

func NewPbService(service *proto.Service) *PbService {
	p := &PbService{
		service: service,
	}

	p.walk()

	return p
}

type PbField interface {
	FieldName() string
	FieldType() string
	IsSlice() bool
}

var _ PbField = (*PbNormalField)(nil)
var _ PbField = (*PbMapField)(nil)
var _ PbField = (*PbEnumField)(nil)

type PbNormalField struct {
	Name  string
	field *proto.NormalField

	// 写在字段上的注释
	comment *PbComment

	// 写在字段后面的注释
	inlineComment *PbComment
}

func (p *PbNormalField) FieldName() string {
	return p.Name
}

func (p *PbNormalField) FieldType() string {
	return p.Type()
}

func (p *PbNormalField) IsSlice() bool {
	return p.field.Repeated
}

func (p *PbNormalField) Field() *proto.NormalField {
	return p.field
}

func (p *PbNormalField) Type() string {
	idx := strings.LastIndex(p.field.Type, ".")
	if idx != -1 {
		return p.field.Type[idx+1:]
	} else {
		return p.field.Type
	}
}

func (p *PbNormalField) FullType() string {
	return p.field.Type
}

func (p *PbNormalField) walk() {
	p.Name = p.field.Name

	if p.field.Comment != nil {
		p.comment = NewPbComment(p.field.Comment)
	}

	if p.field.InlineComment != nil {
		p.inlineComment = NewPbComment(p.field.InlineComment)
	}
}

func NewPbNormalField(f *proto.NormalField) *PbNormalField {
	p := &PbNormalField{
		field: f,
	}

	p.walk()

	return p
}

type PbMapField struct {
	Name  string
	field *proto.MapField
}

func (p *PbMapField) FieldName() string {
	return p.Name
}

func (p *PbMapField) FieldType() string {
	return p.field.Type
}

func (p *PbMapField) IsSlice() bool {
	return false
}

func (p *PbMapField) walk() {

}

func NewPbMapField(f *proto.MapField) *PbMapField {
	p := &PbMapField{
		Name:  f.Name,
		field: f,
	}

	p.walk()

	return p
}

type PbEnumField struct {
	Name     string
	field    *proto.EnumField
	comment  *PbComment
	FullName string
}

func (p *PbEnumField) FieldName() string {
	return p.Name
}

func (p *PbEnumField) FieldType() string {
	return "int32"
}

func (p *PbEnumField) IsSlice() bool {
	return false
}

func (p *PbEnumField) walk() {
	p.Name = p.field.Name
	p.FullName = GetFullName(p.field)

	if p.field.Comment != nil {
		p.comment = NewPbComment(p.field.Comment)
	}
}

func NewPbEnumField(f *proto.EnumField) *PbEnumField {
	p := &PbEnumField{
		field: f,
	}

	p.walk()

	return p
}

type PbMessage struct {
	message *proto.Message

	normalFieldMap map[string]*PbNormalField
	mapFieldMap    map[string]*PbMapField
	enumFieldMap   map[string]*PbEnumField

	normalFieldList []*PbNormalField
	mapFieldList    []*PbMapField
	enumFieldList   []*PbEnumField

	FullName string
	Name     string
}

func (p *PbMessage) PrimaryField() (pkField *PbNormalField) {
	for _, field := range p.normalFieldMap {
		// 先简单粗暴用 id 当作主键，后面再改
		if field.Name == "id" {
			pkField = field
			break
		}
	}

	return pkField
}

func (p *PbMessage) NormalFields() []*PbNormalField {
	return p.normalFieldList
}

func (p *PbMessage) MapFields() []*PbMapField {
	return p.mapFieldList
}

func (p *PbMessage) EnumFields() []*PbEnumField {
	return p.enumFieldList
}

func (p *PbMessage) Message() *proto.Message {
	return p.message
}

func (p *PbMessage) IsTable() bool {
	if strings.HasPrefix(p.FullName, "Model") {
		if !strings.Contains(p.FullName, "_") {
			return true
		}
	}

	// TODO: 允许通过对注释的解析判断时候是表

	return false
}

func (p *PbMessage) NeedOrm() bool {
	if strings.HasPrefix(p.FullName, "Model") {
		return true
	}

	// TODO: 允许通过对注释的解析判断时候是表

	return false
}

func (p *PbMessage) walk() {
	p.FullName = GetFullName(p.message)
	p.Name = p.message.Name

	for _, element := range p.message.Elements {
		var visitor = &ProtoVisitor{}

		element.Accept(visitor)

		for _, field := range visitor.mapFields {
			pterm.Info.Printfln("find map field:%s in %s", field.Name, p.FullName)
			p.mapFieldMap[field.Name] = NewPbMapField(field)
			p.mapFieldList = append(p.mapFieldList, p.mapFieldMap[field.Name])
		}

		for _, field := range visitor.normalFields {
			pterm.Info.Printfln("find normal field:%s in %s", field.Name, p.FullName)
			p.normalFieldMap[field.Name] = NewPbNormalField(field)
			p.normalFieldList = append(p.normalFieldList, p.normalFieldMap[field.Name])
		}

		for _, field := range visitor.enumFields {
			pterm.Info.Printfln("find enum field:%s in %s", field.Name, p.FullName)
			p.enumFieldMap[field.Name] = NewPbEnumField(field)
			p.enumFieldList = append(p.enumFieldList, p.enumFieldMap[field.Name])
		}
	}
}

func NewPbMessage(m *proto.Message) *PbMessage {
	p := &PbMessage{
		message: m,

		normalFieldMap: map[string]*PbNormalField{},
		mapFieldMap:    map[string]*PbMapField{},
		enumFieldMap:   map[string]*PbEnumField{},
	}
	p.walk()
	return p
}

type PbPackage struct {
	protoFilePath string

	proto *proto.Proto

	messages []*PbMessage
	enums    []*PbEnum
	services []*PbService
	rpcs     []*PbRPC
	options  []*PbOption

	msgMap     map[string]*PbMessage
	enumMap    map[string]*PbEnum
	serviceMap map[string]*PbService
	rpcMap     map[string]*PbRPC
	optionMap  map[string]*PbOption

	RawGoPackage   string
	RawPackageName string

	ProtoBuffer string

	// 自定义的一些字段
	Port int64
	Host string
}

func (p *PbPackage) PackageName() string {
	idx := strings.LastIndex(p.RawPackageName, ".")
	if idx > 0 {
		return p.RawPackageName[idx+1:]
	}

	return p.RawPackageName
}

func (p *PbPackage) ProtoFilePath() string {
	return p.protoFilePath
}

func (p *PbPackage) ProtoFileName() string {
	return filepath.Base(p.ProtoFilePath())
}

func (p *PbPackage) Proto() *proto.Proto {
	return p.proto
}

func (p *PbPackage) ProjectRoot() string {
	if state.Config.OutputPath == "" {
		return filepath.Join(filepath.Dir(filepath.Dir(p.protoFilePath)), p.RawGoPackage)
	} else {
		return filepath.Join(state.Config.OutputPath, p.RawGoPackage)
	}
}

func (p *PbPackage) GoPackage() string {
	if state.Config.GoModulePrefix != "" {
		return strings.TrimPrefix(filepath.ToSlash(filepath.Join(state.Config.GoModulePrefix, p.RawGoPackage)), "/")
	} else {
		return strings.TrimPrefix(filepath.ToSlash(p.RawGoPackage), "/")
	}
}

func (p *PbPackage) GoPackageName() string {
	return filepath.Base(p.GoPackage())
}

func (p *PbPackage) Messages() []*PbMessage {
	var messages []*PbMessage
	for _, m := range p.msgMap {
		messages = append(messages, m)
	}
	return messages
}

func (p *PbPackage) GetMessage(msg string) *PbMessage {
	return p.msgMap[msg]
}

func (p *PbPackage) Enums() []*PbEnum {
	return p.enums
}

func (p *PbPackage) GetEnum(e string) *PbEnum {
	return p.enumMap[e]
}

func (p *PbPackage) RPCs() []*PbRPC {
	return p.rpcs
}

func (p *PbPackage) GetRPC(rpc string) *PbRPC {
	return p.rpcMap[rpc]
}

func (p *PbPackage) Services() []*PbService {
	return p.services
}

func (p *PbPackage) GetService(s string) *PbService {
	return p.serviceMap[s]
}

func GetFullNames(e proto.Visitee) []string {
	var names []string

	var walk func(m proto.Visitee)
	walk = func(m proto.Visitee) {
		if m == nil {
			return
		}

		switch x := m.(type) {
		case *proto.Message:
			names = append(names, x.Name)
			walk(x.Parent)

		case *proto.Enum:
			names = append(names, x.Name)
			walk(x.Parent)

		case *proto.EnumField:
			names = append(names, x.Name)

			if x.Parent != nil {
				walk(x.Parent.(*proto.Enum).Parent)
			}

		case *proto.Proto:
		// do nothing

		case *proto.NormalField:
			names = append(names, x.Type)
			walk(x.Parent)

		default:
			log.Panicf("unknown parent type:%T", x)
		}
	}

	walk(e)

	return candy.Reverse(names)
}

func GetFullName(e proto.Visitee) string {
	return strings.Join(GetFullNames(e), "_")
}

func (p *PbPackage) Walk() {
	proto.Walk(p.proto,
		proto.WithMessage(func(m *proto.Message) {
			name := GetFullName(m)

			log.Infof("message:%v", m.Name)
			pterm.Info.Printfln("find message:%s", name)

			p.msgMap[name] = NewPbMessage(m)
			p.messages = append(p.messages, p.msgMap[name])
		}),
		proto.WithService(func(s *proto.Service) {
			log.Infof("service:%v", s.Name)
			pterm.Info.Printfln("find service:%s", s.Name)

			p.serviceMap[s.Name] = NewPbService(s)
			p.services = append(p.services, p.serviceMap[s.Name])
		}),
		proto.WithRPC(func(r *proto.RPC) {
			log.Infof("rpc:%v", r.Name)
			pterm.Info.Printfln("find rpc:%s", r.Name)

			p.rpcMap[r.Name] = NewPbRPC(r)
			p.rpcs = append(p.rpcs, p.rpcMap[r.Name])
		}),
		proto.WithEnum(func(e *proto.Enum) {
			name := GetFullName(e)

			log.Infof("enum:%v", e.Name)
			pterm.Info.Printfln("find enum:%s", name)

			p.enumMap[name] = NewPbEnum(e)
			p.enums = append(p.enums, p.enumMap[name])
		}),
		proto.WithOption(func(option *proto.Option) {
			log.Infof("option:%v", option.Name)
			pterm.Info.Printfln("find option:%s", option.Name)

			p.optionMap[option.Name] = NewPbOption(option)
			p.options = append(p.options, p.optionMap[option.Name])
		}),
		proto.WithPackage(func(pp *proto.Package) {
			log.Infof("package:%v", pp.Name)
			pterm.Info.Printfln("package:%s", pp.Name)

			p.RawPackageName = pp.Name
		}),
	)

	// 处理一下option

	if o, ok := p.optionMap["go_package"]; ok {
		p.RawGoPackage = o.Value
		idx := strings.Index(p.RawGoPackage, ";")
		if idx > 0 {
			p.RawGoPackage = p.RawGoPackage[:idx]
		}
	}

	if o, ok := p.optionMap["(lazygophers.lrpc.core.port)"]; ok {
		port, err := strconv.ParseInt(o.Value, 10, 64)
		if err != nil {
			log.Errorf("err:%v", err)
		} else {
			p.Port = port
		}
	}

	if o, ok := p.optionMap["(lazygophers.lrpc.core.host)"]; ok && o.Value != "" {
		p.Host = o.Value
	} else {
		p.Host = "*"
	}

	// 回填一些默认值
	candy.Each(p.rpcs, func(r *PbRPC) {
		r.setDefaultPackage(p.PackageName())
	})
}

func NewPbPackage(protoFilePath string, p *proto.Proto) *PbPackage {
	return &PbPackage{
		protoFilePath:  protoFilePath,
		proto:          p,
		messages:       nil,
		enums:          nil,
		services:       nil,
		rpcs:           nil,
		options:        nil,
		msgMap:         map[string]*PbMessage{},
		enumMap:        map[string]*PbEnum{},
		serviceMap:     map[string]*PbService{},
		rpcMap:         map[string]*PbRPC{},
		optionMap:      map[string]*PbOption{},
		RawGoPackage:   "",
		RawPackageName: "",
		ProtoBuffer:    "",
		Port:           8080,
	}
}

func ParseProto(path string, cacheFiles ...bool) (*PbPackage, error) {
	protoBuffer, err := os.ReadFile(path)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	pb, err := proto.NewParser(bytes.NewBuffer(protoBuffer)).Parse()
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	p := NewPbPackage(path, pb)

	if len(cacheFiles) > 0 {
		if cacheFiles[0] {
			p.ProtoBuffer = string(protoBuffer)
		}
	}

	p.Walk()

	return p, nil
}

type ProtoVisitor struct {
	msgList      []*proto.Message
	enumFields   []*proto.EnumField
	normalFields []*proto.NormalField
	mapFields    []*proto.MapField
	options      map[string]string
}

func (p *ProtoVisitor) VisitMessage(m *proto.Message) {
	p.msgList = append(p.msgList, m)
}

func (p *ProtoVisitor) VisitService(v *proto.Service) {
}

func (p *ProtoVisitor) VisitSyntax(s *proto.Syntax) {
}

func (p *ProtoVisitor) VisitPackage(pkg *proto.Package) {
}

func (p *ProtoVisitor) VisitOptions(o *proto.Option) {
}

func (p *ProtoVisitor) VisitOption(o *proto.Option) {
	if p.options == nil {
		p.options = map[string]string{}
	}
	p.options[o.Name] = o.Constant.Source
}

func (p *ProtoVisitor) VisitImport(i *proto.Import) {
}

func (p *ProtoVisitor) VisitNormalField(i *proto.NormalField) {
	p.normalFields = append(p.normalFields, i)
}

func (p *ProtoVisitor) VisitEnumField(i *proto.EnumField) {
	p.enumFields = append(p.enumFields, i)
}

func (p *ProtoVisitor) VisitEnum(e *proto.Enum) {
}

func (p *ProtoVisitor) VisitComment(e *proto.Comment) {
}

func (p *ProtoVisitor) VisitOneof(o *proto.Oneof) {
}

func (p *ProtoVisitor) VisitOneofField(o *proto.OneOfField) {
}

func (p *ProtoVisitor) VisitReserved(rs *proto.Reserved) {
}

func (p *ProtoVisitor) VisitRPC(rpc *proto.RPC) {
}

func (p *ProtoVisitor) VisitMapField(f *proto.MapField) {
	p.mapFields = append(p.mapFields, f)
}

func (p *ProtoVisitor) VisitGroup(g *proto.Group) {
}

func (p *ProtoVisitor) VisitExtensions(e *proto.Extensions) {
}

func (p *ProtoVisitor) VisitProto(*proto.Proto) {
}
