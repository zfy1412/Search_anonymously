package compute

import (
	"database/sql"
	"fmt"
	"goweb/send"
	"math"
	"sort"
	"time"
)

type point struct {
	x1     [100]float64
	x2     [100]float64
	year   [100]int
	month  [100]int
	day    [100]int
	hour   [100]int
	minute [100]int
}

var m int         //患者点的个数
var patient point //患者的移动点
var n int         //输入点的个数
var p point       //输入的点
var extrap point  //用函数式求的点
//经纬度距离计算
func EarthDistance(lat1, lng1, lat2, lng2 float64) float64 {
	radius := 6378.137
	rad := math.Pi / 180.0
	lat1 = lat1 * rad
	lng1 = lng1 * rad
	lat2 = lat2 * rad
	lng2 = lng2 * rad
	theta := lng2 - lng1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))
	return dist * radius
}

//e取0.4，点距离比较
func distance(n int, p point) int {

	var m int
	var (
		c, e float64
	)
	var f int
	f = 0
	m = 3
	e = 3
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {

			c = EarthDistance(p.x1[i], p.x2[i], patient.x1[j], patient.x2[j])
			//fmt.Printf("%f\n",c)
			if c <= e {
				f = f + 1
			}
			//fmt.Printf("函数点")

		}
		if i != n-1 {
			hanshu(p.x1[i], p.x1[i+1], p.x2[i], p.x2[i+1])

			for k := 0; k < 5; k++ {
				for j := 0; j < m; j++ {

					c = EarthDistance(extrap.x1[k], extrap.x2[k], patient.x1[j], patient.x2[j])
					//fmt.Printf("%f\n",c)
					if c <= e {
						f = f + 1
					}
				}
			}
		}
	}

	//fmt.Printf("计数点")

	return f
}

//e取0.5面积比较
func area(n int, p point) int {
	var m int
	var (
		s, e float64
	)
	var f int
	f = 0
	m = 80
	e = 0.000005
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			s = (patient.x1[j] - p.x1[i]) * (patient.x2[j] - p.x2[i]) * 0.5
			if s < 0 {
				s = 0 - s
			}
			if s <= e {
				f = f + 1
			}
		}
	}
	return f
}

//函数取点
func hanshu(x1, x2, y1, y2 float64) {
	var (
		k, b, r float64
	)
	k = (y1 - y2) / (x1 - x2)
	b = y1 - k*x1
	r = (x2 - x1) / 6
	for i := 0; i < 5; i++ {
		x1 = x1 + r
		extrap.x1[i] = x1
		extrap.x2[i] = k*extrap.x1[i] + b
	}

}

//函数取时间点
func timehanshu(a1, a2, b1, b2 int) {
	var ho, mi int
	ho = a2 - a1
	mi = b2 - b1
	if ho > 0 {
		mi = mi + 60*ho
	}
	mi = mi / 6
	for i := 0; i < 5; i++ {
		b1 = b1 + mi
		if b1 >= 60 {
			b1 = b1 - 60
			a1++
		}
		extrap.hour[i] = a1
		extrap.minute[i] = b1
	}
}

//时间预判断
func timejudgef(n int, p point) int {
	var judge int
	judge = 0
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			if p.year[i] != patient.year[j] || p.month[i] != patient.month[j] || p.day[i] != patient.day[j] {
				judge = 1
			}
		}
	}
	if judge == 0 {
		return 0
	} else {
		return -1
	}

}

////时间判断
func timejudge(n int, p point) int {

	var hou, min, f int
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			hou = p.hour[i] - patient.hour[j]
			min = (p.minute[i] - patient.minute[j]) + hou*60
			if min < 0 {
				min = 0 - min
			}
			if min < 20 {
				f++
			}
		}
		if i == 0 {
			timehanshu(p.hour[i], p.hour[i+1], p.minute[i], p.minute[i+1])
			for k := 0; k < 5; k++ {
				for j := 0; j < m; j++ {
					hou = extrap.hour[j] - patient.hour[j]
					min = (p.minute[i] - patient.minute[j]) + hou*60
					if min < 0 {
						min = 0 - min
					}
					if min < 20 {
						f++
					}
				}
			}
		}

	}
	//	fmt.Printf("%d\n",f)
	return f
}

//判断危险等级
func judge(fen int) string {
	if fen >= 0 && fen <= 5 {
		return "C"
	} else if fen > 6 && fen <= 20 {
		return "B"
	} else {
		return "A"
	}

}
func parse_timestr_to_datetime(time_str string) time.Time {

	t, error4 := time.Parse("2006-01-02 15:04:05", time_str)
	if error4 != nil {
		panic(error4)
	}
	return t
}

