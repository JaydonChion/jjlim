package main

import "time"

type history struct {
	CreatedAt   time.Time
	Location_id int
}

type food struct {
	Location  string
	FoodType  string
	Amount    int
	Descrp    string
	CreatedAt time.Time
}

type cData struct {
	Tpeople  int
	Nmale    int
	Nfemale  int
	Age      int
	Time     int
	Foodleft int
}

func contributeInfo(tpeopleNum int, nmaleNum int, nfemaleNum int, ageNum int, timeNum int, foodleftNum int) (err error) {
	statement := "insert into predictData(tpeople,nmale,nfemale,age,timeNum,foodleft) values ($1,$2,$3,$4,$5,$6);"
	stmt, err := Db.Prepare(statement)
	_, err = stmt.Exec(tpeopleNum, nmaleNum, nfemaleNum, ageNum, timeNum, foodleftNum)
	defer stmt.Close()

	return
}

func getContribInfo() (cDatas []cData, err error) {
	rows, err := Db.Query("SELECT tpeople,nmale,nfemale,age,timeNum,foodleft FROM predictData limit 100;")
	if err != nil {
		return
	}
	for rows.Next() {
		cd := cData{}
		if err = rows.Scan(&cd.Tpeople, &cd.Nmale, &cd.Nfemale, &cd.Age, &cd.Time, &cd.Foodleft); err != nil {
			return
		}
		cDatas = append(cDatas, cd)
	}
	rows.Close()

	return

}

func getDonatedFood() (foods []food, err error) {
	rows, err := Db.Query("SELECT location,halal,amountNum,descrp,created_at FROM donorData order by created_at desc limit 100;")
	if err != nil {
		return
	}
	for rows.Next() {
		food := food{}
		if err = rows.Scan(&food.Location, &food.FoodType, &food.Amount, &food.Descrp, &food.CreatedAt); err != nil {
			return
		}
		foods = append(foods, food)
	}
	rows.Close()

	return
}

func donateFood(location string, name string, phoneNum int, Mnum string, halal string, amountNum int, descrp string, userId int) (err error) {
	statement := "insert into donorData(location,name,phoneNum,Mnum,halal,amountNum,descrp,user_id,created_at) values ($1,$2,$3,$4,$5,$6,$7,$8,now());"
	stmt, err := Db.Prepare(statement)
	_, err = stmt.Exec(location, name, phoneNum, Mnum, halal, amountNum, descrp, userId)
	defer stmt.Close()

	return
}

func getCupSaved(id int) (cups int, err error) {

	err = Db.QueryRow("SELECT cupSaved FROM users WHERE id = $1;", id).
		Scan(&cups)

	return

}

func activate(mid int) (check bool, err error) {

	err = Db.QueryRow("SELECT EXTRACT(EPOCH FROM (now() - created_at))<2 from rentaldata where location_id = $1 order by created_at desc limit 1;", mid).
		Scan(&check)

	return

}

func getHistory(userId int) (hists []history, err error) {

	rows, err := Db.Query("SELECT created_at,location_id FROM rentaldata where user_id = $1 order by created_at desc", userId)
	if err != nil {
		return
	}
	for rows.Next() {
		hist := history{}
		if err = rows.Scan(&hist.CreatedAt, &hist.Location_id); err != nil {
			return
		}
		hists = append(hists, hist)
	}
	rows.Close()

	return
}

func rentService(userIDnum int, itemID int, locationID int) (err error) {

	statement := "insert into rentalData(user_id, created_at, item_id, location_id) values ($1,now(),$2,$3);"
	stmt, err := Db.Prepare(statement)
	_, err = stmt.Exec(userIDnum, itemID, locationID)
	defer stmt.Close()

	cups, _ := getCupSaved(userIDnum)
	statement2 := "update users set cupSaved = $1 where id = $2;"
	stmt2, _ := Db.Prepare(statement2)
	_, err = stmt2.Exec(cups+1, userIDnum)

	if err != nil {
		panic(err)
	}

	return

}
