package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/cs3238-tsuzu/popcon-sc/ppweb/file_manager"
)

var ContestProblemDir = "contest_problems/"

type JudgeType int

const (
	JudgePerfectMatch JudgeType = 0
	JudgeRunningCode  JudgeType = 1
)

type ContestProblemTestCase struct {
	Pid    int64
	Name   string
	Input  string
	Output string
}

type ContestProblemScoreSetCasesString string

func (cs *ContestProblemScoreSetCasesString) Get() []int64 {
	var results []int64
	json.Unmarshal([]byte(*cs), &results)

	return results
}

func (cs *ContestProblemScoreSetCasesString) Set(a []int64) {
	if b, err := json.Marshal(a); err != nil {
		*cs = "[]"
	} else {
		*cs = ContestProblemScoreSetCasesString(string(b))
	}
}

type ContestProblemScoreSet struct {
	Pid   int64
	Cases ContestProblemScoreSetCasesString
	Score int
}

type ContestProblem struct {
	Pid          int64                    `gorm:"primary_key"`
	Cid          int64                    `gorm:"not null;index;unique_index:cid_and_pidx_index"`
	Pidx         int64                    `gorm:"not null;index;unique_index:cid_and_pidx_index"`
	Name         string                   `gorm:"not null;size:255"`
	Time         int64                    `gorm:"not null"` // Second
	Mem          int64                    `gorm:"not null"` // MB
	LastModified int64                    `gorm:"not null"`
	Score        int                      `gorm:"not null"`
	Type         JudgeType                `gorm:"not null"`       // int->JudgeType
	Cases        []ContestProblemTestCase `gorm:"ForeignKey:Pid"` // json format []TestCase
	Scores       []ContestProblemScoreSet `gorm:"ForeignKey:Pid"` // json format []ScoreSet
}

// TODO: テストケースの情報を乗っけるようにする途中でORMの変更が入ったので中断
// 適当にSQLに乗っけるように変更

func (cp *ContestProblem) UpdateStatement(text string) error {
	return mainFS.Write(FS_CATEGORY_PROBLEM_STATEMENT, strconv.FormatInt(cp.Pid, 10)+"_prob.txt", []byte(text))
}

func (cp *ContestProblem) LoadStatement() (string, error) {
	b, err := mainFS.Read(FS_CATEGORY_PROBLEM_STATEMENT, strconv.FormatInt(cp.Pid, 10)+"_prob.txt")

	if err != nil {
		return "", err
	}

	return string(b), nil
}

type CheckerSavedFormat struct {
	Lid  int64
	Code string
}

func (cp *ContestProblem) UpdateChecker(lid int64, code string) error {
	b, err := json.Marshal(CheckerSavedFormat{lid, code})

	if err != nil {
		return err
	}

	return mainFS.Write(FS_CATEGORY_PROBLEM_STATEMENT, strconv.FormatInt(cp.Pid, 10)+"_checker.txt", b)
}

func (cp *ContestProblem) LoadChecker() (int64, string, error) {
	b, err := mainFS.Read(FS_CATEGORY_PROBLEM_STATEMENT, strconv.FormatInt(cp.Pid, 10)+"_checker.txt")

	if err != nil {
		return 0, "", err
	}

	if len(b) == 0 {
		return 0, "", nil
	}

	var ci CheckerSavedFormat
	err = json.Unmarshal(b, &ci)

	if err != nil {
		return 0, "", err
	}

	return ci.Lid, ci.Code, nil
}

type TestCaseJson struct {
	CaseNames []ContestProblemTestCase `json:"case_names"`
	Scores    []ContestProblemScoreSet `json:"scores"`
}

