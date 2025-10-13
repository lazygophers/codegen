package i18n

import (
	"errors"
	"github.com/lazygophers/codegen/cli/utils"
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/codegen/i18n"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/candy"
	"github.com/lazygophers/utils/osx"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"path/filepath"
	"strings"
)

var allLanguage = candy.Map([]string{
	"af",
	"ak",
	"am",
	"ar",
	"as",
	"ay",
	"az",
	"be",
	"bg",
	"bho",
	"bm",
	"bn",
	"bs",
	"ca",
	"ceb",
	"ckb",
	"co",
	"cs",
	"cy",
	"da",
	"de",
	"doi",
	"dv",
	"ee",
	"el",
	"en",
	"eo",
	"es",
	"et",
	"eu",
	"fa",
	"fi",
	"fr",
	"fy",
	"ga",
	"gd",
	"gl",
	"gn",
	"gom",
	"gu",
	"ha",
	"haw",
	"hi",
	"hmn",
	"hr",
	"ht",
	"hu",
	"hy",
	"id",
	"ig",
	"ilo",
	"is",
	"it",
	"he",
	"iw",
	"ja",
	"jw",
	"jv",
	"ka",
	"kk",
	"km",
	"kn",
	"ko",
	"kri",
	"ku",
	"ky",
	"la",
	"lb",
	"lg",
	"ln",
	"lo",
	"lt",
	"lus",
	"lv",
	"mai",
	"mg",
	"mi",
	"mk",
	"ml",
	"mn",
	"mni-mtei",
	"mr",
	"ms",
	"mt",
	"my",
	"ne",
	"nl",
	"no",
	"nso",
	"ny",
	"om",
	"or",
	"pa",
	"pl",
	"ps",
	"pt",
	"pt-br",
	"qu",
	"ro",
	"ru",
	"rw",
	"sa",
	"sd",
	"si",
	"sk",
	"sl",
	"sm",
	"sn",
	"so",
	"sq",
	"sr",
	"st",
	"su",
	"sv",
	"sw",
	"ta",
	"te",
	"tg",
	"th",
	"ti",
	"tk",
	"tl",
	"tr",
	"ts",
	"tt",
	"ug",
	"uk",
	"ur",
	"uz",
	"vi",
	"xh",
	"yi",
	"yo",
	"zh-cn",
	"zh-tw",
	"zh-chs",
	"zh-hans",
	"zh-hk",
	"zh-mo",
	"zh-sg",
	"zh-cht",
	"zu",
}, i18n.MustParseLanguage)

var tranCmd = &cobra.Command{
	Use:     "tran",
	Short:   "translate",
	Aliases: []string{"t", "translate"},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		mergeTranCmdFlags(cmd)

		srcFile := utils.GetString("src-file", cmd)
		if srcFile == "" {
			pterm.Error.Printfln("src-file is empty")
			return errors.New("src-file is empty")
		}

		if !osx.IsFile(srcFile) {
			pterm.Warning.Printfln("%s is not found", srcFile)
			return errors.New("src-file not found")
		}

		c := &i18n.TranslateConfig{
			SrcFile: srcFile,
			OverwriteKeyPrefix: candy.Map(utils.GetStringSlice("force-key-prefix", cmd), func(key string) string {
				return "." + key
			}),
			Overwrite: state.Config.Overwrite,
			AutoTran:  state.Config.I18n.AutoTran,
		}

		if cmd.Flag("src-language").Changed {
			c.SrcLang, err = i18n.ParseLanguage(utils.GetString("src-language", cmd))
			if err != nil {
				log.Errorf("err:%v", err)
				pterm.Error.Printfln("unknown language %s", c.SrcLang)
				return err
			}
		} else {
			c.SrcLang, err = i18n.ParseLanguage(strings.TrimSuffix(filepath.Base(c.SrcFile), filepath.Ext(c.SrcFile)))
			if err != nil {
				log.Errorf("err:%v", err)
				pterm.Error.Printfln("unknown language %s", c.SrcLang)
				return err
			}
		}

		for _, s := range state.Config.I18n.Languages {
			lang, err := i18n.ParseLanguage(s)
			if err != nil {
				log.Errorf("err:%v", err)
				pterm.Error.Printfln("unknown language %s, skip", lang)
			}

			c.Langs = append(c.Langs, lang)
		}

		if state.Config.I18n.AllLanguages {
			c.Langs = allLanguage
		} else if len(c.Langs) == 0 {
			pterm.Warning.Printfln("no target language")
			return errors.New("no target language")
		} else {
			pterm.Info.Printfln("target languages: \n\t%s", strings.Join(candy.Map(c.Langs, func(l *i18n.Language) string {
				return l.Lang
			}), "\n\t"))
		}

		var ok bool
		c.Localizer, ok = i18n.GetLocalizer(filepath.Ext(srcFile))
		if !ok {
			pterm.Warning.Printfln("no localizer for %s", filepath.Ext(srcFile))
			return errors.New("no localizer for " + filepath.Ext(srcFile))
		}

		c.Translator, ok = i18n.GetTranslator(i18n.TranslateType(state.Config.I18n.Translator))
		if !ok {
			pterm.Warning.Printfln("no translator for %s", state.Config.I18n.Translator)
			return errors.New("no translator be found")
		}

		err = i18n.Translate(c)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		if state.Config.I18n.GenerateConst {
			dstLocalize, err := i18n.LoadLocalize(c.SrcFile, c.Localizer)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

			err = codegen.GenerateI18nConst(
				dstLocalize,
				filepath.Join(
					filepath.Dir(filepath.Dir(c.SrcFile)),
					"i18n_const.gen.go"),
			)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
		} else if state.Config.I18n.GenerateField {
			dstLocalize, err := i18n.LoadLocalize(c.SrcFile, c.Localizer)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

			err = codegen.GenerateI18nField(
				dstLocalize,
				filepath.Join(
					filepath.Dir(filepath.Dir(c.SrcFile)),
					"i18n_field.gen.go"),
			)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
		}

		return nil
	},
}

