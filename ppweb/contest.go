package main

import (
	"time"

	"strconv"

	"io/ioutil"

	"github.com/cs3238-tsuzu/popcon-sc/types"
	"github.com/jinzhu/gorm"
)

type Contest struct {
	Cid             int64               `gorm:"primary_key"`
	Name            string              `gorm:"not null;unique_index"`
	StartTime       time.Time           `gorm:"not null;default:CURRENT_TIMESTAMP;index"`
	FinishTime      time.Time           `gorm:"not null;default:CURRENT_TIMESTAMP;index"`
	Admin           int64               `gorm:"not null"`
	Type            sctypes.ContestType `gorm:"not null"`
	DescriptionFile string              `gorm:"not null"`
}

func (c *Contest) ProblemAdd(pidx int64, name string, time, mem int64, jtype JudgeType) (*ContestProblem, error) {
	pb, err := mainDB.ContestProblemAdd(c.Cid, pidx, name, time, mem, jtype)

	if err != nil {
		return nil, err
	}

	return mainDB.ContestProblemFind(pb)
}

func (c *Contest) DescriptionUpdate(desc string) error {
	var res Contest
	tx := mainDB.db.Begin()

	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Select("description_file").First(&res, c.Cid).Error; err != nil {
		tx.Rollback()
		return err
	}

	newName, err := mainFS.FileUpdate(FS_CATEGORY_CONTEST_DESCRIPTION, res.DescriptionFile, desc)

	if err != nil {
		tx.Rollback()

		return err
	}

	res.Cid = c.Cid
	if err := tx.Model(&res).Update("description_file", newName).Error; err != nil {
		tx.Rollback()

		return err
	}

	tx.Commit()
	return nil
}

func (c *Contest) DescriptionLoad() (string, error) {
	var res Contest
	res.Cid = c.Cid

	if err := mainDB.db.Select("description_file").First(&res, c.Cid).Error; err != nil {
		return "", err
	}

	b, err := mainFS.Read(FS_CATEGORY_CONTEST_DESCRIPTION, res.DescriptionFile)

	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (dm *DatabaseManager) CreateContestTable() error {
	err := dm.db.AutoMigrate(&Contest{}).Error

	if err != nil {
		return err
	}

	return nil
}

func (dm *DatabaseManager) ContestAdd(name string, start time.Time, finish time.Time, admin int64, ctype sctypes.ContestType) (int64, error) {
	contest := Contest{
		Name:       name,
		StartTime:  start,
		FinishTime: finish,
		Admin:      admin,
		Type:       ctype,
	}

	err := dm.db.Create(&contest).Error

	if err != nil {
		return 0, err
	}

	fs, err := mainFS.Open(FS_CATEGORY_CONTEST_DESCRIPTION, "contest_description_"+strconv.FormatInt(contest.Cid, 10)+".txt")

	if err != nil {
		return 0, err
	}
	if err := fs.Close(); err != nil {
		return 0, err
	}

	return contest.Cid, nil
}

func (dm *DatabaseManager) ContestUpdate(cid int64, name string, start time.Time, finish time.Time, admin int64, ctype sctypes.ContestType) error {
	cont := Contest{
		Cid:        cid,
		Name:       name,
		StartTime:  start,
		FinishTime: finish,
		Admin:      admin,
		Type:       ctype,
	}

	err := dm.db.Save(&cont).Error

	if err != nil {
		return err
	}

	return nil
}

func (dm *DatabaseManager) ContestDelete(cid int64) error {
	if cid == 0 {
		return ErrUnknownContest
	}

	err := dm.db.Delete(&Contest{Cid: cid}).Error

	if err != nil {
		return err
	}

	return mainFS.Remove(FS_CATEGORY_CONTEST_DESCRIPTION, "contest_description_"+strconv.FormatInt(cid, 10)+".txt")
}

func (dm *DatabaseManager) ContestFind(cid int64) (*Contest, error) {
	var res Contest

	err := dm.db.First(&res, cid).Error

	if err == gorm.ErrRecordNotFound {
		return nil, ErrUnknownContest
	}

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (dm *DatabaseManager) ContestDescriptionUpdateDisabled(cid int64, desc string) error {
	// TODO:[completed]Support GridFS

	fs, err := mainFS.OpenOnly(FS_CATEGORY_CONTEST_DESCRIPTION, "contest_description_"+strconv.FormatInt(cid, 10)+".txt")

	if err != nil {
		return err
	}

	id := fs.Id()

	if err := fs.Close(); err != nil {
		return err
	}

	fs, err = mainFS.Open(FS_CATEGORY_CONTEST_DESCRIPTION, "contest_description_"+strconv.FormatInt(cid, 10)+".txt")

	if err != nil {
		return err
	}

	_, err = fs.Write([]byte(desc))

	if err != nil {
		return err
	}
	if err := fs.Close(); err != nil {
		return err
	}

	return mainFS.RemoveID(FS_CATEGORY_CONTEST_DESCRIPTION, id)
}

func (dm *DatabaseManager) ContestDescriptionLoadDisabled(cid int64) (string, error) {
	//TODO:[completed]Support GridFS

	fs, err := mainFS.OpenOnly(FS_CATEGORY_CONTEST_DESCRIPTION, "contest_description_"+strconv.FormatInt(cid, 10)+".txt")

	if err != nil {
		return "", err
	}

	defer fs.Close()

	b, err := ioutil.ReadAll(fs)

	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (dm *DatabaseManager) ContestCount(options ...[]interface{}) (int64, error) {
	var count int64

	db := dm.db.Model(&Contest{})
	for i := range options {
		if len(options[i]) > 0 {
			db = db.Where(options[i][0], options[i][1:]...)
		}
	}

	err := db.Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}

// ContestList : if "offset" and "limit" aren't neccesary, set -1
func (dm *DatabaseManager) ContestList(offset, limit int, options ...[]interface{}) (*[]Contest, error) {
	var results []Contest

	db := dm.db
	for i := range options {
		if len(options[i]) > 0 {
			db = db.Where(options[i][0], options[i][1:]...)
		}
	}

	err := db.Offset(offset).Limit(limit).Order("start_time asc").Find(&results).Error

	if err != nil {
		return nil, err
	}

	return &results, nil
}