func (cp *ContestProblem) UpdateTestCaseNames(cases []string, scores []ContestProblemScoreSet) (resErr error) {
	scoreSum := 0
	for i := range scores {
		scoreSum += scores[i].Score
	}

	/*tx, err := mainDB.db.DB().Begin()

	if err != nil {
		return err
	}*/

	/*var casesString, scoresString string
	defer func() {
		if err := recover(); err != nil {
			resErr = err.(error)
			tx.Rollback()
		} else {
			tx.Commit()
			cp.Score = scoreSum
			cp.Scores = scoresString
			cp.Cases = casesString
		}
	}()

	rows, err := tx.Query("select cases, scores from contest_problem where pid=?", cp.Pid)

	if err != nil {
		return err
	}

	rows.Next()
	var oldCases []TestCase
	var oldScores []ScoreSet
	err = rows.Scan(&oldCases, &oldScores)
	rows.Close()

	if err != nil {
		panic(err)
	}

	newCases := make([]TestCase, len(cases))

	for i := range cases {
		newCases[i].Name = cases[i]

		if i < len(oldCases) {
			newCases[i].Input = oldCases[i].Input
			newCases[i].Output = oldCases[i].Output
		}
	}

	for i := len(cases); i < len(oldCases); i++ {
		// TODO: Remove files
		/*
			oldCases[i].Input
	*/
	/*}

	casesBytes, _ := json.Marshal(newCases)
	casesString = string(casesBytes)
	scoresBytes, _ := json.Marshal(scores)
	scoresString = string(scoresBytes)

	_, err = tx.Exec("update contest_problem set cases=?, scores=? where pid=?", string(casesBytes), string(scoresBytes), cp.Pid)

	if err != nil {
		panic(err)
	}
	*/
	return nil
}

func (cp *ContestProblem) CreateUniquelyNamedFile() (*mgo.GridFile, error) {
	id, err := mainRM.UniqueFileID(FS_CATEGORY_TESTCASE_INOUT)

	if err != nil {
		return nil, err
	}

	return mainFS.Open(FS_CATEGORY_TESTCASE_INOUT, "testcase_"+mainFS.TestcaseFileBaseTag+"_"+strconv.FormatInt(id, 10))
}

// ErrUnknownTestcase
func (cp *ContestProblem) UpdateTestCase(isInput bool, caseID int, str string) (retErr error) {
	/*fp, err := cp.CreateUniquelyNamedFile()

	if err != nil {
		return err
	}

	fileName := fp.Name()

	_, err = fp.Write([]byte(str))

	if err != nil {
		return err
	}

	err = fp.Close()

	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			mainFS.db.GridFS(FS_CATEGORY_TESTCASE_INOUT).Remove(fileName)
		}
	}()

	tx, err := mainDB.db.DB().Begin()

	if err != nil {
		return err
	}

	defer func() {
		if err := recover(); err != nil {
			retErr = err.(error)
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	rows, err := tx.Query("select cases from contest_problem where pid=?", cp.Pid)

	if err != nil {
		panic(err)
	}

	rows.Next()
	var results []string
	err = rows.Scan(&str)
	rows.Close()

	if err != nil {
		panic(err)
	}
	if len(results) == 0 {
		panic(ErrUnknownProblem)
	}

	result := results[0]

	var cases []ContestProblemTestCase
	json.Unmarshal([]byte(result), &cases)

	if isInput {
		cases[caseID].Input = fileName
	} else {
		cases[caseID].Output = fileName
	}

	casesBytes, _ := json.Marshal(cases)
	casesString := string(casesBytes)

	_, err = tx.Exec("update contest_problem set cases=? where pid=?", casesString, cp.Pid)

	if err != nil {
		panic(err)
	}*/

	return nil
}

func (cp *ContestProblem) LoadTestCase(isInput bool, caseID int) (string, error) {
	/*fm, err := FileManager.OpenFile(filepath.Join(ContestProblemDir, strconv.FormatInt(cp.Pid, 10)+"/.cases_lock"), os.O_RDONLY, false)

	if err != nil {
		return "", err
	}

	defer fm.Close()

	fileTag := "_in"
	if !isInput {
		fileTag = "_out"
	}

	fp, err := os.OpenFile(filepath.Join(ContestProblemDir, strconv.FormatInt(cp.Pid, 10)+"/cases/"+strconv.FormatInt(int64(caseID), 10)+fileTag), os.O_RDONLY, 0644)

	if err != nil {
		return "", err
	}

	defer fp.Close()

	b, err := ioutil.ReadAll(fp)

	if err != nil {
		return "", nil
	}

	return string(b), err*/
	return "", nil
}

