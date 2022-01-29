# hsp_utils
utils for hsp projects

## stringUtils
Strval: type interface to type string

StrFirstToUpper: 字符串首字母转化为大写 ios_bbbbbbbb -> Ios_bbbbbbbb

StrDBToMapKey：数据库字母转化为原格式 ios_bbbbbbbb -> IosBbbbbbbbb

WriteSTInfoJsonFiles:将interface类型数据写入json文件(可查看测试程序)

ReadSTInfoJsonFiles_Example：将json文件读取到结构体中，属于示例程序，可以在其他模块仿照此程序写一个解析函数 ，因为每个json文件对应的结构体不一样

ReadSTInfoJsonFiles:将json文件读取到map[string]interface{}的map中

DelFiles:删除指定文件

CheckFileExist：确定文件是否存在
## common
CreateRandFloat64List: 通过最小值、最大值和数量，生成随机浮点型数组；

CreateRandIntList：通过最小值、最大值和数量，生成随机整型数组

CreatRandBoolList：通过数量生成bool型随机数组

ConvInterface2Float：将interface型数据转换为float64型数据

ConvInterface2Int64：将interface型数据转换为int64型数据

ConvInterface2Bool：将interface型数据转换为bool型数据

GetCurrentDir:获取当前的路径信息
## protocalUtils
AnalysisSelfProtocalsMsg： 解析协议

PackageSelfProtocalsMsg：打包协议