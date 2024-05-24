package codegen

import (
	"bytes"
	"fmt"
	"github.com/emicklei/proto"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/candy"
	"github.com/pterm/pterm"
	"os"
	"path/filepath"
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
	return &PbOption{
		option: o,
	}
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
	return &PbComment{
		comment: c,
		tags:    map[string]*PbCommentTag{},
	}
}

type PbEnum struct {
	enum *proto.Enum
}

func (p *PbEnum) Enum() *proto.Enum {
	return p.enum
}

func (p *PbEnum) walk() {

}

func NewPbEnum(e *proto.Enum) *PbEnum {
	p := &PbEnum{
		enum: e,
	}

	return p
}

type PbRPC struct {
	rpc  *proto.RPC
	Name string
}

func (p *PbRPC) RPC() *proto.RPC {
	return p.rpc
}

func (p *PbRPC) walk() {
}

func NewPbRPC(rpc *proto.RPC) *PbRPC {
	p := &PbRPC{
		Name: rpc.Name,
		rpc:  rpc,
	}

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

	return p
}

type PbNormalField struct {
	Name  string
	field *proto.NormalField

	// 写在字段上的注释
	comment *PbComment

	// 写在字段后面的注释
	inlineComment *PbComment
}

func (p *PbNormalField) Field() *proto.NormalField {
	return p.field
}

func (p *PbNormalField) Type() string {
	return p.field.Type
}

func (p *PbNormalField) walk() {
	p.Name = p.field.Name

	if p.field.Comment != nil {
		p.comment = NewPbComment(p.field.Comment)
		p.comment.walk()
	}

	if p.field.InlineComment != nil {
		p.inlineComment = NewPbComment(p.field.InlineComment)
		p.inlineComment.walk()
	}
}

func NewPbNormalField(f *proto.NormalField) *PbNormalField {
	return &PbNormalField{
		field: f,
	}
}

type PbMapField struct {
	field *proto.MapField
}

func NewPbMapField(f *proto.MapField) *PbMapField {
	return &PbMapField{
		field: f,
	}
}

type PbEnumField struct {
	field *proto.EnumField
}

func NewPbEnumField(f *proto.EnumField) *PbEnumField {
	return &PbEnumField{
		field: f,
	}
}

type PbMessage struct {
	message      *proto.Message
	normalFields map[string]*PbNormalField
	mapFields    map[string]*PbMapField
	enumFields   map[string]*PbEnumField
	Name         string
}

func (p *PbMessage) Message() *proto.Message {
	return p.message
}

func (p *PbMessage) IsTable() bool {
	if strings.HasPrefix(p.Name, "Model") {
		if !strings.Contains(p.Name, "_") {
			return true
		}
	}

	// TODO: 允许通过对注释的解析判断时候是表

	return false
}

func (p *PbMessage) NeedOrm() bool {
	if strings.HasPrefix(p.Name, "Model") {
		return true
	}

	// TODO: 允许通过对注释的解析判断时候是表

	return false
}

func (p *PbMessage) walk() {
	p.Name = p.message.Name

	for _, element := range p.message.Elements {
		var visitor = &ProtoVisitor{}

		element.Accept(visitor)

		for _, field := range visitor.mapFields {
			pterm.Info.Printfln("find map field:%s in %s", field.Name, p.Name)
			p.mapFields[field.Name] = NewPbMapField(field)
		}

		for _, field := range visitor.normalFields {
			pterm.Info.Printfln("find normal field:%s in %s", field.Name, p.Name)
			p.normalFields[field.Name] = NewPbNormalField(field)
			p.normalFields[field.Name].walk()
		}

		for _, field := range visitor.enumFields {
			pterm.Info.Printfln("find enum field:%s in %s", field.Name, p.Name)
			p.enumFields[field.Name] = NewPbEnumField(field)
		}
	}
}

