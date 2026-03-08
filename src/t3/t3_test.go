package t3

import "testing"

// Test3_0 演示数组的多种定义方式
func Test3_0(t *testing.T) {

	var a = [3]int{}
	t.Log(a[0], a[1])

	var b = [3]int{1, 2, 3}
	t.Log(b[1])

	c := [3]int{1, 2, 3}
	t.Log(c[2])

	d := [...]int{0, 1, 2, 3, 4}
	t.Log(len(d), d[3])

	e := [2][2]int{{1, 1}, {2, 2}}
	t.Log(e[1])
}

// Test3_1 演示数组的遍历
func Test3_1(t *testing.T) {
	a := [...]int{0, 1, 2, 3, 4}

	/* 类似 java 的 for i */
	for i := 0; i < len(a); i++ {
		t.Log(a[i])
	}
	t.Log("===================")
	/* 类似 java 的 for each */
	for idx, e := range a {
		t.Log(idx, e)
	}
	t.Log("===================")
	/* 类似 java 的 for each，并且不关心第一个返回值 */
	for _, e := range a {
		t.Log(e)
	}
}

// Test3_2 演示数组截取（数组截取后得到 slice，而非数组）
func Test3_2(t *testing.T) {
	a := [...]int{0, 1, 2, 3, 4, 5}
	// 取前三个元素 [0 1 2]
	a0 := a[:3]
	t.Log(a0)
	a1 := a[0:3]
	t.Log(a1)
	// 取第三个之后的元素 [3 4 5]
	a2 := a[3:]
	t.Log(a2)
}

// modifyArray 接收数组值拷贝，修改不会影响调用方
func modifyArray(arr [3]int) {
	arr[0] = 999
}

// modifyArrayByPtr 接收数组指针，修改会影响调用方
func modifyArrayByPtr(arr *[3]int) {
	arr[0] = 999
}

// modifySlice 接收 slice（引用类型），修改会影响调用方底层数组
func modifySlice(sl []int) {
	sl[0] = 999
}

// TestArrayIsValueType 演示数组是值类型（传参时拷贝），对比 slice 是引用类型：
//   - 函数接收数组值时，修改副本不影响原数组
//   - 通过指针传递数组时，修改可影响原数组
//   - slice 是引用类型，传参后修改元素会影响原 slice 的底层数组
func TestArrayIsValueType(t *testing.T) {
	// ── 数组是值类型：传参时整体拷贝 ─────────────────────────────────
	original := [3]int{1, 2, 3}
	modifyArray(original)
	t.Logf("值传递后原数组: %v（未被修改）", original)
	if original[0] != 1 {
		t.Errorf("值传递不应修改原数组，期望 original[0]=1，得到 %d", original[0])
	}

	// ── 指针传递数组：修改影响原数组 ─────────────────────────────────
	original2 := [3]int{1, 2, 3}
	modifyArrayByPtr(&original2)
	t.Logf("指针传递后原数组: %v（已被修改）", original2)
	if original2[0] != 999 {
		t.Errorf("指针传递应修改原数组，期望 original2[0]=999，得到 %d", original2[0])
	}

	// ── 对比：slice 是引用类型 ────────────────────────────────────────
	// slice 本质是（指针, len, cap）三元组；传参后修改元素，
	// 因为底层数组是共享的，所以原 slice 可见该修改
	sl := []int{1, 2, 3}
	modifySlice(sl)
	t.Logf("slice 传参后: %v（底层数组已被修改）", sl)
	if sl[0] != 999 {
		t.Errorf("slice 传参后应能看到修改，期望 sl[0]=999，得到 %d", sl[0])
	}

	// ── 数组赋值也是拷贝 ─────────────────────────────────────────────
	arr1 := [3]int{10, 20, 30}
	arr2 := arr1 // 整体拷贝，arr2 是独立副本
	arr2[0] = 0
	t.Logf("arr1=%v, arr2=%v（赋值后各自独立）", arr1, arr2)
	if arr1[0] != 10 {
		t.Error("数组赋值应创建独立副本，arr1 不应被修改")
	}
}

// TestMultiDimArray 演示二维数组：
//   - [3][4]int 定义（3 行 4 列）
//   - 字面量初始化
//   - 两层 for 遍历，并验证元素值
func TestMultiDimArray(t *testing.T) {
	// ── 定义并初始化二维数组 ───────────────────────────────────────────
	// [3][4]int 表示：外层 3 个元素，每个元素是 [4]int
	var matrix [3][4]int
	// 填充：matrix[i][j] = i*10 + j
	for i := 0; i < 3; i++ {
		for j := 0; j < 4; j++ {
			matrix[i][j] = i*10 + j
		}
	}
	t.Logf("matrix: %v", matrix)

	// ── 字面量初始化 ───────────────────────────────────────────────────
	grid := [3][4]int{
		{0, 1, 2, 3},
		{10, 11, 12, 13},
		{20, 21, 22, 23},
	}
	t.Logf("grid: %v", grid)

	// ── 两层 for 遍历并验证 ────────────────────────────────────────────
	for i := 0; i < 3; i++ {
		for j := 0; j < 4; j++ {
			expected := i*10 + j
			if grid[i][j] != expected {
				t.Errorf("grid[%d][%d] 期望 %d，得到 %d", i, j, expected, grid[i][j])
			}
		}
	}

	// ── for-range 遍历二维数组 ─────────────────────────────────────────
	t.Log("for-range 遍历二维数组:")
	for i, row := range grid {
		for j, val := range row {
			t.Logf("  grid[%d][%d] = %d", i, j, val)
		}
	}

	// ── 验证二维数组也是值类型 ─────────────────────────────────────────
	copy2D := grid        // 拷贝整个二维数组
	copy2D[0][0] = 99999 // 修改副本不影响原数组
	if grid[0][0] != 0 {
		t.Errorf("二维数组赋值应创建独立副本，grid[0][0] 期望 0，得到 %d", grid[0][0])
	}
	t.Logf("修改副本后: grid[0][0]=%d（不变），copy2D[0][0]=%d（已改）", grid[0][0], copy2D[0][0])
}