func (cp *ContestProblem) LoadTestCases() ([]ContestProblemTestCase, []ContestProblemScoreSet, error) {
	var scores []ContestProblemScoreSet
	var cases []ContestProblemTestCase

	/*rows, err := mainDB.db.DB().Query("select cases, scores from contest_problem where pid=?", cp.Pid)

	if err != nil {
		return nil, nil, err
	}

	rows.Next()
	err = rows.Scan(&cases, &scores)
	rows.Close()

	if err != nil {
		return nil, nil, ErrUnknownProblem
	}*/

	return cases, scores, nil
}

func (cp *ContestProblem) LoadTestCaseInfo(caseId int) (int64, int64, error) {
	var casesString string
	//var cases []TestCase
	rows, err := mainDB.db.DB().Query("select cases from contest_problem where pid=?", cp.Pid)

	if err != nil {
		return 0, 0, err
	}

	if !rows.Next() {
		return 0, 0, ErrUnknownTestcase
	}
	err = rows.Scan(&casesString)

	return 0, 0, err
}

func (cp *ContestProblem) LoadTestCaseNames() ([]string, []ContestProblemScoreSet, error) {
	var scores []ContestProblemScoreSet
	var cases []string

	fm, err := FileManager.OpenFile(filepath.Join(ContestProblemDir, strconv.FormatInt(cp.Pid, 10)+"/.cases_lock"), os.O_RDONLY, false)

	if err != nil {
		return nil, nil, err
	}

	defer fm.Close()

	fp, err := os.OpenFile(filepath.Join(ContestProblemDir, strconv.FormatInt(cp.Pid, 10)+"/cases/data"), os.O_RDONLY, 0644)

	if err != nil {
		return cases, scores, nil
	}

	b, err := ioutil.ReadAll(fp)

	fp.Close()

	if err != nil {
		return nil, nil, err
	}

	var tcj TestCaseJson

	err = json.Unmarshal(b, &tcj)

	if err != nil {
		return cases, scores, nil
	}

	scores = tcj.Scores

	cases = make([]string, len(tcj.CaseNames))

	// for x := range tcj.CaseNames {
	// 	i, err := strconv.ParseInt(x, 10, 32)

	// 	if err != nil {
	// 		return nil, nil, err
	// 	}

	// 	cases[i] = tcj.CaseNames[x]
	// }
	// return &cases, &scores, nil

	return nil, nil, nil
}