func NewPbMessage(m *proto.Message) *PbMessage {
	return &PbMessage{
		message: m,

		normalFields: map[string]*PbNormalField{},
		mapFields:    map[string]*PbMapField{},
		enumFields:   map[string]*PbEnumField{},
	}
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

	RawGoPackage string
	PackageName  string

	ProtoBuffer string
}

func (p *PbPackage) ProtoFilePath() string {
	return p.protoFilePath
}

func (p *PbPackage) ProtoFileName() string {
	return p.proto.Filename
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
		return fmt.Sprintf("%s/%s",
			strings.TrimSuffix(state.Config.GoModulePrefix, "/"),
			strings.TrimPrefix(p.RawGoPackage, "/"))
	} else {
		return p.RawGoPackage
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

func (p *PbPackage) GetMessageFullName(e *proto.Message) string {
	var names []string

	var walk func(m *proto.Message)
	walk = func(m *proto.Message) {
		if m == nil {
			return
		}

		names = append(names, m.Name)

		if m.Parent == nil {
			return
		}

		switch x := m.Parent.(type) {
		case *proto.Message:
			walk(x)
		default:
			log.Warnf("unknown parent type:%T", x)
		}
	}

	walk(e)

	return strings.Join(candy.Reverse(names), "_")
}

func (p *PbPackage) Walk() {
	proto.Walk(p.proto,
		proto.WithMessage(func(m *proto.Message) {
			m.Name = p.GetMessageFullName(m)

			log.Infof("message:%v", m.Name)
			pterm.Info.Printfln("find message:%s", m.Name)

			p.msgMap[m.Name] = NewPbMessage(m)
			p.msgMap[m.Name].walk()
			p.messages = append(p.messages, p.msgMap[m.Name])
		}),
		proto.WithService(func(s *proto.Service) {
			log.Infof("service:%v", s.Name)
			pterm.Info.Printfln("find service:%s", s.Name)

			p.serviceMap[s.Name] = NewPbService(s)
			p.serviceMap[s.Name].walk()
			p.services = append(p.services, p.serviceMap[s.Name])
		}),
		proto.WithRPC(func(r *proto.RPC) {
			log.Infof("rpc:%v", r.Name)
			pterm.Info.Printfln("find rpc:%s", r.Name)

			p.rpcMap[r.Name] = NewPbRPC(r)
			p.rpcMap[r.Name].walk()
			p.rpcs = append(p.rpcs, p.rpcMap[r.Name])
		}),
		proto.WithEnum(func(e *proto.Enum) {
			log.Infof("enum:%v", e.Name)
			pterm.Info.Printfln("find enum:%s", e.Name)

			p.enumMap[e.Name] = NewPbEnum(e)
			p.enumMap[e.Name].walk()
			p.enums = append(p.enums, p.enumMap[e.Name])
		}),
		proto.WithOption(func(option *proto.Option) {
			log.Infof("option:%v", option.Name)
			pterm.Info.Printfln("find option:%s", option.Name)

			p.optionMap[option.Name] = NewPbOption(option)
			p.optionMap[option.Name].walk()
			p.options = append(p.options, p.optionMap[option.Name])
		}),
		proto.WithPackage(func(pp *proto.Package) {
			log.Infof("package:%v", pp.Name)
			pterm.Info.Printfln("package:%s", pp.Name)

			p.PackageName = pp.Name
		}),
	)

	if o, ok := p.optionMap["go_package"]; ok {
		p.RawGoPackage = o.Value
	}
}

func NewPbPackage(protoFilePath string, p *proto.Proto) *PbPackage {
	return &PbPackage{
		protoFilePath: protoFilePath,
		proto:         p,
		messages:      nil,
		enums:         nil,
		services:      nil,
		rpcs:          nil,
		options:       nil,
		msgMap:        map[string]*PbMessage{},
		enumMap:       map[string]*PbEnum{},
		serviceMap:    map[string]*PbService{},
		rpcMap:        map[string]*PbRPC{},
		optionMap:     map[string]*PbOption{},
		RawGoPackage:  "",
		PackageName:   "",
		ProtoBuffer:   "",
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