func mergeTranCmdFlags(cmd *cobra.Command) {
	if cmd.Flag("generate-const").Changed {
		state.Config.I18n.GenerateConst = utils.GetBool("generate-const", cmd)
	}

	if cmd.Flag("generate-field").Changed {
		state.Config.I18n.GenerateField = utils.GetBool("generate-field", cmd)
	}

	if cmd.Flag("translator").Changed {
		state.Config.I18n.Translator = utils.GetString("translator", cmd)
	}

	if cmd.Flag("all-languages").Changed {
		state.Config.I18n.AllLanguages = utils.GetBool("all-languages", cmd)
	}

	if cmd.Flag("languages").Changed {
		state.Config.I18n.Languages = utils.GetStringSlice("languages", cmd)
	}

	if cmd.Flag("force").Changed {
		state.Config.Overwrite = utils.GetBool("force", cmd)
	}

	if cmd.Flag("auto-tran-enable-record").Changed {
		state.Config.I18n.AutoTran.EnableRecord = utils.GetBool("auto-tran-enable-record", cmd)
	}

	if cmd.Flag("auto-tran-record-path").Changed {
		state.Config.I18n.AutoTran.RecordPath = utils.GetString("auto-tran-record-path", cmd)
	}

}

func initTran() {
	tranCmd.Short = state.Localize(state.I18nTagCliI18nTranShort)
	tranCmd.Long = state.Localize(state.I18nTagCliI18nLong)

	tranCmd.Flags().StringP("src-file", "s", "", state.Localize(state.I18nTagCliI18nTranFlagsSrcFile))
	tranCmd.Flags().String("src-language", "", state.Localize(state.I18nTagCliI18nTranFlagsSrcLanguage))

	tranCmd.Flags().Bool("all-languages", state.Config.I18n.AllLanguages, state.Localize(state.I18nTagCliI18nTranFlagsAllLanguage))
	tranCmd.Flags().StringSliceP("languages", "l", state.Config.I18n.Languages, state.Localize(state.I18nTagCliI18nTranFlagsLanguages))

	tranCmd.Flags().Bool("generate-const", state.Config.I18n.GenerateConst, state.Localize(state.I18nTagCliI18nTranFlagsGenerateConst))
	tranCmd.Flags().Bool("generate-field", state.Config.I18n.GenerateField, state.Localize(state.I18nTagCliI18nTranFlagsGenerateField))
	tranCmd.Flags().String("translator", state.Config.I18n.Translator, state.Localize(state.I18nTagCliI18nTranFlagsTranslatorUsage, map[string]any{
		"Type": map[i18n.TranslateType]string{
			i18n.TranslateTypeGoogleFree: state.Localize(state.I18nTagCliI18nTranFlagsTranslatorGoogleFree),
		},
	}))

	tranCmd.Flags().Bool("auto-tran-enable-record", state.Config.I18n.AutoTran.EnableRecord, state.Localize(state.I18nTagCliI18nTranFlagsAutoTran))
	tranCmd.Flags().String("auto-tran-record-path", state.Config.I18n.AutoTran.RecordPath, state.Localize(state.I18nTagCliI18nTranFlagsAutoTranRecordPath))

	tranCmd.Flags().StringSlice("force-key-prefix", []string{}, state.Localize(state.I18nTagCliI18nTranFlagsForce))

	tranCmd.Flags().BoolP("force", "f", false, state.Localize(state.I18nTagCliI18nTranFlagsForce))

	i18nCmd.AddCommand(tranCmd)
}
