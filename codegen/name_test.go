package codegen_test

import (
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/codegen/state"
	"gotest.tools/v3/assert"
	"testing"
)

func TestNameStyle(t *testing.T) {
	type want struct {
		s string
	}
	type args struct {
		style state.CfgStyleName
	}
	var tests = []struct {
		name string
		want want
		args args
	}{
		{
			name: "default",
			want: want{
				s: "AddUserRecord",
			},
			args: args{
				style: state.CfgStyleNameDefault,
			},
		},

		{
			name: "snake",
			want: want{
				s: "add_user_record",
			},
			args: args{
				style: state.CfgStyleNameSnake,
			},
		},
		{
			name: "pascal",
			want: want{
				s: "AddUserRecord",
			},
			args: args{
				style: state.CfgStyleNamePascal,
			},
		},
		{
			name: "camel",
			want: want{
				s: "addUserRecord",
			},
			args: args{
				style: state.CfgStyleNameCamel,
			},
		},
		{
			name: "kebab",
			want: want{
				s: "add-user-record",
			},
			args: args{
				style: state.CfgStyleNameKebab,
			},
		},

		{
			name: "slash",
			want: want{
				s: "add/user/record",
			},
			args: args{
				style: state.CfgStyleNameSlash,
			},
		},
		{
			name: "slash-snake",
			want: want{
				s: "add/user_record",
			},
			args: args{
				style: state.CfgStyleNameSlashSnake,
			},
		},
		{
			name: "slash-kebab",
			want: want{
				s: "add/user-record",
			},
			args: args{
				style: state.CfgStyleNameSlashKebab,
			},
		},
		{
			name: "slash-camel",
			want: want{
				s: "add/userRecord",
			},
			args: args{
				style: state.CfgStyleNameSlashCamel,
			},
		},
		{
			name: "slash-pascal",
			want: want{
				s: "add/UserRecord",
			},
			args: args{
				style: state.CfgStyleNameSlashPascal,
			},
		},

		{
			name: "slash-reverse",
			want: want{
				s: "record/user/add",
			},
			args: args{
				style: state.CfgStyleNameSlashReverse,
			},
		},
		{
			name: "slash-reverse-snake",
			want: want{
				s: "user_record/add",
			},
			args: args{
				style: state.CfgStyleNameSlashReverseSnake,
			},
		},
		{
			name: "slash-reverse-kebab",
			want: want{
				s: "user-record/add",
			},
			args: args{
				style: state.CfgStyleNameSlashReverseKebab,
			},
		},
		{
			name: "slash-reverse-camel",
			want: want{
				s: "userRecord/add",
			},
			args: args{
				style: state.CfgStyleNameSlashReverseCamel,
			},
		},
		{
			name: "slash-reverse-pascal",
			want: want{
				s: "UserRecord/add",
			},
			args: args{
				style: state.CfgStyleNameSlashReversePascal,
			},
		},

		{
			name: "dot",
			want: want{
				s: "add.user.record",
			},
			args: args{
				style: state.CfgStyleNameDot,
			},
		},
		{
			name: "dot-snake",
			want: want{
				s: "add.user_record",
			},
			args: args{
				style: state.CfgStyleNameDotSnake,
			},
		},
		{
			name: "dot-kebab",
			want: want{
				s: "add.user-record",
			},
			args: args{
				style: state.CfgStyleNameDotKebab,
			},
		},
		{
			name: "dot-camel",
			want: want{
				s: "add.userRecord",
			},
			args: args{
				style: state.CfgStyleNameDotCamel,
			},
		},
		{
			name: "dot-pascal",
			want: want{
				s: "add.UserRecord",
			},
			args: args{
				style: state.CfgStyleNameDotPascal,
			},
		},

		{
			name: "dot-reverse",
			want: want{
				s: "record.user.add",
			},
			args: args{
				style: state.CfgStyleNameDotReverse,
			},
		},
		{
			name: "dot-reverse-snake",
			want: want{
				s: "user_record.add",
			},
			args: args{
				style: state.CfgStyleNameDotReverseSnake,
			},
		},
		{
			name: "dot-reverse-kebab",
			want: want{
				s: "user-record.add",
			},
			args: args{
				style: state.CfgStyleNameDotReverseKebab,
			},
		},
		{
			name: "dot-reverse-camel",
			want: want{
				s: "userRecord.add",
			},
			args: args{
				style: state.CfgStyleNameDotReverseCamel,
			},
		},
		{
			name: "dot-reverse-pascal",
			want: want{
				s: "UserRecord.add",
			},
			args: args{
				style: state.CfgStyleNameDotReversePascal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, codegen.NameStyle[tt.args.style]("AddUserRecord", "ModelUserRecord"), tt.want.s)
		})
	}
}