func Count(now send.Date) []string {
	m = 6

	//设定患者轨迹的基本信息

	fmt.Println("请输入您当日所到的地点数")
	n = len(now.Node)

	fmt.Println("请依次输入您当日所到的地点的起止时间")
	st := parse_timestr_to_datetime(now.Stime)
	ed := parse_timestr_to_datetime(now.Etime)
	fmt.Println(int(st.Month()))
	fmt.Println(ed.Year())
	p.year[0] = st.Year()
	p.month[0] = int(st.Month())
	p.day[0] = st.Day()
	p.hour[0] = st.Hour()
	p.minute[0] = st.Minute()
	p.year[1] = ed.Year()
	p.month[1] = int(ed.Month())
	p.day[1] = ed.Day()
	p.hour[1] = ed.Hour()
	p.minute[1] = ed.Minute()

	fmt.Println("请依次输入您当日所到的地点")
	for i := 0; i < n; i++ {
		p.x1[i] = now.Node[i].X
		p.x2[i] = now.Node[i].Y
	}
	cnt := 1
	var dict = make(map[string]int, 10)
	for cnt <= 4 {
		var db *sql.DB
		var err error
		db, err = sql.Open("mysql", "root:1592933843zzz@tcp(127.0.0.1:3306)/sql_find")

		if err != nil {
			fmt.Printf("connect to db 127.0.0.1:3306 error: %v\n", err)
		}
		sqlStr := "select id, x, y from trajectory where tid = ?"
		rows, err := db.Query(sqlStr, cnt)
		i := 0
		defer rows.Close()
		for rows.Next() {
			var u send.Trajectory
			err := rows.Scan(&u.Id, &u.X, &u.Y)
			fmt.Println(u.Tid)
			if err != nil {
				fmt.Printf("scan failed, err:%v\n", err)
			}
			patient.x1[i] = u.X
			patient.x2[i] = u.Y
			i++
		}

		var name string
		sqlStr = "select did, starttime, endtime from link where tid = ?"
		rows, err = db.Query(sqlStr, cnt)
		defer rows.Close()
		for rows.Next() {
			var u send.Link
			err := rows.Scan(&u.Uid, &u.Starttime, &u.Endtime)
			if err != nil {
				fmt.Printf("scan failed, err:%v\n", err)
			}
			st := parse_timestr_to_datetime(u.Starttime)
			ed := parse_timestr_to_datetime(u.Endtime)
			patient.year[0] = st.Year()
			patient.month[0] = int(st.Month())
			patient.day[0] = st.Day()
			patient.hour[0] = st.Hour()
			patient.minute[0] = st.Minute()
			patient.year[1] = ed.Year()
			patient.month[1] = int(ed.Month())
			patient.day[1] = ed.Day()
			patient.hour[1] = ed.Hour()
			patient.minute[1] = ed.Minute()
			name = u.Uid
		}
		var fen int
		fen = 0
		fen = timejudgef(n, p)
		if fen >= 0 {
			fen = distance(n, p) + (timejudge(n, p) / 2) + area(n, p)/100
			//fen = distance(n, p)
			//return judge(fen)
			fmt.Println("fen")
			fmt.Println(fen)
			if fen != 0 {
				dict[name] = fen
			}

		}
		cnt++
	}
	type peroson struct {
		Name string
		Fun  int
	}

	var lstPerson []peroson
	for k, v := range dict {
		lstPerson = append(lstPerson, peroson{k, v})
	}

	sort.Slice(lstPerson, func(i, j int) bool {
		return lstPerson[i].Fun > lstPerson[j].Fun // 升序
	})
	var ans []string
	if len(lstPerson) >= 3 {
		if lstPerson[0].Fun-lstPerson[1].Fun > 40 {
			ans = append(ans, lstPerson[0].Name)
		} else if lstPerson[1].Fun-lstPerson[2].Fun > 40 {
			ans = append(ans, lstPerson[0].Name, lstPerson[1].Name)
		} else {
			ans = append(ans, lstPerson[0].Name, lstPerson[1].Name, lstPerson[2].Name)
		}
	} else if len(lstPerson) == 2 {
		if lstPerson[0].Fun-lstPerson[1].Fun > 40 {
			ans = append(ans, lstPerson[0].Name)
		} else {
			ans = append(ans, lstPerson[0].Name, lstPerson[1].Name)
		}
	} else {
		ans = append(ans, lstPerson[0].Name)
	}
	return ans

	//fmt.Printf("%d",fen)
}
