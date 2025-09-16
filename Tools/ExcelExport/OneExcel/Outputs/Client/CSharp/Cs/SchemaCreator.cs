namespace Schema
{

	public enum SchemaID
	{
		//SCHEMA_ID_BEGIN
		Schema_AnimationSet,
		Schema_Animation,
		Max,
	}
	public class SchemaCreator
	{
		public static ISchema GetSchema(SchemaID nID)
		{
			switch (nID)
			{
				//SCHEMA_BEGIN
				case SchemaID.Schema_AnimationSet: return new Schema.Schema_AnimationSet();
				case SchemaID.Schema_Animation: return new Schema.Schema_Animation();
				default: return null;
			}
		}
	}
}
