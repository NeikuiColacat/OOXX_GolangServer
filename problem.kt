/*
N 为题目矩阵大小
alpha为超参数，保证生成题目挖空cell数量不超过alpha,用于控制难度
成员变量a为一个二维可变数组 , 表示题目答案
成员变量b为一个二维可变数组，表示生成题目
1表示挖空 , 0 和 1用于表示 OO和XX ， 依照喜好自行定义
*/
import kotlin.random.Random
class problem (var N : Int , val alpha : Int){
    var a = MutableList(N){MutableList(N){-1} }
    var b = MutableList(N){MutableList(N){-1} }
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
    }
}