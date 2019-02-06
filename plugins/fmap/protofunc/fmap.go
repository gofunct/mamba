package protofunc

import "text/template"

var (
	DefaultProtoFmap = template.FuncMap{
		"getProtoFile":                 GetProtoFile,
		"getMessageType":               GetMessageType,
		"getEnumValue":                 GetEnumValue,
		"isFieldMessage":               IsFieldMessage,
		"isFieldMessageTimeStamp":      IsFieldMessageTimeStamp,
		"isFieldRepeated":              IsFieldRepeated,
		"haskellType":                  HaskellType,
		"goType":                       GoType,
		"goZeroValue":                  GoZeroValue,
		"goTypeWithPackage":            GoTypeWithGoPackage,
		"goTypeWithGoPackage":          GoTypeWithGoPackage,
		"jsType":                       JsType,
		"jsSuffixReserved":             JsSuffixReservedKeyword,
		"httpVerb":                     HttpVerb,
		"httpPath":                     HttpPath,
		"httpPathsAdditionalBindings":  HttpPathsAdditionalBindings,
		"httpBody":                     HttpBody,
		"shortType":                    ShortType,
		"urlHasVarsFromMessage":        UrlHasVarsFromMessage,
		"leadingComment":               LeadingComment,
		"trailingComment":              TrailingComment,
		"leadingDetachedComments":      LeadingDetachedComments,
		"stringFieldExtension":         StringFieldExtension,
		"int64FieldExtension":          Int64FieldExtension,
		"stringMethodOptionsExtension": StringMethodOptionsExtension,
		"boolMethodOptionsExtension":   BoolMethodOptionsExtension,
		"boolFieldExtension":           BoolFieldExtension,
		"isFieldMap":                   IsFieldMap,
		"fieldMapKeyType":              FieldMapKeyType,
	}
)
