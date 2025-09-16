package export

import (
	"bytes"
	"text/template"
)

type SchemaCreatorTemplateParam struct {
	EnumContent string
	Content     string
}

const SchemaCreatorTemplate = `namespace Schema
{

	public enum SchemaID
	{
{{.EnumContent}}
		Max,
	}
	public class SchemaCreator
	{
		public static ISchema GetSchema(SchemaID nID)
		{
			switch (nID)
			{
{{.Content}}
				default: return null;
			}
		}
	}
}
`

func (c *SchemaCreatorTemplateParam) GenerateCsharpTemplate() string {
	var buf bytes.Buffer
	t := template.Must(template.New("deserialize").Parse(SchemaCreatorTemplate))
	if err := t.Execute(&buf, c); err != nil {
		panic("template execution failed: " + err.Error())
	}
	return buf.String()
}
