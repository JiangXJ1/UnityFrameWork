using System;
using System.IO;
using tabtoy;

namespace csharptest
{
    class Program
    {

        static void Main(string[] args)
        {            
            using (var stream = new FileStream("../../Config.bin", FileMode.Open))
            {
                stream.Position = 0;

                var reader = new tabtoy.DataReader(stream);
                var config = new table.Sample();

                //int a = reader.ReadInt32();
                //reader.ReadBytes();
                //int b = reader.ReadInt32();
                //byte [] content = reader.ReadBytes();
                //Console.WriteLine(a + b);

                table.Sample.Deserialize(config, reader);

                // 直接通过下标获取或遍历
                var directFetch = config.Data[201];

                var a = directFetch.Key;

                // 添加日志输出或自定义输出
                //config.TableLogger.AddTarget(new tabtoy.DebuggerTarget());

                // 取空时, 当默认值不为空时, 输出日志
                //var nullFetchOutLog = config.GetSampleByID(0);

                string s = typeof(Int32).ToString();
                Console.WriteLine(s);
            }
        }
    }
}
