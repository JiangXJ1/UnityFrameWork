using System;
using System.Collections.Generic;
using System.Threading;

namespace Google.Protobuf
{
    internal static class ObjectPool<T> where T : class, IMessage<T>, new()
    {
        private static T[] m_Pool;
        private static volatile int m_PoolSize;
        private static volatile int m_Count;

        static ObjectPool()
        {
            m_PoolSize = 10;
            m_Count = 0;
            m_Pool = new T[m_PoolSize];
        }

        public static T Get()
        {
            while (m_Count > 0)
            {
                int index = m_Count - 1;
                var obj = m_Pool[index];
                if (Interlocked.CompareExchange(ref m_Count, index, index + 1) == index)
                {
                    if(Interlocked.CompareExchange(ref m_Pool[index], null, obj) == null)
                    {
                        return obj;
                    }
                    break;
                }
            }
            return new T();
        }

        public static void Recycle(T obj)
        {
            if(m_Count < m_PoolSize - 1)
            {
                obj.Clear();
                var index = m_Count + 1;
                if (Interlocked.CompareExchange(ref m_Pool[index], obj, null) == obj)
                {
                    m_Count++;
                }
            }
        }
    }
}
