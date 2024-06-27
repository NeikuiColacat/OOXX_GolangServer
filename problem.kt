/*
N 为题目矩阵大小
alpha为超参数，保证生成题目挖空cell数量不超过alpha,用于控制难度
成员变量a为一个二维可变数组 , 表示题目答案
成员变量b为一个二维可变数组，表示生成题目
成员变量c为一个字符串list，0表示O ， 1表示X ，空格字符表示-1
-1表示挖空 , 0 和 1用于表示 OO和XX ， 依照喜好自行定义
*/
import kotlin.random.Random
class problem (val N : Int , val alpha : Int){
    val a = MutableList(N){MutableList(N){-1} }
    var b = MutableList(N){MutableList(N){-1} }
    val c = MutableList(0){""}
    val s = mutableSetOf<Int>()
    fun ran(mod : Int) :Int{
        return (Random.nextInt()%mod + mod)%mod
    }

    fun getcol(pos:Int,type:Int,a:MutableList<MutableList<Int>>):Int {
        var res = 0
        for(i in 0..N-1) res += if (a[i][pos] == type) 1 else 0
        return res
    }

    fun getrow(pos:Int,type:Int,a:MutableList<MutableList<Int>>):Int {
        var res = 0
        for(i in 0..N-1) res += if (a[pos][i] == type) 1 else 0
        return res
    }

    fun chkconsec(u:Int,v:Int,a:MutableList<MutableList<Int>>) : Boolean{
        val T = a[u][v]
        var (l , r , cnt) = Triple(u , u , 0)
        while(l >= 0 && a[l][v] == T) l--
        while(r < N && a[r][v] == T) r++
        if(r - l - 1 > 2)  return false

        l = v ; r = v
        while(l >= 0 && a[u][l] == T) l--
        while(r < N && a[u][r] == T) r++
        if(r - l - 1 > 2)  return false

        return true
    }

    fun chknum(u:Int,v:Int,a:MutableList<MutableList<Int>>) :Boolean{
        if(getrow(u,1,a)> N/2 || getrow(u,0,a)> N/2) return false
        if(getcol(v,1,a)> N/2 || getcol(v,0,a)> N/2) return false
        return true
    }

    fun chkunique(a:MutableList<MutableList<Int>>):Int{
        s.clear()
        for(i in 0..N-1) {
            var num = 0 ; var num1 = 0
            for(j in 0..N-1){
                num = num or (a[i][j] shl j)
                num1 = num1 or (a[j][i] shl j)
            }
            if(s.contains(num)) return 0
            s.add(num)
            if(s.contains(num1)) return 0
            s.add(num1)
        }
        return 1
    }

    fun dfs(u:Int,v:Int,fl:Int,a:MutableList<MutableList<Int>>):Int{
        var ok = 0
        if(u == N){
            return chkunique(a)
        }
        if(a[u][v]!=-1) {
            ok += if(v == N-1) dfs(u + 1 , 0 , fl , a) else dfs(u , v + 1 , fl , a)
            return ok
        }

        a[u][v] = ran(2)
        if(chknum(u,v,a) && chkconsec(u,v,a)) {
            ok += if(v == N-1) dfs(u + 1 , 0 , fl , a) else dfs(u , v + 1 , fl , a)
            if(ok - fl> 0) return ok
        }

        a[u][v] = a[u][v] xor 1
        if(chknum(u,v,a) && chkconsec(u,v,a)) {
            ok += if(v == N-1) dfs(u + 1 , 0 , fl , a) else dfs(u , v + 1 , fl , a)
            if(ok - fl> 0) return ok
        }

        a[u][v] = -1
        return ok
    }

    fun modify(){
        val cad : MutableList<Pair<Int,Int>> = mutableListOf()
        for(i in 0..N-1)  for(j in 0..N-1) cad.add(Pair(i,j))
        cad.shuffle()
        while(cad.size > 0 && N*N - cad.size < alpha) {
            val (u , v) = cad.last()
            cad.removeAt(cad.size - 1)
            val tmp = b.map{it.toMutableList()}.toMutableList()
            tmp[u][v] = -1
            if(dfs(0,0,1,tmp) > 1) continue
            b[u][v] = -1
        }
    }

    init{
        dfs(0,0,0,a)
        b = a.map { it.toMutableList() }.toMutableList()
        modify()
        for(i in 0..N-1) c.add(b[i].map {if(it == -1)' ' else (if(it == 1) 'X' else 'O')}.joinToString(separator = ""))
    }
}

/*
代码设计思路
这个代码旨在生成一个 OOXX 填字游戏的谜题，满足以下条件：

每一行和每一列中没有超过两个连续的 X 或 O。
每一行和每一列中的 X 和 O 数量相同。
每一行和每一列都是唯一的。

N: 矩阵的大小。
alpha: 超参数，用于控制挖空的单元格数量，以控制谜题的难度。
成员变量
a: 一个二维可变数组，表示谜题的完整答案。
b: 一个二维可变数组，表示生成的谜题（带有挖空）。
c: 一个字符串列表，0 表示 O，1 表示 X，空格字符表示 -1（挖空）。
s: 一个集合，用于检查行和列的唯一性。
成员函数
ran(mod: Int): Int
生成一个随机数，用于决定当前格子是 X 还是 O。

getcol(pos: Int, type: Int, a: MutableList<MutableList<Int>>): Int
计算给定列中指定类型（X 或 O）的数量。

getrow(pos: Int, type: Int, a: MutableList<MutableList<Int>>): Int
计算给定行中指定类型（X 或 O）的数量。

chkconsec(u: Int, v: Int, a: MutableList<MutableList<Int>>): Boolean
检查给定位置是否有超过两个连续的 X 或 O。

chknum(u: Int, v: Int, a: MutableList<MutableList<Int>>): Boolean
检查给定行和列中 X 和 O 的数量是否超过一半。

chkunique(a: MutableList<MutableList<Int>>): Int
检查每一行和每一列是否唯一。

dfs(u: Int, v: Int, fl: Int, a: MutableList<MutableList<Int>>): Int
使用深度优先搜索（DFS）来生成符合规则的矩阵。

modify()
挖空一些单元格以生成谜题，确保挖空后的谜题仍然有唯一解。

初始化过程
在类的初始化部分：

使用 DFS 生成完整的谜题答案 a。
将 a 复制到 b，作为初始谜题。
调用 modify() 函数，挖空部分单元格，生成最终谜题。
将谜题转换为字符串列表 c，以便于输出和显示。


DFS 过程的算法思路和实现
算法思路
深度优先搜索（DFS）在这个代码中的主要目的是生成一个符合条件的谜题答案矩阵。这个过程通过递归遍历矩阵的每一个单元格，并在填充时确保以下三个规则始终得到满足：

没有超过两个连续的 X 或 O： 这一规则防止了单调重复的出现，使谜题更加有趣。
X 和 O 的数量相同： 这一规则确保每一行和每一列中的 X 和 O 数量相同，增加了解题的挑战性。
每一行和每一列都是唯一的： 这一规则保证了谜题的多样性和唯一性。
DFS 的基本思想是从矩阵的左上角开始，逐行逐列地填充矩阵。当一个单元格填充完毕后，递归地处理下一个单元格。在填充每一个单元格时，通过随机选择 X 或 O 并验证当前选择是否符合上述规则。如果不符合规则，则尝试另一种选择。如果两种选择都不符合规则，则回溯到上一个单元格，重新进行选择。

 */