func (dm *DatabaseManager) CreateContestProblemTable() error {
	err := dm.db.AutoMigrate(&ContestProblem{}, &ContestProblemTestCase{}, &ContestProblemScoreSet{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (dm *DatabaseManager) ContestProblemAdd(cid, pidx int64, name string, timeLimit, mem int64, jtype JudgeType) (int64, error) {
	res, err := dm.db.DB().Exec("insert into contest_problem (cid, pidx, name, time, mem, last_modified, score, type) values (?, ?, ?, ?, ?, ?, ?, ?)", cid, pidx, name, timeLimit, mem, time.Now().Unix(), 0, int(jtype))

	if err != nil {
		return 0, err
	}

	id, _ := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	err = os.MkdirAll(filepath.Join(ContestProblemDir, strconv.FormatInt(id, 10)+"/cases/"), os.ModePerm)

	if err != nil {
		dm.ContestProblemDelete(id)

		return 0, err
	}

	fm, err := FileManager.OpenFile(filepath.Join(ContestProblemDir, strconv.FormatInt(id, 10)+"/.cases_lock"), os.O_WRONLY|os.O_CREATE, true)

	if err != nil {
		dm.ContestProblemDelete(id)

		return 0, err
	}

	fm.Close()

	fm, err = FileManager.OpenFile(filepath.Join(ContestProblemDir, strconv.FormatInt(id, 10)+"/checker"), os.O_WRONLY|os.O_CREATE, true)

	if err != nil {
		dm.ContestProblemDelete(id)

		return 0, err
	}

	fm.Close()

	fm, err = FileManager.OpenFile(filepath.Join(ContestProblemDir, strconv.FormatInt(id, 10)+"/prob"), os.O_WRONLY|os.O_CREATE, true)

	if err != nil {
		dm.ContestProblemDelete(id)

		return 0, err
	}

	fm.Close()

	return id, err
}

func (dm *DatabaseManager) ContestProblemUpdate(prob ContestProblem) error {
	return dm.db.Update(&prob).Error
}

func (dm *DatabaseManager) ContestProblemDelete(pid int64) error {
	//timeLimit := ContestProblem{Pid: pid}

	//_, err := dm.db.Delete(&timeLimit)
	err := error(nil)
	if err != nil {
		return err
	}

	fm1, _ := FileManager.OpenFile(filepath.Join(ContestProblemDir, strconv.FormatInt(pid, 10)+"/prob.txt"), os.O_WRONLY|os.O_CREATE, true)
	fm2, _ := FileManager.OpenFile(filepath.Join(ContestProblemDir, strconv.FormatInt(pid, 10)+"/.cases_lock"), os.O_WRONLY|os.O_CREATE, true)
	fm3, _ := FileManager.OpenFile(filepath.Join(ContestProblemDir, strconv.FormatInt(pid, 10)+"/checker"), os.O_WRONLY|os.O_CREATE, true)

	defer func() {
		if fm1 != nil {
			fm1.Close()
		}
		if fm2 != nil {
			fm2.Close()
		}
		if fm3 != nil {
			fm3.Close()
		}
	}()

	err = os.RemoveAll(filepath.Join(ContestProblemDir, strconv.FormatInt(pid, 10)))

	return err
}

func (dm *DatabaseManager) ContestProblemFind(pid int64) (*ContestProblem, error) {
	var resulsts []ContestProblem

	err := dm.db.Select(&resulsts, dm.db.Where("pid", "=", pid)).Error

	if err != nil {
		return nil, err
	}

	if len(resulsts) == 0 {
		return nil, ErrUnknownProblem
	}

	return &resulsts[0], nil
}

func (dm *DatabaseManager) ContestProblemFind2(cid, pidx int64) (*ContestProblem, error) {
	var resulsts []ContestProblem

	//err := dm.db.Select(&resulsts, dm.db.Where("pidx", "=", pidx).And("cid", "=", cid)).Error
	err := error(nil)
	if err != nil {
		return nil, err
	}

	if len(resulsts) == 0 {
		return nil, ErrUnknownProblem
	}

	return &resulsts[0], nil
}

func (dm *DatabaseManager) ContestProblemList(cid int64) (*[]ContestProblem, error) {
	var results []ContestProblem

	//	err := dm.db.Select(&results, dm.db.Where("cid", "=", cid), dm.db.OrderBy("pidx", genmai.ASC))
	err := error(nil)

	if err != nil {
		return nil, err
	}

	return &results, nil
}

func (dm *DatabaseManager) ContestProblemCount(cid int64) (int64, error) {
	var count int64

	// COUNT(*)が重い
	//err := dm.db.Select(&count, dm.db.Count("pid"), dm.db.From(&ContestProblem{}), dm.db.Where("cid", "=", cid))
	err := error(nil)

	if err != nil {
		return 0, err
	}

	return count, nil
}

type ContestProblemLight struct {
	Pidx int64
	Name string
}

func (dm *DatabaseManager) ContestProblemListLight(cid int64) ([]ContestProblemLight, error) {
	results := make([]ContestProblemLight, 0, 50)

	rows, err := dm.db.DB().Query("select pidx, name from contest_problem where cid = ?", cid)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var cpl ContestProblemLight

		rows.Scan(&cpl.Pidx, &cpl.Name)

		results = append(results, cpl)
	}

	return results, nil
}

func (dm *DatabaseManager) ContestProblemRemoveAll(cid int64) error {
	_, err := dm.db.DB().Exec("delete from contest_problem where cid = ?", cid)

	return err
}