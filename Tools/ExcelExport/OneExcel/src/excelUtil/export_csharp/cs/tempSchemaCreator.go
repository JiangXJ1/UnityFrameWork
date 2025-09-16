package cs

const ConfigEnum = `namespace Schema
{

	public enum SchemaID
	{
		//SCHEMA_ID_BEGIN
		Max,
	}
	public class SchemaCreator
	{
		public static ISchema GetSchema(SchemaID nID)
		{
			switch (nID)
			{
				//SCHEMA_BEGIN
				default: return null;
			}
		}
	}
}
`
