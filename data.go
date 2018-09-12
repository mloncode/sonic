package sonic

type Changes struct {
	Old []sonicNode
	New []sonicNode
}

var File1 = Changes{
	Old: []sonicNode{
		{"FuncDecl", "TestIntegration", 8427},
		{"FuncDecl", "TestUastQueries", 1625},
		{"FuncDecl", "TestSquashCorrectness ", 3295},
		{"FuncDecl", "queryResults", 347},
		{"FuncDecl", "TestMissingHeadRefs", 751},
		{"FuncDecl", "BenchmarkQueries", 3601},
		{"FuncDecl", "benchmarkQuery", 252},
		{"FuncDecl", "TestIndexes", 3147},
		{"FuncDecl", "col", 285},
		{"FuncDecl", "createTestIndexes", 1488},
		{"FuncDecl", "createIndex", 311},
		{"FuncDecl", "deleteIndex", 175},
		{"FuncDecl", "setup", 370},
		{"FuncDecl", "newSquashEngine", 351},
		{"FuncDecl", "newBaseEngine", 219},
	},
	New: []sonicNode{
		{"FuncDecl", "TestIntegration", 8427},
		{"FuncDecl", "TestUastQueries", 1625},
		{"FuncDecl", "TestSquashCorrectness", 3295},
		{"FuncDecl", "queryResults", 347},
		{"FuncDecl", "TestMissingHeadRefs", 751},
		{"FuncDecl", "BenchmarkQueries", 3601},
		{"FuncDecl", "benchmarkQuery", 252},
		{"FuncDecl", "TestIndexes", 3147},
		{"FuncDecl", "col", 285},
		{"FuncDecl", "createTestIndexes", 1488},
		{"FuncDecl", "createIndex", 311},
		{"FuncDecl", "deleteIndex", 175},
		{"FuncDecl", "setup", 370},
		{"FuncDecl", "newSquashEngine", 351},
		{"FuncDecl", "newBaseEngine", 219},
	},
}

// [2018-09-11T16:31:14.438801702+02:00]  INFO got change
// old uast nodes:
// {FuncDecl NewUAST 375}
// {FuncDecl IsNullable 153}
// {FuncDecl Resolved 145}
// {FuncDecl Type 61}
// {FuncDecl Children 206}
// {FuncDecl TransformUp 450}
// {FuncDecl String 258}
// {FuncDecl Eval 725}
// {FuncDecl NewUASTMode 106}
// {FuncDecl IsNullable 113}
// {FuncDecl Resolved 105}
// {FuncDecl Type 65}
// {FuncDecl Children 208}
// {FuncDecl TransformUp 398}
// {FuncDecl String 106}
// {FuncDecl Eval 1028}
// {FuncDecl NewUASTXPath 138}
// {FuncDecl Type 64}
// {FuncDecl Eval 1197}
// {FuncDecl String 97}
// {FuncDecl TransformUp 276}
// {FuncDecl exprToString 318}
// {FuncDecl getUAST 1232}
// new uast nodes:
// {FuncDecl NewUAST 375}
// {FuncDecl IsNullable 153}
// {FuncDecl Resolved 145}
// {FuncDecl Type 61}
// {FuncDecl Children 206}
// {FuncDecl TransformUp 450}
// {FuncDecl String 258}
// {FuncDecl Eval 725}
// {FuncDecl NewUASTMode 104}
// {FuncDecl IsNullable 113}
// {FuncDecl Resolved 105}
// {FuncDecl Type 65}
// {FuncDecl Children 208}
// {FuncDecl TransformUp 398}
// {FuncDecl String 106}
// {FuncDecl Eval 1028}
// {FuncDecl NewUASTXPath 138}
// {FuncDecl Type 64}
// {FuncDecl Eval 1197}
// {FuncDecl String 97}
// {FuncDecl TransformUp 276}
// {FuncDecl exprToString 318}
// {FuncDecl getUAST 1232}
// [2018-09-11T16:31:14.803750523+02:00]  INFO got change
// old uast nodes:
// {FuncDecl TestUASTMode 1005}
// {FuncDecl TestUAST 1519}
// {FuncDecl TestUASTXPath 1146}
// {FuncDecl assertUASTBlobs 654}
// {FuncDecl bblfshFixtures 1173}
// {FuncDecl setup 392}
// new uast nodes:
// {FuncDecl TestUASTMode 1005}
// {FuncDecl TestUAST 1519}
// {FuncDecl TestUASTXPath 1146}
// {FuncDecl assertUASTBlobs 654}
// {FuncDecl bblfshFixtures 1173}
// {FuncDecl setup 392}
